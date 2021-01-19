package util

import (
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
)

type KindClusterDetail struct {
	Name       string
	Image      string
	ConfigFile string
}

type ReleaseDetail struct {
	Name       string
	Namespace  string
	Revision   string
	UpdatedZDT string
	Status     string
	Chart      string
	AppVersion string
}

type PodDetail struct {
	Name   string
	Ready  string
	Status string
}

type PodLabel struct {
	Key   string
	Value string
}

type K8ssandraOptions struct {
	Operator OperatorOptions
	Network  NetworkOptions
}

type OptionsContext struct {
	HelmOptions *helm.Options
	KubeOptions *k8s.KubectlOptions
	Namespace   string
	ReleaseName string
	Args        []string
}

type OperatorOptions struct {
	Ctx OptionsContext
}

type NetworkOptions struct {
	Ctx OptionsContext
}

func NewK8ssandraOptions() *K8ssandraOptions {
	return &K8ssandraOptions{}
}
func (oc *OptionsContext) SetArgs(args []string) {
	oc.Args = args
}

func (ko *K8ssandraOptions) SetOperator(options *helm.Options, releaseName string) {
	ko.Operator = createOperatorOptions(options, releaseName)
}

func (ko *K8ssandraOptions) SetNetwork(options *helm.Options, releaseName string) {
	ko.Network = createNetworkOptions(options, releaseName)
}

func (ko *K8ssandraOptions) GetOperatorCtx() OptionsContext {
	return ko.Operator.Ctx
}

func (ko *K8ssandraOptions) GetNetworkCtx() OptionsContext {
	return ko.Network.Ctx
}

func (ko *K8ssandraOptions) GetK8ssandraOptions() K8ssandraOptions {

	return K8ssandraOptions{
		Operator: ko.Operator,
		Network:  ko.Network,
	}
}

func createOperatorOptions(helmOptions *helm.Options, releaseName string) OperatorOptions {

	oc := OptionsContext{
		HelmOptions: helmOptions,
		KubeOptions: helmOptions.KubectlOptions,
		Namespace:   helmOptions.KubectlOptions.Namespace,
		ReleaseName: releaseName}
	return OperatorOptions{Ctx: oc}
}

func createNetworkOptions(helmOptions *helm.Options, releaseName string) NetworkOptions {

	oc := OptionsContext{
		HelmOptions: helmOptions,
		KubeOptions: helmOptions.KubectlOptions,
		Namespace:   helmOptions.KubectlOptions.Namespace,
		ReleaseName: releaseName}

	return NetworkOptions{Ctx: oc}
}
