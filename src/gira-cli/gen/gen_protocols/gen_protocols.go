package gen_protocols

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/lujingwei/gira-cli/proj"

	yaml "gopkg.in/yaml.v3"
)

type field_type int

const (
	field_type_int32 field_type = iota
	field_type_nt64
	field_type_string
	field_type_message
)

var type_name_dict = map[string]field_type{
	"int32":  field_type_int32,
	"int64":  field_type_nt64,
	"string": field_type_string,
}

type message_type int

const (
	message_type_struct message_type = iota
	message_type_request
	message_type_response
	message_type_notify
	message_type_push
)

type message_field struct {
	Tag      int
	Name     string
	Type     field_type
	Array    bool
	TypeName string
	Default  interface{}
}

type Message struct {
	FieldDict map[string]message_field
}

type Packet struct {
	MessageId int
	Name      string
	Type      message_type
	message   Message
	request   Message
	response  Message
	push      Message
	notify    Message
}

func (p *Packet) String() string {
	sb := strings.Builder{}
	var dot string
	if p.Type == message_type_struct {
		dot = "."
	} else {
		dot = ""
	}
	messageStringFunc := func(sb *strings.Builder, m *Message, layer int) {
		for _, v := range m.FieldDict {
			sb.WriteString(fmt.Sprintf("%s %s %d: %s\n", strings.Repeat("    ", layer), v.Name, v.Tag, v.TypeName))
		}
	}
	sb.WriteString(fmt.Sprintf("%s%s %d {\n", dot, p.Name, p.MessageId))
	if p.Type == message_type_struct {
		messageStringFunc(&sb, &p.message, 1)
	} else if p.Type == message_type_request {
		sb.WriteString(fmt.Sprintf("    request {\n"))
		messageStringFunc(&sb, &p.request, 2)
		sb.WriteString(fmt.Sprintf("    }\n"))
		sb.WriteString(fmt.Sprintf("    response {\n"))
		messageStringFunc(&sb, &p.response, 2)
		sb.WriteString(fmt.Sprintf("    }\n"))
	} else if p.Type == message_type_notify {
		sb.WriteString("    response {\n")
		messageStringFunc(&sb, &p.notify, 2)
		sb.WriteString("    }\n")
	} else if p.Type == message_type_push {
		sb.WriteString("    request {\n")
		messageStringFunc(&sb, &p.push, 2)
		sb.WriteString("    }\n")
	}
	sb.WriteString("}\n")
	return sb.String()
}

// 生成协议的状态
type protocol_state struct {
	packetDict map[string]*Packet
}

