# Tuwi
A TUI for managing WiFi on Linux

<p>
    <img src="./assets/demo.gif" alt="tuwi demo" width="800">
</p>

## Features
- **Browse available networks** — scan and view nearby WiFi networks with signal strength, security type and frequency
- **Connect to networks** — connect to open, secured (WPA/WPA2/WPA3) and hidden networks
- **Manage saved connections** — view, connect to and forget saved WiFi profiles
- **Real-time updates** — network list updates automatically via NetworkManager D-Bus signals without manual scanning
- **Security detection** — accurately identifies WEP, WPA, WPA2, WPA3, WPA2/WPA3 transition, Enterprise and OWE networks
- **Keyboard driven** — fully navigable without a mouse

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
