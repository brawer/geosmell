// SPDX-FileCopyrightText: 2019 Sascha Brawer <sascha@brawer.ch>
// SPDX-License-Identifier: MIT

#include <gtest/gtest.h>
#include <s2/s2latlng.h>

#include "chpopstat.h"

namespace geosmell {

// Test cases from: Swiss Confederation, Federal Office of Topography.
// Formeln und Konstanten für die Berechnung der Schweizerischen
// schiefachsigen Zylinderprojektion und der Transformation zwischen
// Koordinatensystemen. October 2018.  https://tinyurl.com/y32bhewm
TEST(SwissGridToLatLng, Bern) {
  const S2LatLng bern = SwissGridToLatLng(600000.000, 200000.000);
  EXPECT_NEAR(bern.lat().degrees(), 46.9510811, 1e-7);
  EXPECT_NEAR(bern.lng().degrees(),  7.4386372, 1e-7);
}

TEST(SwissGridToLatLng, Zimmerwald) {
  const S2LatLng zimmerwald = SwissGridToLatLng(602030.680, 191775.030);
  EXPECT_NEAR(zimmerwald.lat().degrees(), 46.8770923, 1e-7);
  EXPECT_NEAR(zimmerwald.lng().degrees(),  7.4652757, 1e-7);
}

TEST(SwissGridToLatLng, MonteGeneroso) {
  const S2LatLng monteGeneroso = SwissGridToLatLng(722758.810, 87649.670);
  EXPECT_NEAR(monteGeneroso.lat().degrees(), 45.9293009, 1e-7);
  EXPECT_NEAR(monteGeneroso.lng().degrees(),  9.0212199, 1e-7);
}

TEST(SwissRegisterIDToS2Loop, Bern) {
  S2Loop loop;
  SwissRegisterIDToS2Loop(60002000, &loop);
  EXPECT_TRUE(loop.IsValid());
  EXPECT_TRUE(loop.IsNormalized());
  EXPECT_FALSE(loop.is_hole());
  EXPECT_EQ(loop.num_vertices(), 4);
  EXPECT_NEAR(S2LatLng(loop.vertex(0)).lat().degrees(), 46.9510811, 1e-7);
  EXPECT_NEAR(S2LatLng(loop.vertex(0)).lng().degrees(),  7.4386372, 1e-7);
  EXPECT_NEAR(S2LatLng(loop.vertex(1)).lat().degrees(), 46.9510811, 1e-7);
  EXPECT_NEAR(S2LatLng(loop.vertex(1)).lng().degrees(),  7.4399508, 1e-7);
  EXPECT_NEAR(S2LatLng(loop.vertex(2)).lat().degrees(), 46.9519806, 1e-7);
  EXPECT_NEAR(S2LatLng(loop.vertex(2)).lng().degrees(),  7.4399508, 1e-7);
  EXPECT_NEAR(S2LatLng(loop.vertex(3)).lat().degrees(), 46.9519806, 1e-7);
  EXPECT_NEAR(S2LatLng(loop.vertex(3)).lng().degrees(),  7.4386372, 1e-7);
}

TEST(GetOverlapFractions, Bern) {
  OverlapFractions overlaps;
  GetOverlapFractions(60002000, 18, &overlaps);
  double totalOverlap = 0.0;
  for (auto ov : overlaps) {
    LOG(INFO) << "CellId: " << ov.first.ToToken() << ", overlap fraction: " << ov.second;
    EXPECT_EQ(ov.first.level(), 18);
    totalOverlap += ov.second;
  }
  EXPECT_NEAR(totalOverlap, 1.0, 0.01);
}

}  // namespace geosmell

int main(int argc, char **argv) {
  ::testing::InitGoogleTest(&argc, argv);
  return RUN_ALL_TESTS();
}
