package core

import (
	"goMedia/config"
	"goMedia/utils"
	"gopkg.in/yaml.v3"
	"log"
)

// InitConfig 从yaml文件中读取配置信息
func InitConfig() *config.Config {
	c := &config.Config{}
	yamlFile, err := utils.LoadYaml()
	if err != nil {
		log.Fatalf("fail to load yaml file: %v", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("fail to unmarshal yaml: %v", err)
	}
	return c
}
