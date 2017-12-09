/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2016 Load Impact
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package lib

import (
	"crypto/tls"
	"encoding/json"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/loadimpact/k6/stats"
	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v3"
)

func TestOptions(t *testing.T) ***REMOVED***
	t.Run("Paused", func(t *testing.T) ***REMOVED***
		opts := Options***REMOVED******REMOVED***.Apply(Options***REMOVED***Paused: null.BoolFrom(true)***REMOVED***)
		assert.True(t, opts.Paused.Valid)
		assert.True(t, opts.Paused.Bool)
	***REMOVED***)
	t.Run("VUs", func(t *testing.T) ***REMOVED***
		opts := Options***REMOVED******REMOVED***.Apply(Options***REMOVED***VUs: null.IntFrom(12345)***REMOVED***)
		assert.True(t, opts.VUs.Valid)
		assert.Equal(t, int64(12345), opts.VUs.Int64)
	***REMOVED***)
	t.Run("VUsMax", func(t *testing.T) ***REMOVED***
		opts := Options***REMOVED******REMOVED***.Apply(Options***REMOVED***VUsMax: null.IntFrom(12345)***REMOVED***)
		assert.True(t, opts.VUsMax.Valid)
		assert.Equal(t, int64(12345), opts.VUsMax.Int64)
	***REMOVED***)
	t.Run("Duration", func(t *testing.T) ***REMOVED***
		opts := Options***REMOVED******REMOVED***.Apply(Options***REMOVED***Duration: NullDurationFrom(2 * time.Minute)***REMOVED***)
		assert.True(t, opts.Duration.Valid)
		assert.Equal(t, "2m0s", opts.Duration.String())
	***REMOVED***)
	t.Run("Iterations", func(t *testing.T) ***REMOVED***
		opts := Options***REMOVED******REMOVED***.Apply(Options***REMOVED***Iterations: null.IntFrom(1234)***REMOVED***)
		assert.True(t, opts.Iterations.Valid)
		assert.Equal(t, int64(1234), opts.Iterations.Int64)
	***REMOVED***)
	t.Run("Stages", func(t *testing.T) ***REMOVED***
		opts := Options***REMOVED******REMOVED***.Apply(Options***REMOVED***Stages: []Stage***REMOVED******REMOVED***Duration: NullDurationFrom(1 * time.Second)***REMOVED******REMOVED******REMOVED***)
		assert.NotNil(t, opts.Stages)
		assert.Len(t, opts.Stages, 1)
		assert.Equal(t, 1*time.Second, time.Duration(opts.Stages[0].Duration.Duration))
	***REMOVED***)
	t.Run("MaxRedirects", func(t *testing.T) ***REMOVED***
		opts := Options***REMOVED******REMOVED***.Apply(Options***REMOVED***MaxRedirects: null.IntFrom(12345)***REMOVED***)
		assert.True(t, opts.MaxRedirects.Valid)
		assert.Equal(t, int64(12345), opts.MaxRedirects.Int64)
	***REMOVED***)
	t.Run("InsecureSkipTLSVerify", func(t *testing.T) ***REMOVED***
		opts := Options***REMOVED******REMOVED***.Apply(Options***REMOVED***InsecureSkipTLSVerify: null.BoolFrom(true)***REMOVED***)
		assert.True(t, opts.InsecureSkipTLSVerify.Valid)
		assert.True(t, opts.InsecureSkipTLSVerify.Bool)
	***REMOVED***)
	t.Run("TLSCipherSuites", func(t *testing.T) ***REMOVED***
		for suiteName, suiteID := range SupportedTLSCipherSuites ***REMOVED***
			t.Run(suiteName, func(t *testing.T) ***REMOVED***
				opts := Options***REMOVED******REMOVED***.Apply(Options***REMOVED***TLSCipherSuites: &TLSCipherSuites***REMOVED***suiteID***REMOVED******REMOVED***)

				assert.NotNil(t, opts.TLSCipherSuites)
				assert.Len(t, *(opts.TLSCipherSuites), 1)
				assert.Equal(t, suiteID, (*opts.TLSCipherSuites)[0])
			***REMOVED***)
		***REMOVED***
	***REMOVED***)
	t.Run("TLSVersion", func(t *testing.T) ***REMOVED***
		versions := TLSVersions***REMOVED***Min: tls.VersionSSL30, Max: tls.VersionTLS12***REMOVED***
		opts := Options***REMOVED******REMOVED***.Apply(Options***REMOVED***TLSVersion: &versions***REMOVED***)

		assert.NotNil(t, opts.TLSVersion)
		assert.Equal(t, opts.TLSVersion.Min, TLSVersion(tls.VersionSSL30))
		assert.Equal(t, opts.TLSVersion.Max, TLSVersion(tls.VersionTLS12))

		t.Run("JSON", func(t *testing.T) ***REMOVED***
			t.Run("Object", func(t *testing.T) ***REMOVED***
				var opts Options
				jsonStr := `***REMOVED***"tlsVersion":***REMOVED***"min":"ssl3.0","max":"tls1.2"***REMOVED******REMOVED***`
				assert.NoError(t, json.Unmarshal([]byte(jsonStr), &opts))
				assert.Equal(t, &TLSVersions***REMOVED***
					Min: TLSVersion(tls.VersionSSL30),
					Max: TLSVersion(tls.VersionTLS12),
				***REMOVED***, opts.TLSVersion)

				t.Run("Roundtrip", func(t *testing.T) ***REMOVED***
					data, err := json.Marshal(opts.TLSVersion)
					assert.NoError(t, err)
					assert.Equal(t, `***REMOVED***"min":"ssl3.0","max":"tls1.2"***REMOVED***`, string(data))
					var vers2 TLSVersions
					assert.NoError(t, json.Unmarshal(data, &vers2))
					assert.Equal(t, &vers2, opts.TLSVersion)
				***REMOVED***)
			***REMOVED***)
			t.Run("String", func(t *testing.T) ***REMOVED***
				var opts Options
				jsonStr := `***REMOVED***"tlsVersion":"tls1.2"***REMOVED***`
				assert.NoError(t, json.Unmarshal([]byte(jsonStr), &opts))
				assert.Equal(t, &TLSVersions***REMOVED***
					Min: TLSVersion(tls.VersionTLS12),
					Max: TLSVersion(tls.VersionTLS12),
				***REMOVED***, opts.TLSVersion)
			***REMOVED***)
			t.Run("Blank", func(t *testing.T) ***REMOVED***
				var opts Options
				jsonStr := `***REMOVED***"tlsVersion":""***REMOVED***`
				assert.NoError(t, json.Unmarshal([]byte(jsonStr), &opts))
				assert.Equal(t, &TLSVersions***REMOVED******REMOVED***, opts.TLSVersion)
			***REMOVED***)
		***REMOVED***)
	***REMOVED***)
	t.Run("TLSAuth", func(t *testing.T) ***REMOVED***
		tlsAuth := []*TLSAuth***REMOVED***
			***REMOVED***TLSAuthFields***REMOVED***
				Domains: []string***REMOVED***"example.com", "*.example.com"***REMOVED***,
				Cert: "-----BEGIN CERTIFICATE-----\n" +
					"MIIBoTCCAUegAwIBAgIUQl0J1Gkd6U2NIMwMDnpfH8c1myEwCgYIKoZIzj0EAwIw\n" +
					"EDEOMAwGA1UEAxMFTXkgQ0EwHhcNMTcwODE1MTYxODAwWhcNMTgwODE1MTYxODAw\n" +
					"WjAQMQ4wDAYDVQQDEwV1c2VyMTBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABLaf\n" +
					"xEOmBHkzbqd9/0VZX/39qO2yQq2Gz5faRdvy38kuLMCV+9HYrfMx6GYCZzTUIq6h\n" +
					"8QXOrlgYTixuUVfhJNWjfzB9MA4GA1UdDwEB/wQEAwIFoDAdBgNVHSUEFjAUBggr\n" +
					"BgEFBQcDAQYIKwYBBQUHAwIwDAYDVR0TAQH/BAIwADAdBgNVHQ4EFgQUxmQiq5K3\n" +
					"KUnVME945Byt3Ysvkh8wHwYDVR0jBBgwFoAU3qEhcpRgpsqo9V+LFns9a+oZIYww\n" +
					"CgYIKoZIzj0EAwIDSAAwRQIgSGxnJ+/cLUNTzt7fhr/mjJn7ShsTW33dAdfLM7H2\n" +
					"z/gCIQDyVf8DePtxlkMBScTxZmIlMQdNc6+6VGZQ4QscruVLmg==\n" +
					"-----END CERTIFICATE-----",
				Key: "-----BEGIN EC PRIVATE KEY-----\n" +
					"MHcCAQEEIAfJeoc+XgcqmYV0b4owmofx0LXwPRqOPXMO+PUKxZSgoAoGCCqGSM49\n" +
					"AwEHoUQDQgAEtp/EQ6YEeTNup33/RVlf/f2o7bJCrYbPl9pF2/LfyS4swJX70dit\n" +
					"8zHoZgJnNNQirqHxBc6uWBhOLG5RV+Ek1Q==\n" +
					"-----END EC PRIVATE KEY-----",
			***REMOVED***, nil***REMOVED***,
			***REMOVED***TLSAuthFields***REMOVED***
				Domains: []string***REMOVED***"sub.example.com"***REMOVED***,
				Cert: "-----BEGIN CERTIFICATE-----\n" +
					"MIIBojCCAUegAwIBAgIUWMpVQhmGoLUDd2x6XQYoOOV6C9AwCgYIKoZIzj0EAwIw\n" +
					"EDEOMAwGA1UEAxMFTXkgQ0EwHhcNMTcwODE1MTYxODAwWhcNMTgwODE1MTYxODAw\n" +
					"WjAQMQ4wDAYDVQQDEwV1c2VyMTBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABBfF\n" +
					"85gu8fDbNGNlsrtnO+4HvuiP4IXA041jjGczD5kUQ8aihS7hg81tSrLNd1jgxkkv\n" +
					"Po+3TQjzniysiunG3iKjfzB9MA4GA1UdDwEB/wQEAwIFoDAdBgNVHSUEFjAUBggr\n" +
					"BgEFBQcDAQYIKwYBBQUHAwIwDAYDVR0TAQH/BAIwADAdBgNVHQ4EFgQUU0JfPCQb\n" +
					"2YpQZV4j1yiRXBa7J64wHwYDVR0jBBgwFoAU3qEhcpRgpsqo9V+LFns9a+oZIYww\n" +
					"CgYIKoZIzj0EAwIDSQAwRgIhANYDaM18sXAdkjybHccH8xTbBWUNpOYvoHhrGW32\n" +
					"Ov9JAiEA7QKGpm07tQl8p+t7UsOgZu132dHNZUtfgp1bjWfcapU=\n" +
					"-----END CERTIFICATE-----",
				Key: "-----BEGIN EC PRIVATE KEY-----\n" +
					"MHcCAQEEINVilD5qOBkSy+AYfd41X0QPB5N3Z6OzgoBj8FZmSJOFoAoGCCqGSM49\n" +
					"AwEHoUQDQgAEF8XzmC7x8Ns0Y2Wyu2c77ge+6I/ghcDTjWOMZzMPmRRDxqKFLuGD\n" +
					"zW1Kss13WODGSS8+j7dNCPOeLKyK6cbeIg==\n" +
					"-----END EC PRIVATE KEY-----",
			***REMOVED***, nil***REMOVED***,
		***REMOVED***
		opts := Options***REMOVED******REMOVED***.Apply(Options***REMOVED***TLSAuth: tlsAuth***REMOVED***)
		assert.Equal(t, tlsAuth, opts.TLSAuth)

		t.Run("Roundtrip", func(t *testing.T) ***REMOVED***
			optsData, err := json.Marshal(opts)
			assert.NoError(t, err)

			var opts2 Options
			assert.NoError(t, json.Unmarshal(optsData, &opts2))
			if assert.Len(t, opts2.TLSAuth, len(opts.TLSAuth)) ***REMOVED***
				for i := 0; i < len(opts2.TLSAuth); i++ ***REMOVED***
					assert.Equal(t, opts.TLSAuth[i].TLSAuthFields, opts2.TLSAuth[i].TLSAuthFields)
					cert, err := opts2.TLSAuth[i].Certificate()
					assert.NoError(t, err)
					assert.NotNil(t, cert)
				***REMOVED***
			***REMOVED***
		***REMOVED***)
	***REMOVED***)
	t.Run("NoConnectionReuse", func(t *testing.T) ***REMOVED***
		opts := Options***REMOVED******REMOVED***.Apply(Options***REMOVED***NoConnectionReuse: null.BoolFrom(true)***REMOVED***)
		assert.True(t, opts.NoConnectionReuse.Valid)
		assert.True(t, opts.NoConnectionReuse.Bool)
	***REMOVED***)
	t.Run("Thresholds", func(t *testing.T) ***REMOVED***
		opts := Options***REMOVED******REMOVED***.Apply(Options***REMOVED***Thresholds: map[string]stats.Thresholds***REMOVED***
			"metric": ***REMOVED***
				Thresholds: []*stats.Threshold***REMOVED******REMOVED******REMOVED******REMOVED***,
			***REMOVED***,
		***REMOVED******REMOVED***)
		assert.NotNil(t, opts.Thresholds)
		assert.NotEmpty(t, opts.Thresholds)
	***REMOVED***)
	t.Run("External", func(t *testing.T) ***REMOVED***
		opts := Options***REMOVED******REMOVED***.Apply(Options***REMOVED***External: map[string]interface***REMOVED******REMOVED******REMOVED***"a": 1***REMOVED******REMOVED***)
		assert.Equal(t, map[string]interface***REMOVED******REMOVED******REMOVED***"a": 1***REMOVED***, opts.External)
	***REMOVED***)

	t.Run("JSON", func(t *testing.T) ***REMOVED***
		data, err := json.Marshal(Options***REMOVED******REMOVED***)
		assert.NoError(t, err)
		var opts Options
		assert.NoError(t, json.Unmarshal(data, &opts))
		assert.Equal(t, Options***REMOVED******REMOVED***, opts)
	***REMOVED***)
