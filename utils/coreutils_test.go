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
	"errors"
	"fmt"
	"math"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"golang.org/x/crypto/bcrypt"
)

func TestFirstNonEmpty(t *testing.T) {
	//only check with an empty string
	rcv := FirstNonEmpty(EmptyString)
	if rcv != EmptyString {
		t.Errorf("Expecting an empty string, received: %+v", rcv)
	}
	//normal check
	firstElmnt := ""
	sampleMap := make(map[string]string)
	sampleMap["Third"] = "third"
	fourthElmnt := "fourth"
	winnerElmnt := FirstNonEmpty(firstElmnt, sampleMap["second"], sampleMap["Third"], fourthElmnt)
	if winnerElmnt != sampleMap["Third"] {
		t.Error("Wrong elemnt returned: ", winnerElmnt)
	}
}

func TestFirstIntNonEmpty(t *testing.T) {
	//check for default value
	rcv := FirstIntNonEmpty(0)
	if rcv != 0 {
		t.Errorf("Expected 0 \n but received \n %d", rcv)
	}

	//case the first non empty value is on the first position
	rcv = FirstIntNonEmpty(21, 0, 0)
	if rcv != 21 {
		t.Errorf("Expected 21 \n but received \n %d", rcv)
	}

	//case the first non empty value is on the second position
	rcv = FirstIntNonEmpty(0, 21, 0)
	if rcv != 21 {
		t.Errorf("Expected 21 \n but received \n %d", rcv)
	}

	//case the first non empty value is on the third position
	rcv = FirstIntNonEmpty(0, 0, 21)
	if rcv != 21 {
		t.Errorf("Expected 21 \n but received \n %d", rcv)
	}
}

func TestFirstDurationNonEmpty(t *testing.T) {
	//check for default value
	rcv := FirstDurationNonEmpty(0)
	if rcv != 0 {
		t.Errorf("Expected 0 \n but received \n %d", rcv)
	}

	//case the first non empty value is on the first position
	rcv = FirstDurationNonEmpty(2*time.Minute, 0, 0)
	if rcv != 2*time.Minute {
		t.Errorf("Expected 2m \n but received \n %d", rcv)
	}

	//case the first non empty value is on the second position
	rcv = FirstDurationNonEmpty(0, 2*time.Minute, 0)
	if rcv != 2*time.Minute {
		t.Errorf("Expected 2m \n but received \n %d", rcv)
	}

	//case the first non empty value is on the third position
	rcv = FirstDurationNonEmpty(0, 0, 2*time.Minute)
	if rcv != 2*time.Minute {
		t.Errorf("Expected 2m \n but received \n %d", rcv)
	}
}
func TestSha1(t *testing.T) {
	//empty check
	rcv := Sha1(" ")
	eOut := "b858cb282617fb0956d960215c8e84d1ccf909c6"
	if !reflect.DeepEqual(eOut, rcv) {
		t.Errorf("Expecting: %s, received: %s", eOut, rcv)
	}
	//normal check
	rcv = Sha1("teststring")
	eOut = "b8473b86d4c2072ca9b08bd28e373e8253e865c4"
	if !reflect.DeepEqual(eOut, rcv) {
		t.Errorf("Expecting: %s, received: %s", eOut, rcv)
	}
	rcv = Sha1("test1", "test2")
	eOut = "dff964f6e3c1761b6288f5c75c319d36fb09b2b9"
	if !reflect.DeepEqual(eOut, rcv) {
		t.Errorf("Expecting: %s, received: %s", eOut, rcv)
	}
}

func TestSha1ReverseOrder(t *testing.T) {
	rcv := Sha1("test1", "test2")
	revOrd := Sha1("test2", "test1")
	// Sha1 consider order when generating
	if reflect.DeepEqual(revOrd, rcv) {
		t.Errorf("Expecting: %s, received: %s", revOrd, rcv)
	}

	rcv = Sha1("test1")
	revOrd = Sha1("test1")
	if !reflect.DeepEqual(revOrd, rcv) {
		t.Errorf("Expecting: %s, received: %s", revOrd, rcv)
	}
}

func TestUUID(t *testing.T) {
	uuid := GenUUID()
	if len(uuid) == 0 {
		t.Fatalf("GenUUID error %s", uuid)
	}
	uuid2 := GenUUID()
	if len(uuid2) == 0 {
		t.Fatalf("GenUUID error %s", uuid)
	}
	if uuid == uuid2 {
		t.Error("GenUUID error.")
	}
}

func TestUUIDSha1Prefix(t *testing.T) {
	rcv := UUIDSha1Prefix()
	if len(rcv) != 7 {
		t.Errorf("Expected len: 7, received %d", len(rcv))
	}
	rcv2 := UUIDSha1Prefix()
	if len(rcv2) != 7 {
		t.Errorf("Expected len: 7, received %d", len(rcv))
	}
	if rcv == rcv2 {
		t.Error("UUIDSha1Prefix error")
	}
}

func TestRound(t *testing.T) {
	result := Round(12.49, 1, MetaRoundingUp)
	expected := 12.5
	if result != expected {
		t.Errorf("Error rounding up: sould be %v was %v", expected, result)
	}

	result = Round(12.21, 1, MetaRoundingUp)
	expected = 12.3
	if result != expected {
		t.Errorf("Error rounding up: sould be %v was %v", expected, result)
	}

	result = Round(0.0701, 2, MetaRoundingUp)
	expected = 0.08
	if result != expected {
		t.Errorf("Error rounding up: sould be %v was %v", expected, result)
	}

	result = Round(12.49, 1, MetaRoundingDown)
	expected = 12.4
	if result != expected {
		t.Errorf("Error rounding down: sould be %v was %v", expected, result)
	}

	result = Round(12.21, 1, MetaRoundingDown)
	expected = 12.2
	if result != expected {
		t.Errorf("Error rounding up: sould be %v was %v", expected, result)
	}

	//AlredyHavingPrecision
	x := 0.07
	if y := Round(x, 2, MetaRoundingUp); y != x {
		t.Error("Error rounding when already has desired precision: ", y)
	}
	if y := Round(x, 2, MetaRoundingMiddle); y != x {
		t.Error("Error rounding when already has desired precision: ", y)
	}
	if y := Round(x, 2, MetaRoundingDown); y != x {
		t.Error("Error rounding when already has desired precision: ", y)
	}

	result = Round(14.37, 8, MetaRoundingDown)
	expected = 14.37
	if result != expected {
		t.Errorf("Expecting: %v, received:  %v", expected, result)
	}
	result = Round(14.37, 8, "ROUNDING_NOWHERE")
	expected = 14.37
	if result != expected {
		t.Errorf("Expecting: %v, received:  %v", expected, result)
	}

	result = Round(14.37, 0, MetaRoundingMiddle)
	expected = 14
	if result != expected {
		t.Errorf("Expecting: %v, received:  %v", expected, result)
	}

	result = Round(14.37, -1, MetaRoundingMiddle)
	expected = 10
	if result != expected {
		t.Errorf("Expecting: %v, received:  %v", expected, result)
	}
}

