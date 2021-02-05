package kubeapi

import (
	. "fmt"
	. "github.com/onsi/gomega"
	networkingv1 "k8s.io/api/networking/v1beta1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// FindIngressRuleByHttpPath finds an IngressRule and HTTPIngressPath with a path matching the given string from among the given array of IngressRules.
func FindIngressRuleByHttpPath(rules []networkingv1.IngressRule, path string) (*networkingv1.IngressRule, *networkingv1.HTTPIngressPath) {
	for _, ruleCandidate := range rules {
		for _, pathCandidate := range ruleCandidate.HTTP.Paths {
			if pathCandidate.Path == path {
				return &ruleCandidate, &pathCandidate
			}
		}
	}
	return nil, nil
}

// VerifyNoRuleWithPath asserts that none of the given IngressRules have a given HTTP path entry.
func VerifyNoRuleWithPath(rules []networkingv1.IngressRule, path string) {
	authRule, _ := FindIngressRuleByHttpPath(rules, path)
	ExpectWithOffset(1, authRule).To(BeNil())
}

// VerifyIngressRule finds an IngresssRule from the given array with the given path and asserts that
// it exists and has the correct pathType, host, serviceName, and port.
func VerifyIngressRule(rules []networkingv1.IngressRule, path string, pathType *networkingv1.PathType, host, serviceName string, port int) {
	rule, httpPath := FindIngressRuleByHttpPath(rules, path)

	description := Sprintf("rule not found for path %v", path)
	ExpectWithOffset(1, rule).ToNot(BeNil(), description)
	ExpectWithOffset(1, httpPath).ToNot(BeNil(), description)

	description = Sprintf("Missed expectation for rule with path %v", path)
	ExpectWithOffset(1, rule.Host).To(Equal(host), description)
	if pathType == nil {
		ExpectWithOffset(1, httpPath.PathType).To(BeNil(), description)
	} else {
		ExpectWithOffset(1, httpPath.PathType).To(Equal(pathType), description)
	}
	ExpectWithOffset(1, httpPath.Backend.ServiceName).To(Equal(serviceName), description)
	ExpectWithOffset(1, httpPath.Backend.ServicePort).To(Equal(intstr.FromInt(port)), description)
}
