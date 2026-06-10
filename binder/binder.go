package binder

import (
	"fmt"

	"github.com/gin-gonic/gin"
	pop "github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"
	model "github.com/jdcaberoy/mtgstore-backend/models"
)

type BinderService struct {
	DB *pop.Connection
}

type BinderServiceInterface interface {
	CreateBinder(ctx *gin.Context)
}

func NewBinderService(db *pop.Connection) *BinderService {
	return &BinderService{DB: db}
}

func (s *BinderService) CreateBinder(ctx *gin.Context, b model.Binder) error {
	tx, err := s.DB.NewTransaction()
	if err != nil {
		return fmt.Errorf("error creating new transaction: %v", err)
	}
	defer tx.Store.Rollback()

	var newBinder model.Binder
	newBinder.Name = b.Name
	requestorName := ctx.GetHeader("username")
	requestorID := ctx.GetHeader("id")
	newBinder.OwnerID = uuid.FromStringOrNil(requestorID)
	newBinder.OwnerUsername = requestorName
	newBinder.Description = b.Description
	if err = tx.Create(&newBinder); err != nil {
		return fmt.Errorf("error creating new binder: %v", err)
	}
	tx.Store.Commit()
	return nil
}

func (s *BinderService) CreateDefaultBinder(ctx *gin.Context, u model.User) error {
	tx, err := s.DB.NewTransaction()
	if err != nil {
		return fmt.Errorf("error creating new transaction: %v", err)
	}
	var newBinder model.Binder
	newBinder.Name = "_unsorted"
	newBinder.Description = "default bunder where unsorted cards go"
	tx.Create(&newBinder)
	if err = tx.Store.Commit(); err != nil {
		return fmt.Errorf("error commiting changes: %v", err)
	}
	return nil
}
