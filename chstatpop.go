// SPDX-FileCopyrightText: 2019 Sascha Brawer <sascha@brawer.ch>
// SPDX-License-Identifier: MIT

package main

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"time"
)

const chStatPopHost = "https://www.bfs.admin.ch"

const chStatPopListUrl = "/bfs/de/home/dienstleistungen/geostat/geodaten-bundesstatistik/gebaeude-wohnungen-haushalte-personen/bevoelkerung-haushalte-ab-2010/_jcr_content/par/tabs/items/geodaten_statpop/tabpar/ws_parametrized_list.dynamiclist.html"

func findLatestCHStatPop(client *http.Client) (string, *time.Time, error) {
	errNotFound := errors.New("could not find latest STATPOP dataset at bfs.admin.ch")
	listUrl, err := url.Parse(chStatPopHost + chStatPopListUrl)
	if err != nil {
		return "", nil, err
	}

	resp, err := client.Get(listUrl.String())
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", nil, err
	}

	re := regexp.MustCompile("href=\"(/bfs/de/home/dienstleistungen/geostat/geodaten-bundesstatistik/[^\"]+)\"")
	match := re.FindStringSubmatch(string(body))
	if match == nil {
		return "", nil, errNotFound
	}
	relativeDatasetUrl, err := url.Parse(match[1])
	if err != nil {
		return "", nil, err
	}
	datasetUrl := listUrl.ResolveReference(relativeDatasetUrl)

	resp2, err := client.Get(datasetUrl.String())
	defer resp.Body.Close()
	body2, err := ioutil.ReadAll(resp2.Body)
	if err != nil {
		return "", nil, err
	}

	reUrl := regexp.MustCompile(`<a href="([^"]+)"[^>]*?>Download`)
	matchUrl := reUrl.FindStringSubmatch(string(body2))
	if matchUrl == nil {
		return "", nil, errNotFound
	}
	downloadUrl, err := url.Parse(matchUrl[1])
	if err != nil {
		return "", nil, errNotFound
	}
	downloadUrl = datasetUrl.ResolveReference(downloadUrl)

	rePubDate := regexp.MustCompile("<th>Veröffentlicht am</th>\\s*<td>([0-9]{1,2})\\.([0-9]{1,2})\\.(20[0-9]{2})")
	matchPubDate := rePubDate.FindStringSubmatch(string(body2))
	if matchPubDate == nil {
		return "", nil, errNotFound
	}
	year, _ := strconv.Atoi(matchPubDate[3])
	month, _ := strconv.Atoi(matchPubDate[2])
	day, _ := strconv.Atoi(matchPubDate[1])
	pubDate := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)

	return downloadUrl.String(), &pubDate, nil
}
