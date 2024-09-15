package testsupport

import (
	"context"
	"fmt"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"math/rand"
	"strconv"
	"testing"
)

func StartMockServer(t *testing.T, ctx context.Context) Container {
	t.Helper()

	return startTestContainer(t, ctx, MockServer, "mockserver", "started on port")
}

func StartMockServerAndClient(t *testing.T, ctx context.Context) (Container, MockServerClient) {
	t.Helper()

	container := StartMockServer(t, ctx)

	client := NewMockServerClient(t, "localhost", container.GetMappedPort(t, 1080, "tcp"))

	return container, client
}

type containerType string

const (
	MockServer containerType = "mockserver"

	DefaultMockServerPort = 1080
)

func startTestContainer(t *testing.T, ctx context.Context, service containerType, name string, waitForMessage string) Container {
	image := ""
	ports := []string{}

	switch service {
	case MockServer:
		image = "jamesdbloom/mockserver"
		ports = []string{"1080/tcp"}
	}

	if image == "" {
		t.Fatalf("Unknown service type: %s", service)
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Name:         testContainerName(name),
			Image:        image,
			ExposedPorts: ports,
			WaitingFor:   wait.ForLog(waitForMessage),
		},

		Started: true,
	})

	if err != nil {
		t.Fatalf("Failed to start TestContainer: %s", err)
	}

	result := Container{
		TestContainer: container,
	}

	if len(ports) > 0 {
		inspect, err := container.Inspect(ctx)

		if err != nil {
			t.Fatalf("Failed to inspect TestContainer: %s", err)
		}

		result.Ports = inspect.NetworkSettings.Ports
	}

	return result
}

func testContainerName(name string) string {
	return "mqtt-http-bridge-test-container-" + name + "-" + fmt.Sprintf("%06d", rand.Int31n(1000000))
}

type Container struct {
	TestContainer testcontainers.Container
	Ports         nat.PortMap
}

func (c Container) GetMappedPort(t *testing.T, port int, protocol string) int {
	t.Helper()

	port, _ = c.GetMappedPortAndIP(t, port, protocol)

	return port
}

func (c Container) GetMappedPortAndIP(t *testing.T, port int, protocol string) (int, string) {
	t.Helper()

	targetPort := fmt.Sprintf("%d", port)

	for containerPort, hostPort := range c.Ports {
		if len(hostPort) > 0 && containerPort.Port() == targetPort && containerPort.Proto() == protocol {
			port, err := strconv.Atoi(hostPort[0].HostPort)

			if err != nil {
				t.Fatalf("Failed to convert port to int: %s", err)
			}

			return port, hostPort[0].HostIP
		}
	}

	t.Fatalf("Failed to find mapped port for %d/%s", port, protocol)

	return 0, ""
}

func (c Container) GetURL(t *testing.T) string {
	t.Helper()

	port, ip := c.GetMappedPortAndIP(t, DefaultMockServerPort, "tcp")

	return fmt.Sprintf("http://%s:%d", ip, port)
}
