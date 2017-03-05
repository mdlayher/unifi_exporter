// Package unifi implements a client for the Ubiquiti UniFi Controller v4 and
// v5 API.
package unifi

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
)

const (
	// Predefined content types for HTTP requests.
	formEncodedContentType = "application/x-www-form-urlencoded"
	jsonContentType        = "application/json;charset=UTF-8"

	// userAgent is the default user agent this package will report to the UniFi
	// Controller v4 API.
	userAgent = "github.com/mdlayher/unifi"
)

// InsecureHTTPClient creates a *http.Client which does not verify a UniFi
// Controller's certificate chain and hostname.
//
// Please think carefully before using this client: it should only be used
// with self-hosted, internal UniFi Controllers.
func InsecureHTTPClient(timeout time.Duration) *http.Client {
	return &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
}

// A Client is a client for the Ubiquiti UniFi Controller v4 API.
//
// Client.Login must be called and return a nil error before any additional
// actions can be performed with a Client.
type Client struct {
	UserAgent string

	apiURL *url.URL
	client *http.Client
}

// NewClient creates a new Client, using the input API address and an optional
// HTTP client.  If no HTTP client is specified, a default one will be used.
//
// If working with a self-hosted UniFi Controller which does not have a valid
// TLS certificate, InsecureHTTPClient can be used.
//
// Client.Login must be called and return a nil error before any additional
// actions can be performed with a Client.
func NewClient(addr string, client *http.Client) (*Client, error) {
	// Trim trailing slash to ensure sane path creation in other methods
	u, err := url.Parse(strings.TrimRight(addr, "/"))
	if err != nil {
		return nil, err
	}

	if client == nil {
		client = &http.Client{
			Timeout: 10 * time.Second,
		}
	}

	if client.Jar == nil {
		jar, err := cookiejar.New(nil)
		if err != nil {
			return nil, err
		}
		client.Jar = jar
	}

	c := &Client{
		UserAgent: userAgent,

		apiURL: u,
		client: client,
	}

	return c, nil
}

// Login authenticates against the UniFi Controller using the specified
// username and password.  Login must be called and return a nil error before
// any additional actions can be performed.
func (c *Client) Login(username string, password string) error {
	auth := &login{
		Username: username,
		Password: password,
	}

	req, err := c.newRequest(http.MethodPost, "/api/login", auth)
	if err != nil {
		return err
	}

	_, err = c.do(req, nil)
	return err
}

type login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// newRequest creates a new HTTP request, using the specified HTTP method and
// API endpoint. Additionally, it accepts a struct which can be marshaled to
// a JSON body.
func (c *Client) newRequest(method string, endpoint string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}
	u := c.apiURL.ResolveReference(rel)

	hasBody := method == http.MethodPost && body != nil
	var length int64

	// If performing a POST request and body parameters exist, encode
	// them now
	buf := bytes.NewBuffer(nil)
	if hasBody {
		if err := json.NewEncoder(buf).Encode(body); err != nil {
			return nil, err
		}
		length = int64(buf.Len())
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	// For POST requests, add proper headers
	if hasBody {
		req.Header.Add("Content-Type", formEncodedContentType)
		req.ContentLength = length
	}

	req.Header.Add("Accept", jsonContentType)
	req.Header.Add("User-Agent", c.UserAgent)

	return req, nil
}

// do performs an HTTP request using req and unmarshals the result onto v, if
// v is not nil.
func (c *Client) do(req *http.Request, v interface{}) (*http.Response, error) {
	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if err := checkResponse(res); err != nil {
		return res, err
	}

	// If no second parameter was passed, do not attempt to handle response
	if v == nil {
		return res, nil
	}

	return res, json.NewDecoder(res.Body).Decode(v)
}

// checkResponse checks for correct content type in a response and for non-200
// HTTP status codes, and returns any errors encountered.
func checkResponse(res *http.Response) error {
	if cType := res.Header.Get("Content-Type"); cType != jsonContentType {
		return fmt.Errorf("expected %q content type, but received %q", jsonContentType, cType)
	}

	// Check for 200-range status code
	if c := res.StatusCode; 200 <= c && c <= 299 {
		return nil
	}

	return fmt.Errorf("unexpected HTTP status code: %d", res.StatusCode)
}
