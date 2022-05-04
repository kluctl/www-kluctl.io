---
title: "awsSecretsManager"
linkTitle: "awsSecretsManager"
weight: 4
description: >
  Loads secrets from AWS Secrets Manager.
---

[AWS Secrets Manager](https://aws.amazon.com/secrets-manager/) integration. Loads a secrets YAML from an AWS Secrets
Manager secret. The secret can either be specified via an ARN or via a secretName and region combination. An AWS
config profile can also be specified (which must exist while sealing).

The secrets stored in AWS Secrets manager must contain
a valid [secrets file]({{< ref "docs/reference/kluctl-project/secrets-config#format-of-secrets-files" >}}).

Example using an ARN:
```yaml
secretsConfig:
  secretSets:
    - name: prod
      sources:
        - awsSecretsManager:
            secretName: arn:aws:secretsmanager:eu-central-1:12345678:secret:secret-name-XYZ
            profile: my-prod-profile
```

Example using a secret name and region:
```yaml
secretsConfig:
  secretSets:
    - name: prod
      sources:
        - awsSecretsManager:
            secretName: secret-name
            region: eu-central-1
            profile: my-prod-profile
```

The advantage of the latter is that the auto-generated suffix in the ARN (which might not be known at the time of
writing the configuration) doesn't have to be specified.