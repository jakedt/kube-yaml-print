# kube-yaml-print

This utility takes in a big paginated set of yaml kubernetes resources and prints the results as a tree.

## Example

`kustomize build overlays/dev | kube-yaml-print`

```sh
.
├── v1/Secret
│   └── a-secret
├── cert-manager.io/v1/ClusterIssuer
│   └── letsencrypt
└── [namespace]  app
    ├── v1/Service
    │   └── app
    ├── v1/Namespace
    │   └── app
    └── apps/v1/Deployment
        └── app
```
