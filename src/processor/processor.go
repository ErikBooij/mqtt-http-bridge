package processor

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/blues/jsonata-go"
	"log"
	"sync"
	"text/template"
	"zigbee-coordinator/src/datastore"
	"zigbee-coordinator/src/publisher"
	"zigbee-coordinator/src/subscription"
	"zigbee-coordinator/src/utilities"
)

type Processor interface {
	Process(topic string, user string, payload string)
}

func New(store datastore.Store, publisher publisher.Publisher) Processor {
	return &processor{
		publisher: publisher,
		store:     store,

		expressionCache: make(map[string]*jsonata.Expr),
		templateCache:   make(map[string]*template.Template),
	}
}

type processor struct {
	publisher publisher.Publisher
	store     datastore.Store

	expressionCache   map[string]*jsonata.Expr
	expressionCacheMu sync.RWMutex

	templateCache   map[string]*template.Template
	templateCacheMu sync.RWMutex
}

func (p *processor) Process(topic string, user string, payload string) {
	sub, err := p.store.GetSubscriptionForTopic(topic)

	switch {
	case errors.Is(err, datastore.ErrorSubscriptionNotFound):
		return
	case err != nil:
		log.Printf("Error getting subscription for topic %s: %s\n", topic, err)
		return
	}

	extractedParameters := p.extractParametersFromMessage(sub, payload)

	parameters := map[string]any{
		"topic":   topic,
		"client":  user,
		"payload": payload,
		"custom":  extractedParameters,
	}

	if !p.filterMessage(sub, parameters) {
		log.Printf("Message for subscription %s was filtered out (%s)\n", sub.ID, payload)
		return
	}

	requestBody := p.renderTemplate(sub, parameters, payload)

	p.publisher.Publish(requestBody, sub)
}

func (p *processor) cacheExpression(expression string) *jsonata.Expr {
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
		log.Printf("Error compiling expression %s: %s\n", expression, err)
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
		log.Printf("Topic message for sub %s was not JSON: %s\n", sub.ID, err)
		return values
	}

	for key, expression := range sub.Extract {
		value, err := p.extractParameterFromData(data, expression)

		if err != nil {
			log.Printf("Error extracting value for key %s: %s\n", key, err)
			continue
		}

		values[key] = value
	}

	return values
}

func (p *processor) extractParameterFromData(data interface{}, expression string) (any, error) {
	expr := p.cacheExpression(expression)

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

	expr := p.cacheExpression(sub.Filter)

	if expr == nil {
		return true
	}

	res, err := expr.Eval(parameters)

	if err != nil {
		log.Printf("Error evaluating filter expression for subscription %s: %s\n", sub.ID, err)
		return true
	}

	if b, ok := res.(bool); ok && !b {
		// Only if the expression was successfully parsed, and evaluated to false
		return false
	}

	return true
}

func (p *processor) renderTemplate(sub subscription.Subscription, parameters map[string]any, message string) []byte {
	cacheKey := utilities.MD5Hash(sub.Template)

	p.templateCacheMu.RLock()
	tmpl, ok := p.templateCache[cacheKey]
	p.templateCacheMu.RUnlock()

	if !ok {
		tmpl, err := template.New(cacheKey).Parse(sub.Template)

		if err != nil {
			tmpl = nil
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
		log.Printf("Error rendering template for subscription %s: %s\n", sub.ID, err)
		return []byte(message)
	}

	return buf.Bytes()
}
