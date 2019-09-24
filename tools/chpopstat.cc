// SPDX-FileCopyrightText: 2019 Sascha Brawer <sascha@brawer.ch>
// SPDX-License-Identifier: MIT

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
