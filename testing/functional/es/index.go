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

package es

import (
	"encoding/json"
	"fmt"
	"io"
	"testing"

	"github.com/elastic/apm-tools/pkg/espoll"
)

type dataStream struct {
	Name      string
	PreferIlm bool   `json:"prefer_ilm"`
	IlmPolicy string `json:"ilm_policy"`
	Template  string `json:"template"`

	// Indices []struct {
	// 	PreferIlm bool   `json:"prefer_ilm"`
	// 	IlmPolicy string `json:"ilm_policy"`
	// }
}

type GetIndexResponse struct {
	DataStreams []dataStream `json:"data_streams"`
}

func Pointer2Bool(b bool) *bool {
	return &b
}

func GetIndex(t *testing.T, c *espoll.Client, indexName string) GetIndexResponse {
	t.Helper()

	resp, err := c.Indices.GetDataStream(
		c.Indices.GetDataStream.WithName(indexName),
		c.Indices.GetDataStream.WithIncludeDefaults(true))
	if err != nil {
		t.Fatal(fmt.Errorf("cannot create api key: %w", err))
	}
	defer resp.Body.Close()

	// fmt.Printf("%+v\n", resp)

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	var a GetIndexResponse
	err = json.Unmarshal(content, &a)
	if err != nil {
		t.Fatal(err)
	}

	return a
}

func GetIndexSettings(t *testing.T, c *espoll.Client, indexName string) {

}
