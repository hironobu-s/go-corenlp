package connector

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type HttpClient struct {
	Endpoint   string
	Annotators []string
	Username   string
	Password   string

	ctx context.Context
}

func NewHttpClient(ctx context.Context, endpoint string) *HttpClient {
	if ctx == nil {
		ctx = context.Background()
	}
	return &HttpClient{
		Endpoint: endpoint,
		ctx:      ctx,
	}
}

func (c *HttpClient) buildRequest(text string) (*http.Request, error) {
	u, err := url.Parse(c.Endpoint)
	if err != nil {
		return nil, err
	}

	// parameters for CoreNLP server
	params := map[string]interface{}{
		"annotators":   strings.Join(c.Annotators, ","),
		"outputFormat": "json",
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

func (c *HttpClient) Run(text string) (response Response, err error) {
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
