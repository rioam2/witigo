package abi

import (
	"fmt"

	"golang.org/x/text/encoding/unicode"
)

func ReadString(opts AbiOptions, ptr uint32) (result string, err error) {
	strEncoding := opts.StringEncoding
	strAlignment := strEncoding.Alignment()
	taggedCodeUnitSize := strEncoding.CodeUnitSize()
	strPtr, ok := opts.Memory.ReadUint32Le(ptr)
	if !ok {
		return "", fmt.Errorf("failed to read string pointer at %d", ptr)
	}
	taggedCodeUnits, ok := opts.Memory.ReadUint32Le(ptr + 4)
	if !ok {
		return "", fmt.Errorf("failed to read tagged code units at %d", ptr+4)
	}
	strByteLength := taggedCodeUnits * taggedCodeUnitSize
	if strPtr == 0 {
		return "", fmt.Errorf("invalid string pointer: %d", strPtr)
	}
	if strPtr != AlignTo(strPtr, strAlignment) {
		return "", fmt.Errorf("string pointer %d is not aligned to %d bytes", strPtr, strAlignment)
	}
	if strPtr+strByteLength > opts.Memory.Size() {
		return "", fmt.Errorf("string pointer %d with length %d exceeds memory size %d", strPtr, strByteLength, opts.Memory.Size())
	}
	strData, ok := opts.Memory.Read(strPtr, strByteLength)
	if !ok {
		return "", fmt.Errorf("failed to read string data at %d with length %d", strPtr, strByteLength)
	}
	switch strEncoding {
	case StringEncodingUTF8:
		return string(strData), nil
	case StringEncodingUTF16:
		decoder := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewDecoder()
		str, err := decoder.String(string(strData))
		if err != nil {
			return "", err
		}
		return str, nil
	default:
		return "", fmt.Errorf("unsupported string encoding: %s", strEncoding)
	}
}

func WriteString(opts AbiOptions, str string) (ptr uint32, byteSize uint32, codeUnits uint32, err error) {
	strEncoding := opts.StringEncoding
	strAlignment := strEncoding.Alignment()
	codeUnits = uint32(len(str))
	var strData []byte
	switch strEncoding {
	case StringEncodingUTF8:
		strData = []byte(str)
	case StringEncodingUTF16:
		encoder := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewEncoder()
		strData, err = encoder.Bytes([]byte(str))
		if err != nil {
			return 0, 0, 0, err
		}
	default:
		return 0, 0, 0, fmt.Errorf("unsupported string encoding: %s", strEncoding)
	}
	strByteSize := len(strData)
	ptr, _, err = Call(opts, "cabi_realloc", 0, 0, uint64(strAlignment), uint64(strByteSize))
	if err != nil || ptr == 0 {
		return 0, 0, 0, fmt.Errorf("failed to allocate memory for string: %w", err)
	}
	ok := opts.Memory.Write(ptr, strData)
	if !ok {
		return 0, 0, 0, fmt.Errorf("failed to write string data to memory at %d", ptr)
	}
	return ptr, uint32(strByteSize), codeUnits, nil
}

func FreeString(opts AbiOptions, ptr uint32, size uint32) (err error) {
	if ptr == 0 || size == 0 {
		return nil
	}
	_, _, err = Call(opts, "cabi_realloc", uint64(ptr), uint64(size), 0, 0)
	if err != nil {
		return fmt.Errorf("failed to free string memory at %d: %w", ptr, err)
	}
	return nil
}
