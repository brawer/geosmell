// SPDX-FileCopyrightText: 2019 Sascha Brawer <sascha@brawer.ch>
// SPDX-License-Identifier: MIT

package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestWikiCommonsParser(t *testing.T) {
	expected := []float64{
		44.34099960, 8.55650043,
		0.0, 0.0,
		53.14299194, 9.88410444}
	got := parse(`
/*!40000 ALTER TABLE 'geo_tags' DISABLE KEYS */;
INSERT INTO 'geo_tags' VALUES (56,18518224,'earth',1,44.34099960,8.55650043,NULL,NULL,NULL,NULL,NULL),(1509,17171704,'earth',1,0.00000000,0.00000000,NULL,NULL,NULL,NULL,NULL);
INSERT INTO 'geo_tags' VALUES (158915664,42805125,'earth',1,53.14299194,9.88410444,NULL,NULL,NULL,NULL,NULL);`)
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("expected %v, got %v", expected, got)
	}
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
