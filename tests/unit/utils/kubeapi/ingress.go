package kubeapi

import (
	. "fmt"
	. "github.com/onsi/gomega"
	networkingv1 "k8s.io/api/networking/v1beta1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

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

func VerifyNoRuleWithPath(rules []networkingv1.IngressRule, path string) {
	authRule, _ := FindIngressRuleByHttpPath(rules, path)
	ExpectWithOffset(1, authRule).To(BeNil())
}

func VerifyIngressRule(rules []networkingv1.IngressRule, releaseName, path string, pathType *networkingv1.PathType, host string, port int) {
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
	ExpectWithOffset(1, httpPath.Backend.ServicePort).To(Equal(intstr.FromInt(port)), description)
	ExpectWithOffset(1, httpPath.Backend.ServiceName).To(Equal(Sprintf("%s-%s-stargate-service", releaseName, "dc1")), description)
}
