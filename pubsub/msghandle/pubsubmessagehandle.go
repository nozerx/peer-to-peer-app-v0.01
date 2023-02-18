package msghandle

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"strings"

	p2ppubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"peer-to-peer-app-v0.01/filehandling"
	"peer-to-peer-app-v0.01/filehandling/send"
)

type Chatmessage struct {
	Messagecontent string
	Messagefrom    peer.ID
	Authorname     string
}

type Packet struct {
	Type         string
	InnerContent []byte
}

func composeMessage(msg string, host host.Host) *Chatmessage {
	return &Chatmessage{
		Messagecontent: msg,
		Messagefrom:    host.ID(),
		Authorname:     host.ID().ShortString(),
	}
}

func handleInputFromSubscription(ctx context.Context, host host.Host, sub *p2ppubsub.Subscription) {
	inputPacket := &Packet{}
	for {
		inputMsg, err := sub.Next(ctx)
		if err != nil {
			fmt.Println("Error while recieving next message from the subscription")
		} else {
			err := json.Unmarshal(inputMsg.Data, inputPacket)
			if err != nil {
				fmt.Println("Error while unmarshalling at packet level")
			} else {
				if inputPacket.Type == "chat" {
					chatMsg := &Chatmessage{}
					err := json.Unmarshal(inputPacket.InnerContent, chatMsg)
					if err != nil {
						fmt.Println("Error while unmarshalling at chat message level")
					} else {
						fmt.Println("[ BY-> ", inputMsg.ReceivedFrom.Pretty()[len(inputMsg.ReceivedFrom.Pretty())-6:len(inputMsg.ReceivedFrom.Pretty())], "->", chatMsg.Messagecontent)
					}
				}
				if inputPacket.Type == "flsnd" {
					fileRecieveRequest := &filehandling.FileSendReqest{}
					err := json.Unmarshal(inputPacket.InnerContent, fileRecieveRequest)
					if fileRecieveRequest.From == host.ID() {
						continue
					}
					fmt.Println("File recieve request recieved")
					if err != nil {
						fmt.Println("Error while unmarshalling at file recieve request level")
					} else {
						fmt.Println("Request to recieve file ", fileRecieveRequest.FileName, " of type ", fileRecieveRequest.FileType, "of size ", fileRecieveRequest.FileSize)
						fmt.Println("From ", fileRecieveRequest.From)
						fileStream := fileRecieveRequest.CreatFileRecieveStream(ctx, host)
						fileRecieveRequest.HandleFileRecieve(fileStream)
					}
				}
			}
		}
	}
}

func handleInputFromSDI(ctx context.Context, host host.Host, topic *p2ppubsub.Topic) {
	reader := bufio.NewReader(os.Stdin)
	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error during reading from standard input")
		} else {
			if input[:3] == "<s>" {
				fmt.Println("Tag-> <s>")
				escapeSeqLen := 0
				if runtime.GOOS == "windows" {
					escapeSeqLen = 2
				} else {
					escapeSeqLen = 1
				}
				fileName := input[3 : len(input)-escapeSeqLen]
				fileType := strings.Split(fileName, ".")[1]
				fmt.Println("Directive to send the file ", fileName, "of file type :", fileType)
				fmt.Println("the length of file name", len(fileName))
				fileSendObj, err := send.ComposeFileSend(fileName, fileType, host.ID())
				if err != nil {
					fmt.Println(err)
				}
				fileSendReqMsg := fileSendObj.ComposeFileSendRequestMessage()
				inputContent, err := json.Marshal(fileSendReqMsg)
				if err != nil {
					fmt.Println("Error while marshalling at file send request level")
				} else {
					packetContent := &Packet{
						InnerContent: inputContent,
						Type:         "flsnd",
					}
					packet, err := json.Marshal(packetContent)
					if err != nil {
						fmt.Println("Error while marshalling at packet level")
					} else {
						topic.Publish(ctx, packet)
					}
				}
				fmt.Println("Size of file to be sent is :", fileSendObj.FileSize)
				filehandling.FileSendQueue = filehandling.FileSendQueue.Enqueue(fileSendReqMsg)
				fmt.Println("File send queue len", len(filehandling.FileSendQueue))
				for k := 0; k < len(filehandling.FileSendQueue); k++ {
					fmt.Println(filehandling.FileSendQueue[k].FileName)
				}
			} else {
				fmt.Println("Tag-> <c>")
				writeToSubscription(ctx, host, input, topic)
			}
		}
	}
}

func writeToSubscription(ctx context.Context, host host.Host, message string, topic *p2ppubsub.Topic) {
	chatMessage := composeMessage(message, host)
	inputContent, err := json.Marshal(chatMessage)
	if err != nil {
		fmt.Println("Error while marshalling at chatmessage level")
	} else {
		packetContent := &Packet{
			InnerContent: inputContent,
			Type:         "chat",
		}
		packet, err := json.Marshal(packetContent)
		if err != nil {
			fmt.Println("Error while marshalling at packet level")
		} else {
			topic.Publish(ctx, packet)
		}
	}
}

func HandlePubSubMessages(ctx context.Context, host host.Host, sub *p2ppubsub.Subscription, topic *p2ppubsub.Topic) {
	go handleInputFromSubscription(ctx, host, sub)
	handleInputFromSDI(ctx, host, topic)
}
