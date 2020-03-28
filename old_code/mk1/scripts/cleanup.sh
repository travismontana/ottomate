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

rm -rf ${STARTINGWITHOUTSRC}/bin/*
rm -rf ${STARTINGWITHOUTSRC}/pkg
rm -rf ${STARTINGWITHOUTSRC}/src/github.com
rm -rf ${STARTINGWITHOUTSRC}/src/golang.org