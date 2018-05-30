package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"
)

const (
	appName   = "selcvsdt"
	usageJSON = "output datetime params json file."
	usageMsg  = `%s selects lines from a csv file by specified range of date and time items.

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

var (
	// ErrUsage represents usage.
	ErrUsage = errors.New("usage")
)

func usage() {
	fmt.Fprintf(os.Stderr, usageMsg, appName, appName, appName, usageJSON)
}

// Args struct.
type args struct {
	csvfile  string
	basetime time.Time
	incdays  int
	outpath  string
}

func (args *args) validArgs(fmtArg string) error {
	if ln := len(os.Args); ln != 4 && ln != 5 {
		return ErrUsage
	}

	csvfile := os.Args[1]
	if fi, err := os.Stat(csvfile); err != nil || fi.IsDir() {
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

	outpath := "."
	if len(os.Args) == 5 {
		outpath = os.Args[4]
		if fi, err := os.Stat(outpath); err != nil || !fi.IsDir() {
			return err
		}
	}

	args.csvfile = csvfile
	args.basetime = basetime
	args.incdays = incdays
	args.outpath = outpath

	return nil
}

func (args *args) inctime() time.Time {
	return args.basetime.AddDate(0, 0, args.incdays)
}
