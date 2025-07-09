package utils

import (
	"encoding/binary"
	"math"

	"github.com/found-cake/CacheStore/errors"
	"github.com/found-cake/CacheStore/utils/generic"
)

func CheckFloat32Special(value float32) bool {
	value64 := float64(value)
	return CheckFloat64Special(value64)
}

func CheckFloat64Special(value float64) bool {
	return math.IsNaN(value) || math.IsInf(value, 0)
}

func Int16CheckOver(value, delta int16) bool {
	if delta > 0 {
		return value > math.MaxInt16-delta
	}
	if delta < 0 {
		return value < math.MinInt16-delta
	}
	return false
}

func Int32CheckOver(value, delta int32) bool {
	if delta > 0 {
		return value > math.MaxInt32-delta
	}
	if delta < 0 {
		return value < math.MinInt32-delta
	}
	return false
}

func Int64CheckOver(value, delta int64) bool {
	if delta > 0 {
		return value > math.MaxInt64-delta
	}
	if delta < 0 {
		return value < math.MinInt64-delta
	}
	return false
}

func UInt16CheckOverFlow(value, delta uint16) bool {
	return value > math.MaxUint16-delta
}

func UInt32CheckOverFlow(value, delta uint32) bool {
	return value > math.MaxUint32-delta
}

func UInt64CheckOverFlow(value, delta uint64) bool {
	return value > math.MaxUint64-delta
}

func UintCheckUnderFlow[T generic.Unsigned](value, delta T) bool {
	return value < delta
}

func Float32CheckOver(value, delta float32) bool {
	if delta > 0 {
		return value > math.MaxFloat32-delta || (math.MaxFloat32-value) < delta
	}
	if delta < 0 {
		return value < -math.MaxFloat32-delta || (value+math.MaxFloat32) < -delta
	}
	return false
}

func Float64CheckOver(value, delta float64) bool {
	if delta > 0 {
		return value > math.MaxFloat64-delta || (math.MaxFloat64-value) < delta
	}
	if delta < 0 {
		return value < -math.MaxFloat64-delta || (value+math.MaxFloat64) < -delta
	}
	return false
}

func Binary2Int16(data []byte) (int16, error) {
	if ui, err := Binary2UInt16(data); err != nil {
		return 0, err
	} else {
		return int16(ui), nil
	}
}

func Int16toBinary(value int16) []byte {
	return UInt16toBinary(uint16(value))
}

func Binary2Int32(data []byte) (int32, error) {
	if ui, err := Binary2UInt32(data); err != nil {
		return 0, err
	} else {
		return int32(ui), nil
	}
}

func Int32toBinary(value int32) []byte {
	return UInt32toBinary(uint32(value))
}

func Binary2Int64(data []byte) (int64, error) {
	if ui, err := Binary2UInt64(data); err != nil {
		return 0, err
	} else {
		return int64(ui), nil
	}
}

func Int64toBinary(value int64) []byte {
	return UInt64toBinary(uint64(value))
}

func Binary2UInt16(data []byte) (uint16, error) {
	if len(data) != 2 {
		return 0, errors.ErrInvalidDataLength(2, len(data))
	}
	return binary.LittleEndian.Uint16(data), nil
}

func UInt16toBinary(value uint16) []byte {
	buffer := make([]byte, 2)
	binary.LittleEndian.PutUint16(buffer, value)
	return buffer
}

func Binary2UInt32(data []byte) (uint32, error) {
	if len(data) != 4 {
		return 0, errors.ErrInvalidDataLength(4, len(data))
	}
	return binary.LittleEndian.Uint32(data), nil
}

func UInt32toBinary(value uint32) []byte {
	buffer := make([]byte, 4)
	binary.LittleEndian.PutUint32(buffer, value)
	return buffer
}

func Binary2UInt64(data []byte) (uint64, error) {
	if len(data) != 8 {
		return 0, errors.ErrInvalidDataLength(8, len(data))
	}
	return binary.LittleEndian.Uint64(data), nil
}

func UInt64toBinary(value uint64) []byte {
	buffer := make([]byte, 8)
	binary.LittleEndian.PutUint64(buffer, value)
	return buffer
}

func Binary2Float32(data []byte) (float32, error) {
	if ui, err := Binary2UInt32(data); err != nil {
		return 0, err
	} else {
		return math.Float32frombits(ui), nil
	}
}

func Float32toBinary(value float32) []byte {
	return UInt32toBinary(math.Float32bits(value))
}

func Binary2Float64(data []byte) (float64, error) {
	if ui, err := Binary2UInt64(data); err != nil {
		return 0, err
	} else {
		return math.Float64frombits(ui), nil
	}
}

func Float64toBinary(value float64) []byte {
	return UInt64toBinary(math.Float64bits(value))
}
