package sarama

import (
	"net"
	"strconv"
)

type ConsumerMetadataResponse struct ***REMOVED***
	Err             KError
	Coordinator     *Broker
	CoordinatorID   int32  // deprecated: use Coordinator.ID()
	CoordinatorHost string // deprecated: use Coordinator.Addr()
	CoordinatorPort int32  // deprecated: use Coordinator.Addr()
***REMOVED***

func (r *ConsumerMetadataResponse) decode(pd packetDecoder, version int16) (err error) ***REMOVED***
	tmp, err := pd.getInt16()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	r.Err = KError(tmp)

	coordinator := new(Broker)
	if err := coordinator.decode(pd); err != nil ***REMOVED***
		return err
	***REMOVED***
	if coordinator.addr == ":0" ***REMOVED***
		return nil
	***REMOVED***
	r.Coordinator = coordinator

	// this can all go away in 2.0, but we have to fill in deprecated fields to maintain
	// backwards compatibility
	host, portstr, err := net.SplitHostPort(r.Coordinator.Addr())
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	port, err := strconv.ParseInt(portstr, 10, 32)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	r.CoordinatorID = r.Coordinator.ID()
	r.CoordinatorHost = host
	r.CoordinatorPort = int32(port)

	return nil
***REMOVED***

func (r *ConsumerMetadataResponse) encode(pe packetEncoder) error ***REMOVED***
	pe.putInt16(int16(r.Err))
	if r.Coordinator != nil ***REMOVED***
		host, portstr, err := net.SplitHostPort(r.Coordinator.Addr())
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		port, err := strconv.ParseInt(portstr, 10, 32)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		pe.putInt32(r.Coordinator.ID())
		if err := pe.putString(host); err != nil ***REMOVED***
			return err
		***REMOVED***
		pe.putInt32(int32(port))
		return nil
	***REMOVED***
	pe.putInt32(r.CoordinatorID)
	if err := pe.putString(r.CoordinatorHost); err != nil ***REMOVED***
		return err
	***REMOVED***
	pe.putInt32(r.CoordinatorPort)
	return nil
***REMOVED***

func (r *ConsumerMetadataResponse) key() int16 ***REMOVED***
	return 10
***REMOVED***

func (r *ConsumerMetadataResponse) version() int16 ***REMOVED***
	return 0
***REMOVED***

func (r *ConsumerMetadataResponse) requiredVersion() KafkaVersion ***REMOVED***
	return V0_8_2_0
***REMOVED***
