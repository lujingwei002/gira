package gen_resource

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path"
	"strings"

	"github.com/lujingwei002/gira"
	"github.com/lujingwei002/gira/gen"
	"github.com/lujingwei002/gira/log"
	"github.com/lujingwei002/gira/proj"
)

type golang_parser struct {
}

func (p *golang_parser) parseStruct(resource *Resource, s *ast.StructType) error {

	fileArr := make([]*Field, 0)
	for _, f := range s.Fields.List {
		var tag int
		var fieldName string
		var typeStr string
		if tags, err := gen.ExplodeTag(f.Tag); err != nil {
			log.Println("eeee1", err)
			return err
		} else {
			if v, err := tags.Int("tag"); err != nil {
				log.Println("eeee", err)
				return err
			} else {
				tag = v
			}
		}
		log.Println(f.Names)
		fieldName = f.Names[0].Name
		if v, ok := f.Type.(*ast.Ident); ok {
			typeStr = v.Name
		}
		field := &Field{
			FieldName:       fieldName,
			StructFieldName: camelString(fieldName),
			Tag:             tag,
		}
		if typeValue, ok := type_name_dict[typeStr]; ok {
			field.Type = typeValue
			field.GoTypeName = go_type_name_dict[field.Type]
		} else {
			field.Type = field_type_struct
			field.GoTypeName = typeStr
		}
		if comments, err := gen.ExtraComment(f.Doc); err != nil {
			return err
		} else {
			field.Comment = strings.Join(comments, "\n")
		}
		// for _, option := range optionArr {
		// 	log.Println(option)
		// 	optionDict := option.(map[string]interface{})
		// 	if defaultVal, ok := optionDict["default"]; ok {
		// 		field.Default = defaultVal
		// 	}
		// 	if comment, ok := optionDict["comment"]; ok {
		// 		field.Comment = comment.(string)
		// 	}
		// }
		if last, ok := resource.FieldDict[fieldName]; ok {
			resource.FieldDict[fieldName] = field
			// 替换之前的
			for k, v := range resource.FieldArr {
				if v == last {
					resource.FieldArr[k] = field
					break
				}
			}
		} else {
			resource.FieldDict[fieldName] = field
			fileArr = append(fileArr, field)

		}
	}
	// sort.Sort(SortFieldByName(fileArr))
	resource.FieldArr = append(resource.FieldArr, fileArr...)
	return nil
}

func (p *golang_parser) parseBundleStruct(state *gen_state, s *ast.StructType) error {
	for _, v := range s.Fields.List {
		name := v.Names[0].Name
		bundleType := v.Type.(*ast.StructType)
		if annotates, err := gen.ExtraAnnotate(v.Doc.List); err != nil {
			return err
		} else {
			format := ""
			if v, ok := annotates["format"]; ok {
				format = v.Values[0]
			}
			resources := make([]string, 0)
			for _, v2 := range bundleType.Fields.List {
				resources = append(resources, v2.Names[0].Name)
			}
			log.Println("parse bundle ", name, resources)
			bundle := &Bundle{
				BundleType:          format,
				BundleName:          name,
				BundleStructName:    name,
				CapBundleStructName: capLowerString(name),
				ResourceNameArr:     resources,
				ResourceArr:         make([]*Resource, 0),
			}
			state.BundleDict[name] = bundle
			state.BundleArr = append(state.BundleArr, bundle)
		}
	}
	return nil
}

func (p *golang_parser) parseLoaderStruct(state *gen_state, s *ast.StructType) error {
	for _, v := range s.Fields.List {
		name := v.Names[0].Name
		bundles := make([]string, 0)
		loaderType := v.Type.(*ast.StructType)
		for _, v2 := range loaderType.Fields.List {
			bundles = append(bundles, v2.Names[0].Name)
		}
		log.Println("parse loader ", name, bundles)
		loader := &Loader{
			LoaderName:        name,
			LoaderStructName:  name,
			HandlerStructName: "I" + name + "Handler",
			bundleNameArr:     bundles,
			BundleArr:         make([]*Bundle, 0),
		}
		state.LoaderArr = append(state.LoaderArr, loader)
	}
	return nil
}

