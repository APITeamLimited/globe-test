// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"encoding/xml"
	"net/http"
	"os"
	"path"
	"reflect"
	"runtime"
	"strings"
)

const BindKey = "_gin-gonic/gin/bindkey"

func Bind(val interface***REMOVED******REMOVED***) HandlerFunc ***REMOVED***
	value := reflect.ValueOf(val)
	if value.Kind() == reflect.Ptr ***REMOVED***
		panic(`Bind struct can not be a pointer. Example:
	Use: gin.Bind(Struct***REMOVED******REMOVED***) instead of gin.Bind(&Struct***REMOVED******REMOVED***)
`)
	***REMOVED***
	typ := value.Type()

	return func(c *Context) ***REMOVED***
		obj := reflect.New(typ).Interface()
		if c.Bind(obj) == nil ***REMOVED***
			c.Set(BindKey, obj)
		***REMOVED***
	***REMOVED***
***REMOVED***

func WrapF(f http.HandlerFunc) HandlerFunc ***REMOVED***
	return func(c *Context) ***REMOVED***
		f(c.Writer, c.Request)
	***REMOVED***
***REMOVED***

func WrapH(h http.Handler) HandlerFunc ***REMOVED***
	return func(c *Context) ***REMOVED***
		h.ServeHTTP(c.Writer, c.Request)
	***REMOVED***
***REMOVED***

type H map[string]interface***REMOVED******REMOVED***

// MarshalXML allows type H to be used with xml.Marshal
func (h H) MarshalXML(e *xml.Encoder, start xml.StartElement) error ***REMOVED***
	start.Name = xml.Name***REMOVED***
		Space: "",
		Local: "map",
	***REMOVED***
	if err := e.EncodeToken(start); err != nil ***REMOVED***
		return err
	***REMOVED***
	for key, value := range h ***REMOVED***
		elem := xml.StartElement***REMOVED***
			Name: xml.Name***REMOVED***Space: "", Local: key***REMOVED***,
			Attr: []xml.Attr***REMOVED******REMOVED***,
		***REMOVED***
		if err := e.EncodeElement(value, elem); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	if err := e.EncodeToken(xml.EndElement***REMOVED***Name: start.Name***REMOVED***); err != nil ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

func assert1(guard bool, text string) ***REMOVED***
	if !guard ***REMOVED***
		panic(text)
	***REMOVED***
***REMOVED***

func filterFlags(content string) string ***REMOVED***
	for i, char := range content ***REMOVED***
		if char == ' ' || char == ';' ***REMOVED***
			return content[:i]
		***REMOVED***
	***REMOVED***
	return content
***REMOVED***

func chooseData(custom, wildcard interface***REMOVED******REMOVED***) interface***REMOVED******REMOVED*** ***REMOVED***
	if custom == nil ***REMOVED***
		if wildcard == nil ***REMOVED***
			panic("negotiation config is invalid")
		***REMOVED***
		return wildcard
	***REMOVED***
	return custom
***REMOVED***

func parseAccept(acceptHeader string) []string ***REMOVED***
	parts := strings.Split(acceptHeader, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts ***REMOVED***
		index := strings.IndexByte(part, ';')
		if index >= 0 ***REMOVED***
			part = part[0:index]
		***REMOVED***
		part = strings.TrimSpace(part)
		if len(part) > 0 ***REMOVED***
			out = append(out, part)
		***REMOVED***
	***REMOVED***
	return out
***REMOVED***

func lastChar(str string) uint8 ***REMOVED***
	size := len(str)
	if size == 0 ***REMOVED***
		panic("The length of the string can't be 0")
	***REMOVED***
	return str[size-1]
***REMOVED***

func nameOfFunction(f interface***REMOVED******REMOVED***) string ***REMOVED***
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
***REMOVED***

func joinPaths(absolutePath, relativePath string) string ***REMOVED***
	if len(relativePath) == 0 ***REMOVED***
		return absolutePath
	***REMOVED***

	finalPath := path.Join(absolutePath, relativePath)
	appendSlash := lastChar(relativePath) == '/' && lastChar(finalPath) != '/'
	if appendSlash ***REMOVED***
		return finalPath + "/"
	***REMOVED***
	return finalPath
***REMOVED***

func resolveAddress(addr []string) string ***REMOVED***
	switch len(addr) ***REMOVED***
	case 0:
		if port := os.Getenv("PORT"); len(port) > 0 ***REMOVED***
			debugPrint("Environment variable PORT=\"%s\"", port)
			return ":" + port
		***REMOVED***
		debugPrint("Environment variable PORT is undefined. Using port :8080 by default")
		return ":8080"
	case 1:
		return addr[0]
	default:
		panic("too much parameters")
	***REMOVED***
***REMOVED***
