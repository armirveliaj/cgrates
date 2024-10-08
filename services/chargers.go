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

package services

import (
	"fmt"
	"sync"

	"github.com/cgrates/birpc"
	v1 "github.com/cgrates/cgrates/apier/v1"
	"github.com/cgrates/cgrates/config"
	"github.com/cgrates/cgrates/cores"
	"github.com/cgrates/cgrates/engine"
	"github.com/cgrates/cgrates/servmanager"
	"github.com/cgrates/cgrates/utils"
)

// NewChargerService returns the Charger Service
func NewChargerService(cfg *config.CGRConfig, dm *DataDBService,
	cacheS *engine.CacheS, filterSChan chan *engine.FilterS, server *cores.Server,
	internalChargerSChan chan birpc.ClientConnector, connMgr *engine.ConnManager,
	anz *AnalyzerService, srvDep map[string]*sync.WaitGroup) servmanager.Service {
	return &ChargerService{
		connChan:    internalChargerSChan,
		cfg:         cfg,
		dm:          dm,
		cacheS:      cacheS,
		filterSChan: filterSChan,
		server:      server,
		connMgr:     connMgr,
		anz:         anz,
		srvDep:      srvDep,
	}
}

// ChargerService implements Service interface
type ChargerService struct {
	sync.RWMutex
	cfg         *config.CGRConfig
	dm          *DataDBService
	cacheS      *engine.CacheS
	filterSChan chan *engine.FilterS
	server      *cores.Server
	connMgr     *engine.ConnManager

	chrS     *engine.ChargerService
	connChan chan birpc.ClientConnector
	anz      *AnalyzerService
	srvDep   map[string]*sync.WaitGroup
}

// Start should handle the sercive start
func (chrS *ChargerService) Start() error {
	if chrS.IsRunning() {
		return utils.ErrServiceAlreadyRunning
	}

	<-chrS.cacheS.GetPrecacheChannel(utils.CacheChargerProfiles)
	<-chrS.cacheS.GetPrecacheChannel(utils.CacheChargerFilterIndexes)

	filterS := <-chrS.filterSChan
	chrS.filterSChan <- filterS
	dbchan := chrS.dm.GetDMChan()
	datadb := <-dbchan
	dbchan <- datadb

	chrS.Lock()
	defer chrS.Unlock()
	chrS.chrS = engine.NewChargerService(datadb, filterS, chrS.cfg, chrS.connMgr)
	utils.Logger.Info(fmt.Sprintf("<%s> starting <%s> subsystem", utils.CoreS, utils.ChargerS))
	srv, err := engine.NewService(v1.NewChargerSv1(chrS.chrS))
	if err != nil {
		return err
	}
	if !chrS.cfg.DispatcherSCfg().Enabled {
		chrS.server.RpcRegister(srv)
	}
	chrS.connChan <- chrS.anz.GetInternalCodec(srv, utils.ChargerS)
	return nil
}

// Reload handles the change of config
func (chrS *ChargerService) Reload() (err error) {
	return
}

// Shutdown stops the service
func (chrS *ChargerService) Shutdown() (err error) {
	chrS.Lock()
	defer chrS.Unlock()
	chrS.chrS.Shutdown()
	chrS.chrS = nil
	<-chrS.connChan
	return
}

// IsRunning returns if the service is running
func (chrS *ChargerService) IsRunning() bool {
	chrS.RLock()
	defer chrS.RUnlock()
	return chrS.chrS != nil
}

// ServiceName returns the service name
func (chrS *ChargerService) ServiceName() string {
	return utils.ChargerS
}

// ShouldRun returns if the service should be running
func (chrS *ChargerService) ShouldRun() bool {
	return chrS.cfg.ChargerSCfg().Enabled
}