func TestParseTimeDetectLayout(t *testing.T) {
	tmStr := "2013-12-30T15:00:01Z"
	expectedTime := time.Date(2013, 12, 30, 15, 0, 1, 0, time.UTC)
	tm, err := ParseTimeDetectLayout(tmStr, "")
	if err != nil {
		t.Error(err)
	} else if !tm.Equal(expectedTime) {
		t.Errorf("Unexpected time parsed: %v, expecting: %v", tm, expectedTime)
	}
	_, err = ParseTimeDetectLayout(tmStr[1:], "")
	if err == nil {
		t.Errorf("Expecting error")
	}
	tmStr = "2016-04-01T02:00:00+02:00"
	expectedTime = time.Date(2016, 4, 1, 0, 0, 0, 0, time.UTC)
	tm, err = ParseTimeDetectLayout(tmStr, "")
	if err != nil {
		t.Error(err)
	} else if !tm.Equal(expectedTime) {
		t.Errorf("Unexpected time parsed: %v, expecting: %v", tm, expectedTime)
	}
	_, err = ParseTimeDetectLayout(tmStr[1:], "")
	if err == nil {
		t.Errorf("Expecting error")
	}
	sqlTmStr := "2013-12-30 15:00:01"
	expectedTime = time.Date(2013, 12, 30, 15, 0, 1, 0, time.UTC)
	sqlTm, err := ParseTimeDetectLayout(sqlTmStr, "")
	if err != nil {
		t.Error(err)
	} else if !sqlTm.Equal(expectedTime) {
		t.Errorf("Unexpected time parsed: %v, expecting: %v", sqlTm, expectedTime)
	}
	_, err = ParseTimeDetectLayout(sqlTmStr[1:], "")
	if err == nil {
		t.Errorf("Expecting error")
	}
	unixTmStr := "1388415601"
	unixTm, err := ParseTimeDetectLayout(unixTmStr, "")
	if err != nil {
		t.Error(err)
	} else if !unixTm.Equal(expectedTime) {
		t.Errorf("Unexpected time parsed: %v, expecting: %v", unixTm, expectedTime)
	}
	_, err = ParseTimeDetectLayout(unixTmStr[1:], "")
	if err == nil {
		t.Errorf("Expecting error")
	}
	goTmStr := "2013-12-30 15:00:01 +0000 UTC"
	goTm, err := ParseTimeDetectLayout(goTmStr, "")
	if err != nil {
		t.Error(err)
	} else if !goTm.Equal(expectedTime) {
		t.Errorf("Unexpected time parsed: %v, expecting: %v", goTm, expectedTime)
	}
	_, err = ParseTimeDetectLayout(goTmStr[1:], "")
	if err == nil {
		t.Errorf("Expecting error")
	}
	goTmStr = "2013-12-30 15:00:01.000000000 +0000 UTC"
	goTm, err = ParseTimeDetectLayout(goTmStr, "")
	if err != nil {
		t.Error(err)
	} else if !goTm.Equal(expectedTime) {
		t.Errorf("Unexpected time parsed: %v, expecting: %v", goTm, expectedTime)
	}
	_, err = ParseTimeDetectLayout(goTmStr[1:], "")
	if err == nil {
		t.Errorf("Expecting error")
	}
	loc, err := time.LoadLocation("Asia/Kabul")
	if err != nil {
		t.Error(err)
	}
	expectedTime = time.Date(2013, 12, 30, 15, 0, 1, 0, loc)
	goTmStr2 := "2013-12-30 15:00:01 +0430 +0430"
	goTm, err = ParseTimeDetectLayout(goTmStr2, "")
	if err != nil {
		t.Error(err)
	} else if !goTm.Equal(expectedTime) {
		t.Errorf("Unexpected time parsed: %v, expecting: %v", goTm, expectedTime)
	}
	//goTmStr2 = "2013-12-30 15:00:01 +0430"
	//if _, err = ParseTimeDetectLayout(goTmStr2, ""); err != nil {
	//	t.Errorf("Expecting error")
	//}
	fsTmstampStr := "1394291049287234"
	fsTm, err := ParseTimeDetectLayout(fsTmstampStr, "")
	expectedTime = time.Date(2014, 3, 8, 15, 4, 9, 287234000, time.UTC)
	if err != nil {
		t.Error(err)
	} else if !fsTm.Equal(expectedTime) {
		t.Errorf("Unexpected time parsed: %v, expecting: %v", fsTm, expectedTime)
	}
	fsTmstampStr = "9999999999999999"
	fsTm, err = ParseTimeDetectLayout(fsTmstampStr, "")
	if err == nil {
		t.Error("Error expected: 'value out of range', received nil")
	}
	fsTmstampStr = "1394291049287234286"
	fsTm, err = ParseTimeDetectLayout(fsTmstampStr, "")
	expectedTime = time.Date(2014, 3, 8, 15, 4, 9, 287234286, time.UTC)
	if err != nil {
		t.Error(err)
	} else if !fsTm.Equal(expectedTime) {
		t.Errorf("Unexpected time parsed: %v, expecting: %v", fsTm, expectedTime)
	}
	fsTmstampStr = "9999999999999999999"
	fsTm, err = ParseTimeDetectLayout(fsTmstampStr, "")
	if err == nil {
		t.Error("Error expected: 'value out of range', received nil")
	}
	var nilTime time.Time
	fsTmstampStr = "+9999999999999999999"
	fsTm, err = ParseTimeDetectLayout(fsTmstampStr, "")
	if err == nil {
		t.Error("Error expected: 'value out of range', received nil")
	} else if fsTm != nilTime {
		t.Errorf("Expecting nilTime, received: %+v", fsTm)
	}
	fsTmstampStr = "0"
	fsTm, err = ParseTimeDetectLayout(fsTmstampStr, "")
	expectedTime = time.Time{}
	if err != nil {
		t.Error(err)
	} else if !fsTm.Equal(expectedTime) {
		t.Errorf("Unexpected time parsed: %v, expecting: %v", fsTm, expectedTime)
	}
	onelineTmstampStr := "20131023215149"
	olTm, err := ParseTimeDetectLayout(onelineTmstampStr, "")
	expectedTime = time.Date(2013, 10, 23, 21, 51, 49, 0, time.UTC)
	if err != nil {
		t.Error(err)
	} else if !olTm.Equal(expectedTime) {
		t.Errorf("Unexpected time parsed: %v, expecting: %v", olTm, expectedTime)
	}
	oneSpaceTmStr := "08.04.2014 22:14:29"
	tsTm, err := ParseTimeDetectLayout(oneSpaceTmStr, "")
	expectedTime = time.Date(2014, 4, 8, 22, 14, 29, 0, time.UTC)
	if err != nil {
		t.Error(err)
	} else if !tsTm.Equal(expectedTime) {
		t.Errorf("Unexpected time parsed: %v, expecting: %v", tsTm, expectedTime)
	}
	if nowTm, err := ParseTimeDetectLayout(MetaNow, ""); err != nil {
		t.Error(err)
	} else if time.Since(nowTm) > 10*time.Millisecond {
		t.Errorf("Unexpected time parsed: %v", nowTm)
	}
	eamonTmStr := "31/05/2015 14:46:00"
	eamonTmS, err := ParseTimeDetectLayout(eamonTmStr, "")
	expectedTime = time.Date(2015, 5, 31, 14, 46, 0, 0, time.UTC)
	if err != nil {
		t.Error(err)
	} else if !eamonTmS.Equal(expectedTime) {
		t.Errorf("Unexpected time parsed: %v, expecting: %v", eamonTmS, expectedTime)
	}
	broadSoftTmStr := "20160419210007.037"
	broadTmS, err := ParseTimeDetectLayout(broadSoftTmStr, "")
	expectedTime = time.Date(2016, 4, 19, 21, 0, 7, 37000000, time.UTC)
	if err != nil {
		t.Error(err)
	} else if !broadTmS.Equal(expectedTime) {
		t.Errorf("Expecting: %v, received: %v", expectedTime, broadTmS)
	}
	astTimestamp := "2016-09-14T19:37:43.665+0000"
	expectedTime = time.Date(2016, 9, 14, 19, 37, 43, 665000000, time.UTC)
	astTMS, err := ParseTimeDetectLayout(astTimestamp, "")
	if err != nil {
		t.Error(err)
	} else if !astTMS.Equal(expectedTime) {
		t.Errorf("Expecting: %v, received: %v", expectedTime, astTMS)
	}
	nowTimeStr := "+24h"
	start := time.Now().Add(23*time.Hour + 59*time.Minute + 58*time.Second)
	end := start.Add(2 * time.Second)
	parseNowTimeStr, err := ParseTimeDetectLayout(nowTimeStr, "")
	if err != nil {
		t.Error(err)
	} else if parseNowTimeStr.After(start) && parseNowTimeStr.Before(end) {
		t.Errorf("Unexpected time parsed: %v", parseNowTimeStr)
	}

	unixTmMilisecStr := "1534176053410"
	expectedTime = time.Date(2018, 8, 13, 16, 00, 53, 410000000, time.UTC)
	unixTm, err = ParseTimeDetectLayout(unixTmMilisecStr, "")
	if err != nil {
		t.Error(err)
	} else if !unixTm.Equal(expectedTime) {
		t.Errorf("Unexpected time parsed: %v, expecting: %v", unixTm, expectedTime)
	}

	tmStr = "2005-08-26T14:17:34"
	expectedTime = time.Date(2005, 8, 26, 14, 17, 34, 0, time.UTC)
	tm, err = ParseTimeDetectLayout(tmStr, "")
	if err != nil {
		t.Error(err)
	} else if !tm.Equal(expectedTime) {
		t.Errorf("Unexpected time parsed: %v, expecting: %v", tm, expectedTime)
	}

	expected := time.Date(2013, 7, 30, 19, 33, 10, 0, time.UTC)
	date, err := ParseTimeDetectLayout("1375212790", "")
	if err != nil || !date.Equal(expected) {
		t.Error("error parsing date: ", expected.Sub(date))
	}

	date, err = ParseTimeDetectLayout("*unlimited", "")
	if err != nil || !date.IsZero() {
		t.Error("error parsing unlimited date!: ")
	}

	date, err = ParseTimeDetectLayout("", "")
	if err != nil || !date.IsZero() {
		t.Error("error parsing unlimited date!: ")
	}

	date, err = ParseTimeDetectLayout("+20s", "")
	expected = time.Now()
	if err != nil || date.Sub(expected).Seconds() > 20 || date.Sub(expected).Seconds() < 19 {
		t.Error("error parsing date: ", date.Sub(expected).Seconds())
	}

	date, err = ParseTimeDetectLayout("-20s", "")
	expected = time.Now()
	if err != nil || expected.Sub(date).Seconds() > 21 || expected.Sub(date).Seconds() < 20 {
		t.Error("error parsing date: ", expected.Sub(date).Seconds())
	}

	expected = time.Now().AddDate(0, 0, 1)
	if date, err := ParseTimeDetectLayout("*daily", ""); err != nil {
		t.Error(err)
	} else if expected.Sub(date).Seconds() > 1 {
		t.Errorf("received: %+v", date)
	}
	expected = time.Now().AddDate(0, 1, 0)
	if date, err := ParseTimeDetectLayout("*monthly", ""); err != nil {
		t.Error(err)
	} else if expected.Sub(date).Seconds() > 1 {
		t.Errorf("received: %+v", date)
	}

	expected = time.Now().AddDate(0, 1, 0)
	if date, err := ParseTimeDetectLayout("*monthly_estimated", ""); err != nil {
		t.Error(err)
	} else {
		initialMnt := time.Now().Month()
		for expected.Month()-initialMnt > 1 {
			expected = expected.AddDate(0, 0, -1)
		}
		if expected.Day() != date.Day() || expected.Month() != date.Month() || expected.Hour() != date.Hour() {
			t.Errorf("Expected: %+v, received: %+v", expected, date)
		}
	}

	expected = time.Now().AddDate(0, 1, 0)
	if date, err := ParseTimeDetectLayout("*mo", ""); err != nil {
		t.Error(err)
	} else if expected.Sub(date).Seconds() > 1 {
		t.Errorf("received: %+v", date)
	}
	expected = time.Now().AddDate(0, 1, 0).Add(time.Hour + 2*time.Minute)
	if date, err := ParseTimeDetectLayout("*mo+1h2m", ""); err != nil {
		t.Error(err)
	} else if expected.Sub(date).Seconds() > 1 {
		t.Errorf("received: %+v", date)
	}
	expected = time.Now().AddDate(1, 0, 0)
	if date, err := ParseTimeDetectLayout("*yearly", ""); err != nil {
		t.Error(err)
	} else if expected.Sub(date).Seconds() > 1 {
		t.Errorf("received: %+v", date)
	}

	expected = GetEndOfMonth(time.Now())
	if date, err := ParseTimeDetectLayout("*month_end", ""); err != nil {
		t.Error(err)
	} else if !date.Equal(expected) {
		t.Errorf("received: %+v", date)
	}
	expected = GetEndOfMonth(time.Now()).Add(time.Hour).Add(2 * time.Minute)
	if date, err := ParseTimeDetectLayout("*month_end+1h2m", ""); err != nil {
		t.Error(err)
	} else if !date.Equal(expected) {
		t.Errorf("expecting: %+v, received: %+v", expected, date)
	}
	if date, err := ParseTimeDetectLayout("*month_end+xyz", ""); err == nil {
		t.Error("Expecting error 'time: invalid time duration', received: nil")
	} else if date != nilTime {
		t.Errorf("Expecting nilTime, received: %+v", date)
	}

	if date, err := ParseTimeDetectLayout("*mo+xyz", ""); err == nil {
		t.Error("Expecting error 'time: invalid time duration', received: nil")
	} else if date != nilTime {
		t.Errorf("Expecting nilTime, received: %+v", date)
	}

	date, err = ParseTimeDetectLayout("2013-07-30T19:33:10Z", "")
	expected = time.Date(2013, 7, 30, 19, 33, 10, 0, time.UTC)
	if err != nil || !date.Equal(expected) {
		t.Error("error parsing date: ", expected.Sub(date))
	}
	date, err = ParseTimeDetectLayout("2016-04-01T02:00:00+02:00", "")
	expected = time.Date(2016, 4, 1, 0, 0, 0, 0, time.UTC)
	if err != nil || !date.Equal(expected) {
		t.Errorf("Expecting: %v, received: %v", expected, date)
	}

	date, err = ParseTimeDetectLayout("2014-11-25T00:00:00+01:00", "")
	expected = time.Date(2014, 11, 24, 23, 0, 0, 0, time.UTC)
	if err != nil || !date.UTC().Equal(expected.UTC()) {
		t.Errorf("Expecting: %v, received: %v", expected.UTC(), date.UTC())
	}

}

