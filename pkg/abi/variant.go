package abi

import (
	"errors"
	"fmt"
	"reflect"
)

// ReadVariant reads a variant from memory at the specified pointer into the result.
// Memory layout (current implementation):
//
//	[ discriminant (size=SizeOf(Type field)) ][ payload (active case only, aligned) ]
//
// The overall allocation size equals discriminant size + max(size(case_i)) rounded up
// to the maximum alignment of all fields. This mirrors a tagged union representation.
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
	// Empty struct payload -> nothing further to read
	if activeField.Kind() == reflect.Struct && activeField.Type().NumField() == 0 {
		return nil
	}

	// Canonical ABI: the payload area is placed at the discriminant size rounded up to the
	// maximum alignment of all case payload types (not the active case's alignment). The
	// overall allocation size already accounts for this (see SizeOf). Previously we aligned
	// to only the active field alignment which caused smaller-alignment cases (e.g. list/
	// string with alignment 4) to be read 4 bytes too early when a larger-alignment case
	// (e.g. a record with alignment 8) existed. This manifested as corrupt list headers and
	// huge lengths for the bytes arm of complex variants.
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

	discrField := rv.Field(0)
	discrSize := SizeOf(discrField.Interface())
	discrAlign := AlignmentOf(discrField.Interface())
	discrPtr := AlignTo(ptr, discrAlign)

	// Serialize discriminant bytes
	bytes := make([]byte, discrSize)
	if discrField.CanUint() {
		v := discrField.Uint()
		for i := range bytes {
			bytes[i] = byte(v >> (8 * uint(i)))
		}
	} else if discrField.CanInt() {
		v := discrField.Int()
		for i := range bytes {
			bytes[i] = byte(v >> (8 * uint(i)))
		}
	} else {
		return ptr, free, fmt.Errorf("discriminant field type %s not int/uint", discrField.Type())
	}
	if !opts.Memory.Write(discrPtr, bytes) {
		return ptr, free, fmt.Errorf("failed to write discriminant at %d", discrPtr)
	}

	discrVal := uint64(0)
	if discrField.CanUint() {
		discrVal = discrField.Uint()
	} else {
		discrVal = uint64(discrField.Int())
	}
	caseIndex := int(discrVal)
	numCases := rv.NumField() - 1
	if caseIndex < 0 || caseIndex >= numCases {
		return ptr, free, fmt.Errorf("variant discriminant %d out of range [0,%d)", caseIndex, numCases)
	}
	activeField := rv.Field(caseIndex + 1)
	if activeField.Kind() == reflect.Struct && activeField.Type().NumField() == 0 {
		return ptr, free, nil // empty payload
	}
	// Canonical ABI payload placement: align to max alignment across all cases.
	valuePtr := AlignTo(discrPtr+discrSize, maxVariantAlignment(rv))
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

	discrField := rv.Field(0)
	discrUint := uint64(0)
	if discrField.CanUint() {
		discrUint = discrField.Uint()
	} else if discrField.CanInt() {
		discrUint = uint64(discrField.Int())
	} else {
		return params, free, fmt.Errorf("discriminant field type %s not int/uint", discrField.Type())
	}
	caseIndex := int(discrUint)
	numCases := rv.NumField() - 1
	if caseIndex < 0 || caseIndex >= numCases {
		return params, free, fmt.Errorf("variant discriminant %d out of range [0,%d)", caseIndex, numCases)
	}

	// Append discriminant first
	discrParam := Parameter{Value: discrUint, Size: SizeOf(discrField.Interface()), Alignment: AlignmentOf(discrField.Interface())}
	params = append(params, discrParam)

	// Build unified payload shape across all cases (flatten_variant logic approximation)
	type slot struct{ size, align uint64 }
	slots := []slot{}
	for i := 1; i < rv.NumField(); i++ {
		f := rv.Field(i)
		if f.Kind() == reflect.Struct && f.Type().NumField() == 0 {
			continue // empty payload contributes nothing
		}
		zeroVal := reflect.New(f.Type()).Elem().Interface()
		fp, fpFree, e := WriteParameter(opts, zeroVal)
		freeCallbacks = append(freeCallbacks, fpFree)
		if e != nil {
			return params, free, e
		}
		for si, p := range fp {
			if si >= len(slots) {
				slots = append(slots, slot{size: p.Size, align: p.Alignment})
			} else {
				if p.Size > slots[si].size {
					slots[si].size = p.Size
				}
				if p.Alignment > slots[si].align {
					slots[si].align = p.Alignment
				}
			}
		}
	}

	// Real params for active case
	activeField := rv.Field(caseIndex + 1)
	realParams := []Parameter{}
	if !(activeField.Kind() == reflect.Struct && activeField.Type().NumField() == 0) {
		rp, rpFree, e := WriteParameter(opts, activeField.Interface())
		freeCallbacks = append(freeCallbacks, rpFree)
		if e != nil {
			return params, free, fmt.Errorf("failed to write variant payload parameters: %w", e)
		}
		realParams = rp
	}

	for i, s := range slots {
		if i < len(realParams) {
			p := realParams[i]
			p.Size = s.size
			p.Alignment = s.align
			params = append(params, p)
		} else {
			params = append(params, Parameter{Value: 0, Size: s.size, Alignment: s.align})
		}
	}
	return params, free, nil
}
