package traefik

import (
	"fmt"
	. "github.com/onsi/gomega"
	traefik "github.com/traefik/traefik/v2/pkg/provider/kubernetes/crd/traefik/v1alpha1"
)

// VerifyTraefikTCPIngressRoute asserts that the given IngressRouteTCP has the given entrypoint, match, and backend service+port.
func VerifyTraefikTCPIngressRoute(ingress traefik.IngressRouteTCP, entrypoint string, match string, backendService string, backendPort int) {
	ExpectWithOffset(1, ingress.Spec.EntryPoints).To(ContainElement(entrypoint))
	route := findTraefikTCPRouteByMatch(&ingress.Spec, match)
	ExpectWithOffset(1, route).ToNot(BeNil(), fmt.Sprintf("No Route found with match=%s for entrypoint=%s", match, entrypoint))
	service := findTraefikTCPServiceByName(route, backendService)
	ExpectWithOffset(1, service).ToNot(BeNil(), fmt.Sprintf("No BackendService found with name=%s for match=%s and entrypoint=%s", backendService, match, entrypoint))
	ExpectWithOffset(1, service.Port).To(Equal(int32(backendPort)))
}

// findTraefikTCPRouteByMatch finds a Route with the given match from among the given array of IngressRouteTCPSpecs.
func findTraefikTCPRouteByMatch(spec *traefik.IngressRouteTCPSpec, match string) *traefik.RouteTCP {
	for _, routeCandidate := range spec.Routes {
		if routeCandidate.Match == match {
			return &routeCandidate
		}
	}
	return nil
}

// findTraefikTCPServiceByName finds a ServiceTCP with the given name from among the given array of RouteTCPs.
func findTraefikTCPServiceByName(spec *traefik.RouteTCP, name string) *traefik.ServiceTCP {
	for _, serviceCandidate := range spec.Services {
		if serviceCandidate.Name == name {
			return &serviceCandidate
		}
	}
	return nil
}
