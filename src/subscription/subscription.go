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
	// Body is the template to use for rendering the HTTP response body
	Body string `json:"template"`
}
