package filehandling

import (
	"context"
	"fmt"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"peer-to-peer-app-v0.01/filehandling/recieve"
)

const protocolID = "/rex/fileshare"

type FileSendReqest struct {
	FileName string
	FileType string
	FileSize int
	From     peer.ID
}

func (fhd FileSendReqest) CreatFileRecieveStream(ctx context.Context, host host.Host) network.Stream {
	str, err := host.NewStream(ctx, fhd.From, protocolID)
	if err != nil {
		fmt.Println("Error while creating a new file recieving stream, for file", fhd.FileName)
	}
	return str
}

func (fhd FileSendReqest) initializeFileRecieve() (*recieve.Filerecieve, *recieve.FileRecieveRequest) {
	fileRecieveObj := &recieve.Filerecieve{
		FileName: fhd.FileName,
		FileType: fhd.FileType,
		FileSize: fhd.FileSize,
		From:     fhd.From,
	}
	fileRecieveRequestObj := &recieve.FileRecieveRequest{
		FileRfileRecieveObj: fileRecieveObj,
		Processed:           false,
		Recieved:            true,
	}
	return fileRecieveObj, fileRecieveRequestObj

}

func (fhd FileSendReqest) HandleFileRecieve(fileStr network.Stream) {
	fileRecieveObj, fileRecieveRequestObj := fhd.initializeFileRecieve()
	fileRecieveObj.RecieveFile(fileStr, fileRecieveRequestObj)
	fmt.Println("Status of file recieve")
	fmt.Println(fileRecieveObj.FileName, "[Processed :", fileRecieveRequestObj.Processed, "][Recieved :", fileRecieveRequestObj.Recieved, "]")
}
