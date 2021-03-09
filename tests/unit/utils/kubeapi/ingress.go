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
	description := Sprintf("Found a rule with path %v when none was expected", path)
	Expect(authRule).To(BeNil(), description)
}

// VerifyIngressRule finds an IngressRule from the given array with the given path and asserts that
// it exists and has the correct pathType, host, serviceName, and port.
func VerifyIngressRule(rules []networkingv1.IngressRule, path string, host *string, serviceName string, port int) {
	rule, httpPath := FindIngressRuleByHttpPath(rules, path)

	description := Sprintf("rule not found for path %v", path)
	Expect(rule).ToNot(BeNil(), description)
	Expect(httpPath).ToNot(BeNil(), description)

	description = Sprintf("Missed expectation for rule with path %v", path)

	if host == nil {
		Expect(rule.Host).To(Equal(""), description)
	} else {
		Expect(rule.Host).To(Equal(*host), description)
	}
	Expect(httpPath.PathType).To(BeNil(), description) // PathType requires k8s 1.18+; as of now we are supporting 1.16+
	Expect(httpPath.Backend.ServiceName).To(Equal(serviceName), description)
	Expect(httpPath.Backend.ServicePort).To(Equal(intstr.FromInt(port)), description)
}
