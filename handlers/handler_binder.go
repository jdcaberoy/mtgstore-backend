package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/jdcaberoy/mtgstore-backend/binder"
	"github.com/jdcaberoy/mtgstore-backend/user"

	model "github.com/jdcaberoy/mtgstore-backend/models"
)

type BinderServiceHandler struct {
	BinderService *binder.BinderService
	UserService   *user.UserService
}

func NewBinderServiceHandler(svc *binder.BinderService) *BinderServiceHandler {
	return &BinderServiceHandler{BinderService: svc}
}

func (s *BinderServiceHandler) CreateBinder(ctx *gin.Context) {
	var binder model.Binder
	if err := ctx.ShouldBindJSON(&binder); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"error": "could not marshall binder values: " + err.Error(),
		})
		return
	}
	err := s.BinderService.CreateBinder(ctx, binder)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			// ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "Could not create new binder." + err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"binder": binder,
	})
}

func (s *BinderServiceHandler) CreateDefaultBinder(ctx *gin.Context) {
	var binder model.Binder
	user_id := ctx.GetHeader("id")
	user, err := s.UserService.SearchbyUserID(ctx, uuid.FromStringOrNil(user_id))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	if err := ctx.ShouldBindJSON(&binder); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"error": "could not marshall binder values: " + err.Error(),
		})
		return
	}
	err = s.BinderService.CreateDefaultBinder(ctx, user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			// ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "Could not create new binder." + err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"binder": binder,
	})
}
