# Copyright Greg Haskins All Rights Reserved
#
# SPDX-License-Identifier: Apache-2.0
#
FROM ubuntu:16.04
MAINTAINER xl-cetc

COPY payload/normalCrypto /usr/local/bin
COPY payload/libtest.so /usr/local/bin
RUN mkdir /usr/local/bin/templates
ADD payload/templates/ /usr/local/bin/templates
EXPOSE 8081
CMD normalCrypto
