package recieve

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
)

type FileRecieveRequest struct {
	FileRecieveObj *Filerecieve
	Processed      bool
	Recieved       bool
}

type Filerecieve struct {
	FileName string
	FileType string
	FileSize int
	From     peer.ID
}

const buffSize = 1

func (fr Filerecieve) fromStream(str network.Stream, file *os.File, bufferSize int, fileRecReq *FileRecieveRequest) {
	streamReader := bufio.NewReader(str)
	buffer := make([]byte, bufferSize)
	iterationCount := fr.FileSize / bufferSize
	for i := 0; i < iterationCount; i++ {
		_, err := streamReader.Read(buffer)
		if err != nil {
			if err == io.EOF {
				fmt.Println("End of file recieved")
				fileRecReq.Recieved = true
			} else {
				fmt.Println("Error while reading file bytes from stream")
			}
		} else {
			_, err := file.Write(buffer)
			if err != nil {
				fmt.Println("Error while writing to the recieving file")
			}
		}
	}
	leftByte := fr.FileSize % bufferSize
	additonalBuffer := make([]byte, leftByte)
	_, err := streamReader.Read(additonalBuffer)
	if err != nil {
		if err == io.EOF {
			fmt.Println("End of file recieved")
		} else {
			fmt.Println("Error while reading file bytes from stream for last piece")
		}
	} else {
		_, err := file.Write(additonalBuffer)
		if err != nil {
			fmt.Println("Error while writing to the recieving file")
		} else {
			fileRecReq.Recieved = true
			fmt.Println("Completely recieved the file")
		}
	}

}

func (fr Filerecieve) RecieveFile(str network.Stream, fileRecReq *FileRecieveRequest) {
	file, err := os.Create(fr.FileName)
	// time.Sleep(60 * time.Second)
	if err != nil {
		fmt.Println("Error while creating the recieving file ", fr.FileName)
	} else {
		fileRecReq.Processed = true
		fr.fromStream(str, file, buffSize, fileRecReq)
		err := file.Close()
		if err != nil {
			fmt.Println("Error while closing the reciving file")
		}
		str.Close()
	}
}
