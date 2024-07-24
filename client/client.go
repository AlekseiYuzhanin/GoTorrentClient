package client

import (
	"app/bitfield"
	"app/handshake"
	"app/peers"
	"net"
)

type Client struct {
	Conn net.Conn
	Choked bool
	Bitfield bitfield.Bitfield
	peer peers.Peer
	infoHash [20]byte
	peerID [20]byte
}

func completeHandshake(conn net.Conn, infoHash, peerID [20]byte) (*handshake.Handshake,error)