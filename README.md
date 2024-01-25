# Cod(e)Only

Codonly helps to find resources that are not part of infrastructure management tooling.

:warning: Don't use this for any serious use-cases.

## How does codonly work?

1. Read and parse terraform state
2. List resources (databases, k8s clusters, buckets, ....) from third parties
3. Check what resources are not part of the state file
4. Report