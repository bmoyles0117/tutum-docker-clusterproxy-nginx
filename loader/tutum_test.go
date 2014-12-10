package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestExtractApiEndpoint(t *testing.T) {
	Convey("An empty string should not produce an api URL", t, func() {
		apiEndpoint, err := extractApiEndpoint("")
		So(apiEndpoint, ShouldBeNil)
		So(err, ShouldNotBeNil)
	})

	Convey("A string without an equal sign is not a valid environment variable", t, func() {
		apiEndpoint, err := extractApiEndpoint("https://dashboard.tutum.co/api/v1/service/SERVICE_UUID/")
		So(apiEndpoint, ShouldBeNil)
		So(err, ShouldNotBeNil)
	})

	Convey("Extraction specifically looks for TUTUM_API_URL suffixes in the key", t, func() {
		apiEndpoint, err := extractApiEndpoint("TEST=https://dashboard.tutum.co/api/v1/service/SERVICE_UUID/")
		So(apiEndpoint, ShouldBeNil)
		So(err, ShouldNotBeNil)
	})

	Convey("All things being valid, a valid API url should be returned", t, func() {
		apiEndpoint, err := extractApiEndpoint("WEB_TUTUM_API_URL=https://dashboard.tutum.co/api/v1/service/SERVICE_UUID/")
		So(err, ShouldBeNil)
		So(apiEndpoint, ShouldNotBeNil)
		So(apiEndpoint.ServiceName, ShouldEqual, "WEB")
		So(apiEndpoint.Url, ShouldEqual, "https://dashboard.tutum.co/api/v1/service/SERVICE_UUID/")
	})
}

func TestExtractApiEndpoints(t *testing.T) {
	Convey("Empty environment variables should not produces any API urls", t, func() {
		apiEndpoints := extractApiEndpoints([]string{})
		So(apiEndpoints, ShouldBeEmpty)
	})

	Convey("Non-relevant environ variables should not be returned", t, func() {
		apiEndpoints := extractApiEndpoints([]string{
			"TEST=1",
			"TEST2=2",
		})
		So(apiEndpoints, ShouldBeEmpty)
	})

	Convey("Valid environ variables should be returned", t, func() {
		apiEndpoints := extractApiEndpoints([]string{
			"WEB_TUTUM_API_URL=https://dashboard.tutum.co/api/v1/service/SERVICE_UUID1/",
			"TEST=1",
			"RANDOM_TUTUM_API_URL=https://dashboard.tutum.co/api/v1/service/SERVICE_UUID2/",
		})
		So(len(apiEndpoints), ShouldEqual, 2)

		So(apiEndpoints[0].ServiceName, ShouldEqual, "WEB")
		So(apiEndpoints[0].Url, ShouldEqual, "https://dashboard.tutum.co/api/v1/service/SERVICE_UUID1/")

		So(apiEndpoints[1].ServiceName, ShouldEqual, "RANDOM")
		So(apiEndpoints[1].Url, ShouldEqual, "https://dashboard.tutum.co/api/v1/service/SERVICE_UUID2/")
	})
}

func TestExtractBackendRoutes(t *testing.T) {
	Convey("A map without any values should not produce backend routes", t, func() {
		backendRoutes := extractBackendRoutes("_PORT_80_TCP", map[string]string{})
		So(backendRoutes, ShouldBeEmpty)
	})

	Convey("A map with just an addr and not a port should define a backend route accordingly", t, func() {
		backendRoutes := extractBackendRoutes("_PORT_80_TCP", map[string]string{
			"WORDPRESS_STACKABLE_1_PORT_80_TCP_ADDR": "wordpress-stackable-1.9691c44e-admin.node.staging.tutum.io",
		})
		So(len(backendRoutes), ShouldEqual, 1)
		So(backendRoutes["WORDPRESS_STACKABLE_1"].Addr, ShouldEqual, "wordpress-stackable-1.9691c44e-admin.node.staging.tutum.io")
		So(backendRoutes["WORDPRESS_STACKABLE_1"].Port, ShouldEqual, "")
	})

	Convey("A map with just a port and not an addr should define a backend route accordingly", t, func() {
		backendRoutes := extractBackendRoutes("_PORT_80_TCP", map[string]string{
			"WORDPRESS_STACKABLE_1_PORT_80_TCP_PORT": "80",
		})
		So(len(backendRoutes), ShouldEqual, 1)
		So(backendRoutes["WORDPRESS_STACKABLE_1"].Addr, ShouldEqual, "")
		So(backendRoutes["WORDPRESS_STACKABLE_1"].Port, ShouldEqual, "80")
	})

	Convey("A map with both port and addr should define a backend route accordingly", t, func() {
		backendRoutes := extractBackendRoutes("_PORT_80_TCP", map[string]string{
			"WORDPRESS_STACKABLE_1_PORT_80_TCP_ADDR": "wordpress-stackable-1.9691c44e-admin.node.staging.tutum.io",
			"TEST1": "1",
			"TEST2": "2",
			"WORDPRESS_STACKABLE_1_PORT_80_TCP_PORT": "80",
		})
		So(len(backendRoutes), ShouldEqual, 1)
		So(backendRoutes["WORDPRESS_STACKABLE_1"].Addr, ShouldEqual, "wordpress-stackable-1.9691c44e-admin.node.staging.tutum.io")
		So(backendRoutes["WORDPRESS_STACKABLE_1"].Port, ShouldEqual, "80")
	})
}
