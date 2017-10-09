package client

import (
	"context"
	"net/url"

	"github.com/loadimpact/k6/api/v1"
)

var MetricsURL = &url.URL***REMOVED***Path: "/v1/metrics"***REMOVED***

func (c *Client) Metrics(ctx context.Context) (ret []v1.Metric, err error) ***REMOVED***
	return ret, c.call(ctx, "GET", MetricsURL, nil, &ret)
***REMOVED***