***REMOVED***

func TestOptionsEnv(t *testing.T) ***REMOVED***
	testdata := map[struct***REMOVED*** Name, Key string ***REMOVED***]map[string]interface***REMOVED******REMOVED******REMOVED***
		***REMOVED***"Paused", "K6_PAUSED"***REMOVED***: ***REMOVED***
			"":      null.Bool***REMOVED******REMOVED***,
			"true":  null.BoolFrom(true),
			"false": null.BoolFrom(false),
		***REMOVED***,
		***REMOVED***"VUs", "K6_VUS"***REMOVED***: ***REMOVED***
			"":    null.Int***REMOVED******REMOVED***,
			"123": null.IntFrom(123),
		***REMOVED***,
		***REMOVED***"VUsMax", "K6_VUS_MAX"***REMOVED***: ***REMOVED***
			"":    null.Int***REMOVED******REMOVED***,
			"123": null.IntFrom(123),
		***REMOVED***,
		***REMOVED***"Duration", "K6_DURATION"***REMOVED***: ***REMOVED***
			"":    NullDuration***REMOVED******REMOVED***,
			"10s": NullDurationFrom(10 * time.Second),
		***REMOVED***,
		***REMOVED***"Iterations", "K6_ITERATIONS"***REMOVED***: ***REMOVED***
			"":    null.Int***REMOVED******REMOVED***,
			"123": null.IntFrom(123),
		***REMOVED***,
		***REMOVED***"Stages", "K6_STAGES"***REMOVED***: ***REMOVED***
			// "": []Stage***REMOVED******REMOVED***,
			"1s": []Stage***REMOVED******REMOVED***
				Duration: NullDurationFrom(1 * time.Second)***REMOVED***,
			***REMOVED***,
			"1s:100": []Stage***REMOVED***
				***REMOVED***Duration: NullDurationFrom(1 * time.Second), Target: null.IntFrom(100)***REMOVED***,
			***REMOVED***,
			"1s,2s:100": []Stage***REMOVED***
				***REMOVED***Duration: NullDurationFrom(1 * time.Second)***REMOVED***,
				***REMOVED***Duration: NullDurationFrom(2 * time.Second), Target: null.IntFrom(100)***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***"MaxRedirects", "K6_MAX_REDIRECTS"***REMOVED***: ***REMOVED***
			"":    null.Int***REMOVED******REMOVED***,
			"123": null.IntFrom(123),
		***REMOVED***,
		***REMOVED***"InsecureSkipTLSVerify", "K6_INSECURE_SKIP_TLS_VERIFY"***REMOVED***: ***REMOVED***
			"":      null.Bool***REMOVED******REMOVED***,
			"true":  null.BoolFrom(true),
			"false": null.BoolFrom(false),
		***REMOVED***,
		// TLSCipherSuites
		// TLSVersion
		// TLSAuth
		***REMOVED***"NoConnectionReuse", "K6_NO_CONNECTION_REUSE"***REMOVED***: ***REMOVED***
			"":      null.Bool***REMOVED******REMOVED***,
			"true":  null.BoolFrom(true),
			"false": null.BoolFrom(false),
		***REMOVED***,
		***REMOVED***"UserAgent", "K6_USER_AGENT"***REMOVED***: ***REMOVED***
			"":    null.String***REMOVED******REMOVED***,
			"Hi!": null.StringFrom("Hi!"),
		***REMOVED***,
		***REMOVED***"Throw", "K6_THROW"***REMOVED***: ***REMOVED***
			"":      null.Bool***REMOVED******REMOVED***,
			"true":  null.BoolFrom(true),
			"false": null.BoolFrom(false),
		***REMOVED***,
		// Thresholds
		// External
	***REMOVED***
	for field, data := range testdata ***REMOVED***
		os.Clearenv()
		t.Run(field.Name, func(t *testing.T) ***REMOVED***
			for str, val := range data ***REMOVED***
				t.Run(`"`+str+`"`, func(t *testing.T) ***REMOVED***
					assert.NoError(t, os.Setenv(field.Key, str))
					var opts Options
					assert.NoError(t, envconfig.Process("k6", &opts))
					assert.Equal(t, val, reflect.ValueOf(opts).FieldByName(field.Name).Interface())
				***REMOVED***)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***
