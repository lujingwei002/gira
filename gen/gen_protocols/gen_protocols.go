package gen_protocols

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"text/template"
	"unicode"

	"github.com/lujingwei002/gira/log"

	"github.com/lujingwei002/gira/proj"

	yaml "gopkg.in/yaml.v3"
)

type field_type int

var sproto_template = `
<<.Header>>
<<- range .PacketArr>>
# <<.Comment>>
<<if .Type.IsStructType>>.<<.StructName>><<else>><<.Name>> <<.MessageId>><<- end>> {

	<<- if .Type.IsStructType>>
	<<- range .Message.FieldArr>>
	<<- if .HasComment>>
	# <<.Comment>> 
	<<- end>>
    <<.Name>> <<.Tag>> : <<- if .IsArray>>*<<- end>><<.SpTypeName>>
	<<- end>>
	<<- end>>

	<<- if .Type.IsRequestType>>
    request {
		<<- range .Request.FieldArr>>
		<<- if .HasComment>>
		# <<.Comment>> 
		<<- end>>
		<<.Name>> <<.Tag>> : <<- if .IsArray>>*<<- end>><<.SpTypeName>>
		<<- end>>
	}
	response {
		<<- range .Response.FieldArr>>
		<<- if .HasComment>>
		# <<.Comment>> 
		<<- end>>
		<<.Name>> <<.Tag>> : <<- if .IsArray>>*<<- end>><<.SpTypeName>>
		<<- end>>
	}
	<<- end>>

	<<- if .Type.IsNotifyType>>
	response {
		<<- range .Notify.FieldArr>>
		<<- if .HasComment>>
		# <<.Comment>> 
		<<- end>>
		<<.Name>> <<.Tag>> : <<- if .IsArray>>*<<- end>><<.SpTypeName>>
		<<- end>>
	}
	<<- end>>

	<<- if .Type.IsPushType>>
    request {
		<<- range .Push.FieldArr>>
		<<- if .HasComment>>
		# <<.Comment>> 
		<<- end>>
		<<.Name>> <<.Tag>> : <<- if .IsArray>>*<<- end>><<.SpTypeName>>
		<<- end>>
	}
	<<- end>>
}
<<end>>
`

