package protocol

import (
	"fmt"
	"strconv"
	"unicode"
)

func getLength(m Message) (int, error) {
	lenStr, message := "", m.Message()
	for i := 1; i < len(message[1:]); i++ {
		r := rune(message[i])
		if unicode.IsDigit(r) {
			lenStr += string(r)
		} else {
			break
		}
	}
	msgLen, err := strconv.Atoi(lenStr)
	if err != nil {
		return 0, err
	}
	fmt.Println(msgLen)
	return msgLen, nil
}

func checkCRLF(m Message) error {
	message := m.Message()
	if message[len(message)-2:] != "\r\n" {
		return fmt.Errorf("expected message to end with '\r\n'")
	}
	return nil
}

func checkByte(m Message) error {
	firstByte := rune(m.Message()[0])
	if firstByte != m.Value().Type() {
		return fmt.Errorf("expected message to have first byte %v, got %v", m.Value().Type(), firstByte)
	}
	return nil
}

type Message interface {
	Value() RespValue
	Message() string
	Deserialize() (RespValue, error)
}

type IntegerMessage struct {
	IntegerMessage string
	RespValue      Integer
}

func (im IntegerMessage) Value() RespValue {
	return im.RespValue
}

func (im IntegerMessage) Message() string {
	return im.IntegerMessage
}

func (im IntegerMessage) Deserialize() (RespValue, error) {
	if err := checkCRLF(im); err != nil {
		return nil, err
	}
	if err := checkByte(im); err != nil {
		return nil, err
	}
	return im.RespValue, nil
}

type SimpleStringMessage struct {
	SimpleStringMessage string
	RespValue           SimpleString
}

func (ssm SimpleStringMessage) Value() RespValue {
	return ssm.RespValue
}

func (ssm SimpleStringMessage) Message() string {
	return ssm.SimpleStringMessage
}

func (ssm SimpleStringMessage) Deserialize() (RespValue, error) {
	if err := checkCRLF(ssm); err != nil {
		return nil, err
	}
	if err := checkByte(ssm); err != nil {
		return nil, err
	}
	return ssm.RespValue, nil
}

type BulkStringMessage struct {
	BulkStringMessage string
	RespValue         BulkString
}

func (bsm BulkStringMessage) Value() RespValue {
	return bsm.RespValue
}

func (bsm BulkStringMessage) Message() string {
	return bsm.BulkStringMessage
}

func (bsm BulkStringMessage) Deserialize() (RespValue, error) {
	if err := checkCRLF(bsm); err != nil {
		return nil, err
	}
	if err := checkByte(bsm); err != nil {
		return nil, err
	}
	strLen, err := getLength(bsm)
	if err != nil {
		return nil, err
	}
	actualLen := len(bsm.RespValue.Val)
	if strLen != actualLen {
		return nil, fmt.Errorf("expected bulk string of length %d, got length %d", actualLen, strLen)
	}
	return bsm.RespValue, nil
}

type ErrorMessage struct {
	ErrorMessage string
	RespValue    Error
}

func (em ErrorMessage) Value() RespValue {
	return em.RespValue
}

func (em ErrorMessage) Message() string {
	return em.ErrorMessage
}

func (em ErrorMessage) Deserialize() (RespValue, error) {
	if err := checkCRLF(em); err != nil {
		return nil, err
	}
	if err := checkByte(em); err != nil {
		return nil, err
	}
	return em.RespValue, nil
}

type ArrayMessage struct {
	ArrayMessage string
	RespValue    Array
}

func (am ArrayMessage) Value() RespValue {
	return am.RespValue
}

func (am ArrayMessage) Message() string {
	return am.ArrayMessage
}

func (am ArrayMessage) Deserialize() (RespValue, error) {
	if err := checkCRLF(am); err != nil {
		return nil, err
	}
	if err := checkByte(am); err != nil {
		return nil, err
	}
	arrayLen, err := getLength(am)
	if err != nil {
		return nil, err
	}
	actualLen := len(am.RespValue.Val)
	if arrayLen != actualLen {
		return nil, fmt.Errorf("array message must have length of %d, got %d", actualLen, arrayLen)
	}
	return am.RespValue, nil
}
