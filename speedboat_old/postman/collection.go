package postman

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var (
	ErrVariablesNotSupported = errors.New("Variables are not yet implemented")
	ErrScriptSrcNotSupported = errors.New("External scripts are not implemented")

	ErrScriptUnsupportedType = errors.New("Only text/javascript scripts are supported")
	ErrDurationWrongType     = errors.New("Durations must be numbers or strings")
	ErrTimeWrongType         = errors.New("Times must be numbers or strings")
	ErrMissingHeaderKey      = errors.New("Missing key in request header")
)

type Duration time.Duration

func (d *Duration) UnmarshalJSON(b []byte) error ***REMOVED***
	var data interface***REMOVED******REMOVED***
	if err := json.Unmarshal(b, &data); err != nil ***REMOVED***
		return err
	***REMOVED***

	switch v := data.(type) ***REMOVED***
	case string:
		num, err := strconv.ParseInt(v, 10, 64)
		if err != nil ***REMOVED***
			duration, err := time.ParseDuration(v)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			*d = Duration(duration)
			break
		***REMOVED***
		*d = Duration(time.Duration(num) * time.Millisecond)
	case float64:
		*d = Duration(time.Duration(v) * time.Millisecond)
	default:
		return ErrDurationWrongType
	***REMOVED***

	return nil
***REMOVED***

type Time time.Time

func (d *Time) UnmarshalJSON(b []byte) error ***REMOVED***
	var data interface***REMOVED******REMOVED***
	if err := json.Unmarshal(b, &data); err != nil ***REMOVED***
		return err
	***REMOVED***

	switch v := data.(type) ***REMOVED***
	case string:
		// Why.
		if v == "Invalid Date" ***REMOVED***
			*d = Time***REMOVED******REMOVED***
			break
		***REMOVED***

		t, err := time.Parse("Mon Jan 2 2006 15:04:05 GMT-0700 (MST)", v)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		*d = Time(t)
	case float64:
		*d = Time(time.Unix(int64(v), 0))
	default:
		return ErrTimeWrongType
	***REMOVED***

	return nil
***REMOVED***

type ScriptSrc struct***REMOVED******REMOVED***

func (ScriptSrc) UnmarshalJSON(b []byte) error ***REMOVED***
	return ErrScriptSrcNotSupported
***REMOVED***

type ScriptExec string

func (e *ScriptExec) UnmarshalJSON(b []byte) error ***REMOVED***
	var data interface***REMOVED******REMOVED***
	if err := json.Unmarshal(b, &data); err != nil ***REMOVED***
		return err
	***REMOVED***

	switch v := data.(type) ***REMOVED***
	case string:
		*e = ScriptExec(v)
	case []interface***REMOVED******REMOVED***:
		lines := make([]string, 0, len(v))
		for _, val := range v ***REMOVED***
			switch line := val.(type) ***REMOVED***
			case string:
				lines = append(lines, line)
			default:
				lines = append(lines, fmt.Sprint(line))
			***REMOVED***
		***REMOVED***
		*e = ScriptExec(strings.Join(lines, "\n"))
	***REMOVED***

	return nil
***REMOVED***

type ScriptImpl struct ***REMOVED***
	ID   string     `json"id"`
	Type string     `json:"type"`
	Exec ScriptExec `json:"exec"`
	Src  ScriptSrc  `json:"src"`
	Name string     `json:"name"`
***REMOVED***

type Script ScriptImpl

func (s *Script) UnmarshalJSON(b []byte) error ***REMOVED***
	var data interface***REMOVED******REMOVED***
	if err := json.Unmarshal(b, &data); err != nil ***REMOVED***
		return err
	***REMOVED***

	switch v := data.(type) ***REMOVED***
	case string:
		s.Type = "text/javascript"
		s.Exec = ScriptExec(v)
		s.Name = "inline"
	default:
		var impl ScriptImpl
		if err := json.Unmarshal(b, &impl); err != nil ***REMOVED***
			return err
		***REMOVED***

		switch impl.Type ***REMOVED***
		case "text/javascript":
		case "":
			impl.Type = "text/javascript"
		default:
			return ErrScriptUnsupportedType
		***REMOVED***

		*s = Script(impl)
	***REMOVED***

	return nil
***REMOVED***

type Event struct ***REMOVED***
	Listen   string `json:"listen"`
	Script   Script `json:"script"`
	Disabled bool   `json:"disabled"`
***REMOVED***

type HeaderImpl struct ***REMOVED***
	Key   string `json:"key"`
	Value string `json:"value"`
***REMOVED***

type Header HeaderImpl

func (h *Header) UnmarshalJSON(b []byte) error ***REMOVED***
	var impl HeaderImpl
	if err := json.Unmarshal(b, &impl); err != nil ***REMOVED***
		return err
	***REMOVED***

	if impl.Key == "" ***REMOVED***
		return ErrMissingHeaderKey
	***REMOVED***

	*h = Header(impl)
	return nil
***REMOVED***

type Cookie struct ***REMOVED***
	Domain   string   `json:"domain"`
	Expires  Time     `json:"expires"`
	MaxAge   Duration `json:"maxAge"`
	HostOnly bool     `json:"hostOnly"`
	HTTPOnly bool     `json:"httpOnly"`
	Name     string   `json:"name"`
	Path     string   `json:"path"`
	Secure   bool     `json:"secure"`
	Session  bool     `json:"session"`
	Value    string   `json:"value`
	// Not parsing extensions. They are wholly uninteresting to us.
