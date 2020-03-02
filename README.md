# concourse crds

A set of custom resource definitions that aim at controlling Concourse.

ps.: you should not take this repository seriously - the purpose of this thing
is to demonstrate what we could achieve with
https://github.com/concourse/rfcs/pull/44.


### usage

1. build the binary


```console
$ make
```

2. install the crds

```console
$ make install
```

ps.: assumes you have [`kustomize`] set up
pps.: assumes `kubeconfig` is set up

[`kustomize`]: https://github.com/kubernetes-sigs/kustomize


3. run it locally

ps.: will target the apiserver configured in `kubeconfig`.

```console
$ make run
```

4. create a pipeline object


```console
$ make pipeline
```

ps.: assume `concourse` runs on `localhost:8080` (you can do it so following
https://github.com/concourse/concourse-docker).


