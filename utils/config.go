package utils

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
	"reflect"
)

var Config *Configuration

const CONFIGPATH = "conf/conf.properties"

type Configuration struct {
	Logpath      string
	Loglevel     string
	Accesslog    string
	DataDriver   string
	DataSource   string
	Ip           string
	User         string
	Passwd       string
	Port         string
	Dbbase       string
	Proportion   string
}
func SetupConfig() {
	Config, _ = LoadConfig()
	fmt.Println(Config)
	Logger := GetLog()
	Logger.Infoln("SetupConfig")
	fmt.Println("SetupConfig")
}

func DefaultConfiguration() *Configuration {
	cfg := &Configuration{
		Logpath:  "/var/log/catkeeper/api.log",
		Loglevel: "info",
		Accesslog:  "/var/log/catkeeper/access.log",
		DataDriver:   "mysql",
		DataSource:   "root:@tcp(10.72.84.145:3306)/catkeeper",
		Ip: "10.72.84.145",
		User: "root",
		Passwd: "root",
		Port: "3306",
		Dbbase: "catkeeper",
		Proportion: "6",
	}
	return cfg
}

func LoadConfig() (*Configuration, error) {
	rtConfig := DefaultConfiguration()
	if _, err := os.Stat(CONFIGPATH); err != nil {
		fmt.Fprintln(os.Stderr,"config file does exsit,skipped config file")
	} else {
		_, err = toml.DecodeFile(CONFIGPATH, &rtConfig)
		if err != nil {
			fmt.Fprintln(os.Stderr,"failed to decode config file,skipped config file", err)
		}
	}
	mergeConfig(rtConfig, configFromFlag())
	fmt.Println(rtConfig);
	return rtConfig, nil
}

func configFromFlag() *Configuration {
	cfg := &Configuration{}
	
	flag.StringVar(&cfg.Logpath, "Logpath", "", "path for the log file")
	flag.StringVar(&cfg.Loglevel, "Loglevel", "", "using standard go library")
	flag.StringVar(&cfg.Accesslog, "Accesslog", "", "path for access file")
	flag.StringVar(&cfg.DataDriver, "DataDriver", "", "database driver")
	flag.StringVar(&cfg.DataSource, "DataSource", "", "using standard mysql datasource")
	flag.StringVar(&cfg.Ip, "Ip", "", "the database ip")
	flag.StringVar(&cfg.User, "User", "", "the database user")
	flag.StringVar(&cfg.Passwd, "Passwd", "", "the database password")
	flag.StringVar(&cfg.Port, "Port", "", "the database Port")
	flag.StringVar(&cfg.Dbbase, "Dbbase", "", "the database database")
	flag.StringVar(&cfg.Proportion, "Proportion", "", "the Sold proportion")
		
	flag.Parse()
	return cfg
}

func mergeConfig(defaultcfg, filecfg interface{}) {
	v1 := reflect.ValueOf(filecfg).Elem()
	v := reflect.ValueOf(defaultcfg).Elem()
	mergeValue(v, v1)
}

func mergeValue(v, v1 reflect.Value) {
	for i := 0; i < v.NumField(); i++ {
		switch v.Field(i).Kind() {
		case reflect.Ptr:
			if v.Field(i).CanSet() && !v1.Field(i).IsNil() {
				mergeValue(v.Field(i).Elem(), v1.Field(i).Elem())
			} else {
				fmt.Fprintln(os.Stderr,"can not set or value is empty")
			}
		case reflect.Bool:
			if v.Field(i).CanSet() {
				v.Field(i).Set(v1.Field(i))
			} else {
				fmt.Fprintln(os.Stderr,"can not set or value is empty")
			}
		case reflect.Int:
			if v.Field(i).CanSet() && v1.Field(i).Int() != 0 {
				v.Field(i).Set(v1.Field(i))
			} else {
				fmt.Fprintln(os.Stderr,"can not set or value is empty")
			}
		default:
			if v.Field(i).CanSet() && v1.Field(i).Len() != 0 {
				v.Field(i).Set(v1.Field(i))
			} else {
				fmt.Fprintln(os.Stderr,"can not set or value is empty")
			}
		}
	}
}
