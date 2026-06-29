package sdkutil

import (
	"context"

	"github.com/jumpserver-south/jumpserver-sdk-go/internal/core"
)

// List fetches a paginated list by decoding the response into a
// model.XXXPage type. The caller passes a pointer to a typed Page
// struct (e.g. *model.UserPage) so the decoder fills in Results,
// Total, NextURL, and PreviousURL. This helper copies the page
// metadata into resp and returns just the results slice.
func List[T any](ctx context.Context, client core.HTTPClient, listURL string, opts *core.ListOptions) ([]T, *core.Response, error) {
	params := map[string]string{}
	if opts != nil {
		opts.Apply(params)
	}
	path := AppendQuery(listURL, params)
	httpReq, err := client.NewRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, nil, err
	}
	var page core.Page[T]
	resp, err := client.Do(ctx, httpReq, &page)
	if err != nil {
		return nil, resp, err
	}
	if resp != nil {
		resp.Count = page.Total
		resp.NextURL = page.NextURL
		resp.PreviousURL = page.PreviousURL
	}
	return page.Results, resp, nil
}

// Get fetches a single resource by ID and returns a pointer to it.
func Get[T any](ctx context.Context, client core.HTTPClient, detailURL string, id any) (*T, *core.Response, error) {
	httpReq, err := client.NewRequest(ctx, "GET", Spath(detailURL, id), nil)
	if err != nil {
		return nil, nil, err
	}
	var out T
	resp, err := client.Do(ctx, httpReq, &out)
	if err != nil {
		return nil, resp, err
	}
	return &out, resp, nil
}

// Create POSTs req to listURL and returns a pointer to the created resource.
func Create[T, R any](ctx context.Context, client core.HTTPClient, listURL string, req *R) (*T, *core.Response, error) {
	httpReq, err := client.NewRequest(ctx, "POST", listURL, req)
	if err != nil {
		return nil, nil, err
	}
	var out T
	resp, err := client.Do(ctx, httpReq, &out)
	if err != nil {
		return nil, resp, err
	}
	return &out, resp, nil
}

// Update PATCHes a resource by ID and returns a pointer to the updated resource.
func Update[T, R any](ctx context.Context, client core.HTTPClient, detailURL string, id any, req *R) (*T, *core.Response, error) {
	httpReq, err := client.NewRequest(ctx, "PATCH", Spath(detailURL, id), req)
	if err != nil {
		return nil, nil, err
	}
	var out T
	resp, err := client.Do(ctx, httpReq, &out)
	if err != nil {
		return nil, resp, err
	}
	return &out, resp, nil
}

// Replace PUTs a full resource by ID and returns a pointer to the replaced resource.
func Replace[T, R any](ctx context.Context, client core.HTTPClient, detailURL string, id any, req *R) (*T, *core.Response, error) {
	httpReq, err := client.NewRequest(ctx, "PUT", Spath(detailURL, id), req)
	if err != nil {
		return nil, nil, err
	}
	var out T
	resp, err := client.Do(ctx, httpReq, &out)
	if err != nil {
		return nil, resp, err
	}
	return &out, resp, nil
}

// Delete sends a DELETE request for the resource at detailURL/id.
func Delete(ctx context.Context, client core.HTTPClient, detailURL string, id any) (*core.Response, error) {
	httpReq, err := client.NewRequest(ctx, "DELETE", Spath(detailURL, id), nil)
	if err != nil {
		return nil, err
	}
	return client.Do(ctx, httpReq, nil)
}

// Action POSTs a body to an action URL and decodes into T.
func Action[T, R any](ctx context.Context, client core.HTTPClient, actionURL string, body *R) (*T, *core.Response, error) {
	httpReq, err := client.NewRequest(ctx, "POST", actionURL, body)
	if err != nil {
		return nil, nil, err
	}
	var out T
	resp, err := client.Do(ctx, httpReq, &out)
	if err != nil {
		return nil, resp, err
	}
	return &out, resp, nil
}

// MapAction POSTs a body to an action URL and returns a map[string]any response.
func MapAction(ctx context.Context, client core.HTTPClient, actionURL string, body any) (map[string]any, *core.Response, error) {
	httpReq, err := client.NewRequest(ctx, "POST", actionURL, body)
	if err != nil {
		return nil, nil, err
	}
	out := map[string]any{}
	resp, err := client.Do(ctx, httpReq, &out)
	if err != nil {
		return nil, resp, err
	}
	return out, resp, nil
}
