# 0001: Use OpenCost for metrics and cost data

## Status
Accepted

## Context
We need real-time Kubernetes cost/usage data (CPU/memory requested vs actual
usage, cross-referenced with cloud pricing) to calculate idle capacity.
Building this collection layer from scratch would take months and duplicate
work already solved well by an existing open-source project.

## Decision
Use OpenCost (Apache 2.0, CNCF project) as the underlying metrics and pricing
data source. Our backend will query OpenCost's API/Prometheus metrics rather
than collecting raw usage/pricing data ourselves.

## Consequences
- Faster time to v1 — we build on a mature, maintained data layer.
- Our differentiation is the waste-calculation logic, recommendations, and
  simple UX/install experience — not the underlying data collection.
- We take on a dependency on OpenCost's data model and API stability, but
  since it's open source (not a closed vendor API), we retain the option to
  fork or extend it if needed.