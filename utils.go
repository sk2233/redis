/*
@author: sk
@date: 2024/5/11
*/
package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
)

func HandleErr(err error) {
	if err != nil {
		panic(err)
	}
}

// tar必须是指针
func ReadAny(reader io.Reader, tar any) {
	data := ReadData(reader)
	err := json.Unmarshal(data.Data, tar)
	HandleErr(err)
}

func ReadData(reader io.Reader) *ByteData {
	size := ReadUint64(reader)
	return &ByteData{
		Size: size,
		Data: ReadByte(reader, int(size)),
	}
}

func ReadUint64(reader io.Reader) uint64 {
	bs := ReadByte(reader, 8)
	return binary.BigEndian.Uint64(bs)
}

func ReadByte(reader io.Reader, size int) []byte {
	bs := make([]byte, size)
	count, err := reader.Read(bs)
	HandleErr(err)
	if size != count {
		panic(fmt.Sprintf("need size = %d , has count = %d", size, count))
	}
	return bs
}

func WriteAny(writer io.Writer, data any) {
	bs, err := json.Marshal(data)
	HandleErr(err)
	WriteData(writer, NewByteData(bs))
}

func WriteData(writer io.Writer, data *ByteData) {
	WriteUint64(writer, data.Size)
	WriteByte(writer, data.Data)
}

func WriteUint64(writer io.Writer, data uint64) {
	temp := make([]byte, 8)
	binary.BigEndian.PutUint64(temp, data)
	WriteByte(writer, temp)
}

func WriteByte(writer io.Writer, data []byte) {
	count, err := writer.Write(data)
	HandleErr(err)
	if len(data) != count {
		panic(fmt.Sprintf("need size = %d , has count = %d", len(data), count))
	}
}

func Has[T any](data map[string]T, val string) bool {
	_, ok := data[val]
	return ok
}
