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
	"sync"
	"testing"

	"github.com/cgrates/birpc"
	"github.com/cgrates/birpc/context"
	"github.com/cgrates/cgrates/config"
	"github.com/cgrates/cgrates/cores"
	"github.com/cgrates/cgrates/engine"
	"github.com/cgrates/cgrates/utils"
)

func TestNewTrendService(t *testing.T) {
	cfg := &config.CGRConfig{}
	dm := &DataDBService{}
	cacheS := &CacheService{}
	filterSChan := make(chan *engine.FilterS)
	server := &cores.Server{}
	internalTrendSChan := make(chan birpc.ClientConnector)
	connMgr := &engine.ConnManager{}
	anz := &AnalyzerService{}
	srvDep := make(map[string]*sync.WaitGroup)

	service := NewTrendService(cfg, dm, cacheS, filterSChan, server, internalTrendSChan, connMgr, anz, srvDep)

	trendService, ok := service.(*TrendService)
	if !ok {
		t.Errorf("Expected type *TrendService, but got %T", service)
	}

	if trendService.cfg != cfg {
		t.Errorf("Expected cfg to be %v, but got %v", cfg, trendService.cfg)
	}
	if trendService.dm != dm {
		t.Errorf("Expected dm to be %v, but got %v", dm, trendService.dm)
	}
	if trendService.cacheS != cacheS {
		t.Errorf("Expected cacheS to be %v, but got %v", cacheS, trendService.cacheS)
	}

	if trendService.server != server {
		t.Errorf("Expected server to be %v, but got %v", server, trendService.server)
	}
	if trendService.connChan != internalTrendSChan {
		t.Errorf("Expected connChan to be %v, but got %v", internalTrendSChan, trendService.connChan)
	}
	if trendService.connMgr != connMgr {
		t.Errorf("Expected connMgr to be %v, but got %v", connMgr, trendService.connMgr)
	}
	if trendService.anz != anz {
		t.Errorf("Expected anz to be %v, but got %v", anz, trendService.anz)
	}

}

func TestTrendServiceServiceName(t *testing.T) {
	tr := &TrendService{}

	serviceName := tr.ServiceName()

	expected := utils.TrendS
	if serviceName != expected {
		t.Errorf("Expected service name to be %s, but got %s", expected, serviceName)
	}
}

func TestTrendService_Reload(t *testing.T) {
	tr := &TrendService{}

	ctx := context.Background()

	err := tr.Reload(ctx, nil)

	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
}