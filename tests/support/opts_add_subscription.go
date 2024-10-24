package testsupport

func DefaultSubscriptionOptions() AddSubscriptionOptions {
	return AddSubscriptionOptions{
		Name:       "Dummy Subscription",
		Topic:      "dummy/topic",
		HTTPMethod: "POST",
		HTTPURL:    "http://localhost:8080/",
	}
}

type AddSubscriptionOptions struct {
	Name string `json:"name"`

	Topic string `json:"topic"`

	Extract map[string]string `json:"extract"`
	Filter  string            `json:"filter"`

	HTTPMethod       string            `json:"method"`
	HTTPURL          string            `json:"url"`
	HTTPHeaders      map[string]string `json:"headers"`
	HTTPBodyTemplate string            `json:"body"`

	SubscriptionTemplateID         *string        `json:"subscriptionTemplateId"`
	SubscriptionTemplateParameters map[string]any `json:"subscriptionTemplateParameters"`
}

func (aso AddSubscriptionOptions) WithName(name string) AddSubscriptionOptions {
	clone := aso
	clone.Name = name

	return clone
}

func (aso AddSubscriptionOptions) WithTopic(topic string) AddSubscriptionOptions {
	clone := aso
	clone.Topic = topic

	return clone
}

func (aso AddSubscriptionOptions) WithExtract(extract map[string]string) AddSubscriptionOptions {
	clone := aso
	clone.Extract = extract

	return clone
}

func (aso AddSubscriptionOptions) WithFilter(filter string) AddSubscriptionOptions {
	clone := aso
	clone.Filter = filter

	return clone
}

func (aso AddSubscriptionOptions) WithHTTPMethod(method string) AddSubscriptionOptions {
	clone := aso
	clone.HTTPMethod = method

	return clone
}

func (aso AddSubscriptionOptions) WithHTTPURL(url string) AddSubscriptionOptions {
	clone := aso
	clone.HTTPURL = url

	return clone
}

func (aso AddSubscriptionOptions) WithHTTPHeader(name, value string) AddSubscriptionOptions {
	clone := aso

	if clone.HTTPHeaders == nil {
		clone.HTTPHeaders = make(map[string]string)
	}
	clone.HTTPHeaders[name] = value

	return clone
}

func (aso AddSubscriptionOptions) WithHTTPHeaders(headers map[string]string) AddSubscriptionOptions {
	clone := aso
	clone.HTTPHeaders = headers

	return clone
}

func (aso AddSubscriptionOptions) WithHTTPBodyTemplate(template string) AddSubscriptionOptions {
	clone := aso
	clone.HTTPBodyTemplate = template

	return clone
}

func (aso AddSubscriptionOptions) WithSubscriptionTemplateID(id string) AddSubscriptionOptions {
	clone := aso
	clone.SubscriptionTemplateID = &id

	return clone
}

func (aso AddSubscriptionOptions) WithSubscriptionTemplateParameters(parameters map[string]any) AddSubscriptionOptions {
	clone := aso
	clone.SubscriptionTemplateParameters = parameters

	return clone
}