func TestRoundDuration(t *testing.T) {
	minute := time.Minute
	result := RoundDuration(minute, 0)
	expected := 0 * time.Second
	if result != expected {
		t.Errorf("Error rounding to minute1: expected %v was %v", expected, result)
	}
	result = RoundDuration(time.Second, time.Second+500*time.Millisecond)
	expected = 2 * time.Second
	if result != expected {
		t.Errorf("Error rounding to minute1: expected %v was %v", expected, result)
	}
	result = RoundDuration(minute, time.Second)
	expected = minute
	if result != expected {
		t.Errorf("Error rounding to minute2: expected %v was %v", expected, result)
	}
	result = RoundDuration(minute, 5*time.Second)
	expected = minute
	if result != expected {
		t.Errorf("Error rounding to minute3: expected %v was %v", expected, result)
	}
	result = RoundDuration(minute, minute)
	expected = minute
	if result != expected {
		t.Errorf("Error rounding to minute4: expected %v was %v", expected, result)
	}
	result = RoundDuration(minute, 90*time.Second)
	expected = 120 * time.Second
	if result != expected {
		t.Errorf("Error rounding to minute5: expected %v was %v", expected, result)
	}
	result = RoundDuration(60, 120)
	expected = 120.0
	if result != expected {
		t.Errorf("Error rounding to minute5: expected %v was %v", expected, result)
	}
}

func TestRoundStatDuration(t *testing.T) {
	result := RoundStatDuration(time.Second+14565876*time.Nanosecond, 5)
	expected := time.Second + 14570000*time.Nanosecond
	if result != expected {
		t.Errorf("Expected %+v, received %+v", expected, result)
	}

	result = RoundStatDuration(time.Second+14565876*time.Nanosecond, 1)
	expected = time.Second
	if result != expected {
		t.Errorf("Expected %+v, received %+v", expected, result)
	}

	result = RoundStatDuration(time.Second+14565876*time.Nanosecond, 9)
	expected = time.Second + 14565876*time.Nanosecond
	if result != expected {
		t.Errorf("Expected %+v, received %+v", expected, result)
	}

	result = RoundStatDuration(24*time.Second+14565876*time.Nanosecond, -1)
	expected = 20 * time.Second
	if result != expected {
		t.Errorf("Expected %+v, received %+v", expected, result)
	}

	result = RoundStatDuration(24*time.Second+14565876*time.Nanosecond, 0)
	expected = 24 * time.Second
	if result != expected {
		t.Errorf("Expected %+v, received %+v", expected, result)
	}
}

func TestSplitPrefix(t *testing.T) {
	exp := []string{"0123456789", "012345678", "01234567", "0123456", "012345", "01234", "0123", "012", "01", "0"}
	if a := SplitPrefix("0123456789", 1); len(a) != 10 {
		t.Error("Error splitting prefix: ", a)
	} else if !reflect.DeepEqual(a, exp) {
		t.Errorf("Expecting: %v, received: %v", exp, a)
	}
	exp = []string{"0123456789", "012345678", "01234567", "0123456", "012345", "01234"}
	if a := SplitPrefix("0123456789", 5); len(a) != 6 {
		t.Error("Error splitting prefix: ", a)
	} else if !reflect.DeepEqual(a, exp) {
		t.Errorf("Expecting: %v, received: %v", exp, a)
	}
	exp = []string{}
	if a := SplitPrefix("", 1); len(a) != 0 {
		t.Error("Error splitting prefix: ", a)
	} else if !reflect.DeepEqual(a, exp) {
		t.Errorf("Expecting: %v, received: %v", exp, a)
	}
}

func TestSplitSuffix(t *testing.T) {
	exp := []string{"9", "89", "789", "6789", "56789", "456789", "3456789", "23456789", "123456789", "0123456789"}
	if a := SplitSuffix("0123456789"); len(a) != 10 {
		t.Error("Error splitting prefix: ", ToJSON(a))
	} else if !reflect.DeepEqual(a, exp) {
		t.Errorf("Expecting: %v, received: %v", exp, a)
	}
	if a := SplitSuffix(""); len(a) != 0 {
		t.Error("Error splitting prefix: ", a)
	}
}

func TestCopyHour(t *testing.T) {
	var src, dst, eOut time.Time
	if rcv := CopyHour(src, dst); !reflect.DeepEqual(rcv, eOut) {
		t.Errorf("Expecting: %+v, received: %+v", eOut, rcv)
	}
	src = time.Date(2020, time.April, 18, 20, 10, 11, 01, time.UTC)
	dst = time.Date(2019, time.May, 25, 23, 0, 4, 0, time.UTC)
	eOut = time.Date(2019, time.May, 25, 20, 10, 11, 01, time.UTC)
	if rcv := CopyHour(src, dst); !reflect.DeepEqual(rcv, eOut) {
		t.Errorf("Expecting: %+v, received: %+v", eOut, rcv)
	}
}

func TestParseDurationWithSecs(t *testing.T) {
	var durExpected time.Duration
	if rcv, err := ParseDurationWithSecs(""); err != nil {
		t.Error(err)
	} else if !reflect.DeepEqual(rcv, durExpected) {
		t.Errorf("Expecting: 0s, received: %+v", rcv)
	}
	durStr := "2"
	durExpected = 2 * time.Second
	if parsed, err := ParseDurationWithSecs(durStr); err != nil {
		t.Error(err)
	} else if parsed != durExpected {
		t.Error("Parsed different than expected")
	}
	durStr = "2s"
	if parsed, err := ParseDurationWithSecs(durStr); err != nil {
		t.Error(err)
	} else if parsed != durExpected {
		t.Error("Parsed different than expected")
	}
	durStr = "2ms"
	durExpected = 2 * time.Millisecond
	if parsed, err := ParseDurationWithSecs(durStr); err != nil {
		t.Error(err)
	} else if parsed != durExpected {
		t.Error("Parsed different than expected")
	}
	durStr = "0.002"
	durExpected = 2 * time.Millisecond
	if parsed, err := ParseDurationWithSecs(durStr); err != nil {
		t.Error(err)
	} else if parsed != durExpected {
		t.Error("Parsed different than expected")
	}
	durStr = "1.002"
	durExpected = 1002 * time.Millisecond
	if parsed, err := ParseDurationWithSecs(durStr); err != nil {
		t.Error(err)
	} else if parsed != durExpected {
		t.Error("Parsed different than expected")
	}
}

func TestParseDurationWithNanosecs(t *testing.T) {
	var eOut time.Duration
	if rcv, err := ParseDurationWithNanosecs(""); err != nil {
		t.Error(err)
	} else if !reflect.DeepEqual(rcv, eOut) {
		t.Errorf("Expecting: %+v, received: %+v", eOut, rcv)
	}
	eOut, _ = time.ParseDuration("-1ns")
	if rcv, err := ParseDurationWithNanosecs(MetaUnlimited); err != nil {
		t.Error(err)
	} else if !reflect.DeepEqual(rcv, eOut) {
		t.Errorf("Expecting: %+v, received: %+v", eOut, rcv)
	}
	eOut, _ = time.ParseDuration("28ns")
	if rcv, err := ParseDurationWithNanosecs("28"); err != nil {
		t.Error(err)
	} else if !reflect.DeepEqual(rcv, eOut) {
		t.Errorf("Expecting: %+v, received: %+v", eOut, rcv)
	}
}

