package gen_model

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
	filepath.WalkDir(proj.Config.DocModelDir, func(path string, d os.DirEntry, err error) error {
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
	})
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
		GenBinFilePath:        path.Join(proj.Config.SrcGenModelDir, dbName, "bin", fmt.Sprintf("%s.gen.go", dbName)),
		GenBinDir:             path.Join(proj.Config.SrcGenModelDir, dbName, "bin"),
		GenModelDir:           path.Join(proj.Config.SrcGenModelDir, dbName),
		GenModelFilePath:      path.Join(proj.Config.SrcGenModelDir, dbName, fmt.Sprintf("%s.gen.go", dbName)),
		GenProtobufFilePath:   path.Join(proj.Config.GenModelDir, dbName, fmt.Sprintf("%s.gen.proto", dbName)),
		DbName:                dbName,
		DbStructName:          camelString(dbName),
		MongoDriverStructName: fmt.Sprintf("%sMongoDriver", camelString(dbName)),
		RedisDriverStructName: fmt.Sprintf("%sRedisDriver", camelString(dbName)),
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
						if err = p.parseDriver(state, database, filePath, fileContent, fset, s); err != nil {
							return false
						}
					} else if s.Names[0].Name == "Import" {
						if err = p.parseImport(state, database, filePath, fileContent, fset, s); err != nil {
							return false
						}
					}
				case *ast.TypeSpec:
					var comment string
					if x.Doc != nil {
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
	if err != nil {
		return err
	}
	state.databaseArr = append(state.databaseArr, database)
	return nil
}

func (p *golang_parser) parseDriver(state *gen_state, database *Database, filePath string, fileContent []byte, fset *token.FileSet, s *ast.ValueSpec) error {
	if values0, ok := s.Values[0].(*ast.CompositeLit); !ok {
		return p.reportError(fset, filePath, s.Pos(), "driver必须是string数组")
	} else if typ, ok := values0.Type.(*ast.ArrayType); !ok {
		return p.reportError(fset, filePath, values0.Pos(), "driver必须是string数组")
	} else if elt, ok := typ.Elt.(*ast.Ident); !ok {
		return p.reportError(fset, filePath, values0.Pos(), "driver必须是string数组")
	} else if elt.Name != "string" {
		return p.reportError(fset, filePath, values0.Pos(), "driver必须是string数组")
	} else if len(values0.Elts) <= 0 {
		return p.reportError(fset, filePath, values0.Pos(), "至少一个元素")
	} else {
		if v, ok := values0.Elts[0].(*ast.BasicLit); !ok {
			return p.reportError(fset, filePath, values0.Pos(), "至少一个元素")
		} else {
			database.Driver = v.Value[1 : len(v.Value)-1]
			return nil
		}
	}
}

func (p *golang_parser) parseImport(state *gen_state, database *Database, filePath string, fileContent []byte, fset *token.FileSet, s *ast.ValueSpec) error {
	if values0, ok := s.Values[0].(*ast.CompositeLit); !ok {
		return p.reportError(fset, filePath, s.Pos(), "driver必须是string数组")
	} else if typ, ok := values0.Type.(*ast.ArrayType); !ok {
		return p.reportError(fset, filePath, values0.Pos(), "driver必须是string数组")
	} else if elt, ok := typ.Elt.(*ast.Ident); !ok {
		return p.reportError(fset, filePath, values0.Pos(), "driver必须是string数组")
	} else if elt.Name != "string" {
		return p.reportError(fset, filePath, values0.Pos(), "driver必须是string数组")
	} else if len(values0.Elts) <= 0 {
		return p.reportError(fset, filePath, values0.Pos(), "至少一个元素")
	} else {
		if v, ok := values0.Elts[0].(*ast.BasicLit); !ok {
			return p.reportError(fset, filePath, values0.Pos(), "至少一个元素")
		} else {
			database.ImportArr = append(database.ImportArr, v.Value[1:len(v.Value)-1])
			return nil
		}
	}
}

