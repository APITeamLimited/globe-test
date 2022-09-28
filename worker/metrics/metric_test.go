package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMetric(t *testing.T) ***REMOVED***
	t.Parallel()
	testdata := map[string]struct ***REMOVED***
		Type     MetricType
		SinkType Sink
	***REMOVED******REMOVED***
		"Counter": ***REMOVED***Counter, &CounterSink***REMOVED******REMOVED******REMOVED***,
		"Gauge":   ***REMOVED***Gauge, &GaugeSink***REMOVED******REMOVED******REMOVED***,
		"Trend":   ***REMOVED***Trend, &TrendSink***REMOVED******REMOVED******REMOVED***,
		"Rate":    ***REMOVED***Rate, &RateSink***REMOVED******REMOVED******REMOVED***,
	***REMOVED***

	for name, data := range testdata ***REMOVED***
		name, data := name, data
		t.Run(name, func(t *testing.T) ***REMOVED***
			t.Parallel()
			m := newMetric("my_metric", data.Type)
			assert.Equal(t, "my_metric", m.Name)
			assert.IsType(t, data.SinkType, m.Sink)
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestAddSubmetric(t *testing.T) ***REMOVED***
	t.Parallel()
	testdata := map[string]struct ***REMOVED***
		err  bool
		tags map[string]string
	***REMOVED******REMOVED***
		"":                        ***REMOVED***true, nil***REMOVED***,
		"  ":                      ***REMOVED***true, nil***REMOVED***,
		"a":                       ***REMOVED***false, map[string]string***REMOVED***"a": ""***REMOVED******REMOVED***,
		"a:1":                     ***REMOVED***false, map[string]string***REMOVED***"a": "1"***REMOVED******REMOVED***,
		" a : 1 ":                 ***REMOVED***false, map[string]string***REMOVED***"a": "1"***REMOVED******REMOVED***,
		"a,b":                     ***REMOVED***false, map[string]string***REMOVED***"a": "", "b": ""***REMOVED******REMOVED***,
		` a:"",b: ''`:             ***REMOVED***false, map[string]string***REMOVED***"a": "", "b": ""***REMOVED******REMOVED***,
		`a:1,b:2`:                 ***REMOVED***false, map[string]string***REMOVED***"a": "1", "b": "2"***REMOVED******REMOVED***,
		` a : 1, b : 2 `:          ***REMOVED***false, map[string]string***REMOVED***"a": "1", "b": "2"***REMOVED******REMOVED***,
		`a : '1' , b : "2"`:       ***REMOVED***false, map[string]string***REMOVED***"a": "1", "b": "2"***REMOVED******REMOVED***,
		`" a" : ' 1' , b : "2 " `: ***REMOVED***false, map[string]string***REMOVED***" a": " 1", "b": "2 "***REMOVED******REMOVED***, //nolint:gocritic
	***REMOVED***

	for name, expected := range testdata ***REMOVED***
		name, expected := name, expected
		t.Run(name, func(t *testing.T) ***REMOVED***
			t.Parallel()

			m := newMetric("metric", Trend)
			sm, err := m.AddSubmetric(name)
			if expected.err ***REMOVED***
				require.Error(t, err)
				return
			***REMOVED***
			require.NoError(t, err)
			require.NotNil(t, sm)
			assert.EqualValues(t, expected.tags, sm.Tags.tags)
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestParseMetricName(t *testing.T) ***REMOVED***
	t.Parallel()

	tests := []struct ***REMOVED***
		name                 string
		metricNameExpression string
		wantMetricName       string
		wantTags             []string
		wantErr              bool
	***REMOVED******REMOVED***
		***REMOVED***
			name:                 "metric name without tags",
			metricNameExpression: "test_metric",
			wantMetricName:       "test_metric",
			wantErr:              false,
		***REMOVED***,
		***REMOVED***
			name:                 "metric name with single tag",
			metricNameExpression: "test_metric***REMOVED***abc:123***REMOVED***",
			wantMetricName:       "test_metric",
			wantTags:             []string***REMOVED***"abc:123"***REMOVED***,
			wantErr:              false,
		***REMOVED***,
		***REMOVED***
			name:                 "metric name with multiple tags",
			metricNameExpression: "test_metric***REMOVED***abc:123,easyas:doremi***REMOVED***",
			wantMetricName:       "test_metric",
			wantTags:             []string***REMOVED***"abc:123", "easyas:doremi"***REMOVED***,
			wantErr:              false,
		***REMOVED***,
		***REMOVED***
			name:                 "metric name with multiple spaced tags",
			metricNameExpression: "test_metric***REMOVED***abc:123, easyas:doremi***REMOVED***",
			wantMetricName:       "test_metric",
			wantTags:             []string***REMOVED***"abc:123", "easyas:doremi"***REMOVED***,
			wantErr:              false,
		***REMOVED***,
		***REMOVED***
			name:                 "metric name with group tag",
			metricNameExpression: "test_metric***REMOVED***group:::mygroup***REMOVED***",
			wantMetricName:       "test_metric",
			wantTags:             []string***REMOVED***"group:::mygroup"***REMOVED***,
			wantErr:              false,
		***REMOVED***,
		***REMOVED***
			name:                 "metric name with valid name and repeated curly braces tokens in tags definition",
			metricNameExpression: "http_req_duration***REMOVED***name:http://$***REMOVED******REMOVED***.com***REMOVED***",
			wantMetricName:       "http_req_duration",
			wantTags:             []string***REMOVED***"name:http://$***REMOVED******REMOVED***.com"***REMOVED***,
			wantErr:              false,
		***REMOVED***,
		***REMOVED***
			name:                 "metric name with valid name and repeated curly braces and colon tokens in tags definition",
			metricNameExpression: "http_req_duration***REMOVED***name:http://$***REMOVED******REMOVED***.com,url:ssh://github.com:grafana/k6***REMOVED***",
			wantMetricName:       "http_req_duration",
			wantTags:             []string***REMOVED***"name:http://$***REMOVED******REMOVED***.com", "url:ssh://github.com:grafana/k6"***REMOVED***,
			wantErr:              false,
		***REMOVED***,
		***REMOVED***
			name:                 "metric name with tag definition missing `:value`",
			metricNameExpression: "test_metric***REMOVED***easyas***REMOVED***",
			wantErr:              true,
		***REMOVED***,
		***REMOVED***
			name:                 "metric name with tag definition missing value",
			metricNameExpression: "test_metric***REMOVED***easyas:***REMOVED***",
			wantErr:              true,
		***REMOVED***,
		***REMOVED***
			name:                 "metric name with mixed valid and invalid tag definitions",
			metricNameExpression: "test_metric***REMOVED***abc:123,easyas:***REMOVED***",
			wantErr:              true,
		***REMOVED***,
		***REMOVED***
			name:                 "metric name with valid name and unmatched opening tags definition token",
			metricNameExpression: "test_metric***REMOVED***abc:123,easyas:doremi",
			wantErr:              true,
		***REMOVED***,
		***REMOVED***
			name:                 "metric name with valid name and unmatched closing tags definition token",
			metricNameExpression: "test_metricabc:123,easyas:doremi***REMOVED***",
			wantErr:              true,
		***REMOVED***,
		***REMOVED***
			name:                 "metric name with valid name and invalid starting tags definition token",
			metricNameExpression: "test_metric***REMOVED***abc:123,easyas:doremi***REMOVED***",
			wantErr:              true,
		***REMOVED***,
		***REMOVED***
			name:                 "metric name with valid name and invalid curly braces in tags definition",
			metricNameExpression: "test_metric***REMOVED***abc***REMOVED***bar",
			wantErr:              true,
		***REMOVED***,
		***REMOVED***
			name:                 "metric name with valid name and trailing characters after closing curly brace in tags definition",
			metricNameExpression: "test_metric***REMOVED***foo:ba***REMOVED***r",
			wantErr:              true,
		***REMOVED***,
	***REMOVED***
	for _, tt := range tests ***REMOVED***
		tt := tt

		t.Run(tt.name, func(t *testing.T) ***REMOVED***
			t.Parallel()

			gotMetricName, gotTags, gotErr := ParseMetricName(tt.metricNameExpression)

			assert.Equal(t,
				gotErr != nil, tt.wantErr,
				"ParseMetricName() error = %v, wantErr %v", gotErr, tt.wantErr,
			)

			if gotErr != nil ***REMOVED***
				assert.ErrorIs(t,
					gotErr, ErrMetricNameParsing,
					"ParseMetricName() error chain should contain ErrMetricNameParsing",
				)
			***REMOVED***

			assert.Equal(t,
				gotMetricName, tt.wantMetricName,
				"ParseMetricName() gotMetricName = %v, want %v", gotMetricName, tt.wantMetricName,
			)

			assert.Equal(t,
				gotTags, tt.wantTags,
				"ParseMetricName() gotTags = %v, want %v", gotTags, tt.wantTags,
			)
		***REMOVED***)
	***REMOVED***
***REMOVED***
