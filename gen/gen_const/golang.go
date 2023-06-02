package gen_const

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path"
	"regexp"
	"strings"

	"github.com/lujingwei002/gira"
	"github.com/lujingwei002/gira/proj"
)

type golang_parser struct {
}

func (p *golang_parser) extraTag(str string) map[string]string {
	re := regexp.MustCompile(`(\w+):"([^"]*)"`)
	matches := re.FindAllStringSubmatch(str, -1)

	values := make(map[string]string)
	for _, match := range matches {
		if len(match) >= 3 {
			values[match[1]] = match[2]
		}
	}
	return values
}

func (p *golang_parser) extraAnnotate(comments []*ast.Comment) map[string]string {
	re := regexp.MustCompile(`// @(\w+)\((.*)\)`)

	values := make(map[string]string)
	for _, comment := range comments {
		if !strings.HasPrefix(comment.Text, "// @") {
			continue
		}
		fmt.Println(comment.Text)

		matches := re.FindAllStringSubmatch(comment.Text, -1)
		for _, match := range matches {
			fmt.Println(match)

		}
	}

	return values
}

func (p *golang_parser) parseConstStruct(state *const_state, s *ast.StructType) error {
	for _, field := range s.Fields.List {
		p.extraAnnotate(field.Doc.List)
		name := field.Names[0].Name
		tags := p.extraTag(field.Tag.Value)
		var filePath string
		var keyArr []string
		if v, ok := tags["xlsx"]; !ok {
			return gira.ErrTodo.Trace()
		} else {
			filePath = v
		}
		if v, ok := tags["keys"]; !ok {
			return gira.ErrTodo.Trace()
		} else {
			keyArr = strings.Split(v, ",")
		}
		descriptor := &Descriptor{
			Name:     name,
			filePath: filePath,
			keyArr:   keyArr,
		}
		state.constFile.descriptorDict[name] = descriptor
	}
	return nil
}

func (p *golang_parser) parse(state *const_state) error {
	filePath := path.Join(proj.Config.DocDir, "const.go")
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		fmt.Println(err)
		return err
	}
	ast.Inspect(f, func(n ast.Node) bool {
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
	return nil
}
