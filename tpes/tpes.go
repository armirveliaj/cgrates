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

package tpes

import (
	"context"

	"github.com/cgrates/cgrates/config"
	"github.com/cgrates/cgrates/engine"
	"github.com/cgrates/cgrates/utils"
)

func NewTPeS(cfg *config.CGRConfig, cm *engine.ConnManager) (tpE *TPeS) {
	tpE = &TPeS{
		cfg:     cfg,
		connMgr: cm,
		exps:    make(map[string]tpExporter),
	}
	/*
		for expType := range tpExporterTypes {
			if tpE[expType], err = newTPExporter(expType, dm); err != nil {
				return nil, err
			}
		}
	*/
	return
}

// TPeS is managing the TariffPlanExporter
type TPeS struct {
	cfg     *config.CGRConfig
	connMgr *engine.ConnManager
	fltr    *engine.FilterS
	exps    map[string]tpExporter
}

type ArgsExportTP struct {
	Tenant      string
	APIOpts     map[string]interface{}
	ExportItems map[string][]string // map[expType][]string{"itemID1", "itemID2"}
}

// V1ExportTariffPlan is the API executed to export tariff plan items
func (tpE *TPeS) V1ExportTariffPlan(ctx *context.Context, args *ArgsExportTP, reply *utils.Account) (err error) {
	for eType := range args.ExportItems {
		if _, has := tpE.exps[eType]; !has {
			return utils.ErrPrefix(utils.ErrUnsupportedTPExporterType, eType)
		}
	}
	/*
		code to export to zip comes here
		for expType, expItms := range  args.ExportItems {
			if expCnt, err = tpE.exps[expType].exportItems(expItms); err != nil {
				return utils.NewErrServerError(err)
			}
		}
	*/
	return
}