package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/golang/geo/s2"
)

func main() {
	level := flag.Int("level", 17, "Level of S2 cells being aggregated")
	flag.Parse()
	//resp, err := http.Get("http://en.wikipedia.org/")

	// http://ftp.acc.umu.se/mirror/wikimedia.org/dumps/commonswiki/20190820/commonswiki-20190820-geo_tags.sql.gz
	stream, err := os.Open("commonswiki-20190820-geo_tags.sql")
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
	writeCounts(buf)

	if err := c.Err(); err != nil {
		log.Fatal(err)
	}
}

type S2Cells []s2.CellID

func (a S2Cells) Len() int           { return len(a) }
func (a S2Cells) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a S2Cells) Less(i, j int) bool { return a[i] < a[j] }

func writeCounts(counts map[s2.CellID]int64) {
	cells := make(S2Cells, len(counts))
	i := 0
	for cell, _ := range counts {
		cells[i] = cell
		i++
	}
	sort.Sort(cells)
	for _, cell := range cells {
		fmt.Printf("%s,%d\n", cell.ToToken(), counts[cell])
	}
}
