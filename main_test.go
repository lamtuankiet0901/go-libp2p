package main

import (
	"context"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/stretchr/testify/assert"
)

func TestCreateHost(t *testing.T) {
	n, _ := NewNode()
	defer n.host.Close()

	assert.NotNil(t, n.host, "Host should not be nil")
	assert.NotEmpty(t, n.host.Addrs(), "Host should have listen addresses")
}

func TestConnectToPeer(t *testing.T) {
	ctx := context.Background()

	n1, _ := NewNode()
	defer n1.host.Close()

	n2, _ := NewNode()
	defer n2.host.Close()

	addr := n2.host.Addrs()[0].String() + "/p2p/" + n2.host.ID().String()
	n1.Connect(ctx, addr)

	assert.Contains(t, n1.host.Network().Peers(), n2.host.ID(), "Host1 should be connected to Host2")
}

// Assume this function is adjusted for testability
func setupStreamHandlerForTest(h host.Host, messages *[]string) {
	h.SetStreamHandler(protocolID, func(stream network.Stream) {
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
		if err == nil {
			*messages = append(*messages, string(decryptedMessage))
		}
	})
}

func TestMessageSendReceive(t *testing.T) {
	ctx := context.Background()
	messages := make([]string, 0)

	n1, _ := NewNode()
	defer n1.host.Close()

	n2, _ := NewNode()
	defer n2.host.Close()

	n1.StartListening()
	setupStreamHandlerForTest(n2.host, &messages)

	addr := n2.host.Addrs()[0].String() + "/p2p/" + n2.host.ID().String()
	n1.Connect(ctx, addr)

	n1.SendMessage(ctx, n2.host.ID(), "Hello, peer!")

	// In a real test, use synchronization or sleep for demo purposes only
	// In production code, consider using channels, wait groups, or other synchronization mechanisms
	time.Sleep(time.Second * 1)

	assert.Contains(t, messages, "Hello, peer!", "Host2 should receive the message from Host1")
}
