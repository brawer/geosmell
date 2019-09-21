# GeoSmell

Experimental code for computing signals about geographic areas.

To identify geographic areas, we use [S2 cells](http://s2geometry.io/).
Currently, we just count how many pictures have been geo-tagged on
[Wikimedia Commons](https://commons.wikimedia.org/) and aggregate
this to S2 cells. You can choose the desired [aggregation level](https://s2geometry.io/resources/s2cell_statistics); `--level 17` aggregates
to cells of about 70×70 meters in size.

To visualize the geographic region of an S2 cell on a map, you can use
the [Sidewalk Labs Region Coverer](https://s2.sidewalklabs.com/regioncoverer/).  Enter the S2 cell
ID into the “cells” and press the circular button.


## Building and running

```sh
git clone https://github.com/brawer/geosmell.git ; cd geosmell
go build ; go test
./geosmell --level 17
```

This produces a file `wikicommons-20190901.csv.gz` (or similar) with
the number of geocodced pictures on Wikimedia Commons for each S2 cell.


## Further work

It would be nice to collect additional signals about geographic areas.

It would be nice to extend the tool so it stores the resulting file
into a Content Delivery Network such as Digital Ocean Spaces, Microsoft Azure,
or similar. Then, it could be run as a cronjob on Kubernetes.

Probably I won’t have time to do any of this, but feel free to send
pull requests (or fork the project).


## License

Copyright 2019 by Sascha Brawer. Licensed under the MIT license.
