# selcsvdt

selcvsdt selects lines from a csv file by specified range of date and time items.

## Installation

When you have installed the Go, Please execute following `go get` command:

```sh
go get -u github.com/qt-luigi/selcsvdt
```

## Usage

```sh
$ selcvsdt
selcvsdt selects lines from a csv file by specified range of date and time items.

Default columns on the csv file is the date column is 0 and the time column is 1.

When selcsvdt.json file exists in current directory, selcvsdt read it when start.

An output file name is added a dot and a sequence number after the csv file name.

Usage:

	selcvsdt -json
	selcvsdt <csvfile> <basetime> <incdays> [<outpath>]

Each arguments are:

	-json
		output datetime params json file.
	<csvfile>
		a reading csv file.
	<basetime>
		a basetime (YYYYMMDDHHmm).
	<incdays>
		add days to the basetime.
	[<outpath>]
		an output file path. default is ".".

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
```

See the test.csv file.

```sh
$ cat sample/test.csv
"2018-05-28","21:57",1,"foo"
"2018-05-29","20:21",2,"bar"
"2018-05-30","13:46",3,"baz"
```

If you specifiy basetime is 201805292000 (a format is defined by "fmtarg") and incdays is 1, datetime range is 201805292000 - 201805302000.

```sh
$ selcsvdt sample/test.csv 201805292000 1
$ cat test.csv.1 
"2018-05-29","20:21",2,"bar"
"2018-05-30","13:46",3,"baz"
```

If incdays is -1, datetime range is 201805282000 - 201805292000.

```sh
$ selcsvdt sample/test.csv 201805292000 -1
$ cat test.csv.1 
"2018-05-28","21:57",1,"foo"
```

## License

MIT

## Author

Ryuji Iwata

## Note

This tool is mainly using by myself. :-)
