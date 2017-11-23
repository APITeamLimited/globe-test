// Copyright 2017 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package render

import (
	"net/http"

	"github.com/ugorji/go/codec"
)

type MsgPack struct ***REMOVED***
	Data interface***REMOVED******REMOVED***
***REMOVED***

var msgpackContentType = []string***REMOVED***"application/msgpack; charset=utf-8"***REMOVED***

func (r MsgPack) WriteContentType(w http.ResponseWriter) ***REMOVED***
	writeContentType(w, msgpackContentType)
***REMOVED***

func (r MsgPack) Render(w http.ResponseWriter) error ***REMOVED***
	return WriteMsgPack(w, r.Data)
***REMOVED***

func WriteMsgPack(w http.ResponseWriter, obj interface***REMOVED******REMOVED***) error ***REMOVED***
	writeContentType(w, msgpackContentType)
	var h codec.Handle = new(codec.MsgpackHandle)
	return codec.NewEncoder(w, h).Encode(obj)
***REMOVED***
