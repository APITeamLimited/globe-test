package sarama

import (
	"fmt"
	"time"
)

type Timestamp struct ***REMOVED***
	*time.Time
***REMOVED***

func (t Timestamp) encode(pe packetEncoder) error ***REMOVED***
	timestamp := int64(-1)

	if !t.Before(time.Unix(0, 0)) ***REMOVED***
		timestamp = t.UnixNano() / int64(time.Millisecond)
	***REMOVED*** else if !t.IsZero() ***REMOVED***
		return PacketEncodingError***REMOVED***fmt.Sprintf("invalid timestamp (%v)", t)***REMOVED***
	***REMOVED***

	pe.putInt64(timestamp)
	return nil
***REMOVED***

func (t Timestamp) decode(pd packetDecoder) error ***REMOVED***
	millis, err := pd.getInt64()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// negative timestamps are invalid, in these cases we should return
	// a zero time
	timestamp := time.Time***REMOVED******REMOVED***
	if millis >= 0 ***REMOVED***
		timestamp = time.Unix(millis/1000, (millis%1000)*int64(time.Millisecond))
	***REMOVED***

	*t.Time = timestamp
	return nil
***REMOVED***
