package data

// UserBucket .
var UserBucket = []byte("user")

// User 用戶
type User struct {
	Name     string
	Password string
}