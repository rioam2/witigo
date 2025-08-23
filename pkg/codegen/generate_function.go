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
	fn = fn.AddStatements(generator.NewRawStatementf("var flatParams []uint64"))
	fn = fn.AddStatements(generator.NewRawStatementf("var params []abi.Parameter"))
	if w.Returns() != nil {
		fn = fn.AddStatements(generator.NewRawStatementf("var result %s", GenerateTypenameFromType(w.Returns())))
	}
	for idx, param := range w.Params() {
		fn = fn.AddStatements(
			generator.NewRawStatementf("arg%02dArgs, arg%02dFree, err := abi.WriteParameter(i.abiOpts, %s)", idx, idx, textcase.CamelCase(param.Name())),
			generator.NewRawStatementf("defer arg%02dFree()", idx),
			generator.NewRawStatementf("params = append(params, arg%02dArgs...)", idx),
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
	fn = fn.AddStatements(
		generator.NewRawStatement("if len(params) > abi.MAX_FLAT_PARAMS {"),
		generator.NewRawStatement("  flatParam, flatParamFree, err := abi.WriteIndirectParameters(i.abiOpts, params...)"),
	)
	if w.Returns() == nil {
		fn = fn.AddStatements(
			generator.NewRawStatement("  if err != nil {"),
			generator.NewRawStatement("    return fmt.Errorf(\"failed to write indirect parameters: %w\", err)"),
			generator.NewRawStatement("  }"),
		)
	} else {
		fn = fn.AddStatements(
			generator.NewRawStatement("  if err != nil {"),
			generator.NewRawStatement("    return result, fmt.Errorf(\"failed to write indirect parameters: %w\", err)"),
			generator.NewRawStatement("  }"),
		)
	}
	fn = fn.AddStatements(
		generator.NewRawStatement("  flatParams = append(flatParams, flatParam)"),
		generator.NewRawStatement("  defer flatParamFree()"),
		generator.NewRawStatement("} else {"),
		generator.NewRawStatement("  flatParams = make([]uint64, len(params))"),
		generator.NewRawStatement("  for i := range params {"),
		generator.NewRawStatement("    flatParams[i] = params[i].Value"),
		generator.NewRawStatement("  }"),
		generator.NewRawStatement("}"),
	)
	if w.Returns() == nil {
		fn = fn.AddStatements(
			generator.NewRawStatementf("_, postReturn, err := abi.Call(i.abiOpts, \"%s\", flatParams...)", textcase.KebabCase(w.Name())),
			generator.NewRawStatementf("if err != nil {"),
			generator.NewRawStatementf("  return fmt.Errorf(\"failed to call %s: %%w\", err)", textcase.KebabCase(w.Name())),
			generator.NewRawStatementf("}"),
		)
	} else {
		fn = fn.AddStatements(
			generator.NewRawStatementf("ret, postReturn, err := abi.Call(i.abiOpts, \"%s\", flatParams...)", textcase.KebabCase(w.Name())),
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
