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

package engine

import (
	"errors"
	"slices"
	"time"

	"github.com/cgrates/birpc/context"
	"github.com/cgrates/cgrates/utils"
	"github.com/ericlagergren/decimal"
)

// GetFloat64Opts checks the specified option names in order among the keys in APIOpts returning the first value it finds as float64, otherwise it
// returns the config option if at least one filter passes or the default value if none of them do
func GetFloat64Opts(ctx *context.Context, tnt string, ev *utils.CGREvent, fS *FilterS, dynOpts []*utils.DynamicFloat64Opt,
	dftOpt float64, optNames ...string) (cfgOpt float64, err error) {
	for _, optName := range optNames {
		if opt, has := ev.APIOpts[optName]; has {
			return utils.IfaceAsFloat64(opt)
		}
	}
	evDP := ev.AsDataProvider()
	for _, opt := range dynOpts { // iterate through the options
		if !slices.Contains([]string{utils.EmptyString, utils.MetaAny, tnt}, opt.Tenant) {
			continue
		}
		if pass, err := fS.Pass(ctx, tnt, opt.FilterIDs, evDP); err != nil { // check if the filter is passing for the DataProvider and return the option if it does
			return 0, err
		} else if pass {
			return opt.Value, nil
		}
	}
	return dftOpt, nil // return the default value if there are no options and none of the filters pass
}

// GetDurationOpts checks the specified option names in order among the keys in APIOpts returning the first value it finds as time.Duration, otherwise it
// returns the config option if at least one filter passes or the default value if none of them do
func GetDurationOpts(ctx *context.Context, tnt string, ev *utils.CGREvent, fS *FilterS, dynOpts []*utils.DynamicDurationOpt,
	dftOpt time.Duration, optNames ...string) (cfgOpt time.Duration, err error) {
	for _, optName := range optNames {
		if opt, has := ev.APIOpts[optName]; has {
			return utils.IfaceAsDuration(opt)
		}
	}
	evDP := ev.AsDataProvider()
	for _, opt := range dynOpts { // iterate through the options
		if !slices.Contains([]string{utils.EmptyString, utils.MetaAny, tnt}, opt.Tenant) {
			continue
		}
		if pass, err := fS.Pass(ctx, tnt, opt.FilterIDs, evDP); err != nil { // check if the filter is passing for the DataProvider and return the option if it does
			return 0, err
		} else if pass {
			return opt.Value, nil
		}
	}
	return dftOpt, nil // return the default value if there are no options and none of the filters pass
}

// GetStringOpts checks the specified option names in order among the keys in APIOpts returning the first value it finds as string, otherwise it
// returns the config option if at least one filter passes or the default value if none of them do
func GetStringOpts(ctx *context.Context, tnt string, ev *utils.CGREvent, fS *FilterS, dynOpts []*utils.DynamicStringOpt,
	dftOpt string, optNames ...string) (cfgOpt string, err error) {
	for _, optName := range optNames {
		if opt, has := ev.APIOpts[optName]; has {
			return utils.IfaceAsString(opt), nil
		}
	}
	evDP := ev.AsDataProvider()
	for _, opt := range dynOpts { // iterate through the options
		if !slices.Contains([]string{utils.EmptyString, utils.MetaAny, tnt}, opt.Tenant) {
			continue
		}
		if pass, err := fS.Pass(ctx, tnt, opt.FilterIDs, evDP); err != nil { // check if the filter is passing for the DataProvider and return the option if it does
			return utils.EmptyString, err
		} else if pass {
			return opt.Value, nil
		}
	}
	return dftOpt, nil // return the default value if there are no options and none of the filters pass
}

