package main

import (
	"bytes"
)

func getNginxConfig(backendRoutes []*BackendRoute) []byte {
	buffer := bytes.NewBuffer([]byte{})
	buffer.WriteString(`
upstream backend {
`)

	existingAddrs := map[string]bool{}

	for _, backendRoute := range backendRoutes {
		addr := backendRoute.GetAddr()

		// Skip duplicates
		if _, exists := existingAddrs[addr]; exists {
			continue
		}

		existingAddrs[addr] = true

		buffer.WriteString("\tserver " + addr + ";\n")
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