func (m *Message) parse(arr []interface{}) error {
	m.FieldDict = make(map[string]message_field)
	//commaRegexp := regexp.MustCompile("[^,]+")
	spaceRegexp := regexp.MustCompile("[^\\s]+")
	equalRegexp := regexp.MustCompile("[^=]+")

	for _, v := range arr {
		var valueStr string
		var tag int
		var err error
		var fieldName string
		var typeStr string
		var tagStr string
		var optionArr []interface{}
		switch v.(type) {
		case string:
			valueStr = v.(string)
		case map[string]interface{}:
			for k1, v1 := range v.(map[string]interface{}) {
				valueStr = k1
				optionArr = v1.([]interface{})
			}
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
		typeStr = args[0]
		fieldName = args[1]

		field := message_field{
			Name:  fieldName,
			Array: false,
			Tag:   tag,
		}
		//args = commaRegexp.FindAllString(v.(string), -1)
		//if len(args) <= 0 {
		//	fmt.Println(args)
		//	return fmt.Errorf("field args wrong %s %s %s", "name", k, v)
		//}
		field.TypeName = typeStr
		if typeValue, ok := type_name_dict[typeStr]; !ok {
			field.Type = typeValue
		}
		fmt.Println(optionArr)
		for _, option := range optionArr {
			optionDict := option.(map[string]interface{})
			if defaultVal, ok := optionDict["default"]; ok {
				field.Default = defaultVal
			}
		}
		m.FieldDict[fieldName] = field
	}
	return nil
}

func genProtocols1(protocolState *protocol_state, filePathArr []string) error {
	fmt.Println(filePathArr)
	// 提取 request, response, push, notify

	spaceRegexp := regexp.MustCompile("[^\\s]+")

	for _, filePath := range filePathArr {
		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			fmt.Printf("err: %v\n", err)
			return err
		}
		result := make(map[string]interface{})
		if err := yaml.Unmarshal(data, result); err != nil {
			fmt.Printf("err: %v\n", err)
			return err
		}
		for k, v := range result {
			args := spaceRegexp.FindAllString(k, -1)
			if len(args) == 1 && strings.HasPrefix(args[0], ".") {
				name := args[0][1:]
				packet := &Packet{
					Name: name,
					Type: message_type_struct,
				}
				if arr, ok := v.([]interface{}); !ok {
					return fmt.Errorf("%s invalid", k)
				} else {
					if err := packet.message.parse(arr); err != nil {
						return err
					}
					protocolState.packetDict[name] = packet
				}

			} else if len(args) == 2 {
				name := args[0]
				packet := &Packet{
					Name: name,
				}
				if message_id, err := strconv.Atoi(args[1]); err != nil {
					return err
				} else {
					packet.MessageId = message_id
				}
				if dict, ok := v.(map[string]interface{}); !ok {
					return fmt.Errorf("%s invalid", k)
				} else {
					hasRequest := false
					hasResponse := false
					if _, ok := dict["request"]; ok {
						if _, ok := dict["request"].([]interface{}); !ok {
							return fmt.Errorf("%s invalid", k)
						}
						hasRequest = true
					}
					if _, ok := dict["response"]; ok {
						if _, ok := dict["response"].([]interface{}); !ok {
							return fmt.Errorf("%s invalid", k)
						}
						hasResponse = true
					}
					if hasRequest && hasResponse {
						packet.Type = message_type_request
					} else if hasRequest {
						packet.Type = message_type_push
					} else if hasResponse {
						packet.Type = message_type_notify
					} else {
						return fmt.Errorf("%s request and response not found", name)
					}
					if packet.Type == message_type_request {
						if err := packet.request.parse(dict["request"].([]interface{})); err != nil {
							return err
						}
						if err := packet.response.parse(dict["response"].([]interface{})); err != nil {
							return err
						}
					} else if packet.Type == message_type_notify {
						if err := packet.notify.parse(dict["response"].([]interface{})); err != nil {
							return err
						}
					} else if packet.Type == message_type_push {
						if err := packet.push.parse(dict["request"].([]interface{})); err != nil {
							return err
						}
					}
					protocolState.packetDict[name] = packet
				}

			}
		}
	}
	//for _, v := range protocolState.packetDict {
	//fmt.Printf("%s", v)
	//}

	return nil
}

func genProtocols2(state *proj.State, protocolState *protocol_state) error {
	filePath := path.Join(state.GenDir, "protocols.sproto")
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	for _, v := range protocolState.packetDict {
		file.WriteString(v.String())
		file.WriteString("\n")
	}
	return nil
}

// 生成协议
func Gen(state *proj.State) error {
	protocolFilePathArr := make([]string, 0)
	protocolsFilePath := path.Join(state.DocDir, "protocols.yaml")
	if _, err := os.Stat(protocolsFilePath); !os.IsNotExist(err) {
		protocolFilePathArr = append(protocolFilePathArr, protocolsFilePath)
	}
	protocolsDir := path.Join(state.DocDir, "protocols")
	filepath.WalkDir(protocolsDir, func(path string, d os.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		if filepath.Ext(path) == ".yaml" {
			protocolFilePathArr = append(protocolFilePathArr, path)
		}
		return nil
	})
	protocolState := &protocol_state{
		packetDict: make(map[string]*Packet),
	}
	if err := genProtocols1(protocolState, protocolFilePathArr); err != nil {
		log.Println(err)
		return err
	}
	if err := genProtocols2(state, protocolState); err != nil {
		log.Println(err)
		return err
	}
	return nil
}
