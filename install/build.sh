#!/bin/bash

set -e # exit on error

if [ "$#" -ne 2 ]; then
  echo "Usage: $0 <architecture> <version>"
  exit 1
fi

ARCH=$1
VERSION=$2
PACKAGE=pvault_${ARCH}_${VERSION}

## make the package
mkdir -p $PACKAGE/DEBIAN

## copy and update control info
echo "|> Setting control info (Architecture: ${ARCH}, Version: ${VERSION})"
cp config/control $PACKAGE/DEBIAN/control

sed -i "0,/^Architecture:/s/^Architecture:.*/Architecture: ${ARCH}/" ${PACKAGE}/DEBIAN/control
sed -i "0,/^Version:/s/^Version:.*/Version: ${VERSION}/" ${PACKAGE}/DEBIAN/control

## build app
echo "|> Building app (Architecture: ${ARCH})"
GOOS=linux GOARCH=$ARCH go build -C .. -ldflags="-s -w" -o pvault_${VERSION}

mkdir -p $PACKAGE/usr/local/bin
mv ../pvault_${VERSION} $PACKAGE/usr/local/bin/pvault

## build deb
dpkg-deb --build $PACKAGE
