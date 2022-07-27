package util

import (
	"CodeGenerationGo/Template"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
)

func ReadYamlConfig(path string) (*Template.MatchExpression, error) {
	conf := &Template.MatchExpression{}
	if f, err := os.Open(path); err != nil {
		return nil, err
	} else {
		fmt.Println(f)
		content, _ := ioutil.ReadAll(f)
		fmt.Println(string(content))
		yaml.NewDecoder(f).Decode(conf)
	}
	fmt.Println("conf: ", conf)

	return conf, nil
}

func ReadConfigYaml(path string) (*Template.Config, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var conf Template.Config
	err = yaml.Unmarshal(buf, &conf)
	if err != nil {
		return nil, fmt.Errorf("in file %q: %v", path, err)
	}
	return &conf, nil
}
