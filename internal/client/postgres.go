package postgres

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/jinzhu/gorm"
)

type DBClient struct {
	client *gorm.DB
}

type dbConfig struct {
	Addr     string
	Port     int
	Username string
	Name     string
	Password string
}

func getDbConfig() *dbConfig {
	config := dbConfig{}

	file := "../../assets/configs/config.json"
	data, err := ioutil.ReadFile(file)

	err = json.Unmarshal(data, &config)
	if err != nil {
		panic(err)
	}

	return &config
}

func (m *DBClient) Connect() {
	config := getDbConfig()

	connStr := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable", config.Addr, config.Port, config.Username, config.Name, config.Password)
	client, err := gorm.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	m.client = client
}