func TestMinDuration(t *testing.T) {
	d1, _ := time.ParseDuration("1m")
	d2, _ := time.ParseDuration("59s")
	minD1 := MinDuration(d1, d2)
	minD2 := MinDuration(d2, d1)
	if minD1 != d2 || minD2 != d2 {
		t.Error("Error getting min duration: ", minD1, minD2)
	}
}

func TestParseZeroRatingSubject(t *testing.T) {
	subj := []string{"", "*zero1024", "*zero1s", "*zero5m", "*zero10h"}
	dur := []time.Duration{time.Second, 1024,
		time.Second, 5 * time.Minute, 10 * time.Hour}
	dfltRatingSubject := map[string]string{
		MetaAny:   "*zero1ns",
		MetaVoice: "*zero1s",
	}
	for i, s := range subj {
		if i == 0 {
			if d, err := ParseZeroRatingSubject(MetaVoice, s, dfltRatingSubject, false); err == nil {
				t.Error("Error parsing rating subject: ", s, d, err)
			}

		} else {
			if d, err := ParseZeroRatingSubject(MetaVoice, s, dfltRatingSubject, true); err != nil || d != dur[i] {
				t.Error("Error parsing rating subject: ", s, d, err)
			}

		}
		if d, err := ParseZeroRatingSubject(MetaData, EmptyString, dfltRatingSubject, true); err != nil || d != time.Nanosecond {
			t.Error("Error parsing rating subject: ", EmptyString, d, err)
		}
		if d, err := ParseZeroRatingSubject(MetaSMS, EmptyString, dfltRatingSubject, true); err != nil || d != time.Nanosecond {
			t.Error("Error parsing rating subject: ", EmptyString, d, err)
		}
		if d, err := ParseZeroRatingSubject(MetaMMS, EmptyString, dfltRatingSubject, true); err != nil || d != time.Nanosecond {
			t.Error("Error parsing rating subject: ", EmptyString, d, err)
		}
		if d, err := ParseZeroRatingSubject(MetaMonetary, EmptyString, dfltRatingSubject, true); err != nil || d != time.Nanosecond {
			t.Error("Error parsing rating subject: ", EmptyString, d, err)
		}
		expecting := "malformed rating subject: test"
		if _, err := ParseZeroRatingSubject(MetaMonetary, "test", dfltRatingSubject, true); err == nil || err.Error() != expecting {
			t.Errorf("Expecting: %+v, received: %+v ", expecting, err)
		}
	}
}
func TestConcatenatedKey(t *testing.T) {
	if key := ConcatenatedKey("a"); key != "a" {
		t.Error("Unexpected key value received: ", key)
	}
	if key := ConcatenatedKey("a", "b"); key != fmt.Sprintf("a%sb", ConcatenatedKeySep) {
		t.Error("Unexpected key value received: ", key)
	}
	if key := ConcatenatedKey("a", "b", "c"); key != fmt.Sprintf("a%sb%sc", ConcatenatedKeySep, ConcatenatedKeySep) {
		t.Error("Unexpected key value received: ", key)
	}
}

func TestSplitConcatenatedKey(t *testing.T) {
	key := "test1:test2:test3"
	eOut := []string{"test1", "test2", "test3"}
	if rcv := SplitConcatenatedKey(key); !reflect.DeepEqual(eOut, rcv) {
		t.Errorf("Expecting: %+v, received: %+v", eOut, rcv)
	}
}

func TestInfieldJoin(t *testing.T) {
	if rcv := InfieldJoin(""); rcv != "" {
		t.Errorf("Expecting: empty string, received: %+v", ToJSON(rcv))
	}
	key := "test1;test2;test3"

	eOut := "test1;test2;test3"
	if rcv := InfieldJoin(key); !reflect.DeepEqual(eOut, rcv) {
		t.Errorf("Expecting: %+v, received: %+v", ToJSON(eOut), ToJSON(rcv))
	}
	key2 := "test10;test12"

	eOut = "test1;test2;test3;test10;test12"
	if rcv := InfieldJoin(key, key2); !reflect.DeepEqual(eOut, rcv) {
		t.Errorf("Expecting: %+v, received: %+v", ToJSON(eOut), ToJSON(rcv))
	}
}

func TestInfieldSplit(t *testing.T) {
	key := "test1;test2;test3"
	eOut := []string{"test1", "test2", "test3"}
	if rcv := InfieldSplit(key); !reflect.DeepEqual(eOut, rcv) {
		t.Errorf("Expecting: %+v, received: %+v", eOut, rcv)
	}
}

func TestFmtFieldWidth(t *testing.T) {
	//Mandatory0
	if _, err := FmtFieldWidth("", "", 0, "", "", true); err == nil {
		t.Errorf("Failed to detect mandatory value")
	}
	//width = 0
	if result, err := FmtFieldWidth("", "test", 0, "", "", false); err != nil {
		t.Error(err)
	} else if result != "test" {
		t.Errorf("Expecting 'test', received %+q", result)
	}
	//MaxLen
	if result, err := FmtFieldWidth("", "test", 4, "", "", false); err != nil {
		t.Error(err)
	} else if result != "test" {
		t.Errorf("Expected \"test\" received: \"%s\"", result)
	}
	//RPadding
	if result, err := FmtFieldWidth("", "test", 8, "", "*right", false); err != nil {
		t.Error(err)
	} else if result != "test    " {
		t.Errorf("Expected <\"test    \"> \" received: \"%s\"", result)
	}
	//PaddingFiller
	expected := "        "
	if result, err := FmtFieldWidth("", "", 8, "", "*right", false); err != nil {
		t.Error(err)
	} else if result != expected {
		t.Errorf("Expected \"%s \" received: \"%s\"", expected, result)
	}
	//LPadding
	expected = "    test"
	if result, err := FmtFieldWidth("", "test", 8, "", "*left", false); err != nil {
		t.Error(err)
	} else if result != expected {
		t.Errorf("Expected \"%s \" received: \"%s\"", expected, result)
	}
	//ZeroLPadding
	expected = "0000test"
	if result, err := FmtFieldWidth("", "test", 8, "", "*zeroleft", false); err != nil {
		t.Error(err)
	} else if result != expected {
		t.Errorf("Expected \"%s \" received: \"%s\"", expected, result)
	}
	//RStrip
	expected = "te"
	if result, err := FmtFieldWidth("", "test", 2, "*right", "", false); err != nil {
		t.Error(err)
	} else if result != expected {
		t.Errorf("Expected \"%s \" received: \"%s\"", expected, result)
	}
	//XRStrip
	expected = "tex"
	if result, err := FmtFieldWidth("", "test", 3, "*xright", "", false); err != nil {
		t.Error(err)
	} else if result != expected {
		t.Errorf("Expected \"%s \" received: \"%s\"", expected, result)
	}
	//LStrip
	expected = "st"
	if result, err := FmtFieldWidth("", "test", 2, "*left", "", false); err != nil {
		t.Error(err)
	} else if result != expected {
		t.Errorf("Expected \"%s \" received: \"%s\"", expected, result)
	}
	//XLStrip
	expected = "xst"
	if result, err := FmtFieldWidth("", "test", 3, "*xleft", "", false); err != nil {
		t.Error(err)
	} else if result != expected {
		t.Errorf("Expected \"%s \" received: \"%s\"", expected, result)
	}
	//StripNotAllowed
	if _, err := FmtFieldWidth("", "test", 3, "", "", false); err == nil {
		t.Error("Expected error")
	}
	//PaddingNotAllowed
	if _, err := FmtFieldWidth("", "test", 5, "", "", false); err == nil {
		t.Error("Expected error")
	}
	//
	expected = "test"
	if result, err := FmtFieldWidth("", "test", 3, "wrong", "", false); err != nil {
		t.Error(err)
	} else if result != expected {
		t.Errorf("Expected \"%s \" received: \"%s\"", expected, result)
	}
}

func TestEndOfMonth(t *testing.T) {
	eom := GetEndOfMonth(time.Date(2016, time.February, 5, 10, 1, 2, 3, time.UTC))
	expected := time.Date(2016, time.February, 29, 23, 59, 59, 0, time.UTC)
	if !eom.Equal(expected) {
		t.Errorf("Expected %v was %v", expected, eom)
	}
	eom = GetEndOfMonth(time.Date(2015, time.February, 5, 10, 1, 2, 3, time.UTC))
	expected = time.Date(2015, time.February, 28, 23, 59, 59, 0, time.UTC)
	if !eom.Equal(expected) {
		t.Errorf("Expected %v was %v", expected, eom)
	}
	eom = GetEndOfMonth(time.Date(2016, time.January, 31, 10, 1, 2, 3, time.UTC))
	expected = time.Date(2016, time.January, 31, 23, 59, 59, 0, time.UTC)
	if !eom.Equal(expected) {
		t.Errorf("Expected %v was %v", expected, eom)
	}
	eom = GetEndOfMonth(time.Date(2016, time.December, 31, 10, 1, 2, 3, time.UTC))
	expected = time.Date(2016, time.December, 31, 23, 59, 59, 0, time.UTC)
	if !eom.Equal(expected) {
		t.Errorf("Expected %v was %v", expected, eom)
	}
	eom = GetEndOfMonth(time.Date(2016, time.July, 31, 23, 59, 59, 0, time.UTC))
	expected = time.Date(2016, time.July, 31, 23, 59, 59, 0, time.UTC)
	if !eom.Equal(expected) {
		t.Errorf("Expected %v was %v", expected, eom)
	}
	// Date ref represents zero time instant, will return time.Now().
	eom = GetEndOfMonth(time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC))
	if currentDateRef := time.Now(); currentDateRef.Before(eom) ||
		time.Since(eom) > time.Millisecond {
		t.Errorf("Expected GetEndOfMonth to return the current date info, received: %s", eom)
	}
}