var sproto_go_template = `
// Code generated by github.com/lujingwei002/gira. DO NOT EDIT.
// Code generated by github.com/lujingwei002/gira. DO NOT EDIT.
// Code generated by github.com/lujingwei002/gira. DO NOT EDIT.

package <<.Module>>

import (
	"reflect"
	"context"
	gosproto "github.com/xjdrew/gosproto"
	"github.com/lujingwei002/gira"
	"github.com/lujingwei002/gira/log"
	"github.com/lujingwei002/gira/sproto"
	"sync/atomic"
	"sync"
)

var _ reflect.Type

<<- range .PacketArr>>
	<<- /* struct 类型 */>>
	<<- if .Type.IsStructType>>
type <<.Message.StructName>> struct {
	<<- range .Message.FieldArr>>
	// <<.Comment>>
    <<.CamelName>> <<if .IsArray>>[]<<end>><<.GoTypeName>> <<quote>>sproto:"<<.SpWireTypeName>>,<<.Tag>><<if .IsArray>>,array<<end>>"<<quote>>
	<<- end>>
}
	<<end>>
	<<- /* request 类型 */>>
	<<- if .Type.IsRequestType>>
type <<.Request.StructName>> struct {
	<<- range .Request.FieldArr>>
	// <<.Comment>>
    <<.CamelName>> <<if .IsArray>>[]<<end>><<.GoTypeName>> <<quote>>sproto:"<<.SpWireTypeName>>,<<.Tag>><<if .IsArray>>,array<<end>>"<<quote>>
	<<- end>>
}

type <<.Response.StructName>> struct {
	<<- range .Response.FieldArr>>
	// <<.Comment>>
    <<.CamelName>> <<if .IsArray>>[]<<end>><<.GoTypeName>> <<quote>>sproto:"<<.SpWireTypeName>>,<<.Tag>><<if .IsArray>>,array<<end>>"<<quote>>
	<<- end>>
}

func (self *<<.Request.StructName>>) GetRequestName() string {
	return "<<.StructName>>"
}

<<- range .Request.FieldArr>>
// <<.Comment>>
func (self *<<.Message.StructName>>) Get<<.CamelName>>() <<if .IsArray>>[]<<end>><<.GoTypeName>> {
	return self.<<.CamelName>> 
}
<<- end>>

func (self *<<.Response.StructName>>) SetErrorCode(v int32) {
	self.ErrorCode = v
}
func (self *<<.Response.StructName>>) SetErrorMsg(v string) {
	self.ErrorMsg = v
}
<<end>>
	<<- /* notify 类型 */>>
	<<- if .Type.IsNotifyType>>
type <<.Notify.StructName>> struct {
	<<- range .Notify.FieldArr>>
	// <<.Comment>>
    <<.CamelName>> <<if .IsArray>>[]<<end>><<.GoTypeName>> <<quote>>sproto:"<<.SpWireTypeName>>,<<.Tag>><<if .IsArray>>,array<<end>>"<<quote>>
	<<- end>>
}

	<<end>>
	<<- /* push 类型 */>>
	<<- if .Type.IsPushType>>
type <<.Push.StructName>> struct {
	<<- range .Push.FieldArr>>
	// <<.Comment>>
    <<.CamelName>> <<if .IsArray>>[]<<end>><<.GoTypeName>> <<quote>>sproto:"<<.SpWireTypeName>>,<<.Tag>><<if .IsArray>>,array<<end>>"<<quote>>
	<<- end>>
}

func (self *<<.Push.StructName>>) GetPushName() string {
	return "<<.StructName>>"
}
	<<end>>
<<- end>>

var Protocols []*gosproto.Protocol = []*gosproto.Protocol {
<<- range .PacketArr>>
<<- if .Type.IsStructType>>
<<- else>>
	{
		Type: 		<<.MessageId>>,
		Name: 		"<<.StructName>>",
		MethodName: "<<.StructName>>",
	<<- if .Type.IsPushType>>
		Request: reflect.TypeOf(&<<.Push.StructName>>{}),
	<<- end>>
	<<- if .Type.IsNotifyType>>
		Request: reflect.TypeOf(&<<.Notify.StructName>>{}),
	<<- end>>
	<<- if .Type.IsRequestType>>
		Request: reflect.TypeOf(&<<.Request.StructName>>{}),
		Response: reflect.TypeOf(&<<.Response.StructName>>{}),
	<<- end>>
	},
<<- end>>
<<- end>>
}


type Client struct {
	proto			*sproto.Sproto
	conn			gira.GateClient
	reqId			uint64
	ctx				context.Context
	requestDict     sync.Map
	// 末处理的push, 只保留最后一个
	pendingPushDict sync.Map
	waitPushDict    sync.Map
	onPushDict		map[string]*sync.Map

	// push 消息的侦听函数 
<<- range .PacketArr>>
<<- if .Type.IsPushType>>
	<< capLower .Push.StructName>>Func	sync.Map
	<< capLower .Push.StructName>>Id	uint64
<<- end>>
<<- end>>
	// push 消息的侦听函数 
}

type PushMethod interface {
	Call(req sproto.SprotoPush, err error)
}

<<- range .PacketArr>>
<<- if .Type.IsPushType>>
type <<.Push.StructName>>Func func (req *<<.Push.StructName>>, err error) 

type <<.Push.StructName>>Method struct {
	Func <<.Push.StructName>>Func
}

func (self* <<.Push.StructName>>Method) Call(req sproto.SprotoPush, err error) {
	r := req.(*<<.Push.StructName>>)
	self.Func(r, err)
}

<<- end>>
<<- end>>

type waitRequest struct {
	route	string
	data	[]byte
	err		error
    caller	chan*waitRequest
	reqId	uint64
	mode	gosproto.RpcMode
	name    string
	session int32
	resp    interface{}
}

type waitPush struct {
	route	string
	data	[]byte
	err		error
    caller	chan*waitPush
	mode	gosproto.RpcMode
	name    string
	session int32
	resp    interface{}
}

func NewClient(ctx context.Context, conn gira.GateClient, sproto *sproto.Sproto) *Client {
	self := &Client{
		proto:		sproto,
		conn:		conn,
		ctx:		ctx,
		onPushDict:	make(map[string]*sync.Map, 0),
	}
<<- range .PacketArr>>
<<- if .Type.IsPushType>>
	self.onPushDict["<<.StructName>>"] = &self.<<capLower .Push.StructName>>Func 
<<- end>>
<<- end>>
	go self.readRoutine()
	return self
}

func (self *Client) readRoutine() error {
	cancelCtx, cancelFunc := context.WithCancel(self.ctx)
	defer cancelFunc()
	for {
		typ, route, reqId, data, err := self.conn.Recv(cancelCtx)
		if err != nil {
			//log.Info("client read routine exit", err)
			self.requestDict.Range(func(k any, v any) bool {
				wait := v.(*waitRequest)
				wait.err = err
				wait.caller <- wait	
				return true
			})
			return err
		}
		switch typ {
			case gira.GateMessageType_RESPONSE:
				if v, ok := self.requestDict.Load(reqId); ok {
					self.requestDict.Delete(reqId)
					wait := v.(*waitRequest)	
					wait.data = data
					// 上下文信息
					wait.route = route
					wait.reqId = reqId
					// 解析message
					wait.mode, wait.name, wait.session, wait.resp, wait.err = self.proto.ResponseDecode(data)
					wait.caller <- wait
				}
			case gira.GateMessageType_PUSH:
				mode, name, _, resp, err := self.proto.PushDecode(data)
				log.Infow("recv push message", "name", name, "data", resp)
				if mode != gosproto.RpcRequestMode {
					continue
				}
				handlerCount := 0
				if dict, ok := self.onPushDict[name]; ok {
					dict.Range(func(k any, v any) bool {
						handlerCount = handlerCount + 1
						f := v.(PushMethod)
						f.Call(resp.(sproto.SprotoPush), err)
						return true
					})
				}
				if v, ok := self.waitPushDict.Load(name); ok {
					handlerCount = handlerCount + 1
					self.waitPushDict.Delete(name)
					wait := v.(*waitPush)	
					wait.data = data
					// 上下文信息
					wait.route = route
					// 解析message
					wait.mode, wait.name, wait.session, wait.resp, wait.err = self.proto.PushDecode(data)
					wait.caller <- wait
				}
				if handlerCount == 0 {
					self.pendingPushDict.Store(name, resp)
					// log.Warnw("client proto push not register", "name", name)
				}
		}
	}
}

<<- range .PacketArr>>

<<- if .Type.IsRequestType>>
/// <<.Comment>>
func (self *Client) <<.StructName>>(ctx context.Context, req *<<.Request.StructName>>) (*<<.Response.StructName>>, error) {
	reqId := atomic.AddUint64(&self.reqId, 1)
	data, err := self.proto.RequestEncode("<<.StructName>>", int32(reqId), req)
	if err != nil {
		return nil, err
	}
	wait := &waitRequest {
		caller: make(chan* waitRequest, 1),
	}
	if _, loaded := self.requestDict.LoadOrStore(reqId, wait); loaded {
		return nil, gira.ErrSprotoReqIdConflict
	}
	if err := self.conn.Request("", reqId, data); err != nil {
		return nil, err
	}
	defer close(wait.caller)
	select {
		case v := <- wait.caller:
			if v.err != nil {
				return nil, v.err
			}
	        if wait.mode != gosproto.RpcResponseMode {
				return nil, gira.ErrSprotoResponseConversion
			}
	        if uint64(wait.session) != wait.reqId {
				return nil, gira.ErrSprotoResponseConversion 
			}
			resp1, ok := wait.resp.(*<<.Response.StructName>>)
			if !ok {
				return nil, gira.ErrSprotoResponseConversion
			}
			return resp1, nil
		case <-ctx.Done():
			return nil, gira.ErrSprotoReqTimeout
	}
}
<<- end>>


<<- if .Type.IsPushType>>

/// <<.Comment>>
func (self *Client) Wait<<.StructName>>Push(ctx context.Context) (*<<.Push.StructName>>, error) {
	wait := &waitPush {
		caller: make(chan* waitPush, 1),
	}
	if v, ok := self.pendingPushDict.Load("<<.StructName>>"); ok {
		if resp, ok := v.(*<<.Push.StructName>>); ok {
			return resp, nil
		}
	} 
	if _, loaded := self.waitPushDict.LoadOrStore("<<.StructName>>", wait); loaded {
		return nil, gira.ErrSprotoWaitPushConflict
	}
	defer close(wait.caller)
	select {
		case v := <- wait.caller:
			if v.err != nil {
				return nil, v.err
			}
	        if wait.mode != gosproto.RpcRequestMode {
				return nil, gira.ErrSprotoPushConversion
			}
			resp1, ok := wait.resp.(*<<.Push.StructName>>)
			if !ok {
				return nil, gira.ErrSprotoPushConversion
			}
			return resp1, nil
		case <-ctx.Done():
			return nil, gira.ErrSprotoWaitPushTimeout
	}
}

func (self *Client) On<<.StructName>>Push(f <<.Push.StructName>>Func) uint64 {
	id := atomic.AddUint64(&self.<< capLower .Push.StructName>>Id, 1)
	method := &<<.Push.StructName>>Method {
		Func: f,
	}
	self.<<capLower .Push.StructName>>Func.Store(id, method)
	return id
}

func (self *Client) Off<<.StructName>>Push(id uint64) error {
	self.<<capLower .Push.StructName>>Func.Delete(id)
	return nil
}

<<- end>>
<<- end>>

`

