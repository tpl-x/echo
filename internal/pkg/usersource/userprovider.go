package usersource

type UserProvider interface {
	// UserId  returns the user uniq id
	UserId() any
	// UserName returns  username ,normally is the login name
	UserName() string
}
