package internal

import (
	"fmt"
	"strconv"
	"time"

	"github.com/APITeamLimited/redis/v9/internal/util"
)

func AppendArg(b []byte, v interface***REMOVED******REMOVED***) []byte ***REMOVED***
	switch v := v.(type) ***REMOVED***
	case nil:
		return append(b, "<nil>"...)
	case string:
		return appendUTF8String(b, util.StringToBytes(v))
	case []byte:
		return appendUTF8String(b, v)
	case int:
		return strconv.AppendInt(b, int64(v), 10)
	case int8:
		return strconv.AppendInt(b, int64(v), 10)
	case int16:
		return strconv.AppendInt(b, int64(v), 10)
	case int32:
		return strconv.AppendInt(b, int64(v), 10)
	case int64:
		return strconv.AppendInt(b, v, 10)
	case uint:
		return strconv.AppendUint(b, uint64(v), 10)
	case uint8:
		return strconv.AppendUint(b, uint64(v), 10)
	case uint16:
		return strconv.AppendUint(b, uint64(v), 10)
	case uint32:
		return strconv.AppendUint(b, uint64(v), 10)
	case uint64:
		return strconv.AppendUint(b, v, 10)
	case float32:
		return strconv.AppendFloat(b, float64(v), 'f', -1, 64)
	case float64:
		return strconv.AppendFloat(b, v, 'f', -1, 64)
	case bool:
		if v ***REMOVED***
			return append(b, "true"...)
		***REMOVED***
		return append(b, "false"...)
	case time.Time:
		return v.AppendFormat(b, time.RFC3339Nano)
	default:
		return append(b, fmt.Sprint(v)...)
	***REMOVED***
***REMOVED***

func appendUTF8String(dst []byte, src []byte) []byte ***REMOVED***
	dst = append(dst, src...)
	return dst
***REMOVED***