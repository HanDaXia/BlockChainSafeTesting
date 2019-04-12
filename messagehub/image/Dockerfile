# Copyright Greg Haskins All Rights Reserved
#
# SPDX-License-Identifier: Apache-2.0
#
FROM ubuntu:16.04
MAINTAINER xl-cetc

COPY payload/messagehub /usr/local/bin
RUN chmod a+x /usr/local/bin/messagehub
EXPOSE 8000
CMD messagehub > ~/messagehub.log
