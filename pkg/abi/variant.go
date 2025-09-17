package abi

import (
	"errors"
	"fmt"
	"reflect"
)

// ReadVariant reads a variant from memory at the specified pointer into the result.
func ReadVariant(opts AbiOptions, ptr uint64, result any) error {
	rv := reflect.ValueOf(result)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return errors.New("must pass a non-nil pointer result")
	}
	rv = rv.Elem()
	if !isStructVariantType(rv) {
		return fmt.Errorf("result must be a variant pointer, got %s", rv.Type().Name())
	}
	if rv.NumField() == 0 {
		return errors.New("variant struct missing fields")
	}

	// Read discriminant into first field
	discriminantField := rv.Field(0)
	discriminantSize := SizeOf(discriminantField.Interface())
	discriminantAlign := AlignmentOf(discriminantField.Interface())
	discriminantPtr := AlignTo(ptr, discriminantAlign)

	// allocate a temporary variable to read into then set (so we invoke int logic)
	tmpPtrVal := reflect.New(discriminantField.Type())
	if err := ReadInt(opts, discriminantPtr, tmpPtrVal.Interface()); err != nil {
		return fmt.Errorf("failed to read variant discriminant: %w", err)
	}
	discriminantField.Set(tmpPtrVal.Elem())

	discrVal := uint64(0)
	if discriminantField.CanUint() {
		discrVal = discriminantField.Uint()
	} else if discriminantField.CanInt() {
		discrVal = uint64(discriminantField.Int())
	} else {
		return fmt.Errorf("discriminant field type %s not int/uint", discriminantField.Type())
	}

	caseIndex := int(discrVal)
	numCases := rv.NumField() - 1
	if caseIndex < 0 || caseIndex >= numCases { // invalid discriminant
		return fmt.Errorf("variant discriminant %d out of range [0,%d)", caseIndex, numCases)
	}

	// Active case field is offset +1 from Type field
	activeField := rv.Field(caseIndex + 1)

	// Empty variant case (struct{}) has no payload to read
	if isAnonymousEmptyStruct(activeField) {
		return nil
	}

	valuePtr := AlignTo(discriminantPtr+discriminantSize, maxVariantAlignment(rv))
	return Read(opts, valuePtr, activeField.Addr().Interface())
}

// WriteVariant writes a variant value to linear memory and returns the pointer & free callback.
func WriteVariant(opts AbiOptions, value any, ptrHint *uint64) (ptr uint64, free AbiFreeCallback, err error) {
	ptr = 0
	freeCallbacks := []AbiFreeCallback{}
	free = wrapFreeCallbacks(&freeCallbacks)

	rv := reflect.ValueOf(value)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if !rv.IsValid() {
		return ptr, free, errors.New("must pass a valid variant value")
	}
	if !isStructVariantType(rv) {
		return ptr, free, fmt.Errorf("value must be a variant, got %s", rv.Kind())
	}
	if rv.NumField() == 0 {
		return ptr, free, errors.New("variant struct missing fields")
	}

	size := SizeOf(value)
	alignment := AlignmentOf(value)
	if ptrHint != nil && *ptrHint != 0 {
		ptr = AlignTo(*ptrHint, alignment)
	} else {
		var freeVar AbiFreeCallback
		ptr, freeVar, err = abiMalloc(opts, size, alignment)
		if err != nil {
			return ptr, free, err
		}
		freeCallbacks = append(freeCallbacks, freeVar)
	}

	discriminantField := rv.Field(0)
	discriminantSize := SizeOf(discriminantField.Interface())
	discriminantAlign := AlignmentOf(discriminantField.Interface())
	discriminantPtr := AlignTo(ptr, discriminantAlign)

	// Serialize discriminant bytes into linear memory
	bytes := make([]byte, discriminantSize)
	if discriminantField.CanUint() {
		v := discriminantField.Uint()
		for i := range bytes {
			bytes[i] = byte(v >> (8 * uint(i)))
		}
	} else if discriminantField.CanInt() {
		v := discriminantField.Int()
		for i := range bytes {
			bytes[i] = byte(v >> (8 * uint(i)))
		}
	} else {
		return ptr, free, fmt.Errorf("discriminant field type %s not int/uint", discriminantField.Type())
	}
	if !opts.Memory.Write(discriminantPtr, bytes) {
		return ptr, free, fmt.Errorf("failed to write discriminant at %d", discriminantPtr)
	}

	discriminantVal := uint64(0)
	if discriminantField.CanUint() {
		discriminantVal = discriminantField.Uint()
	} else {
		discriminantVal = uint64(discriminantField.Int())
	}
	caseIndex := int(discriminantVal)
	numCases := rv.NumField() - 1
	if caseIndex < 0 || caseIndex >= numCases {
		return ptr, free, fmt.Errorf("variant discriminant %d out of range [0,%d)", caseIndex, numCases)
	}
	activeField := rv.Field(caseIndex + 1)
	if isAnonymousEmptyStruct(activeField) {
		return ptr, free, nil
	}
	valuePtr := AlignTo(discriminantPtr+discriminantSize, maxVariantAlignment(rv))
	_, valueFree, err := Write(opts, activeField.Interface(), &valuePtr)
	freeCallbacks = append(freeCallbacks, valueFree)
	if err != nil {
		return ptr, free, fmt.Errorf("failed to write variant payload: %w", err)
	}
	return ptr, free, nil
}

