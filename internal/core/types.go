package core

import (
	"context"
	"io"
	"net/http"
	"strconv"
)

// HTTPClient is the interface that service sub-packages use to make
// HTTP requests. The root *Client satisfies this interface.
type HTTPClient interface {
	NewRequest(ctx context.Context, method, path string, body any) (*http.Request, error)
	Do(ctx context.Context, req *http.Request, v any) (*Response, error)
	DoRaw(ctx context.Context, req *http.Request, w io.Writer) (*Response, error)
}

// Response is the typed wrapper returned by every service call. It
// carries the underlying *http.Response alongside JumpServer's
// pagination metadata for list endpoints.
type Response struct {
	*http.Response

	Count       int
	NextURL     string
	PreviousURL string
}

// HasNextPage reports whether a subsequent page is available.
func (r *Response) HasNextPage() bool {
	return r != nil && r.NextURL != ""
}

// HasPreviousPage reports whether a previous page is available.
func (r *Response) HasPreviousPage() bool {
	return r != nil && r.PreviousURL != ""
}

// ListOptions configures pagination and search for list endpoints.
type ListOptions struct {
	Limit  int    `url:"limit,omitempty"`
	Offset int    `url:"offset,omitempty"`
	Search string `url:"search,omitempty"`
	Order  string `url:"order,omitempty"`
}

// Apply appends the list options to a params map.
func (o *ListOptions) Apply(v map[string]string) {
	if o == nil {
		return
	}
	if o.Limit > 0 {
		v["limit"] = strconv.Itoa(o.Limit)
	}
	if o.Offset > 0 {
		v["offset"] = strconv.Itoa(o.Offset)
	}
	if o.Search != "" {
		v["search"] = o.Search
	}
	if o.Order != "" {
		v["order"] = o.Order
	}
}

// Next returns a copy of o with Offset advanced by Limit.
func (o *ListOptions) Next() *ListOptions {
	if o == nil || o.Limit <= 0 {
		return nil
	}
	next := *o
	next.Offset += next.Limit
	return &next
}

// Page is a generic paginated list envelope matching JumpServer's
// {count, next, previous, results} shape. It is used internally by
// the crud helpers to decode list responses generically.
type Page[T any] struct {
	Total       int    `json:"count"`
	NextURL     string `json:"next"`
	PreviousURL string `json:"previous"`
	Results     []T    `json:"results"`
}

// PageFetcher paginates through all pages of a list endpoint.
type PageFetcher func(ctx context.Context, opts *ListOptions) (*Response, error)
