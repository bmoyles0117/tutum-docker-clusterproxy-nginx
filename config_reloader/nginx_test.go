package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestGetNginxConfig(t *testing.T) {
	Convey("A config writer with nil backend routes should produce a basic config", t, func() {
		config := getNginxConfig(nil)
		So(string(config), ShouldEqual, `
upstream backend {

}

server {
	location / {
		proxy_pass http://backend;
	}
}`)
	})

	Convey("A config writer with a single backend route should produce a config with one entry", t, func() {
		config := getNginxConfig([]*BackendRoute{
			&BackendRoute{
				ContainerName: "WEB_1",
				Addr:          "127.0.0.1",
				Port:          "80",
			},
		})

		So(string(config), ShouldEqual, `
upstream backend {
	server 127.0.0.1:80;

}

server {
	location / {
		proxy_pass http://backend;
	}
}`)
	})

	Convey("A config writer with multiple backend routes should produce a config with multiple entries", t, func() {
		config := getNginxConfig([]*BackendRoute{
			&BackendRoute{
				ContainerName: "WEB_1",
				Addr:          "127.0.0.0",
				Port:          "80",
			},
			&BackendRoute{
				ContainerName: "WEB_2",
				Addr:          "127.0.0.1",
				Port:          "81",
			},
		})

		So(string(config), ShouldEqual, `
upstream backend {
	server 127.0.0.0:80;
	server 127.0.0.1:81;

}

server {
	location / {
		proxy_pass http://backend;
	}
}`)
	})

	Convey("A config writer with duplicate backend routes by address should eliminate duplicates", t, func() {
		config := getNginxConfig([]*BackendRoute{
			&BackendRoute{
				ContainerName: "WEB_1",
				Addr:          "127.0.0.0",
				Port:          "80",
			},
			&BackendRoute{
				ContainerName: "WEB_1",
				Addr:          "127.0.0.0",
				Port:          "80",
			},
			&BackendRoute{
				ContainerName: "WEB_2",
				Addr:          "127.0.0.1",
				Port:          "81",
			},
		})

		So(string(config), ShouldEqual, `
upstream backend {
	server 127.0.0.0:80;
	server 127.0.0.1:81;

}

server {
	location / {
		proxy_pass http://backend;
	}
}`)
	})
}
