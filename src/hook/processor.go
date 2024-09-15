package hook

import (
	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/packets"
	"mqtt-http-bridge/src/processor"
)

func ProcessorHook(processor processor.Processor) mqtt.Hook {
	return &processorHook{
		processor: processor,
	}
}

type processorHook struct {
	mqtt.HookBase

	processor processor.Processor
}

func (p *processorHook) ID() string {
	return "processor-hook"
}

func (p *processorHook) Provides(b byte) bool {
	return mqtt.OnPublished == b
}

func (p *processorHook) OnPublished(cl *mqtt.Client, pk packets.Packet) {
	p.processor.Process(pk.TopicName, string(cl.Properties.Username), string(pk.Payload))
}
