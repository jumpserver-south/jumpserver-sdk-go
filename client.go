package jumpserver

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"maps"
	"math/rand/v2"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/jumpserver-south/jumpserver-sdk-go/accounts"
	"github.com/jumpserver-south/jumpserver-sdk-go/acls"
	"github.com/jumpserver-south/jumpserver-sdk-go/assets"
	"github.com/jumpserver-south/jumpserver-sdk-go/audits"
	"github.com/jumpserver-south/jumpserver-sdk-go/auth"
	"github.com/jumpserver-south/jumpserver-sdk-go/internal/core"
	"github.com/jumpserver-south/jumpserver-sdk-go/labels"
	"github.com/jumpserver-south/jumpserver-sdk-go/orgs"
	"github.com/jumpserver-south/jumpserver-sdk-go/perms"
	"github.com/jumpserver-south/jumpserver-sdk-go/rbac"
	"github.com/jumpserver-south/jumpserver-sdk-go/settings"
	"github.com/jumpserver-south/jumpserver-sdk-go/terminal"
	"github.com/jumpserver-south/jumpserver-sdk-go/tickets"
	"github.com/jumpserver-south/jumpserver-sdk-go/users"
	"github.com/jumpserver-south/jumpserver-sdk-go/xpack"
)

const (
	orgHeaderKey = "X-JMS-ORG"
)

// HTTPClient is the interface that service sub-packages use to make
// HTTP requests. *Client satisfies this interface.
type HTTPClient = core.HTTPClient

// Client talks to a JumpServer instance. Construct one with [NewClient].
type Client struct {
	baseURL    *url.URL
	httpClient *http.Client
	auth       Authenticator
	headers    map[string]string
	cookies    map[string]string
	userAgent  string
	orgID      string

	maxRetries   int
	retryMinWait time.Duration
	retryMaxWait time.Duration

	logger        Logger
	debugRequests bool

	Auth       *auth.Service
	Users      *users.Service
	UserGroups *users.GroupsService
	Roles      *rbac.Service
	Assets     *assets.AssetsService
	Hosts      *assets.CategoryService
	Devices    *assets.CategoryService
	Databases  *assets.CategoryService
	Webs       *assets.CategoryService
	Clouds     *assets.CategoryService
	Customs    *assets.CategoryService
	Nodes      *assets.NodesService
	Platforms  *assets.PlatformsService
	Zones      *assets.ZonesService
	Gateways         *assets.GatewaysService
	Labels           *labels.Service
	Accounts         *accounts.Service
	AccountTemplates *accounts.TemplatesService
	ChangeSecrets    *accounts.ChangeSecretService
	AccountBackups   *accounts.BackupService
	Organizations    *orgs.Service
	Permissions      *perms.Service
	CommandFilters   *acls.CommandFiltersService
	LoginACLs        *acls.LoginACLsService
	Audits           *audits.Service
	Terminal         *terminal.Service
	Tickets          *tickets.Service
	Settings         *settings.Service
	Xpack            *xpack.Service
}

// NewClient returns a ready-to-use [Client].
func NewClient(opts ...Option) *Client {
	cfg := defaultConfig()
	for _, o := range opts {
		o(cfg)
	}

	u, err := url.Parse(strings.TrimRight(cfg.baseURL, "/"))
	if err != nil {
		panic(fmt.Sprintf("jumpserver: invalid base URL %q: %v", cfg.baseURL, err))
	}

	c := &Client{
		baseURL:       u,
		httpClient:    cfg.buildHTTPClient(),
		auth:          cfg.auth,
		headers:       cfg.headers,
		cookies:       cfg.cookies,
		userAgent:     cfg.userAgent,
		orgID:         cfg.orgID,
		maxRetries:    cfg.maxRetries,
		retryMinWait:  cfg.retryMinWait,
		retryMaxWait:  cfg.retryMaxWait,
		logger:        cfg.logger,
		debugRequests: cfg.debugRequests,
	}

	c.initServices()
	return c
}

