---
title: "Installation"
linkTitle: "Installation"
weight: 2
description: >
    Installing kluctl.
---

## Additional tools needed

You need to install a set of command line tools to fully use kluctl. These are:

1. [helm](https://github.com/helm/helm/releases):
   Download/Install the binaries for your system, make them executable and make them globally available
   (modify your PATH or copy it into /usr/local/bin). The version of Helm required depends on the Helm Charts that you
   plan to use. If these have a minimum version of Helm, this becomes the minimum version for kluctl.

All of these tools must be in your PATH, so that kluctl can easily invoke them.

## Installation of kluctl

kluctl can currently only be installed this way:
1. Download a standalone binary from the latest release and make it available in your PATH, either by copying it into `/usr/local/bin` or by modifying the PATH variable

Future releases will include packaged releases for homebrew and other established package managers (contributions are welcome).