***REMOVED***

type Param struct ***REMOVED***
	Key     string `json:"key"`
	Value   string `json:"value"`
	Enabled bool   `json:"enabled"`
***REMOVED***

type RequestImpl struct ***REMOVED***
	URL    string   `json:"url"` // TODO: Decompose into net/url.URL structs, handle maps
	Auth   Auth     `json:"auth"`
	Method string   `json:"method"`
	Header []Header `json:"header"` // Docs aren't clear on what a string here means?
	Body   struct ***REMOVED***
		Mode       string  `json:"mode"`
		Raw        string  `json:"raw"`
		URLEncoded []Param `json:"urlencoded"`
		FormData   []Param `json:"formdata"`
	***REMOVED*** `json:"body"`
***REMOVED***

type Request RequestImpl

func (r *Request) UnmarshalJSON(b []byte) error ***REMOVED***
	var data interface***REMOVED******REMOVED***
	if err := json.Unmarshal(b, &data); err != nil ***REMOVED***
		return err
	***REMOVED***

	switch v := data.(type) ***REMOVED***
	case string:
		r.URL = v
		r.Method = "GET"
	default:
		var impl RequestImpl
		if err := json.Unmarshal(b, &impl); err != nil ***REMOVED***
			return err
		***REMOVED***
		if impl.Method == "" ***REMOVED***
			impl.Method = "GET"
		***REMOVED***
		*r = Request(impl)
	***REMOVED***

	return nil
***REMOVED***

type Response struct ***REMOVED***
	OriginalRequest Request  `json:"originalRequest"`
	ResponseTime    Duration `json:"responseTime"`
	Header          []Header `json:"header"`
	Cookie          []Cookie `json:"cookie"`
	Body            string   `json:"body"`
	Status          string   `json:"status"`
	Code            int      `json:"code"`
***REMOVED***

// The docs for this are vague and I can't find a UI for it anywhere in the Postman app.
type Variable struct ***REMOVED***
***REMOVED***

func (Variable) UnmarshalJSON(b []byte) error ***REMOVED***
	return ErrVariablesNotSupported
***REMOVED***

type Auth struct ***REMOVED***
	Type string `json:"type"`

	AWSv4 struct ***REMOVED***
		AccessKey string `json:"accessKey"`
		SecretKey string `json:"secretKey"`
		Region    string `json:"region"`
		Service   string `json:"service"`
	***REMOVED*** `json:"awsv4"`

	Basic struct ***REMOVED***
		Username string `json:"username"`
		Password string `json:"password"`
	***REMOVED*** `json:"basic"`

	Digest struct ***REMOVED***
		Username    string `json:"username"`
		Realm       string `json:"realm"`
		Password    string `json:"password"`
		Nonce       string `json:"nonce"`
		NonceCount  string `json:"nonceCount"`
		Algorithm   string `json:"algorithm"`
		QOP         string `json:"qop"`
		ClientNonce string `json:"clientNonce"`
	***REMOVED*** `json:"digest"`

	Hawk struct ***REMOVED***
		AuthID     string `json:"authId"`
		AuthKey    string `json:"authKey"`
		Algorithm  string `json:"algorithm"`
		User       string `json:"user"`
		Nonce      string `json:"nonce"`
		ExtraData  string `json:"extraData"`
		AppID      string `json:"appId"`
		Delegation string `json:"delegation"`
	***REMOVED*** `json:"hawk"`

	OAuth1 struct ***REMOVED***
		ConsumerKey     string `json:"consumerKey"`
		ConsumerSecret  string `json:"consumerSecret"`
		Token           string `json:"token"`
		TokenSecret     string `json:"tokenSecret"`
		SignatureMethod string `json:"signatureMethod"`
		Timestamp       string `json:"timeStamp"`
		Nonce           string `json:"nonce"`
		Version         string `json:"version"`
		Realm           string `json:"realm"`
		EncodeOAuthSign string `json:"encodeOAuthSign"`
	***REMOVED*** `json:"oauth1"`

	OAuth2 struct ***REMOVED***
		AddTokenTo     string `json:"addTokenTo"`
		CallbackURL    string `json:"callBackUrl"`
		AuthURL        string `json:"authUrl"`
		AccessTokenURL string `json:"accessTokenUrl"`
		ClientID       string `json:"clientId"`
		ClientSecret   string `json:"clientSecret"`
		Scope          string `json:"scope"`

		RequestAccessTokenLocally string `json:"requestAccessTokenLocally"`
	***REMOVED*** `json:"oauth2"`
***REMOVED***

type Item struct ***REMOVED***
	// Items + Folders
	Name string `json:"name"`

	// Items
	ID       string     `json:"id"`
	Event    []Event    `json:"event"`
	Request  Request    `json:"request"`
	Response []Response `json:"response"`

	// Folders
	Description string `json:"description"`
	Item        []Item `json:"item"`
	Auth        Auth   `json:"auth"`
***REMOVED***

type Information struct ***REMOVED***
	Name string `json:"name"`
***REMOVED***

type Collection struct ***REMOVED***
	Info     Information `json:"info"`
	Item     []Item      `json:"item"`
	Event    []Event     `json:"event"`
	Variable []Variable  `json:"variable"`
	Auth     Auth        `json:"auth"`
***REMOVED***
