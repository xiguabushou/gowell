package utils

import (
	"goMedia/global"
	"gopkg.in/yaml.v3"
	"io/fs"
	"os"
)

const configFile = "config.yaml"

// LoadYaml 从文件中加载配置
func LoadYaml() ([]byte, error) {
	return os.ReadFile(configFile)
}

// SaveYaml 将配置保存到文件中
func SaveYaml() error {
	byteData, err := yaml.Marshal(global.Config)
	if err != nil {
		return err
	}
	return os.WriteFile(configFile, byteData, fs.ModePerm)
}
