# Go Dynamic DNS Client

## Prerequisites

- CloudFlare Account with registered domain
- Go (tested with 1.17.6)
- dpkg-deb

## Build

To build a debian package:

```bash
./build.sh
```


## Install

Copy the generated `build/goddns-<version>.deb` file to the Raspberry Pi.

Run the following commands:

```bash
sudo apt-get install ./goddns.deb
```

Update `/home/goddns/.goddns/config.yaml` to include the CloudFlare API Key and Hosted Zone ID.

Finally, enable the service:

```bash
sudo systemctl enable goddns --now
```
