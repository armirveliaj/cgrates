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
package utils

import (
	"testing"
)

func TestDynamicDataProviderProccesFieldPath(t *testing.T) {
	dp := MapStorage{
		MetaCgrep: MapStorage{
			"Stir": MapStorage{
				"CHRG_ROUTE1_END": "Identity1",
				"CHRG_ROUTE2_END": "Identity2",
				"CHRG_ROUTE3_END": "Identity3",
				"CHRG_ROUTE4_END": "Identity4",
			},
			"Routes": MapStorage{
				"SortedRoutes": []MapStorage{
					{"ID": "ROUTE1"},
					{"ID": "ROUTE2"},
					{"ID": "ROUTE3"},
					{"ID": "ROUTE4"},
				},
			},
			"AccountField": "Account",
			"Account":      "1001",
			"HalfAccount":  "01",
			"BestRoute":    0,
		},
	}
	newpath, err := ProcessFieldPath("~*cgrep.Stir.<CHRG_;~*cgrep.Routes.SortedRoutes[1].ID;_END>", InfieldSep, dp)
	if err != nil {
		t.Fatal(err)
	}
	expectedPath := "~*cgrep.Stir.CHRG_ROUTE2_END"
	if newpath != expectedPath {
		t.Errorf("Expected: %q,received %q", expectedPath, newpath)
	}
	_, err = ProcessFieldPath("~*cgrep.Stir.<CHRG_;~*cgrep.Routes.SortedRoutes[1].ID;_END", InfieldSep, dp)
	if err != ErrWrongPath {
		t.Errorf("Expected error %s received %v", ErrWrongPath, err)
	}

	_, err = ProcessFieldPath("~*cgrep.Stir<CHRG_", InfieldSep, dp)
	if err != ErrWrongPath {
		t.Errorf("Expected error %s received %v", ErrWrongPath, err)
	}

	_, err = ProcessFieldPath("~*cgrep.Stir<CHRG_;~*cgrep.Routes.SortedRoutes[1].ID2;_END>", InfieldSep, dp)
	if err != ErrNotFound {
		t.Errorf("Expected error %s received %v", ErrNotFound, err)
	}
	newpath, err = ProcessFieldPath("~*cgrep.Stir[1]", InfieldSep, dp)
	if err != nil {
		t.Fatal(err)
	}
	if newpath != EmptyString {
		t.Errorf("Expected: %q,received %q", EmptyString, newpath)
	}
	newpath, err = ProcessFieldPath("<*string:~*req.+~*cgrep.AccountField+:10+~*cgrep.HalfAccount+>", "+", dp)
	if err != nil {
		t.Fatal(err)
	}
	expectedPath = "*string:~*req.Account:1001"
	if newpath != expectedPath {
		t.Errorf("Expected: %q,received %q", expectedPath, newpath)
	}
	newpath, err = ProcessFieldPath("<*exists:~*req.+~*cgrep.AccountField+>:", "+", dp)
	if err != nil {
		t.Fatal(err)
	}
	expectedPath = "*exists:~*req.Account:"
	if newpath != expectedPath {
		t.Errorf("Expected: %q,received %q", expectedPath, newpath)
	}
	newpath, err = ProcessFieldPath("*exists:~*req.Account:", "+", dp)
	if err != nil {
		t.Fatal(err)
	}
	expectedPath = ""
	if newpath != expectedPath {
		t.Errorf("Expected: %q,received %q", expectedPath, newpath)
	}
}

func TestDynamicDataProviderGetFullFieldPath(t *testing.T) {
	dp := MapStorage{
		MetaCgrep: MapStorage{
			"Stir": MapStorage{
				"CHRG_ROUTE1_END": "Identity1",
				"CHRG_ROUTE2_END": "Identity2",
				"CHRG_ROUTE3_END": "Identity3",
				"CHRG_ROUTE4_END": "Identity4",
			},
			"Routes": MapStorage{
				"SortedRoutes": []MapStorage{
					{"ID": "ROUTE1"},
					{"ID": "ROUTE2"},
					{"ID": "ROUTE3"},
					{"ID": "ROUTE4"},
				},
			},
			"BestRoute": 0,
		},
	}
	newpath, err := GetFullFieldPath("~*cgrep.Stir.<CHRG_;~*cgrep.Routes.SortedRoutes[1].ID;_END>.Something", dp)
	if err != nil {
		t.Fatal(err)
	}
	expectedPath := "~*cgrep.Stir.CHRG_ROUTE2_END.Something"
	if newpath.Path != expectedPath {
		t.Errorf("Expected: %q,received %q", expectedPath, newpath.Path)
	}
	_, err = GetFullFieldPath("~*cgrep.Stir<CHRG_;~*cgrep.Routes.SortedRoutes[1].ID;_END.Something", dp)
	if err != ErrWrongPath {
		t.Errorf("Expected error %s received %v", ErrWrongPath, err)
	}

	_, err = GetFullFieldPath("~*cgrep.Stir<CHRG_", dp)
	if err != ErrWrongPath {
		t.Errorf("Expected error %s received %v", ErrWrongPath, err)
	}

	_, err = GetFullFieldPath("~*cgrep.Stir<CHRG_;~*cgrep.Routes.SortedRoutes[1].ID2;_END>", dp)
	if err != ErrNotFound {
		t.Errorf("Expected error %s received %v", ErrNotFound, err)
	}
	newpath, err = GetFullFieldPath("~*cgrep.Stir[1]", dp)
	if err != nil {
		t.Fatal(err)
	}
	if newpath != nil {
		t.Errorf("Expected: %v,received %q", nil, newpath)
	}

	newpath, err = GetFullFieldPath("~*cgrep.Stir", dp)
	if err != nil {
		t.Fatal(err)
	}
	if newpath != nil {
		t.Errorf("Expected: %v,received %q", nil, newpath)
	}

}
