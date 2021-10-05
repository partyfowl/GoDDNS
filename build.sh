rm -rf build
mkdir build

cp -r packaging build/goddns

mkdir -p build/goddns/usr/local/bin/
GOOS=linux GOARCH=arm go build -o build/goddns/usr/local/bin/

cd build
dpkg-deb --build goddns
