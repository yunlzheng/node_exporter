// Copyright 2015 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package collector includes all individual collectors to gather and export system metrics.
package collector

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"io/ioutil"
	"net/http"
	"os"
)

// Namespace defines the common namespace to be used by all metrics.
const Namespace = "node"

// Factories contains the list of all available collectors.
var Factories = make(map[string]func() (Collector, error))

var (
	metadataServer = getEnv("RANCHER_METADATA", "http://169.254.169.250")
)

// getEnv - Allows us to supply a fallback option if nothing specified
func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func warnDeprecated(collector string) {
	log.Warnf("The %s collector is deprecated and will be removed in the future!", collector)
}

// Collector is the interface a collector has to implement.
type Collector interface {
	// Get new metrics and expose them via prometheus registry.
	Update(ch chan<- prometheus.Metric) error
}

type typedDesc struct {
	desc      *prometheus.Desc
	valueType prometheus.ValueType
}

var agentIP string
var environmentUUID string

func init() {
	environmentUUID, _ = getMetadata("environment_uuid")
	agentIP, _ = getMetadata("agent_ip")
	fmt.Printf("init() current agent_ip: %s environment_uuid: %s", agentIP, environmentUUID)
}

func getMetadata(key string) (string, error) {
	resp, err := http.Get(metadataServer + "/latest/self/host/" + key)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	data, _ := ioutil.ReadAll(resp.Body)
	return string(data), nil
}

func (d *typedDesc) mustNewConstMetric(value float64, labels ...string) prometheus.Metric {
	return prometheus.MustNewConstMetric(d.desc, d.valueType, value, labels...)
}