const (
	field_type_int field_type = iota

	field_type_int32
	field_type_int64
	field_type_string
	field_type_message
	field_type_bool
)

var type_name_dict = map[string]field_type{
	"int":    field_type_int,
	"int32":  field_type_int32,
	"int64":  field_type_int64,
	"string": field_type_string,
	"bool":   field_type_bool,
}

var sp_type_name_dict = map[field_type]string{
	field_type_int:    "integer",
	field_type_int32:  "integer",
	field_type_int64:  "integer",
	field_type_string: "string",
	field_type_bool:   "boolean",
}

var go_type_name_dict = map[field_type]string{
	field_type_int:    "int",
	field_type_int32:  "int32",
	field_type_int64:  "int64",
	field_type_string: "string",
	field_type_bool:   "bool",
}

type MessageType int

const (
	message_type_struct MessageType = iota
	message_type_request
	message_type_response
	message_type_notify
	message_type_push
)

type Field struct {
	Tag            int
	Name           string
	CamelName      string
	Type           field_type
	Array          bool
	TypeName       string
	SpTypeName     string
	SpWireTypeName string
	GoTypeName     string
	Default        interface{}
	IsArray        bool
	Comment        string
	HasComment     bool
	Message        *Message
}

type Message struct {
	StructName string
	FieldDict  map[string]*Field
	FieldArr   []*Field
}

