package postgres

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"

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

func (m *DBClient) List(paging string, sorting string) (*entity.ListUsersResponse, error) {
	var users []entity.UserTable
	var getUser []entity.GetUser

	var res = m.client

	if sorting != "" {
		fmt.Println(sorting)
		res = res.Order(sorting)
		fmt.Println(sorting)
	}

	if paging != "" {
		// max pageSize = 10
		page, _ := strconv.Atoi(paging)

		pageSize := 10
		offset := (page - 1) * pageSize

		res = res.Offset(offset).Limit(pageSize)
	}

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

func (m *DBClient) Get(fullname string) (*entity.GetUser, error) {
	record := m.getUserRecordByFullname(fullname)
	if record == nil {
		return nil, fmt.Errorf("record not found")
	}

	return &entity.GetUser{Acct: record.Acct, Fullname: record.Fullname}, nil
}

func (m *DBClient) GetUserDetail(account string) (*entity.UserTable, error) {
	record := m.getUserRecordByAcct(account)
	if record == nil {
		return nil, fmt.Errorf("record not found")
	}
	return record, nil
}

func (m *DBClient) Insert(user entity.CreateUserRequest) error {
	res := m.client.Create(&user)
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (m *DBClient) Validate(session entity.UserSessionRequest) error {
	err := m.validateAccount(session.Acct, session.Pwd)
	if err != nil {
		return err
	}

	return nil
}

func (m *DBClient) Delete(account string) error {
	record := m.getUserRecordByAcct(account)
	if record == nil {
		return fmt.Errorf("record not found")
	}

	res := m.client.Where("acct = ?", account).Delete(&record)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (m *DBClient) Update(account string, user entity.UserTable) error {
	record := m.getUserRecordByAcct(account)
	if record == nil {
		return fmt.Errorf("record not found")
	}

	if record.Fullname == user.Fullname && record.Pwd == user.Pwd {
		return fmt.Errorf("cannot update with same value")
	}

	now := time.Now()
	m.client.Model(&entity.UserTable{}).Where("acct = ?", account).Updates(map[string]interface{}{"pwd": user.Pwd, "fullname": user.Fullname, "updated_at": now})

	return nil
}

func (m *DBClient) UpdateFullname(account string, updateFullname entity.UpdateFullnameRequest) error {
	record := m.getUserRecordByAcct(account)
	if record == nil {
		return fmt.Errorf("record not found")
	}

	if record.Fullname == updateFullname.Fullname {
		return fmt.Errorf("cannot update with same fullname")
	}

	m.client.Model(&entity.UserTable{}).Where("acct = ?", account).Updates(map[string]interface{}{"fullname": updateFullname.Fullname})

	return nil
}

func (m *DBClient) getUserRecordByFullname(fullname string) *entity.UserTable {
	var user entity.UserTable

	if err := m.client.Where("fullname = ?", fullname).First(&user).Error; err != nil {
		return nil
	} else {
		return &user
	}
}

func (m *DBClient) getUserRecordByAcct(account string) *entity.UserTable {
	var user entity.UserTable

	if err := m.client.Where("acct = ?", account).First(&user).Error; err != nil {
		return nil
	} else {
		return &user
	}
}

func (m *DBClient) validateAccount(account string, passwd string) error {
	var user entity.UserTable

	if err := m.client.Where("acct = ?", account).First(&user).Error; err != nil {
		return fmt.Errorf("not found")
	}

	if user.Pwd != passwd {
		return fmt.Errorf("incorrect passwd")
	}

	return nil
}
