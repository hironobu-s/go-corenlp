package provider

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

func (p *HttpClient) buildRequest(text string) (*http.Request, error) {
	u, err := url.Parse(p.Endpoint)
	if err != nil {
		return nil, err
	}

	// parameters for CoreNLP server
	params := map[string]interface{}{
		"annotators":   strings.Join(p.Annotators, ","),
		"outputFormat": "json",
	}

	deadline, ok := p.ctx.Deadline()
	if ok {
		d := deadline.Sub(time.Now()).Seconds()
		params["timeout"] = d / float64(time.Millisecond)
	}

	if p.Username != "" && p.Password != "" {
		params["username"] = p.Username
		params["password"] = p.Password
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

func (p *HttpClient) Run(text string) (response Response, err error) {
	req, err := p.buildRequest(text)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(p.ctx)

	tr := &http.Transport{}
	cli := &http.Client{Transport: tr}

	rs, err := cli.Do(req)
	if err == nil {
		response = rs.Body
	}
	return response, err
}