func (s MessageType) IsStructType() bool {
	return s == message_type_struct
}

func (s MessageType) IsRequestType() bool {
	return s == message_type_request
}

func (s MessageType) IsSResponseType() bool {
	return s == message_type_response
}

func (s MessageType) IsNotifyType() bool {
	return s == message_type_notify
}

func (s MessageType) IsPushType() bool {
	return s == message_type_push
}

type Packet struct {
	MessageId  int
	Comment    string
	Name       string
	FullName   string
	StructName string
	Type       MessageType
	Message    Message
	Request    Message
	Response   Message
	Push       Message
	Notify     Message
}

func capUpperString(s string) string {
	return strings.ToUpper(s[0:1]) + s[1:]
}

func capLowerString(s string) string {
	return strings.ToLower(s[0:1]) + s[1:]
}
func QuoteChar() interface{} {
	return "`"
}
func camelString(s string) string {
	data := make([]byte, 0, len(s))
	j := false
	k := false
	num := len(s) - 1
	for i := 0; i <= num; i++ {
		d := s[i]
		if k == false && d >= 'A' && d <= 'Z' {
			k = true
		}
		if d >= 'a' && d <= 'z' && (j || k == false) {
			d = d - 32
			j = false
			k = true
		}
		if k && d == '_' && num > i && s[i+1] >= 'a' && s[i+1] <= 'z' {
			j = true
			continue
		}
		data = append(data, d)
	}
	return string(data[:])
}

