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

	AddSubscriptionTemplate(opts AddSubscriptionTemplateOptions) SubscriptionTemplateResponse
	GetSubscriptionTemplate(id string) SubscriptionTemplateResponse
	ListSubscriptionTemplates() SubscriptionTemplatesResponse
	UpdateSubscriptionTemplate(id string, opts UpdateSubscriptionTemplateOptions) SubscriptionTemplateResponse
	DeleteSubscriptionTemplate(id string)

	SetGlobalParameter(parameter string, value any)
	GetGlobalParameter(parameter string) any
	ListGlobalParameters() GlobalParameters

	AddSubscriptionFromTemplate(opts AddSubscriptionFromTemplateOptions) SubscriptionResponse
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

func (a *apiClient) AddSubscriptionTemplate(opts AddSubscriptionTemplateOptions) (resp SubscriptionTemplateResponse) {
	a.client.doAssign(&resp, http.MethodPost, "/subscription-templates", bodyJson(opts), successStatuses(http.StatusCreated))

	return resp
}

func (a *apiClient) GetSubscriptionTemplate(id string) (resp SubscriptionTemplateResponse) {
	a.client.doAssign(&resp, http.MethodGet, fmt.Sprintf("/subscription-templates/%s", id), nil)

	return resp
}

func (a *apiClient) ListSubscriptionTemplates() (resp SubscriptionTemplatesResponse) {
	a.client.doAssign(&resp, http.MethodGet, "/subscription-templates", nil)

	return resp
}

func (a *apiClient) UpdateSubscriptionTemplate(id string, opts UpdateSubscriptionTemplateOptions) (resp SubscriptionTemplateResponse) {
	a.client.doAssign(&resp, http.MethodPut, fmt.Sprintf("/subscription-templates/%s", id), bodyJson(opts))

	return resp
}

func (a *apiClient) DeleteSubscriptionTemplate(id string) {
	a.client.do(http.MethodDelete, fmt.Sprintf("/subscription-templates/%s", id), nil)
}

func (a *apiClient) SetGlobalParameter(parameter string, value any) {
	a.client.do(http.MethodPost, "/global-parameters", bodyJson(SetGlobalParameterOptions{
		Parameter: parameter,
		Value:     value,
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

func (a *apiClient) AddSubscriptionFromTemplate(opts AddSubscriptionFromTemplateOptions) (resp SubscriptionResponse) {
	a.client.doAssign(&resp, http.MethodPost, "/subscriptions", bodyJson(opts))

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

	HTTPMethod       string            `json:"method"`
	HTTPURL          string            `json:"path"`
	HTTPHeaders      map[string]string `json:"httpHeaders"`
	HTTPBodyTemplate string            `json:"httpBodyTemplate"`

	TemplateID         *string        `json:"templateId"`
	TemplateParameters map[string]any `json:"templateParameters"`
}

type UpdateSubscriptionOptions struct {
	Name string `json:"name"`

	Topic string `json:"topic"`

	Extract map[string]string `json:"extract"`
	Filter  string            `json:"filter"`

	HTTPMethod       string            `json:"method"`
	HTTPURL          string            `json:"path"`
	HTTPHeaders      map[string]string `json:"headers"`
	HTTPBodyTemplate string            `json:"bodyTemplate"`
}

type SubscriptionTemplateResponse struct {
	SubscriptionTemplate SubscriptionTemplate `json:"subscriptionTemplate"`
}

type SubscriptionTemplatesResponse struct {
	SubscriptionTemplates []SubscriptionTemplate `json:"subscriptionTemplates"`
}

type SubscriptionTemplate struct {
	ID   string `json:"id"`
	Name string `json:"name"`

	Topic string `json:"topic"`

	Extract map[string]string `json:"extract"`
	Filter  string            `json:"filter"`

	HTTPMethod       string            `json:"method"`
	HTTPURL          string            `json:"path"`
	HTTPHeaders      map[string]string `json:"httpHeaders"`
	HTTPBodyTemplate string            `json:"httpBodyTemplate"`

	RequiredParameters []string `json:"requiredParameters"`
}

type AddSubscriptionTemplateOptions struct {
	Name string `json:"name"`

	Topic string `json:"topic"`

	Extract map[string]string `json:"extract"`
	Filter  string            `json:"filter"`

	HTTPMethod       string            `json:"method"`
	HTTPURL          string            `json:"url"`
	HTTPHeaders      map[string]string `json:"headers"`
	HTTPBodyTemplate string            `json:"bodyTemplate"`

	RequiredParameters []string `json:"requiredParameters"`
}

type UpdateSubscriptionTemplateOptions struct {
	Name string `json:"name"`

	Topic string `json:"topic"`

	Extract map[string]string `json:"extract"`
	Filter  string            `json:"filter"`

	HTTPMethod       string            `json:"method"`
	HTTPURL          string            `json:"path"`
	HTTPHeaders      map[string]string `json:"headers"`
	HTTPBodyTemplate string            `json:"bodyTemplate"`

	RequiredParameters []string `json:"requiredParameters"`
}

type SetGlobalParameterOptions struct {
	Parameter string `json:"parameter"`
	Value     any    `json:"value"`
}

type GlobalParameters struct {
	Parameters map[string]any `json:"parameters"`
}

type AddSubscriptionFromTemplateOptions struct {
	TemplateID         string         `json:"templateId"`
	TemplateParameters map[string]any `json:"templateParameters"`
}
