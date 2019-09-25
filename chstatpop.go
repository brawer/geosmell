// SPDX-FileCopyrightText: 2019 Sascha Brawer <sascha@brawer.ch>
// SPDX-License-Identifier: MIT

package main

import (
	"archive/zip"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"strconv"
	"time"
)

const chStatPopHost = "https://www.bfs.admin.ch"

const chStatPopListUrl = "/bfs/de/home/dienstleistungen/geostat/geodaten-bundesstatistik/gebaeude-wohnungen-haushalte-personen/bevoelkerung-haushalte-ab-2010/_jcr_content/par/tabs/items/geodaten_statpop/tabpar/ws_parametrized_list.dynamiclist.html"

type SwissPopulationStatistics struct {
	client *http.Client
}

func (s SwissPopulationStatistics) FindUpstreamVersion() (*time.Time, error) {
	_, t, err := findLatestCHStatPop(s.client)
	return t, err
}

func (s SwissPopulationStatistics) Process(s2Level int, outpath string) error {
	return errors.New("not yet implemented")
}

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

func fetchCHStatPop(client *http.Client, url string) (string, error) {
	tempDir, err := ioutil.TempDir("", "geosmell-chstatpop")
	fetchedPath := path.Join(tempDir, "fetched.zip")
	extractedPath := path.Join(tempDir, "extracted.csv")
	if err != nil {
		return "", err
	}

	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	fetchedFile, err := os.Create(fetchedPath)
	if err != nil {
		return "", err
	}
	defer fetchedFile.Close()

	_, err = io.Copy(fetchedFile, resp.Body)

	if err != nil {
		return "", err
	}

	extractedFile, err := os.Create(extractedPath)
	if err != nil {
		return "", err
	}
	defer extractedFile.Close()

	zipFile, err := zip.OpenReader(fetchedFile.Name())
	if err != nil {
		return "", err
	}

	re := regexp.MustCompile(`STATPOP20[0-9]{2}G\.csv$`)
	for _, file := range zipFile.File {
		if re.FindString(file.Name) != "" {
			statFile, err := file.Open()
			if err != nil {
				return "", err
			}
			defer statFile.Close()
			_, err = io.Copy(extractedFile, statFile)
			if err != nil {
				return "", err
			}
		}
	}

	// TODO: Extract CSV with raw statistics, write to exractedFile.
	return extractedPath, nil
}
