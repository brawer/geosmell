package main

import (
	"bufio"
	"bytes"
	"io"
	"regexp"
	"strconv"
)

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
