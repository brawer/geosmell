// SPDX-FileCopyrightText: 2019 Sascha Brawer <sascha@brawer.ch>
// SPDX-License-Identifier: MIT

package main

import (
	"errors"
	"net/http"
	"time"
)

type Dataset interface {
	FindUpstreamVersion() (*time.Time, error)
	Process(s2Level int, outpath string) error
}

func NewDataset(name string, client *http.Client) (Dataset, error) {
	switch name {
	case "chpopstat":
		return SwissPopulationStatistics{client: client}, nil

	case "wikicommons":
		return WikiCommons{client: client}, nil

	default:
		return nil, errors.New("unknown dataset: " + name +
			"; expected one of {chpopstat, wikicommons}")
	}
}
