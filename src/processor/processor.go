package processor

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/blues/jsonata-go"
	"log"
	"mqtt-http-bridge/src/publisher"
	"mqtt-http-bridge/src/subscription"
	"mqtt-http-bridge/src/utilities"
	"sync"
	"text/template"
)

const InternalBroker = "internal"

type Processor interface {
	Process(message MQTTMessage)
}

type MQTTMessage struct {
	Server  string
	Topic   string
	Payload string
	// Internal Server Only
	User string
}

func New(store subscription.Service, publisher publisher.Publisher, mqttMessageChan chan<- MQTTMessage, logger *log.Logger) Processor {
	return &processor{
		logger:          logger,
		mqttMessageChan: mqttMessageChan,
		publisher:       publisher,
		service:         store,

		expressionCache: make(map[string]*jsonata.Expr),
		templateCache:   make(map[string]*template.Template),
	}
}

type processor struct {
	logger          *log.Logger
	mqttMessageChan chan<- MQTTMessage
	publisher       publisher.Publisher
	service         subscription.Service

	expressionCache   map[string]*jsonata.Expr
	expressionCacheMu sync.RWMutex

	templateCache   map[string]*template.Template
	templateCacheMu sync.RWMutex
}

func (p *processor) Process(message MQTTMessage) {
	p.mqttMessageChan <- message

	subs, err := p.service.GetSubscriptionsForTopic(message.Topic)

	switch {
	case err != nil:
		p.logger.Printf("Error getting subscriptions for topic %s: %s\n", message.Topic, err)
		return
	}

	globalParams, err := p.service.GetGlobalParameters()
	if err != nil {
		p.logger.Printf("Error getting global parameters: %s\n", err)
		return
	}

	for _, sub := range subs {
		go func() {
			parameters := map[string]any{
				"meta": map[string]any{
					"topic":   message.Topic,
					"client":  message.User,
					"payload": message.Payload,
				},
				"global":  globalParams,
				"extract": p.extractParametersFromMessage(sub, message.Payload),
			}

			sub, err = p.service.ApplyPlaceholdersOnSubscription(sub, parameters)

			if err != nil {
				p.logger.Printf("Error applying placeholders to subscription %s: %s\n", sub.ID, err)
				return
			}

			if !p.filterMessage(sub, parameters) {
				p.logger.Printf("Message for subscription %s was filtered out\n", sub.ID)
				return
			}

			requestBody := p.renderTemplate(sub, parameters, message.Payload)

			p.publisher.Publish(requestBody, sub)
		}()
	}
}

func (p *processor) cacheExpression(expression string, context string) *jsonata.Expr {
	if p.expressionCache == nil {
		p.expressionCache = make(map[string]*jsonata.Expr)
	}

	cacheKey := utilities.MD5Hash(expression)

	p.expressionCacheMu.RLock()
	expr, ok := p.expressionCache[cacheKey]
	p.expressionCacheMu.RUnlock()

	if ok {
		return expr
	}

	expr, err := jsonata.Compile(expression)

	if err != nil {
		p.logger.Printf("Error compiling expression `%s` in context %s: %s\n", expression, context, err)
	}

	p.expressionCacheMu.Lock()
	p.expressionCache[cacheKey] = expr
	p.expressionCacheMu.Unlock()

	return expr
}

func (p *processor) extractParametersFromMessage(sub subscription.Subscription, message string) map[string]any {
	values := make(map[string]any)

	if len(sub.Extract) == 0 {
		return values
	}

	var data interface{}

	if err := json.Unmarshal([]byte(message), &data); err != nil {
		p.logger.Printf("Topic message for sub %s was not JSON: %s\n", sub.ID, err)
		return values
	}

	for key, expression := range sub.Extract {
		value, err := p.extractParameterFromData(data, expression, fmt.Sprintf("parameter[%s]", key))

		if err != nil && !errors.Is(err, jsonata.ErrUndefined) {
			p.logger.Printf("Error extracting value for key %s: %s\n", key, err)
			continue
		}

		values[key] = value
	}

	return values
}

func (p *processor) extractParameterFromData(data interface{}, expression string, context string) (any, error) {
	expr := p.cacheExpression(expression, context)

	if expr == nil {
		return nil, errors.New("expression invalid")
	}

	res, err := expr.Eval(data)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (p *processor) filterMessage(sub subscription.Subscription, parameters map[string]any) bool {
	if sub.Filter == "" {
		return true
	}

	expr := p.cacheExpression(sub.Filter, "filter")

	if expr == nil {
		return true
	}

	res, err := expr.Eval(parameters)

	if err != nil {
		p.logger.Printf("Error evaluating filter expression for subscription %s: %s\n", sub.ID, err)
		return true
	}

	if b, ok := res.(bool); ok && !b {
		// Only if the expression was successfully parsed, and evaluated to false
		return false
	}

	return true
}

func (p *processor) renderTemplate(sub subscription.Subscription, parameters map[string]any, message string) []byte {
	cacheKey := utilities.MD5Hash(sub.Body)

	p.templateCacheMu.RLock()
	tmpl, ok := p.templateCache[cacheKey]
	p.templateCacheMu.RUnlock()

	if !ok {
		if sub.Body == "" {
			tmpl = nil
		} else {
			var err error
			tmpl, err = template.New(cacheKey).Parse(sub.Body)

			if err != nil {
				tmpl = nil
			}
		}

		p.templateCacheMu.Lock()
		p.templateCache[cacheKey] = tmpl
		p.templateCacheMu.Unlock()
	}

	if tmpl == nil {
		return []byte(message)
	}

	buf := new(bytes.Buffer)

	if err := tmpl.Execute(buf, parameters); err != nil {
		p.logger.Printf("Error rendering template for subscription %s: %s\n", sub.ID, err)
		return []byte(message)
	}

	return buf.Bytes()
}
