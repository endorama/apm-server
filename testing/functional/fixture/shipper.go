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

package fixture

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

type IntakeResponse struct {
	Accepted int
}

func New(t *testing.T, c *http.Client, f, u string) Shipper {
	return Shipper{c: c, t: t, testdataFolder: f, srvURL: u}
}

type Shipper struct {
	c *http.Client

	testdataFolder string
	srvURL         string
	t              *testing.T
}

func (fs *Shipper) NewRequest(urlPath, payload string) *http.Request {
	fs.t.Helper()

	f := openFile(fs.t, path.Join(fs.testdataFolder, fmt.Sprintf("%s.ndjson", payload)))

	u, _ := url.Parse(fs.srvURL + urlPath)
	query := u.Query()
	query.Set("verbose", "true")
	u.RawQuery = query.Encode()

	req, _ := http.NewRequest("POST", u.String(), f)
	req.Header.Add("Content-Type", "application/x-ndjson")

	return req
}

func (fs *Shipper) Send(req *http.Request) IntakeResponse {
	return _send(fs.t, fs.c, req)
}

func _send(t *testing.T, c *http.Client, req *http.Request) IntakeResponse {
	t.Helper()

	resp, err := c.Do(req)
	if err != nil {
		t.Fatal(fmt.Errorf("cannot perform request: %w", err))
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Equal(t, http.StatusAccepted, resp.StatusCode, string(respBody))

	var response IntakeResponse
	err = json.Unmarshal(respBody, &response)
	require.NoError(t, err)
	return response
}

func openFile(t *testing.T, p string) *os.File {
	f, err := os.Open(p)
	if err != nil {
		t.Fatal(err)
	} else {
		t.Cleanup(func() {
			f.Close()
		})
	}
	return f
}
