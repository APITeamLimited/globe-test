// Copyright 2017 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import (
	"net/http"

	"github.com/ugorji/go/codec"
)

type msgpackBinding struct***REMOVED******REMOVED***

func (msgpackBinding) Name() string ***REMOVED***
	return "msgpack"
***REMOVED***

func (msgpackBinding) Bind(req *http.Request, obj interface***REMOVED******REMOVED***) error ***REMOVED***

	if err := codec.NewDecoder(req.Body, new(codec.MsgpackHandle)).Decode(&obj); err != nil ***REMOVED***
		//var decoder *codec.Decoder = codec.NewDecoder(req.Body, &codec.MsgpackHandle)
		//if err := decoder.Decode(&obj); err != nil ***REMOVED***
		return err
	***REMOVED***
	return validate(obj)

***REMOVED***
