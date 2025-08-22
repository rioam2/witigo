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
	fn = fn.AddStatements(generator.NewRawStatementf("var args []uint64"))
	if w.Returns() != nil {
		fn = fn.AddStatements(generator.NewRawStatementf("var result %s", GenerateTypenameFromType(w.Returns())))
	}
	for idx, param := range w.Params() {
		fn = fn.AddStatements(
			generator.NewRawStatementf("arg%02dArgs, arg%02dFree, err := abi.WriteParameter(i.abiOpts, %s)", idx, idx, textcase.CamelCase(param.Name())),
			generator.NewRawStatementf("defer arg%02dFree()", idx),
			generator.NewRawStatementf("args = append(args, arg%02dArgs...)", idx),
		)
		if w.Returns() == nil {
			fn = fn.AddStatements(
				generator.NewRawStatementf("if err != nil {"),
				generator.NewRawStatementf("  return fmt.Errorf(\"failed to write parameter %d: %%w\", err)", idx),
				generator.NewRawStatementf("}"),
			)
		} else {
			fn = fn.AddStatements(
				generator.NewRawStatementf("if err != nil {"),
				generator.NewRawStatementf("  return result, fmt.Errorf(\"failed to write parameter %d: %%w\", err)", idx),
				generator.NewRawStatementf("}"),
			)
		}
	}
	if w.Returns() == nil {
		fn = fn.AddStatements(
			generator.NewRawStatementf("_, postReturn, err := abi.Call(i.abiOpts, \"%s\", args...)", textcase.KebabCase(w.Name())),
			generator.NewRawStatementf("if err != nil {"),
			generator.NewRawStatementf("  return fmt.Errorf(\"failed to call %s: %%w\", err)", textcase.KebabCase(w.Name())),
			generator.NewRawStatementf("}"),
		)
	} else {
		fn = fn.AddStatements(
			generator.NewRawStatementf("ret, postReturn, err := abi.Call(i.abiOpts, \"%s\", args...)", textcase.KebabCase(w.Name())),
			generator.NewRawStatementf("if err != nil {"),
			generator.NewRawStatementf("  return result, fmt.Errorf(\"failed to call %s: %%w\", err)", textcase.KebabCase(w.Name())),
			generator.NewRawStatementf("}"),
		)
	}
	fn = fn.AddStatements(
		generator.NewRawStatementf("defer postReturn()"),
	)
	if w.Returns() != nil {
		fn = fn.AddStatements(
			generator.NewRawStatementf("err = abi.Read(i.abiOpts, ret, &result)"),
			generator.NewRawStatementf("if err != nil {"),
			generator.NewRawStatementf("  return result, fmt.Errorf(\"failed to read result: %%w\", err)"),
			generator.NewRawStatementf("}"),
		)
	}
	if w.Returns() == nil {
		fn = fn.AddStatements(generator.NewRawStatement("return nil"))
	} else {
		fn = fn.AddStatements(generator.NewRawStatement("return result, nil"))
	}
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
	signature := generator.NewFuncSignature(textcase.PascalCase(w.Name())).
		AddParameters(parameters...)
	if w.Returns() != nil {
		signature = signature.AddReturnTypes(GenerateTypenameFromType(w.Returns()))
	}
	signature = signature.AddReturnTypes("error")
	return signature
}
