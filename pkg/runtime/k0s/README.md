# k0s Runtime Design
## basefs

```shell
.
├── amd64
│   ├── bin
│   │   ├── k0s
│   │   ├── kubectl
│   │   ├── nerdctl
│   │   └── seautil
│   ├── images
│   │   └── registry.tar.gz
│   └── Metadata
├── imageList
├── Kubefile
└── rootfs
    ├── etc
    │   ├── dump-config.toml
    │   └── registry.yml
    └── scripts
        ├── containerd.sh
        ├── init-registry.sh
        └── init.sh
```

## introduce
We define the k0s runtime has 5 phases to install/scale/reset the cluster.

basefs contains binary、shell script、config file and image. See more about [sealerio/basefs](https://github.com/sealerio/basefs)

Runtime before filesystem lead cluster install through execute [k0s](https://github.com/k0sproject/k0s) command.

+ Init
  + When sealer leads to cluster install first, init phase copy rootfs/bin to /usr/bin in init.sh script
  + create bootstrap config /etc/k0s/k0s.yaml to lead controller init
  + generate k0s join token /etc/k0s/worker-token and /etc/k0s/controller-token, also private registry cert
  + init controller node
  + fetch config for ~/.kube/config to manage the cluster.
+ Join
  + Join phase prepare the registry certs, and use `k0s join` to scale up the cluster.
+ Delete
  + Delete is same as join, but it recycles anything installed by join phase.
+ Reset
  + Reset through k0s to remove this cluster and remove anything about the cluster generated by sealer.
+ Upgrade
  + Upgrade can upgrade the k0s cluster.