// SPDX-FileCopyrightText: 2019 Sascha Brawer <sascha@brawer.ch>
// SPDX-License-Identifier: MIT

package main

import (
	//"fmt"
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
