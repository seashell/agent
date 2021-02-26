package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/hashicorp/go-cleanhttp"
)

// Client provides a client to the Seashell API
type Client struct {
	config     Config
	headers    map[string]string
	httpClient *http.Client
}

// NewClient returns a new Seashell API client
func NewClient(config *Config) (*Client, error) {

	config = DefaultConfig().Merge(config)

	if !strings.HasPrefix(config.Address, "http") {
		config.Address = "http://" + config.Address
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}

	client := &Client{
		config:     *config,
		headers:    map[string]string{},
		httpClient: cleanhttp.DefaultClient(),
	}

	return client, nil
}

// WithHeaders returns a new Client that will use the specified headers
// in addition to those present in the original client in any future request.
// In case of a collision, the header in the original client will be overwritten.
func (c *Client) WithHeaders(headers map[string]string) *Client {

	nc := &Client{
		config:     c.config,
		headers:    c.headers,
		httpClient: c.httpClient,
	}

	for k, v := range headers {
		nc.headers[k] = v
	}

	return nc
}

func (c *Client) get(path string, id string, receiver interface{}) error {

	u, err := url.Parse(c.config.Address)
	if err != nil {
		return err
	}

	u.Path += path
	if len(id) > 0 {
		u.Path += "/" + id
	}

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return err
	}

	c.addHeaders(req)

	return c.do(req, receiver)
}

func (c *Client) post(path string, sender interface{}, receiver interface{}) error {

	u, err := url.Parse(c.config.Address)
	if err != nil {
		return err
	}

	u.Path += path

	b := &bytes.Buffer{}
	json.NewEncoder(b).Encode(sender)

	req, err := http.NewRequest("POST", u.String(), b)
	if err != nil {
		return err
	}

	c.addHeaders(req)

	return c.do(req, receiver)
}

func (c *Client) patch(id, path string, sender interface{}, receiver interface{}) error {

	base, err := url.Parse(c.config.Address)
	if err != nil {
		return err
	}

	base.Path += path
	base.Path += "/" + id

	b := &bytes.Buffer{}
	json.NewEncoder(b).Encode(sender)

	req, err := http.NewRequest("PATCH", base.String(), b)
	if err != nil {
		return err
	}

	c.addHeaders(req)
	return c.do(req, receiver)
}

func (c *Client) delete(id, path string, receiver interface{}) error {

	u, err := url.Parse(c.config.Address)
	if err != nil {
		return err
	}

	u.Path += path
	u.Path += "/" + id

	req, err := http.NewRequest("DELETE", u.String(), nil)
	if err != nil {
		return err
	}

	c.addHeaders(req)
	return c.do(req, receiver)
}

func (c *Client) addHeaders(req *http.Request) {
	for k, v := range c.headers {
		req.Header.Set(k, v)
	}
	req.Header.Set("Content-Type", "application/json")
}

func (c *Client) addQuery(filters map[string]string, req *http.Request) {
	q := req.URL.Query()

	if len(filters) > 0 {
		for k, v := range filters {
			q.Add(k, v)
		}
	}

	req.URL.RawQuery = q.Encode()
}

func (c *Client) do(req *http.Request, receiver interface{}) error {
	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if ok := res.StatusCode >= 200 && res.StatusCode < 300; !ok {
		resBody, _ := ioutil.ReadAll(res.Body)
		return fmt.Errorf("%v: %v", res.Status, string(resBody))
	}

	if receiver != nil {
		return json.NewDecoder(res.Body).Decode(receiver)
	}

	return nil
}
