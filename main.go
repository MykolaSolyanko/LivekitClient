package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/MykolaSolyanko/LivekitClient/cli"
	"github.com/MykolaSolyanko/LivekitClient/client"

	lksdk "github.com/livekit/server-sdk-go"
)

func main() {
	// lksdk.SetLogger(logger.LogRLogger(logr.Discard()))

	selectOption := cli.New()
	clients := make(map[string]*lksdk.Room)

	apiKey := os.Getenv("LIVEKIT_API_KEY")
	if apiKey == "" {
		fmt.Printf("LIVEKIT_API_KEY is required")
		return
	}

	secretKey := os.Getenv("LIVEKIT_SECRET_KEY")
	if secretKey == "" {
		fmt.Printf("LIVEKIT_SECRET_KEY is required")
		return
	}

	url := os.Getenv("LIVEKIT_URL")
	if url == "" {
		fmt.Printf("LIVEKIT_URL is required")
		return
	}

	for {
		option := selectOption.GetSelectOption()

		switch option {
		// case cli.EnterMessage:
		// 	message := selectOption.EnterMessage()

		// 	participantSid := selectOption.EnterParticipantSid()

		// 	c, ok := clients[participantSid]
		// 	if !ok {
		// 		fmt.Println("Participant not found")
		// 		continue
		// 	}

		// 	if err := client.SendDataMessage(message, c); err != nil {
		// 		fmt.Println(err)
		// 	}

		// case cli.SendReaction:
		// 	fmt.Println("Send Reaction")

		// 	participantSid := selectOption.EnterParticipantSid()

		// 	c, ok := clients[participantSid]
		// 	if !ok {
		// 		fmt.Println("Participant not found")
		// 		continue
		// 	}

		// 	if err := client.SendReaction(c); err != nil {
		// 		fmt.Println(err)
		// 	}

		// case cli.RequestChatHistory:
		// 	fmt.Println("Request Chat History")

		// 	participantSid := selectOption.EnterParticipantSid()

		// 	c, ok := clients[participantSid]
		// 	if !ok {
		// 		fmt.Println("Participant not found")
		// 		continue
		// 	}

		// 	if err := client.SendRequestChatHistory(c); err != nil {
		// 		fmt.Println(err)
		// 	}

		case cli.JoinNewParticipantToRoom:
			fmt.Println("Enter how many participants you want to join to the room: ")
			participamntCount := selectOption.ReadFromStdin()

			count, err := strconv.Atoi(participamntCount)
			if err != nil {
				fmt.Println(err)
				continue
			}

			expectedCount := count * 2
			receiveCount := 0

			notification := make(chan struct{}, 1)
			wait := make(chan struct{})

			go func() {
				for {
					select {
					case <-notification:
						if receiveCount++; receiveCount == expectedCount {
							fmt.Println("All participants joined to the room")
							wait <- struct{}{}
							return
						}
					case <-time.After(5 * time.Second):
						fmt.Println("Timeout")
						wait <- struct{}{}
						return
					}
				}
			}()

			for i := 1; i <= count; i++ {
				room, sid, err := client.JoinToRoom(apiKey, secretKey, url, notification)
				if err != nil {
					fmt.Println(err)

					continue
				}

				clients[sid] = room
			}

			<-wait
			fmt.Printf("Received %d notifications\n", receiveCount)

		case cli.DisconnectParticipantFromRoom:
			participantSid := selectOption.EnterParticipantSid()
			fmt.Println("Disconnect Participant from Room: ", participantSid)

			c, ok := clients[participantSid]
			if !ok {
				fmt.Println("Participant not found")
				continue
			}

			c.Disconnect()
			delete(clients, participantSid)

		case cli.Exit:
			fmt.Println("Exit")

			for _, c := range clients {
				c.Disconnect()
			}

			return
		default:
			fmt.Println("Invalid option")
		}
	}
}
