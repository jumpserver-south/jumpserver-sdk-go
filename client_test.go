package jumpserver

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/jumpserver-south/jumpserver-sdk-go/model"
)

func newTestClient(t *testing.T, srv *httptest.Server) *Client {
	t.Helper()
	return NewClient(
		WithBaseURL(srv.URL),
		WithHTTPClient(srv.Client()),
		WithTimeout(5*time.Second),
		WithRetry(1, 10*time.Millisecond, 50*time.Millisecond),
		WithInsecureSkipVerify(true),
	)
}

func TestNewClientDefaults(t *testing.T) {
	c := NewClient(WithBaseURL("https://example.com"))
	if c.baseURL.String() != "https://example.com" {
		t.Errorf("unexpected base URL: %s", c.baseURL.String())
	}
	if c.Users == nil || c.Assets == nil || c.Accounts == nil {
		t.Error("services not wiring")
	}
}

func TestUsersService_Profile(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/users/profile/" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Header.Get("X-JMS-ORG") != "ROOT" {
			t.Errorf("missing default org header")
		}
		_ = json.NewEncoder(w).Encode(model.User{ID: "abc", Username: "alice", Name: "Alice"})
	}))
	defer srv.Close()

	c := newTestClient(t, srv)
	user, resp, err := c.Users.Profile(context.Background())
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if resp == nil || resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %v", resp)
	}
	if user.Username != "alice" {
		t.Errorf("unexpected user: %+v", user)
	}
}

func TestUsersService_List(t *testing.T) {
	var baseURL string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("limit") != "10" {
			t.Errorf("limit not forwarded: %s", r.URL.RawQuery)
		}
		if r.URL.Query().Get("search") != "alice" {
			t.Errorf("search not forwarded: %s", r.URL.RawQuery)
		}
		_ = json.NewEncoder(w).Encode(model.UserPage{
			Total:   1,
			NextURL: baseURL + "/api/v1/users/users/?limit=10&offset=10",
			Results: []model.User{{ID: "u1", Username: "alice"}},
		})
	}))
	defer srv.Close()
	baseURL = srv.URL

	c := newTestClient(t, srv)
	usrs, resp, err := c.Users.List(context.Background(),
		map[string]string{"search": "alice"},
		&ListOptions{Limit: 10},
	)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if len(usrs) != 1 || usrs[0].Username != "alice" {
		t.Errorf("unexpected users: %+v", usrs)
	}
	if !resp.HasNextPage() {
		t.Error("expected HasNextPage to be true")
	}
	if resp.Count != 1 {
		t.Errorf("Count should be 1, got %d", resp.Count)
	}
}

func TestAssetsService_Create(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		var req model.AssetRequest
		_ = json.NewDecoder(r.Body).Decode(&req)
		if req.Name != "db01" {
			t.Errorf("unexpected name: %s", req.Name)
		}
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(model.Asset{ID: "new", Name: "db01"})
	}))
	defer srv.Close()

	c := newTestClient(t, srv)
	asset, resp, err := c.Databases.Create(context.Background(), &model.AssetRequest{
		Name:      "db01",
		Address:   "10.0.0.1",
		Platform:  22,
		Protocols: []model.NamePort{{Name: "mysql", Port: 3306}},
	})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if resp.StatusCode != 201 {
		t.Errorf("expected 201, got %d", resp.StatusCode)
	}
	if asset.ID != "new" {
		t.Errorf("unexpected asset: %+v", asset)
	}
}

