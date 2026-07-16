# Tuwi
A TUI for managing WiFi on Linux

<p>
    <img src="./assets/demo.gif" alt="tuwi demo" width="800">
</p>

## Features
- Manage saved WiFi connections.
- Scan, browse and connect to nearby WiFi networks.
- Support hidden networks.

## Prerequisites
- A Linux based OS.
- NetworkManager running.

> [!INFO] Support for iwd backend has not been added yet.

## Installation
### Binary releases (recommended)
Download a binary from the Releases page [here](https://github.com/thecentinol/tuwi/releases)

### Build from source
- Make sure [Go](https://go.dev/doc/install) is installed.

```bash
git clone https://github.com/thecentinol/tuwi
cd tuwi
make install
```
`make install` installs the `tuwi` binary into `/usr/bin`. To install to a different location, edit the `PREFIX` and `BINDIR` variables in the [Makefile](https://github.com/thecentinol/tuwi/Makefile).

## Usage
Run:
```bash
tuwi
```
Key bindings are displayed at the bottom of the interface.
