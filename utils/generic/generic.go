package generic

type Unsigned interface {
	~uint16 | ~uint32 | ~uint64
}

type Numberic interface {
	~int16 | ~int32 | ~int64 | ~float32 | ~float64 | Unsigned
}
