package backend_test

import (
	"context"
	"fmt"
	testsupport "mqtt-http-bridge/tests/support"
	"testing"
)

/**
 * This test is a simulation of an external consumer of the application that:
 * 1. Starts the application
 * 2. Starts two mock servers to receive webhooks
 * 3. Sets up subscriptions through the HTTP JSON API
 * 4. Publishes a number of messages to the MQTT broker
 * 5. Checks that the recorded requests on the mock servers match the expected values
 */
func TestWebhooksArriveAfterSettingUpPlainSubscriptions(t *testing.T) {
	// Run the application
	testApp := testsupport.ConfigureAndRunApp(t, context.Background(), testsupport.DefaultTestAppOptions(2))

	// Set a global parameter to be used by templates
	testApp.SetGlobalParameter("authTokenService1", "test-1234")
	testApp.SetGlobalParameter("authTokenService2", "test-5678")

	// Configure a basic subscription
	subscriptionBase := testsupport.DefaultSubscriptionOptions()

	testApp.AddSubscription(
		subscriptionBase.
			WithName("Test subscription 1").
			WithTopic("test/topic").
			WithExtract(map[string]string{"actionParam": "action", "valueParam": "value"}).
			WithFilter("extract.valueParam = true").
			WithHTTPMethod("POST").
			WithHTTPURL(fmt.Sprintf("%s/webhooks", testApp.MockServers[0].URL)).
			WithHTTPHeader("Authorization", "Bearer {{.global.authTokenService1}}").
			WithHTTPHeader("Content-Type", "application/json").
			WithHTTPBody(`{"action":"{{.extract.actionParam}}","translated":true}`),
	)
	testApp.AddSubscription(
		subscriptionBase.
			WithName("Test subscription 2").
			WithTopic("some/other/topic").
			WithHTTPMethod("PUT").
			WithHTTPURL(fmt.Sprintf("%s/callback/5", testApp.MockServers[1].URL)).
			WithHTTPHeader("Authorization", "Bearer {{.global.authTokenService2}}"),
	)

	// Publish test messages to the MQTT broker
	testApp.Publish("test/topic", `{"action":"button-1-pressed"}`)
	testApp.Publish("test/topic", `{"action":"button-1-pressed","value":null}`)
	testApp.Publish("test/topic", `{"action":"button-1-pressed","value":true}`)
	testApp.Publish("test/topic", `{"action":"button-1-pressed","value":true}`)
	testApp.Publish("test/topic", `{"action":"button-1-pressed","value":false}`)
	testApp.Publish("some/other/topic", `{"action":"button-pressed"}`)

	// Check recorded requests on mock servers
	testApp.MockServers[0].AssertRequestRecorded(
		"POST",
		"/webhooks",
		testsupport.Header("Content-Type", "application/json"),
		testsupport.Header("Authorization", "Bearer test-1234"),
		testsupport.Header("Subscription-Name", "Test subscription 1"),
		testsupport.Body(`{"action":"button-1-pressed","translated":true}`),
		testsupport.Exactly(2),
	)

	testApp.MockServers[1].AssertRequestRecorded(
		"PUT",
		"/callback/5",
		testsupport.Header("Authorization", "Bearer test-5678"),
		testsupport.Header("Subscription-Name", "Test subscription 2"),
		testsupport.Body(`{"action":"button-pressed"}`),
		testsupport.Exactly(1),
	)
}
