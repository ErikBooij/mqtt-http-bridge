package testsupport

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"mqtt-http-bridge/src/utilities"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"testing"
)

type testHTTPClient struct {
	baseURL string
	client  *http.Client
	t       *testing.T
}

func newTestHTTPClient(t *testing.T, baseURL string) *testHTTPClient {
	return &testHTTPClient{
		baseURL: strings.TrimRight(baseURL, "/"),
		client:  &http.Client{},
		t:       t,
	}
}

type requestOptions struct {
	body    []byte
	headers map[string]string

	successStatuses []int
}

type requestOption func(*testing.T, *requestOptions)

func header(key, value string) requestOption {
	return func(t *testing.T, o *requestOptions) {
		if o.headers == nil {
			o.headers = make(map[string]string)
		}

		o.headers[key] = value
	}
}

func bodyJson(body any) requestOption {
	return func(t *testing.T, o *requestOptions) {
		b, err := json.Marshal(body)

		if err != nil {
			t.Fatalf("Error marshalling JSON: %s", err)
		}

		o.body = b

		header("Content-Type", "application/json")(t, o)
	}
}

//lint:ignore U1000 // Keeping this around as a convenience for future tests
func bodyText(body string) requestOption {
	return func(t *testing.T, o *requestOptions) {
		o.body = []byte(body)
	}
}

func successStatuses(statuses ...int) requestOption {
	return func(t *testing.T, o *requestOptions) {
		o.successStatuses = statuses
	}
}

func (c *testHTTPClient) doAssign(target interface{}, method, path string, options ...requestOption) {
	c.t.Helper()

	resp := c.do(method, path, options...)

	if resp.Body != nil {
		defer resp.Body.Close()

		if err := json.NewDecoder(resp.Body).Decode(target); err != nil && !errors.Is(err, io.EOF) {
			c.t.Fatalf("Error decoding response: %s", err)
		}
	}
}

func (c *testHTTPClient) do(method, path string, options ...requestOption) *http.Response {
	c.t.Helper()

	opts := &requestOptions{
		successStatuses: []int{http.StatusOK},
	}

	for _, opt := range options {
		opt(c.t, opts)
	}

	var body io.Reader

	if opts.body != nil {
		body = bytes.NewReader(opts.body)
	}

	req, err := http.NewRequest(method, c.baseURL+"/"+strings.TrimLeft(path, "/"), body)

	if err != nil {
		c.t.Fatalf("Error creating request: %s", err)
	}

	if opts.headers != nil {
		for k, v := range opts.headers {
			req.Header.Add(k, v)
		}
	}

	resp, err := c.client.Do(req)

	if err != nil {
		c.t.Fatalf("Error making request: %s", err)
	}

	if !slices.Contains(opts.successStatuses, resp.StatusCode) {
		body, _ := io.ReadAll(resp.Body)
		c.t.Fatalf("Unexpected status code: %d (expected %s)\n\nBody: %s", resp.StatusCode, strings.Join(utilities.MapSlice(opts.successStatuses, strconv.Itoa), "/"), string(body))
	}

	return resp
}
