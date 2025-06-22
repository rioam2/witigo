package witigo

type StringEncoding string

const (
	StringEncodingUTF8        StringEncoding = "utf8"
	StringEncodingUTF16       StringEncoding = "utf16"
	StringEncodingLatin1Utf16 StringEncoding = "latin1+utf16"
)

type AbiOptions struct {
	StringEncoding StringEncoding
	Memory         RuntimeMemory
	Realloc        func(origPtr int32, origSize int32, alignment int32, newSize int32) (int32, error)
}
