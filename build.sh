set -e

rm -rf build
mkdir build

chmod 775 -R packaging/DEBIAN/*inst

cp -r packaging build/goddns

mkdir -p build/goddns/usr/local/bin/
go get
GOOS=linux GOARCH=arm go build -o build/goddns/usr/local/bin/

cd build
dpkg-deb --build goddns
