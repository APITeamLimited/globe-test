package ntlmssp

import (
	"encoding/base64"
	"strings"
)

type authheader []string

func (h authheader) IsBasic() bool ***REMOVED***
	for _, s := range h ***REMOVED***
		if strings.HasPrefix(string(s), "Basic ") ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

func (h authheader) Basic() string ***REMOVED***
	for _, s := range h ***REMOVED***
		if strings.HasPrefix(string(s), "Basic ") ***REMOVED***
			return s
		***REMOVED***
	***REMOVED***
	return ""
***REMOVED***

func (h authheader) IsNegotiate() bool ***REMOVED***
	for _, s := range h ***REMOVED***
		if strings.HasPrefix(string(s), "Negotiate") ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

func (h authheader) IsNTLM() bool ***REMOVED***
	for _, s := range h ***REMOVED***
		if strings.HasPrefix(string(s), "NTLM") ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

func (h authheader) GetData() ([]byte, error) ***REMOVED***
	for _, s := range h ***REMOVED***
		if strings.HasPrefix(string(s), "NTLM") || strings.HasPrefix(string(s), "Negotiate") || strings.HasPrefix(string(s), "Basic ") ***REMOVED***
			p := strings.Split(string(s), " ")
			if len(p) < 2 ***REMOVED***
				return nil, nil
			***REMOVED***
			return base64.StdEncoding.DecodeString(string(p[1]))
		***REMOVED***
	***REMOVED***
	return nil, nil
***REMOVED***

func (h authheader) GetBasicCreds() (username, password string, err error) ***REMOVED***
	d, err := h.GetData()
	if err != nil ***REMOVED***
		return "", "", err
	***REMOVED***
	parts := strings.SplitN(string(d), ":", 2)
	return parts[0], parts[1], nil
***REMOVED***
