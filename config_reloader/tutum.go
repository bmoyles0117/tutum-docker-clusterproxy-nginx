package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"sort"
	"strings"
)

type ServiceEndpoint struct {
	ServiceName string
	Url         string

	roundTripper http.RoundTripper
}

func (se *ServiceEndpoint) GetService(tutumAuth string) (*Service, error) {
	req, err := http.NewRequest("GET", se.Url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", tutumAuth)

	client := http.Client{Transport: se.roundTripper}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	service := &Service{}
	if err := json.NewDecoder(res.Body).Decode(service); err != nil {
		return nil, err
	}

	return service, err
}

type BackendRoute struct {
	ContainerName string
	Addr          string
	Port          string
}

func (br *BackendRoute) GetAddr() string {
	return br.Addr + ":" + br.Port
}

type Service struct {
	LinkVariables map[string]string `json:"link_variables"`
}

func extractServiceEndpoint(environ string) (*ServiceEndpoint, error) {
	parts := strings.SplitN(environ, "=", 2)
	if len(parts) != 2 {
		return nil, errors.New("Invalid environment string passed in, should be KEY=VALUE")
	}

	index := strings.Index(parts[0], "_TUTUM_API_URL")
	if index == -1 {
		return nil, errors.New("Not a valid environment variable. Should match {SERVICE_NAME}_TUTUM_API_URL")
	}

	return &ServiceEndpoint{
		ServiceName: parts[0][:index],
		Url:         parts[1],
	}, nil
}

func extractServiceEndpoints(environs []string) []*ServiceEndpoint {
	serviceEndpoints := []*ServiceEndpoint{}

	for i := range environs {
		serviceEndpoint, err := extractServiceEndpoint(environs[i])
		if err == nil {
			serviceEndpoints = append(serviceEndpoints, serviceEndpoint)
		}
	}

	return serviceEndpoints
}

type BackendRoutesByContainerName []*BackendRoute

func (l BackendRoutesByContainerName) Len() int {
	return len(l)
}

func (l BackendRoutesByContainerName) Less(i, j int) bool {
	return l[i].ContainerName < l[j].ContainerName
}

func (l BackendRoutesByContainerName) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func extractBackendRoutes(compare string, values map[string]string) []*BackendRoute {
	backendRoutes := []*BackendRoute{}
	containerMap := map[string]int{}

	compareLen := len(compare)

	var backendRoute *BackendRoute

	for k, v := range values {
		index := strings.Index(k, compare)
		if index == -1 {
			continue
		}

		containerName := k[:index]

		if index, exists := containerMap[containerName]; exists {
			backendRoute = backendRoutes[index]
		} else {
			containerMap[containerName] = len(backendRoutes)

			backendRoute = &BackendRoute{
				ContainerName: containerName,
			}

			backendRoutes = append(backendRoutes, backendRoute)
		}

		if k[index+compareLen:] == "_ADDR" {
			backendRoute.Addr = v
		}

		if k[index+compareLen:] == "_PORT" {
			backendRoute.Port = v
		}
	}

	sort.Sort(BackendRoutesByContainerName(backendRoutes))

	return backendRoutes
}