func (c *Client) initServices() {
	c.Auth = auth.NewService(c)
	c.Users = users.NewService(c)
	c.UserGroups = users.NewGroupsService(c)
	c.Roles = rbac.NewService(c)
	c.Assets = assets.NewAssetsService(c)
	c.Hosts = assets.NewCategoryService(c, "host")
	c.Devices = assets.NewCategoryService(c, "device")
	c.Databases = assets.NewCategoryService(c, "database")
	c.Webs = assets.NewCategoryService(c, "web")
	c.Clouds = assets.NewCategoryService(c, "cloud")
	c.Customs = assets.NewCategoryService(c, "custom")
	c.Nodes = assets.NewNodesService(c)
	c.Platforms = assets.NewPlatformsService(c)
	c.Zones = assets.NewZonesService(c)
	c.Gateways = assets.NewGatewaysService(c)
	c.Labels = labels.NewService(c)
	c.Accounts = accounts.NewService(c)
	c.AccountTemplates = accounts.NewTemplatesService(c)
	c.ChangeSecrets = accounts.NewChangeSecretService(c)
	c.AccountBackups = accounts.NewBackupService(c)
	c.Organizations = orgs.NewService(c)
	c.Permissions = perms.NewService(c)
	c.CommandFilters = acls.NewCommandFiltersService(c)
	c.LoginACLs = acls.NewLoginACLsService(c)
	c.Audits = audits.NewService(c)
	c.Terminal = terminal.NewService(c)
	c.Tickets = tickets.NewService(c)
	c.Settings = settings.NewService(c)
	c.Xpack = xpack.NewService(c)
}

// BaseURL returns the configured base URL.
func (c *Client) BaseURL() *url.URL { return c.baseURL }

// NewRequest builds an *http.Request against the client's base URL.
// Body is JSON-encoded when non-nil.
func (c *Client) NewRequest(ctx context.Context, method, path string, body any) (*http.Request, error) {
	u := *c.baseURL
	rel, err := url.Parse(path)
	if err != nil {
		return nil, fmt.Errorf("jumpserver: invalid request path %q: %w", path, err)
	}
	u = *u.ResolveReference(rel)

	var buf io.Reader
	if body != nil {
		raw, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("jumpserver: encode body: %w", err)
		}
		buf = bytes.NewReader(raw)
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.userAgent)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.orgID != "" {
		req.Header.Set(orgHeaderKey, c.orgID)
	}
	for k, v := range c.headers {
		req.Header.Set(k, v)
	}
	for k, v := range c.cookies {
		req.AddCookie(&http.Cookie{Name: k, Value: v})
	}
	if c.auth != nil {
		if err := c.auth.Authenticate(req); err != nil {
			return nil, fmt.Errorf("jumpserver: authenticate: %w", err)
		}
	}
	return req, nil
}

// Do sends req and decodes the response into v. When v is a
// *bytes.Buffer the raw body is copied into it instead.
func (c *Client) Do(ctx context.Context, req *http.Request, v any) (*Response, error) {
	if c.debugRequests && c.logger != nil {
		if dump, err := httputil.DumpRequestOut(req, true); err == nil {
			c.logger.Printf("jumpserver -> %s", dump)
		}
	}

	start := time.Now()
	resp, err := c.doWithRetry(ctx, req)
	elapsed := time.Since(start)
	if c.logger != nil {
		c.logger.Printf("jumpserver: %s %s -> %v (%s)", req.Method, req.URL, statusOf(resp), elapsed)
	}
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return wrapResponse(resp), err
	}
	if c.debugRequests && c.logger != nil {
		c.logger.Printf("jumpserver <- %s", body)
	}

	wrapped := wrapResponse(resp)
	if ct := resp.Header.Get("Content-Type"); strings.Contains(ct, "application/json") {
		parsePagination(body, wrapped)
	}

	if resp.StatusCode >= 400 {
		return wrapped, &APIError{
			StatusCode: resp.StatusCode,
			Method:     req.Method,
			URL:        req.URL.String(),
			Body:       body,
			Message:    extractAPIErrorMessage(body),
			Response:   resp,
		}
	}

	if v == nil || resp.StatusCode == http.StatusNoContent || len(body) == 0 {
		return wrapped, nil
	}
	if buf, ok := v.(*bytes.Buffer); ok {
		buf.Write(body)
		return wrapped, nil
	}
	if err := json.Unmarshal(body, v); err != nil {
		return wrapped, fmt.Errorf("jumpserver: decode %s %s: %w", req.Method, req.URL, err)
	}
	return wrapped, nil
}

