package config

type Database struct {
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	User         string `yaml:"user"`
	Password     string `yaml:"password"`
	Name         string `yaml:"name"`
	MaxIdleConns int    `yaml:"max_idle_conns"`
	MaxOpenConns int    `yaml:"max_open_conns"`
	Prefix       string `yaml:"prefix"`
}

type Redis struct {
	Type         string   `yaml:"type"`
	Host         string   `yaml:"host"`
	Password     string   `yaml:"password"`
	Database     int      `yaml:"database"`
	PoolSize     int      `yaml:"poolSize"`
	MinIdleConns int      `yaml:"minIdleConns"`
	MaxIdleTime  int      `yaml:"maxIdleTime"`
	ClusterHosts []string `yaml:"clusterHosts"`
}

type Google struct {
	Id     string `yaml:"id"`
	Secret string `yaml:"secret"`
}

type Mail struct {
	Host      string `yaml:"host"`
	Port      int    `yaml:"port"`
	From      string `yaml:"from"`
	User      string `yaml:"user"`
	Passwd    string `yaml:"passwd"`
	VerifyUrl string `yaml:"verifyUrl"`
}

type AirWallex struct {
	Id            string `yaml:"id"`
	Key           string `yaml:"key"`
	Url           string `yaml:"url"`
	WebhookSecret string `yaml:"webhookSecret"`
}

type Server struct {
	Env       string    `yaml:"env"`
	Port      int       `yaml:"port"`
	Queue     string    `yaml:"queue"`
	Google    Google    `yaml:"google"`
	Mail      Mail      `yaml:"mail"`
	AirWallex AirWallex `yaml:"airwallex"`
}

type Config struct {
	Database Database `yaml:"database"`
	Redis    Redis    `yaml:"redis"`
	Server   Server   `yaml:"server"`
}
