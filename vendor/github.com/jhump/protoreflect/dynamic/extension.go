package dynamic

import (
	"fmt"

	"github.com/golang/protobuf/proto"

	"github.com/jhump/protoreflect/codec"
	"github.com/jhump/protoreflect/desc"
)

// SetExtension sets the given extension value. If the given message is not a
// dynamic message, the given extension may not be recognized (or may differ
// from the compiled and linked in version of the extension. So in that case,
// this function will serialize the given value to bytes and then use
// proto.SetRawExtension to set the value.
func SetExtension(msg proto.Message, extd *desc.FieldDescriptor, val interface***REMOVED******REMOVED***) error ***REMOVED***
	if !extd.IsExtension() ***REMOVED***
		return fmt.Errorf("given field %s is not an extension", extd.GetFullyQualifiedName())
	***REMOVED***

	if dm, ok := msg.(*Message); ok ***REMOVED***
		return dm.TrySetField(extd, val)
	***REMOVED***

	md, err := desc.LoadMessageDescriptorForMessage(msg)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := checkField(extd, md); err != nil ***REMOVED***
		return err
	***REMOVED***

	val, err = validFieldValue(extd, val)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	var b codec.Buffer
	b.SetDeterministic(defaultDeterminism)
	if err := b.EncodeFieldValue(extd, val); err != nil ***REMOVED***
		return err
	***REMOVED***
	proto.SetRawExtension(msg, extd.GetNumber(), b.Bytes())
	return nil
***REMOVED***
