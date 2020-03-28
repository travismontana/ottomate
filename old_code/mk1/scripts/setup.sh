#!/bin/bash

# Usage:
# ./setup.sh <full path to ottomate/src>

if [ "$#" == "0" ]; then
  echo "Not enough arguments given."
  exit 1
fi  

# First, where am I 
# (I'll need you to tell me where ottomate/src is in reference to the ${CWD})
STARTINGDIR=$1
STARTINGWITHOUTSRC=`echo ${STARTINGDIR} | sed -e 's;/src;;'`

# Now where is go
GOBINFULL=`which go`
GOBINDIR=${GOBINFULL%/*}
GOINSTALLDIR=`echo ${GOBINDIR} | sed -e 's;/bin;;'`


# Setup variables
echo "GOPATH=${STARTINGWITHOUTSRC}"
export GOPATH=${STARTINGWITHOUTSRC}

echo "GOROOT=${GOINSTALLDIR}"
export GOROOT=${GOINSTALLDIR}

# Walk through the subdir's in src and go get what we need
for DIR in $(find ${STARTINGDIR}/* -type d); do
  cd ${DIR}
  go get ./...
done
#
# Install the required stuff 

# What package manger do we use here?

for PKGMGR in yum apt-get; do
  which ${PKGMGR} > /dev/null 2>&1
  if [ "x$?" == "x0" ]; then
    PKGR=${PKGMGR}
  fi
done
#
echo "${PKGR}"
