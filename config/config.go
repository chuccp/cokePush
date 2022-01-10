package config

import (
	"github.com/magiconair/properties"
	"log"
	"strconv"
	"sync"
)

type Config struct {
	data map[string]string
	lock *sync.RWMutex
}

var defaultConfig = &Config{lock: new(sync.RWMutex),data:make(map[string]string)}

func DefaultConfig() *Config {
	return defaultConfig
}
func LoadFile(fileName string, charset properties.Encoding) (*Config, error) {
	pro, err := properties.LoadFile(fileName, charset)
	if err == nil {
		var config = &Config{lock: new(sync.RWMutex)}
		config.data = pro.Map()
		return config, err
	}
	return nil, err
}
func LoadFiles(fileNames []string, charset properties.Encoding) (*Config, error) {
	pro, err := properties.LoadFiles(fileNames, charset, false)
	if err == nil {
		var config = &Config{lock: new(sync.RWMutex)}
		config.data = pro.Map()
		return config, err
	}
	return nil, err
}
func (config *Config) getValue(key string) string {
	config.lock.RLock()
	defer config.lock.RUnlock()
	v := config.data[key]
	return v
}
func (config *Config) GetString(key string) string {
	iSection := config.getValue(key)
	return iSection
}
func (config *Config) GetInt(key string) int {
	v := config.GetString(key)
	if len(v) > 0 {
		iv, err := strconv.Atoi(v)
		if err == nil {
			return iv
		} else {
			log.Panic("格式错误  key:", key)
		}
	}
	return -1
}
func (config *Config) GetBool(key string) bool {
	v := config.GetString(key)
	if v == "true" {
		return true
	}
	return false
}
func (config *Config) GetBoolOrDefault(key string, defaultValue bool) bool {
	v := config.GetString(key)
	if len(v) == 0 {
		return defaultValue
	}
	if v == "true" {
		return true
	}
	return false
}

func (config *Config) SetString(key string, value string) {
	config.lock.Lock()
	config.data[key] = value
	config.lock.Unlock()
}

func (config *Config) GetStringOrDefault(key string, defaultValue string) string {
	v := config.getValue(key)
	if len(v) == 0 {
		return defaultValue
	}
	return v
}
func (config *Config) GetIntOrDefault(key string, defaultValue int) int {
	v := config.GetString(key)
	if len(v) == 0 {
		return defaultValue
	}
	if len(v) > 0 {
		iv, err := strconv.Atoi(v)
		if err == nil {
			return iv
		} else {
			log.Panic("格式错误  key:", key)
		}
	}
	return -1
}
