#!/bin/bash -eu
#
# Copyright IBM Corp. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0

GOROOT='/opt/go'
<<<<<<< HEAD
GO_VERSION=1.16.7
=======
GO_VERSION=1.18.7
>>>>>>> 867fbedd06c667ac880ebe82b5a18eddc059ec33

# ----------------------------------------------------------------
# Install Golang
# ----------------------------------------------------------------
GO_URL=https://storage.googleapis.com/golang/go${GO_VERSION}.linux-amd64.tar.gz
mkdir -p $GOROOT
curl -sL "$GO_URL" | (cd $GOROOT && tar --strip-components 1 -xz)

# ----------------------------------------------------------------
# Setup environment
# ----------------------------------------------------------------
cat <<EOF >/etc/profile.d/goroot.sh
export GOROOT=$GOROOT
export PATH=\$PATH:$GOROOT/bin
EOF