func TestAPIError_NotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"detail":"Not found."}`))
	}))
	defer srv.Close()

	c := newTestClient(t, srv)
	_, _, err := c.Users.Get(context.Background(), "missing")
	if err == nil {
		t.Fatal("expected error")
	}
	if !IsNotFound(err) {
		t.Errorf("expected IsNotFound, got %v", err)
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.Message != "Not found." {
		t.Errorf("unexpected message: %q", apiErr.Message)
	}
}

func TestSignatureAuth_SignsRequest(t *testing.T) {
	a := &SignatureAuth{KeyID: "kid", SecretID: "supersecret"}
	req, _ := http.NewRequest("GET", "http://example.com/api/v1/users/", nil)
	req.Header.Set("Date", "Sun, 09 Jun 2024 12:00:00 GMT")
	if err := a.Authenticate(req); err != nil {
		t.Fatalf("sign err: %v", err)
	}
	got := req.Header.Get("Authorization")
	if got == "" {
		t.Fatal("Authorization header not set")
	}
	if !bytes.Contains([]byte(got), []byte(`Signature keyId="kid"`)) {
		t.Errorf("unexpected Authorization: %s", got)
	}
	if !bytes.Contains([]byte(got), []byte("hmac-sha256")) {
		t.Errorf("algorithm missing: %s", got)
	}
}

func TestWithOrgScope(t *testing.T) {
	var seen string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seen = r.Header.Get("X-JMS-ORG")
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	c := newTestClient(t, srv)
	_, _, err := c.Users.Profile(context.Background())
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if seen != "ROOT" {
		t.Errorf("expected ROOT, got %q", seen)
	}

	scoped := c.WithOrgScope("org-123")
	_, _, err = scoped.Users.Profile(context.Background())
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if seen != "org-123" {
		t.Errorf("expected org-123, got %q", seen)
	}

	_, _, err = c.Users.Profile(context.Background())
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if seen != "ROOT" {
		t.Errorf("expected ROOT unchanged, got %q", seen)
	}
}

func TestWalkPages(t *testing.T) {
	var calls atomic.Int32
	var baseURL string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := calls.Add(1)
		offset := r.URL.Query().Get("offset")
		hasNext := offset == "" || offset == "0"
		body := map[string]any{
			"count":   3,
			"results": []map[string]string{{"id": "u1"}, {"id": "u2"}},
		}
		if hasNext {
			body["next"] = baseURL + "/api/v1/users/users/?limit=2&offset=2"
		}
		_ = json.NewEncoder(w).Encode(body)
		_ = n
	}))
	defer srv.Close()
	baseURL = srv.URL

	c := newTestClient(t, srv)
	var all []model.User
	err := WalkPages(context.Background(), &ListOptions{Limit: 2}, 0,
		func(ctx context.Context, opts *ListOptions) (*Response, error) {
			usrs, resp, err := c.Users.List(ctx, nil, opts)
			if err != nil {
				return resp, err
			}
			all = append(all, usrs...)
			return resp, nil
		})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if calls.Load() < 2 {
		t.Errorf("expected at least 2 calls, got %d", calls.Load())
	}
	if len(all) < 2 {
		t.Errorf("expected aggregated users, got %d", len(all))
	}
}

func TestDoRaw_StreamsBinary(t *testing.T) {
	payload := []byte{0x00, 0x01, 0x02, 0xff}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/octet-stream")
		_, _ = w.Write(payload)
	}))
	defer srv.Close()

	c := newTestClient(t, srv)
	buf := &bytes.Buffer{}
	req, _ := c.NewRequest(context.Background(), "GET", "/some/file", nil)
	_, err := c.DoRaw(context.Background(), req, buf)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !bytes.Equal(buf.Bytes(), payload) {
		t.Errorf("raw body mismatch")
	}
}

var _ io.Writer = (*bytes.Buffer)(nil)

func TestAuthEndpoint(t *testing.T) {
	var seenPath string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seenPath = r.URL.Path
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(model.Token{Token: "v4-token", User: "alice"})
	}))
	defer srv.Close()

	c := NewClient(
		WithBaseURL(srv.URL),
		WithHTTPClient(srv.Client()),
		WithInsecureSkipVerify(true),
	)

	tok, _, err := c.Auth.CreateToken(context.Background(), &model.TokenRequest{
		Username: "alice",
		Password: "secret",
	})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if tok.Token != "v4-token" {
		t.Errorf("unexpected token: %s", tok.Token)
	}
	if seenPath != "/api/v1/authentication/tokens/" {
		t.Errorf("expected tokens path, got %s", seenPath)
	}
}

func TestNewServicesWired(t *testing.T) {
	c := NewClient(WithBaseURL("https://example.com"))
	if c.AccountTemplates == nil {
		t.Error("AccountTemplates not wired")
	}
	if c.ChangeSecrets == nil {
		t.Error("ChangeSecrets not wired")
	}
	if c.AccountBackups == nil {
		t.Error("AccountBackups not wired")
	}
}
