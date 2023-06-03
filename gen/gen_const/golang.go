package gen_const

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path"

	"github.com/lujingwei002/gira"
	"github.com/lujingwei002/gira/gen"
	"github.com/lujingwei002/gira/proj"
)

type golang_parser struct {
}

func (p *golang_parser) parseConstStruct(state *const_state, s *ast.StructType) error {
	for _, field := range s.Fields.List {
		if annotates, err := gen.ExtraAnnotate(field.Doc.List); err != nil {
			return err
		} else {
			name := field.Names[0].Name
			var filePath string
			var keyArr []string
			if v, ok := annotates["excel"]; !ok {
				return gira.ErrTodo.Trace()
			} else {
				filePath = v.Values[0]
			}
			if v, ok := annotates["keys"]; !ok {
				return gira.ErrTodo.Trace()
			} else {
				keyArr = v.Values
			}
			descriptor := &Descriptor{
				Name:     name,
				filePath: filePath,
				keyArr:   keyArr,
			}
			state.constFile.descriptorDict[name] = descriptor
		}
	}
	return nil
}

func (p *golang_parser) parse(state *const_state) (err error) {
	filePath := path.Join(proj.Config.SrcDocDir, "const.go")
	fset := token.NewFileSet()
	var f *ast.File
	f, err = parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return
	}
	ast.Inspect(f, func(n ast.Node) bool {
		if err != nil {
			return false
		}
		switch x := n.(type) {
		case *ast.TypeSpec:
			if x.Name.Name == "Const" {
				if err = p.parseConstStruct(state, x.Type.(*ast.StructType)); err != nil {
					return false
				}
			}
		}
		return true
	})
	// ast.Print(fset, f)
	return
}
