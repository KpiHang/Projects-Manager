package config

import (
	"github.com/spf13/viper"
	"log"
	"os"
	"test.com/project-common/logs"
)

var Conf = InitConfig()

type Config struct {
	viper      *viper.Viper
	SC         *ServerConifg
	GC         *GrpcConfig
	EtcdConfig *EtcdConfig
}

type ServerConifg struct {
	Name string
	Addr string
}

type GrpcConfig struct {
	Name string
	Addr string
}

type EtcdConfig struct {
	Addrs []string
}

func InitConfig() *Config {
	config := &Config{viper: viper.New()}
	workDir, _ := os.Getwd() // 返回的是你运行命令时所在的工作目录，而不是 main.go 文件所在的目录。
	config.viper.SetConfigName("config")
	config.viper.SetConfigType("yaml")
	config.viper.AddConfigPath("/etc/ms_project/user") // 可以添加多路径；这个路径不存在，仅test；
	config.viper.AddConfigPath(workDir + "/config")
	err := config.viper.ReadInConfig()
	if err != nil {
		log.Fatalln(err)
	}
	config.GetServerConfig() // 读取Server配置；
	config.InitZapLog()      // 加载日志配置；加载=读取+用；直接用上了；
	config.ReadEtcdConfig()
	return config
}

// GetServerConfig 读取Server配置；被InitConfig()调用
func (c *Config) GetServerConfig() {
	sc := &ServerConifg{}
	sc.Name = c.viper.GetString("server.name")
	sc.Addr = c.viper.GetString("server.Addr")
	c.SC = sc
}

func (c *Config) InitZapLog() {
	lc := &logs.LogConfig{
		DebugFileName: c.viper.GetString("zap.debugFileName"),
		InfoFileName:  c.viper.GetString("zap.infoFileName"),
		WarnFileName:  c.viper.GetString("zap.warnFileName"),
		MaxSize:       c.viper.GetInt("zap.maxSize"),
		MaxAge:        c.viper.GetInt("zap.maxAge"),
		MaxBackups:    c.viper.GetInt("zap.MaxBackups"),
	}

	err := logs.InitLogger(lc)
	if err != nil {
		log.Fatalln(err)
	}
}

func (c *Config) ReadEtcdConfig() {
	ec := &EtcdConfig{}
	var addrs []string
	err := c.viper.UnmarshalKey("etcd.Addrs", &addrs)
	if err != nil {
		log.Fatalln(err)
	}
	ec.Addrs = addrs
	c.EtcdConfig = ec
}
