package cards

import (
	"fmt"

	"github.com/gin-gonic/gin"
	pop "github.com/gobuffalo/pop/v6"
	model "github.com/jdcaberoy/mtgstore-backend/models"
)

type CardsService struct {
	DB *pop.Connection
}

func NewCardsService(db *pop.Connection) *CardsService {
	return &CardsService{DB: db}
}

type queryparams struct {
	name string
}

func (s *CardsService) FindCards(ctx *gin.Context, q queryparams) (model.Cards, error) {
	tx, err := s.DB.NewTransaction()
	if err != nil {
		return model.Cards{}, fmt.Errorf("error creating new transaction")
	}
	var foundCards model.Cards
	// TODO: allow admin to search hidden cards
	// query = tx.Q()
	// if !isRequestorAdmin(ctx){
	//
	// }
	if err := tx.Q().Where("name = ?", q.name).Where("hidden = ?", false).All(&foundCards); err != nil {
		return model.Cards{}, fmt.Errorf("error query cards: %v", err)
	}
	return foundCards, nil
}
