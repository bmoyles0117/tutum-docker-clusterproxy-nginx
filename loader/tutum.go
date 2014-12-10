package main

import (
	"errors"
	"strings"
)

type ApiEndpoint struct {
	ServiceName string
	Url         string
}

func extractApiEndpoint(environ string) (*ApiEndpoint, error) {
	parts := strings.SplitN(environ, "=", 2)
	if len(parts) != 2 {
		return nil, errors.New("Invalid environment string passed in, should be KEY=VALUE")
	}

	index := strings.Index(parts[0], "_TUTUM_API_URL")
	if index == -1 {
		return nil, errors.New("Not a valid environment variable. Should match {SERVICE_NAME}_TUTUM_API_URL")
	}

	return &ApiEndpoint{
		ServiceName: parts[0][:index],
		Url:         parts[1],
	}, nil
}

func extractApiEndpoints(environs []string) []*ApiEndpoint {
	apiEndpoints := []*ApiEndpoint{}

	for i := range environs {
		apiEndpoint, err := extractApiEndpoint(environs[i])
		if err == nil {
			apiEndpoints = append(apiEndpoints, apiEndpoint)
		}
	}

	return apiEndpoints
}
