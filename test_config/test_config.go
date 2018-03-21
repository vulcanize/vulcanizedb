package test_config

import (
	"log"

	"os"

	"github.com/spf13/viper"
	"github.com/vulcanize/vulcanizedb/pkg/config"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

var TestConfig *viper.Viper
var DBConfig config.Database
var Infura *viper.Viper
var InfuraClient config.Client
var ABIFilePath string

func init() {
	setTestConfig()
	setInfuraConfig()
	setABIPath()
}

func setTestConfig() {
	TestConfig = viper.New()
	TestConfig.SetConfigName("private")
	TestConfig.AddConfigPath("$GOPATH/src/github.com/vulcanize/vulcanizedb/environments/")
	err := TestConfig.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
	hn := TestConfig.GetString("database.hostname")
	port := TestConfig.GetInt("database.port")
	name := TestConfig.GetString("database.name")
	DBConfig = config.Database{
		Hostname: hn,
		Name:     name,
		Port:     port,
	}
}

func setInfuraConfig() {
	Infura = viper.New()
	Infura.SetConfigName("infura")
	Infura.AddConfigPath("$GOPATH/src/github.com/vulcanize/vulcanizedb/environments/")
	err := Infura.ReadInConfig()
	ipc := Infura.GetString("client.ipcpath")
	if err != nil {
		log.Fatal(err)
	}
	InfuraClient = config.Client{
		IPCPath: ipc,
	}
}

func setABIPath() {
	gp := os.Getenv("GOPATH")
	ABIFilePath = gp + "/src/github.com/vulcanize/vulcanizedb/pkg/geth/testing/"
}

func NewTestDB(node core.Node) *postgres.DB {
	db, _ := postgres.NewDB(DBConfig, node)
	db.MustExec("DELETE FROM watched_contracts")
	db.MustExec("DELETE FROM transactions")
	db.MustExec("DELETE FROM blocks")
	db.MustExec("DELETE FROM logs")
	db.MustExec("DELETE FROM receipts")
	db.MustExec("DELETE FROM log_filters")
	return db
}