// DoRaw sends req and streams the response body into w without JSON
// decoding. Use it for binary downloads (e.g. session replays, files).
func (c *Client) DoRaw(ctx context.Context, req *http.Request, w io.Writer) (*Response, error) {
	if c.debugRequests && c.logger != nil {
		if dump, err := httputil.DumpRequestOut(req, true); err == nil {
			c.logger.Printf("jumpserver -> %s", dump)
		}
	}
	start := time.Now()
	resp, err := c.doWithRetry(ctx, req)
	elapsed := time.Since(start)
	if c.logger != nil {
		c.logger.Printf("jumpserver: %s %s -> %v (%s)", req.Method, req.URL, statusOf(resp), elapsed)
	}
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	wrapped := wrapResponse(resp)
	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return wrapped, &APIError{
			StatusCode: resp.StatusCode,
			Method:     req.Method,
			URL:        req.URL.String(),
			Body:       body,
			Message:    extractAPIErrorMessage(body),
			Response:   resp,
		}
	}
	if w != nil {
		if _, err := io.Copy(w, resp.Body); err != nil {
			return wrapped, err
		}
	}
	return wrapped, nil
}

// ---------- Retry ----------

// doWithRetry executes req with automatic retries on transient network
// errors and 408/429/5xx responses.
//
// Key properties:
//   - Uses time.NewTimer with defer Stop() to prevent timer leaks
//   - Only retries on transient network errors (timeout, connection reset,
//     temporary DNS) — permanent errors like TLS cert issues are not retried
//   - Respects Retry-After header, capped at maxWait
//   - Uses math/rand/v2 for concurrency-safe jitter
func (c *Client) doWithRetry(ctx context.Context, req *http.Request) (*http.Response, error) {
	var bodyBytes []byte
	if req.Body != nil && req.Body != http.NoBody {
		var err error
		bodyBytes, err = io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		req.Body.Close()
	}

	var (
		resp *http.Response
		err  error
	)
	attempts := c.maxRetries + 1
	for i := range attempts {
		if bodyBytes != nil {
			req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
			req.ContentLength = int64(len(bodyBytes))
			req.GetBody = func() (io.ReadCloser, error) {
				return io.NopCloser(bytes.NewReader(bodyBytes)), nil
			}
		}

		resp, err = c.httpClient.Do(req)
		if err != nil {
			if i+1 >= attempts || !isTransientError(err) {
				return nil, err
			}
			if err := c.retrySleep(ctx, i, nil); err != nil {
				return nil, err
			}
			continue
		}

		if !isRetryableStatus(resp.StatusCode) {
			return resp, nil
		}
		// Drain and close the body so the connection can be reused.
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()

		if i+1 >= attempts {
			return resp, nil
		}
		if err := c.retrySleep(ctx, i, resp); err != nil {
			return resp, err
		}
	}
	return resp, nil
}

