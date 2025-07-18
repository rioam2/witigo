package codegen

import (
	"github.com/rioam2/witigo/pkg/wit"
)

func GenerateFromFile(path string) (string, error) {
	witDefinition, err := wit.NewFromFile(path)
	if err != nil {
		return "", err
	}
	codeGen := GenerateFromWorld(witDefinition.Worlds()[0], witDefinition.Name())
	return codeGen.EnableSyntaxChecking().Gofmt().Generate(0)
}
