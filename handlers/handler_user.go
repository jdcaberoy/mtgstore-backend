package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jdcaberoy/mtgstore-backend/binder"
	"github.com/jdcaberoy/mtgstore-backend/user"

	model "github.com/jdcaberoy/mtgstore-backend/models"
)

type UserServiceHandler struct {
	UserService   *user.UserService
	BinderService *binder.BinderService
}

func NewUserServiceHandler(usersvc *user.UserService, bindersvc *binder.BinderService) *UserServiceHandler {
	return &UserServiceHandler{UserService: usersvc, BinderService: bindersvc}
}

func (s *BinderServiceHandler) CreateUser(ctx *gin.Context) {
	var user model.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"error": "could not marshall user values: " + err.Error(),
		})
		return
	}
	user, err := s.UserService.CreateUser(ctx, user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Could not create new user." + err.Error(),
		})
		return
	}
	if err = s.BinderService.CreateDefaultBinder(ctx, user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Could not create default binder." + err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}
