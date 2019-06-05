package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strings"

	"gopkg.in/yaml.v2"
)

type config struct {
	Debug               bool   `yaml:"debug"`
	Port                string `yaml:"port"`
	Hostname            string `yaml:"hostname"`
	GitUser             string `yaml:"git_user"`
	GitPassword         string `yaml:"git_password"`
	GitTmpDir           string `yaml:"git_tmp_dir"`
	DbPath              string `yaml:"db_path"`
	DefaultGitServer    string `yaml:"default_git_server"`
	DefaultGitNamespace string `yaml:"default_git_namespace"`
}

type settings struct {
	GitServer    string `json:"git_server"`
	GitNamespace string `json:"git_namespace"`
}

var defaultConfig = config{
	Debug:               true,
	Port:                "8080",
	Hostname:            "",
	GitUser:             "",
	GitPassword:         "",
	GitTmpDir:           "/tmp/nebular/repos",
	DbPath:              "/tmp/nebular",
	DefaultGitServer:    "http://github.com",
	DefaultGitNamespace: "robrotheram",
}
var configPath = "config/config.yaml"
var configSamplepath = "config/config.sample.yaml"

var configuration config

func configFromFile() {
	//Create Sample config if it does not exists
	if ex, err := exists(configSamplepath); err != nil || !ex {
		fmt.Println(defaultConfig)
		y, err := yaml.Marshal(defaultConfig)
		if err != nil {
			fmt.Println(err)
		}
		f, err := os.Create(configSamplepath)
		defer f.Close()
		n2, err := f.Write(y)
		fmt.Printf("wrote %d bytes\n", n2)
		f.Sync()
	}

	if ex, err := exists(configPath); err == nil && ex {
		yamlFile, err := ioutil.ReadFile(configPath)
		if err != nil {
			log.Printf("yamlFile.Get err   #%v ", err)
		}
		err = yaml.Unmarshal(yamlFile, configuration)
		if err != nil {
			log.Fatalf("Unmarshal: %v", err)
		}
	} else {
		fmt.Println("No Configfile found using defaults")
		configuration = defaultConfig
	}
}

func configFromEnv() {
	t := reflect.TypeOf(configuration)
	// Iterate over all available fields and read the tag value
	for i := 0; i < t.NumField(); i++ {
		// Get the field, returns https://golang.org/pkg/reflect/#StructField
		field := t.Field(i)
		// Get the field tag value
		tag := field.Tag.Get("yaml")
		envName := fmt.Sprintf("NEBULAR_%v", strings.ToUpper(tag))
		if configuration.Debug {
			fmt.Println("Looking for ENV: " + envName)
		}
		value, found := os.LookupEnv(envName)
		if found {
			v := reflect.ValueOf(&configuration).Elem().FieldByName(field.Name)
			if v.IsValid() {
				v.SetString(value)
			}
		}
	}
}
