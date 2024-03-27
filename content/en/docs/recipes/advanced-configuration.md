---
title: "Advanced configuration"
linkTitle: "Advanced configuration"
weight: 30
---

This recipe will try to give best practices on how to achieve advanced configuration that keeps being maintainable.

## Args as entrypoint

Kluctl offers multiple ways to introduce configuration args into your deployment. These are all accessible via
[templating]({{% ref "docs/kluctl/templating" %}}) by referencing the global `args` variable, e.g. `{{ args.my_arg }}`.

Args can be passed via [command line arguments]({{% ref "docs/kluctl/commands/common-arguments#project-arguments" %}}),
[target definitions]({{% ref "docs/kluctl/kluctl-project/targets#args" %}}) and GitOps
[KluctlDeployment spec]({{% ref "docs/gitops/spec/v1beta1/kluctldeployment#args" %}}).

It might however be tempting to provide all necessary configuration via args, which can easily end up clogging things up
in a very unmaintainable way.

## Combining args with vars sources

The better and much more maintainable approach is to combine `args` with 
[variable sources]({{% ref "docs/kluctl/templating/variable-sources" %}}). You could for example
introduce an arg that is later used to load further configuration from YAML files or even external vars sources (e.g. git).

Consider the following example:

```yaml
# .kluctl.yaml
targets:
  - name: prod
    context: prod.example.com
    args:
      environment_type: prod
      environment_name: prod
  - name: test
    context: test.example.com
    args:
      environment_type: non-prod
      environment_name: test
  - name: dev
    context: test.example.com
    args:
      environment_type: non-prod
      environment_name: dev
```

```yaml
# root deployment.yaml
vars:
  - file: config/{{ args.environment_type }}.yaml

deployments:
  - include: my-include
  - path: my-deployment
```

The above `deployment.yaml` will load different configuration, depending on the passed `environment_type` argument.

This means, you'll also need the following configuration files:

```yaml
# config/prod.yaml
myApp:
  replicas: 3
```

```yaml
# config/non-prod.yaml
myApp:
  replicas: 1
```

This way, you don't have to bloat up the `.kluctl.yaml` with some ever-growing amount of configuration but instead can
move such configuration into dedicated configuration files.

The resulting configuration can then be used via templating, e.g. `{{ myApp.replicas }}`

## Layering configuration on top of each other

Kluctl merges already loaded configuration with freshly loaded configuration. It does this for every item in `vars`.
At the same time, Kluctl allows to use templating with the previously loaded configuration context in each loaded
vars source. This means, that configuration that was loaded by a vars item before the current one can already be used
in the current one.

All [deployment items]({{% ref "docs/kluctl/deployments/deployment-yml#deployments" %}}) will then be provided with the
final merged configuration. If deployment items also define [vars]({{% ref "docs/kluctl/deployments/deployment-yml#vars-deployment-item" %}}),
these are merged as well, but only for the context of the specific deployment item.

Consider the following example:

```yaml
# root deployment.yaml
vars:
  - file: config/common.yaml
  - file: config/{{ args.environment_type }}.yaml
  - file: config/monitoring.yaml
```

```yaml
# config/common.yaml
myApp:
  monitoring:
    enabled: false
```

```yaml
# config/prod.yaml
myApp:
  replicas: 3
  monitoring:
    enabled: true
```

```yaml
# config/non-prod.yaml
myApp:
  replicas: 1
```

The merged configuration for `prod` environments will have `myApp.monitoring.enabled` set to `true`, while all other
environments will have it set to `false`.

## Putting configuration into the target cluster

Kluctl supports many different [variable sources]({{% ref "docs/kluctl/templating/variable-sources" %}}), which means
you are not forced to store all configuration in files which are part of the project.

You can also store configuration inside the target cluster and access this configuration via the
[clusterConfigMap or clusterSecret]({{% ref "docs/kluctl/templating/variable-sources#clusterconfigmap" %}}) variable
sources. This configuration could for example be part of the cluster provisioning stage and contain information about
networking info, cloud info, DNS info, and so on, so that this can then be re-used wherever needed (e.g. in ingresses).

Consider the following example ConfigMap, which was already deployed to your target cluster:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: cluster-info
  namespace: kube-system
data:
  vars: |
    clusterInfo:
      baseDns: test.example.com
      aws:
        accountId: 12345
        irsaPrefix: test-example-com
```

Your deployment:

```yaml
# root deployment.yaml
vars:
  - clusterConfigMap:
      name: cluster-info
      namespace: kube-system
      key: vars
  - file: ... # some other configuration, as usual

deployments:
  # as usual
  - ...
```

```yaml
# some/example/ingress.yaml
# look at the DNS name
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: my-ingress
  namespace: my-namespace
spec:
  rules:
    - host: my-ingress.{{ clusterInfo.baseDns }}
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: my-service
                port:
                  number: 80
  tls:
    - hosts:
        - 'my-ingress.{{ clusterInfo.baseDns }}'
      secretName: 'ssl-cert'
```

```yaml
# some/example/irso-service-account.yaml
# Assuming you're using IRSA (https://docs.aws.amazon.com/eks/latest/userguide/iam-roles-for-service-accounts.html)
# for external-dns
apiVersion: v1
kind: ServiceAccount
metadata:
  name: external-dns
  annotations:
    eks.amazonaws.com/role-arn: arn:aws:iam::{{ clusterInfo.aws.accountId }}:role/{{ clusterInfo.aws.irsaPrefix }}-external-dns
```
