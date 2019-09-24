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

static const char* kTestCSV =
    "RELI,X_KOORD,Y_KOORD,E_KOORD,N_KOORD,B18BTOT,B18B11,B18B12,B18B13,"
    "B18B14,B18B15,B18B16,B18B21,B18B22,B18B23,B18B24,B18B25,B18B26,B18B27,"
    "B18B28,B18B29,B18B30,B18BMTOT,B18BWTOT,B18B41,B18B42,B18B43,B18B44,"
    "B18B45,B18B46,B18B51,B18B52,B18B53,B18B54,B18B55,B18B56\n"
    "49221163,492200,116300,2492200,1116300,48,39,9,6,3,3,0,37,0,33,4,0,11,"
    "5,3,5,0,23,25,7,12,3,12,14,0,40,3,3,3,3,3\n"
    "60002000,600000,200000,2600000,1200000,4,3,3,3,0,0,0,3,3,3,0,0,3,3,0,"
    "0,0,3,3,0,0,0,3,3,0,4,0,0,0,0,0\n";

TEST(ConvertSwissPopulationStats, ShouldProcessTestFile) {
  std::istringstream input(kTestCSV);
  std::ostringstream output;
  ConvertSwissPopulationStats(&input, 16, &output);
  EXPECT_EQ(
      output.str(),
      "S2CellId,TotalPopulation,FemalePopulation,MalePopulation\n"
      "478c7d241,9,4,5\n"
      "478c7d243,22,11,11\n"
      "478c7d269,12,6,6\n"
      "478c7d26b,5,3,2\n"
      "478e39be5,4,3,1\n");
}

TEST(CSVParser, TwoLines) {
  std::istringstream stream(kTestCSV);
  CSVParser parser(&stream);
  uint64_t areaId;
  PopulationStats stats;

  ASSERT_TRUE(parser.Next(&areaId, &stats));
  EXPECT_EQ(areaId, 49221163);
  EXPECT_EQ(stats.numTotal, 48);
  EXPECT_EQ(stats.numFemale, 25);
  EXPECT_EQ(stats.numMale, 23);

  // For privacy, the data source never states less than 3 for any metric.
  ASSERT_TRUE(parser.Next(&areaId, &stats));
  EXPECT_EQ(areaId, 60002000);
  EXPECT_EQ(stats.numTotal, 4);
  EXPECT_EQ(stats.numFemale, 3);
  EXPECT_EQ(stats.numMale, 3);

  ASSERT_FALSE(parser.Next(&areaId, &stats));
  EXPECT_EQ(areaId, 0);
  EXPECT_EQ(stats.numTotal, 0);
  EXPECT_EQ(stats.numFemale, 0);
  EXPECT_EQ(stats.numMale, 0);
}

}  // namespace geosmell

int main(int argc, char **argv) {
  ::testing::InitGoogleTest(&argc, argv);
  return RUN_ALL_TESTS();
}
