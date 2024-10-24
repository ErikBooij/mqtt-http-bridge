package testsupport

import (
	"context"
	"fmt"
	"mqtt-http-bridge/src/config"
	"mqtt-http-bridge/src/process"
	"net"
	"testing"
)

type TestAppOptions struct {
	exposeAppLogs bool
	mockServers   int
}

func DefaultTestAppOptions(mockServers int) TestAppOptions {
	return TestAppOptions{
		exposeAppLogs: false,
		mockServers:   mockServers,
	}
}

func (tao TestAppOptions) ExposeAppLogs(expose bool) TestAppOptions {
	tao.exposeAppLogs = expose

	return tao
}

type TestApp struct {
	APIClient
	MQTTClient
	MockServers []TestMockServer

	Ports TestPorts
}

type TestMockServer struct {
	MockServerClient

	Container Container
	URL       string
}

type TestPorts struct {
	MQTTBroker int
	HTTPServer int
}

// ConfigureAndRunApp starts the app and returns a TestApp instance. This automatically picks free ports for required
// listeners, starts the requested number of mock servers, and configures the API and MQTT clients.
func ConfigureAndRunApp(t *testing.T, ctx context.Context, opts TestAppOptions) *TestApp {
	mqttBrokerPort, httpServerPort := GetFreePort(t), GetFreePort(t)

	testApp := &TestApp{
		APIClient:   NewAPIClient(t, fmt.Sprintf("http://localhost:%d", httpServerPort)),
		MockServers: make([]TestMockServer, 0, opts.mockServers),
		MQTTClient:  nil, // To be configured after starting the app, since the client will immediately try to connect.
		Ports: TestPorts{
			MQTTBroker: mqttBrokerPort,
			HTTPServer: httpServerPort,
		},
	}

	// Start the app (SUT)
	runApp(t, ctx, RunAppOptions{
		MQTTPort:       testApp.Ports.MQTTBroker,
		HTTPPort:       testApp.Ports.HTTPServer,
		AllowAppOutput: opts.exposeAppLogs,
	})

	testApp.MQTTClient = CreateMQTTClient(t, MQTTClientOptions{Port: mqttBrokerPort})

	testMockServerC := make(chan TestMockServer)

	for range opts.mockServers {
		// Start the mock servers in parallel.
		go func() {
			container, client := StartMockServerAndClient(t, ctx)

			testMockServerC <- TestMockServer{
				Container:        container,
				MockServerClient: client,
				URL:              container.GetURL(t),
			}
		}()
	}

	for tms := range testMockServerC {
		testApp.MockServers = append(testApp.MockServers, tms)

		if len(testApp.MockServers) == opts.mockServers {
			close(testMockServerC)
		}
	}

	return testApp
}

type RunAppOptions struct {
	MQTTPort       int
	HTTPPort       int
	AllowAppOutput bool
}

func runApp(t *testing.T, ctx context.Context, options RunAppOptions) {
	appStartErr := make(chan error)

	go process.Start(ctx, &config.Config{
		AppEnv: "test",
		Broker: config.BrokerConfig{
			Address:  "127.0.0.1",
			Port:     options.MQTTPort,
			OpenAuth: true,
		},
		Server: config.ServerConfig{
			Address: "127.0.0.1",
			Port:    options.HTTPPort,
		},
		Silent: !options.AllowAppOutput,
		Storage: config.StorageConfig{
			Driver: "memory",
		},
	}, appStartErr)

	err := <-appStartErr

	if err != nil {
		t.Fatalf("Error starting process: %s", err)
	}
}

func GetFreePort(t *testing.T) int {
	t.Helper()

	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("Error obtaining a free port: %s", err)
	}
	defer listener.Close()

	addr := listener.Addr().(*net.TCPAddr)

	return addr.Port
}
