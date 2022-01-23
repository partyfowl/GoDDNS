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


if [[ -n "${GITHUB_REF_NAME}" ]]
then
  if [[ "${GITHUB_REF_NAME}" == "main" ]]
  then
    mv goddns.deb goddns-$(dpkg-deb -f goddns.deb Version).deb
  else
    mv goddns.deb goddns-$(dpkg-deb -f goddns.deb Version)-${GITHUB_REF_NAME}.deb
  fi
fi
