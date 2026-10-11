package facebookautopost

import (
	"context"
	"errors"
	"net/http"

	"github.com/dualsubstrate/blog/internal/outbound"
)

func graphDo(ctx context.Context, c *Client, class outbound.RetryClass, endpoint string, newReq func() (*http.Request, error)) (*http.Response, []byte, error) {
	resp, err := outbound.Do(ctx, outbound.Config{
		HTTP:  c.HTTP,
		Name:  "facebook-graph",
		Class: class,
	}, newReq)
	if err != nil {
		var status *outbound.StatusError
		if errors.As(err, &status) {
			return nil, nil, newGraphHTTPError(endpoint, status.StatusCode, status.Body)
		}
		return nil, nil, err
	}
	defer resp.Body.Close()
	body, err := readGraphResponse(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		if endpoint == "" && resp.Request != nil && resp.Request.URL != nil {
			endpoint = resp.Request.URL.String()
		}
		return nil, nil, newGraphHTTPError(endpoint, resp.StatusCode, string(body))
	}
	return resp, body, nil
}

func graphStatus(ctx context.Context, c *Client, class outbound.RetryClass, endpoint string, newReq func() (*http.Request, error)) error {
	_, _, err := graphDo(ctx, c, class, endpoint, newReq)
	return err
}

func graphRead(ctx context.Context, c *Client, class outbound.RetryClass, endpoint string, newReq func() (*http.Request, error)) ([]byte, error) {
	_, body, err := graphDo(ctx, c, class, endpoint, newReq)
	return body, err
}
