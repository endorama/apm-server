// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package main_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/elastic/apm-server/testing/functional/es"
	esclient "github.com/elastic/apm-server/testing/functional/es/client"
	"github.com/elastic/apm-server/testing/functional/fixture"
	"github.com/elastic/apm-tools/pkg/apmclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpgrade(t *testing.T) {
	// create infra, manually with terraform apply
	// extract credentials with terraform output --json
	infra := readInfra(t, "./infra.json")

	// setup ES client
	client := esclient.New(t,
		esclient.WithAddresses([]string{infra.ElasticsearchURL.Value}),
		esclient.WithUsernamePassword(infra.ElasticsearchUsername.Value, infra.ElasticsearchPassword.Value),
	)
	// apm is already setup
	// ingest events with apmsoak/apmtelemetrygen

	aclient, err := apmclient.New(apmclient.Config{
		ElasticsearchURL: infra.ElasticsearchURL.Value,
		Username:         infra.ElasticsearchUsername.Value,
		Password:         infra.ElasticsearchPassword.Value,
	})
	require.NoError(t, err)
	ak, err := aclient.CreateAgentAPIKey(context.Background(), 0)
	require.NoError(t, err)

	// TODO: would be nice supporting writing to different data stream namespaces, to allow
	// parallel tests to ingest data in the same cluster but different datastreams.
	s := fixture.New(t, http.DefaultClient, "../../../testdata/intake-v2", infra.APMServerURL.Value)
	req := s.NewRequest("/intake/v2/events", "errors")
	fmt.Printf("%+v\n", req)
	req.Header.Add("Authorization", "ApiKey "+ak)
	fmt.Printf("%+v\n", req)
	response := s.Send(req)
	fmt.Printf("%+v\n", response)
	require.Equal(t, 5, response.Accepted)
	// check indices

	r := es.GetIndex(t, client, "logs-apm.error-default")
	assert.True(t, r.DataStreams[0].PreferIlm)
	fmt.Println(r.DataStreams[0].IlmPolicy)

	// perform upgrade
	// ingest data
	// check indices
}

type setting struct {
	// Sensitive bool
	// Type      string
	Value string
}

type settings struct {
	APMSecretToken        setting `json:"apm_secret_token"`
	APMServerURL          setting `json:"apm_server_url"`
	ElasticsearchPassword setting `json:"elasticsearch_password"`
	ElasticsearchURL      setting `json:"elasticsearch_url"`
	ElasticsearchUsername setting `json:"elasticsearch_username"`
	KibanaURL             setting `json:"kibana_url"`
	StackVersion          setting `json:"stack_version"`
	DeploymentID          setting `json:"deployment_id"`
}

func readInfra(t *testing.T, path string) settings {
	t.Helper()

	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	content, err := io.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}
	var s settings
	err = json.Unmarshal(content, &s)
	if err != nil {
		t.Fatal(err)
	}
	return s
}
