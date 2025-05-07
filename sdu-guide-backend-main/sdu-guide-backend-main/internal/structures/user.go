package structures

import "time"

type UserRegister struct {
	Login    string `json:"login"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
type User struct {
	ID               uint64    `bson:"_id,omitempty" json:"_id,omitempty"`
	Username         string    `bson:"username" json:"username"`
	Email            string    `bson:"email" json:"email"`
	PasswordHash     string    `bson:"passwordHash" json:"-"`
	FirstName        string    `bson:"firstName" json:"firstName"`
	LastName         string    `bson:"lastName" json:"lastName"`
	RegistrationDate time.Time `bson:"registrationDate" json:"registrationDate"`
	LastLogin        time.Time `bson:"lastLogin" json:"lastLogin"`
	ImageHash        string    `bson:"imageHash" json:"imageHash"`
}

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Session struct {
	ID          int
	UserID      int
	Token       string
	ExpiredDate time.Time
}

type UpdateUser struct {
	Username  *string `bson:"username" json:"username"`
	Email     *string `bson:"email" json:"email"`
	FirstName *string `bson:"firstName" json:"firstName"`
	LastName  *string `bson:"lastName" json:"lastName"`
	Age       *int    `bson:"age" json:"age"`
	Gender    *string `bson:"gender" json:"gender"`
	ImageHash *string `bson:"imageHash" json:"imageHash"`
}
