package util

import (
	"CodeGenerationGo/template"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

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
