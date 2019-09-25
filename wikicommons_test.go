// SPDX-FileCopyrightText: 2019 Sascha Brawer <sascha@brawer.ch>
// SPDX-License-Identifier: MIT

package main

import (
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
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

	case "/mirror/wikimedia.org/dumps/commonswiki/20190820/commonswiki-20190820-geo_tags.sql.gz":
		return os.Open("testdata/wikicommons/geo_tags.sql.gz")
	}
	return os.Open("testdata/wikicommons/notfound.html")
}

func NewTestClient() *http.Client {
	fs := &FakeServer{}
	return &http.Client{Transport: http.NewFileTransport(fs)}
}

func TestWikiCommonsFindUpstreamVersion(t *testing.T) {
	d, _ := NewDataset("wikicommons", NewTestClient())
	if version, err := d.FindUpstreamVersion(); err == nil {
		equals(t, "2019-08-20", version.String()[:10])
	} else {
		t.Error(err)
	}
}

func TestWikiCommonsProcess(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "geosmell-test")
	if err != nil {
		t.Error(err)
		return
	}
	defer os.RemoveAll(tempDir)

	filePath := path.Join(tempDir, "out.gz")
	d, err := NewDataset("wikicommons", NewTestClient())
	if err != nil {
		t.Error(err)
		return
	}

	if err := d.Process(17, filePath); err != nil {
		t.Error(err)
		return
	}

	stream, err := os.Open(filePath)
	if err != nil {
		t.Error(err)
		return
	}
	defer stream.Close()

	gzstream, err := gzip.NewReader(stream)
	if err != nil {
		t.Error(err)
		return
	}
	defer gzstream.Close()

	contentBytes, err := ioutil.ReadAll(gzstream)
	if err != nil {
		t.Error(err)
		return
	}
	result := string(contentBytes)

	if !strings.Contains(result, "\n1027878a84,21\n") {
		t.Error("Expected result to contain '1027878a84,21'; got " + result)
	}
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
