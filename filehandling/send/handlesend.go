package send

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
)

const buffSize = 1000

type filesend struct {
	fileName string
	fileType string
	To       peer.ID
}

func (fl filesend) getFileSize() (int, error) {
	file, err := os.Stat(fl.fileName)
	if err != nil {
		fmt.Println("Error while opening the file :", fl.fileName)
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
	fmt.Println("Sending  file to :", fl.To)
	filesize, err := fl.getFileSize()
	if err != nil {
		fmt.Println("Error while retrieving the file size")
		fmt.Println("Cannot conitnue further without file size, try again later")
	} else {
		fmt.Println("[File:", fl.fileName, "][Size:", filesize, "]")
		file, err := os.Open(fl.fileName)
		if err != nil {
			fmt.Println("Error while trying to open the file, cannot proceed further")
		} else {
			fl.toStream(str, file, buffSize, filesize)
		}
	}
}
