package postman

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestUnmarshalDuration(t *testing.T) ***REMOVED***
	var d Duration
	assert.NoError(t, json.Unmarshal([]byte(`12345`), &d))
	assert.Equal(t, Duration(12345*time.Millisecond), d)
***REMOVED***

func TestUnmarshalDurationString(t *testing.T) ***REMOVED***
	var d Duration
	assert.NoError(t, json.Unmarshal([]byte(`"12345"`), &d))
	assert.Equal(t, Duration(12345*time.Millisecond), d)
***REMOVED***

func TestUnmarshalDurationStringUnit(t *testing.T) ***REMOVED***
	var d Duration
	assert.NoError(t, json.Unmarshal([]byte(`"12345s"`), &d))
	assert.Equal(t, Duration(12345*time.Second), d)
***REMOVED***

func TestUnmarshalDurationStringInvalid(t *testing.T) ***REMOVED***
	var d Duration
	assert.Error(t, json.Unmarshal([]byte(`***REMOVED******REMOVED***`), &d))
***REMOVED***

func TestUnmarshalDurationStringInvalidJSON(t *testing.T) ***REMOVED***
	var d Duration
	assert.Error(t, json.Unmarshal([]byte(`e`), &d))
***REMOVED***

func TestUnmarshalTime(t *testing.T) ***REMOVED***
	var tm Time
	assert.NoError(t, json.Unmarshal([]byte(`12345`), &tm))
	assert.Equal(t, int64(12345), time.Time(tm).Unix())
***REMOVED***

func TestUnmarshalTimeString(t *testing.T) ***REMOVED***
	var tm Time
	assert.NoError(t, json.Unmarshal([]byte(`"12345"`), &tm))
	assert.Equal(t, int64(12345), time.Time(tm).Unix())
***REMOVED***

func TestUnmarshalTimeInvalid(t *testing.T) ***REMOVED***
	var tm Time
	assert.Error(t, json.Unmarshal([]byte(`***REMOVED******REMOVED***`), &tm))
***REMOVED***

func TestUnmarshalTimeInvalidJSON(t *testing.T) ***REMOVED***
	var tm Time
	assert.Error(t, json.Unmarshal([]byte(`e`), &tm))
***REMOVED***

func TestUnmarshalInfo(t *testing.T) ***REMOVED***
	j := []byte(`***REMOVED***
		"info": ***REMOVED*** "name": "Name" ***REMOVED***
	***REMOVED***`)

	c := Collection***REMOVED******REMOVED***
	assert.NoError(t, json.Unmarshal(j, &c))
	assert.Equal(t, "Name", c.Info.Name)
***REMOVED***

func TestUnmarshalItem(t *testing.T) ***REMOVED***
	j := []byte(`***REMOVED***
		"item": [
			***REMOVED*** "id": "1", "name": "A" ***REMOVED***,
			***REMOVED*** "id": "2", "name": "B" ***REMOVED***
		]
	***REMOVED***`)

	c := Collection***REMOVED******REMOVED***
	assert.NoError(t, json.Unmarshal(j, &c))
	assert.Len(t, c.Item, 2)
	assert.Equal(t, "1", c.Item[0].ID)
	assert.Equal(t, "A", c.Item[0].Name)
	assert.Equal(t, "2", c.Item[1].ID)
	assert.Equal(t, "B", c.Item[1].Name)
***REMOVED***

func TestUnmarshalItemNested(t *testing.T) ***REMOVED***
	j := []byte(`***REMOVED***
		"item": [***REMOVED***
			"name": "Folder",
			"description": "Lorem ipsum",
			"item": [
				***REMOVED*** "id": "1", "name": "A" ***REMOVED***,
				***REMOVED*** "id": "2", "name": "B" ***REMOVED***
			]
		***REMOVED***]
	***REMOVED***`)

	c := Collection***REMOVED******REMOVED***
	assert.NoError(t, json.Unmarshal(j, &c))
	assert.Len(t, c.Item, 1)
	assert.Equal(t, "Folder", c.Item[0].Name)
	assert.Equal(t, "Lorem ipsum", c.Item[0].Description)
	assert.Len(t, c.Item[0].Item, 2)
	assert.Equal(t, "1", c.Item[0].Item[0].ID)
	assert.Equal(t, "A", c.Item[0].Item[0].Name)
	assert.Equal(t, "2", c.Item[0].Item[1].ID)
	assert.Equal(t, "B", c.Item[0].Item[1].Name)
