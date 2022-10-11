package util

import (
	"CodeGenerationGo/template"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	v1 "k8s.io/api/core/v1"
)

func ReadmatchlsrFile(path string) (*template.LabelSelectorRequirement, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var conf template.LabelSelectorRequirement
	err = yaml.Unmarshal(buf, &conf)
	if err != nil {
		return nil, fmt.Errorf("in file %q: %v", path, err)
	}

	return &conf, nil
}

func ReadAffinityYamlFile(path string) (*template.Affinity, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var conf template.Affinity
	err = yaml.Unmarshal(buf, &conf)
	if err != nil {
		return nil, fmt.Errorf("in file %q: %v", path, err)
	}

	return &conf, nil
}

func ReadPodYamlFile(path string) (*template.Pod, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var conf template.Pod
	err = yaml.Unmarshal(buf, &conf)
	if err != nil {
		return nil, fmt.Errorf("in file %q: %v", path, err)
	}

	return &conf, nil
}

func ReadYamlFile[T any](path string) (*T, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var conf T
	err = yaml.Unmarshal(buf, &conf)
	if err != nil {
		return nil, fmt.Errorf("in file %q: %v", path, err)
	}

	return &conf, nil
}
func FetchYaml(sourceYaml []byte) (*yaml.Node, error) {
	rootNode := yaml.Node{}
	err := yaml.Unmarshal(sourceYaml, &rootNode)
	if err != nil {
		return nil, err
	}
	return &rootNode, nil
}

func ReadAffinityJson(path string) (*v1.Affinity, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var conf v1.Affinity
	err = json.Unmarshal(buf, &conf)
	if err != nil {
		return nil, fmt.Errorf("in file %q: %v", path, err)
	}

	return &conf, nil
}
