package config

import (
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"log"
)

type Viper struct {
	vip *viper.Viper
}

func (v *Viper) Init() *Viper {
	v.vip = viper.New()
	v.vip.SetConfigName(".env")
	v.vip.SetConfigType("env")
	v.vip.AddConfigPath(".")
	err := v.vip.ReadInConfig()
	if err != nil {
		log.Fatalln("Please init .env file.")
		return v
	}
	v.vip.SetEnvPrefix("appenv")
	v.vip.AutomaticEnv()

	return v
}

func (v *Viper) Map(config map[string]interface{}) map[string]interface{} {
	return config
}

// Env 读取环境变量，支持默认值
func (v *Viper) Env(envName string, defaultValue ...interface{}) interface{} {
	if len(defaultValue) > 0 {
		return v.Get(envName, defaultValue[0])
	}
	return v.Get(envName)
}

// Add 新增配置项
func (v *Viper) Add(name string, configuration map[string]interface{}) {
	v.vip.Set(name, configuration)
}

// Get 获取配置项，允许使用点式获取，如：app.name
func (v *Viper) Get(path string, defaultValue ...interface{}) interface{} {
	// 不存在的情况
	if !v.vip.IsSet(path) {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return nil
	}
	return v.vip.Get(path)
}

// GetString 获取 String 类型的配置信息
func (v *Viper) GetString(path string, defaultValue ...interface{}) string {
	return cast.ToString(v.Get(path, defaultValue...))
}

// GetInt 获取 Int 类型的配置信息
func (v *Viper) GetInt(path string, defaultValue ...interface{}) int {
	return cast.ToInt(v.Get(path, defaultValue...))
}

// GetBool 获取 Bool 类型的配置信息
func (v *Viper) GetBool(path string, defaultValue ...interface{}) bool {
	return cast.ToBool(v.Get(path, defaultValue...))
}
