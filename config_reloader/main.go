package main

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"
)

var (
	logger           = log.New(os.Stdout, "[nginx-proxy] ", log.Ldate|log.Ltime)
	PORT             = os.Getenv("PORT")
	TUTUM_AUTH       = os.Getenv("TUTUM_AUTH")
	POLLING_PERIOD   = 30
	LINK_ENV_PATTERN = ""
)

func writeConfig(path string, config []byte) (int, error) {
	file, err := os.OpenFile(path, os.O_WRONLY, os.ModePerm)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	if err := file.Truncate(0); err != nil {
		return 0, err
	}

	return file.Write(config)
}

func deployNginxConfig(config []byte) error {
	if _, err := writeConfig("/etc/nginx/conf.d/default.conf", config); err != nil {
		return err
	}

	return exec.Command("service", "nginx", "restart").Start()
}

func main() {
	if PORT == "" {
		PORT = "80"
	}

	if pollingPeriod, err := strconv.Atoi(os.Getenv("POLLING_PERIOD")); err == nil {
		logger.Println("Using custom polling period", pollingPeriod)
		POLLING_PERIOD = pollingPeriod
	}

	LINK_ENV_PATTERN = "_PORT_" + PORT + "_TCP"

	currentConfig := []byte{}

	for {
		for _, serviceEndpoint := range extractServiceEndpoints(os.Environ()) {
			service, err := serviceEndpoint.GetService(TUTUM_AUTH)
			if err != nil {
				logger.Fatalf("Failed to get service - %s\n", err)
			}

			backendRoutes := extractBackendRoutes(LINK_ENV_PATTERN, service.LinkVariables)
			newConfig := getNginxConfig(backendRoutes)

			if bytes.Compare(currentConfig, newConfig) != 0 {
				logger.Println("Deploying config")
				logger.Println(string(newConfig))

				if err := deployNginxConfig(newConfig); err != nil {
					logger.Fatalf("Failed to deploy nginx config - %s\n", err)
				}

				currentConfig = newConfig
			}
		}

		time.Sleep(time.Duration(POLLING_PERIOD) * time.Second)
	}
}