// WriteParameterVariant flattens the discriminant plus the active case payload parameters.
func WriteParameterVariant(opts AbiOptions, value any) (params []Parameter, free AbiFreeCallback, err error) {
	params = []Parameter{}
	freeCallbacks := []AbiFreeCallback{}
	free = wrapFreeCallbacks(&freeCallbacks)

	rv := reflect.ValueOf(value)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if !rv.IsValid() {
		return params, free, errors.New("must pass a valid variant value")
	}
	if !isStructVariantType(rv) {
		return params, free, fmt.Errorf("value must be a variant, got %s", rv.Kind())
	}
	if rv.NumField() == 0 {
		return params, free, errors.New("variant struct missing fields")
	}

	discriminantField := rv.Field(0)
	discriminantUint := uint64(0)
	if discriminantField.CanUint() {
		discriminantUint = discriminantField.Uint()
	} else if discriminantField.CanInt() {
		discriminantUint = uint64(discriminantField.Int())
	} else {
		return params, free, fmt.Errorf("discriminant field type %s not int/uint", discriminantField.Type())
	}
	caseIndex := int(discriminantUint)
	numCases := rv.NumField() - 1
	if caseIndex < 0 || caseIndex >= numCases {
		return params, free, fmt.Errorf("variant discriminant %d out of range [0,%d)", caseIndex, numCases)
	}

	// Append discriminant first
	discriminantParam := Parameter{
		Value:     discriminantUint,
		Size:      SizeOf(discriminantField.Interface()),
		Alignment: AlignmentOf(discriminantField.Interface()),
	}
	params = append(params, discriminantParam)

	// Build unified payload shape across all cases (flatten_variant logic approximation)
	type slot struct {
		size  uint64
		align uint64
	}
	slots := []slot{}
	for fieldIndex := 1; fieldIndex < rv.NumField(); fieldIndex++ {
		field := rv.Field(fieldIndex)
		if isAnonymousEmptyStruct(field) {
			continue // empty payload contributes nothing
		}
		zeroVal := reflect.New(field.Type()).Elem().Interface()
		fieldParams, fieldFree, e := WriteParameter(opts, zeroVal)
		freeCallbacks = append(freeCallbacks, fieldFree)
		if e != nil {
			return params, free, e
		}
		for slotIndex, param := range fieldParams {
			if slotIndex >= len(slots) {
				slots = append(slots, slot{size: param.Size, align: param.Alignment})
			} else {
				if param.Size > slots[slotIndex].size {
					slots[slotIndex].size = param.Size
				}
				if param.Alignment > slots[slotIndex].align {
					slots[slotIndex].align = param.Alignment
				}
			}
		}
	}

	// Real params for active case
	activeField := rv.Field(caseIndex + 1)
	realParams := []Parameter{}
	if !isAnonymousEmptyStruct(activeField) {
		activeFieldParams, activeFieldFree, e := WriteParameter(opts, activeField.Interface())
		freeCallbacks = append(freeCallbacks, activeFieldFree)
		if e != nil {
			return params, free, fmt.Errorf("failed to write variant payload parameters: %w", e)
		}
		realParams = activeFieldParams
	}

	for slotIndex, slot := range slots {
		if slotIndex < len(realParams) {
			p := realParams[slotIndex]
			p.Size = slot.size
			p.Alignment = slot.align
			params = append(params, p)
		} else {
			params = append(params, Parameter{
				Value:     0,
				Size:      slot.size,
				Alignment: slot.align,
			})
		}
	}

	return params, free, nil
}
