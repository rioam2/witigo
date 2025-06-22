package witigo

type StringEncoding string

const (
	StringEncodingUTF8  StringEncoding = "utf8"
	StringEncodingUTF16 StringEncoding = "utf16"
)

func (e StringEncoding) CodeUnitSize() int32 {
	switch e {
	case StringEncodingUTF8:
		return 1
	case StringEncodingUTF16:
		return 2
	default:
		return 1
	}
}

func (e StringEncoding) Alignment() int32 {
	switch e {
	case StringEncodingUTF8:
		return 1
	case StringEncodingUTF16:
		return 2
	default:
		return 1
	}
}

type AbiOptions struct {
	StringEncoding StringEncoding
	Memory         RuntimeMemory
	Realloc        func(origPtr int32, origSize int32, alignment int32, newSize int32) (int32, error)
}
