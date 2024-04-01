package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/libp2p/go-libp2p/core/peer"
)

func main() {
	ctx := context.Background()

	node, err := NewNode()
	if err != nil {
		fmt.Println("Error creating node:", err)
		return
	}

	node.StartListening()

	go func() {
		peerChan := initMDNS(node.host, serviceTag) // Initiate mDNS in a separate goroutine
		for pr := range peerChan {                  // Continuously listen for new peers
			if pr.ID > node.host.ID() {
				fmt.Print("Found peer:", pr, " id is greater than us, wait for it to connect to us\n")
				continue
			}

			fmt.Print("Found peer:", pr, ", connecting\n")
			if err := node.host.Connect(ctx, pr); err != nil {
				fmt.Println("Connection failed:", err)
				continue
			}
			fmt.Print("Connected to:", pr.String()+"\n")
			bufio.NewReader(os.Stdin).ReadBytes('\n')

		}
	}()

	printInstructions()
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input == "/quit" {
			return
		} else if strings.HasPrefix(input, "/connect") {
			parts := strings.Split(input, " ")
			if len(parts) != 2 {
				fmt.Println("Invalid command")
				continue
			}
			addr := parts[1]
			err := node.Connect(ctx, addr)
			if err != nil {
				fmt.Println("Error connecting to node:", err)
			}
		} else if input == "/peers" {
			fmt.Println("Connected peers:")
			for _, p := range node.host.Network().Peers() {
				fmt.Println("-", p)
			}
		} else if strings.HasPrefix(input, "/send") {
			parts := strings.Split(input, " ")
			if len(parts) < 3 {
				fmt.Println("Invalid command")
				continue
			}
			peerID, err := peer.Decode(parts[1])
			if err != nil {
				fmt.Println("Invalid peer ID:", err)
				continue
			}

			msg := strings.Join(parts[2:], " ")

			if err := node.SendMessage(ctx, peerID, msg); err != nil {
				fmt.Println("Error sending message:", err)
			}
			fmt.Printf("Send message to peer ID %s with messasge: %s \n", peerID, msg)

			continue

		} else {
			printInstructions()
		}
	}
}

func printInstructions() {
	fmt.Println("\nEnter '/connect <multiaddress>' to connect to another node")
	fmt.Println("Enter '/peers' to display connected peers")
	fmt.Println("Enter '/send <peerId> <message>' to send message to another node")
	fmt.Println("Enter '/quit' to exit")
}
