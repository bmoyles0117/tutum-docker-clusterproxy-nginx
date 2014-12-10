package main

import (
	"bytes"
)

func getNginxConfig(backendRoutes []*BackendRoute) []byte {
	buffer := bytes.NewBuffer([]byte{})
	buffer.WriteString(`
upstream backend {
`)

	for _, backendRoute := range backendRoutes {
		buffer.WriteString("\tserver " + backendRoute.Addr + ":" + backendRoute.Port + ";\n")
	}

	buffer.WriteString(`
}

server {
	location / {
		proxy_pass http://backend;
	}
}`)

	return buffer.Bytes()
}
