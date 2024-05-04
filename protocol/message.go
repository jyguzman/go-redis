package protocol

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

func getLength(m Deserializer, idx int) (int, int, error) {
	idx += 1
	lenStr, message := "", m.Serialized()[idx:]
	if message[0:2] == "-1" {
		return -1, 5, nil
	}
	for _, r := range message {
		if unicode.IsDigit(r) {
			lenStr += string(r)
			idx += 1
		} else {
			break
		}
	}
	msgLen, err := strconv.Atoi(lenStr)
	if err != nil {
		return 0, idx, err
	}
	return msgLen, idx + 2, nil
}

func checkCRLF(m Deserializer, idx int) error {
	message := m.Serialized()[idx:]
	if message[len(message)-2:] != "\r\n" {
		return fmt.Errorf("expected message to end with '\r\n'")
	}
	return nil
}

func checkByte(m Deserializer, idx int) error {
	firstByte := rune(m.Serialized()[idx])
	if firstByte != m.Type() {
		return fmt.Errorf("expected message to have first byte %s, got %s", string(m.Type()), string(firstByte))
	}
	return nil
}

func DeserializeMessage(message string) (RespValue, int, error) {
	switch message[0] {
	case RespInteger:

		return IntDeserializer{Message: message}.Deserialize(0)
	case RespSimpleString:

		return SimpleStringDeserializer{Message: message}.Deserialize(0)
	case RespError:
		return ErrorDeserializer{Message: message}.Deserialize(0)
	case RespBulkString:
		return BulkStringDeserializer{Message: message}.Deserialize(0)
	case RespArray:
		return ArrayDeserializer{Message: message}.Deserialize(0)
	}
	return nil, -1, fmt.Errorf("recevied invalid RESP type byte %s", string(message[0]))
}

type Deserializer interface {
	Type() rune
	Serialized() string
	Deserialize(int) (RespValue, int, error)
}

type IntDeserializer struct {
	Message string
}

func (i IntDeserializer) Type() rune {
	return RespInteger
}

func (i IntDeserializer) Serialized() string {
	return i.Message
}

func (i IntDeserializer) Deserialize(idx int) (RespValue, int, error) {
	if err := checkCRLF(i, idx); err != nil {
		return nil, -1, err
	}
	if err := checkByte(i, idx); err != nil {
		return nil, -1, err
	}
	numStr := ""
	idx += 1
	for _, r := range i.Message[idx:] {
		if r != '\r' {
			numStr += string(r)
			idx += 1
		} else {
			break
		}
	}
	val, err := strconv.Atoi(numStr)
	if err != nil {
		return nil, -1, err
	}
	idx += 2
	return Integer{Val: val}, idx, nil
}

type SimpleStringDeserializer struct {
	Message string
}

func (s SimpleStringDeserializer) Type() rune {
	return RespSimpleString
}

func (s SimpleStringDeserializer) Serialized() string {
	return s.Message
}

func (s SimpleStringDeserializer) Deserialize(idx int) (RespValue, int, error) {
	if err := checkCRLF(s, idx); err != nil {
		return nil, -1, err
	}
	if err := checkByte(s, idx); err != nil {
		return nil, -1, err
	}
	idx += 1
	str := ""
	for _, r := range s.Message[idx:] {
		if r != '\r' {
			str += string(r)
			idx += 1
		} else {
			break
		}
	}
	idx += 2
	return SimpleString{Val: str}, idx, nil
}

type ErrorDeserializer struct {
	Message string
}

func (e ErrorDeserializer) Type() rune {
	return RespError
}

func (e ErrorDeserializer) Serialized() string {
	return e.Message
}

func (e ErrorDeserializer) Deserialize(idx int) (RespValue, int, error) {
	if err := checkCRLF(e, idx); err != nil {
		return nil, -1, err
	}
	if err := checkByte(e, idx); err != nil {
		return nil, -1, err
	}
	errStr := ""
	for _, r := range e.Message[idx:] {
		if r != '\r' {
			errStr += string(r)
			idx += 1
		} else {
			break
		}
	}
	idx += 2
	return Error{Val: errStr}, idx, nil
}

type BulkStringDeserializer struct {
	Message string
}

func (bs BulkStringDeserializer) Type() rune {
	return RespBulkString
}

func (bs BulkStringDeserializer) Serialized() string {
	return bs.Message
}

func (bs BulkStringDeserializer) Deserialize(idx int) (RespValue, int, error) {
	if err := checkCRLF(bs, idx); err != nil {
		return nil, -1, err
	}
	if err := checkByte(bs, idx); err != nil {
		return nil, -1, err
	}
	msg := bs.Message[idx:]
	//length, skip, err := getLength(bs, idx)
	fields := strings.Split(msg, "\r\n")
	lengthStr, txt := fields[0][1:], fields[1]
	length, err := strconv.Atoi(lengthStr)
	if err != nil {
		return nil, -1, err
	}
	if length == -1 {
		return Nil{NilType: RespBulkString}, 5, nil
	}
	if length != len(txt) {
		return nil, -1, fmt.Errorf("expected bulk string length %v, got %v", len(msg), length)
	}
	idx += len(fields[0]) + 2 + len(txt) + 2
	return BulkString{Val: txt}, idx, nil
}

type ArrayDeserializer struct {
	Message string
}

func (a ArrayDeserializer) Type() rune {
	return RespArray
}

func (a ArrayDeserializer) Serialized() string {
	return a.Message
}

func (a ArrayDeserializer) Deserialize(idx int) (RespValue, int, error) {
	if err := checkCRLF(a, idx); err != nil {
		return nil, -1, err
	}
	if err := checkByte(a, idx); err != nil {
		return nil, -1, err
	}
	numElements, msgIdx, err := getLength(a, idx)
	if numElements == -1 {
		return Nil{NilType: RespArray}, msgIdx, nil
	}
	if err != nil {
		return nil, -1, err
	}
	arr, arrIdx := make([]RespValue, numElements), 0
	for arrIdx < numElements {
		respVal, skip, err := DeserializeMessage(a.Message[msgIdx:])
		if err != nil {
			return nil, -1, err
		}
		arr[arrIdx] = respVal
		arrIdx += 1
		msgIdx += skip
	}
	return &Array{Val: arr}, msgIdx, nil
}
