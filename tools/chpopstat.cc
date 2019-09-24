// SPDX-FileCopyrightText: 2019 Sascha Brawer <sascha@brawer.ch>
// SPDX-License-Identifier: MIT

#include <ctype.h>
#include <string.h>
#include <cmath>
#include <cstdlib>
#include <iomanip>
#include <memory>
#include <utility>
#include <vector>

#include <glog/logging.h>

#include <s2/s2cell_id.h>
#include <s2/s2latlng.h>
#include <s2/s2loop.h>
#include <s2/s2polygon.h>
#include <s2/s2region_coverer.h>

#include "chpopstat.h"

namespace geosmell {

void ConvertSwissPopulationStats(std::istream *input, int level, std::ostream *output) {
  CellPopulationStats cellStats;
  CSVParser parser(input);
  uint64_t regionId;
  PopulationStats regionStats;
  while (parser.Next(&regionId, &regionStats)) {
    LOG_EVERY_N(INFO, 1000) << "Processing statistical region: " << regionId;
    DistributeStats(regionId, regionStats, level, &cellStats);
  }
  *output << "S2CellId,TotalPopulation,FemalePopulation,MalePopulation\n";
  for (auto s : cellStats) {
      S2CellId cellId = s.first;
      const PopulationStats& stat = s.second;
      long numTotal = std::round(stat.numTotal);
      long numFemale = std::round(stat.numFemale);
      // Prevent impossible output in case of (rare) rounding errors.
      if (numFemale > numTotal) {
          numFemale = numTotal;
      }
      long numMale = numTotal - numFemale;
      // Suppress all-zero cells; this can happen due to rounding.
      if (numTotal > 0) {
          *output << cellId.ToToken()
                  << ',' << numTotal
                  << ',' << numFemale
                  << ',' << numMale
                  << '\n';
    }
  }
  output->flush();
}

void DistributeStats(uint64_t regionId,
                     const PopulationStats& stats,
                     int level, CellPopulationStats *cellStats) {
  OverlapFractions overlaps;
  GetOverlapFractions(regionId, level, &overlaps);
  for (auto overlap : overlaps) {
    S2CellId cellId = overlap.first;
    double fraction = overlap.second;
    PopulationStats* s = &((*cellStats)[cellId]);
    s->numTotal += fraction * stats.numTotal;
    s->numFemale += fraction * stats.numFemale;
    s->numMale += fraction * stats.numMale;
  }
}

CSVParser::CSVParser(std::istream* stream) :
  stream_(stream),
  columnAreaId_(-1), columnTotal_(-1), columnFemale_(-1), columnMale_(-1)
{
  std::string header;
  if (!std::getline(*stream_, header)) {
    return;
  }
  size_t start = 0, end = 0;
  int columnId = 0;
  while ((end = header.find(",", start)) != std::string::npos) {
    const std::string columnName = header.substr(start, end - start);
    if (columnName == "RELI") {
      columnAreaId_ = columnId;
    }
    if (columnName.size() > 4 &&
	columnName[0] == 'B' && isdigit(columnName[1]) && isdigit(columnName[2])) {
      std::string h = columnName.substr(3);
      if (h == "BTOT") {
	columnTotal_ = columnId;
      } else if (h == "BWTOT") {
	columnFemale_ = columnId;
      } else if (h == "BMTOT") {
	columnMale_ = columnId;
      }
    }
    start = end + 1;
    columnId += 1;
  }
}

CSVParser::~CSVParser() {
}

bool CSVParser::Next(uint64_t *areaId, PopulationStats *stats) {
  std::string header;
  *areaId = 0;
  memset(static_cast<void*>(stats), 0, sizeof(*stats));
  if (!std::getline(*stream_, header)) {
    return false;
  }
  size_t start = 0, end = 0;
  int columnId = 0;
  while ((end = header.find(",", start)) != std::string::npos) {
    const int value = std::atoi(header.substr(start, end - start).c_str());
    if (columnId == columnAreaId_) {
      *areaId = value;
    } else if (columnId == columnTotal_) {
      stats->numTotal += value;
    } else if (columnId == columnFemale_) {
      stats->numFemale += value;
    } else if (columnId == columnMale_) {
      stats->numMale += value;
    }
    start = end + 1;
    columnId += 1;
  }
  return true;
}

S2LatLng SwissGridToLatLng(double y, double x) {
  // Conversion from: Swiss Confederation, Federal Office of Topography.
  // Formeln und Konstanten für die Berechnung der Schweizerischen
  // schiefachsigen Zylinderprojektion und der Transformation zwischen
  // Koordinatensystemen. October 2018.  https://tinyurl.com/y32bhewm
  y = (y - 600000) / 1000000.0;
  x = (x - 200000) / 1000000.0;
  double lat =
      16.9023892
      + (3.238272 * x)
      - (0.270978 * y * y)
      - (0.002528 * x * x)
      - (0.0447 * y * y * x)
      - (0.0140 * x * x * x);
  double lng =
      2.6779094
      + (4.728982 * y)
      + (0.791484 * y * x)
      + (0.1306 * y * x * x)
      - (0.0436 * y * y * y);
  return S2LatLng::FromDegrees(lat * 100.0 / 36.0, lng * 100.0 / 36.0);
}

void SwissRegisterIDToS2Loop(uint64_t regid, S2Loop *loop) {
  const double y = (regid / 10000) * 100.0;
  const double x = (regid % 10000) * 100.0;
  const double y1 = y + 100.0;
  const double x1 = x + 100.0;

  std::vector<S2Point> vertices;
  vertices.reserve(4);
  vertices.push_back(SwissGridToLatLng(y, x).ToPoint());
  vertices.push_back(SwissGridToLatLng(y1, x).ToPoint());
  vertices.push_back(SwissGridToLatLng(y1, x1).ToPoint());
  vertices.push_back(SwissGridToLatLng(y, x1).ToPoint());

  loop->Init(vertices);
}

void GetOverlapFractions(uint64_t regid, int level, OverlapFractions *out) {
  std::unique_ptr<S2Loop> regLoop(new S2Loop());
  SwissRegisterIDToS2Loop(regid, regLoop.get());
  S2Polygon regPoly;
  regPoly.Init(std::move(regLoop));

  S2RegionCoverer::Options options;
  options.set_min_level(level);
  options.set_max_level(level);
  options.set_max_cells(10000);

  std::vector<S2CellId> cells;
  S2RegionCoverer coverer(options);
  coverer.GetCovering(regPoly, &cells);

  out->clear();
  out->reserve(cells.size());
  for (S2CellId cellId : cells) {
    S2Polygon cellPoly;
    std::unique_ptr<S2Loop> cellLoop(new S2Loop(S2Cell(cellId)));
    cellPoly.Init(std::move(cellLoop));
    std::pair<double, double> overlap = S2Polygon::GetOverlapFractions(&regPoly, &cellPoly);
    out->push_back(std::make_pair(cellId, overlap.first));
  }
}

};  // namespace geosmell
