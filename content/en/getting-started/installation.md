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

### Install kubectl binary with curl on Linux

1. Download the latest release with the command:

   ```bash
   # To download a specific version, replace the kluctl_version with the version you need for example kluctl_version={{< param "fullversion" >}}
   kluctl_version=$(curl --silent -H "Accept: application/vnd.github.v3+json" https://api.github.com/repos/kluctl/kluctl/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
   curl -LO "https://github.com/kluctl/kluctl/releases/download/${kluctl_version}/kluctl-linux-amd64"
   ```

2. **Optional**: Validate the kluctl binary *only available as of version >= v2.7.0*

   Download the kluctl checksum file:
   ```bash
   kluctl_version=$(curl --silent -H "Accept: application/vnd.github.v3+json" https://api.github.com/repos/kluctl/kluctl/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
   curl -LO "https://github.com/kluctl/kluctl/releases/download/${kluctl_version}/checksums.txt"
   ```
   Validate the kluctl binary against the checksum file:
   ```
   echo $(grep kluctl-linux-amd64 checksums.txt) | sha256sum --check
   ```
   If valid, the output is:
   ```
   kluctl-linux-amd64: OK
   ```

3. Install kubectl

   ```bash
   sudo install -o root -g root -m 0755 kluctl-linux-amd64 /usr/local/bin/kluctl
   ```

   {{< alert >}}
   If you do not have root access on the target system, you can still install kluctl to the `~/.local/bin` directory:

   ```bash
   chmod +x kluctl-linux-amd64
   mkdir -p ~/.local/bin
   mv ./kluctl-linux-amd64 ~/.local/bin/kluctl
   # and then append (or prepend) ~/.local/bin to $PATH
   ```
   {{< /alert >}}

4. Test to ensure the version you installed is up-to-date:

   ```bash
   kluctl version
   ```
