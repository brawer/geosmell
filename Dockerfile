# SPDX-FileCopyrightText: 2019 Sascha Brawer <sascha@brawer.ch>
# SPDX-License-Identifier: MIT

FROM alpine:3.10.2 AS builder

RUN apk add --no-cache build-base cmake g++ git go ninja openssl-dev

RUN git clone --branch v2.2.2 https://github.com/gflags/gflags.git /src/gflags
WORKDIR /src/gflags
RUN cmake -H. -Bbuild -GNinja
RUN cmake --build build
RUN cmake --build build --target install

RUN git clone --branch v0.4.0 https://github.com/google/glog.git /src/glog
WORKDIR /src/glog
RUN cmake -H. -Bbuild -GNinja
RUN cmake --build build
RUN cmake --build build --target test
RUN cmake --build build --target install

RUN git clone --branch release-1.8.1 https://github.com/google/googletest.git /src/gtest
WORKDIR /src/gtest
RUN cmake -H. -Bbuild -GNinja
RUN cmake --build build
RUN cmake --build build --target install

RUN git clone --branch v0.9.0 https://github.com/google/s2geometry.git /src/s2geometry
WORKDIR /src/s2geometry
RUN cmake -H. -Bbuild -GNinja -DBUILD_SHARED_LIBS=OFF -DWITH_GFLAGS=ON -DWITH_GLOG=ON
RUN cmake --build build
RUN cmake --build build --target install

COPY . /src/geosmell
WORKDIR /src/geosmell/tools
RUN cmake -H. -Bbuild -GNinja
RUN cmake --build build
RUN cmake --build build --target test
RUN cp /src/geosmell/tools/build/chpopstat_convert /usr/local/bin

WORKDIR /src/geosmell
RUN go mod download
RUN CGO_ENABLED=1 go build -a -o geosmell .
RUN CGO_ENABLED=1 go test

FROM alpine:3.10.2
WORKDIR /run
RUN apk add --no-cache ca-certificates
COPY --from=builder /src/geosmell/tools/build/chpopstat_convert /usr/local/bin
COPY --from=builder /src/geosmell/geosmell /usr/local/bin

CMD /bin/sh
