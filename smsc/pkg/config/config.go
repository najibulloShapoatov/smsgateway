package config

import (
	"smsc/pkg/db"
	"smsc/pkg/log"
	"strconv"

	"github.com/sasbury/mini"
)

//KannelConfig ...
type KannelConfig struct {
	Link     string
	UserName string
	Password string
	SmsC     string
}

var config *mini.Config
var err error

//Init configs from file
func Init(filepath string) {

	config, err = mini.LoadConfiguration(filepath)
	if err != nil {
		log.Error("error not get conf >>", err)
		panic(err)
	}
}

//LoadDBConfigs func
func LoadDBConfigs() *db.Config {

	SectionName := "DB"
	return &db.Config{
		Driver:          config.String("DB_DRIVER", "pgx"),
		ApplicationName: config.String("APPLICATION_NAME", ""),
		Host:            loadStringFromSection(SectionName, config, "Host", "127.0.0.1"),
		Port:            loadStringFromSection(SectionName, config, "Port", "5432"),
		Dbname:          loadStringFromSection(SectionName, config, "Dbname", ""),
		SslMode:         loadStringFromSection(SectionName, config, "SslMode", ""),
		User:            loadStringFromSection(SectionName, config, "User", ""),
		Pass:            loadStringFromSection(SectionName, config, "Pass", ""),
		ConnMaxLifetime: loadIntFromSection(SectionName, config, "ConnMaxLifetime", "30000"),
		MaxOpenConns:    loadIntFromSection(SectionName, config, "MaxOpenConns", "16"),
		MaxIdleConns:    loadIntFromSection(SectionName, config, "MaxIdleConns", "8"),
	}
}

//GetHTTPPort  ....
func GetHTTPPort() string {
	SectionName := "HTTP_SERVER"
	return loadStringFromSection(SectionName, config, "PORT", "")
}

//GetHTTPPort  ....
func GetHTTPHost() string {
	SectionName := "HTTP_SERVER"
	return loadStringFromSection(SectionName, config, "HTTPHost", "")
}

//GetKannelCfg ...
func GetKannelCfg() *KannelConfig {
	sName := "KANNEL"
	return &KannelConfig{
		Link:     loadStringFromSection(sName, config, "LINK", ""),
		UserName: loadStringFromSection(sName, config, "USERNAME", ""),
		Password: loadStringFromSection(sName, config, "PASSWORD", ""),
		SmsC:     loadStringFromSection(sName, config, "SMSC", ""),
	}
}

// loadIntFromSection load int paparameter and log err
func loadIntFromSection(sectionName string, pgcfg *mini.Config, name string, defval string) int {
	strVal := pgcfg.StringFromSection(sectionName, name, defval)
	if defval == "" && strVal == "" {
		log.ErrorDepth("Missing mandatory: Section, Parameter", 1, sectionName, name)
		return 0
	}
	intVal, err := strconv.Atoi(strVal)
	if err != nil {
		log.ErrorDepth("ncorrect integer: Section, Parameter, Value", 1, err, sectionName, name, strVal)
		return 0
	}

	//log.InfoDepth("Load config parameter: Section, Parameter, Value", 1, sectionName, name, intVal)

	return abs(intVal)
}

// loadStringFromSection load str paparameter and log err
func loadStringFromSection(sectionName string, pgcfg *mini.Config, name string, defval string) string {
	strVal := pgcfg.StringFromSection(sectionName, name, defval)
	if defval == "" && strVal == "" {
		log.ErrorDepth("Missing mandatory: Section, Parameter", 1, sectionName, name)
		return ""
	}
	//log.InfoDepth("Load config parameter: Section, Parameter, Value", 1, sectionName, name, strVal)

	return strVal
}

// loadBoolFromSection load bool paparameter and log err
func loadBoolFromSection(sectionName string, pgcfg *mini.Config, name string, defval string) bool {
	var boolVal bool
	strVal := pgcfg.StringFromSection(sectionName, name, defval)
	if defval == "" && strVal == "" {
		log.ErrorDepth("Missing mandatory: Section, Parameter", 1, sectionName, name)
		return false
	}

	if strVal != "" {
		switch strVal {
		case "true":
			boolVal = true
		case "false":
			boolVal = false
		default:
			log.ErrorDepth("Incorrect boolean, оnly avaliable: 'true', 'false': Section, Parameter, Value", 1, sectionName, name, strVal)
			return false
		}
	}

	log.InfoDepth("Load config parameter: Section, Parameter, Value", 1, sectionName, name, boolVal)

	return boolVal
}

func abs(a int) int {
	if a < 0 {
		return a * -1
	}
	return a
}
