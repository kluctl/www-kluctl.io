---
description: Running the Kluctl Webui locally
github_branch: main
github_repo: https://github.com/kluctl/kluctl
lastmod: "2024-01-05T11:52:13+01:00"
linkTitle: Running locally
path_base_for_github_subdir:
    from: .*
    to: docs/webui/running-locally.md
title: Running locally
weight: 20
---





The Kluctl Webui can be run locally by simply invoking [`kluctl webui run`](../kluctl/commands/webui-run.md).
It will by default connect to your local Kubeconfig Context and expose the Webui on `localhost`. It will also open
the browser for you.

## Multiple Clusters

The Webui can already handle multiple clusters. Simply pass `--context <context-name>` multiple times to `kluctl webui run`.
This will cause the Webui to listen for status updates on all passed clusters.

As noted in [State of the Webui](./#state-of-the-webui), the Webui is still in early stage and thus currently
lacks sorting and filtering for clusters. This will be implemented in future releases.
