// Copyright (c) Alex Ellis 2017, Ken Fukuyama 2017. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"

	"github.com/openfaas-incubator/faas-rancher/rancher"
	"github.com/openfaas/faas/gateway/requests"
	client "github.com/rancher/go-rancher/v2"
)

// ValidateDeployRequest validates that the service name is valid for Kubernetes
func ValidateDeployRequest(request *requests.CreateFunctionRequest) error {
	var validDNS = regexp.MustCompile(`^[a-zA-Z\-]+$`)
	matched := validDNS.MatchString(request.Service)
	if matched {
		return nil
	}

	return fmt.Errorf("(%s) must be a valid DNS entry for service name", request.Service)
}

// MakeDeployHandler creates a handler to create new functions in the cluster
func MakeDeployHandler(client rancher.BridgeClient) VarsHandler {
	return func(w http.ResponseWriter, r *http.Request, vars map[string]string) {

		defer r.Body.Close()

		body, _ := ioutil.ReadAll(r.Body)

		request := requests.CreateFunctionRequest{}
		err := json.Unmarshal(body, &request)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := ValidateDeployRequest(&request); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		serviceSpec := makeServiceSpec(request)

		_, err = client.CreateService(serviceSpec)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		log.Println("Created service - " + request.Service)
		log.Println(string(body))

		w.WriteHeader(http.StatusAccepted)

	}
}

func makeServiceSpec(request requests.CreateFunctionRequest) *client.Service {

	envVars := make(map[string]interface{})
	for k, v := range request.EnvVars {
		envVars[k] = v
	}

	if len(request.EnvProcess) > 0 {
		envVars["fprocess"] = request.EnvProcess
	}

	labels := make(map[string]interface{})
	labels[FaasFunctionLabel] = request.Service
	labels["io.rancher.container.pull_image"] = "always"

	launchConfig := &client.LaunchConfig{
		Environment: envVars,
		ImageUuid:   "docker:" + request.Image, // not sure if it's ok to just prefix with 'docker:'
		Labels:      labels,
	}

	serviceSpec := &client.Service{
		Name:          request.Service,
		Scale:         1,
		StartOnCreate: true,
		LaunchConfig:  launchConfig,
	}

	return serviceSpec
}
