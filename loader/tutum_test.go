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
