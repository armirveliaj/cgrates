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

package v1

import (
	"testing"

	"github.com/cgrates/birpc/context"
	"github.com/cgrates/cgrates/config"
	"github.com/cgrates/cgrates/engine"
	"github.com/cgrates/cgrates/utils"
)

func TestAccountsSetGet(t *testing.T) {
	cfg := config.NewDefaultCGRConfig()
	db, err := engine.NewInternalDB(nil, nil, true, nil, cfg.DataDbCfg().Items)
	if err != nil {
		t.Error(err)
	}
	apierSv1 := &APIerSv1{
		DataManager: engine.NewDataManager(db, config.CgrConfig().CacheCfg(), nil),
		Config:      cfg,
	}
	cgrTenant := "cgrates.org"
	iscTenant := "itsyscom.com"
	b10 := &engine.Balance{Value: 10, Weight: 10}
	cgrAcnt1 := &engine.Account{ID: utils.ConcatenatedKey(cgrTenant, "account1"),
		BalanceMap: map[string]engine.Balances{utils.MetaMonetary + utils.MetaOut: {b10}}}
	cgrAcnt2 := &engine.Account{ID: utils.ConcatenatedKey(cgrTenant, "account2"),
		BalanceMap: map[string]engine.Balances{utils.MetaMonetary + utils.MetaOut: {b10}}}
	cgrAcnt3 := &engine.Account{ID: utils.ConcatenatedKey(cgrTenant, "account3"),
		BalanceMap: map[string]engine.Balances{utils.MetaMonetary + utils.MetaOut: {b10}}}
	iscAcnt1 := &engine.Account{ID: utils.ConcatenatedKey(iscTenant, "account1"),
		BalanceMap: map[string]engine.Balances{utils.MetaMonetary + utils.MetaOut: {b10}}}
	iscAcnt2 := &engine.Account{ID: utils.ConcatenatedKey(iscTenant, "account2"),
		BalanceMap: map[string]engine.Balances{utils.MetaMonetary + utils.MetaOut: {b10}}}
	for _, account := range []*engine.Account{cgrAcnt1, cgrAcnt2, cgrAcnt3, iscAcnt1, iscAcnt2} {
		if err := db.SetAccountDrv(account); err != nil {
			t.Error(err)
		}
	}

	var accounts []any
	var attrs utils.AttrGetAccounts
	if err := apierSv1.GetAccounts(context.Background(), &utils.AttrGetAccounts{Tenant: "cgrates.org"}, &accounts); err != nil {
		t.Error("Unexpected error", err.Error())
	} else if len(accounts) != 3 {
		t.Errorf("Accounts returned: %+v", accounts)
	}
	attrs = utils.AttrGetAccounts{Tenant: "itsyscom.com"}
	if err := apierSv1.GetAccounts(context.Background(), &attrs, &accounts); err != nil {
		t.Error("Unexpected error", err.Error())
	} else if len(accounts) != 2 {
		t.Errorf("Accounts returned: %+v", accounts)
	}
	attrs = utils.AttrGetAccounts{Tenant: "cgrates.org", AccountIDs: []string{"account1"}}
	if err := apierSv1.GetAccounts(context.Background(), &attrs, &accounts); err != nil {
		t.Error("Unexpected error", err.Error())
	} else if len(accounts) != 1 {
		t.Errorf("Accounts returned: %+v", accounts)
	}
	attrs = utils.AttrGetAccounts{Tenant: "itsyscom.com", AccountIDs: []string{"INVALID"}}
	if err := apierSv1.GetAccounts(context.Background(), &attrs, &accounts); err != nil {
		t.Error("Unexpected error", err.Error())
	} else if len(accounts) != 0 {
		t.Errorf("Accounts returned: %+v", accounts)
	}
	attrs = utils.AttrGetAccounts{Tenant: "INVALID"}
	if err := apierSv1.GetAccounts(context.Background(), &attrs, &accounts); err != nil {
		t.Error("Unexpected error", err.Error())
	} else if len(accounts) != 0 {
		t.Errorf("Accounts returned: %+v", accounts)
	}
}