// 生成协议的状态
type GenState struct {
	Module     string
	Header     string
	PacketDict map[string]*Packet
	PacketArr  []*Packet
}

func (m *Message) parse(attrs map[string]interface{}) error {
	m.FieldDict = make(map[string]*Field)
	m.FieldArr = make([]*Field, 0)
	//commaRegexp := regexp.MustCompile("[^,]+")
	spaceRegexp := regexp.MustCompile("[^\\s]+")
	equalRegexp := regexp.MustCompile("[^=]+")

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
			return fmt.Errorf("%+v invalid11 %s", v, valueStr)
		}
		args := equalRegexp.FindAllString(valueStr, -1)
		if len(args) != 2 {
			return fmt.Errorf("%s invalid22", valueStr)
		}
		tagStr = strings.TrimSpace(args[1])
		if tag, err = strconv.Atoi(tagStr); err != nil {
			return err
		}
		args = spaceRegexp.FindAllString(args[0], -1)
		if len(args) != 2 {
			return fmt.Errorf("%s invalid33", valueStr)
		}
		typeStr = args[1]
		isArray := false
		fieldName = args[0]
		if strings.HasPrefix(typeStr, "[]") {
			isArray = true
			typeStr = strings.Replace(typeStr, "[]", "", 1)
			typeStr = strings.TrimSpace(typeStr)
		}

		field := &Field{
			Name:      fieldName,
			CamelName: camelString(fieldName),
			Array:     false,
			Tag:       tag,
			IsArray:   isArray,
			Message:   m,
		}
		//args = commaRegexp.FindAllString(v.(string), -1)
		//if len(args) <= 0 {
		//	fmt.Println(args)
		//	return fmt.Errorf("field args wrong %s %s %s", "name", k, v)
		//}
		field.TypeName = typeStr
		if typeValue, ok := type_name_dict[typeStr]; ok {
			field.Type = typeValue
		} else {
			field.Type = field_type_message
		}
		if typeName, ok := sp_type_name_dict[field.Type]; ok {
			field.SpTypeName = typeName
		} else {
			field.SpTypeName = typeStr
		}
		if typeName, ok := sp_type_name_dict[field.Type]; ok {
			field.SpWireTypeName = typeName
		} else {
			field.SpWireTypeName = "struct"
		}
		if typeName, ok := go_type_name_dict[field.Type]; ok {
			field.GoTypeName = typeName
		} else {
			field.GoTypeName = fmt.Sprintf("*%s", typeStr)
		}

		// fmt.Println(optionArr)
		for _, option := range optionArr {
			optionDict := option.(map[string]interface{})
			if defaultVal, ok := optionDict["default"]; ok {
				field.Default = defaultVal
			}
			if comment, ok := optionDict["comment"]; ok {
				field.Comment = comment.(string)
				field.HasComment = true
			}
		}
		m.FieldDict[fieldName] = field
		m.FieldArr = append(m.FieldArr, field)
	}
	return nil
}

