package users

import (
	"time"
	"errors"

	"golang.org/x/crypto/bcrypt"
	"github.com/jinzhu/gorm"

	"github.com/ld100/goblet/environment"
	"github.com/ld100/goblet/log"
)

// TODO: Remove this test func
func MigrateUsers() {
	environment.GDB.DropTable(&User{})
	environment.GDB.AutoMigrate(&User{})

	user := User{
		FirstName: "John the",
		LastName:  "Doe",
		Email:     "you@example.com",
		Password:  "password",
	}

	var errs []error
	errs = user.CreateUser()
	if errs != nil {
		log.Fatal(errs)
	}
	log.Debug("Created user entity with ID: ", user.ID)

	user.FirstName = "Peter"
	errs = user.SaveUser()
	if errs != nil {
		log.Fatal(errs)
	}
	log.Debug("Updated user entity with ID: ", user.ID)

	userCopy := User{ID: user.ID}
	errs = userCopy.FindUserByID()
	if errs != nil {
		log.Fatal(errs)
	}
	log.Debug("Fetched another instance of user : ", userCopy.FirstName)

	userEmailCopy := User{Email: user.Email}
	errs = userEmailCopy.FindUserByEmail()
	if errs != nil {
		log.Fatal(errs)
	}
	log.Debug("Fetched another instance of user by e-mail : ", userEmailCopy.Email)

	errs = user.DeleteUser()
	if errs != nil {
		log.Fatal(errs)
	}
	log.Debug("Deleted user entity with ID: ", user.ID)

	var users []*User
	users = FindAllUsers()
	log.Debug("Users found: ", len(users))

	defer environment.GDB.Close()
}

func (u *User) CreateUser() []error {
	if environment.GDB.NewRecord(u) {
		return environment.GDB.Create(&u).GetErrors()
	}
	return nil
}

func (u *User) DeleteUser() []error {
	if !environment.GDB.NewRecord(u) {
		return environment.GDB.Delete(&u).GetErrors()
	}
	return nil
}

func (u *User) SaveUser() []error {
	return environment.GDB.Save(&u).GetErrors()
}

func (u *User) FindUserByID() []error {
	return environment.GDB.First(&u, u.ID).GetErrors()
}

func (u *User) FindUserByEmail() []error {
	return environment.GDB.Where("email = ?", u.Email).First(&u).GetErrors()
}

func FindAllUsers() []*User {
	var users []*User
	errs := environment.GDB.Find(&users).GetErrors()
	if errs != nil {
		log.Fatal(errs)
	}
	return users
}

type User struct {
	ID        uint   `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	FirstName string `gorm:"size:255" valid:"optional"`              // Default size for string is 255, reset it with this tag
	LastName  string `gorm:"size:255" valid:"optional"`              // Default size for string is 255, reset it with this tag
	Email     string `gorm:"not null;unique" valid:"required,email"` // Set field as not nullable and unique
	Password  string `gorm:"size:255"`
}

// GORM callback: Encode password before create
func (u *User) BeforeCreate() (err error) {
	u.Password, err = HashPassword(u.Password)
	if err != nil {
		err = errors.New("cannot hash user password")
		log.Fatal(err)
	}
	return
}

func (u User) Validate(db *gorm.DB) {
	if u.FirstName == "Joe" {
		db.AddError(errors.New("go change your name, it is invalid"))
	}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
