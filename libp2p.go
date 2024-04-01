package main

import (
	"context"
	"fmt"
	"io"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
)

// Node represents a P2P network node.
type Node struct {
	host host.Host
}

// NewNode creates a new P2P node.
func NewNode() (*Node, error) {
	// Set up libp2p host
	host, err := libp2p.New(
		libp2p.ListenAddrStrings("/ip4/127.0.0.1/tcp/0"),
	)
	if err != nil {
		return nil, err
	}

	return &Node{host: host}, nil
}

// StartListening starts listening for incoming connections.
func (n *Node) StartListening() {
	n.host.SetStreamHandler(protocolID, n.handleStream)
	fmt.Printf("Node with ID %s is listening on addresses:\n", n.host.ID())
	for _, addr := range n.host.Addrs() {
		fmt.Printf("- %s\n", addr)
	}
}

// Connect connects to another node using multiaddress.
func (n *Node) Connect(ctx context.Context, addr string) error {
	maddr, err := multiaddr.NewMultiaddr(addr)
	if err != nil {
		return err
	}
	peerInfo, err := peer.AddrInfoFromP2pAddr(maddr)
	if err != nil {
		return err
	}
	err = n.host.Connect(ctx, *peerInfo)
	if err != nil {
		return err
	}
	fmt.Printf("Connected to peer %s\n", peerInfo.ID)
	return nil
}

// SendMessage sends a message to a specific peer.
func (n *Node) SendMessage(ctx context.Context, peerID peer.ID, message string) error {
	// Convert the message to bytes
	messageBytes := []byte(message)

	// Encrypt the message
	encryptedMessage, err := Encrypt(messageBytes, key)
	if err != nil {
		return fmt.Errorf("failed to encrypt message: %w", err)
	}

	stream, err := n.host.NewStream(ctx, peerID, protocolID)
	if err != nil {
		return err
	}
	defer stream.Close()
	_, err = stream.Write(encryptedMessage)
	if err != nil {
		fmt.Println("Error writing message:", err)
	}
	fmt.Println("Message sent")
	return nil
}

// handleStream handles incoming streams.
func (n *Node) handleStream(stream network.Stream) {
	fmt.Printf("Received message from %s: ", stream.Conn().RemotePeer()+"\n")
	defer stream.Close()

	// Read the encrypted message from the stream
	encryptedMessage, err := io.ReadAll(stream)
	if err != nil {
		fmt.Println("Failed to read from stream:", err)
		return
	}

	// Decrypt the message
	decryptedMessage, err := Decrypt(encryptedMessage, key)
	if err != nil {
		fmt.Println("Failed to decrypt message:", err)
		return
	}
	fmt.Print(string(decryptedMessage))
}