func genProtocols1(genState *GenState, mainFilePath string, filePathArr []string) error {

	// 主文件的注释
	f, err := os.Open(mainFilePath)
	if err != nil {
		return err
	}
	defer f.Close()
	reader := bufio.NewReader(f)

	var results []string
	// 按行处理txt
	for {
		b, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		line := string(b)
		if strings.HasPrefix(line, "### ") {
			line = strings.Replace(line, "### ", "", 1)
			results = append(results, line)
		}
	}
	genState.Header = strings.Join(results, "\n")

	// 提取 request, response, push, notify
	spaceRegexp := regexp.MustCompile("[^\\s]+")
	for _, filePath := range filePathArr {
		log.Info("处理文件", filePath)
		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			return err
		}
		result := make(map[string]interface{})
		if err := yaml.Unmarshal(data, result); err != nil {
			return err
		}
		for k, v := range result {
			var dict map[string]interface{}
			var ok bool
			if dict, ok = v.(map[string]interface{}); !ok {
				return fmt.Errorf("%s invalid111", k)
			}
			args := spaceRegexp.FindAllString(k, -1)
			name := args[0]
			packet := &Packet{
				Name:       name,
				FullName:   name,
				StructName: camelString(name),
			}
			if c, ok := dict["comment"]; ok {
				packet.Comment = c.(string)
			}
			if len(args) == 1 {
				log.Info("读取协议", name)
				packet.Type = message_type_struct
				if _, ok := dict["struct"]; ok {
					if _, ok := dict["struct"].(map[string]interface{}); ok {
						if err := packet.Message.parse(dict["struct"].(map[string]interface{})); err != nil {
							return err
						}
					}
				}
				packet.Message.StructName = fmt.Sprintf("%s", camelString(name))
				genState.PacketDict[name] = packet
			} else if len(args) == 2 {
				if message_id, err := strconv.Atoi(args[1]); err != nil {
					return err
				} else {
					packet.MessageId = message_id
				}
				log.Info("读取协议", name)
				hasRequest := false
				hasResponse := false
				if v, ok := dict["request"]; ok {
					if v == nil {
					} else if _, ok := v.(map[string]interface{}); !ok {
						return fmt.Errorf("%s invalid66 %#v", k, v)
					}
					hasRequest = true
				}
				if v, ok := dict["response"]; ok {
					if v == nil {
					} else if _, ok := v.(map[string]interface{}); !ok {
						return fmt.Errorf("%s invalid77", k)
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
					if v, ok := dict["request"].(map[string]interface{}); ok {
						if err := packet.Request.parse(v); err != nil {
							return err
						}
					}
					if v, ok := dict["response"].(map[string]interface{}); ok {
						if err := packet.Response.parse(v); err != nil {
							return err
						}
					}
					packet.Request.StructName = fmt.Sprintf("%sRequest", camelString(name))
					packet.Response.StructName = fmt.Sprintf("%sResponse", camelString(name))
				} else if packet.Type == message_type_notify {
					if v, ok := dict["response"].(map[string]interface{}); ok {
						if err := packet.Notify.parse(v); err != nil {
							return err
						}
					}
					packet.Notify.StructName = fmt.Sprintf("%sNotify", camelString(name))
				} else if packet.Type == message_type_push {
					if v, ok := dict["request"].(map[string]interface{}); ok {
						if err := packet.Push.parse(v); err != nil {
							return err
						}
					}
					packet.Push.StructName = fmt.Sprintf("%sPush", camelString(name))
				}
				genState.PacketDict[name] = packet
			}
		}
	}
	// 对packetarr排序
	for _, filePath := range filePathArr {
		f, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer f.Close()
		reader := bufio.NewReader(f)
		// 按行处理txt
		for {
			b, _, err := reader.ReadLine()
			if err == io.EOF {
				break
			}
			line := string(b)
			if len(line) > 0 {
				if !unicode.IsSpace([]rune(line)[0]) && !strings.HasPrefix(line, "#") {
					args := spaceRegexp.FindAllString(line, -1)
					name := args[0]
					name = strings.Replace(name, ":", "", 1)
					if p, ok := genState.PacketDict[name]; ok {
						genState.PacketArr = append(genState.PacketArr, p)
					}
				}
			}
		}
	}
	return nil
}

