package gen_application

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path"

	"github.com/lujingwei002/gira/proj"
)

type golang_parser struct {
}

func (p *golang_parser) parseApplicationStruct(state *gen_state, s *ast.StructType) error {
	for _, field := range s.Fields.List {
		state.applications = append(state.applications, application{
			ApplicationName: field.Names[0].Name,
			ModuleName:      proj.Config.Module,
		})
	}
	return nil
}

func (p *golang_parser) parse(state *gen_state) (err error) {
	applicationFilePath := path.Join(proj.Config.SrcDocDir, "application.go")
	fset := token.NewFileSet()
	var f *ast.File
	f, err = parser.ParseFile(fset, applicationFilePath, nil, parser.ParseComments)
	if err != nil {
		fmt.Println(err)
		return
	}
	ast.Inspect(f, func(n ast.Node) bool {
		if err != nil {
			return false
		}
		switch x := n.(type) {
		case *ast.TypeSpec:
			if x.Name.Name == "Application" {
				if err = p.parseApplicationStruct(state, x.Type.(*ast.StructType)); err != nil {
					return false
				}
			}
		}
		return true
	})
	// ast.Print(fset, f)
	return
}
