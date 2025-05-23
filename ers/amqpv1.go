/*
Real-time Online/Offline Charging System (OCS) for Telecom & ISP environments
Copyright (C) ITsysCOM GmbH

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>
*/

package ers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	amqpv1 "github.com/Azure/go-amqp"
	"github.com/cgrates/cgrates/agents"
	"github.com/cgrates/cgrates/config"
	"github.com/cgrates/cgrates/engine"
	"github.com/cgrates/cgrates/utils"
)

// NewAMQPv1ER return a new amqpv1 event reader
func NewAMQPv1ER(cfg *config.CGRConfig, cfgIdx int,
	rdrEvents, partialEvents chan *erEvent, rdrErr chan error,
	fltrS *engine.FilterS, rdrExit chan struct{}) (er EventReader, err error) {
	rdr := &AMQPv1ER{
		cgrCfg:        cfg,
		cfgIdx:        cfgIdx,
		fltrS:         fltrS,
		rdrEvents:     rdrEvents,
		partialEvents: partialEvents,
		rdrExit:       rdrExit,
		rdrErr:        rdrErr,
	}
	if concReq := rdr.Config().ConcurrentReqs; concReq != -1 {
		rdr.cap = make(chan struct{}, concReq)
	}

	if amqOpts := rdr.Config().Opts.AMQP; amqOpts != nil {
		if amqOpts.QueueID != nil {
			rdr.queueID = "/" + *amqOpts.QueueID
		}
		if amqOpts.Username != nil && amqOpts.Password != nil {
			rdr.connOpts = &amqpv1.ConnOptions{
				SASLType: amqpv1.SASLTypePlain(*amqOpts.Username, *amqOpts.Password),
			}
		}
	}
	return rdr, nil
}

// AMQPv1ER implements EventReader interface for amqpv1 message
type AMQPv1ER struct {
	// sync.RWMutex
	cgrCfg *config.CGRConfig
	cfgIdx int // index of config instance within ERsCfg.Readers
	fltrS  *engine.FilterS

	queueID string

	rdrEvents     chan *erEvent // channel to dispatch the events created to
	partialEvents chan *erEvent // channel to dispatch the partial events created to
	rdrExit       chan struct{}
	rdrErr        chan error
	cap           chan struct{}

	conn     *amqpv1.Conn
	connOpts *amqpv1.ConnOptions
	ses      *amqpv1.Session
}

// Config returns the curent configuration
func (rdr *AMQPv1ER) Config() *config.EventReaderCfg {
	return rdr.cgrCfg.ERsCfg().Readers[rdr.cfgIdx]
}

// Serve will start the gorutines needed to watch the amqpv1 topic
func (rdr *AMQPv1ER) Serve() (err error) {
	if rdr.conn, err = amqpv1.Dial(context.TODO(), rdr.Config().SourcePath, rdr.connOpts); err != nil {
		return
	}
	if rdr.ses, err = rdr.conn.NewSession(context.TODO(), nil); err != nil {
		rdr.close()
		return
	}
	if rdr.Config().RunDelay == time.Duration(0) { // 0 disables the automatic read, maybe done per API
		return
	}

	var receiver *amqpv1.Receiver
	if receiver, err = rdr.ses.NewReceiver(context.TODO(), rdr.queueID,
		nil); err != nil {
		return
	}
	go func() {
		<-rdr.rdrExit
		receiver.Close(context.Background())
		rdr.close()
	}()

	go rdr.readLoop(receiver) // read until the connection is closed
	return
}

func (rdr *AMQPv1ER) readLoop(recv *amqpv1.Receiver) (err error) {
	if rdr.Config().StartDelay > 0 {
		select {
		case <-time.After(rdr.Config().StartDelay):
		case <-rdr.rdrExit:
			return
		}
	}
	for {
		if rdr.Config().ConcurrentReqs != -1 {
			rdr.cap <- struct{}{} // do not try to read if the limit is reached
		}
		ctx := context.Background()
		var msg *amqpv1.Message
		if msg, err = recv.Receive(ctx, nil); err != nil {
			if err.Error() == "amqp: link closed" {
				err = nil
				return
			}
			rdr.rdrErr <- err
			return
		}
		if err = recv.AcceptMessage(ctx, msg); err != nil {
			utils.Logger.Warning(
				fmt.Sprintf("<%s> unable to accept message error: %s",
					utils.ERs, err.Error()))
			continue
		}

		go func(msg *amqpv1.Message) {
			body := msg.GetData()
			if err := rdr.processMessage(body); err != nil {
				utils.Logger.Warning(
					fmt.Sprintf("<%s> processing message error: %s",
						utils.ERs, err.Error()))
			}
			if rdr.Config().ConcurrentReqs != -1 {
				<-rdr.cap
			}
		}(msg)
	}
}

func (rdr *AMQPv1ER) processMessage(msg []byte) (err error) {
	var decodedMessage map[string]any
	if err = json.Unmarshal(msg, &decodedMessage); err != nil {
		return
	}

	reqVars := &utils.DataNode{Type: utils.NMMapType, Map: map[string]*utils.DataNode{utils.MetaReaderID: utils.NewLeafNode(rdr.cgrCfg.ERsCfg().Readers[rdr.cfgIdx].ID)}}

	agReq := agents.NewAgentRequest(
		utils.MapStorage(decodedMessage), reqVars,
		nil, nil, nil, rdr.Config().Tenant,
		rdr.cgrCfg.GeneralCfg().DefaultTenant,
		utils.FirstNonEmpty(rdr.Config().Timezone,
			rdr.cgrCfg.GeneralCfg().DefaultTimezone),
		rdr.fltrS, nil) // create an AgentRequest
	var pass bool
	if pass, err = rdr.fltrS.Pass(agReq.Tenant, rdr.Config().Filters,
		agReq); err != nil || !pass {
		return
	}
	if err = agReq.SetFields(rdr.Config().Fields); err != nil {
		return
	}
	cgrEv := utils.NMAsCGREvent(agReq.CGRRequest, agReq.Tenant, utils.NestingSep, agReq.Opts)
	rdrEv := rdr.rdrEvents
	if _, isPartial := cgrEv.APIOpts[utils.PartialOpt]; isPartial {
		rdrEv = rdr.partialEvents
	}
	rdrEv <- &erEvent{
		cgrEvent: cgrEv,
		rdrCfg:   rdr.Config(),
	}
	return
}

func (rdr *AMQPv1ER) close() (err error) {
	if rdr.ses != nil {
		if err = rdr.ses.Close(context.Background()); err != nil {
			return
		}
	}
	return rdr.conn.Close()
}
