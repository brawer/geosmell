# SPDX-FileCopyrightText: 2019 Sascha Brawer <sascha@brawer.ch>
# SPDX-License-Identifier: MIT

cmake_minimum_required(VERSION 3.0.2)
project(geosmell VERSION 0.1)

find_package(gflags 2.2.0 REQUIRED)
find_package(glog 0.4.0 REQUIRED)

find_package(OpenSSL REQUIRED)
find_package(Threads REQUIRED)

add_definitions(-DS2_USE_GLOG)
#add_definitions(-Wno-attributes)
#add_definitions(-Wno-deprecated-declarations)

set(CMAKE_FIND_LIBRARY_SUFFIXES ".a")
set(CMAKE_EXE_LINKER_FLAGS "-static-libgcc -static-libstdc++")


add_executable(chpopstat_convert chpopstat_convert.cc)
target_link_libraries(chpopstat_convert glog::glog s2 ${OPENSSL_LIBRARIES} -static)
