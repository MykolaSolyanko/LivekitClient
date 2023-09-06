package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Cli struct {
	reader *bufio.Reader
}

const (
	JoinNewParticipantToRoom      = "1"
	DisconnectParticipantFromRoom = "2"
	EnterMessage                  = "3"
	SendReaction                  = "4"
	RequestChatHistory            = "5"
	Exit                          = "6"
)

func New() *Cli {
	return &Cli{
		reader: bufio.NewReader(os.Stdin),
	}
}

const showMenu = `
1. Join new Participant to Room
2. Disconnect Participant from Room
3. Enter Message
4. Send Reaction
5. Request Chat History
6. Exit
`

func (c *Cli) GetSelectOption() string {
	fmt.Print(showMenu)
	fmt.Print("Select an option: ")

	option, _ := c.reader.ReadString('\n')

	return strings.TrimSpace(option)
}

func (c *Cli) EnterMessage() string {
	fmt.Print("Enter message: ")

	message, _ := c.reader.ReadString('\n')

	return strings.TrimSpace(message)
}

func (c *Cli) EnterParticipantSid() string {
	fmt.Print("Enter Participant SID: ")

	participantSid, _ := c.reader.ReadString('\n')

	return strings.TrimSpace(participantSid)
}

func (c *Cli) ReadFromStdin() string {
	text, _ := c.reader.ReadString('\n')

	return strings.TrimSpace(text)
}
