package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"readfiles/dao"
)

const (
	DBHost     = 0
)

func init() {
	readConfig()
}

// Load config file
func Load(category int, manifest interface{}) error {
	switch category {
	case DBHost:
		loadDBHost(*(manifest.(*map[string]*dao.DBHost)))
		return nil
	default:
		return fmt.Errorf("current system does't support this category of configuration")
	}
}

func readConfig() {
	// set config file
	registerConfigFile()
}

func loadDBHost(dbh map[string]*dao.DBHost) {
	dbh["group"] = &dao.DBHost{}
	dbh["group"].Host = viper.GetString("datastore.mysql.group.host")
	dbh["group"].Port = viper.GetString("datastore.mysql.group.port")
	dbh["group"].Username = viper.GetString("datastore.mysql.group.username")
	dbh["group"].Password = viper.GetString("datastore.mysql.group.password")
	dbh["group"].DBName = viper.GetString("datastore.mysql.group.dbname")
}


func registerConfigFile() {
	viper.SetConfigType("json")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Println(err)
	}
}
