// SPDX-FileCopyrightText: 2019 Sascha Brawer <sascha@brawer.ch>
// SPDX-License-Identifier: MIT

#ifndef GEOSMELL_TOOLS_CHPOPSTAT_H_
#define GEOSMELL_TOOLS_CHPOPSTAT_H_

#include <utility>
#include <vector>

#include <s2/s2latlng.h>
#include <s2/s2loop.h>

namespace geosmell {

// Convert Swiss Grid coordinates [https://epsg.io/21781] to WGS 84.
S2LatLng SwissGridToLatLng(double y, double x);

// Convert a Swiss statistical register ID to an S2 loop.
void SwissRegisterIDToS2Loop(uint64_t regid, S2Loop *loop);

typedef std::vector<std::pair<S2CellId, double> > OverlapFractions;

// Given a Swiss statistical register ID, Compute which S2 cells are overlapping
// that are by how much. For example, regid=60002000 (a cell near the former
// astronomical observatory in Bern, Switzerland) at level=17 returns seven
// S2CellIds whose overlaps are ranging from 0.01 to 0.36; the total overlap is 1.0.
void GetOverlapFractions(uint64_t regid, int level, OverlapFractions *out);

};  // namespace geosmell

#endif  // GEOSMELL_TOOLS_CHPOPSTAT_H_
