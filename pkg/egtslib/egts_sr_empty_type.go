package egts

import (
	"bytes"
	"fmt"
)

type SrEmptyDataType struct {
	firstByte  uint8 `json:"fB"`
	secondByte uint8 `json:"sB"`
	thirdByte  uint8 `json:"tB"`
	fourthByte uint8 `json:"foB"`
	fifthByte  uint8 `json:"fiB"`
}

func (e *SrEmptyDataType) Decode(content []byte) error {

	var (
		err error
	)

	buf := bytes.NewReader(content)

	if e.firstByte, err = buf.ReadByte(); err != nil {
		return fmt.Errorf("Не удалось получить значение первого байта: %v", err)
	}

	if e.secondByte, err = buf.ReadByte(); err != nil {
		return fmt.Errorf("Не удалось получить значение второго байта: %v", err)
	}

	if e.thirdByte, err = buf.ReadByte(); err != nil {
		return fmt.Errorf("Не удалось получить значение третьего байта: %v", err)
	}

	if e.fourthByte, err = buf.ReadByte(); err != nil {
		return fmt.Errorf("Не удалось получить значение четвертого байта: %v", err)
	}

	if e.fifthByte, err = buf.ReadByte(); err != nil {
		return fmt.Errorf("Не удалось получить значение пятого байта: %v", err)
	}

	return err
}

func (e *SrEmptyDataType) Encode() ([]byte, error) {
	var (
		err    error
		result []byte
	)
	buf := new(bytes.Buffer)

	if err = buf.WriteByte(e.firstByte); err != nil {
		return result, fmt.Errorf("Не удалось записать значение первого байта: %v", err)
	}

	if err = buf.WriteByte(e.secondByte); err != nil {
		return result, fmt.Errorf("Не удалось записать значение первого байта: %v", err)
	}

	if err = buf.WriteByte(e.thirdByte); err != nil {
		return result, fmt.Errorf("Не удалось записать значение первого байта: %v", err)
	}

	if err = buf.WriteByte(e.fourthByte); err != nil {
		return result, fmt.Errorf("Не удалось записать значение первого байта: %v", err)
	}

	if err = buf.WriteByte(e.fifthByte); err != nil {
		return result, fmt.Errorf("Не удалось записать значение первого байта: %v", err)
	}

	result = buf.Bytes()

	return result, err
}

func (e *SrEmptyDataType) Length() uint16 {
	var result uint16

	if recBytes, err := e.Encode(); err != nil {
		result = uint16(0)
	} else {
		result = uint16(len(recBytes))
	}

	return result
}
