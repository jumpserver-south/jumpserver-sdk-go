package jumpserver

import (
	"crypto/tls"
	"net/http"
	"time"
)

// Option configures a [Client]. Pass any number of options to [NewClient].
type Option func(*clientConfig)

// clientConfig is the internal bag that Options mutate. It is applied to
// the Client after all options have been collected.
type clientConfig struct {
	baseURL    string
	httpClient *http.Client
	auth       Authenticator
	headers    map[string]string
	cookies    map[string]string
	userAgent  string
	orgID      string

	timeout            time.Duration
	maxRetries         int
	retryMinWait       time.Duration
	retryMaxWait       time.Duration
	insecureSkipVerify bool

	logger        Logger
	debugRequests bool
}

// Logger is a minimal interface satisfied by *log.Logger. Pass any
// compatible implementation to [WithLogger].
type Logger interface {
	Printf(format string, v ...any)
}

func defaultConfig() *clientConfig {
	return &clientConfig{
		baseURL:            "http://127.0.0.1:8080",
		userAgent:          userAgent,
		headers:            map[string]string{},
		cookies:            map[string]string{},
		orgID:              "00000000-0000-0000-0000-000000000002",
		timeout:            30 * time.Second,
		maxRetries:         3,
		retryMinWait:       500 * time.Millisecond,
		retryMaxWait:       15 * time.Second,
		insecureSkipVerify: true,
	}
}

// WithBaseURL sets the JumpServer base URL, e.g.
// "https://jumpserver.example.com". Trailing slashes are trimmed.
func WithBaseURL(u string) Option {
	return func(c *clientConfig) { c.baseURL = u }
}

// WithHTTPClient replaces the default *http.Client. When set, the
// timeout/retry/InsecureSkipVerify options are still respected for
// retry logic but the caller owns the transport.
func WithHTTPClient(hc *http.Client) Option {
	return func(c *clientConfig) { c.httpClient = hc }
}

// WithTimeout sets the per-request timeout.
func WithTimeout(d time.Duration) Option {
	return func(c *clientConfig) { c.timeout = d }
}

// WithUserAgent overrides the default User-Agent header.
func WithUserAgent(ua string) Option {
	return func(c *clientConfig) { c.userAgent = ua }
}

// WithHeader adds an extra header sent on every request.
func WithHeader(key, value string) Option {
	return func(c *clientConfig) { c.headers[key] = value }
}

// WithCookie adds a cookie sent on every request.
func WithCookie(name, value string) Option {
	return func(c *clientConfig) { c.cookies[name] = value }
}

// WithOrg sets the default X-JMS-ORG header value. JumpServer routes
// most endpoints through an organization; use "ROOT" (default) for the
// default organization, "00000000-0000-0000-0000-000000000002" for the
// global org, or a specific org ID.
func WithOrg(id string) Option {
	return func(c *clientConfig) { c.orgID = id }
}

// WithRetry configures automatic retries. Retries occur on:
//   - HTTP 408, 429, 500, 502, 503, 504 responses
//   - Transient network errors (timeout, connection reset, DNS temporary failure)
//
// Permanent errors (TLS failures, context cancellation) are never retried.
// Backoff uses exponential delay with full jitter, honouring Retry-After headers.
func WithRetry(maxRetries int, minWait, maxWait time.Duration) Option {
	return func(c *clientConfig) {
		if maxRetries < 0 {
			maxRetries = 0
		}
		c.maxRetries = maxRetries
		c.retryMinWait = minWait
		c.retryMaxWait = maxWait
	}
}

// WithInsecureSkipVerify controls TLS certificate verification. The
// default is true for compatibility with self-signed JumpServer
// deployments, but production code should set this to false.
func WithInsecureSkipVerify(skip bool) Option {
	return func(c *clientConfig) { c.insecureSkipVerify = skip }
}

// WithLogger enables request/response logging via the supplied Logger.
func WithLogger(l Logger) Option {
	return func(c *clientConfig) { c.logger = l }
}

// WithDebugRequests toggles verbose request/response dumps in the
// logger output.
func WithDebugRequests(on bool) Option {
	return func(c *clientConfig) { c.debugRequests = on }
}

// WithAccessKeyAuth authenticates using HMAC-SHA256 HTTP Signature
// (JumpServer "Access Key"). This is the recommended way to talk to
// the API with a service account.
func WithAccessKeyAuth(keyID, secretID string) Option {
	return func(c *clientConfig) {
		c.auth = &SignatureAuth{KeyID: keyID, SecretID: secretID}
	}
}

// WithBearerToken authenticates with "Authorization: Bearer <token>".
func WithBearerToken(token string) Option {
	return func(c *clientConfig) { c.auth = &BearerTokenAuth{Token: token} }
}

// WithPrivateToken authenticates with "Authorization: Token <token>",
// the legacy JumpServer private-token scheme.
func WithPrivateToken(token string) Option {
	return func(c *clientConfig) { c.auth = &PrivateTokenAuth{Token: token} }
}

// WithBasicAuth authenticates with HTTP Basic. Useful to obtain a
// Bearer token through the authentication endpoint.
func WithBasicAuth(username, password string) Option {
	return func(c *clientConfig) { c.auth = &BasicAuth{Username: username, Password: password} }
}

// WithAuthenticator installs a fully custom [Authenticator].
func WithAuthenticator(a Authenticator) Option {
	return func(c *clientConfig) { c.auth = a }
}

func (c *clientConfig) buildHTTPClient() *http.Client {
	if c.httpClient != nil {
		return c.httpClient
	}
	tr := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   10,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: c.insecureSkipVerify},
	}
	return &http.Client{Transport: tr, Timeout: c.timeout}
}
