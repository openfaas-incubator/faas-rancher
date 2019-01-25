// Copyright (c) Alex Ellis 2017, Ken Fukuyama 2017. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/openfaas-incubator/faas-rancher/rancher"
	"github.com/openfaas/faas/gateway/requests"
	client "github.com/rancher/go-rancher/v2"
)

// MakeUpdateHandler creates a handler to create new functions in the cluster
func MakeUpdateHandler(client rancher.BridgeClient) VarsHandler {
	return func(w http.ResponseWriter, r *http.Request, vars map[string]string) {

		defer r.Body.Close()

		body, _ := ioutil.ReadAll(r.Body)

		request := requests.CreateFunctionRequest{}
		err := json.Unmarshal(body, &request)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if len(request.Service) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		serviceSpec, findErr := client.FindServiceByName(request.Service)
		if findErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		} else if serviceSpec == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		fmt.Println(serviceSpec.State)

		if serviceSpec.State != "active" {
			fmt.Println("Service to upgrade not in active state.")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Service is not active. Maybe another update is currently in progress?"))
			return
		}

		upgradeSpec := makeUpgradeSpec(request)
		_, err = client.UpgradeService(serviceSpec, upgradeSpec)
		if err != nil {
			fmt.Println("UPGRADE ERROR", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		go func() {
			fmt.Println("Waiting for upgrade to finish")
			for pollCounter := 20; pollCounter > 0; pollCounter-- {
				pollResult, pollErr := client.FindServiceByName(request.Service)
				fmt.Println(pollResult.State)
				if pollErr != nil {
					fmt.Println("POLL ERROR")
					continue
				}
				time.Sleep(1 * time.Second)

				if pollResult.State == "upgraded" {
					fmt.Println("Finishing upgrade")
					_, err = client.FinishUpgradeService(pollResult)
					if err != nil {
						fmt.Println("FINISH ERROR", err)
						return
					}
					fmt.Println("Upgrade finished")
					return

				}
			}
			fmt.Println("Poll timeout!")
		}()

		fmt.Println("Updated service - " + request.Service)
		w.WriteHeader(http.StatusAccepted)
	}
}

func makeUpgradeSpec(request requests.CreateFunctionRequest) *client.ServiceUpgrade {

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

	spec := &client.ServiceUpgrade{
		InServiceStrategy: &client.InServiceUpgradeStrategy{
			BatchSize:              1,
			StartFirst:             true,
			LaunchConfig:           launchConfig,
			SecondaryLaunchConfigs: []client.SecondaryLaunchConfig{},
		},
	}

	return spec
}
