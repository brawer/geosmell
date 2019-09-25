// SPDX-FileCopyrightText: 2019 Sascha Brawer <sascha@brawer.ch>
// SPDX-License-Identifier: MIT

package main

import (
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}

func processDataset(dataset Dataset, s2Level int) (string, error) {
	tempDir, err := ioutil.TempDir("", "geosmell-test")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(tempDir)

	filePath := path.Join(tempDir, "out.gz")
	if err := dataset.Process(s2Level, filePath); err != nil {
		return "", err
	}

	stream, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer stream.Close()

	gzstream, err := gzip.NewReader(stream)
	if err != nil {
		return "", err
	}
	defer gzstream.Close()

	contentBytes, err := ioutil.ReadAll(gzstream)
	if err != nil {
		return "", err
	}
	return string(contentBytes), nil
}
