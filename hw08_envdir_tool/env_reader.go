package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"sync"
)

type Environment map[string]string

type Data struct {
	env Environment
	mx  sync.Mutex
}

func readFileLine(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	line, err := reader.ReadBytes('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return "", err
	}

	// replace \n, when file contains multiple lines
	line = bytes.ReplaceAll(line, []byte("\n"), []byte(""))
	line = bytes.ReplaceAll(line, []byte("\x00"), []byte("\n"))

	return strings.TrimRight(string(line), " \t"), nil
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	var wg sync.WaitGroup
	var data Data
	data.env = Environment{}
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, fileInfo := range files {
		if !fileInfo.Mode().IsRegular() {
			continue
		}
		wg.Add(1)
		go func(fileInfo os.FileInfo) {
			var (
				line string
				// Nice bug - if I try to use "err" without redeclaring, this variable will be used from
				// upper scope and race condition happens
				err error
			)
			defer wg.Done()
			fileName := fileInfo.Name()
			if fileInfo.Size() > 0 {
				line, err = readFileLine(path.Join(dir, fileName))
				if err != nil {
					fmt.Printf("Couldn't read data from %s\n", fileName)

					return
				}
			}
			data.mx.Lock()
			data.env[fileName] = line
			data.mx.Unlock()
		}(fileInfo)
	}
	wg.Wait()

	return data.env, nil
}
