package config

import "github.com/magiconair/properties"

type Config struct {
	data map[string]string
}

func DefaultConfig() *Config {
	return nil
}
func LoadFile(fileName string,charset string) *Config {
	properties.LoadFile(fileName,properties.UTF8)
	return nil
}
func LoadFiles(fileNames[] string,charset string) *Config {
	properties.LoadFiles(fileNames,properties.UTF8,false)
	return nil
}