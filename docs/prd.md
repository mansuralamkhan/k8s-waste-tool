Product Requirements Doc — K8s Waste Detector (working name)

Problem statement
Small-to-mid-sized engineering teams running Kubernetes (EKS/GKE) don't have visibility into how much money they're wasting on over-provisioned pods and underutilized nodes. Existing tools like Kubecost solve this but are complex to set up and price out of reach for teams without enterprise budgets (reports of $70k+/year for multi-cluster monitoring alone).

Target user
A DevOps/platform engineer or engineering lead at a startup or small company (roughly 5-50 engineers) already running Kubernetes in production, without a dedicated FinOps team or budget for enterprise cost tooling.

Core user flow (v1)

User installs the tool via a single Helm command on their cluster.
Within minutes, a dashboard shows total monthly spend, broken into "used" vs "wasted," by namespace/team.
User sees a short list of specific, dollar-quantified recommendations (e.g., "reduce X's memory request, save $Y/month").
User acts on a recommendation manually (v1 doesn't auto-apply changes).

Out of scope for v1 (explicitly, to prevent scope creep)

Multi-cloud support (AWS/EKS only first)
Automated remediation (recommend only, don't auto-change configs)
Historical trend charts beyond a basic time range
Team accounts / multi-user permissions
Alerting (Slack/email) — add after core dashboard is validated
Anything beyond CPU/memory waste (no storage/network cost analysis yet)

Success metrics for v1

5 real installs from people outside your own testing
At least 2 of those return to check the dashboard a second time without being prompted
At least 1 person says they changed a resource request/limit because of a recommendation
At least 1 person says explicitly what they'd pay for this