func TestSizeFmt(t *testing.T) {
	if str := SizeFmt(0, EmptyString); str != "0.0B" {
		t.Errorf("Expecting: 0.0B, received: %+q", str)
	}
	if str := SizeFmt(1023, EmptyString); str != "1023.0B" {
		t.Errorf("Expecting: 1023.0B, received: %+q", str)
	}
	if str := SizeFmt(1024, EmptyString); str != "1.0KiB" {
		t.Errorf("Expecting: 1.0KiB, received: %+q", str)
	}
	if str := SizeFmt(1048575, EmptyString); str != "1024.0KiB" {
		t.Errorf("Expecting: 1024.0KiB, received: %+q", str)
	}
	if str := SizeFmt(1048576, EmptyString); str != "1.0MiB" {
		t.Errorf("Expecting: 1.0MiB, received: %+q", str)
	}
	if str := SizeFmt(1073741823, EmptyString); str != "1024.0MiB" {
		t.Errorf("Expecting: 1024.0MiB, received: %+q", str)
	}
	if str := SizeFmt(1073741824, EmptyString); str != "1.0GiB" {
		t.Errorf("Expecting: 1.0GiB, received: %+q", str)
	}
	if str := SizeFmt(1099511627775, EmptyString); str != "1024.0GiB" {
		t.Errorf("Expecting: 1024.0GiB, received: %+q", str)
	}
	if str := SizeFmt(1099511627776, EmptyString); str != "1.0TiB" {
		t.Errorf("Expecting: 1.0TiB, received: %+q", str)
	}
	if str := SizeFmt(1125899906842623, EmptyString); str != "1024.0TiB" {
		t.Errorf("Expecting: 1024.0TiB, received: %+q", str)
	}
	if str := SizeFmt(1125899906842624, EmptyString); str != "1.0PiB" {
		t.Errorf("Expecting: 1.0PiB, received: %+q", str)
	}
	if str := SizeFmt(1152921504606847000, EmptyString); str != "1.0EiB" {
		t.Errorf("Expecting: 1.0EiB, received: %+q", str)
	}
	if str := SizeFmt(1180591620717411303424, EmptyString); str != "1.0ZiB" {
		t.Errorf("Expecting: 1.0ZiB, received: %+q", str)
	}
	if str := SizeFmt(9000000000000000000000000, EmptyString); str != "7.4YiB" {
		t.Errorf("Expecting: 7.4YiB, received: %+q", str)
	}

}

func TestParseHierarchyPath(t *testing.T) {
	eHP := HierarchyPath([]string{"Root", "CGRateS"})
	if hp := ParseHierarchyPath("/Root/CGRateS/", ""); !reflect.DeepEqual(hp, eHP) {
		t.Errorf("Expecting: %+v, received: %+v", eHP, hp)
	}
	if hp := ParseHierarchyPath("Root.CGRateS", ""); !reflect.DeepEqual(hp, eHP) {
		t.Errorf("Expecting: %+v, received: %+v", eHP, hp)
	}
}

func TestHierarchyPathAsString(t *testing.T) {
	tests := []struct {
		name     string
		hP       HierarchyPath
		sep      string
		prefix   bool
		expected string
	}{
		{
			name:     "Empty HierarchyPath without prefix",
			hP:       HierarchyPath{},
			sep:      "/",
			prefix:   false,
			expected: ".",
		},
		{
			name:     "Empty HierarchyPath with prefix",
			hP:       HierarchyPath{},
			sep:      "/",
			prefix:   true,
			expected: "/",
		},
		{
			name:     "Single element without prefix",
			hP:       HierarchyPath{"element"},
			sep:      "/",
			prefix:   false,
			expected: "element",
		},
		{
			name:     "Single element with prefix",
			hP:       HierarchyPath{"element"},
			sep:      "/",
			prefix:   true,
			expected: "/element",
		},
		{
			name:     "Multiple elements without prefix",
			hP:       HierarchyPath{"a", "b", "c"},
			sep:      "/",
			prefix:   false,
			expected: "a/b/c",
		},
		{
			name:     "Multiple elements with prefix",
			hP:       HierarchyPath{"a", "b", "c"},
			sep:      "/",
			prefix:   true,
			expected: "/a/b/c",
		},
		{
			name:     "Custom separator without prefix",
			hP:       HierarchyPath{"a", "b", "c"},
			sep:      "->",
			prefix:   false,
			expected: "a->b->c",
		},
		{
			name:     "Custom separator with prefix",
			hP:       HierarchyPath{"a", "b", "c"},
			sep:      "->",
			prefix:   true,
			expected: "->a->b->c",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if received := tt.hP.AsString(tt.sep, tt.prefix); received != tt.expected {
				t.Errorf("expected <%s>, received <%s>", tt.expected, received)
			}
		})
	}
}

func TestParseHierarchyClone(t *testing.T) {
	eHP := HierarchyPath([]string{"Root", "CGRateS"})
	rcv := eHP.Clone()
	if !reflect.DeepEqual(eHP, rcv) {
		t.Errorf("Expected: %+v\nReceived: %+v", ToJSON(eHP), ToJSON(rcv))
	}
	if rcv[0] = ""; eHP[0] != "Root" {
		t.Errorf("Expected clone to not modify the cloned")
	}
	eHP = nil
	rcv = eHP.Clone()
	if !reflect.DeepEqual(eHP, rcv) {
		t.Errorf("Expected: %+v\nReceived: %+v", ToJSON(eHP), ToJSON(rcv))
	}
}

func TestMaskSuffix(t *testing.T) {
	dest := "+4986517174963"
	if destMasked := MaskSuffix(dest, 3); destMasked != "+4986517174***" {
		t.Error("Unexpected mask applied", destMasked)
	}
	if destMasked := MaskSuffix(dest, -1); destMasked != dest {
		t.Error("Negative maskLen should not modify destination", destMasked)
	}
	if destMasked := MaskSuffix(dest, 0); destMasked != dest {
		t.Error("Zero maskLen should not modify destination", destMasked)
	}
	if destMasked := MaskSuffix(dest, 100); destMasked != "**************" {
		t.Error("High maskLen should return complete mask", destMasked)
	}
}

func TestTimeIs0h(t *testing.T) {
	t1, err := time.Parse(time.RFC3339, "2012-11-01T22:08:41+00:00")
	if err != nil {
		t.Error("time parsing error")
	}
	result := TimeIs0h(t1)
	if result != false {
		t.Error("time is 0 when it's supposed to be", t1)
	}
}

func TestToJSON(t *testing.T) {
	if outNilObj := ToJSON(nil); outNilObj != "null" {
		t.Errorf("Expecting null, received: <%q>", outNilObj)
	}
}

func TestClone(t *testing.T) {
	a := 15
	var b int
	err := Clone(a, &b)
	if err != nil {
		t.Error("Cloning failed")
	}
	if b != a {
		t.Error("Expected:", a, ", received:", b)
	}
	// Clone from an interface
	c := "mystr"
	ifaceC := any(c)
	clndIface := reflect.Indirect(reflect.New(reflect.TypeOf(ifaceC))).Interface().(string)
	if err := Clone(ifaceC, &clndIface); err != nil {
		t.Error(err)
	} else if !reflect.DeepEqual(ifaceC, clndIface) {
		t.Errorf("Expecting: %+v, received: %+v", ifaceC, clndIface)
	}

	if err := Clone(math.NaN, nil); err == nil {
		t.Error("Expected error")
	}
}

func TestFib(t *testing.T) {
	fib := Fib()
	if tmp := fib(); tmp != 1 {
		t.Error("Expecting: 1s, received ", tmp)
	}
	if tmp := fib(); tmp != 1 {
		t.Error("Expecting: 1s, received ", tmp)
	}
	if tmp := fib(); tmp != 2 {
		t.Error("Expecting: 2s, received ", tmp)
	}
	if tmp := fib(); tmp != 3 {
		t.Error("Expecting: 3s, received ", tmp)
	}
	if tmp := fib(); tmp != 5 {
		t.Error("Expecting: 5s, received ", tmp)
	}
}

func TestStringPointer(t *testing.T) {
	result := StringPointer("*zero")
	if *result != EmptyString {
		t.Errorf("Expecting: `%+q`, received `%+q`", EmptyString, *result)
	}
	str := "test_string"
	result = StringPointer(str)
	expected := &str
	if *result != *expected {
		t.Errorf("Expecting: %+v, received: %+v", &str, result)
	}
}
func TestIntPointer(t *testing.T) {
	t1 := 14
	result := IntPointer(t1)
	expected := &t1
	if *expected != *result {
		t.Errorf("Expecting: %+v, received: %+v", expected, result)
	}
}

func TestInt64Pointer(t *testing.T) {
	var t1 int64 = 19
	result := Int64Pointer(t1)
	expected := &t1
	if *expected != *result {
		t.Error("Expected:", expected, ", received:", result)
	}
}

func TestFloat64Pointer(t *testing.T) {
	var t1 float64 = 11.5
	result := Float64Pointer(t1)
	expected := &t1
	if *expected != *result {
		t.Error("Expected:", expected, ", received:", result)
	}
}

func TestBoolPointer(t *testing.T) {
	t1 := true
	result := BoolPointer(t1)
	expected := &t1
	if *expected != *result {
		t.Error("Expected:", expected, ", received:", result)
	}
}

func TestStringMapPointer(t *testing.T) {
	t1 := map[string]bool{"cgr1": true, "cgr2": true}
	expected := &t1
	result := StringMapPointer(t1)
	if *result == nil {
		t.Error("Expected:", expected, ", received: nil")
	}
}

