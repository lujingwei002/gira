package gen_behavior

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/lujingwei002/gira/gen"
	"github.com/lujingwei002/gira/log"
	"github.com/lujingwei002/gira/proj"
)

type golang_parser struct {
}

func (p *golang_parser) parse(state *gen_state) error {
	filePathArr := make([]string, 0)
	if err := filepath.WalkDir(path.Join(proj.Config.DocDir, "behavior"), func(path string, d os.DirEntry, err error) error {
		if d == nil {
			return nil
		}
		if d.IsDir() {
			return nil
		}
		if filepath.Ext(d.Name()) == ".go" {
			filePathArr = append(filePathArr, path)
		}
		return nil
	}); err != nil {
		return err
	}
	sort.Strings(filePathArr)
	for _, filePath := range filePathArr {
		if err := p.parseFile(state, filePath); err != nil {
			return err
		}
	}
	return nil
}

func (p *golang_parser) reportError(fset *token.FileSet, filePath string, pos token.Pos, msg string) error {
	return fmt.Errorf("%s %v %s", filePath, fset.Position(pos), msg)
}

func (p *golang_parser) parseFile(state *gen_state, filePath string) (err error) {
	log.Info("处理文件", filePath)
	fileName := path.Base(filePath)
	dbName := strings.Replace(fileName, ".go", "", 1)
	database := &Database{
		Module:                proj.Config.Module,
		CollectionArr:         make([]*Collection, 0),
		GenModelDir:           path.Join(proj.Config.SrcGenBehaviorDir, dbName),
		GenModelFilePath:      path.Join(proj.Config.SrcGenBehaviorDir, dbName, fmt.Sprintf("%s.gen.go", dbName)),
		GenBinFilePath:        path.Join(proj.Config.SrcGenBehaviorDir, dbName, "bin", fmt.Sprintf("%s.gen.go", dbName)),
		GenBinDir:             path.Join(proj.Config.SrcGenBehaviorDir, dbName, "bin"),
		DbName:                dbName,
		DbStructName:          camelString(dbName),
		MongoDriverStructName: fmt.Sprintf("%sMongoDriver", camelString(dbName)),
		DriverInterfaceName:   fmt.Sprintf("%sDriver", camelString(dbName)),
	}
	fset := token.NewFileSet()
	var f *ast.File
	var fileContent []byte
	fileContent, err = ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	f, err = parser.ParseFile(fset, "", fileContent, parser.ParseComments)
	if err != nil {
		return
	}
	ast.Inspect(f, func(n ast.Node) bool {
		if err != nil {
			return false
		}
		switch x := n.(type) {
		case *ast.GenDecl:
			for _, spec := range x.Specs {
				switch s := spec.(type) {
				case *ast.ValueSpec:
					if s.Names[0].Name == "Driver" {
						// database.Driver = v.(string)
						if values0, ok := s.Values[0].(*ast.CompositeLit); !ok {
							err = p.reportError(fset, filePath, s.Pos(), "driver必须是string数组")
							return false
						} else if typ, ok := values0.Type.(*ast.ArrayType); !ok {
							err = p.reportError(fset, filePath, values0.Pos(), "driver必须是string数组")
							return false
						} else if elt, ok := typ.Elt.(*ast.Ident); !ok {
							err = p.reportError(fset, filePath, values0.Pos(), "driver必须是string数组")
							return false
						} else if elt.Name != "string" {
							err = p.reportError(fset, filePath, values0.Pos(), "driver必须是string数组")
							return false
						} else if len(values0.Elts) <= 0 {
							err = p.reportError(fset, filePath, values0.Pos(), "至少一个元素")
							return false
						} else {
							if v, ok := values0.Elts[0].(*ast.BasicLit); !ok {
								err = p.reportError(fset, filePath, values0.Pos(), "至少一个元素")
								return false
							} else {
								database.Driver = v.Value
							}
						}
					}
				case *ast.TypeSpec:
					var comment string
					if x.Doc != nil {
						// if annotates, err := gen.ExtraAnnotate(x.Doc.List); err == nil {
						// 	log.Println("ccccc", annotates)
						// }
						if comments, err := gen.ExtraComment(x.Doc); err == nil {
							comment = strings.Join(comments, ",")
						}
					}
					if typ, ok := s.Type.(*ast.StructType); ok {
						var coll *Collection
						if coll, err = p.parseCollection(state, database, fileContent, s.Name.Name, fset, typ); err != nil {
							return false
						} else {
							coll.Comment = comment
							database.CollectionArr = append(database.CollectionArr, coll)
						}
					}
				}
			}
		}
		return true
	})
	// ast.Print(fset, f)
	if err != nil {
		return err
	}
	// for _, coll := range database.CollectionArr {
	// 	log.Println(coll.CollName)
	// 	for _, v := range coll.FieldArr {
	// 		log.Println("==", v.CamelName, v.GoTypeName)
	// 	}
	// 	for _, v := range coll.IndexArr {
	// 		log.Println("==", v.FullName)
	// 	}
	// }
	state.databaseArr = append(state.databaseArr, database)
	return nil
}

