package entity

import "time"

type ListUsersRequest struct {
	Paging  int
	Sorting string
}

type ListUsersResponse struct {
	Users []GetUser
}

type User struct {
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

type UserSessionRequest struct {
	Acct string `json:"account"`
	Pwd  string `json:"password"`
}

type UserSessionResponse struct {
	Jwt string `json:"jtw"`
}

type DeleteUserRequest struct {
}

type DeleteUserResponse struct {
}

type UpdateUserResponse struct {
}

type UpdateFullnameRequest struct {
	Fullname string `json:"fullname"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
