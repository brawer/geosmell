// SPDX-FileCopyrightText: 2019 Sascha Brawer <sascha@brawer.ch>
// SPDX-License-Identifier: MIT

package main

import (
	//"fmt"
	"io/ioutil"
	"net/http"
	"os"
	//"path/filepath"
	//"reflect"
	//"runtime"
	"strings"
	"testing"
)

type fakeCHStatPopServer struct{}

func (s fakeCHStatPopServer) Open(path string) (f http.File, e error) {
	if strings.HasSuffix(path, "dynamiclist.html") {
		return os.Open("testdata/chstatpop/list.html")
	} else if strings.HasSuffix(path, "9606372.html") {
		return os.Open("testdata/chstatpop/9606372.html")
	} else if strings.HasSuffix(path, "/assets/9606372/master") {
		return os.Open("testdata/chstatpop/statpop2018.zip")
	} else {
		return os.Open("testdata/chstatpop/notfound.html")
	}
}

func TestFindLatestCHStatPop(t *testing.T) {
	fs := &fakeCHStatPopServer{}
	client := &http.Client{Transport: http.NewFileTransport(fs)}
	if url, timestamp, err := findLatestCHStatPop(client); err == nil {
		equals(t, "https://www.bfs.admin.ch/bfsstatic/dam/assets/9606372/master", url)
		equals(t, "2019-08-27", timestamp.String()[:10])
	} else {
		t.Error(err)
	}
}

func TestFetchCHStatPop(t *testing.T) {
	fs := &fakeCHStatPopServer{}
	client := &http.Client{Transport: http.NewFileTransport(fs)}
	filepath, err := fetchCHStatPop(client, "https://www.bfs.admin.ch/bfsstatic/dam/assets/9606372/master")
	if err != nil {
		t.Error(err)
		return
	}
	if !strings.HasSuffix(filepath, "extracted.csv") {
		t.Error("expected suffix \"extracted.csv\", got " + filepath)
		return
	}

	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		t.Error(err)
		return
	}
	if !strings.Contains(string(content), "RELI,X_KOORD") {
		t.Error("expected " + filepath + " to contain RELI,X_KOORD")
		return
	}
}
