// SPDX-FileCopyrightText: 2019 Sascha Brawer <sascha@brawer.ch>
// SPDX-License-Identifier: MIT

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

const baseUrl = "https://ftp.acc.umu.se/mirror/wikimedia.org/dumps/commonswiki/"

func fetchWikiCommons(client *http.Client, version time.Time) (*http.Response, error) {
	date := fmt.Sprintf("%04d%02d%02d", version.Year(), version.Month(), version.Day())
	url := fmt.Sprintf("%s/%s/commonswiki-%s-geo_tags.sql.gz", baseUrl, date, date)
	return client.Get(url)
}

func findLatestWikiCommons(client *http.Client) time.Time {
	resp, err := client.Get(baseUrl)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
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
		resp, err := client.Get(baseUrl + "/" + date + "/")
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		if strings.Contains(string(body), "geo_tags.sql.gz") {
			year, _ := strconv.Atoi(date[0:4])
			month, _ := strconv.Atoi(date[4:6])
			day, _ := strconv.Atoi(date[6:8])
			return time.Date(year, time.Month(month), day, 0, 0, 0, 0,
				time.UTC)
		}
	}

	log.Fatal("cannot find any Wikimedia Commons dump")
	return time.Date(2009, time.Month(11), 25, 0, 0, 0, 0, time.UTC)
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
