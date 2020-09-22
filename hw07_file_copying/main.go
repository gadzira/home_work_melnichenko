package main

import (
	"flag"
	"fmt"

	"github.com/gadzira/home_work_melnichenko/hw07_file_copying/copy"
)

var (
	from, to      string
	limit, offset int64
)

func init() {
	flag.StringVar(&from, "from", "", "file to read from")
	flag.StringVar(&to, "to", "", "file to write to")
	flag.Int64Var(&limit, "limit", 0, "limit of bytes to copy")
	flag.Int64Var(&offset, "offset", 0, "offset in input file")
}

func main() {
	flag.Parse()
	fromPath := from
	toPath := to
	offset := offset
	limit := limit

	err := copy.Copy(fromPath, toPath, offset, limit)
	if err != nil {
		fmt.Println(err)
	}
}
