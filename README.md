# slack-operator

Kubernetes operator for Slack

## About

Slack operator is used to automate the process of setting up a Slack channel for alertmanager in a k8s native way. By using CRDs it lets you:

1. Manage Channels
2. Configure Issues

It uses [Slack Api](https://api.slack.com/methods) in it's underlying layer and can be extended to perform other tasks that are supported via the REST API.

## Usage

### Prerequisites

- Slack account
- API Token to access Slack API (https://api.slack.com/)

### Create secret

Create the following secret which is required for slack-operator:

```yaml
kind: Secret
apiVersion: v1
metadata:
  name: slack-secret
type: Opaque
data:
  APIToken: <SLACK_API_TOKEN>
```

### Deploy operator

- Make sure that [certman](https://cert-manager.io/) is deployed in your cluster since webhooks require certman to generate valid certs since webhooks serve using HTTPS
- To install certman

```terminal
$ kubectl apply --validate=false -f https://github.com/jetstack/cert-manager/releases/download/v1.3.1/cert-manager.yaml
```

- Deploy operator

```terminal
$ oc apply -f bundle/manifests
```

## Local Development

- [Operator-sdk v1.7.2](https://github.com/operator-framework/operator-sdk/releases/tag/v1.7.2) is required for local development.

1. Create `slack-secret` secret
2. Run `make run ENABLE_WEBHOOKS=false WATCH_NAMESPACE=default OPERATOR_NAMESPACE=default` where `WATCH_NAMESPACE` denotes the namespaces that the operator is supposed to watch and `OPERATOR_NAMESPACE` is the namespace in which it's supposed to be deployed.

3. Before committing your changes run the following to ensure that everything is verified and up-to-date:
   - `make verify`

## Running Tests

### Pre-requisites:

1. Create a namespace with the name `test`
2. Create `slack-secret` secret in test namespace

### To run tests:

Use the following command to run tests:
`make test OPERATOR_NAMESPACE=test USE_EXISTING_CLUSTER=true`
