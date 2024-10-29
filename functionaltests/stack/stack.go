package stack

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func Start(t testing.TB) (testcontainers.Container, error) {
	t.Helper()

	rolesfile, err := filepath.Abs("./testing/docker/elasticsearch/roles.yml")
	require.NoError(t, err)

	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "docker.elastic.co/elasticsearch/elasticsearch:8.13.0",
		ExposedPorts: []string{"9200/tcp"},
		Env: map[string]string{
			"ES_JAVA_OPTS":   "-Xms1g -Xmx1g",
			"network.host":   "0.0.0.0",
			"transport.host": "127.0.0.1",
			"http.host":      "0.0.0.0",
			"cluster.routing.allocation.disk.threshold_enabled": "false",
			"discovery.type":                                   "single-node",
			"xpack.security.authc.anonymous.roles":             "remote_monitoring_collector",
			"xpack.security.authc.realms.file.file1.order":     "0",
			"xpack.security.authc.realms.native.native1.order": "1",
			"xpack.security.enabled":                           "true",
			"xpack.license.self_generated.type":                "trial",
			"xpack.security.authc.token.enabled":               "true",
			"xpack.security.authc.api_key.enabled":             "true",
			"logger.org.elasticsearch":                         "error",
			"action.destructive_requires_name":                 "false",
		},
		Files: []testcontainers.ContainerFile{
			// volumes:
			//   - "./testing/docker/elasticsearch/roles.yml:/usr/share/elasticsearch/config/roles.yml"
			//   - "./testing/docker/elasticsearch/users:/usr/share/elasticsearch/config/users"
			//   - "./testing/docker/elasticsearch/users_roles:/usr/share/elasticsearch/config/users_roles"
			//   - "./testing/docker/elasticsearch/ingest-geoip:/usr/share/elasticsearch/config/ingest-geoip"
			{
				HostFilePath:      rolesfile,
				ContainerFilePath: "/usr/share/elasticsearch/config/roles.yml",
				FileMode:          0666,
			},
		},
		// healthcheck:
		//   test: ["CMD-SHELL", "curl -s http://localhost:9200/_cluster/health?wait_for_status=yellow&timeout=500ms"]
		//   retries: 300
		//   interval: 1s
		WaitingFor: wait.ForExec([]string{
			"curl", "-s", "http://localhost:9200/_cluster/health?wait_for_status=yellow&timeout=500ms"}),
	}

	esC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	return esC, err
}
