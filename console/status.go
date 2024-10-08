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

package console

import (
	"github.com/cgrates/cgrates/cores"
	"github.com/cgrates/cgrates/utils"
)

func init() {
	c := &CmdStatus{
		name:      "status",
		rpcMethod: utils.CoreSv1Status,
		rpcParams: &cores.V1StatusParams{},
	}
	commands[c.Name()] = c
	c.CommandExecuter = &CommandExecuter{c}
}

type CmdStatus struct {
	name      string
	rpcMethod string
	rpcParams *cores.V1StatusParams
	*CommandExecuter
}

func (self *CmdStatus) Name() string {
	return self.name
}

func (self *CmdStatus) RpcMethod() string {
	return self.rpcMethod
}

func (self *CmdStatus) RpcParams(reset bool) any {
	if reset || self.rpcParams == nil {
		self.rpcParams = &cores.V1StatusParams{
			APIOpts: make(map[string]any),
		}
	}
	return self.rpcParams
}

func (self *CmdStatus) PostprocessRpcParams() error {
	return nil
}

func (self *CmdStatus) RpcResult() any {
	var s map[string]any
	return &s
}
