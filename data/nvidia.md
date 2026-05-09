# Nvidia Driver Installation

Helpers for installing Nvidia GPU drivers on Debian-family distros.

## Install headers
`sudo apt install linux-headers-$(uname -r)`

Install kernel headers matching the running kernel — required before building proprietary modules.
