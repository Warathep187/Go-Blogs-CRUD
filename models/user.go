package models

type User struct {
	ID       string `bson:"_id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type UserSessionData struct {
	ID    string `json:"_id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}
