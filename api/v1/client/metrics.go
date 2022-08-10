package client

import (
	"context"
	"net/http"
	"net/url"

	v1 "go.k6.io/k6/api/v1"
)

// Metrics returns the current metrics summary.
func (c *Client) Metrics(ctx context.Context) (ret []v1.Metric, err error) ***REMOVED***
	var resp v1.MetricsJSONAPI

	err = c.CallAPI(ctx, http.MethodGet, &url.URL***REMOVED***Path: "/v1/metrics"***REMOVED***, nil, &resp)
	if err != nil ***REMOVED***
		return ret, err
	***REMOVED***

	return resp.Metrics(), nil
***REMOVED***
