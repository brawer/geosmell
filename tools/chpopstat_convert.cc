// SPDX-FileCopyrightText: 2019 Sascha Brawer <sascha@brawer.ch>
// SPDX-License-Identifier: MIT

#include <iostream>
#include <fstream>

#include <gflags/gflags.h>
#include <glog/logging.h>
#include <s2/s2latlng.h>

#include "chpopstat.h"

int main(int argc, char* argv[]) {
  google::InitGoogleLogging(argv[0]);
  gflags::SetUsageMessage("Convert Swiss population statistics "
                          "from Swiss statistical regions to S2 cells");
  gflags::ParseCommandLineFlags(&argc, &argv, true);

  if (argc != 4) {
    std::cerr << "Usage: "
	      << gflags::ProgramInvocationShortName()
	      << " path/to/input.csv path/to/output.csv 17"
	      << std::endl;
    return 1;
  }

  // TODO: Use gflags. Somehow gflags don’t get recognized when declared locally.
  std::ifstream input(argv[1]);
  std::ofstream output(argv[2]);
  int level = std::atoi(argv[3]);
  geosmell::ConvertSwissPopulationStats(&input, level, &output);
  return 0;
}
