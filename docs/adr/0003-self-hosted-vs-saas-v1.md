# 0003: Self-hosted and free for v1, not hosted SaaS

## Status
Accepted

## Context
We need to decide how users first get access to the tool: a hosted SaaS
(sign up, we run it, they connect their cluster) vs. self-hosted (they run
it in their own cluster via Helm, no accounts needed).

Validated via user conversations: "free/self-hosted model might be a good
start" — users are more likely to try something with zero signup friction
and no data leaving their cluster (a real trust concern for cost/usage data).

## Decision
v1 ships as a self-hosted Helm chart, free, no accounts or billing system.
Users install it directly into their own cluster and never send data
outside their infrastructure.

## Consequences
- No auth, multi-tenancy, or billing code needed for v1 — meaningfully
  smaller scope.
- Removes the biggest install friction (data leaving the cluster, signup
  flow) which matters for a "one command to see value" wedge against
  Kubecost's more involved setup.
- Willingness to install for free doesn't prove willingness to pay — that
  remains a separate open question to test later, likely via an optional
  hosted/paid tier once free adoption is validated.
- Revisit this decision once there's real usage data; a hosted tier can be
  added later without breaking the self-hosted option.