func (p *golang_parser) parseResourceStruct(state *gen_state, s *ast.StructType) error {
	for _, field := range s.Fields.List {
		name := field.Names[0].Name
		log.Println("parse resource ", name)

		if annotates, err := gen.ExtraAnnotate(field.Doc.List); err != nil {

			return err
		} else {
			var filePath string
			var keyArr []string
			if v, ok := annotates["excel"]; ok {
				filePath = v.Values[0]
			}
			// 类型
			var typ resource_type = resource_type_array
			if v, ok := annotates["map"]; ok {
				typ = resource_type_map
				keyArr = v.Values
			}
			if v, ok := annotates["object"]; ok {
				typ = resource_type_object
				keyArr = v.Values
			}

			yamlFileName := fmt.Sprintf("%s.yaml", name)
			camelName := camelString(name)
			r := &Resource{
				FieldDict:      make(map[string]*Field, 0),
				FieldArr:       make([]*Field, 0),
				ValueArr:       make([][]interface{}, 0),
				Type:           typ,
				ResourceName:   name,
				StructName:     camelName,
				TableName:      name,
				CapStructName:  capLowerString(camelName),
				ArrTypeName:    camelName + "Arr",
				MapTypeName:    camelName + "Map",
				ObjectTypeName: camelName + "Object",
				FilePath:       filePath,
				YamlFileName:   yamlFileName,
				MapKeyArr:      keyArr,
				ObjectKeyArr:   keyArr,
			}
			if typ == resource_type_array {
				r.WrapStructName = r.ArrTypeName
			} else if typ == resource_type_map {
				r.WrapStructName = r.MapTypeName
			} else if typ == resource_type_object {
				r.WrapStructName = r.ObjectTypeName
			}
			// 解析excel
			// filePath := path.Join(proj.Config.ExcelDir, v.FilePath)
			if err := r.readExcel(path.Join(proj.Config.ExcelDir, filePath)); err != nil {
				return err
			}
			// 解析struct
			if v, ok := field.Type.(*ast.StructType); ok {
				if err := p.parseStruct(r, v); err != nil {
					return err
				}
			}
			// 处理不同的转换类型， 转map, 转object
			r.GoMapTypeName = ""
			for _, v := range r.MapKeyArr {
				if field, ok := r.FieldDict[v]; ok {
					r.GoMapTypeName = r.GoMapTypeName + fmt.Sprintf(`map[%s]`, field.GoTypeName)
				} else {
					return fmt.Errorf("resource %s key %s not found\n", r.StructName, v)
				}
			}
			r.GoMapTypeName = r.GoMapTypeName + fmt.Sprintf(" *%s", r.StructName)
			// GoObjectTypeName
			r.GoObjectTypeName = ""
			for _, v := range r.ObjectKeyArr {
				if field, ok := r.FieldDict[v]; ok {
					r.GoObjectTypeName = r.GoObjectTypeName + fmt.Sprintf(`map[%s]`, field.GoTypeName)
				} else {
					return fmt.Errorf("resource %s key %s not found\n", r.StructName, v)
				}
			}
			r.GoObjectTypeName = r.GoObjectTypeName + fmt.Sprintf(" *%s", r.StructName)
			// WrapTypeName
			if r.Type == resource_type_map {
				r.WrapTypeName = r.GoMapTypeName
			} else if r.Type == resource_type_array {
				r.WrapTypeName = fmt.Sprintf("[]* %s", r.StructName)
			} else if r.Type == resource_type_object {
				for k, v := range r.FieldArr {
					if v.FieldName == r.ObjectKeyArr[0] {
						r.ObjectKeyIndex = k
						break
					}
				}
			}
			// make key时用的
			mapKeyGoTypeNameArr := make([]string, 0)
			for _, k := range r.MapKeyArr {
				f, _ := r.FieldDict[k]
				mapKeyGoTypeNameArr = append(mapKeyGoTypeNameArr, f.GoTypeName)
			}
			r.MapKeyGoTypeNameArr = mapKeyGoTypeNameArr

			// make key时用的
			objectKeyGoTypeNameArr := make([]string, 0)
			for _, k := range r.ObjectKeyArr {
				f, _ := r.FieldDict[k]
				objectKeyGoTypeNameArr = append(objectKeyGoTypeNameArr, f.GoTypeName)
			}
			r.ObjectKeyGoTypeNameArr = objectKeyGoTypeNameArr
			state.ResourceDict[name] = r

			state.ResourceArr = append(state.ResourceArr, r)
		}
	}
	return nil
}

func (p *golang_parser) parse(state *gen_state) (err error) {
	filePath := path.Join(proj.Config.DocDir, "resource.go")
	fset := token.NewFileSet()
	var f *ast.File
	f, err = parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		fmt.Println(err)
		return
	}
	ast.Print(fset, f)
	return gira.ErrTodo
	ast.Inspect(f, func(n ast.Node) bool {
		if err != nil {
			return false
		}
		switch x := n.(type) {
		case *ast.TypeSpec:
			if x.Name.Name == "Resource" {
				if err = p.parseResourceStruct(state, x.Type.(*ast.StructType)); err != nil {
					return false
				}
			} else if x.Name.Name == "Bundle" {
				if err = p.parseBundleStruct(state, x.Type.(*ast.StructType)); err != nil {
					return false
				}
			} else if x.Name.Name == "Loader" {
				if err = p.parseLoaderStruct(state, x.Type.(*ast.StructType)); err != nil {
					return false
				}
			}
		}
		return true
	})
	for _, v := range state.BundleArr {
		for _, name := range v.ResourceNameArr {
			if r, ok := state.ResourceDict[name]; ok {
				v.ResourceArr = append(v.ResourceArr, r)
			}
		}
	}
	for _, v := range state.LoaderArr {
		for _, name := range v.bundleNameArr {
			if r, ok := state.BundleDict[name]; ok {
				v.BundleArr = append(v.BundleArr, r)

			}
		}
	}
	return
}
