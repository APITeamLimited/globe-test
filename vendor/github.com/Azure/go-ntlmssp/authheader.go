package ntlmssp

import (
	"encoding/base64"
	"strings"
)

type authheader string

func (h authheader) IsBasic() bool ***REMOVED***
	return strings.HasPrefix(string(h), "Basic ")
***REMOVED***

func (h authheader) IsNegotiate() bool ***REMOVED***
	return strings.HasPrefix(string(h), "Negotiate")
***REMOVED***

func (h authheader) IsNTLM() bool ***REMOVED***
	return strings.HasPrefix(string(h), "NTLM")
***REMOVED***

func (h authheader) GetData() ([]byte, error) ***REMOVED***
	p := strings.Split(string(h), " ")
	if len(p) < 2 ***REMOVED***
		return nil, nil
	***REMOVED***
	return base64.StdEncoding.DecodeString(string(p[1]))
***REMOVED***

func (h authheader) GetBasicCreds() (username, password string, err error) ***REMOVED***
	d, err := h.GetData()
	if err != nil ***REMOVED***
		return "", "", err
	***REMOVED***
	parts := strings.SplitN(string(d), ":", 2)
	return parts[0], parts[1], nil
***REMOVED***
