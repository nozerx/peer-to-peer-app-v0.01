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
	// To       peer.ID
}

type FileQueue []FileSendReqest

var FileSendQueue FileQueue

func (fq FileQueue) Enqueue(fsr FileSendReqest) FileQueue {
	fq = append(fq, fsr)
	return fq
}

func (fq FileQueue) Remove(frs FileSendReqest) FileQueue {
	for g := 0; g < len(fq); g++ {
		if fq[g].FileName == frs.FileName {
			if g == 0 {
				fq = fq[1:]
			} else {
				fq = fq[:g+1]
				fq_rest := fq[g+1:]
				for o := 0; o < len(fq_rest); o++ {
					fq = append(fq, fq_rest[o])
				}
			}
		}
	}
	return fq
}

func (fq FileQueue) Dequeue() (FileQueue, FileSendReqest) {
	element := fq[0]
	fq = fq[1:]
	return fq, element
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
		FileRecieveObj: fileRecieveObj,
		Processed:      false,
		Recieved:       false,
	}
	return fileRecieveObj, fileRecieveRequestObj

}

func (fhd FileSendReqest) HandleFileRecieve(fileStr network.Stream) {
	fileRecieveObj, fileRecieveRequestObj := fhd.initializeFileRecieve()
	fileRecieveObj.RecieveFile(fileStr, fileRecieveRequestObj)
	fmt.Println("Status of file recieve")
	fmt.Println(fileRecieveObj.FileName, "[Processed :", fileRecieveRequestObj.Processed, "][Recieved :", fileRecieveRequestObj.Recieved, "]")
}
