# Go Dynamic DNS Client

## Prerequisites

- CloudFlare Account with registered domain
- Go (tested with 1.13.8)

## Build

To build a debian package:

```bash
./build.sh
```


## Install

Copy the generated `build/goddns.deb` file to the Raspberry Pi.

Run the following commands:

```bash
sudo apt-get install ./goddns.deb
```

Update `/home/goddns/.goddns/config.yaml` to include the CloudFlare API Key and Hosted Zone ID.

Finally, enable the service:

```bash
sudo systemctl enable goddns --now
```
