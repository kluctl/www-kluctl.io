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

### Install kubectl binary with curl on Linux

1. Download the latest release with the command:

   ```bash
   curl -LO "https://github.com/kluctl/kluctl/releases/download/$(curl --silent -H "Accept: application/vnd.github.v3+json" https://api.github.com/repos/kluctl/kluctl/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')/kluctl-linux-amd64" 
   ```

   {{< note >}}

   To download a specific version, replace the `$(curl --silent -H "Accept: application/vnd.github.v3+json" https://api.github.com/repos/kluctl/kluctl/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')` portion of the command with the specific version.
   For example, to download version {{< param "fullversion" >}} on Linux, type:

   ```bash
   curl -LO https://github.com/kluctl/kluctl/releases/download/{{< param "fullversion" >}}/kluctl-linux-amd64
   ```
   {{< /note >}}

3. Install kubectl

   ```bash
   sudo install -o root -g root -m 0755 kluctl-linux-amd64 /usr/local/bin/kluctl
   ```

   {{< note >}}
   If you do not have root access on the target system, you can still install kluctl to the `~/.local/bin` directory:

   ```bash
   chmod +x kluctl-linux-amd64
   mkdir -p ~/.local/bin
   mv ./kluctl-linux-amd64 ~/.local/bin/kluctl
   # and then append (or prepend) ~/.local/bin to $PATH
   ```

   {{< /note >}}

4. Test to ensure the version you installed is up-to-date:

   ```bash
   kluctl version
   ```