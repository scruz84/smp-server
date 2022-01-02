package smpserver

type FixedLengthFragmentEncoder struct {
	nextHandler  ChannelWriteHandler
	fragmentSize int
}

func NewFixedLengthFragmentEncoder(nextHandler ChannelWriteHandler, fragmentSize int) *FixedLengthFragmentEncoder {
	h := FixedLengthFragmentEncoder{nextHandler: nextHandler, fragmentSize: fragmentSize}
	return &h
}

func (h FixedLengthFragmentEncoder) Write(message []byte, channel Channel) {
	encodedMessage := message

	//calculate the bytes to set to zero to fill the fragment size
	var zeros int
	if len(message) <= h.fragmentSize {
		zeros = h.fragmentSize - len(message)
	} else {
		zeros = h.fragmentSize - len(message)%h.fragmentSize
	}

	if zeros > 0 {
		encodedMessage = append(encodedMessage, make([]byte, zeros)...)
	}

	if h.nextHandler != nil {
		h.nextHandler.Write(encodedMessage, channel)
	} else {
		channel.write(encodedMessage)
	}
}

func (h FixedLengthFragmentEncoder) NextHandler(handler ChannelWriteHandler) {
	h.nextHandler = handler
}
