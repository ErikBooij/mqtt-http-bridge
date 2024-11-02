package testsupport

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
)

type APIClient interface {
	AddSubscription(opts AddSubscriptionOptions) SubscriptionResponse
	GetSubscription(id string) SubscriptionResponse
	ListSubscriptions() SubscriptionsResponse
	UpdateSubscription(id string, opts UpdateSubscriptionOptions) SubscriptionResponse
	DeleteSubscription(id string)

	SetGlobalParameter(parameter string, value any)
	GetGlobalParameter(parameter string) any
	ListGlobalParameters() GlobalParameters
}

func NewAPIClient(t *testing.T, host string) APIClient {
	return &apiClient{
		client: newTestHTTPClient(t, strings.TrimRight(host, "/")+"/api/v1"),

		t: t,
	}
}

type apiClient struct {
	client *testHTTPClient

	t *testing.T
}

func (a *apiClient) AddSubscription(opts AddSubscriptionOptions) (resp SubscriptionResponse) {
	a.client.doAssign(&resp, http.MethodPost, "/subscriptions", bodyJson(opts), successStatuses(http.StatusCreated))

	return resp
}

func (a *apiClient) GetSubscription(id string) (resp SubscriptionResponse) {
	a.client.doAssign(&resp, http.MethodGet, fmt.Sprintf("/subscriptions/%s", id), nil)

	return resp
}

func (a *apiClient) ListSubscriptions() (resp SubscriptionsResponse) {
	a.client.doAssign(&resp, http.MethodGet, "/subscriptions", nil)

	return resp
}

func (a *apiClient) UpdateSubscription(id string, opts UpdateSubscriptionOptions) (resp SubscriptionResponse) {
	a.client.doAssign(&resp, http.MethodPut, fmt.Sprintf("/subscriptions/%s", id), bodyJson(opts))

	return resp
}

func (a *apiClient) DeleteSubscription(id string) {
	a.client.do(http.MethodDelete, fmt.Sprintf("/subscriptions/%s", id), nil)
}

func (a *apiClient) SetGlobalParameter(parameter string, value any) {
	a.client.do(http.MethodPost, "/global-parameters", bodyJson(SetGlobalParameterOptions{
		Key:   parameter,
		Value: value,
	}))
}

func (a *apiClient) GetGlobalParameter(parameter string) any {
	var globalParameters GlobalParameters

	a.client.doAssign(&globalParameters, http.MethodGet, "/global-parameters")

	if value, ok := globalParameters.Parameters[parameter]; ok {
		return value
	}

	return nil
}

func (a *apiClient) ListGlobalParameters() (resp GlobalParameters) {
	a.client.doAssign(&resp, http.MethodGet, "/global-parameters")

	return resp
}

type SubscriptionResponse struct {
	Subscription Subscription `json:"subscription"`
}

type SubscriptionsResponse struct {
	Subscriptions []Subscription `json:"subscriptions"`
}

type Subscription struct {
	ID   string `json:"id"`
	Name string `json:"name"`

	Topic string `json:"topic"`

	Extract map[string]string `json:"extract"`
	Filter  string            `json:"filter"`

	HTTPMethod  string            `json:"method"`
	HTTPURL     string            `json:"path"`
	HTTPHeaders map[string]string `json:"httpHeaders"`
	HTTPBody    string            `json:"httpBodyTemplate"`
}

type UpdateSubscriptionOptions struct {
	Name string `json:"name"`

	Topic string `json:"topic"`

	Extract map[string]string `json:"extract"`
	Filter  string            `json:"filter"`

	HTTPMethod  string            `json:"method"`
	HTTPURL     string            `json:"path"`
	HTTPHeaders map[string]string `json:"headers"`
	HTTPBody    string            `json:"body"`
}

type SetGlobalParameterOptions struct {
	Key   string `json:"key"`
	Value any    `json:"value"`
}

type GlobalParameters struct {
	Parameters map[string]any `json:"parameters"`
}
