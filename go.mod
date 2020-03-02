module github.com/cirocosta/crds

go 1.13

require (
	github.com/caarlos0/env/v6 v6.2.1
	github.com/concourse/concourse v1.6.1-0.20200228185904-aac7b2f461ed
	github.com/go-logr/logr v0.1.0
	github.com/onsi/ginkgo v1.11.0
	github.com/onsi/gomega v1.8.1
	golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45
	k8s.io/apimachinery v0.17.2
	k8s.io/client-go v11.0.0+incompatible
	sigs.k8s.io/controller-runtime v0.5.0
)

replace k8s.io/client-go v11.0.0+incompatible => k8s.io/client-go v0.17.2
