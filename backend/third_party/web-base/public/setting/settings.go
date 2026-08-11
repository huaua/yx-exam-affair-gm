package setting

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

var Conf = new(AppConfig)

type AppConfig struct {
	Name            string `mapstructure:"name"`
	Mode            string `mapstructure:"mode"`
	Version         string `mapstructure:"version"`
	Port            int    `mapstructure:"port"`
	Host            string `mapstructure:"host"`
	IsDocker        bool
	UiInside        string `mapstructure:"ui_inside"`    // 前端资源放入服务端访问，"embed":编译时打包嵌入到可执行程序内部；"outdir":访问静态资源目录
	OpenBrowser     bool   `mapstructure:"open_browser"` // server服务启动后，自动打开浏览器访问系统
	*TokenConfig    `mapstructure:"token"`
	*LogConfig      `mapstructure:"log"`
	*Datasource     `mapstructure:"datasource"`
	*UploadFile     `mapstructure:"upload_file"`
	*SecurityConfig `mapstructure:"security"`
}

type SecurityConfig struct {
	EnablePprof       bool     `mapstructure:"enable_pprof"`
	AllowedOrigins    []string `mapstructure:"allowed_origins"`
	TrustedProxies    []string `mapstructure:"trusted_proxies"`
	MaxRequestBodyMB  int64    `mapstructure:"max_request_body_mb"`
	MaxDownloadMB     int64    `mapstructure:"max_download_mb"`
	ReadTimeoutSec    int      `mapstructure:"read_timeout_sec"`
	WriteTimeoutSec   int      `mapstructure:"write_timeout_sec"`
	IdleTimeoutSec    int      `mapstructure:"idle_timeout_sec"`
	JWTMaxLifetimeMin int64    `mapstructure:"jwt_max_lifetime_minutes"`
	LoginMaxAttempts  int      `mapstructure:"login_max_attempts"`
	LoginWindowMin    int      `mapstructure:"login_window_minutes"`
}

type TokenConfig struct {
	ExpireTime int64  `mapstructure:"expire_time"`
	Secret     string `mapstructure:"secret"`
	Issuer     string `mapstructure:"issuer"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
}

type Datasource struct {
	Master   *Master   `mapstructure:"master"`
	Slave    *Slave    `mapstructure:"slave"`
	Redis    *Redis    `mapstructure:"redis"`
	RabbitMQ *RabbitMQ `mapstructure:"rabbitmq"`
	Kafka    *Kafka    `mapstructure:"kafka"`
	Mongodb  *Mongodb  `mapstructure:"mongodb"`
}
type Master struct {
	DriverName   string `mapstructure:"driver_name"`
	Host         string `mapstructure:"host"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	DB           string `mapstructure:"dbname"`
	Port         int    `mapstructure:"port"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}

type Slave struct {
	Count        int      `mapstructure:"count"`
	DriverName   string   `mapstructure:"driver_name"`
	Hosts        []string `mapstructure:"hosts"`
	Users        []string `mapstructure:"users"`
	Passwords    []string `mapstructure:"passwords"`
	DBs          []string `mapstructure:"dbnames"`
	Ports        []int    `mapstructure:"ports"`
	MaxOpenConns int      `mapstructure:"max_open_conns"`
	MaxIdleConns int      `mapstructure:"max_idle_conns"`
}

type Redis struct {
	Enabled      bool   `mapstructure:"enabled"`
	Host         string `mapstructure:"host"`
	Password     string `mapstructure:"password"`
	Port         int    `mapstructure:"port"`
	DB           int    `mapstructure:"db"`
	PoolSize     int    `mapstructure:"pool_size"`
	MinIdleConns int    `mapstructure:"min_idle_conns"`
}

type RabbitMQ struct {
	Enabled  bool   `mapstructure:"enabled"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
}

type Kafka struct {
	Enabled bool   `mapstructure:"enabled"`
	Addrs   string `mapstructure:"addrs"`
}

