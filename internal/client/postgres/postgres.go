package postgres

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/a23667788/ui-assignment/internal/entity"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
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

	data, err := ioutil.ReadFile("configs/config.json")
	if err != nil {
		panic(err)
	}

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

func (m *DBClient) Disconnect() {
	m.client.Close()
}

func (m *DBClient) List() (*entity.ListUsersResponse, error) {
	var users []entity.UserTable
	var getUser []entity.GetUser
	var res = m.client

	res = res.Find(&users)

	if res.Error != nil {
		return nil, res.Error
	}

	for _, user := range users {
		entiity := entity.GetUser{
			Acct:     user.Acct,
			Fullname: user.Fullname,
		}
		getUser = append(getUser, entiity)
	}

	return &entity.ListUsersResponse{Users: getUser}, nil
}
