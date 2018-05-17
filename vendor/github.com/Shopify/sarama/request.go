package sarama

import (
	"encoding/binary"
	"fmt"
	"io"
)

type protocolBody interface ***REMOVED***
	encoder
	versionedDecoder
	key() int16
	version() int16
	requiredVersion() KafkaVersion
***REMOVED***

type request struct ***REMOVED***
	correlationID int32
	clientID      string
	body          protocolBody
***REMOVED***

func (r *request) encode(pe packetEncoder) (err error) ***REMOVED***
	pe.push(&lengthField***REMOVED******REMOVED***)
	pe.putInt16(r.body.key())
	pe.putInt16(r.body.version())
	pe.putInt32(r.correlationID)
	err = pe.putString(r.clientID)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	err = r.body.encode(pe)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return pe.pop()
***REMOVED***

func (r *request) decode(pd packetDecoder) (err error) ***REMOVED***
	var key int16
	if key, err = pd.getInt16(); err != nil ***REMOVED***
		return err
	***REMOVED***
	var version int16
	if version, err = pd.getInt16(); err != nil ***REMOVED***
		return err
	***REMOVED***
	if r.correlationID, err = pd.getInt32(); err != nil ***REMOVED***
		return err
	***REMOVED***
	r.clientID, err = pd.getString()

	r.body = allocateBody(key, version)
	if r.body == nil ***REMOVED***
		return PacketDecodingError***REMOVED***fmt.Sprintf("unknown request key (%d)", key)***REMOVED***
	***REMOVED***
	return r.body.decode(pd, version)
***REMOVED***

func decodeRequest(r io.Reader) (req *request, bytesRead int, err error) ***REMOVED***
	lengthBytes := make([]byte, 4)
	if _, err := io.ReadFull(r, lengthBytes); err != nil ***REMOVED***
		return nil, bytesRead, err
	***REMOVED***
	bytesRead += len(lengthBytes)

	length := int32(binary.BigEndian.Uint32(lengthBytes))
	if length <= 4 || length > MaxRequestSize ***REMOVED***
		return nil, bytesRead, PacketDecodingError***REMOVED***fmt.Sprintf("message of length %d too large or too small", length)***REMOVED***
	***REMOVED***

	encodedReq := make([]byte, length)
	if _, err := io.ReadFull(r, encodedReq); err != nil ***REMOVED***
		return nil, bytesRead, err
	***REMOVED***
	bytesRead += len(encodedReq)

	req = &request***REMOVED******REMOVED***
	if err := decode(encodedReq, req); err != nil ***REMOVED***
		return nil, bytesRead, err
	***REMOVED***
	return req, bytesRead, nil
***REMOVED***

func allocateBody(key, version int16) protocolBody ***REMOVED***
	switch key ***REMOVED***
	case 0:
		return &ProduceRequest***REMOVED******REMOVED***
	case 1:
		return &FetchRequest***REMOVED******REMOVED***
	case 2:
		return &OffsetRequest***REMOVED***Version: version***REMOVED***
	case 3:
		return &MetadataRequest***REMOVED******REMOVED***
	case 8:
		return &OffsetCommitRequest***REMOVED***Version: version***REMOVED***
	case 9:
		return &OffsetFetchRequest***REMOVED******REMOVED***
	case 10:
		return &ConsumerMetadataRequest***REMOVED******REMOVED***
	case 11:
		return &JoinGroupRequest***REMOVED******REMOVED***
	case 12:
		return &HeartbeatRequest***REMOVED******REMOVED***
	case 13:
		return &LeaveGroupRequest***REMOVED******REMOVED***
	case 14:
		return &SyncGroupRequest***REMOVED******REMOVED***
	case 15:
		return &DescribeGroupsRequest***REMOVED******REMOVED***
	case 16:
		return &ListGroupsRequest***REMOVED******REMOVED***
	case 17:
		return &SaslHandshakeRequest***REMOVED******REMOVED***
	case 18:
		return &ApiVersionsRequest***REMOVED******REMOVED***
	case 19:
		return &CreateTopicsRequest***REMOVED******REMOVED***
	case 20:
		return &DeleteTopicsRequest***REMOVED******REMOVED***
	case 22:
		return &InitProducerIDRequest***REMOVED******REMOVED***
	case 24:
		return &AddPartitionsToTxnRequest***REMOVED******REMOVED***
	case 25:
		return &AddOffsetsToTxnRequest***REMOVED******REMOVED***
	case 26:
		return &EndTxnRequest***REMOVED******REMOVED***
	case 28:
		return &TxnOffsetCommitRequest***REMOVED******REMOVED***
	case 29:
		return &DescribeAclsRequest***REMOVED******REMOVED***
	case 30:
		return &CreateAclsRequest***REMOVED******REMOVED***
	case 31:
		return &DeleteAclsRequest***REMOVED******REMOVED***
	case 32:
		return &DescribeConfigsRequest***REMOVED******REMOVED***
	case 33:
		return &AlterConfigsRequest***REMOVED******REMOVED***
	case 37:
		return &CreatePartitionsRequest***REMOVED******REMOVED***
	***REMOVED***
	return nil
***REMOVED***
