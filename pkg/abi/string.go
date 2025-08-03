package abi

import (
	"fmt"

	"golang.org/x/text/encoding/unicode"
)

func ReadString(opts AbiOptions, ptr uint32) (string, error) {
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

func WriteString(opts AbiOptions, str string) (uint32, uint32, error) {
	strEncoding := opts.StringEncoding
	strAlignment := strEncoding.Alignment()
	var codeUnits uint32 = uint32(len(str))
	var strData []byte
	var err error
	switch strEncoding {
	case StringEncodingUTF8:
		strData = []byte(str)
	case StringEncodingUTF16:
		encoder := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewEncoder()
		strData, err = encoder.Bytes([]byte(str))
		if err != nil {
			return 0, 0, err
		}
	default:
		return 0, 0, fmt.Errorf("unsupported string encoding: %s", strEncoding)
	}
	ptr, err := Call(opts, "cabi_realloc", 0, 0, uint64(strAlignment), uint64(len(strData)))
	if err != nil || ptr == 0 {
		return 0, 0, fmt.Errorf("failed to allocate memory for string: %w", err)
	}
	ok := opts.Memory.Write(ptr, strData)
	if !ok {
		return 0, 0, fmt.Errorf("failed to write string data to memory at %d", ptr)
	}
	return ptr, codeUnits, nil
}