func (p *golang_parser) parseCollection(state *gen_state, database *Database, fileContent []byte, collName string, fset *token.FileSet, s *ast.StructType) (*Collection, error) {
	coll := &Collection{
		MongoDriverStructName: database.MongoDriverStructName,
		CollName:              collName,
		StructName:            camelString(collName),
		PbStructName:          fmt.Sprintf("%sPb", camelString(collName)),
		ArrStructName:         fmt.Sprintf("%sArr", camelString(collName)),
		DataStructName:        fmt.Sprintf("%sData", camelString(collName)),
		MongoDaoStructName:    fmt.Sprintf("%sMongoDao", camelString(collName)),
		RedisDaoStructName:    fmt.Sprintf("%sRedisDao", camelString(collName)),
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
		} else if f.Names[0].Name == "Derive" {
			if derive, ok := f.Type.(*ast.StructType); ok {
				if err := p.parseDerive(state, coll, fileContent, fset, derive); err != nil {
					return nil, err
				}
			}
		} else if f.Names[0].Name == "Capped" {
			capped := strings.Replace(f.Tag.Value, "`", "", 2)
			if v, err := strconv.Atoi(capped); err != nil {
				return nil, err
			} else {
				coll.Capped = int64(v)
			}
		}
	}
	if coll.Derive == "userarr" {
		if len(coll.KeyArr) != 2 {
			return nil, fmt.Errorf("collection %s derive userarr need 2 key", coll.CollName)
		}
		primaryKey := coll.KeyArr[0]
		secondaryKey := coll.KeyArr[1]
		coll.PrimaryKey = primaryKey
		coll.CamelPrimaryKey = camelString(primaryKey)
		coll.CapCamelPrimaryKey = capLowerString(coll.CamelPrimaryKey)
		coll.SecondaryKey = secondaryKey
		coll.CamelSecondaryKey = camelString(secondaryKey)
		coll.CapCamelSecondaryKey = capLowerString(coll.CamelSecondaryKey)
		if field, ok := coll.FieldDict[primaryKey]; ok {
			coll.PrimaryKeyField = field
			field.IsPrimaryKey = true
		} else {
			return nil, fmt.Errorf("collection %s derive userarr, but primary key %s not found", coll.CollName, primaryKey)
		}
		if field, ok := coll.FieldDict[secondaryKey]; ok {
			coll.SecondaryKeyField = field
			field.IsSecondaryKey = true
		} else {
			return nil, fmt.Errorf("collection %s derive userarr, but secondary key %s not found", coll.CollName, secondaryKey)
		}
	}
	return coll, nil
}

func (p *golang_parser) parseDerive(state *gen_state, coll *Collection, fileContent []byte, fset *token.FileSet, s *ast.StructType) error {
	for _, f := range s.Fields.List {
		if f.Names[0].Name == "User" {
			coll.Derive = "user"
			if keys, ok := f.Type.(*ast.StructType); ok {
				if err := p.parseDeriveKeys(state, coll, fileContent, fset, keys); err != nil {
					return err
				}
			}
		} else if f.Names[0].Name == "UserArr" {
			coll.Derive = "userarr"
			if keys, ok := f.Type.(*ast.StructType); ok {
				if err := p.parseDeriveKeys(state, coll, fileContent, fset, keys); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (p *golang_parser) parseDeriveKeys(state *gen_state, coll *Collection, fileContent []byte, fset *token.FileSet, s *ast.StructType) error {
	for _, f := range s.Fields.List {
		// log.Println(f.Names[0].Name)
		coll.KeyArr = append(coll.KeyArr, f.Names[0].Name)
	}
	return nil
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
		if f.Tag == nil {
			return p.reportError(fset, "fff", f.Pos(), "tag error")
		}
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
			Array:     false,
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
			field.ProtobufTypeName = protobuf_type_name_dict[field.Type]
		} else {
			field.Type = field_type_struct
			field.GoTypeName = typeStr
			field.ProtobufTypeName = "bytes"
		}
		coll.FieldDict[fieldName] = field
		coll.FieldArr = append(coll.FieldArr, field)
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
			Name:     indexName,
			Tag:      tag,
			FullName: indexName,
		}
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
	sort.Sort(sort_index_by_name(coll.IndexArr))
	return nil
}
