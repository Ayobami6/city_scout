package authservice

type User struct {
	Username string `bson:"username" json:"username"`
	Password string `bson:"password" json:"password"`
	Email    string `bson:"email" json:"email"`
	ApiKey   string `bson:"apiKey" json:"apiKey"`
}

func NewUser(username string, password string, email string, apiKey string) *User {
	return &User{
		Username: username,
		Password: password,
		Email:    email,
		ApiKey:   apiKey,
	}
}
