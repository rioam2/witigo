package codegen

import (
	"github.com/golang-cz/textcase"
	"github.com/moznion/gowrtr/generator"
	"github.com/rioam2/witigo/pkg/wit"
)

func GenerateFromFunction(w wit.WitFunction, receiver *generator.FuncReceiver) *generator.Func {
	fn := generator.NewFunc(
		receiver,
		GenerateSignatureFromFunction(w),
	)
	fn = fn.AddStatements(generator.NewRawStatementf("var args []uint32"))
	fn = fn.AddStatements(generator.NewRawStatementf("var result %s", GenerateTypenameFromType(w.Returns())))
	for idx, param := range w.Params() {
		fn = fn.AddStatements(
			generator.NewRawStatementf("arg%02dArgs, arg%02dFree, err := abi.WriteParameter(i.abiOpts, %s)", idx, idx, textcase.CamelCase(param.Name())),
			generator.NewRawStatementf("defer arg%02dFree()", idx),
			generator.NewRawStatementf("if err != nil {"),
			generator.NewRawStatementf("  return result, fmt.Errorf(\"failed to write parameter %d: %%w\", err)", idx),
			generator.NewRawStatementf("}"),
			generator.NewRawStatementf("args = append(args, arg%02dArgs...)", idx),
		)
	}
	fn = fn.AddStatements(
		generator.NewRawStatementf("ret, postReturn, err := abi.Call(i.abiOpts, \"%s\", args...)", textcase.KebabCase(w.Name())),
		generator.NewRawStatementf("if err != nil {"),
		generator.NewRawStatementf("  return result, fmt.Errorf(\"failed to call %s: %%w\", err)", textcase.KebabCase(w.Name())),
		generator.NewRawStatementf("}"),
		generator.NewRawStatementf("defer postReturn()"),
		generator.NewRawStatementf("err = abi.Read(i.abiOpts, ret, &result)"),
		generator.NewRawStatementf("if err != nil {"),
		generator.NewRawStatementf("  return result, fmt.Errorf(\"failed to read result: %%w\", err)"),
		generator.NewRawStatementf("}"),
	)
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
