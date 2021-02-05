package traefik

import (
	"fmt"
	"github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	traefik "github.com/traefik/traefik/v2/pkg/provider/kubernetes/crd/traefik/v1alpha1"
)

func VerifyTraefikHTTPIngressRoute(ingress traefik.IngressRoute, entrypoint string, match string, backendService string, backendPort int) {
	ExpectWithOffset(1, ingress.Spec.EntryPoints).To(ContainElement(entrypoint))
	route := findTraefikHTTPRouteByMatch(&ingress.Spec, match)
	ExpectWithOffset(1, route).ToNot(BeNil(), fmt.Sprintf("No Route found with match=%s for entrypoint=%s", match, entrypoint))
	service := findTraefikHTTPServiceByName(route, backendService)
	ExpectWithOffset(1, route).ToNot(BeNil(), fmt.Sprintf("No BackendService found with name=%s for match=%s and entrypoint=%s", backendService, match, entrypoint))
	ExpectWithOffset(1, service.Port).To(Equal(int32(backendPort)))
}

func findTraefikHTTPRouteByMatch(spec *traefik.IngressRouteSpec, match string) *traefik.Route {
	fmt.Fprintf(ginkgo.GinkgoWriter, "Looking for route with match=%s\n", match)
	for _, routeCandidate := range spec.Routes {
		fmt.Fprintf(ginkgo.GinkgoWriter, "Checking route: %v\n", routeCandidate)
		if routeCandidate.Match == match {
			fmt.Fprintln(ginkgo.GinkgoWriter, "Match!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
			return &routeCandidate
		}
	}
	fmt.Fprintln(ginkgo.GinkgoWriter, "NO MATCH FOUND. XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX")
	return nil
}

func findTraefikHTTPServiceByName(spec *traefik.Route, name string) *traefik.Service {
	for _, serviceCandidate := range spec.Services {
		if serviceCandidate.Name == name {
			return &serviceCandidate
		}
	}
	return nil
}
