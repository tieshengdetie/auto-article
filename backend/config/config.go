package config

// DataSourceConfig 数据库的结构体
type DataSourceConfig struct {
	DriverName string `mapstructure:"driverName"`
	Host       string `mapstructure:"host"`
	Port       int    `mapstructure:"port"`
	Database   string `mapstructure:"database"`
	Username   string `mapstructure:"username"`
	Password   string `mapstructure:"password"`
	Charset    string `mapstructure:"charset"`
	Loc        string `mapstructure:"loc"`
}

// ServerConfig 整个项目的配置
type ServerConfig struct {
	Port          int               `mapstructure:"port"`
	Salt          string            `mapstructure:"salt"`
	Deploy        string            `mapstructure:"deploy"`
	DataSource    DataSourceConfig  `mapstructure:"datasource"`
	Kafka         KafkaConfig       `mapstructure:"kafka"`
	Redis         RedisConfig       `mapstructure:"redis"`
	AliYunModel   AliYunModelConfig `mapstructure:"aliYunModel"`
	RabbitMq      RabbitMqConfig    `mapstructure:"rabbitMq"`
	TianApiConfig TianApiConfig     `mapstructure:"tianApiConfig"`
}
type RabbitMqConfig struct {
	Host     string `mapstructure:"host"`
	UserName string `mapstructure:"userName"`
	Password string `mapstructure:"password"`
}
type KafkaConfig struct {
	Brokers  string `mapstructure:"brokers"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
}

// DataWarehouseConfig 数仓配置
type DataWarehouseConfig struct {
	BaseUrl   string `mapstructure:"baseUrl"`
	AppKey    string `mapstructure:"appKey"`
	AppSecret string `mapstructure:"appSecret"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Password string `mapstructure:"password"`
}

type AliYunModelConfig struct {
	BaseUrl string `mapstructure:"baseUrl"`
	AppKey  string `mapstructure:"appKey"`
}

type TianApiConfig struct {
	BaseUrl string `mapstructure:"baseUrl"`
	AppKey  string `mapstructure:"appKey"`
}