***REMOVED***

func TestUnmarshalItemRequest(t *testing.T) ***REMOVED***
	j := []byte(`***REMOVED***
		"item": [***REMOVED***
			"request": ***REMOVED***
				"url": "http://example.com/",
				"method": "POST",
				"header": [***REMOVED***"key": "Content-Type", "value": "text/plain"***REMOVED***],
				"body": "lorem ipsum"
			***REMOVED***
		***REMOVED***]
	***REMOVED***`)

	c := Collection***REMOVED******REMOVED***
	assert.NoError(t, json.Unmarshal(j, &c))
	assert.Len(t, c.Item, 1)
	assert.Equal(t, "http://example.com/", c.Item[0].Request.URL)
	assert.Equal(t, "POST", c.Item[0].Request.Method)
	assert.Len(t, c.Item[0].Request.Header, 1)
	assert.Equal(t, "Content-Type", c.Item[0].Request.Header[0].Key)
	assert.Equal(t, "text/plain", c.Item[0].Request.Header[0].Value)
	assert.Equal(t, "lorem ipsum", c.Item[0].Request.Body)
***REMOVED***

func TestUnmarshalItemRequestImplicitMethod(t *testing.T) ***REMOVED***
	j := []byte(`***REMOVED***
		"item": [***REMOVED***
			"request": ***REMOVED***
				"url": "http://example.com/"
			***REMOVED***
		***REMOVED***]
	***REMOVED***`)

	c := Collection***REMOVED******REMOVED***
	assert.NoError(t, json.Unmarshal(j, &c))
	assert.Len(t, c.Item, 1)
	assert.Equal(t, "http://example.com/", c.Item[0].Request.URL)
	assert.Equal(t, "GET", c.Item[0].Request.Method)
***REMOVED***

func TestUnmarshalItemRequestString(t *testing.T) ***REMOVED***
	j := []byte(`***REMOVED***
		"item": [***REMOVED***
			"request": "http://example.com/"
		***REMOVED***]
	***REMOVED***`)

	c := Collection***REMOVED******REMOVED***
	assert.NoError(t, json.Unmarshal(j, &c))
	assert.Len(t, c.Item, 1)
	assert.Equal(t, "http://example.com/", c.Item[0].Request.URL)
	assert.Equal(t, "GET", c.Item[0].Request.Method)
	assert.Len(t, c.Item[0].Request.Header, 0)
	assert.Equal(t, "", c.Item[0].Request.Body)
***REMOVED***

func TestUnmarshalItemRequestMissingHeaderKey(t *testing.T) ***REMOVED***
	j := []byte(`***REMOVED***
		"item": [***REMOVED***
			"request": ***REMOVED***
				"url": "http://example.com/",
				"method": "POST",
				"header": [***REMOVED***"value": "text/plain"***REMOVED***]
			***REMOVED***
		***REMOVED***]
	***REMOVED***`)

	c := Collection***REMOVED******REMOVED***
	assert.Equal(t, ErrMissingHeaderKey, json.Unmarshal(j, &c))
***REMOVED***

