package send

import (
	"fmt"

	"github.com/libp2p/go-libp2p/core/network"
	"peer-to-peer-app-v0.01/filehandling"
)

func HandleInputStream(str network.Stream) {
	fmt.Println("Inuput File stream handler called")
	_, fileSendRequest := filehandling.FileSendQueue.Dequeue()
	fileSendObj, err := ComposeFileSend(fileSendRequest.FileName, fileSendRequest.FileType, fileSendRequest.From)
	if err != nil {
		fmt.Println("Error while trying to send file", fileSendObj.FileName)
	} else {
		go fileSendObj.SendFile(str)
		filehandling.FileSendQueue, _ = filehandling.FileSendQueue.Dequeue()
	}
	// This part is reserved for future use case, where we can properly implement the file send and recieve queue

	// for h := 0; h < len(filehandling.FileSendQueue); h++ {
	// 	if str.Conn().RemotePeer() == filehandling.FileSendQueue[h].From {
	// 		fileSendObj, err := ComposeFileSend(filehandling.FileSendQueue[h].FileName, filehandling.FileSendQueue[h].FileType, filehandling.FileSendQueue[h].From)
	// 		if err != nil {
	// 			fmt.Println("Error while trying to send file ", filehandling.FileSendQueue[h].FileName)
	// 		} else {
	// 			go fileSendObj.SendFile(str)
	// 			go filehandling.FileSendQueue.Remove(filehandling.FileSendReqest(fileSendObj))
	// 		}
	// 	}
	// }
}
