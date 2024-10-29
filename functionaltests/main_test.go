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

package functionaltests

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/elastic/apm-server/functionaltests/stack"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
)

func TestEs(t *testing.T) {
	c, err := stack.Start(t)
	require.NoError(t, err)

	defer testcontainers.CleanupContainer(t, c)
	time.Sleep(180)

	_, r, err := c.Exec(context.Background(), []string{"cat", "/usr/share/elasticsearch/config/roles.yml"})
	require.NoError(t, err)
	b, err := io.ReadAll(r)
	require.NoError(t, err)
	t.Log(string(b))

	c.Stop(context.Background(), nil)
	c.Start(context.Background())

	_, r, err = c.Exec(context.Background(), []string{"cat", "/usr/share/elasticsearch/config/roles.yml"})
	require.NoError(t, err)
	b, err = io.ReadAll(r)
	require.NoError(t, err)
	t.Log(string(b))
}
