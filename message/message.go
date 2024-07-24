package message

import (
	"encoding/binary"
	"fmt"
	"io"
)

type messageID uint8

const (
	MsgChoke messageID = 0
	
	MsgUnchoke messageID = 1

	MsgIntrested messageID = 2

	MsgNotIntrested messageID = 3

	MsgHave messageID = 4

	MsgBitfield messageID = 5

	MsgRequest messageID = 6

	MsgPiece messageID = 7

	MsgCancel messageID = 8
)

type Message struct {
	ID messageID
	Payload []byte
}

func FormatRequest(index, begin, length int) *Message{
	payload := make([]byte, 12)
	binary.BigEndian.PutUint32(payload[0:4], uint32(index))
	binary.BigEndian.PutUint32(payload[4:8], uint32(begin))
	binary.BigEndian.PutUint32(payload[8:12], uint32(length))
	return &Message{ID: MsgHave, Payload: payload}
}

func FormatHave(index int) *Message {
	payload := make([]byte, 4)
	binary.BigEndian.PutUint32(payload, uint32(index))
	return &Message{ID: MsgHave, Payload: payload}
}

func Read(r io.Reader) (*Message, error){
	lengthBuf := make([]byte, 4)
	_, err := io.ReadFull(r, lengthBuf)
	if err != nil{
		return nil, err
	}
	length := binary.BigEndian.Uint32(lengthBuf)

	if length == 0{
		return nil, nil
	}

	messageBuf := make([]byte, length)
	_, err = io.ReadFull(r, messageBuf)
	if err != nil{
		return nil, err
	}

	m := Message{
		ID: messageID(messageBuf[0]),
		Payload: messageBuf[1:],
	}

	return &m, nil
}

func (m *Message) name() string{
	if m == nil {
		return "KeepAlive"
	}
	switch m.ID{
		case MsgChoke:
			return "Choke"
		case MsgUnchoke:
			return "Unchoke"
		case MsgIntrested:
			return "Intrested"
		case MsgNotIntrested:
			return "NotIntrested"
		case MsgBitfield:
			return "Bitfield"
		case MsgRequest:
			return "Request"
		case MsgPiece:
			return "Piece"
		case MsgCancel:
			return "Cancel"
		default:
			return fmt.Sprintf("Unknown #%d", m.ID) 
	}
}

func (m *Message) String() string {
	if m == nil{
		return m.name()
	}

	return fmt.Sprintf("%s [%d]", m.name(), len(m.Payload))
}