package util

import (
	"CodeGenerationGo/template"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
)

func ReadYamlConfig(path string) (*template.MatchExpression, error) {
	conf := &template.MatchExpression{}
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

func ReadConfigYaml(path string) (*template.Config, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var conf template.Config
	err = yaml.Unmarshal(buf, &conf)
	if err != nil {
		return nil, fmt.Errorf("in file %q: %v", path, err)
	}
	return &conf, nil
}
