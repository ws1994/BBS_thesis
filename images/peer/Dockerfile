#
# Copyright contributors to the Hyperledger Fabric project
#
# SPDX-License-Identifier: Apache-2.0
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at:
#
# 	  http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

###############################################################################
# Build image
###############################################################################

ARG UBUNTU_VER
FROM ubuntu:${UBUNTU_VER} as builder

ARG TARGETARCH
ARG TARGETOS
ARG FABRIC_VER
ARG GO_VER
ARG GO_TAGS

RUN apt update && apt install -y \
    git \
    gcc \
    curl \
    make

RUN curl -sL https://go.dev/dl/go${GO_VER}.${TARGETOS}-${TARGETARCH}.tar.gz | tar zxf - -C /usr/local
ENV PATH="/usr/local/go/bin:$PATH"

ADD . .

RUN make peer GO_TAGS=${GO_TAGS} FABRIC_VER=${FABRIC_VER}
RUN make ccaasbuilder


###############################################################################
# Runtime image
###############################################################################

ARG UBUNTU_VER
FROM ubuntu:${UBUNTU_VER}

ARG TARGETOS
ARG TARGETARCH
ARG FABRIC_VER

# set up nsswitch.conf for Go's "netgo" implementation
# - https://github.com/golang/go/blob/go1.9.1/src/net/conf.go#L194-L275
# - docker run --rm debian:stretch grep '^hosts:' /etc/nsswitch.conf
RUN echo 'hosts: files dns' > /etc/nsswitch.conf

<<<<<<< HEAD
FROM golang:${GO_VER}-alpine${ALPINE_VER} as golang
RUN apk add --no-cache \
	bash \
	gcc \
	git \
	make \
	musl-dev;

ADD . $GOPATH/src/github.com/hyperledger/fabric
WORKDIR $GOPATH/src/github.com/hyperledger/fabric

FROM golang as peer
RUN apk add --no-cache \
	bash \
	g++ \
	m4 \
	flex \
	bison \
	libc-dev \
	gawk \
	mawk \
	gcc \
	git \
	make \
	musl-dev;

RUN ls
COPY ./pbc-0.5.14.tar.gz ./
COPY ./gmp-6.2.1.tar ./
RUN	tar -zxvf ./pbc-0.5.14.tar.gz pbc-0.5.14
RUN tar -xvf ./gmp-6.2.1.tar gmp-6.2.1
RUN cd ./gmp-6.2.1 && ./configure && make && make install
RUN	cd ./pbc-0.5.14 && ./configure && make && make install
ARG GO_TAGS
RUN make peer GO_TAGS=${GO_TAGS}

FROM peer-base
ENV FABRIC_CFG_PATH /etc/hyperledger/fabric
VOLUME /etc/hyperledger/fabric
VOLUME /var/hyperledger
COPY ./parameter.txt ./root/
COPY ./userPK_ag.dat ./root/
COPY ./userPK_asg.dat ./root/
COPY ./serPK_g.dat ./root/
COPY ./serPK_sg.dat ./root/
COPY ./timeKeyUpdate.dat ./root/
COPY ./userSK_a.dat ./root/
COPY --from=peer /go/src/github.com/hyperledger/fabric/build/bin /usr/local/bin
COPY --from=peer /usr/local/lib /usr/local/lib
COPY --from=peer /usr/local/include /usr/local/include
COPY --from=peer /usr/local/share /usr/local/share
COPY --from=peer /go/src/github.com/hyperledger/fabric/pbc-0.5.14 /home/pbc-0.5.14
COPY --from=peer /go/src/github.com/hyperledger/fabric/gmp-6.2.1 /home/gmp-6.2.1
COPY --from=peer /go/src/github.com/hyperledger/fabric/sampleconfig/msp ${FABRIC_CFG_PATH}/msp
COPY --from=peer /go/src/github.com/hyperledger/fabric/sampleconfig/core.yaml ${FABRIC_CFG_PATH}
EXPOSE 7051
CMD ["peer","node","start"]
=======
ENV     FABRIC_CFG_PATH /var/hyperledger/fabric/config
ENV     FABRIC_VER      ${FABRIC_VER}

COPY    --from=builder  build/bin/peer          /usr/local/bin
COPY    --from=builder  sampleconfig/msp        ${FABRIC_CFG_PATH}/msp
COPY    --from=builder  sampleconfig/core.yaml  ${FABRIC_CFG_PATH}/core.yaml

COPY    --from=builder  release/${TARGETOS}-${TARGETARCH}/builders/ccaas/bin /opt/hyperledger/ccaas_builder/bin

VOLUME  /etc/hyperledger/fabric
VOLUME  /var/hyperledger

EXPOSE  7051

CMD     [ "peer", "node", "start" ]
>>>>>>> 867fbedd06c667ac880ebe82b5a18eddc059ec33