// retrySleep waits for the backoff duration, respecting context cancellation.
// Uses time.NewTimer with proper cleanup to prevent timer leaks.
func (c *Client) retrySleep(ctx context.Context, attempt int, resp *http.Response) error {
	wait := c.backoff(attempt, resp)
	timer := time.NewTimer(wait)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

// isTransientError reports whether err is a transient network error
// worth retrying. Permanent errors (TLS failures, DNS not found,
// invalid URL) are not retried.
func isTransientError(err error) bool {
	if err == nil {
		return false
	}
	// Context cancellation/deadline exceeded — never retry.
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return false
	}
	// Check net.Error for timeout / temporary flags.
	var netErr net.Error
	if errors.As(err, &netErr) {
		return netErr.Timeout()
	}
	// Check for common transient errors.
	var opErr *net.OpError
	if errors.As(err, &opErr) {
		return opErr.Temporary() || opErr.Timeout()
	}
	// DNS temporary failures, connection reset, etc.
	var dnsErr *net.DNSError
	if errors.As(err, &dnsErr) {
		return dnsErr.Temporary()
	}
	// Conservatively retry unknown wrapped errors that contain
	// "connection reset" or "broken pipe".
	msg := err.Error()
	return strings.Contains(msg, "connection reset") ||
		strings.Contains(msg, "broken pipe") ||
		strings.Contains(msg, "connection refused")
}

// isRetryableStatus reports whether the HTTP status code indicates a
// transient failure worth retrying.
func isRetryableStatus(code int) bool {
	switch code {
	case http.StatusRequestTimeout, http.StatusTooManyRequests,
		http.StatusInternalServerError, http.StatusBadGateway,
		http.StatusServiceUnavailable, http.StatusGatewayTimeout:
		return true
	}
	return false
}

// backoff computes the wait duration for the given attempt using
// exponential backoff with full jitter. Honours Retry-After headers,
// capped at maxWait to prevent server-induced delays.
func (c *Client) backoff(attempt int, resp *http.Response) time.Duration {
	minW := c.retryMinWait
	maxW := c.retryMaxWait
	if minW <= 0 {
		minW = 500 * time.Millisecond
	}
	if maxW <= 0 {
		maxW = 15 * time.Second
	}

	// Honour Retry-After if present, but cap at maxWait.
	if resp != nil {
		if d := parseRetryAfter(resp); d > 0 {
			if d > maxW {
				return maxW
			}
			return d
		}
	}

	// Exponential backoff: minW * 2^attempt, capped at maxW.
	wait := minW << attempt
	if wait > maxW {
		wait = maxW
	}
	// Full jitter: uniform random in [wait/2, wait].
	half := wait / 2
	return half + rand.N(half+1)
}

// parseRetryAfter honours either a delta-seconds or an HTTP-date value.
func parseRetryAfter(resp *http.Response) time.Duration {
	if resp == nil {
		return 0
	}
	v := resp.Header.Get("Retry-After")
	if v == "" {
		return 0
	}
	if n, err := strconv.Atoi(v); err == nil && n > 0 {
		return time.Duration(n) * time.Second
	}
	if t, err := time.Parse(http.TimeFormat, v); err == nil {
		if d := time.Until(t); d > 0 {
			return d
		}
	}
	return 0
}

// ---------- Helpers ----------

func wrapResponse(r *http.Response) *Response {
	if r == nil {
		return nil
	}
	return &Response{Response: r}
}

func statusOf(r *http.Response) string {
	if r == nil {
		return "<no response>"
	}
	return r.Status
}

func parsePagination(body []byte, resp *Response) {
	if len(body) == 0 || body[0] != '{' {
		return
	}
	var raw struct {
		Count    *int    `json:"count"`
		Next     *string `json:"next"`
		Previous *string `json:"previous"`
	}
	if err := json.Unmarshal(body, &raw); err != nil {
		return
	}
	if raw.Count != nil {
		resp.Count = *raw.Count
	}
	if raw.Next != nil {
		resp.NextURL = *raw.Next
	}
	if raw.Previous != nil {
		resp.PreviousURL = *raw.Previous
	}
}

// WithOrgScope returns a copy of c whose default X-JMS-ORG header is
// overridden to id.
func (c *Client) WithOrgScope(id string) *Client {
	cc := *c
	cc.headers = maps.Clone(c.headers)
	cc.headers[orgHeaderKey] = id
	cc.orgID = id
	cc.initServices()
	return &cc
}

// Compile-time check: *Client implements HTTPClient.
var _ HTTPClient = (*Client)(nil)