// GetTimeOpts checks the specified option names in order among the keys in APIOpts returning the first value it finds as time.Time, otherwise it
// returns the config option if at least one filter passes or the default value if none of them do
func GetTimeOpts(ctx *context.Context, tnt string, ev *utils.CGREvent, fS *FilterS, dynOpts []*utils.DynamicStringOpt,
	tmz string, dftOpt string, optNames ...string) (_ time.Time, err error) {
	for _, optName := range optNames {
		if opt, has := ev.APIOpts[optName]; has {
			return utils.IfaceAsTime(opt, tmz)
		}
	}
	evDP := ev.AsDataProvider()
	for _, opt := range dynOpts { // iterate through the options
		if !slices.Contains([]string{utils.EmptyString, utils.MetaAny, tnt}, opt.Tenant) {
			continue
		}
		var pass bool
		if pass, err = fS.Pass(ctx, tnt, opt.FilterIDs, evDP); err != nil { // check if the filter is passing for the DataProvider and return the option if it does
			return
		} else if pass {
			return utils.ParseTimeDetectLayout(opt.Value, tmz)
		}
	}
	return utils.ParseTimeDetectLayout(dftOpt, tmz) // return the default value if there are no options and none of the filters pass
}

// GetStringSliceOpts checks the specified option names in order among the keys in APIOpts returning the first value it finds as []string, otherwise it
// returns the config option if at least one filter passes or the default value if none of them do
func GetStringSliceOpts(ctx *context.Context, tnt string, ev *utils.CGREvent, fS *FilterS, dynOpts []*utils.DynamicStringSliceOpt,
	dftOpt []string, optNames ...string) (cfgOpt []string, err error) {
	for _, optName := range optNames {
		if opt, has := ev.APIOpts[optName]; has {
			return utils.IfaceAsStringSlice(opt)
		}
	}
	evDP := ev.AsDataProvider()
	for _, opt := range dynOpts { // iterate through the options
		if !slices.Contains([]string{utils.EmptyString, utils.MetaAny, tnt}, opt.Tenant) {
			continue
		}
		if pass, err := fS.Pass(ctx, tnt, opt.FilterIDs, evDP); err != nil { // check if the filter is passing for the DataProvider and return the option if it does
			return nil, err
		} else if pass {
			return opt.Value, nil
		}
	}
	return dftOpt, nil // return the default value if there are no options and none of the filters pass
}

// GetIntOpts checks the specified option names in order among the keys in APIOpts returning the first value it finds as int, otherwise it
// returns the config option if at least one filter passes or the default value if none of them do
func GetIntOpts(ctx *context.Context, tnt string, ev *utils.CGREvent, fS *FilterS, dynOpts []*utils.DynamicIntOpt,
	dftOpt int, optNames ...string) (cfgOpt int, err error) {
	for _, optName := range optNames {
		if opt, has := ev.APIOpts[optName]; has {
			var value int64
			if value, err = utils.IfaceAsTInt64(opt); err != nil {
				return 0, err
			}
			return int(value), nil
		}
	}
	evDP := ev.AsDataProvider()
	for _, opt := range dynOpts { // iterate through the options
		if !slices.Contains([]string{utils.EmptyString, utils.MetaAny, tnt}, opt.Tenant) {
			continue
		}
		if pass, err := fS.Pass(ctx, tnt, opt.FilterIDs, evDP); err != nil { // check if the filter is passing for the DataProvider and return the option if it does
			return 0, err
		} else if pass {
			return opt.Value, nil
		}
	}
	return dftOpt, nil // return the default value if there are no options and none of the filters pass
}

