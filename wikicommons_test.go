// SPDX-FileCopyrightText: 2019 Sascha Brawer <sascha@brawer.ch>
// SPDX-License-Identifier: MIT

package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

type FakeServer struct{}

func (s FakeServer) Open(path string) (f http.File, e error) {
	switch path {
	case "/mirror/wikimedia.org/dumps/commonswiki":
		return os.Open("testdata/wikicommons/index.html")

	case "/mirror/wikimedia.org/dumps/commonswiki/20190820":
		return os.Open("testdata/wikicommons/20190820.html")
	}
	return os.Open("testdata/wikicommons/notfound.html")
}

func NewTestClient() *http.Client {
	fs := &FakeServer{}
	return &http.Client{Transport: http.NewFileTransport(fs)}
}

func TestFindLatestWikiCommons(t *testing.T) {
	tc := NewTestClient()
	equals(t, "2019-08-20", findLatestWikiCommons(tc).String()[:10])
}

func TestWikiCommonsParser(t *testing.T) {
	expected := []float64{
		44.34099960, 8.55650043,
		0.0, 0.0,
		53.14299194, 9.88410444}
	got := parse(`
/*!40000 ALTER TABLE 'geo_tags' DISABLE KEYS */;
INSERT INTO 'geo_tags' VALUES (56,18518224,'earth',1,44.34099960,8.55650043,NULL,NULL,NULL,NULL,NULL),(1509,17171704,'earth',1,0.00000000,0.00000000,NULL,NULL,NULL,NULL,NULL);
INSERT INTO 'geo_tags' VALUES (158915664,42805125,'earth',1,53.14299194,9.88410444,NULL,NULL,NULL,NULL,NULL);`)
	equals(t, expected, got)
}

func parse(s string) []float64 {
	r := make([]float64, 0)
	p := NewWikiCommonsParser(strings.NewReader(s))
	for p.Next() {
		r = append(r, p.Lat)
		r = append(r, p.Lon)
	}
	return r
}

func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}
