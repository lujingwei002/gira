package configs
import (
	yaml "gopkg.in/yaml.v3"
	"io/ioutil"
	"path/filepath"
)

type Config struct {
    ErrorCodeArr ErrorCodeArr
    FlowReasonArr FlowReasonArr
    AchievementTaskArr AchievementTaskArr
}

func (self *Config) LoadFromYaml(dir string) error {
	var filePath string

	filePath = filepath.Join(dir, "AchievementTask.yaml")
	if err := self.AchievementTaskArr.Load(filePath); err != nil {
		return err
	}
	filePath = filepath.Join(dir, "ErrorCode.yaml")
	if err := self.ErrorCodeArr.Load(filePath); err != nil {
		return err
	}
	filePath = filepath.Join(dir, "FlowReason.yaml")
	if err := self.FlowReasonArr.Load(filePath); err != nil {
		return err
	}
	return nil
}

type ErrorCode struct {
    Id int `yaml:id`
    Key string `yaml:key`
    Value string `yaml:value`
}

type ErrorCodeArr []ErrorCode

func (self *ErrorCodeArr) Load(filePath string) error {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(data, self); err != nil {
		return err
	}
	return nil
}


type FlowReason struct {
    Id int `yaml:id`
    Key string `yaml:key`
}

type FlowReasonArr []FlowReason

func (self *FlowReasonArr) Load(filePath string) error {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(data, self); err != nil {
		return err
	}
	return nil
}


type AchievementTask struct {
    Id int `yaml:id`
    GroupId int `yaml:group_id`
    Order int `yaml:order`
    Hook string `yaml:hook`
    CompleteType int `yaml:complete_type`
    Arg1 int `yaml:arg1`
    Arg2 interface{} `yaml:arg2`
    Count int `yaml:count`
    Bonus interface{} `yaml:bonus`
}

type AchievementTaskArr []AchievementTask

func (self *AchievementTaskArr) Load(filePath string) error {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(data, self); err != nil {
		return err
	}
	return nil
}