// GetBoolOpts checks the specified option names in order among the keys in APIOpts returning the first value it finds as bool, otherwise it
// returns the config option if at least one filter passes or the default value if none of them do
func GetBoolOpts(ctx *context.Context, tnt string, dP utils.DataProvider, fS *FilterS, dynOpts []*utils.DynamicBoolOpt,
	dftOpt bool, optNames ...string) (cfgOpt bool, err error) {
	values, err := dP.FieldAsInterface([]string{utils.MetaOpts})
	if err != nil {
		return false, err
	}
	opts, err := ConvertOptsToMapStringAny(values)
	if err != nil {
		return
	}
	for _, optName := range optNames {
		if opt, has := opts[optName]; has {
			return utils.IfaceAsBool(opt)
		}
	}
	for _, opt := range dynOpts { // iterate through the options
		if !slices.Contains([]string{utils.EmptyString, utils.MetaAny, tnt}, opt.Tenant) {
			continue
		}
		if pass, err := fS.Pass(ctx, tnt, opt.FilterIDs, dP); err != nil { // check if the filter is passing for the DataProvider and return the option if it does
			return false, err
		} else if pass {
			return opt.Value, nil
		}
	}
	return dftOpt, nil // return the default value if there are no options and none of the filters pass
}

// GetDecimalBigOpts checks the specified option names in order among the keys in APIOpts returning the first value it finds as *decimal.Big, otherwise it
// returns the config option if at least one filter passes or the default value if none of them do
func GetDecimalBigOpts(ctx *context.Context, tnt string, ev *utils.CGREvent, fS *FilterS, dynOpts []*utils.DynamicDecimalBigOpt,
	dftOpt *decimal.Big, optNames ...string) (cfgOpt *decimal.Big, err error) {
	for _, optName := range optNames {
		if opt, has := ev.APIOpts[optName]; has {
			return utils.IfaceAsBig(opt)
		}
	}
	evDP := ev.AsDataProvider()
	for _, opt := range dynOpts { // iterate through the options
		if !slices.Contains([]string{utils.EmptyString, utils.MetaAny, tnt}, opt.Tenant) {
			continue
		}
		if pass, err := fS.Pass(ctx, tnt, opt.FilterIDs, evDP); err != nil { // check if the filter is passing for the DataProvider and return the option if it does
			return nil, err
		} else if pass {
			return opt.Value, nil
		}
	}
	return dftOpt, nil // return the default value if there are no options and none of the filters pass
}

// GetInterfaceOpts checks the specified option names in order among the keys in APIOpts returning the first value it finds as any, otherwise it
// returns the config option if at least one filter passes or the default value if none of them do
func GetInterfaceOpts(ctx *context.Context, tnt string, ev *utils.CGREvent, fS *FilterS, dynOpts []*utils.DynamicInterfaceOpt,
	dftOpt any, optNames ...string) (cfgOpt any, err error) {
	for _, optName := range optNames {
		if opt, has := ev.APIOpts[optName]; has {
			return opt, nil
		}
	}
	evDP := ev.AsDataProvider()
	for _, opt := range dynOpts { // iterate through the options
		if !slices.Contains([]string{utils.EmptyString, utils.MetaAny, tnt}, opt.Tenant) {
			continue
		}
		if pass, err := fS.Pass(ctx, tnt, opt.FilterIDs, evDP); err != nil { // check if the filter is passing for the DataProvider and return the option if it does
			return false, err
		} else if pass {
			return opt.Value, nil
		}
	}
	return dftOpt, nil // return the default value if there are no options and none of the filters pass
}

// GetIntPointerOpts checks the specified option names in order among the keys in APIOpts returning the first value it finds as *int, otherwise it
// returns the config option if at least one filter passes or NOT_FOUND if none of them do
func GetIntPointerOpts(ctx *context.Context, tnt string, ev *utils.CGREvent, fS *FilterS, dynOpts []*utils.DynamicIntPointerOpt,
	optNames ...string) (cfgOpt *int, err error) {
	for _, optName := range optNames {
		if opt, has := ev.APIOpts[optName]; has {
			var value int64
			if value, err = utils.IfaceAsTInt64(opt); err != nil {
				return nil, err
			}
			return utils.IntPointer(int(value)), nil
		}
	}
	evDP := ev.AsDataProvider()
	for _, opt := range dynOpts { // iterate through the options
		if !slices.Contains([]string{utils.EmptyString, utils.MetaAny, tnt}, opt.Tenant) {
			continue
		}
		if pass, err := fS.Pass(ctx, tnt, opt.FilterIDs, evDP); err != nil { // check if the filter is passing for the DataProvider and return the option if it does
			return nil, err
		} else if pass {
			return opt.Value, nil
		}
	}
	return nil, nil
}

