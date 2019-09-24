# SPDX-FileCopyrightText: 2019 Sascha Brawer <sascha@brawer.ch>
# SPDX-License-Identifier: MIT

cmake_minimum_required(VERSION 3.0.2)
project(geosmell VERSION 0.1)

enable_testing()

find_package(gflags 2.2.0 REQUIRED)
find_package(glog 0.4.0 REQUIRED)
find_package(GTest 1.8.0 REQUIRED)

find_package(OpenSSL REQUIRED)
find_package(Threads REQUIRED)

add_definitions(-DS2_USE_GLOG)

set(CMAKE_FIND_LIBRARY_SUFFIXES ".a")

link_libraries(glog::glog s2 ${OPENSSL_LIBRARIES})

add_library(chpopstat chpopstat.cc chpopstat.h)
add_executable(chpopstat_convert chpopstat_convert.cc)
target_link_libraries(chpopstat_convert -static chpopstat)

add_executable(chpopstat_test chpopstat_test.cc)
target_link_libraries(chpopstat_test chpopstat s2testing ${GTEST_LIBRARIES})
add_test(NAME chpopstat_test COMMAND chpopstat_test --logtostderr)
