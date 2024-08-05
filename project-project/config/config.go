package config

import (
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"log"
	"os"
	"test.com/project-common/logs"
)

var Conf = InitConfig()

type Config struct {
	viper       *viper.Viper
	SC          *ServerConifg
	GC          *GrpcConfig
	EtcdConfig  *EtcdConfig
	MysqlConfig *MysqlConfig
	JwtConfig   *JwtConfig
}

type ServerConifg struct {
	Name string
	Addr string
}

type GrpcConfig struct {
	Name    string
	Addr    string
	Version string
	Weight  int64
}

type EtcdConfig struct {
	Addrs []string
}

type MysqlConfig struct {
	Username string
	Password string
	Host     string
	Port     int
	Db       string
}

type JwtConfig struct {
	AccessSecret  string
	AccessExp     int
	RefreshSecret string
	RefreshExp    int
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
	config.GetGrpcConfig()
	config.ReadEtcdConfig()
	config.ReadMysqlConfig()
	config.ReadJwtConfig()
	return config
}

// GetServerConfig 读取Server配置；被InitConfig()调用；Get不合适，毕竟是无返回值的，最好叫Read
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

func (c *Config) GetRedisConfig() *redis.Options {
	return &redis.Options{
		Addr:     c.viper.GetString("redis.host") + ":" + c.viper.GetString("redis.port"),
		Password: c.viper.GetString("redis.password"),
		DB:       c.viper.GetInt("redis.db"),
	}
}

// GetGrpcConfig Get不合适，毕竟是无返回值的，最好叫Read
func (c *Config) GetGrpcConfig() {
	gc := &GrpcConfig{}
	gc.Name = c.viper.GetString("grpc.name")
	gc.Addr = c.viper.GetString("grpc.addr")
	gc.Version = c.viper.GetString("grpc.version")
	gc.Weight = c.viper.GetInt64("grpc.weight")
	c.GC = gc
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

func (c *Config) ReadMysqlConfig() {
	mc := &MysqlConfig{
		Username: c.viper.GetString("mysql.username"),
		Password: c.viper.GetString("mysql.password"),
		Host:     c.viper.GetString("mysql.host"),
		Port:     c.viper.GetInt("mysql.port"),
		Db:       c.viper.GetString("mysql.db"),
	}
	c.MysqlConfig = mc
}

func (c *Config) ReadJwtConfig() {
	jc := &JwtConfig{}
	jc.AccessSecret = c.viper.GetString("jwt.accessSecret")
	jc.AccessExp = c.viper.GetInt("jwt.accessExp")
	jc.RefreshSecret = c.viper.GetString("jwt.refreshSecret")
	jc.RefreshExp = c.viper.GetInt("jwt.refreshExp")
	c.JwtConfig = jc
}
