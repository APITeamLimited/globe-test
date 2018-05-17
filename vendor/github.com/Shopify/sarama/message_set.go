package sarama

type MessageBlock struct ***REMOVED***
	Offset int64
	Msg    *Message
***REMOVED***

// Messages convenience helper which returns either all the
// messages that are wrapped in this block
func (msb *MessageBlock) Messages() []*MessageBlock ***REMOVED***
	if msb.Msg.Set != nil ***REMOVED***
		return msb.Msg.Set.Messages
	***REMOVED***
	return []*MessageBlock***REMOVED***msb***REMOVED***
***REMOVED***

func (msb *MessageBlock) encode(pe packetEncoder) error ***REMOVED***
	pe.putInt64(msb.Offset)
	pe.push(&lengthField***REMOVED******REMOVED***)
	err := msb.Msg.encode(pe)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return pe.pop()
***REMOVED***

func (msb *MessageBlock) decode(pd packetDecoder) (err error) ***REMOVED***
	if msb.Offset, err = pd.getInt64(); err != nil ***REMOVED***
		return err
	***REMOVED***

	if err = pd.push(&lengthField***REMOVED******REMOVED***); err != nil ***REMOVED***
		return err
	***REMOVED***

	msb.Msg = new(Message)
	if err = msb.Msg.decode(pd); err != nil ***REMOVED***
		return err
	***REMOVED***

	if err = pd.pop(); err != nil ***REMOVED***
		return err
	***REMOVED***

	return nil
***REMOVED***

type MessageSet struct ***REMOVED***
	PartialTrailingMessage bool // whether the set on the wire contained an incomplete trailing MessageBlock
	Messages               []*MessageBlock
***REMOVED***

func (ms *MessageSet) encode(pe packetEncoder) error ***REMOVED***
	for i := range ms.Messages ***REMOVED***
		err := ms.Messages[i].encode(pe)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (ms *MessageSet) decode(pd packetDecoder) (err error) ***REMOVED***
	ms.Messages = nil

	for pd.remaining() > 0 ***REMOVED***
		magic, err := magicValue(pd)
		if err != nil ***REMOVED***
			if err == ErrInsufficientData ***REMOVED***
				ms.PartialTrailingMessage = true
				return nil
			***REMOVED***
			return err
		***REMOVED***

		if magic > 1 ***REMOVED***
			return nil
		***REMOVED***

		msb := new(MessageBlock)
		err = msb.decode(pd)
		switch err ***REMOVED***
		case nil:
			ms.Messages = append(ms.Messages, msb)
		case ErrInsufficientData:
			// As an optimization the server is allowed to return a partial message at the
			// end of the message set. Clients should handle this case. So we just ignore such things.
			ms.PartialTrailingMessage = true
			return nil
		default:
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (ms *MessageSet) addMessage(msg *Message) ***REMOVED***
	block := new(MessageBlock)
	block.Msg = msg
	ms.Messages = append(ms.Messages, block)
***REMOVED***
