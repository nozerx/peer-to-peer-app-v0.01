package send

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"peer-to-peer-app-v0.01/filehandling"
)

const buffSize = 1000

type filesend struct {
	FileName string
	FileType string
	FileSize int
	From     peer.ID
}

func ComposeFileSend(filename string, filetype string, host host.Host) (filesend, error) {
	fileSendObj := &filesend{
		FileName: filename,
		FileType: filetype,
		FileSize: 0,
		From:     host.ID(),
	}
	filesize, err := fileSendObj.getFileSize()
	fileSendObj.FileSize = filesize
	if err != nil {
		fmt.Println("Error while retrieving file size for file send reuest object")
		fmt.Println("Cannot continue with file send request")
		return *fileSendObj, fmt.Errorf("Error while trying to send the file %s", fileSendObj.FileName)
	}
	return *fileSendObj, nil
}

func (fl filesend) ComposeFileSendRequestMessage() filehandling.FileSendReqest {
	fileSendRequestObj := &filehandling.FileSendReqest{
		FileName: fl.FileName,
		FileType: fl.FileType,
		FileSize: fl.FileSize,
		From:     fl.From,
	}
	return *fileSendRequestObj
}

func (fl filesend) getFileSize() (int, error) {
	file, err := os.Stat(fl.FileName)
	if err != nil {
		fmt.Println("Error while opening the file :", fl.FileName)
		return 0, err
	}
	return int(file.Size()), nil
}

func (fl filesend) toStream(str network.Stream, file *os.File, bufferSize int, fileSize int) {
	iterationcount := fileSize / bufferSize
	buffer := make([]byte, bufferSize)
	streamWriter := bufio.NewWriter(str)
	for i := 0; i < iterationcount; i++ {
		_, err := file.Read(buffer)
		if err == io.EOF {
			fmt.Println("File send to stream completely")
			break
		}
		if err != nil {
			fmt.Println("Error while reading from the file")
		}
		sendByte, err := streamWriter.Write(buffer)
		if err != nil {
			fmt.Println("Error while sending the buffer to the stream")
		} else {
			fmt.Println("Send ", sendByte, " bytes to stream")
		}

	}
	leftByte := fileSize % bufferSize
	additionalBuffer := make([]byte, leftByte)
	_, err := file.Read(additionalBuffer)
	if err != nil {
		if err == io.EOF {
			fmt.Println("File send to stream completely")
		} else {
			fmt.Println("Error while reading the last buffer")

		}
	}
	sendByte, err := streamWriter.Write(additionalBuffer)
	if err != nil {
		fmt.Println("Error while sending the buffer to the stream")
	} else {
		fmt.Println("Send ", sendByte, " bytes to stream")
	}

}

func (fl filesend) SendFile(str network.Stream) {
	fmt.Println("Sending  file from :", fl.From)
	filesize, err := fl.getFileSize()
	if err != nil {
		fmt.Println("Error while retrieving the file size")
		fmt.Println("Cannot conitnue further without file size, try again later")
	} else {
		fmt.Println("[File:", fl.FileName, "][Size:", filesize, "]")
		file, err := os.Open(fl.FileName)
		if err != nil {
			fmt.Println("Error while trying to open the file, cannot proceed further")
		} else {
			fl.toStream(str, file, buffSize, filesize)
			err := file.Close()
			if err != nil {
				fmt.Println("Error while trying to close ", fl.FileName)
			}
		}
	}
}
