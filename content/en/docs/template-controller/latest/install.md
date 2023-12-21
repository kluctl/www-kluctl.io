---
description: Installation documentation
github_repo: https://github.com/kluctl/template-controller
lastmod: "2023-05-17T23:13:11+02:00"
path_base_for_github_subdir:
    from: .*
    to: main/docs/install.md
title: Installation
weight: 10
---





The Template Controller can currently only be installed via kustomize:

```sh
kubectl create ns kluctl-system
kustomize build "https://github.com/kluctl/template-controller/config/install?ref=v0.7.1" | kubectl apply -f-
```

## Helm
A Helm Chart for the controller is also available [here](https://github.com/kluctl/charts/tree/main/charts/template-controller).
To install the controller via Helm, run:
```shell
$ helm repo add kluctl https://kluctl.github.io/charts
$ helm install template-controller kluctl/template-controller
```
