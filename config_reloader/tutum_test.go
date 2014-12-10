package main

import (
	"bytes"
	"errors"
	. "github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"net/http"
	"testing"
)

type MockRoundTripper func(*http.Request) (*http.Response, error)

func (mrt MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return mrt(req)
}

func TestBackendRouteGetAddr(t *testing.T) {
	Convey("The backend route should always produce ADDR:PORT", t, func() {
		backendRoute := &BackendRoute{Addr: "127.0.0.1", Port: "80"}
		So(backendRoute.GetAddr(), ShouldEqual, "127.0.0.1:80")

		backendRoute = &BackendRoute{Port: "80"}
		So(backendRoute.GetAddr(), ShouldEqual, ":80")

		backendRoute = &BackendRoute{Addr: "127.0.0.1"}
		So(backendRoute.GetAddr(), ShouldEqual, "127.0.0.1:")
	})
}

func TestServiceEndpointGetService(t *testing.T) {
	Convey("Returning an error should be returned by GetService", t, func() {
		mrt := MockRoundTripper(func(req *http.Request) (*http.Response, error) {
			So(req.Header.Get("Authorization"), ShouldEqual, "SAMPLE_AUTH")

			return nil, errors.New("Sample Error")
		})

		serviceEndpoint := &ServiceEndpoint{
			ServiceName:  "Testing",
			Url:          "http://www.google.com",
			roundTripper: mrt,
		}

		service, err := serviceEndpoint.GetService("SAMPLE_AUTH")
		So(err.Error(), ShouldEqual, "Get http://www.google.com: Sample Error")
		So(service, ShouldBeNil)
	})

	Convey("Returning a valid response should render a valid service", t, func() {
		mrt := MockRoundTripper(func(req *http.Request) (*http.Response, error) {
			So(req.Header.Get("Authorization"), ShouldEqual, "SAMPLE_AUTH")

			return &http.Response{
				Body: ioutil.NopCloser(bytes.NewBuffer([]byte(`{"actions":["/api/v1/action/ee921c03-0f52-438f-8b7d-7331df2a9c18/","/api/v1/action/ed4324c3-c3a3-4dad-a780-fc11ad55975d/"],"autodestroy":"OFF","autorestart":"ON_FAILURE","bindings":[{"host_path":null,"container_path":"/tmp","rewritable":true,"volume_group":"/api/v1/volumegroup/2f4f54e5-9d3b-4ac1-85ad-a2d4ff25a173/"},{"host_path":"/etc","container_path":"/etc","rewritable":true,"volume_group":null}],"container_envvars":[{"key":"DB_PASS","value":"test"}],"container_ports":[{"inner_port":80,"outer_port":80,"port_name":"http","protocol":"tcp","published":true}],"containers":["/api/v1/container/6f8ee454-9dc3-4387-80c3-57aac1be3cc6/","/api/v1/container/fdf9c116-7c08-4a60-b0ce-c54ca72c2f25/"],"cpu_shares":100,"current_num_containers":2,"deployed_datetime":"Mon, 13 Oct 2014 11:01:43 +0000","destroyed_datetime":null,"entrypoint":"","image_name":"tutum/wordpress-stackable:latest","image_tag":"/api/v1/image/tutum/wordpress-stackable/tag/latest/","link_variables":{"WORDPRESS_STACKABLE_1_ENV_DB_HOST":"**LinkMe**","WORDPRESS_STACKABLE_1_ENV_DB_NAME":"wordpress","WORDPRESS_STACKABLE_1_ENV_DB_PASS":"szVaPz925B7I","WORDPRESS_STACKABLE_1_ENV_DB_PORT":"**LinkMe**","WORDPRESS_STACKABLE_1_ENV_DB_USER":"admin","WORDPRESS_STACKABLE_1_ENV_DEBIAN_FRONTEND":"noninteractive","WORDPRESS_STACKABLE_1_ENV_HOME":"/","WORDPRESS_STACKABLE_1_ENV_PATH":"/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin","WORDPRESS_STACKABLE_1_PORT":"tcp://wordpress-stackable-1.9691c44e-admin.node.staging.tutum.io:49153","WORDPRESS_STACKABLE_1_PORT_80_TCP":"tcp://wordpress-stackable-1.9691c44e-admin.node.staging.tutum.io:49153","WORDPRESS_STACKABLE_1_PORT_80_TCP_ADDR":"wordpress-stackable-1.9691c44e-admin.node.staging.tutum.io","WORDPRESS_STACKABLE_1_PORT_80_TCP_PORT":"49153","WORDPRESS_STACKABLE_1_PORT_80_TCP_PROTO":"tcp","WORDPRESS_STACKABLE_2_ENV_DB_HOST":"**LinkMe**","WORDPRESS_STACKABLE_2_ENV_DB_NAME":"wordpress","WORDPRESS_STACKABLE_2_ENV_DB_PASS":"szVaPz925B7I","WORDPRESS_STACKABLE_2_ENV_DB_PORT":"**LinkMe**","WORDPRESS_STACKABLE_2_ENV_DB_USER":"admin","WORDPRESS_STACKABLE_2_ENV_DEBIAN_FRONTEND":"noninteractive","WORDPRESS_STACKABLE_2_ENV_HOME":"/","WORDPRESS_STACKABLE_2_ENV_PATH":"/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin","WORDPRESS_STACKABLE_2_PORT":"tcp://wordpress-stackable-2.2686cf70-admin.node.staging.tutum.io:49154","WORDPRESS_STACKABLE_2_PORT_80_TCP":"tcp://wordpress-stackable-2.2686cf70-admin.node.staging.tutum.io:49154","WORDPRESS_STACKABLE_2_PORT_80_TCP_ADDR":"wordpress-stackable-2.2686cf70-admin.node.staging.tutum.io","WORDPRESS_STACKABLE_2_PORT_80_TCP_PORT":"49154","WORDPRESS_STACKABLE_2_PORT_80_TCP_PROTO":"tcp","WORDPRESS_STACKABLE_ENV_DB_HOST":"**LinkMe**","WORDPRESS_STACKABLE_ENV_DB_NAME":"wordpress","WORDPRESS_STACKABLE_ENV_DB_PASS":"szVaPz925B7I","WORDPRESS_STACKABLE_ENV_DB_PORT":"**LinkMe**","WORDPRESS_STACKABLE_ENV_DB_USER":"admin","WORDPRESS_STACKABLE_ENV_DEBIAN_FRONTEND":"noninteractive","WORDPRESS_STACKABLE_ENV_HOME":"/","WORDPRESS_STACKABLE_ENV_PATH":"/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin","WORDPRESS_STACKABLE_PORT":"tcp://wordpress-stackable-1.9691c44e-admin.node.staging.tutum.io:49153","WORDPRESS_STACKABLE_PORT_80_TCP":"tcp://wordpress-stackable-1.9691c44e-admin.node.staging.tutum.io:49153","WORDPRESS_STACKABLE_PORT_80_TCP_ADDR":"wordpress-stackable-1.9691c44e-admin.node.staging.tutum.io","WORDPRESS_STACKABLE_PORT_80_TCP_PORT":"49153","WORDPRESS_STACKABLE_PORT_80_TCP_PROTO":"tcp","WORDPRESS_STACKABLE_TUTUM_API_URL":"https://app-test.tutum.co/api/v1/service/adeebc1b-1b81-4af0-b8f2-cefffc69d7fb/"},"linked_from_service":[],"linked_to_service":[{"from_service":"/api/v1/service/09cbcf8d-a727-40d9-b420-c8e18b7fa55b/","name":"DB","to_service":"/api/v1/service/72f175bd-390b-46e3-9463-830aca32ce3e/"}],"memory":2048,"name":"wordpress-stackable","privileged":false,"resource_uri":"/api/v1/service/09cbcf8d-a727-40d9-b420-c8e18b7fa55b/","roles":["global"],"run_command":"/run-wordpress.sh","running_num_containers":1,"sequential_deployment":false,"started_datetime":"Mon, 13 Oct 2014 11:01:43 +0000","state":"Partly running","stopped_datetime":null,"stopped_num_containers":0,"target_num_containers":2,"unique_name":"wordpress-stackable","uuid":"09cbcf8d-a727-40d9-b420-c8e18b7fa55b","tags":[{"name":"tag_one"},{"name":"tag-two"},{"name":"tagthree3"}]}`))),
			}, nil
		})

		serviceEndpoint := &ServiceEndpoint{
			ServiceName:  "Testing",
			Url:          "http://www.google.com",
			roundTripper: mrt,
		}

		service, err := serviceEndpoint.GetService("SAMPLE_AUTH")
		So(err, ShouldBeNil)
		So(service, ShouldNotBeNil)
		So(len(service.LinkVariables), ShouldEqual, 40)
	})
}

