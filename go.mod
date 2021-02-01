module github.com/k8ssandra/k8ssandra

go 1.15

require (
	github.com/cloudflare/cfssl v0.0.0-20180726162950-56268a613adf // indirect
	github.com/codedellemc/goscaleio v0.0.0-20170830184815-20e2ce2cf885 // indirect
	github.com/coreos/etcd v3.3.13+incompatible // indirect
	github.com/coreos/rkt v1.30.0 // indirect
	github.com/d2g/dhcp4 v0.0.0-20170904100407-a1d1b6c41b1c // indirect
	github.com/d2g/dhcp4client v0.0.0-20170829104524-6e570ed0a266 // indirect
	github.com/datastax/cass-operator v1.4.1
	github.com/golang/groupcache v0.0.0-20191227052852-215e87163ea7 // indirect
	github.com/google/certificate-transparency-go v1.0.21 // indirect
	github.com/heketi/rest v0.0.0-20180404230133-aa6a65207413 // indirect
	github.com/heketi/utils v0.0.0-20170317161834-435bc5bdfa64 // indirect
	github.com/jteeuwen/go-bindata v0.0.0-20151023091102-a0ff2567cfb7 // indirect
	github.com/kr/fs v0.0.0-20131111012553-2788f0dbd169 // indirect
	github.com/mholt/caddy v0.0.0-20180213163048-2de495001514 // indirect
	github.com/onsi/ginkgo v1.14.2
	github.com/onsi/gomega v1.10.3
	github.com/pkg/sftp v0.0.0-20160930220758-4d0e916071f6 // indirect
	github.com/shurcooL/sanitized_anchor_name v0.0.0-20151028001915-10ef21a441db // indirect
	github.com/sigma/go-inotify v0.0.0-20181102212354-c87b6cf5033d // indirect
	github.com/stretchr/testify v1.6.1 // indirect
	github.com/vmware/photon-controller-go-sdk v0.0.0-20170310013346-4a435daef6cc // indirect
	github.com/xanzy/go-cloudstack v0.0.0-20160728180336-1e2cbf647e57 // indirect
	golang.org/x/lint v0.0.0-20191125180803-fdd1cda4f05f // indirect
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d // indirect
	k8s.io/api v0.18.6
	k8s.io/apimachinery v0.18.6
	k8s.io/client-go v12.0.0+incompatible
	k8s.io/kubernetes v1.17.4
	sigs.k8s.io/controller-runtime v0.6.4
)

replace (
	github.com/datastax/cass-operator => github.com/jsanda/cass-operator v1.0.1-0.20201111154347-91205f4d8f1e
	github.com/k8ssandra/reaper-operator => github.com/k8ssandra/reaper-operator v0.0.0-20201204054951-d4834fcfac4d
	k8s.io/api => k8s.io/api v0.18.6
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.18.6
	k8s.io/apimachinery => k8s.io/apimachinery v0.18.6
	k8s.io/apiserver => k8s.io/apiserver v0.18.6
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.18.6
	k8s.io/client-go => k8s.io/client-go v0.18.6
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.18.6
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.18.6
	k8s.io/code-generator => k8s.io/code-generator v0.18.6
	k8s.io/component-base => k8s.io/component-base v0.18.6
	k8s.io/cri-api => k8s.io/cri-api v0.18.6
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.18.6
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.18.6
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.18.6
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.18.6
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.18.6
	k8s.io/kubectl => k8s.io/kubectl v0.18.6
	k8s.io/kubelet => k8s.io/kubelet v0.18.6
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.18.6
	k8s.io/metrics => k8s.io/metrics v0.18.6
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.18.6
	sigs.k8s.io/controller-runtime => sigs.k8s.io/controller-runtime v0.6.4
)
