// SPDX-FileCopyrightText: 2019 Sascha Brawer <sascha@brawer.ch>
// SPDX-License-Identifier: MIT

package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/golang/geo/s2"
)

type WikiCommons struct {
	client          *http.Client
	upstreamVersion *time.Time
}

const baseUrl = "https://ftp.acc.umu.se/mirror/wikimedia.org/dumps/commonswiki/"

func fetchWikiCommons(client *http.Client, version time.Time) (*http.Response, error) {
	date := fmt.Sprintf("%04d%02d%02d", version.Year(), version.Month(), version.Day())
	url := fmt.Sprintf("%s/%s/commonswiki-%s-geo_tags.sql.gz", baseUrl, date, date)
	return client.Get(url)
}

func (c WikiCommons) FindUpstreamVersion() (*time.Time, error) {
	if c.upstreamVersion != nil {
		return c.upstreamVersion, nil
	}

	resp, err := c.client.Get(baseUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	dateSet := make(map[string]bool)
	re := regexp.MustCompile("<a href=\"(2[0-9]{7})/\">")
	for _, match := range re.FindAllStringSubmatch(string(body), -1) {
		dateSet[match[1]] = true
	}
	dates := make([]string, 0, len(dateSet))
	for date, _ := range dateSet {
		dates = append(dates, date)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(dates)))
	for _, date := range dates {
		resp, err := c.client.Get(baseUrl + "/" + date + "/")
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		if strings.Contains(string(body), "geo_tags.sql.gz") {
			year, _ := strconv.Atoi(date[0:4])
			month, _ := strconv.Atoi(date[4:6])
			day, _ := strconv.Atoi(date[6:8])
			date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
			c.upstreamVersion = &date
			return c.upstreamVersion, nil
		}
	}

	return nil, errors.New("cannot find any Wikimedia Commons dump")
}

func (c WikiCommons) Process(s2Level int, outpath string) error {
	version, err := c.FindUpstreamVersion()
	if err != nil {
		return err
	}

	resp, err := fetchWikiCommons(c.client, *version)
	if err != nil {
		return err
	}

	stream, err := gzip.NewReader(resp.Body)
	if err != nil {
		return err
	}

	parser := NewWikiCommonsParser(stream)
	buf := make(map[s2.CellID]int64, 7000000) // 5797447 at level 17
	for parser.Next() {
		if parser.Lat == 0 && parser.Lon == 0 {
			continue
		}
		latLng := s2.LatLngFromDegrees(parser.Lat, parser.Lon)
		cellID := s2.CellIDFromLatLng(latLng).Parent(s2Level)
		buf[cellID] += 1
	}
	if err := parser.Err(); err != nil {
		return err
	}

	out, err := os.Create(outpath)
	if err != nil {
		return err
	}
	defer out.Close()

	gzout := gzip.NewWriter(out)
	defer gzout.Close()

	writeCounts(buf, gzout)
	return nil
}

type WikiCommonsParser struct {
	scanner  *bufio.Scanner
	Lat, Lon float64
}

func NewWikiCommonsParser(r io.Reader) *WikiCommonsParser {
	scanner := bufio.NewScanner(r)
	scanner.Split(splitParen)
	return &WikiCommonsParser{scanner, 0.0, 0.0}
}

func splitParen(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.IndexByte(data, ')'); i >= 0 {
		start := bytes.IndexByte(data, '(')
		if start < 0 || start >= i {
			start = 0
		} else {
			start = start + 1
		}
		return i + 1, data[start:i], nil
	}
	if atEOF {
		return len(data), data, nil
	}
	return 0, nil, nil
}

var extractCoords = regexp.MustCompile("^[0-9]+,[0-9]+,'earth',[0-9]+,([0-9\\-.]+),([0-9\\-.]+),")

func (c *WikiCommonsParser) Next() bool {
	for c.scanner.Scan() {
		m := extractCoords.FindSubmatch([]byte(c.scanner.Text()))
		if len(m) == 3 {
			lat, _ := strconv.ParseFloat(string(m[1]), 64)
			lon, _ := strconv.ParseFloat(string(m[2]), 64)
			c.Lat = lat
			c.Lon = lon
			return true
		}
	}
	c.Lat = 0
	c.Lon = 0
	return false
}

func (c *WikiCommonsParser) Err() error {
	return c.scanner.Err()
}

type S2Cells []s2.CellID

func (a S2Cells) Len() int           { return len(a) }
func (a S2Cells) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a S2Cells) Less(i, j int) bool { return a[i] < a[j] }

func writeCounts(counts map[s2.CellID]int64, out io.Writer) {
	cells := make(S2Cells, len(counts))
	i := 0
	for cell, _ := range counts {
		cells[i] = cell
		i++
	}
	sort.Sort(cells)
	for _, cell := range cells {
		fmt.Fprintf(out, "%s,%d\n", cell.ToToken(), counts[cell])
	}
}
