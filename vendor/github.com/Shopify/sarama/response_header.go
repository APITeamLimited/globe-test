package sarama

import "fmt"

type responseHeader struct ***REMOVED***
	length        int32
	correlationID int32
***REMOVED***

func (r *responseHeader) decode(pd packetDecoder) (err error) ***REMOVED***
	r.length, err = pd.getInt32()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if r.length <= 4 || r.length > MaxResponseSize ***REMOVED***
		return PacketDecodingError***REMOVED***fmt.Sprintf("message of length %d too large or too small", r.length)***REMOVED***
	***REMOVED***

	r.correlationID, err = pd.getInt32()
	return err
***REMOVED***
