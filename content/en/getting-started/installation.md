---
title: "Installation"
linkTitle: "Installation"
weight: 2
description: >
    Installing kluctl.
---

## Additional tools needed

You need to install a set of command line tools to fully use kluctl. These are:

1. [helm](https://helm.sh/). Follow the [installation guide](https://helm.sh/docs/intro/install/) of Helm to install it.
   The version of Helm required depends on the Helm Charts that you
   plan to use. If these have a minimum version of Helm, this becomes the minimum version for kluctl.

All of these tools must be in your PATH, so that kluctl can easily invoke them.

## Install kluctl

The kluctl CLI is available as a binary executable for all major platforms,
the binaries can be downloaded form GitHub
[releases page](https://github.com/kluctl/kluctl/releases).

{{% tabs %}}
{{% tab "Homebrew" %}}

With [Homebrew](https://brew.sh) for macOS and Linux:

```sh
brew install kluctl/tap/kluctl
```

{{% /tab %}}
{{% tab "bash" %}}

With [Bash](https://www.gnu.org/software/bash/) for macOS and Linux:

```sh
curl -s https://kluctl.io/install.sh | bash
```

{{% /tab %}}

<!-- TODO uncomment when chocolatey support is implemented
{{% tab "Chocolatey" %}}

With [Chocolatey](https://chocolatey.org/) for Windows:

```powershell
choco install kluctl
```

{{% /tab %}}
-->
{{% /tabs %}}

<!-- TODO uncomment this when completion is implemented
To configure your shell to load `kluctl` [bash completions](./cmd/kluctl_completion_bash.md) add to your profile:

```sh
. <(kluctl completion bash)
```

[`zsh`](./cmd/kluctl_completion_zsh.md), [`fish`](./cmd/kluctl_completion_fish.md),
and [`powershell`](./cmd/kluctl_completion_powershell.md)
are also supported with their own sub-commands.

-->

## Container images

A container image with `helm` and `kluctl` is available on GitHub:

* `ghcr.io/kluctl/kluctl:<version>`