func TestMapStringStringPointer(t *testing.T) {
	mp := map[string]string{"string1": "string2"}
	result := MapStringStringPointer(mp)
	expected := &mp
	if *result == nil {
		t.Errorf("Expected: %+q, received: nil", expected)
	}
}

func TestTimePointer(t *testing.T) {
	t1, err := time.Parse(time.RFC3339, "2012-11-01T22:08:41+00:00")
	if err != nil {
		t.Error("time parsing error")
	}
	result := TimePointer(t1)
	expected := &t1
	if *expected != *result {
		t.Error("Expected:", expected, ", received:", result)
	}
}

func TestDurationPointer(t *testing.T) {
	duration := 10 * time.Nanosecond
	result := DurationPointer(duration)
	expected := &duration
	if *expected != *result {
		t.Errorf("Expected: %+q, received: %+q", expected, result)
	}
}

func TestSliceStringPointer(t *testing.T) {
	sl := []string{"test"}
	result := SliceStringPointer(sl)
	exp := result
	if !reflect.DeepEqual(result, exp) {
		t.Errorf("Expeced:%+q but received %+q", exp, result)
	}

}

func TestToIJSON(t *testing.T) {
	str := "string"
	received := ToIJSON(str)
	expected := "\"string\""
	if !reflect.DeepEqual(received, expected) {
		t.Errorf("Expected: %+q, received: %+q", expected, received)
	}
}

func TestLen(t *testing.T) {
	t1 := Int64Slice{2112443, 414241412, 41231241}
	result := t1.Len()
	expected := 3
	if result != expected {
		t.Error("Expected:", expected, ", received:", result)
	}
}

func TestSwap(t *testing.T) {
	t1 := Int64Slice{414241412, 41231241}
	t1.Swap(0, 1)
	var expected int64 = 414241412
	if t1[1] != expected {
		t.Error("Expected:", expected, ", received:", t1[1])
	}
}

func TestLess(t *testing.T) {
	t1 := Int64Slice{414241412, 41231241}
	expected := false
	if t1.Less(0, 1) != expected {
		t.Error("Expected:", expected, ", received:", t1.Less(1, 2))
	}
}

func TestGetCGRVersion(t *testing.T) {
	tests := []struct {
		name     string
		date     string
		hash     string
		want     string
		wantErr  string
		parseErr bool
	}{
		{
			name: "successful version build (git2.45+ date format)",
			date: "2016-12-30T19:48:09Z",
			hash: "73014daa0c1d7edcb532d5fe600b8a20d588cdf8",
			want: "CGRateS@" + Version + "-20161230194809-73014daa0c1d",
		},
		{
			name: "successful version build (previous date format)",
			date: "2016-12-30T19:48:09+01:00",
			hash: "73014daa0c1d7edcb532d5fe600b8a20d588cdf8",
			want: "CGRateS@" + Version + "-20161230184809-73014daa0c1d",
		},
		{
			name: "unknown verb %cI (git version is too old)",
			date: "%cI",
			hash: "73014daa0c1d7edcb532d5fe600b8a20d588cdf8",
			want: "CGRateS@v0.11.0~dev-00010101000000-73014daa0c1d",
		},
		{
			name: "successful version build (default)",
			date: "",
			hash: "",
			want: "CGRateS@" + Version,
		},
		{
			name:     "wrong date format",
			date:     "wrong format",
			hash:     "73014daa0c1d7edcb532d5fe600b8a20d588cdf8",
			parseErr: true,
		},
		{
			name:     "wrong hash",
			date:     "2016-12-30T19:48:09+01:00",
			hash:     "73014DAA0C1D7EDCB532D5FE600B8A20D588CDF8",
			want:     "CGRateS@" + Version + "-20161230184809-73014daa0c1d",
			wantErr:  "version build error: commit hash does not match expected format",
			parseErr: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			GitCommitDate = test.date
			GitCommitHash = test.hash
			got, err := GetCGRVersion()
			switch {
			case test.parseErr:
				var parseError *time.ParseError
				if !errors.As(err, &parseError) {
					t.Fatalf("\nCommitDate: %q\nCommitHash: %q\nGetCGRVersion() err %q, want error of type *time.ParseError",
						test.date, test.hash, err)
				}
			case err != nil:
				if test.wantErr == "" {
					t.Fatalf("\nCommitDate: %q\nCommitHash: %q\nGetCGRVersion() err %s, want nil",
						test.date, test.hash, err)
				}
				if err.Error() != test.wantErr {
					t.Fatalf("\nCommitDate: %q\nCommitHash: %q\nGetCGRVersion() err %q, want %q",
						test.date, test.hash, err, test.wantErr)
				}
			case test.wantErr != "":
				t.Fatalf("\nCommitDate: %q\nCommitHash: %q\nGetCGRVersion() err nil, want %q",
					test.date, test.hash, test.wantErr)
			case got != test.want:
				t.Errorf("\nCommitDate: %q\nCommitHash: %q\nGetCGRVersion() = %q, want %q",
					test.date, test.hash, got, test.want)
			}
		})
	}
}

func TestNewTenantID(t *testing.T) {
	eOut := &TenantID{ID: EmptyString}
	if rcv := NewTenantID(EmptyString); *rcv != *eOut {
		t.Errorf("Expecting: %+v, received %+v", eOut, rcv)
	}
	eOut = &TenantID{Tenant: "Test"}
	if rcv := NewTenantID("Test:"); *rcv != *eOut {
		t.Errorf("Expecting: %+v, received %+v", eOut, rcv)
	}
	eOut = &TenantID{Tenant: "cgrates.org", ID: "id"}
	if rcv := NewTenantID("cgrates.org:id"); *rcv != *eOut {
		t.Errorf("Expecting: %+v, received %+v", eOut, rcv)
	}
}

func TestTenantID(t *testing.T) {
	tID := &TenantID{Tenant: EmptyString, ID: EmptyString}
	eOut := ":"
	if rcv := tID.TenantID(); rcv != eOut {
		t.Errorf("Expecting: %q, received: %q", eOut, rcv)
	}
	tID = &TenantID{Tenant: "cgrates.org", ID: "id"}
	eOut = "cgrates.org:id"
	if rcv := tID.TenantID(); rcv != eOut {
		t.Errorf("Expecting: %q, received: %q", eOut, rcv)
	}
}

func TestTenantIDWithCache(t *testing.T) {
	tID := &TenantIDWithAPIOpts{TenantID: &TenantID{Tenant: EmptyString, ID: EmptyString}}
	eOut := ":"
	if rcv := tID.TenantID.TenantID(); rcv != eOut {
		t.Errorf("Expecting: %q, received: %q", eOut, rcv)
	}
	tID = &TenantIDWithAPIOpts{TenantID: &TenantID{Tenant: "cgrates.org", ID: "id"}}
	eOut = "cgrates.org:id"
	if rcv := tID.TenantID.TenantID(); rcv != eOut {
		t.Errorf("Expecting: %q, received: %q", eOut, rcv)
	}
}

type TestRPC struct {
}

func (tRPC *TestRPC) V1Copy(args *string, reply *string) error {
	*reply = *args
	return nil
}

func (tRPC *TestRPC) V1Error(args *string, reply *string) error {
	return errors.New("V1_err_test")
}

func (tRPC *TestRPC) Call(args any, reply any) error {
	return nil
}

func (tRPC *TestRPC) V1Error2(args any, reply any) (int, error) {
	return 0, nil
}

func (tRPC *TestRPC) V1Error3(args any, reply any) int {
	return 0
}

type TestRPC2 struct {
}

func (tRPC *TestRPC2) Copy(args *string, reply *string) error {
	*reply = *args
	return nil
}

func (tRPC *TestRPC2) Error(args *string, reply *string) error {
	return errors.New("V1_err_test")
}

func (tRPC *TestRPC2) Call(args any, reply any) error {
	return nil
}

func (tRPC *TestRPC2) Error2(args any, reply any) (int, error) {
	return 0, nil
}

func (tRPC *TestRPC2) Error3(args any, reply any) int {
	return 0
}

func TestCounter(t *testing.T) {
	var cmax int64 = 10000
	var i int64
	cnter := NewCounter(0, cmax)
	for i = 1; i <= cmax; i++ {
		if i != cnter.Next() {
			t.Error("Counter has unexpected value")
		}
	}
	if cnter.Next() != 0 {
		t.Error("Counter did not reset")
	}
}

func TestCounterConcurrent(t *testing.T) {
	var nmax int64 = 10000
	ch := make(chan int64, nmax)
	wg := new(sync.WaitGroup)
	cnter := NewCounter(0, nmax-1)
	var i int64
	for i = 1; i <= nmax; i++ {
		wg.Add(1)
		go func() {
			ch <- cnter.Next()
			wg.Done()
		}()
	}
	wg.Wait()
	m := make(map[int64]bool)
	for i = 1; i <= nmax; i++ {
		m[<-ch] = true
	}
	for i = 1; i <= nmax-1; i++ {
		if !m[i] {
			t.Errorf("Missing value: %d", i)
		}
	}
	if cnter.Value() != 0 {
		t.Error("Counter was not reseted to 0")
	}
}

func TestReverseString(t *testing.T) {
	if rcv := ReverseString(EmptyString); rcv != EmptyString {
		t.Errorf("Expecting <%+q>, received: <%+q>", EmptyString, rcv)
	}
	if rcv := ReverseString("test"); rcv != "tset" {
		t.Errorf("Expecting <tset>, received: <%+q>", rcv)
	}
}

