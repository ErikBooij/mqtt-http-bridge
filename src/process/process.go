package process

import (
	"context"
	"fmt"
	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/listeners"
	"io"
	"log"
	"log/slog"
	"mqtt-http-bridge/src/config"
	"mqtt-http-bridge/src/datastore"
	"mqtt-http-bridge/src/dev"
	"mqtt-http-bridge/src/hook"
	"mqtt-http-bridge/src/processor"
	"mqtt-http-bridge/src/publisher"
	"mqtt-http-bridge/src/server"
	"mqtt-http-bridge/src/subscription"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Start(ctx context.Context, cfg *config.Config, appStartErr chan error) {
	defer close(appStartErr)

	if ctx == nil {
		ctx = context.Background()
	}

	logger := log.New(os.Stdout, "", log.LstdFlags)

	if cfg.Silent {
		logger.SetOutput(io.Discard)
	}

	logStartWithConfig(cfg, logger)

	store, err := setUpStore(cfg)

	if err != nil {
		appStartErr <- fmt.Errorf("unable to load store: %w", err)
		return
	}

	service := subscription.NewService(store)

	if cfg.IsDevelopment() && cfg.PrepareData {
		logger.Println("Development mode detected, preparing data store.")

		if err := dev.PopulateDataStore(service, logger); err != nil {
			appStartErr <- fmt.Errorf("unable to populate data store: %w", err)
			return
		}
	}

	mqttMessageChan := make(chan processor.MQTTMessage, 100)

	proc := processor.New(service, setUpPublisher(ctx, 10, logger), mqttMessageChan, logger)

	// Create signals channel to run broker until interrupted
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		done <- true
	}()

	go func() {
		<-ctx.Done()
		done <- true
	}()

	// Create the new MQTT Server.
	broker := mqtt.New(&mqtt.Options{
		ClientNetWriteBufferSize: 4096,
		ClientNetReadBufferSize:  4096,
		SysTopicResendInterval:   10,
		InlineClient:             false,
	})

	broker.Log = slog.New(slog.NewTextHandler(io.Discard, nil))

	if err := attachHooks(broker, proc, cfg); err != nil {
		appStartErr <- fmt.Errorf("unable to attach hooks: %w", err)
		return
	}

	if err := attachListeners(broker, cfg); err != nil {
		appStartErr <- fmt.Errorf("unable to attach listeners: %w", err)
		return
	}

	httpServer := setUpServer(service, mqttMessageChan, cfg)

	go func() {
		err := broker.Serve()
		if err != nil {
			appStartErr <- fmt.Errorf("unable to start broker: %w", err)
			return
		}
	}()

	go func() {
		err := httpServer.Start(fmt.Sprintf("%s:%d", cfg.Server.Address, cfg.Server.Port))
		if err != nil {
			appStartErr <- fmt.Errorf("unable to start httpServer: %w", err)
			return
		}
	}()

	appStartErr <- nil

	<-done

	logger.Println("Shutting down MQTT forwarder...")
}

func attachHooks(server *mqtt.Server, processor processor.Processor, cfg *config.Config) error {
	authHook := hook.Authentication(cfg.Broker.OpenAuth)

	if !cfg.Broker.OpenAuth {
		for _, user := range cfg.Broker.Users {
			authHook.AddUser(user.Username, user.Password)
		}
	}

	if err := server.AddHook(authHook, nil); err != nil {
		return err
	}

	processorHook := hook.ProcessorHook(processor)

	if err := server.AddHook(processorHook, nil); err != nil {
		return err
	}

	return nil
}

func attachListeners(server *mqtt.Server, cfg *config.Config) error {
	tcp := listeners.NewTCP(listeners.Config{ID: "t1", Address: fmt.Sprintf("%s:%d", cfg.Broker.Address, cfg.Broker.Port)})

	if err := server.AddListener(tcp); err != nil {
		return err
	}

	return nil
}

func logStartWithConfig(cfg *config.Config, logger *log.Logger) {
	a := "without authentication"

	if !cfg.Broker.OpenAuth {
		a = fmt.Sprintf("with %d configured users", len(cfg.Broker.Users))
	}

	logger.Printf("Starting MQTT broker on %s:%d %s\n", cfg.Broker.Address, cfg.Broker.Port, a)
	logger.Printf("Starting HTTP server on %s:%d\n", cfg.Server.Address, cfg.Server.Port)
	logger.Printf("Using %s storage driver\n", cfg.Storage.Driver)
}

func setUpPublisher(ctx context.Context, parallel int, logger *log.Logger) publisher.Publisher {
	return publisher.New(ctx, parallel, func() *http.Client {
		return &http.Client{}
	}, logger)
}

func setUpServer(service subscription.Service, mqttMessageChan <-chan processor.MQTTMessage, cfg *config.Config) server.HTTPServer {
	return server.New(service, mqttMessageChan, cfg)
}

func setUpStore(cfg *config.Config) (datastore.Store, error) {
	switch cfg.Storage.Driver {
	case "memory":
		return datastore.Memory()
	case "file":
		storageConfig, err := cfg.StorageConfigFile()

		if err != nil {
			return nil, err
		}

		return datastore.File(storageConfig.File, 5*time.Second)
	}

	return nil, fmt.Errorf("unknown storage driver: %s", cfg.Storage.Driver)
}
