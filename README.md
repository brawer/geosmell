# GeoSmell

Experimental code for computing signals about geographic areas.

To identify geographic areas, we use [S2 cells](http://s2geometry.io/).
Currently, we just count how many pictures have been geo-tagged on
[Wikimedia Commons](https://commons.wikimedia.org/) and aggregate
this to S2 cells; their level can be selected at runtime.

To visualize the geographic region of an S2 cell on a map, you can use
the [Sidewalk Labs Region Coverer](https://s2.sidewalklabs.com/regioncoverer/).  Enter the S2 cell
ID into the “cells” and press the circular button.


## Building and running

```sh
git clone https://github.com/brawer/geosmell.git ; cd geosmell
go build ; go test
curl -L http://ftp.acc.umu.se/mirror/wikimedia.org/dumps/commonswiki/20190901/commonswiki-20190901-geo_tags.sql.gz -o geo_tags.sql.gz
./geosmell --level 17 | bzip2 -9 >wikicommons-20190901.csv.bz2
```


## Further work

It would be nice to collect additional signals about geographic areas.

It would be nice to completely automate the manual steps above, so
that the tool could run in the cloud and periodically store the output
into a Content Delivery Network such as Digital Ocean Spaces or
similar.  The tool should automatically find and download
the most recent dump from Wikimedia Commons; likewise, it should
automatically create the output file in compressed form with the
date-tagged name and then finally put the output into storage.

Probably I won’t have time to do any of this, but feel free to send
pull requests (or fork the project).


## License

Copyright 2019 by Sascha Brawer. Licensed under the MIT license.
