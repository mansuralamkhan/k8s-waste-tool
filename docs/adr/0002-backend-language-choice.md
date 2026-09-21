# 0002: Use Go for the backend

## Status
Accepted

## Context
We need to choose a backend language before building the waste-calculation
logic and API layer. Main candidates considered: Go and Python (FastAPI).

Key factors:
- Distribution matters a lot for this product — "one command to see value"
  (ADR 0003) means the backend ideally ships as a small, dependency-free
  single binary that runs cleanly inside a Helm-installed pod.
- The ecosystem we're integrating with (Kubernetes client libraries,
  OpenCost itself, most CNCF tooling) is Go-native, so client libraries and
  examples are more direct to work with.
- Tradeoff accepted: less prior familiarity with Go means a slower initial
  build velocity than Python would offer, in exchange for a better fit with
  the deployment model and ecosystem long-term.

## Decision
Backend will be written in Go.

## Consequences
- Smaller, faster container images and simpler Helm chart packaging (no
  Python runtime/dependency management inside the cluster).
- Better alignment with OpenCost's own Go client libraries and Kubernetes'
  native tooling.
- Expect a slower first few tickets while ramping up on Go idioms — worth
  budgeting extra time for the earliest scaffolding steps rather than
  assuming Python-speed velocity.
- Frontend/dashboard remains a separate decision (likely TypeScript/React),
  since Go's strength here is backend/API and CLI tooling, not UI.