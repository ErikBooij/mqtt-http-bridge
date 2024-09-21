package testsupport

import (
	"fmt"
	"net/http"
	"testing"
)

type MockServerClient interface {
	AssertRequestRecorded(method, path string, assertions ...Assertion)
	Reset()
}

func NewMockServerClient(t *testing.T, host string, port int) MockServerClient {
	return &mockServerClient{
		httpClient: newTestHTTPClient(t, fmt.Sprintf("http://%s:%d", host, port)),
		t:          t,
	}
}

type mockServerClient struct {
	httpClient *testHTTPClient

	t *testing.T
}

func (c *mockServerClient) AssertRequestRecorded(method, path string, assertions ...Assertion) {
	assert := assertion{
		minTimes: 1,
	}

	for _, a := range assertions {
		a(&assert)
	}

	verificationReq := verification{
		Assertion: verificationRequest{
			Method: method,
			Path:   path,
		},
		Times: verificationTimes{
			AtLeast: assert.minTimes,
		},
	}

	if assert.body != nil {
		verificationReq.Assertion.Body = verificationBody{
			Type:   "string",
			String: assert.body,
		}
	}

	if assert.headers != nil && len(*assert.headers) > 0 {
		verificationReq.Assertion.Headers = make(multiValueMap, 0, len(*assert.headers))

		for k, v := range *assert.headers {
			verificationReq.Assertion.Headers = append(verificationReq.Assertion.Headers, multiValueMapEntry{
				Name:   k,
				Values: v,
			})
		}
	}

	if assert.queryParams != nil && len(*assert.queryParams) > 0 {
		verificationReq.Assertion.QueryStringParameters = make(multiValueMap, 0, len(*assert.queryParams))

		for k, v := range *assert.queryParams {
			verificationReq.Assertion.QueryStringParameters = append(verificationReq.Assertion.QueryStringParameters, multiValueMapEntry{
				Name:   k,
				Values: v,
			})
		}
	}

	if assert.maxTimes != nil {
		verificationReq.Times.AtMost = *assert.maxTimes
	}

	c.httpClient.do(http.MethodPut, "/mockserver/verify", bodyJson(verificationReq), successStatuses(http.StatusAccepted))
}

func (c *mockServerClient) Reset() {
	c.t.Fatalf("Reset not implemented")
}

type RecordedRequest struct {
	Method  string
	Path    string
	Body    string
	Headers map[string]string
}

type verification struct {
	Assertion           verificationRequest `json:"httpRequest,omitempty"`
	Times               verificationTimes   `json:"times,omitempty"`
	MaxReturnedRequests int                 `json:"maximumNumberOfRequestToReturnInVerificationFailure,omitempty"`
}

type verificationRequest struct {
	Method                string           `json:"method"`
	Path                  string           `json:"path"`
	QueryStringParameters multiValueMap    `json:"queryStringParameters,omitempty"`
	Body                  verificationBody `json:"body"`
	Headers               multiValueMap    `json:"headers,omitempty"`
}

type verificationBody struct {
	Type   string `json:"type"`
	String any    `json:"string"`
}

type verificationTimes struct {
	AtLeast int `json:"atLeast"`
	AtMost  int `json:"atMost,omitempty"`
}

type multiValueMap []multiValueMapEntry

type multiValueMapEntry struct {
	Name   string   `json:"name"`
	Values []string `json:"values"`
}

type assertion struct {
	body        any
	headers     *map[string][]string
	queryParams *map[string][]string

	minTimes int // minTimes is required, maxTimes is not
	maxTimes *int
}

type Assertion func(*assertion)

func Body(body any) Assertion {
	return func(a *assertion) {
		a.body = &body
	}
}

func Header(name, value string) Assertion {
	return func(a *assertion) {
		if a.headers == nil {
			h := make(map[string][]string)
			a.headers = &h
		}

		if (*a.headers)[name] == nil {
			(*a.headers)[name] = []string{value}
		} else {
			(*a.headers)[name] = append((*a.headers)[name], value)
		}
	}
}

func QueryParam(name, value string) Assertion {
	return func(a *assertion) {
		if a.queryParams == nil {
			q := make(map[string][]string)
			a.queryParams = &q
		}

		if (*a.queryParams)[name] == nil {
			(*a.queryParams)[name] = []string{value}
		} else {
			(*a.queryParams)[name] = append((*a.queryParams)[name], value)
		}
	}
}

func Exactly(times int) Assertion {
	return func(a *assertion) {
		AtLeast(times)(a)
		AtMost(times)(a)
	}
}

func AtLeast(times int) Assertion {
	return func(a *assertion) {
		a.minTimes = times
	}
}

func AtMost(times int) Assertion {
	return func(a *assertion) {
		a.maxTimes = &times
	}
}
