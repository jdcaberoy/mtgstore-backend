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

	BinderServiceInterface
}

type BinderServiceInterface interface {
	CreateBinder(ctx *gin.Context)
	AddCardstoBinder(ctx *gin.Context, binder model.Binder, cards model.Cards)
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
		return fmt.Errorf("error committing changes: %v", err)
	}
	return nil
}

type BinderQueryParams struct {
	OwnerUsername string
	BinderName    string
	Pages         int
}

func (s *BinderService) GetBinders(ctx *gin.Context, params BinderQueryParams) (model.Binders, error) {
	tx, err := s.DB.NewTransaction()
	if err != nil {
		return model.Binders{}, fmt.Errorf("error creating new transaction: %v", err)
	}
	var Binders model.Binders
	// to test fuzze search results
	if err = tx.Eager("Cards").Q().Where("OwnerUsername ILIKE ?", "%"+params.OwnerUsername+"%").All(&Binders); err != nil {
		return model.Binders{}, fmt.Errorf("error binder query: %v", err)
	}
	return Binders, nil
}

func (s *BinderService) AddCardstoBinder(ctx *gin.Context, binder model.Binder, cards model.Cards) error {
	tx, err := s.DB.NewTransaction()
	if err != nil {
		return fmt.Errorf("error creating new transaction: %v", err)
	}
	var foundBinder model.Binder
	tx.Eager("Cards").Q().Where("id = ?", binder.ID).First(&foundBinder)
	if foundBinder.ID.IsNil() {
		return fmt.Errorf("no binder found")
	}
	for k := range cards {
		cards[k].Hidden = binder.Hidden
		cards[k].BinderID = foundBinder.ID
		cards[k].BinderName = foundBinder.Name
	}
	foundBinder.Cards = append(foundBinder.Cards, cards...)
	tx.Eager("Cards").Update(&foundBinder)
	if err = tx.Store.Commit(); err != nil {
		return fmt.Errorf("error committing changes: %v", err)
	}
	return nil
}

func (s *BinderService) UpdateBinder(ctx *gin.Context, b model.Binder) error {
	tx, err := s.DB.NewTransaction()
	if err != nil {
		return fmt.Errorf("error creating new transaction: %v", err)
	}
	requestorID := ctx.GetHeader("id")

	// only the binder owner, admins, and superadmins users can update binders
	if requestorID != b.OwnerID.String() {
		var foundUser model.User
		tx.Q().Where("id = ?", requestorID).First(&foundUser)
		if foundUser.Type != model.Admin && foundUser.Type != model.SuperAdmin {
			return fmt.Errorf("error user cannot update binder")
		}
	}

	var foundBinder model.Binder
	if err = tx.Eager("Cards").Q().Where("id = ?", b.ID).First(&foundBinder.ID); err != nil {
		return fmt.Errorf("error binder query: %v", err)
	}
	if foundBinder.ID.IsNil() {
		return fmt.Errorf("no binder found")
	}
	foundBinder.Name = b.Name
	foundBinder.Description = b.Description
	if foundBinder.Hidden != b.Hidden {
		foundBinder.Hidden = b.Hidden
		for k := range foundBinder.Cards {
			foundBinder.Cards[k].Hidden = b.Hidden
		}
	}
	tx.Eager("Cards").Update(&foundBinder)
	tx.Store.Commit()
	return nil
}

func (s *BinderService) MoveCardsBindertoBinder(ctx *gin.Context, from, to model.Binder, cards model.Cards) error {
	tx, err := s.DB.NewTransaction()
	if err != nil {
		return fmt.Errorf("error creating new transaction: %v", err)
	}
	defer tx.Store.Rollback()
	var source, destination model.Binder

	tx.Eager("Cards").Q().Where("id = ?", from.ID).First(&source)
	tx.Eager("Cards").Q().Where("id = ?", to.ID).First(&destination)

	var cardsToTransfer model.Cards
	for sourceK := range source.Cards {
		for moveK := range cards {
			if source.Cards[sourceK] == cards[moveK] {
				cardsToTransfer = append(cardsToTransfer, source.Cards[sourceK])
				tx.Destroy(source.Cards[sourceK])
			}
		}
	}
	destination.Cards = append(destination.Cards, cardsToTransfer...)

	tx.Eager("Cards").Update(&source)
	tx.Eager("Cards").Update(&destination)
	tx.Store.Commit()
	return nil
}

func (s *BinderService) RemoveCardsfromBinder(ctx *gin.Context, binderID uuid.UUID, cards model.Cards) error {
	tx, err := s.DB.NewTransaction()
	if err != nil {
		return fmt.Errorf("error creating new transaction: %v", err)
	}
	defer tx.Store.Rollback()

	var foundBinder model.Binder
	if err = tx.Eager("Cards").Q().Where("id = ?", binderID).First(&foundBinder); err != nil {
		return fmt.Errorf("error query: %v", err)
	}
	if foundBinder.ID.IsNil() {
		return fmt.Errorf("no binder found")
	}
	if len(foundBinder.Cards) < len(cards) {
		return fmt.Errorf("card quantity to be removed exceeds binder card quantity")
	}
	foundBinder.Cards = removeCards(foundBinder.Cards, cards)
	if err = tx.Eager("Cards").Update(&foundBinder); err != nil {
		return fmt.Errorf("error updating changes; %v", err)
	}
	tx.Store.Commit()

	return nil
}

func removeCards(inBinder, target model.Cards) model.Cards {
	targetIDs := make(map[uuid.UUID]struct{}, len(target))

	for _, card := range target {
		targetIDs[card.ID] = struct{}{}
	}

	result := make(model.Cards, 0, len(inBinder))

	for _, card := range inBinder {
		if _, exists := targetIDs[card.ID]; !exists {
			result = append(result, card)
		}
	}

	return result
}
