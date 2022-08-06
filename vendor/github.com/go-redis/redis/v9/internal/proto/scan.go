package proto

import (
	"encoding"
	"fmt"
	"net"
	"reflect"
	"time"

	"github.com/go-redis/redis/v9/internal/util"
)

// Scan parses bytes `b` to `v` with appropriate type.
//nolint:gocyclo
func Scan(b []byte, v interface***REMOVED******REMOVED***) error ***REMOVED***
	switch v := v.(type) ***REMOVED***
	case nil:
		return fmt.Errorf("redis: Scan(nil)")
	case *string:
		*v = util.BytesToString(b)
		return nil
	case *[]byte:
		*v = b
		return nil
	case *int:
		var err error
		*v, err = util.Atoi(b)
		return err
	case *int8:
		n, err := util.ParseInt(b, 10, 8)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		*v = int8(n)
		return nil
	case *int16:
		n, err := util.ParseInt(b, 10, 16)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		*v = int16(n)
		return nil
	case *int32:
		n, err := util.ParseInt(b, 10, 32)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		*v = int32(n)
		return nil
	case *int64:
		n, err := util.ParseInt(b, 10, 64)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		*v = n
		return nil
	case *uint:
		n, err := util.ParseUint(b, 10, 64)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		*v = uint(n)
		return nil
	case *uint8:
		n, err := util.ParseUint(b, 10, 8)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		*v = uint8(n)
		return nil
	case *uint16:
		n, err := util.ParseUint(b, 10, 16)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		*v = uint16(n)
		return nil
	case *uint32:
		n, err := util.ParseUint(b, 10, 32)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		*v = uint32(n)
		return nil
	case *uint64:
		n, err := util.ParseUint(b, 10, 64)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		*v = n
		return nil
	case *float32:
		n, err := util.ParseFloat(b, 32)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		*v = float32(n)
		return err
	case *float64:
		var err error
		*v, err = util.ParseFloat(b, 64)
		return err
	case *bool:
		*v = len(b) == 1 && b[0] == '1'
		return nil
	case *time.Time:
		var err error
		*v, err = time.Parse(time.RFC3339Nano, util.BytesToString(b))
		return err
	case *time.Duration:
		n, err := util.ParseInt(b, 10, 64)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		*v = time.Duration(n)
		return nil
	case encoding.BinaryUnmarshaler:
		return v.UnmarshalBinary(b)
	case *net.IP:
		*v = b
		return nil
	default:
		return fmt.Errorf(
			"redis: can't unmarshal %T (consider implementing BinaryUnmarshaler)", v)
	***REMOVED***
***REMOVED***

func ScanSlice(data []string, slice interface***REMOVED******REMOVED***) error ***REMOVED***
	v := reflect.ValueOf(slice)
	if !v.IsValid() ***REMOVED***
		return fmt.Errorf("redis: ScanSlice(nil)")
	***REMOVED***
	if v.Kind() != reflect.Ptr ***REMOVED***
		return fmt.Errorf("redis: ScanSlice(non-pointer %T)", slice)
	***REMOVED***
	v = v.Elem()
	if v.Kind() != reflect.Slice ***REMOVED***
		return fmt.Errorf("redis: ScanSlice(non-slice %T)", slice)
	***REMOVED***

	next := makeSliceNextElemFunc(v)
	for i, s := range data ***REMOVED***
		elem := next()
		if err := Scan([]byte(s), elem.Addr().Interface()); err != nil ***REMOVED***
			err = fmt.Errorf("redis: ScanSlice index=%d value=%q failed: %w", i, s, err)
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func makeSliceNextElemFunc(v reflect.Value) func() reflect.Value ***REMOVED***
	elemType := v.Type().Elem()

	if elemType.Kind() == reflect.Ptr ***REMOVED***
		elemType = elemType.Elem()
		return func() reflect.Value ***REMOVED***
			if v.Len() < v.Cap() ***REMOVED***
				v.Set(v.Slice(0, v.Len()+1))
				elem := v.Index(v.Len() - 1)
				if elem.IsNil() ***REMOVED***
					elem.Set(reflect.New(elemType))
				***REMOVED***
				return elem.Elem()
			***REMOVED***

			elem := reflect.New(elemType)
			v.Set(reflect.Append(v, elem))
			return elem.Elem()
		***REMOVED***
	***REMOVED***

	zero := reflect.Zero(elemType)
	return func() reflect.Value ***REMOVED***
		if v.Len() < v.Cap() ***REMOVED***
			v.Set(v.Slice(0, v.Len()+1))
			return v.Index(v.Len() - 1)
		***REMOVED***

		v.Set(reflect.Append(v, zero))
		return v.Index(v.Len() - 1)
	***REMOVED***
***REMOVED***
