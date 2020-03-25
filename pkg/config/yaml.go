package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

var ConfigData Config

type Config struct {
	HttpPort 		string `yaml:"httpPort"`
	JwtEncrtpy 		string `yaml:"jwtEncrtpy"`
	LoginExpire 	string `yaml:"loginExpire"`

	DbType 			string `yaml:"dbType"`
	DbUser 			string `yaml:"dbUser"`
	DbHost 			string `yaml:"dbHost"`
	DbPassword 		string `yaml:"dbPassword"`
	DbName 			string `yaml:"dbName"`

	RedisAddr 		string `yaml:"redisAddr"`
	RedisPassword 	string `yaml:"redisPassword"`
	RedisDefaultDB 	int `yaml:"redisDefaultDB"`
}

func Setup(path string){
	yamlFile, readErr := ioutil.ReadFile(path)
	if readErr != nil {
		fmt.Errorf("Config File Loaded Failed ###", readErr.Error())
		return
	}
	err := yaml.Unmarshal(yamlFile, &ConfigData)
	if err != nil {
		fmt.Errorf("Config Setup Failed ###", path )
		return
	}

	fmt.Println("Config Setup Success ...")
}