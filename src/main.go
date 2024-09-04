package main

import (
	"context"
	"fmt"
	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/listeners"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"zigbee-coordinator/src/config"
	"zigbee-coordinator/src/hook"
	"zigbee-coordinator/src/publisher"
	"zigbee-coordinator/src/store"
	"zigbee-coordinator/src/subscription"
)

func main() {
	ctx, cancelCtx := context.WithCancel(context.Background())

	cfg, err := config.Load()

	if err != nil {
		log.Printf("Unable to load config. Stopping.\n\n%s\n", err)
		os.Exit(1)
		return
	}

	logStartWithConfig(cfg)

	store, err := setUpStore(cfg)

	if err != nil {
		log.Printf("Unable to load store. Stopping.\n\n%s\n", err)
		os.Exit(1)
		return
	}

	store.AddSubscription(subscription.Subscription{
		Topic: "zigbee2mqtt/shortcut-button-001",
		Extract: map[string]string{
			"action":  "action",
			"battery": "battery",
		},
		Filter:   "custom.action='1_short_release'",
		Template: `{"action":"{{.custom.action}}","battery":"{{.custom.battery}}"}`,
		URL:      "https://straight-application-12.webhook.cool",
		Method:   "PATCH",
		Headers: map[string]string{
			"Authorization": "Bearer 123",
			"Content-Type":  "application/json",
		},
	})

	// Create signals channel to run server until interrupted
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		cancelCtx()
		done <- true
	}()

	// Create the new MQTT Server.
	server := mqtt.New(&mqtt.Options{
		ClientNetWriteBufferSize: 4096,
		ClientNetReadBufferSize:  4096,
		SysTopicResendInterval:   10,
		InlineClient:             false,
	})

	server.Log = slog.New(slog.NewTextHandler(io.Discard, nil))

	if err := attachHooks(ctx, server, store, cfg); err != nil {
		log.Printf("Unable to attach hooks. Stopping.\n\n%s\n", err)
		os.Exit(1)
		return
	}

	if err := attachListeners(server, cfg); err != nil {
		log.Printf("Unable to attach listeners. Stopping.\n\n%s\n", err)
		os.Exit(1)
		return
	}

	go func() {
		err := server.Serve()
		if err != nil {
			log.Printf("Unable to start server. Stopping.\n\n%s\n", err)
			os.Exit(1)
			return
		}
	}()

	<-done

	log.Println("Shutting down MQTT forwarder...")

	<-ctx.Done()

	log.Println("MQTT forwarder stopped.")
}

func attachHooks(ctx context.Context, server *mqtt.Server, store store.Store, cfg config.Config) error {
	authHook := hook.Authentication(cfg.OpenAuth)

	if !cfg.OpenAuth {
		for _, user := range cfg.UsersParsed {
			authHook.AddUser(user.Username, user.Password)
		}
	}

	if err := server.AddHook(authHook, nil); err != nil {
		return err
	}

	processorHook := hook.Processor(store, setUpPublisher(ctx, 10))

	if err := server.AddHook(processorHook, nil); err != nil {
		return err
	}

	return nil
}

func attachListeners(server *mqtt.Server, cfg config.Config) error {
	tcp := listeners.NewTCP(listeners.Config{ID: "t1", Address: fmt.Sprintf("%s:%d", cfg.BindAddress, cfg.BindPort)})

	if err := server.AddListener(tcp); err != nil {
		return err
	}

	return nil
}

func logStartWithConfig(cfg config.Config) {
	a := "without authentication"

	if !cfg.OpenAuth {
		a = fmt.Sprintf("with %d configured users", len(cfg.UsersParsed))
	}

	log.Printf("Starting MQTT forwarder on %s:%d %s\n", cfg.BindAddress, cfg.BindPort, a)
}

func setUpPublisher(ctx context.Context, parallel int) publisher.Publisher {
	return publisher.New(ctx, parallel, func() *http.Client {
		return &http.Client{}
	})
}

func setUpStore(cfg config.Config) (store.Store, error) {
	switch cfg.StorageDriver {
	case "memory":
		return store.Memory(), nil
	}

	return nil, fmt.Errorf("unknown storage driver: %s", cfg.StorageDriver)
}
