package data

import (
	"encoding/csv"
	"fmt"
	"io"
	"konnex/pkg/errors"
	"os"
)

type Node struct {
	ServerURI string
	NodeID    string
}

const columns = 2

var (
	errNotFound = errors.New("file not found")
	// errWriteFile = errors.New("failed de write file")
	errOpenFile  = errors.New("failed to open file")
	errReadFile  = errors.New("failed to read file")
	errEmptyLine = errors.New("empty or incomplete line found in file")
)

func Save(serverUri, nodeID string) error {
	path := "/store/data.csv"
	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
	if err != nil {
		fmt.Println("ERROR | ", err)
	}
	csvWriter := csv.NewWriter(file)
	csvWriter.Write([]string{serverUri, nodeID})
	csvWriter.Flush()

	return nil
}

func ReadAll() ([]Node, error) {
	path := "/store/data.csv"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, errors.Wrap(errNotFound, err)
	}

	file, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, errors.Wrap(errOpenFile, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	nodes := []Node{}
	for {
		l, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, errors.Wrap(errReadFile, err)
		}

		if len(l) < columns {
			return nil, errEmptyLine
		}

		nodes = append(nodes, Node{l[0], l[1]})
	}

	return nodes, nil
}