func TestGetUrlRawArguments(t *testing.T) {
	eOut := map[string]string{}
	if rcv := GetUrlRawArguments(EmptyString); !reflect.DeepEqual(rcv, eOut) {
		t.Errorf("Expectinc: %+v, received: %+v", eOut, rcv)
	}
	if rcv := GetUrlRawArguments("test"); !reflect.DeepEqual(rcv, eOut) {
		t.Errorf("Expectinc: %+v, received: %+v", eOut, rcv)
	}
	if rcv := GetUrlRawArguments("test?"); !reflect.DeepEqual(rcv, eOut) {
		t.Errorf("Expectinc: %+v, received: %+v", eOut, rcv)
	}
	eOut = map[string]string{"test": "1"}
	if rcv := GetUrlRawArguments("?test=1"); !reflect.DeepEqual(rcv, eOut) {
		t.Errorf("Expectinc: %+v, received: %+v", eOut, rcv)
	}
	eOut = map[string]string{"test": "1", "test2": "2"}
	if rcv := GetUrlRawArguments("?test=1&test2=2"); !reflect.DeepEqual(rcv, eOut) {
		t.Errorf("Expectinc: %+v, received: %+v", eOut, rcv)
	}
	eOut = map[string]string{"test": "1", "test2": "2", EmptyString: "5"}
	if rcv := GetUrlRawArguments("?test=1&test2=2&=5"); !reflect.DeepEqual(rcv, eOut) {
		t.Errorf("Expectinc: %+v, received: %+v", eOut, rcv)
	}
	eOut = map[string]string{"test": "1", "test2": "2"}
	if rcv := GetUrlRawArguments("?test=1&test2=2&5"); !reflect.DeepEqual(rcv, eOut) {
		t.Errorf("Expectinc: %+v, received: %+v", eOut, rcv)
	}
}

func TestWarnExecTime(t *testing.T) {
	//without Log
	WarnExecTime(time.Now(), "MyTestFunc", time.Second)
	//With Log
	WarnExecTime(time.Now(), "MyTestFunc", time.Nanosecond)
}

func TestCastRPCErr(t *testing.T) {
	err := errors.New("test")
	if rcv := CastRPCErr(err); rcv != err {
		t.Errorf("Expecting: %+q, received %+q", err, rcv)
	}
	if rcv := CastRPCErr(ErrNoMoreData); rcv.Error() != ErrNoMoreData.Error() {
		t.Errorf("Expecting: %+v, received %+v", ErrNoMoreData.Error(), rcv)
	}
}

func TestRandomInteger(t *testing.T) {
	a := RandomInteger(0, 100)
	b := RandomInteger(0, 100)
	c := RandomInteger(0, 100)
	if a == b && b == c {
		t.Errorf("same result over 3 attempts")
	}
	if a >= 100 || b >= 100 || c >= 100 {
		t.Errorf("one of the numbers equals or it's above the max limit")
	}
	if a < 0 || b < 0 || c < 0 {
		t.Errorf("one of the numbers are below min limit")
	}
}

func TestGetPathIndex(t *testing.T) {
	if rcv, _ := GetPathIndex(EmptyString); rcv != EmptyString {
		t.Errorf("Expecting: \"\"(EmptyString), received: \"%+v\"", rcv)
	}
	eOut := "test"
	if rcv, _ := GetPathIndex("test"); rcv != eOut {
		t.Errorf("Expecting: %+v, received: %+v", eOut, rcv)
	}
	if rcv, index := GetPathIndex("test[10]"); rcv != eOut {
		t.Errorf("Expecting: %+v, received: %+v", eOut, rcv)
	} else if *index != 10 {
		t.Errorf("Expecting: %+v, received: %+v", 10, *index)
	}
	if rcv, index := GetPathIndex("test[0]"); rcv != eOut {
		t.Errorf("Expecting: %+v, received: %+v", eOut, rcv)
	} else if *index != 0 {
		t.Errorf("Expecting: %+v, received: %+v", 0, *index)
	}
	eOut = "test[notanumber]"
	if rcv, index := GetPathIndex("test[notanumber]"); rcv != eOut {
		t.Errorf("Expecting: %+v, received: %+v", eOut, rcv)
	} else if index != nil {
		t.Errorf("Expecting: nil, received: %+v", *index)
	}
	eOut = "test[]"
	if rcv, index := GetPathIndex("test[]"); rcv != eOut {
		t.Errorf("Expecting: %+v, received: %+v", eOut, rcv)
	} else if index != nil {
		t.Errorf("Expecting: nil, received: %+v", *index)
	}
	eOut = "[]"
	if rcv, index := GetPathIndex("[]"); rcv != eOut {
		t.Errorf("Expecting: %+v, received: %+v", eOut, rcv)
	} else if index != nil {
		t.Errorf("Expecting: nil, received: %+v", *index)
	}
}

func TestIsURL(t *testing.T) {
	urls := map[string]bool{
		"/etc/usr/":                           false,
		"https://github.com/cgrates/cgrates/": true,
		"http://github.com/cgrates/cgrates/i": true,
	}
	for url, expected := range urls {
		if rply := IsURL(url); rply != expected {
			t.Errorf("For: %q ,expected %v received: %v", url, expected, rply)
		}
	}
}

func TestComputeHashMatch(t *testing.T) {
	lns, _ := ComputeHash("test1;test2;test3")
	if err := VerifyHash(lns, "test1;test2;test3"); err != true {
		t.Errorf("Expected <true> received: <%v>", err)
	}
}

func TestVerifyHash(t *testing.T) {
	lns, _ := ComputeHash("test1;test2;test3")
	lns2, _ := ComputeHash("test1;test2;test3")
	verify1 := VerifyHash(lns, "test1;test2;test3")
	verify2 := bcrypt.CompareHashAndPassword([]byte(lns2), []byte(ConcatenatedKey("test1;test2;test3")))
	verify3 := false
	if verify2 == nil {
		verify3 = true
	}
	if !reflect.DeepEqual(verify3, verify1) {
		t.Errorf("Expecting: <%+v>, received: <%+v>", verify3, verify1)
	}

}

func TestAESEncryptErrorNil(t *testing.T) {
	encKey := "6368616e676520746869732070617373776f726420746f206120736563726574"
	_, err := AESEncrypt("exampleText", encKey)
	if err != nil {
		t.Errorf("Expecting: <nil>, received: <%+v>", err)
	}
}

func TestAESEncryptError1(t *testing.T) {
	encKey := "1"
	_, err := AESEncrypt("exampleText", encKey)
	if err == nil || err.Error() != "crypto/aes: invalid key size 0" {
		t.Errorf("Expecting error: <crypto/aes: invalid key size 0>, received: <%+v>", err)
	}
}

func TestAESDecryptErrorNil(t *testing.T) {
	encKey := "6368616e676520746869732070617373776f726420746f206120736563726574"
	eString, _ := AESEncrypt("exampleText", encKey)
	_, err := AESDecrypt(eString, encKey)
	if err != nil {
		t.Errorf("Expecting: <nil>, received: <%+v>", err)
	}
}

func TestAESDecryptError1(t *testing.T) {
	encKey := "1"
	eString, _ := AESEncrypt("exampleText", encKey)
	_, err := AESDecrypt(eString, encKey)
	if err == nil || err.Error() != "crypto/aes: invalid key size 0" {
		t.Errorf("Expecting: <crypto/aes: invalid key size 0>, received: <%+v>", err)
	}
}

func TestAESEncryptDecrypt(t *testing.T) {
	encKey := "6368616e676520746869732070617373776f726420746f206120736563726574"
	eString, _ := AESEncrypt("exampleText", encKey)
	dString, _ := AESDecrypt(eString, encKey)
	if !reflect.DeepEqual("exampleText", dString) {
		t.Errorf("Expecting: <exampleText>, received: <%+v>", dString)
	}
}

func TestBoolGenerator(t *testing.T) {
	boolTest := BoolGenerator().RandomBool()
	if boolTest != true && boolTest != false {
		t.Errorf("Needs to be bool")
	}
}

func TestMonthlyEstimated(t *testing.T) {
	t1 := time.Date(2021, 1, 31, 0, 0, 0, 0, time.UTC)
	expectedTime := time.Date(2021, 2, 28, 0, 0, 0, 0, time.UTC)
	if rcv, err := monthlyEstimated(t1); err != nil {
		t.Error(err)
	} else if !reflect.DeepEqual(rcv, expectedTime) {
		t.Errorf("Expected %+v, received %+v", expectedTime, rcv)
	}
}

func TestCoreUtilsGenerateDBItemOpts(t *testing.T) {
	apiKey := "testKey1"
	routeID := "testKey2"
	cache := "testKey3"
	rmtHost := "testKey4"

	expected := map[string]any{
		OptsAPIKey:    apiKey,
		OptsRouteID:   routeID,
		CacheOpt:      cache,
		RemoteHostOpt: rmtHost,
	}
	received := GenerateDBItemOpts(apiKey, routeID, cache, rmtHost)

	if len(received) != len(expected) {
		t.Fatalf("The maps differ in length")
	}
	for key, value := range received {
		if expected[key] != value {
			t.Errorf("\nExpected: <%+v>, \nReceived: <%+v>", expected[key], value)
		}
	}
}

func TestTenantIDConcatenated(t *testing.T) {
	tnt := &TenantID{
		Tenant: "cgrates.org",
		ID:     "1",
	}
	tntOpts := &TenantIDWithAPIOpts{
		TenantID: tnt,
	}
	rcvExpect := "cgrates.org:1"
	concTnt := tntOpts.TenantIDConcatenated()
	if concTnt != rcvExpect {
		t.Errorf("\nExpected: <%+v>, \nReceived: <%+v>", rcvExpect, concTnt)
	}
}

