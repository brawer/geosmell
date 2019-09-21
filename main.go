// SPDX-FileCopyrightText: 2019 Sascha Brawer <sascha@brawer.ch>
// SPDX-License-Identifier: MIT

package main

import (
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sort"

	"github.com/golang/geo/s2"
)

func main() {
	level := flag.Int("level", 17, "Level of S2 cells being aggregated")
	flag.Parse()

	client := &http.Client{}
	commonsVersion := findLatestWikiCommons(client)
	fmt.Printf("Fetching geotags of Wikimedia Commons, using version: %s\n",
		commonsVersion.String()[:10])
	resp, err := fetchWikiCommons(client, commonsVersion)
	if err != nil {
		log.Fatal(err)
	}

	stream, err := gzip.NewReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	c := NewWikiCommonsParser(stream)
	buf := make(map[s2.CellID]int64, 7000000) // 5797447 at level 17
	for c.Next() {
		if c.Lat == 0 && c.Lon == 0 {
			continue
		}
		latLng := s2.LatLngFromDegrees(c.Lat, c.Lon)
		cellID := s2.CellIDFromLatLng(latLng).Parent(*level)
		buf[cellID] += 1
	}
	filename := fmt.Sprintf("wikicommons-%04d%02d%02d.csv",
		commonsVersion.Year(), commonsVersion.Month(),
		commonsVersion.Day())
	out, err := os.Create(filename + ".gz")
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	gzout := gzip.NewWriter(out)
	gzout.Name = filename
	writeCounts(buf, gzout)
	gzout.Close()
	if err := c.Err(); err != nil {
		log.Fatal(err)
	}
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
