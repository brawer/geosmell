// SPDX-FileCopyrightText: 2019 Sascha Brawer <sascha@brawer.ch>
// SPDX-License-Identifier: MIT

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path"
)

func main() {
	datasetName := flag.String("dataset", "", "Name of dataset to be fetched")
	level := flag.Int("level", 17, "Level of S2 cells being aggregated")
	keys := flag.String("keys", "", "Path to JSON file with access keys for S3 storage")
	force := flag.Bool("force", false, "Force data processing even if stored version is up to date")
	flag.Parse()

	client := &http.Client{}
	dataset, err := NewDataset(*datasetName, client)
	if err != nil {
		log.Fatal(err)
	}

	s3Client, err := newS3Client(*keys)
	if err != nil {
		log.Fatal(err)
	}

	storedVersion, err := findStoredVersion(s3Client, *datasetName)
	if err != nil {
		log.Fatal(err)
	}
	if storedVersion != nil {
		log.Printf("Dataset version in storage: %s\n", storedVersion.String()[:10])
	}

	upstreamVersion, err := dataset.FindUpstreamVersion()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Dataset version at upstream: %s\n", upstreamVersion.String()[:10])
	needsWork := *force || storedVersion == nil || storedVersion.Before(*upstreamVersion)
	if !needsWork {
		log.Printf("No work needed, stored dataset is already up to date")
		return
	}

	tempDir, err := ioutil.TempDir("", "geosmell")
	if err != nil {
		log.Fatal(err)
	}
	fileName := fmt.Sprintf(
		"%s-%04d%02d%02d.csv.gz", *datasetName,
		upstreamVersion.Year(), upstreamVersion.Month(), upstreamVersion.Day())
	filePath := path.Join(tempDir, fileName)

	log.Printf("Processing data, output in " + filePath)
	if err := dataset.Process(*level, filePath); err != nil {
		log.Fatal(err)
	}

	if s3Client != nil {
		log.Printf("Uploading data to storage as " + fileName)
		if err := upload(s3Client, fileName, filePath); err != nil {
			log.Fatal(err)
		}
	}
	log.Printf("Done")
}
