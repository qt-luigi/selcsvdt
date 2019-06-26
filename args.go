package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"
)

const (
	usageJSON = "output datetime params json file."
	usageMsg  = `selcsvdt selects lines from a csv file by specified range of date and time items.

Default columns on the csv file is the date column is 0 and the time column is 1.

When selcsvdt.json file exists in current directory, selcvsdt read it when start.

Usage:

	selcsvdt -json
	selcsvdt <csvpath> <basetime> <incdays>

Each arguments are:

	-json
		%s

	<csvpath>
		a reading csv file.
	<basetime>
		a basetime. default format is YYYYMMDDHHmm.
		To change fmtarg item in json file.
	<incdays>
		add days to the basetime.

When you execute with -json switch, A parameter json file is created in current directory. 

{
	"fmtarg": "200601021504",
	"date": {
		"column": 0,
		"format": "2006-01-02"
	},
	"time": {
		"column": 1,
		"format": "15:04"
	}
}
`
)

var (
	// ErrUsage represents usage.
	ErrUsage = errors.New("usage")
)

func usage() {
	fmt.Fprintf(os.Stderr, usageMsg, usageJSON)
}

// Args struct.
type args struct {
	csvpath  string
	basetime time.Time
	incdays  int
}

func (args *args) validArgs(fmtArg string) error {
	if len(os.Args) != 4 {
		return ErrUsage
	}

	csvpath := os.Args[1]
	if _, err := os.Stat(csvpath); err != nil {
		return err
	}

	basetime, err := time.ParseInLocation(fmtArg, os.Args[2], time.Local)
	if err != nil {
		return err
	}

	incdays, err := strconv.Atoi(os.Args[3])
	if err != nil {
		return err
	}

	args.csvpath = csvpath
	args.basetime = basetime
	args.incdays = incdays

	return nil
}

func (args *args) inctime() time.Time {
	return args.basetime.AddDate(0, 0, args.incdays)
}
