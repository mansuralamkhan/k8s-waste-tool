# K8s Waste Detector

A lightweight, self-hosted tool that shows you exactly how much idle capacity
(and money) your Kubernetes cluster is wasting — broken down by
namespace/team, with specific, actionable recommendations.

Built on top of [OpenCost](https://www.opencost.io/) for cost/usage data.

## Why this exists

Existing cost tools (like Kubecost) are powerful but complex and can get
expensive at scale. This is a simpler, cheaper, self-hosted alternative for
small teams who just want to know: what's idle, and what should I change?

## Prerequisites

- A Kubernetes cluster you can deploy to (tested on [kind](https://kind.sigs.k8s.io/))
- [Prometheus](https://github.com/prometheus-community/helm-charts) installed in-cluster
- [OpenCost](https://www.opencost.io/docs/installation/helm) installed in-cluster,
  pointed at that Prometheus instance
- `kubectl` and `helm` installed locally, configured against your cluster
- Docker, if building the image yourself (no public image published yet)

## Quickstart (local development)

1. Build the image:

       docker build -f deploy/docker/backend.Dockerfile -t waste-detector:dev .

2. Make the image available to your cluster. How you do this depends on
   your cluster type:

   **kind** (Kubernetes-in-Docker, e.g. on Docker Desktop):

       kind load docker-image waste-detector:dev --name <your-cluster-name>

   **A cluster using containerd directly** (e.g. k3s, Killercoda,
   most managed single-node labs) — Docker and the cluster's runtime
   don't share an image cache, so import it manually:

       docker save waste-detector:dev -o waste-detector.tar
       ctr -n=k8s.io images import waste-detector.tar

   Then set `imagePullPolicy: Never` when installing (see step 3), so
   Kubernetes uses the locally-imported image instead of trying to pull.

   **A real cluster with a container registry** (EKS, GKE, AKS, or any
   cluster pulling from ECR/GCR/Docker Hub):

       docker tag waste-detector:dev <your-registry>/waste-detector:v0.1.0
       docker push <your-registry>/waste-detector:v0.1.0

   Then set `image.repository` and `image.tag` in `deploy/helm/values.yaml`
   (or via `--set` at install time) to match.

3. Install the chart:

       helm install waste-detector ./deploy/helm --namespace <your-namespace>

   If you used the containerd import method in step 2, add
   `--set image.pullPolicy=Never` so Kubernetes doesn't attempt a
   registry pull:

       helm install waste-detector ./deploy/helm --namespace <your-namespace> \
         --set image.pullPolicy=Never

4. Access the dashboard:

       kubectl port-forward -n <your-namespace> svc/waste-detector 8080:8080

   Then open http://localhost:8080

## Status

Early development — not yet ready for production use, no published Helm
repo or container registry image yet. See [docs/prd.md](docs/prd.md).

## Docs

- [Product requirements](docs/prd.md)
- [Architecture decisions](docs/adr)