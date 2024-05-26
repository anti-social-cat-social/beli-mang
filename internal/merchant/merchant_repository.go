package merchant

import (
	"database/sql"
	"errors"
	// "fmt"
	localError "belimang/pkg/error"
	"log"
	// "strings"

	// "github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type IMerchantRepository interface {
	CreateMerchant(entity Merchant) *localError.GlobalError
	FindMerchantById(merchantId string) (*Merchant, *localError.GlobalError)
	CreateItem(entity Item) *localError.GlobalError
}

type merchantRepository struct {
	db *sqlx.DB
}

func NewMerchantRepository(db *sqlx.DB) IMerchantRepository {
	return &merchantRepository{
		db: db,
	}
}

// This can be use for authentication process
func (u *merchantRepository) FindMerchantById(merchantId string) (*Merchant, *localError.GlobalError) {
	merchant := Merchant{}

	if err := u.db.Get(&merchant, "SELECT * FROM merchants where id=$1", merchantId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, localError.ErrNotFound("Merchant data not found", err)
		}

		log.Println(err)

		return nil, &localError.GlobalError{
			Code:    400,
			Message: "Not found",
			Error:   err,
		}

	}

	return &merchant, nil
}

// Store new merchant to database
func (u *merchantRepository) CreateMerchant(entity Merchant) *localError.GlobalError {
	q := "INSERT INTO merchants (id, name, merchant_category, image_url, location_lat, location_long) values (:id, :name, :merchant_category, :image_url, :location_lat, :location_long);"

	// Insert into database
	_, err := u.db.NamedExec(q, &entity)
	if err != nil {
		return localError.ErrInternalServer(err.Error(), err)
	}

	return nil
}

// Store new item to database
func (u *merchantRepository) CreateItem(entity Item) *localError.GlobalError {
	q := "INSERT INTO items (id, merchant_id, name, product_category, price, image_url) values (:id, :merchant_id, :name, :product_category, :price, :image_url);"

	// Insert into database
	_, err := u.db.NamedExec(q, &entity)
	if err != nil {
		return localError.ErrInternalServer(err.Error(), err)
	}

	return nil
}