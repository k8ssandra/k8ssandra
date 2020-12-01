module github.com/k8ssandra/k8ssandra/tests

go 1.15

require (
	github.com/gruntwork-io/terratest v0.30.23
	github.com/stretchr/testify v1.6.1
	k8s.io/api v0.19.3
	k8s.io/apimachinery v0.19.3
	sigs.k8s.io/yaml v1.2.0

)

replace github.com/k8ssandra/k8ssandra/tests/integration/util => ./integration/util