func TestUnmarshalItemResponse(t *testing.T) ***REMOVED***
	j := []byte(`***REMOVED***
		"item": [***REMOVED***
			"response": [***REMOVED***
				"originalRequest": "http://example.com/",
				"responseTime": 100,
				"header": [***REMOVED***"key": "Content-Type", "value": "text/plain"***REMOVED***],
				"body": "lorem ipsum",
				"status": "200 OK",
				"code": 200
			***REMOVED***]
		***REMOVED***]
	***REMOVED***`)

	c := Collection***REMOVED******REMOVED***
	assert.NoError(t, json.Unmarshal(j, &c))
	assert.Len(t, c.Item, 1)
	assert.Len(t, c.Item[0].Response, 1)
	assert.Equal(t, "http://example.com/", c.Item[0].Response[0].OriginalRequest.URL)
	assert.Equal(t, "GET", c.Item[0].Response[0].OriginalRequest.Method)
	assert.Equal(t, 100*time.Millisecond, time.Duration(c.Item[0].Response[0].ResponseTime))
	assert.Len(t, c.Item[0].Response[0].Header, 1)
	assert.Equal(t, "Content-Type", c.Item[0].Response[0].Header[0].Key)
	assert.Equal(t, "text/plain", c.Item[0].Response[0].Header[0].Value)
	assert.Equal(t, "lorem ipsum", c.Item[0].Response[0].Body)
	assert.Equal(t, "200 OK", c.Item[0].Response[0].Status)
	assert.Equal(t, 200, c.Item[0].Response[0].Code)
***REMOVED***

func TestUnmarshalItemResponseCookie(t *testing.T) ***REMOVED***
	j := []byte(`***REMOVED***
		"item": [***REMOVED***
			"response": [***REMOVED***
				"originalRequest": "http://example.com/",
				"cookie": [***REMOVED***
					"domain": "example.com",
					"expires": 1,
					"maxAge": 123,
					"hostOnly": true,
					"httpOnly": true,
					"name": "name",
					"path": "/",
					"secure": true,
					"session": true,
					"value": "value"
				***REMOVED***]
			***REMOVED***]
		***REMOVED***]
	***REMOVED***`)

	c := Collection***REMOVED******REMOVED***
	assert.NoError(t, json.Unmarshal(j, &c))
	assert.Len(t, c.Item, 1)
	assert.Len(t, c.Item[0].Response, 1)
	assert.Len(t, c.Item[0].Response[0].Cookie, 1)
	assert.Equal(t, "example.com", c.Item[0].Response[0].Cookie[0].Domain)
	assert.Equal(t, int64(1), time.Time(c.Item[0].Response[0].Cookie[0].Expires).Unix())
	assert.Equal(t, 123*time.Millisecond, time.Duration(c.Item[0].Response[0].Cookie[0].MaxAge))
	assert.Equal(t, true, c.Item[0].Response[0].Cookie[0].HostOnly)
	assert.Equal(t, true, c.Item[0].Response[0].Cookie[0].HTTPOnly)
	assert.Equal(t, "name", c.Item[0].Response[0].Cookie[0].Name)
	assert.Equal(t, "/", c.Item[0].Response[0].Cookie[0].Path)
	assert.Equal(t, true, c.Item[0].Response[0].Cookie[0].Secure)
	assert.Equal(t, true, c.Item[0].Response[0].Cookie[0].Session)
	assert.Equal(t, "value", c.Item[0].Response[0].Cookie[0].Value)
***REMOVED***

func TestUmmarshalEvent(t *testing.T) ***REMOVED***
	j := []byte(`***REMOVED***
		"event": [***REMOVED***
			"listen": "test",
			"script": ***REMOVED***
				"id": "script1",
				"type": "text/javascript",
				"exec": "var v = 1 + 1;\nconsole.log(v);",
				"name": "script1.js"
			***REMOVED***
		***REMOVED***]
	***REMOVED***`)

	c := Collection***REMOVED******REMOVED***
	assert.NoError(t, json.Unmarshal(j, &c))
	assert.Len(t, c.Event, 1)
	assert.Equal(t, "test", c.Event[0].Listen)
	assert.Equal(t, "script1", c.Event[0].Script.ID)
	assert.Equal(t, "text/javascript", c.Event[0].Script.Type)
	assert.Equal(t, ScriptExec("var v = 1 + 1;\nconsole.log(v);"), c.Event[0].Script.Exec)
	assert.Equal(t, "script1.js", c.Event[0].Script.Name)
	assert.Equal(t, false, c.Event[0].Disabled)
