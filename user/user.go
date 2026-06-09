package user

import (
	"fmt"

	"github.com/gin-gonic/gin"
	pop "github.com/gobuffalo/pop/v6"
	model "github.com/jdcaberoy/mtgstore-backend/models"
)

type UserService struct {
	DB *pop.Connection
}

func NewUserService(db *pop.Connection) *UserService {
	return &UserService{DB: db}
}

func (s *UserService) SaveUser(ctx *gin.Context, u model.User) error {
	if !u.HasEssentials() {
		return fmt.Errorf("missing username or Password")
	}
	u.UserType = model.Guest
	tx, err := s.DB.NewTransaction()
	if err != nil {
		return fmt.Errorf("error creating new transaction: %v", err)
	}
	defer tx.Store.Rollback()
	if err = tx.Create(&u); err != nil {
		return fmt.Errorf("error creating new user: %v", err)
	}
	if err = tx.Store.Commit(); err != nil {
		return fmt.Errorf("error committ saving user")
	}
	return nil
}

func (s *UserService) UpdateUser(ctx *gin.Context, u model.User) error {
	if u.ID.IsNil() {
		return fmt.Errorf("no user ID found")
	}
	tx, err := s.DB.NewTransaction()
	if err != nil {
		return fmt.Errorf("error creating new transaction: %v", err)
	}
	defer tx.Store.Rollback()
	var foundUser model.User
	if err = tx.Q().Where("id = ?", u.ID).First(&foundUser); err != nil {
		return fmt.Errorf("error finding user: %v", err)
	}
	foundUser.Address = u.Address
	foundUser.Username = u.Username
	tx.Update(&foundUser)
	tx.Store.Commit()
	return nil
}
