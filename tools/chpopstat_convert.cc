// SPDX-FileCopyrightText: 2019 Sascha Brawer <sascha@brawer.ch>
// SPDX-License-Identifier: MIT

#include <gflags/gflags.h>
#include <glog/logging.h>
#include <s2/s2cell_id.h>


// TODO: The point of this program is just to make sure that we can
// use glog, gflags, and s2 at the same time in a statically linked
// Linux binary that can run independently (without all the build-time
// dependencies) inside a Docker comtainer. Currently, it does not
// do anything useful yet.

int main(int argc, char* argv[]) {
  gflags::ParseCommandLineFlags(&argc, &argv, true);
  google::InitGoogleLogging(argv[0]);

  const S2CellId i = S2CellId::FromFace(4);
  LOG(INFO) << "Hello world, here are " << 12 << " cookies for you; "
	    << i.ToToken();
  return 0;
}
