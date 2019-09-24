// SPDX-FileCopyrightText: 2019 Sascha Brawer <sascha@brawer.ch>
// SPDX-License-Identifier: MIT

#ifndef GEOSMELL_TOOLS_CHPOPSTAT_H_
#define GEOSMELL_TOOLS_CHPOPSTAT_H_

#include <iostream>
#include <map>
#include <utility>
#include <vector>

#include <s2/s2latlng.h>
#include <s2/s2loop.h>

namespace geosmell {

// Converts Swiss population statistics from the format of the
// Swiss Federal Statistical Office to our own CSV format.
//
// In the input format, geographic regions are identified
// by eight-digit numbers. The most significant four digits
// correspond to the y coordinate in the Swiss national grid
// divided by 100, and the least significant four digits
// correspond to the x coordinate divided by 100.
//
// In the output format, geographic regions are identified 
// by cell IDs in the S2 geometry library.
void ConvertSwissPopulationStats(std::istream *input,
                                 int level, std::ostream *output);

// Population statistics about a geographic area.
//
// We allow for fractional counts to better handle rounding errors
// that result from converting Swiss statistical regions to S2 cells.
// For example, Swiss statistical region 60002000 is an area near
// the former astronomical observaty in Bern. At S2 subdivision level 17,
// the statistical regionis is overlapping seven cells. We distribute
// the counts for the statistical region to those seven cells according
// to how much the S2 cells geometrically overlap with the statistical
// region. For example, when an S2 cell is overlapping by 17.3%, we
// account 17.3% of the population in the statistical region towards
// that cell.
struct PopulationStats {
  PopulationStats() : numTotal(0), numFemale(0), numMale(0) {}
  double numTotal;   // Property BTOT: Permanent resident population, total
  double numFemale;  // Property BWTOT: Permanent resident population, female
  double numMale;    // Property BMTOT: Permanent resident population, male
};

typedef std::map<S2CellId, PopulationStats> CellPopulationStats;

// Distribute population statistics for a Swiss statistical region
// to S2 cells of a given level. Finds all S2 cells that geographically
// overlap the given statistical region, and distributes the population
// counts to each overlapping cell according to the overlap fraction.
void DistributeStats(uint64_t regionId,
                     const PopulationStats& regionStats,
                     int level, CellPopulationStats *cellStats);

// Utility class for parsing CSV files with populations statistics
// from bfs.admin.ch.
class CSVParser {
 public:
  explicit CSVParser(std::istream* stream);
  virtual ~CSVParser();
  bool Next(uint64_t *areaId, PopulationStats *stats);

 private:
  std::istream* stream_;
  int columnAreaId_, columnTotal_, columnFemale_, columnMale_;
};


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
