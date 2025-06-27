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

func (t DataType) String() string {
	switch t {
	case RAW:
		return "Raw"
	case BOOLEAN:
		return "Boolean"
	case INT16:
		return "Integer16"
	case INT32:
		return "Integer32"
	case INT64:
		return "Integer64"
	case UINT16:
		return "unsinged integer16"
	case UINT32:
		return "Unsinged Integer32"
	case UINT64:
		return "Unsigned Integer64"
	case FLOAT32:
		return "Float32"
	case FLOAT64:
		return "Float64"
	case STRING:
		return "String"
	case TIME:
		return "Time"
	case JSON:
		return "Json"
	}
}
