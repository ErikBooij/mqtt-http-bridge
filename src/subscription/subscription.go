package subscription

type Subscription struct {
	// ID is the unique identifier for the subscription
	ID string
	// Topic is the MQTT topic the subscription is for
	Topic string
	// Extract is a map of variable names to JSONata expressions
	Extract map[string]string
	// Filter is a JSONata expression to filter messages, returning true if the message should be processed
	Filter string
	// Template is the template to use for rendering the HTTP response body
	Template string

	// URL is the URL to send the HTTP request to
	URL string
	// Method is the HTTP method to use for the request
	Method string
	// Headers is a map of headers to include in the request
	Headers map[string]string
}
