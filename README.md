# GeoSmell

Experimental code for computing signals about geographic areas.

To identify geographic areas, we use [S2 cells](http://s2geometry.io/).
You can choose the desired
[aggregation level](https://s2geometry.io/resources/s2cell_statistics);
`--level=17` aggregates to cells of about 70×70 meters in size.
To visualize the geographic region of an S2 cell on a map, you can use
the [Sidewalk Labs Region Coverer](https://s2.sidewalklabs.com/regioncoverer/).
Enter the S2 cell ID into the “cells” field and press the circular button.

Currently, we process the following datasets:

* **wikicommons**: geo-tagged pictures at [Wikimedia Commons](https://commons.wikimedia.org/)



## Building and running

```sh
git clone https://github.com/brawer/geosmell.git ; cd geosmell
go build ; go test
./geosmell --dataset=wikicommons --level=17
```


## Storage

If you pass the access keys to an S3-compatible cloud storage system, the tool
will store its output into the cloud. Also, it will check whether a previously
stored version of a dataset is current; if so, it will exit early without doing
work. This makes it possible to use geosmell as a cronjob, for example on
Kubernetes or another cloud container system. Pass `--keys=/path/to/keys.json`
with the following content:

```json
{
    "S3Host": "your.preferred.cloud.storage.example.com",
    "S3AccessKey": "YourAccessKey",
    "S3SecretKey": "YourSecret"
}
```


## Future work

It would be nice to collect additional signals about geographic areas.
Feel free to send pull requests.


## License

The code in this repository is Copyright 2019 by Sascha Brawer,
licensed under the MIT license.  The processed data sets have
their own licences.