***REMOVED***

func TestUmmarshalEventArrayExec(t *testing.T) ***REMOVED***
	j := []byte(`***REMOVED***
		"event": [***REMOVED***
			"listen": "test",
			"script": ***REMOVED***
				"id": "script1",
				"type": "text/javascript",
				"exec": [
					"var v = 1 + 1;",
					"console.log(v);"
				],
				"name": "script1.js"
			***REMOVED***
		***REMOVED***]
	***REMOVED***`)

	c := Collection***REMOVED******REMOVED***
	assert.NoError(t, json.Unmarshal(j, &c))
	assert.Len(t, c.Event, 1)
	assert.Equal(t, "test", c.Event[0].Listen)
	assert.Equal(t, "script1", c.Event[0].Script.ID)
	assert.Equal(t, "text/javascript", c.Event[0].Script.Type)
	assert.Equal(t, ScriptExec("var v = 1 + 1;\nconsole.log(v);"), c.Event[0].Script.Exec)
	assert.Equal(t, "script1.js", c.Event[0].Script.Name)
	assert.Equal(t, false, c.Event[0].Disabled)
***REMOVED***

func TestUmmarshalEventScriptImplicitType(t *testing.T) ***REMOVED***
	j := []byte(`***REMOVED***
		"event": [***REMOVED***
			"listen": "test",
			"script": ***REMOVED***
				"exec": "var v = 1 + 1;\nconsole.log(v);"
			***REMOVED***
		***REMOVED***]
	***REMOVED***`)

	c := Collection***REMOVED******REMOVED***
	assert.NoError(t, json.Unmarshal(j, &c))
	assert.Len(t, c.Event, 1)
	assert.Equal(t, "test", c.Event[0].Listen)
	assert.Equal(t, "text/javascript", c.Event[0].Script.Type)
	assert.Equal(t, ScriptExec("var v = 1 + 1;\nconsole.log(v);"), c.Event[0].Script.Exec)
***REMOVED***

func TestUmmarshalEventWrongScriptType(t *testing.T) ***REMOVED***
	j := []byte(`***REMOVED***
		"event": [***REMOVED***
			"listen": "test",
			"script": ***REMOVED***
				"type": "text/vbscript",
				"exec": "/* I don't actually know VBScript lol */"
			***REMOVED***
		***REMOVED***]
	***REMOVED***`)

	c := Collection***REMOVED******REMOVED***
	assert.Equal(t, ErrScriptUnsupportedType, json.Unmarshal(j, &c))
***REMOVED***

func TestUmmarshalEventScriptString(t *testing.T) ***REMOVED***
	j := []byte(`***REMOVED***
		"event": [***REMOVED***
			"listen": "test",
			"script": "var v = 1 + 1;\nconsole.log(v);"
		***REMOVED***]
	***REMOVED***`)

	c := Collection***REMOVED******REMOVED***
	assert.NoError(t, json.Unmarshal(j, &c))
	assert.Len(t, c.Event, 1)
	assert.Equal(t, "test", c.Event[0].Listen)
	assert.Equal(t, "text/javascript", c.Event[0].Script.Type)
	assert.Equal(t, ScriptExec("var v = 1 + 1;\nconsole.log(v);"), c.Event[0].Script.Exec)
***REMOVED***

func TestUmmarshalEventScriptInvalid(t *testing.T) ***REMOVED***
	j := []byte(`***REMOVED***
		"event": [***REMOVED***
			"listen": "test",
			"script": 12345
		***REMOVED***]
	***REMOVED***`)

	c := Collection***REMOVED******REMOVED***
	assert.Error(t, json.Unmarshal(j, &c))
***REMOVED***

func TestUnmarshalVariableNotImplemented(t *testing.T) ***REMOVED***
	j := []byte(`***REMOVED*** "variable": [***REMOVED******REMOVED***] ***REMOVED***`)
	c := Collection***REMOVED******REMOVED***
	assert.Equal(t, ErrVariablesNotSupported, json.Unmarshal(j, &c))
***REMOVED***
