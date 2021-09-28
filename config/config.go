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

func DefaultConfig() *Config {
	return nil
}
func LoadFile(fileName string, charset properties.Encoding) (*Config, error) {
	pro, err := properties.LoadFile(fileName, charset)
	if err == nil {
		var config = &Config{lock:new(sync.RWMutex)}
		config.data = pro.Map()
		return config, err
	}
	return nil, err
}
func LoadFiles(fileNames []string, charset string) *Config {
	properties.LoadFiles(fileNames, properties.UTF8, false)
	return nil
}
func (config *Config) getValue(key string) string {
	config.lock.RLock()
	defer config.lock.RUnlock()
	return config.data[key]
}
func (config *Config) ReadString(key string) string {
	iSection := config.getValue(key)
	return iSection
}
func (config *Config) ReadInt(key string) int {
	v := config.ReadString(key)
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
func (config *Config) ReadBool(key string) bool {
	v := config.ReadString(key)
	if v == "true" {
		return true
	}
	return false
}