func genSproto(genState *GenState) error {
	dir := filepath.Join(proj.Config.GenProtocolDir, genState.Module)
	if _, err := os.Stat(proj.Config.GenProtocolDir); os.IsNotExist(err) {
		if err := os.Mkdir(proj.Config.GenProtocolDir, 0755); err != nil {
			return err
		}
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.Mkdir(dir, 0755); err != nil {
			return err
		}
	}
	filePath := filepath.Join(dir, fmt.Sprintf("%s.sproto", genState.Module))
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	log.Info("生成文件", filePath)
	file.Truncate(0)
	defer file.Close()

	funcMap := template.FuncMap{
		"join":        strings.Join,
		"capUpper":    capUpperString,
		"camelString": camelString,
		"quote":       QuoteChar,
	}
	tmpl := template.New("sproto").Delims("<<", ">>")
	tmpl.Funcs(funcMap)
	if tmpl, err := tmpl.Parse(sproto_template); err != nil {
		return err
	} else {
		if err := tmpl.Execute(file, genState); err != nil {
			return err
		}
	}
	return nil
}

func genSprotoGo(genState *GenState) error {
	dir := filepath.Join(proj.Config.SrcGenProtocolDir, genState.Module)
	if _, err := os.Stat(proj.Config.SrcGenProtocolDir); os.IsNotExist(err) {
		if err := os.Mkdir(proj.Config.SrcGenProtocolDir, 0755); err != nil {
			return err
		}
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.Mkdir(dir, 0755); err != nil {
			return err
		}
	}
	filePath := filepath.Join(dir, fmt.Sprintf("%s.go", genState.Module))
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	log.Info("生成文件", filePath)
	file.Truncate(0)
	defer file.Close()

	funcMap := template.FuncMap{
		"join":        strings.Join,
		"capUpper":    capUpperString,
		"capLower":    capLowerString,
		"camelString": camelString,
		"quote":       QuoteChar,
	}
	tmpl := template.New("sproto").Delims("<<", ">>")
	tmpl.Funcs(funcMap)
	if tmpl, err := tmpl.Parse(sproto_go_template); err != nil {
		return err
	} else {
		if err := tmpl.Execute(file, genState); err != nil {
			return err
		}
	}
	return nil
}

// 生成协议
func Gen() error {
	log.Info("===============gen protocol start===============")

	nameArr := make([]string, 0)
	if d, err := os.Open(proj.Config.DocProtocolDir); err != nil {
		return err
	} else {
		if files, err := d.ReadDir(0); err != nil {
			return err
		} else {
			for _, file := range files {
				if !file.IsDir() && filepath.Ext(file.Name()) == ".yaml" {
					log.Info(file.Name())
					nameArr = append(nameArr, strings.Replace(file.Name(), ".yaml", "", 1))
				}
			}
		}
		d.Close()
	}

	for _, name := range nameArr {
		dir := filepath.Join(proj.Config.DocProtocolDir, name)
		protocolFilePathArr := make([]string, 0)
		protocolFilePath := filepath.Join(proj.Config.DocProtocolDir, fmt.Sprintf("%s.yaml", name))
		protocolFilePathArr = append(protocolFilePathArr, protocolFilePath)
		if _, err := os.Stat(dir); err == nil {
			filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
				if d.IsDir() {
					return nil
				}
				if filepath.Ext(path) == ".yaml" {
					protocolFilePathArr = append(protocolFilePathArr, path)
				}
				return nil
			})
		}

		protocolState := &GenState{
			Module:     name,
			PacketDict: make(map[string]*Packet),
			PacketArr:  make([]*Packet, 0),
		}
		if err := genProtocols1(protocolState, protocolFilePath, protocolFilePathArr); err != nil {
			log.Info(err)
			return err
		}
		if err := genSproto(protocolState); err != nil {
			log.Info(err)
			return err
		}
		if err := genSprotoGo(protocolState); err != nil {
			log.Info(err)
			return err
		}
	}
	log.Info("===============gen protocol finished===============")
	return nil
}