func (p *golang_parser) parseCollection(state *gen_state, database *Database, fileContent []byte, collName string, fset *token.FileSet, s *ast.StructType) (*Collection, error) {
	coll := &Collection{
		MongoDriverStructName: database.MongoDriverStructName,
		CollName:              collName,
		StructName:            camelString(collName),
		MongoDaoStructName:    fmt.Sprintf("%sMongoDao", camelString(collName)),
	}
	for _, f := range s.Fields.List {
		if f.Names[0].Name == "Model" {
			if model, ok := f.Type.(*ast.StructType); ok {
				if err := p.parseStruct(state, coll, fileContent, fset, model); err != nil {
					return nil, err
				}
			}
		} else if f.Names[0].Name == "Index" {
			if index, ok := f.Type.(*ast.StructType); ok {
				if err := p.parseIndex(state, coll, fileContent, fset, index); err != nil {
					return nil, err
				}
			}
		}
	}
	return coll, nil
}

func (p *golang_parser) parseStruct(state *gen_state, coll *Collection, fileContent []byte, fset *token.FileSet, s *ast.StructType) error {
	coll.FieldDict = make(map[string]*Field)
	coll.FieldArr = make([]*Field, 0)
	for _, f := range s.Fields.List {
		var tag int
		var err error
		var fieldName string
		var typeStr string
		var tagStr string
		tagStr = f.Tag.Value
		tagStr = strings.Replace(tagStr, "`", "", 2)
		if tag, err = strconv.Atoi(tagStr); err != nil {
			return err
		}

		// ast.Print(fset, s)
		fieldName = f.Names[0].Name
		typeStr = string(fileContent[f.Type.Pos()-1 : f.Type.End()-1])
		field := &Field{
			Coll:      coll,
			Name:      fieldName,
			CamelName: camelString(fieldName),
			Tag:       tag,
		}
		if f.Doc != nil {
			if comments, err := gen.ExtraComment(f.Doc); err == nil {
				field.CommentArr = comments
			}
		}
		if fieldName == "id" {
			field.Name = "_id"
			field.CamelName = "Id"
		} else {
			field.Name = fieldName
			field.CamelName = camelString(fieldName)
		}
		field.TypeName = typeStr
		if typeValue, ok := type_name_dict[typeStr]; ok {
			field.Type = typeValue
			field.GoTypeName = go_type_name_dict[field.Type]
		} else {
			field.Type = field_type_struct
			field.GoTypeName = typeStr
		}
		coll.FieldDict[fieldName] = field
		coll.FieldArr = append(coll.FieldArr, field)
	}

	for _, f := range coll.FieldArr {
		if f.Name == "log_time" && f.Type == field_type_int64 {
			f.Coll.HasLogTimeField = true
		}
		if f.Name == "create_time" && f.Type == field_type_int64 {
			f.Coll.HasCreateTimeField = true
		}
	}
	sort.Sort(sort_field_by_name(coll.FieldArr))
	return nil
}

