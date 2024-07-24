package client

import (
	"app/bitfield"
	"app/handshake"
	"app/message"
	"app/peers"
	"bytes"
	"fmt"
	"net"
	"time"
)

type Client struct {
	Conn net.Conn
	Choked bool
	Bitfield bitfield.Bitfield
	peer peers.Peer
	infoHash [20]byte
	peerID [20]byte
}

func completeHandshake(conn net.Conn, infoHash, peerID [20]byte) (*handshake.Handshake,error) {
	conn.SetDeadline(time.Now().Add(3 * time.Second))
	defer conn.SetDeadline(time.Time{})

	req := handshake.New(infoHash, peerID)
	_, err := conn.Write(req.Serialize())
	if err != nil{
		return nil, err
	}
	res, err := handshake.Read(conn)
	if err != nil{
		return nil, err
	}
	if !bytes.Equal(res.InfoHash[:], infoHash[:]){
		return nil, fmt.Errorf("expected infohash %x but got %x", res.InfoHash, infoHash)
	}

	return res,nil
}

func recvBitfield(conn net.Conn) (bitfield.Bitfield, error){
	conn.SetDeadline(time.Now().Add(5 * time.Second))
	defer conn.SetDeadline(time.Time{})

	msg, err := message.Read(conn)
	if err != nil{
		return nil, err
	}
	if msg == nil{
		err := fmt.Errorf("expected bitfield but got %s", msg)
		return nil,err
	}
	if msg.ID != message.MsgBitfield{
		err := fmt.Errorf("expected bitfield but got ID %d", msg.ID)
		return nil,err
	}
	return msg.Payload,nil
}

