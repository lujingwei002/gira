package gen_resource

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/lujingwei002/gira/gen"
	"github.com/lujingwei002/gira/log"
	"github.com/lujingwei002/gira/proj"
)

type golang_parser struct {
}

func (p *golang_parser) parseStruct(resource *Resource, attrs map[string]interface{}) error {
	spaceRegexp := regexp.MustCompile("[^\\s]+")
	equalRegexp := regexp.MustCompile("[^=]+")
	fileArr := make([]*Field, 0)

	for valueStr, v := range attrs {
		var tag int
		var err error
		var fieldName string
		var typeStr string
		var tagStr string
		var optionArr []interface{}
		switch v.(type) {
		case nil:
			break
		case []interface{}:
			optionArr = v.([]interface{})
		default:
			return fmt.Errorf("%+v invalid11", v)
		}
		args := equalRegexp.FindAllString(valueStr, -1)
		if len(args) != 2 {
			return fmt.Errorf("%s invalid", valueStr)
		}
		tagStr = strings.TrimSpace(args[1])
		if tag, err = strconv.Atoi(tagStr); err != nil {
			return err
		}
		args = spaceRegexp.FindAllString(args[0], -1)
		if len(args) != 2 {
			return fmt.Errorf("%s invalid", valueStr)
		}
		fieldName = args[0]
		typeStr = args[1]
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
		for _, option := range optionArr {
			log.Println(option)
			optionDict := option.(map[string]interface{})
			if defaultVal, ok := optionDict["default"]; ok {
				field.Default = defaultVal
			}
			if comment, ok := optionDict["comment"]; ok {
				field.Comment = comment.(string)
			}
		}
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
	sort.Sort(SortFieldByName(fileArr))
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
		if annotates, err := gen.ExtraAnnotate(field.Doc.List); err != nil {
			return err
		} else {
			name := field.Names[0].Name
			var filePath string
			var keyArr []string
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
			// if len(v.Struct) > 0 {
			// 	if err := p.parseStruct(r, v.Struct); err != nil {
			// 		return err
			// 	}
			// }
			// 处理不同的转换类型， 转map, 转object
			//GoMapTypeName
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

func (p *golang_parser) parse(state *gen_state) error {
	filePath := path.Join(proj.Config.DocDir, "resource.go")
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		fmt.Println(err)
		return err
	}
	ast.Inspect(f, func(n ast.Node) bool {
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
	ast.Print(fset, f)
	return nil
}
