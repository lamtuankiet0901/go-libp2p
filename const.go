package main

import (
	libp2pProtocol "github.com/libp2p/go-libp2p/core/protocol"
)

const protocolID = libp2pProtocol.ID("/myp2pnode/1.0.0")
const serviceTag = "p2pnode-discovery"

var key = []byte("examplekey123456")
