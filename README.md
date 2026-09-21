
# K8s Waste Detector

A lightweight, self-hosted tool that shows you exactly how much idle capacity
(and money) your Kubernetes cluster is wasting — broken down by namespace/team,
with specific, actionable recommendations.

Built on top of [OpenCost](https://www.opencost.io/) for cost/usage data.

## Why this exists

Existing cost tools (like Kubecost) are powerful but complex and can get
expensive at scale. This is a simpler, cheaper, self-hosted alternative for
small teams who just want to know: what's idle, and what should I change?

## Quickstart

    helm repo add k8s-waste-tool <TBD>
    helm install waste-detector k8s-waste-tool/waste-detector

## Status

Early development — not yet ready for production use. See [docs/prd.md](docs/prd.md).

## Docs

- [Product requirements](docs/prd.md)
- [Architecture decisions](docs/adr/)