func TestExtractServiceEndpoint(t *testing.T) {
	Convey("An empty string should not produce an api URL", t, func() {
		serviceEndpoint, err := extractServiceEndpoint("")
		So(serviceEndpoint, ShouldBeNil)
		So(err, ShouldNotBeNil)
	})

	Convey("A string without an equal sign is not a valid environment variable", t, func() {
		serviceEndpoint, err := extractServiceEndpoint("https://dashboard.tutum.co/api/v1/service/SERVICE_UUID/")
		So(serviceEndpoint, ShouldBeNil)
		So(err, ShouldNotBeNil)
	})

	Convey("Extraction specifically looks for TUTUM_API_URL suffixes in the key", t, func() {
		serviceEndpoint, err := extractServiceEndpoint("TEST=https://dashboard.tutum.co/api/v1/service/SERVICE_UUID/")
		So(serviceEndpoint, ShouldBeNil)
		So(err, ShouldNotBeNil)
	})

	Convey("All things being valid, a valid API url should be returned", t, func() {
		serviceEndpoint, err := extractServiceEndpoint("WEB_TUTUM_API_URL=https://dashboard.tutum.co/api/v1/service/SERVICE_UUID/")
		So(err, ShouldBeNil)
		So(serviceEndpoint, ShouldNotBeNil)
		So(serviceEndpoint.ServiceName, ShouldEqual, "WEB")
		So(serviceEndpoint.Url, ShouldEqual, "https://dashboard.tutum.co/api/v1/service/SERVICE_UUID/")
	})
}

