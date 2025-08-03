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
				generator.NewRawStatementf("arg%02dPtr, arg%02dUnits, err := abi.WriteString(i.abiOpts, %s)", idx, idx, textcase.CamelCase(param.Name())),
				generator.NewRawStatementf("if err != nil {"),
				generator.NewRawStatementf("  panic(fmt.Errorf(\"failed to write string: %%w\", err))"),
				generator.NewRawStatementf("}"),
				generator.NewRawStatementf("args = append(args, uint64(arg%02dPtr), uint64(arg%02dUnits))", idx, idx),
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
			generator.NewRawStatementf("ret, err := abi.Call(i.abiOpts, \"%s\", args...)", textcase.KebabCase(w.Name())),
			generator.NewRawStatementf("if err != nil {"),
			generator.NewRawStatementf("  panic(fmt.Errorf(\"failed to call %s: %%w\", err))", textcase.KebabCase(w.Name())),
			generator.NewRawStatementf("}"),
			generator.NewRawStatementf("result, err = abi.ReadString(i.abiOpts, ret)"),
			generator.NewRawStatementf("if err != nil {"),
			generator.NewRawStatementf("  panic(fmt.Errorf(\"failed to read string result: %%w\", err))"),
			generator.NewRawStatementf("}"),
		)
	}
	fn = fn.AddStatements(generator.NewRawStatement("return result"))
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
		AddReturnTypes(GenerateTypenameFromType(w.Returns()))
}
