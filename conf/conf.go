package conf

import (
	"log"

	"github.com/spf13/viper"
)

type AppConf struct {
	Kafka  KafkaConf
	Etcd   EtcdConf
	Center CenterConf
}

type KafkaConf struct {
	Addr        string
	Topic       string
	ChanMaxSize int
}

type EtcdConf struct {
	Address string
	Timeout int
	Key     string
}

type CenterConf struct {
	Addr string
}

var (
	Appconf AppConf
)

func init() {
	// 初始化 Viper。
	viper.SetConfigName("conf")    // 配置文件名称 (无扩展名)
	viper.SetConfigType("yaml")    // 如果配置文件的内容是 JSON，这里需要改为 "json"
	viper.AddConfigPath(".")       // 查找配置文件所在的路径，这里是项目根目录
	viper.AddConfigPath("./conf/") // 设置另一个可能的配置路径

	// viper.SetConfigFile("/path/to/config.yaml")

	// 读取配置文件。
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
	Appconf = AppConf{
		Kafka: KafkaConf{
			Addr:        viper.GetString("kafka.address"),
			Topic:       viper.GetString("kafka.topic"),
			ChanMaxSize: viper.GetInt("kafka.chan_max_size"),
		},
		Etcd: EtcdConf{
			Address: viper.GetString("etcd.address"),
			Timeout: viper.GetInt("etcd.timeout"),
			Key:     viper.GetString("etcd.collect_log_key"),
		},
		Center: CenterConf{
			Addr: viper.GetString("center.address"),
		},
	}

	// 监听和重新加载配置文件
	// viper.WatchConfig()
}