func TestCoreUtilsSplitPathTT(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		sep      byte
		count    int
		expected []string
	}{
		{
			name:     "split with InInFieldSep",
			input:    "*string:~*req.Account:1001",
			sep:      InInFieldSep[0],
			count:    -1,
			expected: []string{"*string", "~*req.Account", "1001"},
		},
		{
			name:     "split with NestingSep",
			input:    "~*req.Account",
			sep:      NestingSep[0],
			count:    -1,
			expected: []string{"~*req", "Account"},
		},
		{
			name:     "split with InInFieldSep and empty last element",
			input:    "*exists:~*vars.*processedProfileIDs[cgrates.org:ATTR]:",
			sep:      InInFieldSep[0],
			count:    -1,
			expected: []string{"*exists", "~*vars.*processedProfileIDs[cgrates.org:ATTR]", ""},
		},
		{
			name:     "split with NestingSep and complex second element",
			input:    "~*vars.*processedProfileIDs[cgrates.org:ATTR]",
			sep:      NestingSep[0],
			count:    -1,
			expected: []string{"~*vars", "*processedProfileIDs[cgrates.org:ATTR]"},
		},
		{
			name:     "split with InInFieldSep and dynamic path",
			input:    "*notexists:~*vars.*processedProfileIDs[<~*vars.apTenantID>]:",
			sep:      InInFieldSep[0],
			count:    -1,
			expected: []string{"*notexists", "~*vars.*processedProfileIDs[<~*vars.apTenantID>]", ""},
		},
		{
			name:     "split with NestingSep and dynamic path",
			input:    "~*vars.*processedProfileIDs[<~*vars.apTenantID>]",
			sep:      NestingSep[0],
			count:    -1,
			expected: []string{"~*vars", "*processedProfileIDs[<~*vars.apTenantID>]"},
		},
		{
			name:     "split with NestingSep (contains '#')",
			input:    "~*stats.SQ_1002.*sum#~*req.Usage",
			sep:      NestingSep[0],
			count:    -1,
			expected: []string{"~*stats", "SQ_1002", "*sum#~*req.Usage"},
		},
		{
			name:     "split with NestingSep (contains '#')",
			input:    "~*stats.SQ_1002.*sum#~*req.CostDetails.Cost",
			sep:      NestingSep[0],
			count:    -1,
			expected: []string{"~*stats", "SQ_1002", "*sum#~*req.CostDetails.Cost"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SplitPath(tt.input, tt.sep, tt.count)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("SplitPath(%q, %q, %d) = %v, want %v",
					tt.input, tt.sep, tt.count, result, tt.expected)
			}
		})
	}
}

func TestCoreUtilsFibSeqNrOverflow(t *testing.T) {
	fib := Fib()
	for i := 0; i < 92; i++ { // the 93rd fibonacci number in the sequence would normally overflow
		fib()
	}
	exp := fib()
	for i := 0; i < 100; i++ {
		if rcv := fib(); rcv != exp {
			t.Errorf("expected: <%+v>, \nreceived: <%+v>", exp, rcv)
		}
	}
}

func TestCoreUtilsFibDurationSeqNrOverflow(t *testing.T) {
	fib := FibDuration(time.Second, 0)
	for i := 0; i < 49; i++ { // the 50th fibonacci number in the sequence would normally overflow when multiplied with time.Second
		fib()
	}
	for i := 0; i < 100; i++ {
		if rcv := fib(); rcv != math.MaxInt {
			t.Errorf("expected: <%+v>, \nreceived: <%+v>", math.MaxInt, rcv)
		}
	}
	fib = FibDuration(time.Second, 6)
	if rcv := fib(); rcv != 6 {
		t.Errorf("expected: <%+v>, \nreceived: <%+v>", math.MaxInt, rcv)
	}
}

func TestCoreUtilsMapStringSlicePointer(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string][]string
		expected *map[string][]string
	}{
		{
			name:     "Empty map",
			input:    map[string][]string{},
			expected: &map[string][]string{},
		},
		{
			name: "Non-empty map",
			input: map[string][]string{
				"key1": {"value1", "value2"},
				"key2": {"value3"},
			},
			expected: &map[string][]string{
				"key1": {"value1", "value2"},
				"key2": {"value3"},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := MapStringSlicePointer(test.input)
			if diff := cmp.Diff(result, test.expected); diff != "" {
				t.Errorf("Test case %q failed: (-got +want)\n%s", test.name, diff)
			}
		})
	}
}

func TestCoreutilsParseHierarchyPath(t *testing.T) {
	tests := []struct {
		path     string
		sep      string
		expected HierarchyPath
	}{
		{
			path:     "",
			sep:      "/",
			expected: nil,
		},
		{
			path:     "parent/child/grandchild",
			sep:      "",
			expected: HierarchyPath{"parent", "child", "grandchild"},
		},
		{
			path:     "root-section-subsection",
			sep:      "-",
			expected: HierarchyPath{"root", "section", "subsection"},
		},
	}
	for _, test := range tests {
		result := ParseHierarchyPath(test.path, test.sep)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("For path: %s and sep: %s, expected: %v, but got: %v", test.path, test.sep, test.expected, result)
		}
	}
}

func TestMonthlyEstimatedCases(t *testing.T) {
	tests := []struct {
		name     string
		t1       time.Time
		expected time.Time
	}{
		{
			name:     "Non-Leap Year: January 31 to February 28",
			t1:       time.Date(2021, 1, 31, 0, 0, 0, 0, time.UTC),
			expected: time.Date(2021, 2, 28, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "Leap Year: January 29 to February 29",
			t1:       time.Date(2020, 1, 29, 0, 0, 0, 0, time.UTC),
			expected: time.Date(2020, 2, 29, 0, 0, 0, 0, time.UTC),
		},

		{
			name:     "Non-Leap Year: February 28 to March 28",
			t1:       time.Date(2021, 2, 28, 0, 0, 0, 0, time.UTC),
			expected: time.Date(2021, 3, 28, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "October 31 to November 30",
			t1:       time.Date(2021, 10, 31, 0, 0, 0, 0, time.UTC),
			expected: time.Date(2021, 11, 30, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "November 30 to December 30",
			t1:       time.Date(2021, 11, 30, 0, 0, 0, 0, time.UTC),
			expected: time.Date(2021, 12, 30, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "November 25 to December 25",
			t1:       time.Date(2021, 11, 25, 0, 0, 0, 0, time.UTC),
			expected: time.Date(2021, 12, 25, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := monthlyEstimated(tt.t1)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if !got.Equal(tt.expected) {
				t.Errorf("Expected: %v, got: %v", tt.expected, got)
			}
		})
	}
}

func TestGetEndOfMonth(t *testing.T) {
	tests := []struct {
		name     string
		ref      time.Time
		expected time.Time
	}{
		{
			name:     "End of January (non-leap year)",
			ref:      time.Date(2023, time.January, 15, 0, 0, 0, 0, time.UTC),
			expected: time.Date(2023, time.January, 31, 23, 59, 59, 0, time.UTC),
		},
		{
			name:     "End of February (non-leap year)",
			ref:      time.Date(2023, time.February, 10, 0, 0, 0, 0, time.UTC),
			expected: time.Date(2023, time.February, 28, 23, 59, 59, 0, time.UTC),
		},
		{
			name:     "End of February (leap year)",
			ref:      time.Date(2024, time.February, 10, 0, 0, 0, 0, time.UTC),
			expected: time.Date(2024, time.February, 29, 23, 59, 59, 0, time.UTC),
		},
		{
			name:     "End of December",
			ref:      time.Date(2023, time.December, 25, 0, 0, 0, 0, time.UTC),
			expected: time.Date(2023, time.December, 31, 23, 59, 59, 0, time.UTC),
		},
		{
			name:     "End of April ",
			ref:      time.Date(2023, time.April, 15, 0, 0, 0, 0, time.UTC),
			expected: time.Date(2023, time.April, 30, 23, 59, 59, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetEndOfMonth(tt.ref)
			if !result.Equal(tt.expected) {
				t.Errorf("GetEndOfMonth(%v) = %v, want %v", tt.ref, result, tt.expected)
			}
		})
	}
}

func TestParseBinarySize(t *testing.T) {
	type args struct {
		size string
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
		expErr  error
	}{
		{
			name:    "OK1GB",
			args:    args{size: "1GB"},
			want:    int64(1073741824),
			wantErr: false,
		},
		{
			name:    "OK1MB",
			args:    args{size: "1MB"},
			want:    int64(1048576),
			wantErr: false,
		},
		{
			name:    "OK1KB",
			args:    args{size: "1KB"},
			want:    int64(1024),
			wantErr: false,
		},
		{
			name:    "OK1B",
			args:    args{size: "1B"},
			want:    int64(1),
			wantErr: false,
		},
		{
			name:    "Empty",
			args:    args{size: ""},
			want:    int64(0),
			wantErr: true,
			expErr:  fmt.Errorf("invalid size format: "),
		},
		{
			name:    "badSize",
			args:    args{size: "badSize"},
			want:    int64(0),
			wantErr: true,
			expErr:  fmt.Errorf("invalid size format: badSize"),
		},
		{
			name:    "badUnit",
			args:    args{size: "1badunit"},
			want:    int64(0),
			wantErr: true,
			expErr:  fmt.Errorf("unknown unit: badunit"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseBinarySize(tt.args.size)
			if (err != nil) != tt.wantErr {
				t.Error(err)
				return
			}
			if tt.expErr != nil && (err == nil || err.Error() != tt.expErr.Error()) {
				t.Errorf("ParseBinarySize() error = <%v>, expErr <%v>", err, tt.expErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseBinarySize() got <%v>, want <%v>", got, tt.want)
			}
		})
	}
}
