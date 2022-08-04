package influxdb

import (
	"fmt"
	"strings"

	client "github.com/influxdata/influxdb1-client/v2"
	"gopkg.in/guregu/null.v3"
)

func MakeClient(conf Config) (client.Client, error) ***REMOVED***
	if strings.HasPrefix(conf.Addr.String, "udp://") ***REMOVED***
		return client.NewUDPClient(client.UDPConfig***REMOVED***
			Addr:        strings.TrimPrefix(conf.Addr.String, "udp://"),
			PayloadSize: int(conf.PayloadSize.Int64),
		***REMOVED***)
	***REMOVED***
	if conf.Addr.String == "" ***REMOVED***
		conf.Addr = null.StringFrom("http://localhost:8086")
	***REMOVED***
	return client.NewHTTPClient(client.HTTPConfig***REMOVED***
		Addr:               conf.Addr.String,
		Username:           conf.Username.String,
		Password:           conf.Password.String,
		UserAgent:          "k6",
		InsecureSkipVerify: conf.Insecure.Bool,
	***REMOVED***)
***REMOVED***

func MakeBatchConfig(conf Config) client.BatchPointsConfig ***REMOVED***
	if !conf.DB.Valid || conf.DB.String == "" ***REMOVED***
		conf.DB = null.StringFrom("k6")
	***REMOVED***
	return client.BatchPointsConfig***REMOVED***
		Precision:        conf.Precision.String,
		Database:         conf.DB.String,
		RetentionPolicy:  conf.Retention.String,
		WriteConsistency: conf.Consistency.String,
	***REMOVED***
***REMOVED***

func checkDuplicatedTypeDefinitions(fieldKinds map[string]FieldKind, tag string) error ***REMOVED***
	if _, found := fieldKinds[tag]; found ***REMOVED***
		return fmt.Errorf("a tag name (%s) shows up more than once in InfluxDB field type configurations", tag)
	***REMOVED***
	return nil
***REMOVED***

// MakeFieldKinds reads the Config and returns a lookup map of tag names to
// the field type their values should be converted to.
func MakeFieldKinds(conf Config) (map[string]FieldKind, error) ***REMOVED***
	fieldKinds := make(map[string]FieldKind)
	for _, tag := range conf.TagsAsFields ***REMOVED***
		var fieldName, fieldType string
		s := strings.SplitN(tag, ":", 2)
		if len(s) == 1 ***REMOVED***
			fieldName, fieldType = s[0], "string"
		***REMOVED*** else ***REMOVED***
			fieldName, fieldType = s[0], s[1]
		***REMOVED***

		err := checkDuplicatedTypeDefinitions(fieldKinds, fieldName)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		switch fieldType ***REMOVED***
		case "string":
			fieldKinds[fieldName] = String
		case "bool":
			fieldKinds[fieldName] = Bool
		case "float":
			fieldKinds[fieldName] = Float
		case "int":
			fieldKinds[fieldName] = Int
		default:
			return nil, fmt.Errorf("an invalid type (%s) is specified for an InfluxDB field (%s)",
				fieldType, fieldName)
		***REMOVED***
	***REMOVED***

	return fieldKinds, nil
***REMOVED***