func (p *golang_parser) parseIndex(state *gen_state, coll *Collection, fileContent []byte, fset *token.FileSet, s *ast.StructType) error {
	coll.IndexDict = make(map[string]*Index)
	coll.IndexArr = make([]*Index, 0)

	for _, f := range s.Fields.List {
		var indexName string
		var tagStr string
		var tag int
		var err error
		tagStr = f.Tag.Value
		tagStr = strings.Replace(tagStr, "`", "", 2)
		if tag, err = strconv.Atoi(tagStr); err != nil {
			return err
		}
		indexName = f.Names[0].Name
		index := &Index{
			Name: indexName,
			Tag:  tag,
		}
		// log.Println(index, indexName)
		// ast.Print(fset, s)
		if st, ok := f.Type.(*ast.StructType); ok {
			var fullName string
			index.KeyDict = make(map[string]*Key)
			for _, f1 := range st.Fields.List {
				k := f1.Names[0].Name
				v := f1.Tag.Value
				v = strings.Replace(v, "`", "", 2)
				indexKey := &Key{Key: k, Value: v}
				index.KeyDict[k] = indexKey
				index.KeyArr = append(index.KeyArr, indexKey)
				if len(index.KeyDict) == 1 {
					fullName = fmt.Sprintf("%s_%v", k, v)
				} else {
					fullName = fmt.Sprintf("%s_%s_%v", fullName, k, v)
				}
			}
			index.FullName = fullName
		}
		coll.IndexDict[indexName] = index
		coll.IndexArr = append(coll.IndexArr, index)
	}
	return nil

	// for valueStr, v := range attrs {
	// 	var indexName string
	// 	var tagStr string
	// 	var tag int
	// 	var err error
	// 	var optionArr []interface{}
	// 	switch v.(type) {
	// 	case nil:
	// 		break
	// 	case []interface{}:
	// 		optionArr = v.([]interface{})
	// 	default:
	// 		return fmt.Errorf("%+v invalid11", v)
	// 	}
	// 	args := equalRegexp.FindAllString(valueStr, -1)
	// 	if len(args) != 2 {
	// 		return fmt.Errorf("%s invalid", valueStr)
	// 	}
	// 	tagStr = strings.TrimSpace(args[1])
	// 	if tag, err = strconv.Atoi(tagStr); err != nil {
	// 		return err
	// 	}
	// 	indexName = args[0]
	// 	index := &Index{
	// 		Name: indexName,
	// 		Tag:  tag,
	// 	}
	// 	for _, option := range optionArr {
	// 		optionDict := option.(map[string]interface{})
	// 		if keyArr, ok := optionDict["key"]; ok {
	// 			var fullName string
	// 			index.KeyDict = make(map[string]*Key)
	// 			if keyArr2, ok := keyArr.([]interface{}); ok {
	// 				for _, keyObject := range keyArr2 {
	// 					if keyObject2, ok := keyObject.(map[string]interface{}); ok {
	// 						for k, v := range keyObject2 {
	// 							indexKey := &Key{Key: k, Value: v}
	// 							index.KeyDict[k] = indexKey
	// 							index.KeyArr = append(index.KeyArr, indexKey)
	// 							if len(index.KeyDict) == 1 {
	// 								fullName = fmt.Sprintf("%s_%v", k, v)
	// 							} else {
	// 								fullName = fmt.Sprintf("%s_%s_%v", fullName, k, v)
	// 							}
	// 						}
	// 					}
	// 				}
	// 			}
	// 			index.FullName = fullName
	// 		}
	// 	}
	// 	// log.Println(index.KeyDict)
	// 	// log.Println(index.FullName)
	// 	coll.IndexDict[indexName] = index
	// 	coll.IndexArr = append(coll.IndexArr, index)
	// }
	// sort.Sort(sort_index_by_name(coll.IndexArr))
	return nil
}
