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

var jsn bool

func init() {
	flag.BoolVar(&jsn, "json", false, usagejson)
}

func main() {
	flag.Parse()

	jsonfile := filepath.Base(os.Args[0]) + ".json"

	if jsn {
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
	if err := args.validArgs(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	layout := datetime{params.Date.Format, params.Time.Format}.String()
	period := sort(layout, args.Basetime, args.inctime())

	lines := make([]string, 0)

	columns := columns{params.Date.Column, params.Time.Column}
	lines, err := selectRows(lines, args.Csvfile, columns, period)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	outfile := outName(filepath.Base(args.Csvfile), 1)
	if err := output(args.Outpath, outfile, lines, rtncd()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// Period struct.
type period struct {
	from string
	to   string
}

// Columns struct.
type columns struct {
	date int
	time int
}

// Datetime struct.
type datetime struct {
	date string
	time string
}

func (dt datetime) String() string {
	return fmt.Sprintf("%s %s", dt.date, dt.time)
}

func sort(layout string, t1, t2 time.Time) period {
	s1 := t1.Format(layout)
	s2 := t2.Format(layout)
	if s1 > s2 {
		s1, s2 = s2, s1
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

	dt := datetime{rcd[cols.date], rcd[cols.time]}.String()
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
