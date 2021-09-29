# Hollow CoreDHCP Plugin & DHCP Server Build

[![codecov](https://codecov.io/gh/metal-toolbox/hollow-dhcpserver/branch/main/graph/badge.svg?token=xXNVOjjWJ7)](https://codecov.io/gh/metal-toolbox/hollow-dhcpserver)


This provides a plugin that can be used for serving dhcp request from hollow data sources.

## How to use

You will need to store the DHCP information on a server or instance attribute with the namespace `sh.hollow.dhcpserver.lease`.

The format of this data needs to look like:

```json
{
  "ipv4": [
    {
      "boot_file": "ipxe.efi",
      "boot_server": "10.0.0.1",
      "cidr": "10.0.0.20/24",
      "gateway": "10.0.0.1",
      "mac_address": "02:42:ac:13:00:05",
      "resolvers": [
        "1.1.1.1",
        "8.8.8.8"
      ]
    },
    {
      "boot_file": "",
      "boot_server": "",
      "cidr": "10.0.0.100/24",
      "gateway": "10.0.0.1",
      "mac_address": "02:42:ac:13:00:01",
      "resolvers": [
        "1.1.1.1",
        "8.8.8.8"
      ]
    }
  ]
}
```
