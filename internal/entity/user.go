package entity

import "time"

type ListUsersRequest struct {
}

type ListUsersResponse struct {
	Users []GetUser
}

type UserTable struct {
	// The account of user
	Acct string `json:"account"`

	// The fullname of user
	Fullname string `json:"fullname"`

	// The password of user
	Pwd string `json:"password"`

	// The created time of user
	Created_at time.Time `json:"created_at"`

	// The updatedtime of user
	Updated_at time.Time `json:"updated_at"`
} // @name User

type GetUserRequest struct {
}
type GetUser struct {
	Acct     string `json:"account"`
	Fullname string `json:"fullname"`
}

type GetUserDetailRequest struct {
}

type CreateUserRequest struct {
	// The account of user
	Acct string `json:"account"`

	// The fullname of user
	Fullname string `json:"fullname"`

	// The password of user
	Pwd string `json:"password"`
}

type CreateUserResponse struct {
}
