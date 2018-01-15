package config

import (
	"io/ioutil"

	"github.com/THE108/scrapper/info"

	"gopkg.in/yaml.v2"
)

type Item struct {
	Url             string                          `yaml:"url"`
	DeliveryOptions map[string]info.DeliveryDetails `yaml:"delivery"`
}

func ParseConfig(filename string) ([]Item, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var result []Item
	if err := yaml.Unmarshal(content, &result); err != nil {
		return nil, err
	}

	return result, nil
}
