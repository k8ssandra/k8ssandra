package traefik

import (
	"fmt"
	. "github.com/onsi/gomega"
	traefik "github.com/traefik/traefik/v2/pkg/provider/kubernetes/crd/traefik/v1alpha1"
)

// VerifyTraefikHTTPIngressRoute asserts that the given IngressRoute has the given entrypoint, match, and backend service+port.
func VerifyTraefikHTTPIngressRoute(ingress traefik.IngressRoute, entrypoint string, match string, backendService string, backendPort int) {
	ExpectWithOffset(1, ingress.Spec.EntryPoints).To(ContainElement(entrypoint))
	route := findTraefikHTTPRouteByMatch(&ingress.Spec, match)
	ExpectWithOffset(1, route).ToNot(BeNil(), fmt.Sprintf("No Route found with match=%s for entrypoint=%s", match, entrypoint))
	service := findTraefikHTTPServiceByName(route, backendService)
	ExpectWithOffset(1, route).ToNot(BeNil(), fmt.Sprintf("No BackendService found with name=%s for match=%s and entrypoint=%s", backendService, match, entrypoint))
	ExpectWithOffset(1, service.Port).To(Equal(int32(backendPort)))
}

// findTraefikHTTPRouteByMatch finds a Route with the given match from among the given array of IngressRouteSpecs.
func findTraefikHTTPRouteByMatch(spec *traefik.IngressRouteSpec, match string) *traefik.Route {
	for _, routeCandidate := range spec.Routes {
		if routeCandidate.Match == match {
			return &routeCandidate
		}
	}
	return nil
}

// findTraefikHTTPRouteByMatch finds a Service with the given name from among the given array of Routes.
func findTraefikHTTPServiceByName(spec *traefik.Route, name string) *traefik.Service {
	for _, serviceCandidate := range spec.Services {
		if serviceCandidate.Name == name {
			return &serviceCandidate
		}
	}
	return nil
}
