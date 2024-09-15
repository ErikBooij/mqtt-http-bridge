package subscription

type Subscription struct {
	// ID is the unique identifier for the subscription
	ID string `json:"id"`
	// Name is the name of the subscription
	Name string `json:"name"`

	// Topic is the MQTT topic the subscription is for
	Topic string `json:"topic"`

	// Extract is a map of variable names to JSONata expressions
	Extract map[string]string `json:"extract"`
	// Filter is a JSONata expression to filter messages, returning true if the message should be processed
	Filter string `json:"filter"`

	// Method is the HTTP method to use for the request
	Method string `json:"method"`
	// URL is the URL to send the HTTP request to
	URL string `json:"url"`
	// Headers is a map of headers to include in the request
	Headers map[string]string `json:"headers"`
	// BodyTemplate is the template to use for rendering the HTTP response body
	BodyTemplate string `json:"template"`

	// SubscriptionTemplateID is the optional ID of the template this subscription was derived from
	SubscriptionTemplateID *string `json:"subscriptionTemplateId"`
	// SubscriptionTemplateParameters is a map of parameters to use when deriving a subscription from a template
	SubscriptionTemplateParameters map[string]any `json:"subscriptionTemplateParameters"`
}

type SubscriptionTemplate struct {
	Subscription

	// RequiredParameters is a list of parameters that must be provided when deriving a subscription from a template.
	RequiredParameters []string `json:"requiredParameters"`
}
