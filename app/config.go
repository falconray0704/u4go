package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

// Config represents the structure of the yaml file
type Config struct {
	DBsInfos map[string]DBInfo		`yaml:"dbs_infos"`
	ZapLog	ZapLogInfos				`yaml:"zap_log"`
}

type DBInfo struct {
	UserName	string	`yaml:"user_name"`
	Password	string	`yaml:"password"`
	Url 		string	`yaml:"url"`
	Port 		int16	`yaml:"port"`
}

// Zap details for zap
type ZapLogInfos struct {
	LogsPath	string		`yaml:"logsPath"`
	ConsoleFileOut string	`yaml:"consoleFileOut"`
	ConsoleFileErr string	`yaml:"consoleFileErr"`
	JsonFileOut string		`yaml:"jsonFileOut"`
	JsonFileErr string		`yaml:"jsonFileErr"`
}


// App Config factory takes a path to a yaml file and produces a parsed Config
func NewConfig(path string) (*Config, error) {
	var c Config

	data, err := ioutil.ReadFile(path)
	//fmt.Printf("rawDatas:%s\n\n", data)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(data, &c); err != nil {
		return nil, err
	}

	return &c, err
}
