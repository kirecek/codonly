# Cod(e)Only

Simple tool that helps to find resources that are not part of infrastructure management tooling.

:warning: Don't use this for any serious use-cases. I have no idea if this is bullshit or not.

Currently it supports only 1 state, terraform and limited set of gcloud resources.

TODOs:
- more states
- more resources
- more clouds/providers
- more IaaC tools
- reporters (stdout, slack, ...)
- connect to git

## How does codonly work?

1. Read and parse terraform state
2. List resources (databases, k8s clusters, buckets, ....) from third parties
3. Check what resources are not part of the state file
<<<<<<< HEAD
4. Report
=======
4. Report
>>>>>>> a31946d (Add initial terraform and gcp provider)
