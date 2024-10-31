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

package client

import (
	"fmt"
	"testing"
	"time"

	"github.com/elastic/apm-tools/pkg/espoll"
	"github.com/elastic/go-elasticsearch/v8"
)

const (
	defaultMaxElasticsearchBackoff = 10 * time.Second
)

func defaultConfig() elasticsearch.Config {
	return elasticsearch.Config{
		MaxRetries: 5,
		RetryBackoff: func(attempt int) time.Duration {
			backoff := (500 * time.Millisecond) * (1 << (attempt - 1))
			if backoff > defaultMaxElasticsearchBackoff {
				backoff = defaultMaxElasticsearchBackoff
			}
			return backoff
		},
	}

}

func New(t *testing.T, opts ...Option) *espoll.Client {
	t.Helper()

	cfg := defaultConfig()
	for _, o := range opts {
		o(&cfg)
	}

	err := validate(cfg)
	if err != nil {
		t.Fatal(fmt.Errorf("cannot validate elasticsearch config: %w", err))
	}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		t.Fatal(fmt.Errorf("cannot create elasticsearch client: %w", err))
	}

	return espoll.WrapClient(client)
}

func validate(c elasticsearch.Config) error {
	// FIXME: add validation
	return nil
}

type Option func(*elasticsearch.Config)

func WithAddresses(a []string) Option {
	return func(c *elasticsearch.Config) {
		c.Addresses = a
	}
}

func WithUsernamePassword(u, p string) Option {
	return func(c *elasticsearch.Config) {
		c.Username = u
		c.Password = p
		c.APIKey = ""
	}
}

func WithAPIKey(a string) Option {
	return func(c *elasticsearch.Config) {
		c.Username = ""
		c.Password = ""
		c.APIKey = a
	}
}
