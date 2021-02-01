module github.com/k8ssandra/k8ssandra/tests

go 1.15

require (
	github.com/gruntwork-io/terratest v0.30.23
	github.com/onsi/ginkgo v1.14.2
	github.com/onsi/gomega v1.10.3
	github.com/stretchr/testify v1.6.1
	k8s.io/api v0.20.2
	k8s.io/apimachinery v0.20.2
	sigs.k8s.io/yaml v1.2.0

)

replace github.com/k8ssandra/k8ssandra/tests/integration/util => ./integration/util
