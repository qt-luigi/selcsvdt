package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
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

	lines := make([]string, 0)

	columns := columns{params.Date.Column, params.Time.Column}
	dt := datetime(params.Date.Format, params.Time.Format)
	period := sort(args.basetime, args.inctime(), dt)
	lines, err := selectRows(lines, args.csvfile, columns, period)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	outfile := outName(filepath.Base(args.csvfile), 1)
	if err := output(args.outpath, outfile, lines, rtncd()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
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

func selectRows(lines []string, csvfile string, cols columns, period period) ([]string, error) {
	f, err := os.Open(csvfile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

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

func outName(basename string, seq int) string {
	return basename + "." + strconv.Itoa(seq)
}

func rtncd() string {
	if runtime.GOOS == "windows" {
		return "\r\n"
	}
	return "\n"
}

func output(path, file string, datas []string, rtncd string) error {
	f, err := os.Create(filepath.Join(path, file))
	if err != nil {
		return err
	}
	defer f.Close()

	for _, data := range datas {
		if _, err := f.WriteString(data + rtncd); err != nil {
			return err
		}
	}

	return nil
}
