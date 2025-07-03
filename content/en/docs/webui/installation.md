---
description: Installing the Kluctl Webui
github_branch: main
github_repo: https://github.com/kluctl/kluctl
lastmod: "2025-07-03T13:24:15+02:00"
linkTitle: Installation
path_base_for_github_subdir:
    from: .*
    to: docs/webui/installation.md
title: Installation
weight: 10
---





The Kluctl Webui can be installed by using a [Git Include](../kluctl/deployments/deployment-yml.md#git-includes) that refers
to the [webui deployment project](https://github.com/kluctl/kluctl/tree/main/install/webui). Example:

```yaml
deployments:
  - git:
      url: https://github.com/kluctl/kluctl.git
      subDir: install/webui
      ref:
        tag: v2.27.0
```

## Login

### Static Users

By default, the Webui will automatically generate an static credentials for an admin and for a viewer user. These
credentials can be extracted from the `kluctl-system/webui-secret` Secret after the Webui has started up for the first
time. To get the admin password, invoke:

```shell
$ kubectl -n kluctl-system get secret webui-secret -o jsonpath='{.data.admin-password}' | base64 -d
```

For the viewer password, invoke:

```shell
$ kubectl -n kluctl-system get secret webui-secret -o jsonpath='{.data.viewer-password}' | base64 -d
```

If you do not want to rely on the Webui to generate those secrets, simply use your typical means of creating/updating
the `webui-secret` Secret. The secret must contain values for `admin-password`, `viewer-password`.

### OIDC Integration

The Webui offers an OIDC integration, which can be configured via [CLI arguments](#passing-arguments).

For an example of an OIDC provider configurations, see [Azure AD Integration](./oidc-azure-ad.md).

## Customization

### Serving under a different path

By default, the webui is served under the `/`path. To change the path, pass the `--prefix-path` argument to the webui:

```yaml
deployments:
  - git:
      url: https://github.com/kluctl/kluctl.git
      subDir: install/webui
      ref:
        tag: v2.27.0
    vars:
      - values:
          webui_args:
            - --path-prefix=/my-custom-prefix
```

### Overriding the version

The image version of the Webui can be overriden with the `kluctl_version` arg:

```yaml
deployments:
  - git:
      url: https://github.com/kluctl/kluctl.git
      subDir: install/webui
      ref:
        tag: main
    vars:
      - values:
          args:
            kluctl_version: v2.27.0
```

### Passing arguments

You can pass arbitrary command line arguments to the webui by providing the `webui_args` arg:

```yaml
deployments:
  - git:
      url: https://github.com/kluctl/kluctl.git
      subDir: install/webui
      ref:
        tag: v2.27.0
    vars:
      - values:
          webui_args:
            - --gops-agent
```
