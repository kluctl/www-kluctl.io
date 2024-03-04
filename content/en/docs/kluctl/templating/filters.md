---
description: Available filters.
github_repo: https://github.com/kluctl/kluctl
lastmod: "2024-03-04T11:05:56+01:00"
linkTitle: Filters
path_base_for_github_subdir:
    from: .*
    to: main/docs/kluctl/templating/filters.md
title: Filters
weight: 3
---





In addition to the [builtin Jinja2 filters](https://jinja.palletsprojects.com/en/3.1.x/templates/#list-of-builtin-filters),
kluctl provides a few additional filters:

### b64encode
Encodes the input value as base64. Example: `{{ "test" | b64encode }}` will result in `dGVzdA==`.

### b64decode
Decodes an input base64 encoded string. Example `{{ my.source.var | b64decode }}`.

### from_yaml
Parses a yaml string and returns an object. Please note that json is valid yaml, meaning that you can also use this
filter to parse json.

### to_yaml
Converts a variable/object into its yaml representation. Please note that in most cases the resulting string will not
be properly indented, which will require you to also use the `indent` filter. Example:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: my-config
data:
  config.yaml: |
    {{ my_config | to_yaml | indent(4) }}
```

### to_json
Same as `to_yaml`, but with json as output. Please note that json is always valid yaml, meaning that you can also use
`to_json` in yaml files. Consider the following example:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-deployment
spec:
  template:
    spec:
      containers:
      - name: c1
        image: my-image
        env: {{ my_list_of_env_entries | to_json }}
```

This would render json into a yaml file, which is still a valid yaml file. Compare this to how this would have to be
solved with `to_yaml`:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-deployment
spec:
  template:
    spec:
      containers:
      - name: c1
        image: my-image
        env:
          {{ my_list_of_env_entries | to_yaml | indent(10) }}
```

The required indention filter is the part that makes this error-prone and hard to maintain. Consider using `to_json`
whenever you can.

### render
Same as the global [render function](./functions.md#rendertemplate), but deprecated now. `render` being a filter turned out to
not work well with local variables, as these are not accessible in filters. Please only use the global function.

### sha256(digest_len)
Calculates the sha256 digest of the input string. Example:
```
{{ "some-string" | sha256 }}
```

`digest_len` is an optional parameter that allows to limit the length of the returned hex digest. Example:
```
{{ "some-string" | sha256(6) }}
```

### slugify
Slugify a string based on [python-slugify](https://github.com/un33k/python-slugify).
