package facades

var Config ConfigFacade

type ConfigFacade interface {
	Map(map[string]interface{}) map[string]interface{}
	Env(envName string, defaultValue ...interface{}) interface{}
	Add(name string, configuration map[string]interface{})
	Get(path string, defaultValue ...interface{}) interface{}
	GetString(path string, defaultValue ...interface{}) string
	GetInt(path string, defaultValue ...interface{}) int
	GetBool(path string, defaultValue ...interface{}) bool
}