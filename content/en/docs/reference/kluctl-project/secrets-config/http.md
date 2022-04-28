---
title: "http"
linkTitle: "http"
weight: 3
description: >
  Loads secrets from an HTTP resource.
---

The http secrets source allows to load secrets from an arbitrary HTTP resource by performing a GET (or any other
configured HTTP method) on the URL. Example:

```yaml
secretsConfig:
  secretSets:
    - name: prod
      sources:
        - http:
            url: https://example.com/path/to/my/secrets
```

The above secrets source will load a [secrets file]({{< ref "docs/reference/kluctl-project/secrets-config#format-of-secrets-files" >}}) 
from the given URL and make it available to the templating context while sealing.

The following additional properties are supported for http sources:

### method
Specifies the HTTP method to be used when requesting the given resource. Defaults to `GET`.

### body
The body to send along with the request. If not specified, nothing is sent.

### headers
A map of key/values pairs representing the header entries to be added to the request. If not specified, nothing is added.

### jsonPath
Can be used to select a nested element from the yaml/json document returned by the HTTP request. This is useful in case
some REST api is used which does not directly return the secrets file. Example:

```yaml
secretsConfig:
  secretSets:
    - name: prod
      sources:
        - http:
            url: https://example.com/path/to/my/secrets
            jsonPath: $[0].data
```

The above example would successfully use the following json document as secrets source:

```json
[{"data": {"secrets": {"secret": "value1"}}}]
```

### Authentication

Kluctl currently supports BASIC and NTLM authentication. It will prompt for credentials when needed.
