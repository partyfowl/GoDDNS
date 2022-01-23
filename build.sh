set -e

rm -rf build
mkdir build

if [[ -z "${GITHUB_RUN_NUMBER}" ]]; then
    PACKAGE_NAME=goddns-${GITHUB_RUN_NUMBER}.deb
    sed -i "s/^\(Version:\s*\).*$/\1${GITHUB_RUN_NUMBER}/" packaging/DEBIAN/control
else
    PACKAGE_NAME=goddns.deb
fi

chmod 775 -R packaging/DEBIAN/*inst

cp -r packaging build/$PACKAGE_NAME

mkdir -p build/$PACKAGE_NAME/usr/local/bin/
go get
GOOS=linux GOARCH=arm go build -o build/$PACKAGE_NAME/usr/local/bin/

cd build
dpkg-deb --build $PACKAGE_NAME
