package client

import (
	"errors"
	"fmt"
	"log"
	"time"

	pb "github.com/MykolaSolyanko/LivekitClient/protocol"
	"github.com/livekit/protocol/livekit"
	lksdk "github.com/livekit/server-sdk-go"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var i int

func GenerateUUID() (string, error) {
	i++
	testClient := fmt.Sprintf("testClient%d", i)

	return testClient, nil
}

func JoinToRoom(api, secret, url string, notification chan struct{}) (room *lksdk.Room, uuid string, err error) {
	if uuid, err = GenerateUUID(); err != nil {
		return nil, "", err
	}

	roomPublisherCallback := &lksdk.RoomCallback{
		ParticipantCallback: lksdk.ParticipantCallback{
			OnDataReceived: func(data []byte, rp *lksdk.RemoteParticipant) {
				dataMessage := &pb.Message{}
				if err := proto.Unmarshal(data, dataMessage); err != nil {
					fmt.Printf("Error: %s", err)

					return
				}

				responseHandler(dataMessage, rp, notification)
			},
		},
	}
	if room, err = lksdk.ConnectToRoom(url, lksdk.ConnectInfo{
		APIKey:              api,
		APISecret:           secret,
		RoomName:            "DemoClient",
		ParticipantIdentity: uuid,
		ParticipantName:     uuid,
	}, roomPublisherCallback); err != nil {
		return nil, "", err
	}

	fmt.Printf("Participant %s connected to room %s\n", uuid, room.Name())

	fmt.Printf("Participant %s SID: %s\n", uuid, room.LocalParticipant.SID())

	return room, room.LocalParticipant.Identity(), nil
}

func SendDataMessage(text string, room *lksdk.Room) error {
	id := room.LocalParticipant.Identity()

	chat := &pb.Chat{
		Message:  text,
		DateTime: timestamppb.New(time.Now().UTC()),
		Sender:   &id,
	}

	msg := &pb.Message{
		MessageType: &pb.Message_Chat{ // Note the field name and how it's used here
			Chat: chat,
		},
	}

	data, err := proto.Marshal(msg)
	if err != nil {
		log.Fatal("marshaling error: ", err)
	}

	if err := room.LocalParticipant.PublishData(data, livekit.DataPacket_RELIABLE, nil); err != nil {
		return err
	}

	return nil
}

func SendRequestChatHistory(room *lksdk.Room) error {
	chatHistory := &pb.ChatHistoryRequest{
		Offset: -1,
		Limit:  50,
	}

	msg := &pb.Message{
		MessageType: &pb.Message_ChatHistoryRequest{ // Note the field name and how it's used here
			ChatHistoryRequest: chatHistory,
		},
	}

	data, err := proto.Marshal(msg)
	if err != nil {
		return err
	}

	if err := room.LocalParticipant.PublishData(data, livekit.DataPacket_RELIABLE, nil); err != nil {
		return err
	}

	return nil
}

func SendReaction(room *lksdk.Room) error {
	reaction := &pb.Reaction{
		Emoji: "👍",
	}

	msg := &pb.Message{
		MessageType: &pb.Message_Reaction{ // Note the field name and how it's used here
			Reaction: reaction,
		},
	}

	data, err := proto.Marshal(msg)
	if err != nil {
		return err
	}

	if err := room.LocalParticipant.PublishData(data, livekit.DataPacket_RELIABLE, nil); err != nil {
		return err
	}

	return nil
}

func responseHandler(dataMessage *pb.Message, rp *lksdk.RemoteParticipant, notification chan struct{}) error {
	switch msg := dataMessage.MessageType.(type) {
	case *pb.Message_RoomState:
		fmt.Printf("Receive from participant %s room state message: %s\n", rp.Identity(), msg.RoomState)
		notification <- struct{}{}

	case *pb.Message_ChatHistoryResponse:
		fmt.Printf("Receive from participant %s chat history response message: %s\n", rp.Identity(), msg.ChatHistoryResponse)

	case *pb.Message_Statistics:
		fmt.Printf("Receive from participant %s statistics message: %s\n", rp.Identity(), msg.Statistics)
		notification <- struct{}{}

	default:
		return errors.New(fmt.Sprintf("invalid message type %T", dataMessage.MessageType))
	}

	return nil
}
