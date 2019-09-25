// SPDX-FileCopyrightText: 2019 Sascha Brawer <sascha@brawer.ch>
// SPDX-License-Identifier: MIT

package main

import (
	"encoding/json"
	"io/ioutil"
	"regexp"
	"strconv"
	"time"

	"github.com/minio/minio-go/v6"
)

func newS3Client(keys string) (*minio.Client, error) {
	type KeyConfig struct {
		S3Host, S3AccessKey, S3SecretKey string
	}

	blob, err := ioutil.ReadFile(keys)
	if err != nil {
		return nil, err
	}

	var config KeyConfig
	if err := json.Unmarshal(blob, &config); err != nil {
		return nil, err
	}

	return minio.New(config.S3Host, config.S3AccessKey, config.S3SecretKey, true)
}

func findStoredVersion(s3Client *minio.Client, dataset string) (*time.Time, error) {
	if dataset == "" {
		return nil, nil
	}

	re := regexp.MustCompile(`^([a-zA-Z0-9]+)-(20\d{2})(\d{2})(\d{2})\.csv.(gz|bz2)$`)
	doneCh := make(chan struct{})
	defer close(doneCh)
	for object := range s3Client.ListObjects("geosmell", dataset+"-", true, doneCh) {
		if object.Err != nil {
			return nil, object.Err
		}
		if match := re.FindStringSubmatch(object.Key); match != nil && match[1] == dataset {
			year, _ := strconv.Atoi(match[2])
			month, _ := strconv.Atoi(match[3])
			day, _ := strconv.Atoi(match[4])
			t := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
			return &t, nil
		}
	}

	return nil, nil
}

func upload(s3Client *minio.Client, filename string, path string) error {
	_, err := s3Client.FPutObject("geosmell", filename, path, minio.PutObjectOptions{})
	return err
}
