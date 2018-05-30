package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
)

var params Params

func init() {
	params = Params{
		Fmtarg: "200601021504",
		Date:   Item{Column: 0, Format: "2006-01-02"},
		Time:   Item{Column: 1, Format: "15:04"},
	}
}

// Item struct.
type Item struct {
	Column int    `json:"column"`
	Format string `json:"format"`
}

// Params struct.
type Params struct {
	Fmtarg string `json:"fmtarg"`
	Date   Item   `json:"date"`
	Time   Item   `json:"time"`
}

func (params *Params) readJSON(file string) error {
	if b, err := ioutil.ReadFile(file); err != nil {
		return err
	} else if err := json.Unmarshal(b, params); err != nil {
		return err
	} else {
		return nil
	}
}

func (params *Params) writeJSON(file string) error {
	b, err := json.Marshal(params)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	err = json.Indent(&buf, b, "", "    ")
	if err != nil {
		return err
	}

	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.Write(buf.Bytes()); err != nil {
		return err
	}

	return nil
}
