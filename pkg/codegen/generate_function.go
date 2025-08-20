package codegen

import (
	"github.com/golang-cz/textcase"
	"github.com/moznion/gowrtr/generator"
	witigo "github.com/rioam2/witigo/pkg"
	"github.com/rioam2/witigo/pkg/wit"
)

func GenerateFromFunction(w wit.WitFunction, receiver *generator.FuncReceiver) *generator.Func {
	fn := generator.NewFunc(
		receiver,
		GenerateSignatureFromFunction(w),
	)
	fn = fn.AddStatements(generator.NewRawStatementf("var args []uint64"))
	fn = fn.AddStatements(generator.NewRawStatementf("var result %s", GenerateTypenameFromType(w.Returns())))
	for idx, param := range w.Params() {
		switch param.Type().Kind() {
		case witigo.AbiTypeString:
			fn = fn.AddStatements(
				generator.NewRawStatementf("arg%02dPtr, arg%02dFree, err := abi.Write(i.abiOpts, %s, nil)", idx, idx, textcase.CamelCase(param.Name())),
				generator.NewRawStatementf("defer arg%02dFree()", idx),
				generator.NewRawStatementf("if err != nil {"),
				generator.NewRawStatementf("  return \"\", fmt.Errorf(\"failed to write string: %%w\", err)"),
				generator.NewRawStatementf("}"),
				generator.NewRawStatementf("arg%02dStringDataPtr, _ := i.memory.ReadUint32Le(arg%02dPtr)", idx, idx),
				generator.NewRawStatementf("arg%02dStringLength, _ := i.memory.ReadUint32Le(arg%02dPtr + 4)", idx, idx),
				generator.NewRawStatementf("args = append(args, uint64(arg%02dStringDataPtr), uint64(arg%02dStringLength))", idx, idx),
			)
		default:
			fn = fn.AddStatements(
				generator.NewRawStatementf("_ = args"),
			)
		}
	}
	switch w.Returns().Kind() {
	case witigo.AbiTypeString:
		fn = fn.AddStatements(
			generator.NewRawStatementf("ret, postReturn, err := abi.Call(i.abiOpts, \"%s\", args...)", textcase.KebabCase(w.Name())),
			generator.NewRawStatementf("if err != nil {"),
			generator.NewRawStatementf("  panic(fmt.Errorf(\"failed to call %s: %%w\", err))", textcase.KebabCase(w.Name())),
			generator.NewRawStatementf("}"),
			generator.NewRawStatementf("defer postReturn()"),
			generator.NewRawStatementf("err = abi.ReadString(i.abiOpts, ret, &result)"),
			generator.NewRawStatementf("if err != nil {"),
			generator.NewRawStatementf("  panic(fmt.Errorf(\"failed to read string result: %%w\", err))"),
			generator.NewRawStatementf("}"),
		)
	}
	fn = fn.AddStatements(generator.NewRawStatement("return result, nil"))
	return fn
}

func GenerateSignatureFromFunction(w wit.WitFunction) *generator.FuncSignature {
	parameters := make([]*generator.FuncParameter, len(w.Params()))
	for idx, param := range w.Params() {
		parameters[idx] = generator.NewFuncParameter(
			textcase.CamelCase(param.Name()),
			GenerateTypenameFromType(param.Type()),
		)
	}
	return generator.NewFuncSignature(textcase.PascalCase(w.Name())).
		AddParameters(parameters...).
		AddReturnTypes(GenerateTypenameFromType(w.Returns()), "error")
}
