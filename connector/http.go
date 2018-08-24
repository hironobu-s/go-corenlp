package connector

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// HTTPClient is a client of Stanford CoreNLP server.
// https://stanfordnlp.github.io/CoreNLP/corenlp-server.html
type HTTPClient struct {
	Endpoint   string
	Annotators []string
	Username   string
	Password   string

	ctx context.Context
}

// NewHTTPClient initializes HttpClient and returns it.
func NewHTTPClient(ctx context.Context, endpoint string) *HTTPClient {
	if ctx == nil {
		ctx = context.Background()
	}
	return &HTTPClient{
		Endpoint: endpoint,
		ctx:      ctx,
	}
}

func (c *HTTPClient) buildRequest(text string) (*http.Request, error) {
	u, err := url.Parse(c.Endpoint)
	if err != nil {
		return nil, err
	}

	// parameters for CoreNLP server
	params := map[string]interface{}{
		"outputFormat": "json",
	}

	if len(c.Annotators) > 0 {
		params["annotators"] = strings.Join(c.Annotators, ",")
	}

	deadline, ok := c.ctx.Deadline()
	if ok {
		d := deadline.Sub(time.Now()).Seconds()
		params["timeout"] = d / float64(time.Millisecond)
	}

	if c.Username != "" && c.Password != "" {
		params["username"] = c.Username
		params["password"] = c.Password
	}

	props, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	// build http request
	q := u.Query()
	q.Add("properties", string(props))
	u.RawQuery = q.Encode()

	req, err := http.NewRequest("POST", u.String(), strings.NewReader(text))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	return req, err
}

// Run marshals Connector interface implementation.
func (c *HTTPClient) Run(text string) (response Response, err error) {
	req, err := c.buildRequest(text)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(c.ctx)

	tr := &http.Transport{}
	cli := &http.Client{Transport: tr}

	rs, err := cli.Do(req)
	if err == nil {
		response = rs.Body
	}
	return response, err
}
