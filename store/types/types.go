package types

type DataType uint8

const (
	UNKOWN DataType = iota
	RAW
	BOOLEAN
	INT16
	INT32
	INT64
	UINT16
	UINT32
	UINT64
	FLOAT32
	FLOAT64
	STRING
	TIME
	JSON
)