func TestExtractServiceEndpoints(t *testing.T) {
	Convey("Empty environment variables should not produces any API urls", t, func() {
		serviceEndpoints := extractServiceEndpoints([]string{})
		So(serviceEndpoints, ShouldBeEmpty)
	})

	Convey("Non-relevant environ variables should not be returned", t, func() {
		serviceEndpoints := extractServiceEndpoints([]string{
			"TEST=1",
			"TEST2=2",
		})
		So(serviceEndpoints, ShouldBeEmpty)
	})

	Convey("Valid environ variables should be returned", t, func() {
		serviceEndpoints := extractServiceEndpoints([]string{
			"WEB_TUTUM_API_URL=https://dashboard.tutum.co/api/v1/service/SERVICE_UUID1/",
			"TEST=1",
			"RANDOM_TUTUM_API_URL=https://dashboard.tutum.co/api/v1/service/SERVICE_UUID2/",
		})
		So(len(serviceEndpoints), ShouldEqual, 2)

		So(serviceEndpoints[0].ServiceName, ShouldEqual, "WEB")
		So(serviceEndpoints[0].Url, ShouldEqual, "https://dashboard.tutum.co/api/v1/service/SERVICE_UUID1/")

		So(serviceEndpoints[1].ServiceName, ShouldEqual, "RANDOM")
		So(serviceEndpoints[1].Url, ShouldEqual, "https://dashboard.tutum.co/api/v1/service/SERVICE_UUID2/")
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
		So(backendRoutes[0].ContainerName, ShouldEqual, "WORDPRESS_STACKABLE_1")
		So(backendRoutes[0].Addr, ShouldEqual, "wordpress-stackable-1.9691c44e-admin.node.staging.tutum.io")
		So(backendRoutes[0].Port, ShouldEqual, "")
	})

	Convey("A map with just a port and not an addr should define a backend route accordingly", t, func() {
		backendRoutes := extractBackendRoutes("_PORT_80_TCP", map[string]string{
			"WORDPRESS_STACKABLE_1_PORT_80_TCP_PORT": "80",
		})
		So(len(backendRoutes), ShouldEqual, 1)
		So(backendRoutes[0].ContainerName, ShouldEqual, "WORDPRESS_STACKABLE_1")
		So(backendRoutes[0].Addr, ShouldEqual, "")
		So(backendRoutes[0].Port, ShouldEqual, "80")
	})

	Convey("A map with both port and addr should define a backend route accordingly", t, func() {
		backendRoutes := extractBackendRoutes("_PORT_80_TCP", map[string]string{
			"WORDPRESS_STACKABLE_1_PORT_80_TCP_ADDR": "wordpress-stackable-1.9691c44e-admin.node.staging.tutum.io",
			"TEST1": "1",
			"TEST2": "2",
			"WORDPRESS_STACKABLE_1_PORT_80_TCP_PORT": "80",
		})
		So(len(backendRoutes), ShouldEqual, 1)
		So(backendRoutes[0].ContainerName, ShouldEqual, "WORDPRESS_STACKABLE_1")
		So(backendRoutes[0].Addr, ShouldEqual, "wordpress-stackable-1.9691c44e-admin.node.staging.tutum.io")
		So(backendRoutes[0].Port, ShouldEqual, "80")
	})
}