type Mongodb struct {
	Enabled  bool   `mapstructure:"enabled"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	DBName   string `mapstructure:"dbname"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
}

func Init(filePath string) {
	// 方式1：直接指定配置文件路径（相对路径或者绝对路径）
	// 相对路径：相对执行的可执行文件的相对路径
	//viper.SetConfigFile("./conf/config.yaml")
	// 绝对路径：系统中实际的文件路径
	//viper.SetConfigFile("/Users/xxx/Desktop/bluebell/conf/config.yaml")

	// 方式2：指定配置文件名和配置文件的位置，viper自行查找可用的配置文件
	// 配置文件名不需要带后缀
	// 配置文件位置可配置多个
	//viper.SetConfigName("config") // 指定配置文件名（不带后缀）
	//viper.AddConfigPath(".") // 指定查找配置文件的路径（这里使用相对路径）
	//viper.AddConfigPath("./conf")      // 指定查找配置文件的路径（这里使用相对路径）

	// 基本上是配合远程配置中心使用的，告诉viper当前的数据使用什么格式去解析
	//viper.SetConfigType("json")

	viper.SetConfigFile(filePath)

	err := viper.ReadInConfig() // 读取配置信息
	if err != nil {
		// 读取配置信息失败
		fmt.Printf("viper.ReadInConfig failed, err:%v\n", err)
		panic(err)
	}

	// 把读取到的配置信息反序列化到 Conf 变量中
	if err = viper.Unmarshal(Conf); err != nil {
		fmt.Printf("viper.Unmarshal failed, err:%v\n", err)
		panic(err)
	}
	applySecurityDefaults()
	applyEnvironmentOverrides()

	if _, err = os.Stat("/.dockerenv"); err == nil {
		Conf.IsDocker = true
	} else {
		Conf.IsDocker = false
	}

	// viper.WatchConfig()
	// viper.OnConfigChange(func(in fsnotify.Event) {
	// 	fmt.Println("配置文件修改了...")
	// 	if err := viper.Unmarshal(Conf); err != nil {
	// 		fmt.Printf("viper.Unmarshal failed, err:%v\n", err)
	// 		panic(err)
	// 	}
	// })
}

func applySecurityDefaults() {
	if Conf.SecurityConfig == nil {
		Conf.SecurityConfig = &SecurityConfig{}
	}
	if Conf.MaxRequestBodyMB <= 0 {
		Conf.MaxRequestBodyMB = 20
	}
	if Conf.MaxDownloadMB <= 0 {
		Conf.MaxDownloadMB = 100
	}
	if Conf.ReadTimeoutSec <= 0 {
		Conf.ReadTimeoutSec = 15
	}
	if Conf.WriteTimeoutSec <= 0 {
		Conf.WriteTimeoutSec = 60
	}
	if Conf.IdleTimeoutSec <= 0 {
		Conf.IdleTimeoutSec = 60
	}
	if Conf.JWTMaxLifetimeMin <= 0 {
		Conf.JWTMaxLifetimeMin = 480
	}
	if Conf.LoginMaxAttempts <= 0 {
		Conf.LoginMaxAttempts = 5
	}
	if Conf.LoginWindowMin <= 0 {
		Conf.LoginWindowMin = 10
	}
}

func applyEnvironmentOverrides() {
	if v := os.Getenv("GM_DB_PASSWORD"); v != "" {
		Conf.Datasource.Master.Password = v
	}
	if v := os.Getenv("GM_TOKEN_SECRET"); v != "" {
		Conf.TokenConfig.Secret = v
	}
	if v := os.Getenv("GM_ENABLE_PPROF"); v != "" {
		Conf.EnablePprof, _ = strconv.ParseBool(v)
	}
	if v := os.Getenv("GM_ALLOWED_ORIGINS"); v != "" {
		Conf.AllowedOrigins = strings.Split(v, ",")
	}
}
