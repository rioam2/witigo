package codegen

import (
	"github.com/golang-cz/textcase"
	"github.com/moznion/gowrtr/generator"
	"github.com/rioam2/witigo/pkg/wit"
)

func GenerateFromFunction(w wit.WitFunction) *generator.Func {
	parameters := make([]*generator.FuncParameter, len(w.Params()))
	for idx, param := range w.Params() {
		parameters[idx] = generator.NewFuncParameter(
			textcase.CamelCase(param.Name()),
			GenerateTypenameFromType(param.Type()),
		)
	}
	fn := generator.NewFunc(
		nil,
		generator.NewFuncSignature(textcase.PascalCase(w.Name())).
			AddParameters(parameters...).
			AddReturnTypes(GenerateTypenameFromType(w.Returns())),
	)
	fn = fn.AddStatements(
		generator.NewRawStatement("// TODO: Implement function body"),
		generator.NewRawStatementf("var result %s", GenerateTypenameFromType(w.Returns())),
		generator.NewRawStatement("return result"),
	)
	return fn
}
