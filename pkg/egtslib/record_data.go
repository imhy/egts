package egts

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

//RecordData структура секции подзапси у записи ServiceDataRecord
type RecordData struct {
	SubrecordType   byte       `json:"SRT"`
	SubrecordLength uint16     `json:"SRL"`
	SubrecordData   BinaryData `json:"SRD"`
}

//RecordDataSet описывает массив с подзаписями протокола ЕГТС
type RecordDataSet []RecordData

//Decode разбирает байты в структуру подзаписи
func (rds *RecordDataSet) Decode(recDS []byte) error {
	var (
		err error
	)
	buf := bytes.NewBuffer(recDS)
	for buf.Len() > 0 {
		rd := RecordData{}
		if rd.SubrecordType, err = buf.ReadByte(); err != nil {
			return fmt.Errorf("Не удалось получить тип записи subrecord data: %v", err)
		}

		tmpIntBuf := make([]byte, 2)
		if _, err = buf.Read(tmpIntBuf); err != nil {
			return fmt.Errorf("Не удалось получить длину записи subrecord data: %v", err)
		}
		rd.SubrecordLength = binary.LittleEndian.Uint16(tmpIntBuf)

		subRecordBytes := buf.Next(int(rd.SubrecordLength))

		switch rd.SubrecordType {
		case SrPosDataType:
			rd.SubrecordData = &SrPosData{}
		case SrTermIdentityType:
			rd.SubrecordData = &SrTermIdentity{}
		case SrRecordResponseType:
			rd.SubrecordData = &SrResponse{}
		case SrResultCodeType:
			rd.SubrecordData = &SrResultCode{}
		case SrExtPosDataType:
			rd.SubrecordData = &SrExtPosData{}
		case SrAdSensorsDataType:
			rd.SubrecordData = &SrAdSensorsData{}
		case SrType20:
			// признак косвенный в спецификациях его нет
			if rd.SubrecordLength == uint16(5) {
				rd.SubrecordData = &SrStateData{}
			} else {
				// TODO: добавить секцию EGTS_SR_ACCEL_DATA
				return fmt.Errorf("Не реализованная секция EGTS_SR_ACCEL_DATA: %d. Длина: %d. Содержимое: %X", rd.SubrecordType, rd.SubrecordLength, subRecordBytes)
			}
		case SrStateDataType:
			rd.SubrecordData = &SrStateData{}
		case SrLiquidLevelSensorType:
			rd.SubrecordData = &SrLiquidLevelSensor{}
		case SrAbsCntrDataType:
			rd.SubrecordData = &SrAbsCntrData{}
		case SrAuthInfoType:
			rd.SubrecordData = &SrAuthInfo{}
		case SrCountersDataType:
			rd.SubrecordData = &SrCountersData{}
		case SrEgtsPlusDataType:
			rd.SubrecordData = &StorageRecord{}
		case SrAbsAnSensDataType:
			rd.SubrecordData = &SrAbsAnSensData{}
		case SrEmptyType:
			rd.SubrecordData = &SrEmptyDataType{}
		default:
			return fmt.Errorf("Не известный тип подзаписи: %d. Длина: %d. Содержимое: %X", rd.SubrecordType, rd.SubrecordLength, subRecordBytes)
		}

		if err = rd.SubrecordData.Decode(subRecordBytes); err != nil {
			return err
		}
		*rds = append(*rds, rd)
	}

	return err
}

//Encode преобразовывает подзапись в набор байт
func (rds *RecordDataSet) Encode() ([]byte, error) {
	var (
		result []byte
		err    error
	)
	buf := new(bytes.Buffer)

	for _, rd := range *rds {
		if err := binary.Write(buf, binary.LittleEndian, rd.SubrecordType); err != nil {
			return result, err
		}

		if err := binary.Write(buf, binary.LittleEndian, rd.SubrecordLength); err != nil {
			return result, err
		}

		srd, err := rd.SubrecordData.Encode()
		if err != nil {
			return result, err
		}

		buf.Write(srd)
	}

	result = buf.Bytes()

	return result, err
}

//Length получает длину массива записей
func (rds *RecordDataSet) Length() uint16 {
	var result uint16

	if recBytes, err := rds.Encode(); err != nil {
		result = uint16(0)
	} else {
		result = uint16(len(recBytes))
	}

	return result
}
