package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var flagJSON bool

func init() {
	flag.BoolVar(&flagJSON, "json", false, usageJSON)
}

func main() {
	flag.Parse()

	jsonfile := jsonName()

	if flagJSON {
		if len(os.Args) != 2 {
			usage()
			os.Exit(2)
		}
		if err := params.writeJSON(jsonfile); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
		os.Exit(0)
	}

	if fi, err := os.Stat(jsonfile); err == nil && !fi.IsDir() {
		if err := params.readJSON(jsonfile); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
	}

	var args args
	if err := args.validArgs(params.Fmtarg); err != nil {
		if err == ErrUsage {
			usage()
		} else {
			fmt.Fprintln(os.Stderr, err)
		}
		os.Exit(2)
	}

	columns := columns{params.Date.Column, params.Time.Column}
	dt := datetime(params.Date.Format, params.Time.Format)
	period := sort(args.basetime, args.inctime(), dt)
	if lines, err := selectRows(args.csvpath, columns, period); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	} else {
		for _, line := range lines {
			fmt.Println(line)
		}
	}
}

func jsonName() string {
	base := filepath.Base(os.Args[0])
	return base[:len(base)-len(filepath.Ext(base))] + ".json"
}

func datetime(date, time string) string {
	return fmt.Sprintf("%s %s", date, time)
}

type columns struct {
	date int
	time int
}

type period struct {
	from string
	to   string
}

func sort(t1, t2 time.Time, layout string) period {
	s1 := t1.Format(layout)
	s2 := t2.Format(layout)
	if s1 > s2 {
		return period{s2, s1}
	}
	return period{s1, s2}
}

func selectRows(csvpath string, cols columns, period period) ([]string, error) {
	if fi, err := os.Stat(csvpath); err != nil {
		return nil, err
	} else if !fi.IsDir() {
		return read(csvpath, cols, period)
	} else {
		fis, err := ioutil.ReadDir(csvpath)
		if err != nil {
			return nil, err
		}
		lines := make([]string, 0)
		for _, fi := range fis {
			if !fi.IsDir() {
				fn := filepath.Join(csvpath, fi.Name())
				if ss, err := read(fn, cols, period); err != nil {
					return nil, err
				} else {
					lines = append(lines, ss...)
				}
			}
		}
		return lines, nil
	}
}

func read(filename string, cols columns, period period) ([]string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	lines := make([]string, 0)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		text := scanner.Text()
		if ok, err := target(text, cols, period); err != nil {
			return nil, err
		} else if ok {
			lines = append(lines, text)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

func target(text string, cols columns, period period) (bool, error) {
	rcd, err := csv.NewReader(strings.NewReader(text)).Read()
	if err == io.EOF {
		return false, nil
	} else if err != nil {
		return false, err
	}

	if ln := len(rcd); ln == 0 {
		return false, nil
	} else if ln <= cols.date || ln <= cols.time {
		return false, nil
	}

	dt := datetime(rcd[cols.date], rcd[cols.time])
	if dt < period.from || period.to < dt {
		return false, nil
	}

	return true, nil
}
