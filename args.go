package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

const (
	appName   = "selcvsdt"
	usagejson = "output datetime params json file."
	usage     = `%s selects lines from a csv file by specified range of date and time items.

Default columns on the csv file is the date column is 0 and the time column is 1.

An output file name is added a dot and a sequence number after the csv file name.

Usage:

	%s -json
	%s <csvfile> <basetime> <incdays> [<outpath>]

Each arguments are:

	-json
		%s
	<csvfile>
		a reading csv file.
	<basetime>
		a basetime (YYYYMMDDHHmm).
	<incdays>
		add days to the basetime.
	[<outpath>]
		an output file path. default is ".".

`
)

// Args struct.
type args struct {
	Csvfile  string
	Basetime time.Time
	Incdays  int
	Outpath  string
}

func (args *args) validArgs() error {
	if ln := len(os.Args); ln != 4 && ln != 5 {
		fmt.Fprintf(os.Stderr, usage, appName, appName, appName, usagejson)
		return nil
	}

	csvfile := os.Args[1]
	if fi, err := os.Stat(csvfile); err != nil || fi.IsDir() {
		return err
	}

	basetime, err := time.ParseInLocation(params.Fmtarg, os.Args[2], time.Local)
	if err != nil {
		return err
	}

	incdays, err := strconv.Atoi(os.Args[3])
	if err != nil {
		return err
	}

	outpath := "."
	if len(os.Args) == 5 {
		outpath = os.Args[4]
		if fi, err := os.Stat(outpath); err != nil || !fi.IsDir() {
			return err
		}
	}

	args.Csvfile = csvfile
	args.Basetime = basetime
	args.Incdays = incdays
	args.Outpath = outpath

	return nil
}

func (args *args) inctime() time.Time {
	return args.Basetime.AddDate(0, 0, args.Incdays)
}
