# slack-operator

A Helm chart to deploy slack-operator

## Pre-requisites

- Make sure that [certman](https://cert-manager.io/) is deployed in your cluster since webhooks require certman to generate valid certs since webhooks serve using HTTPS

```terminal
$ kubectl apply --validate=false -f https://github.com/jetstack/cert-manager/releases/download/v1.0.1/cert-manager.yaml
```

## Installing the chart

```sh
helm repo add stakater https://stakater.github.io/stakater-charts/
helm repo update
helm install stakater/slack-operator --namespace slack-operator
```