// GetDurationPointerOptsFromMultipleMaps checks the specified option names in order among the keys in APIOpts, then in startOpts, returning the first value it finds as *time.Duration,
// otherwise it returns the config option if at least one filter passes or NOT_FOUND if none of them do
func GetDurationPointerOptsFromMultipleMaps(ctx *context.Context, tnt string, eventStart, apiOpts, startOpts map[string]any, fS *FilterS,
	dynOpts []*utils.DynamicDurationPointerOpt, optName string) (cfgOpt *time.Duration, err error) {
	var value time.Duration
	if opt, has := apiOpts[optName]; has {
		if value, err = utils.IfaceAsDuration(opt); err != nil {
			return nil, err
		}
		return utils.DurationPointer(value), nil
	} else if opt, has = startOpts[optName]; has {
		if value, err = utils.IfaceAsDuration(opt); err != nil {
			return nil, err
		}
		return utils.DurationPointer(value), nil
	}
	evMS := utils.MapStorage{
		utils.MetaOpts: apiOpts,
		utils.MetaReq:  eventStart,
	}
	for _, opt := range dynOpts { // iterate through the options
		if !slices.Contains([]string{utils.EmptyString, utils.MetaAny, tnt}, opt.Tenant) {
			continue
		}
		if pass, err := fS.Pass(ctx, tnt, opt.FilterIDs, evMS); err != nil { // check if the filter is passing for the DataProvider and return the option if it does
			return nil, err
		} else if pass {
			return opt.Value, nil
		}
	}
	return nil, nil
}

// GetDurationOptsFromMultipleMaps checks the specified option names in order among the keys in APIOpts, then in startOpts, returning the first value it finds as time.Duration,
// otherwise it returns the config option if at least one filter passes or the default one if none of them do
func GetDurationOptsFromMultipleMaps(ctx *context.Context, tnt string, eventStart, apiOpts, startOpts map[string]any, fS *FilterS, dynOpts []*utils.DynamicDurationOpt,
	dftOpt time.Duration, optName string) (cfgOpt time.Duration, err error) {
	var value time.Duration
	if opt, has := apiOpts[optName]; has {
		if value, err = utils.IfaceAsDuration(opt); err != nil {
			return 0, err
		}
		return value, nil
	} else if opt, has = startOpts[optName]; has {
		if value, err = utils.IfaceAsDuration(opt); err != nil {
			return 0, err
		}
		return value, nil
	}
	evMS := utils.MapStorage{
		utils.MetaOpts: apiOpts,
		utils.MetaReq:  eventStart,
	}
	for _, opt := range dynOpts { // iterate through the options
		if !slices.Contains([]string{utils.EmptyString, utils.MetaAny, tnt}, opt.Tenant) {
			continue
		}
		if pass, err := fS.Pass(ctx, tnt, opt.FilterIDs, evMS); err != nil { // check if the filter is passing for the DataProvider and return the option if it does
			return 0, err
		} else if pass {
			return opt.Value, nil
		}
	}
	return dftOpt, nil // return the default value if there are no options and none of the filters pass
}

func ConvertOptsToMapStringAny(in any) (map[string]any, error) {
	out := make(map[string]any)
	switch val := in.(type) {
	case MapEvent:
		for k, v := range val {
			out[k] = v
		}
	case utils.MapStorage:
		for k, v := range val {
			out[k] = v
		}
	case map[string]any:
		return val, nil
	default:
		return nil, errors.New("cannot convert to map[string]any")
	}
	return out, nil
}