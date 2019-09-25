// SPDX-FileCopyrightText: 2019 Sascha Brawer <sascha@brawer.ch>
// SPDX-License-Identifier: MIT

package main

import (
	"net/http"
	"os"
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

func TestSwissPopulationStatistics(t *testing.T) {
	fs := &fakeCHStatPopServer{}
	client := &http.Client{Transport: http.NewFileTransport(fs)}
	dataset, err := NewDataset("chpopstat", client)
	if err != nil {
		t.Error(err)
		return
	}

	version, err := dataset.FindUpstreamVersion()
	if err != nil {
		t.Error(err)
		return
	}

	equals(t, version.String()[:10], "2019-08-27")
	result, err := processDataset(dataset, 17)
	if err != nil {
		if strings.Contains(err.Error(), "chpopstat_convert") &&
			strings.Contains(err.Error(), "executable file not found") {
			t.Skip("executable file \"chpopstat_convert\" not found; skipping tests on conversion results")
			return
		}
		t.Error(err)
		return
	}

	if !strings.Contains(result, "S2CellId,TotalPopulation") {
		t.Error("expected S2CellId,TotalPopulation in result; got: " + result)
	}
}
