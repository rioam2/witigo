package witigo

type RuntimeMemory interface {
	Read(ptr int32, size int32) ([]byte, error)
	Write(ptr int32, data []byte) error
	Size() int32
}
