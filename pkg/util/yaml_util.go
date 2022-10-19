package util

import (
	template2 "CodeGenerationGo/pkg/template"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	v1 "k8s.io/api/core/v1"
	"log"
	"os"
)

func ReadmatchlsrFile(path string) (*template2.LabelSelectorRequirement, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var conf template2.LabelSelectorRequirement
	err = yaml.Unmarshal(buf, &conf)
	if err != nil {
		return nil, fmt.Errorf("in file %q: %v", path, err)
	}

	return &conf, nil
}

func ReadAffinityYamlFile(path string) (*template2.Affinity, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var conf template2.Affinity
	err = yaml.Unmarshal(buf, &conf)
	if err != nil {
		return nil, fmt.Errorf("in file %q: %v", path, err)
	}

	return &conf, nil
}

func ReadPodYamlFile(path string) (*template2.Pod, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var conf template2.Pod
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

func WriteObject2Yaml(in interface{}, outPath string) {
	yamlByte, _ := yaml.Marshal(in)

	if err := os.WriteFile(outPath, yamlByte, 0666); err != nil {
		log.Fatal(err)
	}
}

func DeleteFile(filepath string) {
	err := os.Remove(filepath)
	if err != nil {
		log.Fatal(err)
	}
}
