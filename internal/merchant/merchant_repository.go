package merchant

import (
	localError "belimang/pkg/error"
	"database/sql"
	"errors"
	"fmt"
	"log"

	// "strings"

	// "github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type IMerchantRepository interface {
	FindAllMerchants(params GetMerchantQueryParams) ([]Merchant, *localError.GlobalError)
	FindMerchantById(merchantId string) (*Merchant, *localError.GlobalError)
	CreateMerchant(entity Merchant) *localError.GlobalError
	FindAllItem(params GetItemQueryParam, merchantId string) ([]Item, *localError.GlobalError)
	CreateItem(entity Item) *localError.GlobalError
	CheckMerchantIDs(IDs []string) ([]Merchant, *localError.GlobalError)
	CheckItemIDs(IDs []string) ([]Item, *localError.GlobalError)
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

// List all item from database
func (r *merchantRepository) FindAllItem(param GetItemQueryParam, merchantId string) ([]Item, *localError.GlobalError) {
	// Define emtpy maps of item
	items := []Item{}

	query := "SELECT * FROM items where 1=1"

	// Filter by merhat ID
	query += fmt.Sprintf(" AND merchant_id = '%s'", merchantId)

	// Filter by ID
	if param.ItemID != "" {
		query += fmt.Sprintf(" AND id = '%s'", param.ItemID)
	}

	// Filter by Name
	if param.Name != "" {
		query += fmt.Sprintf(" name ILIKE '%%%s%%'", param.Name)
	}

	// Filter by Category
	validCategories := []ProductCategories{
		Beverage,
		Food,
		Snack,
		Condiments,
		Additions,
	}

	categoryExists := false

	if string(param.ProductCategory) != "" {
		for _, validCategory := range validCategories {
			if string(validCategory) == string(param.ProductCategory) {
				categoryExists = true
			}
		}
	}

	// Filter if category is valid
	if categoryExists {
		query += fmt.Sprintf(" AND product_category = '%s'", string(param.ProductCategory))
	}

	// Sort by created at
	if param.CreatedAt != "" {
		query += fmt.Sprintf(" order by created_at %s", string(param.CreatedAt))
	}

	// Set limit & offset
	if param.Limit != 0 {
		query += fmt.Sprintf(" LIMIT %d", param.Limit)
	} else {
		query += " LIMIT 5"
	}
	if param.Offset != 0 {
		query += fmt.Sprintf(" OFFSET %d", param.Offset)
	} else {
		query += " OFFSET 0"
	}

	err := r.db.Select(&items, query)
	if err != nil {
		return items, localError.ErrInternalServer(err.Error(), err)
	}

	return items, nil
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

func (r *merchantRepository) FindAllMerchants(params GetMerchantQueryParams) ([]Merchant, *localError.GlobalError) {
	merchants := []Merchant{}

	query := "SELECT * FROM merchants"
	nwhere := 0

	if params.MerchantID != "" {
		nwhere += 1
		query += fmt.Sprintf(" WHERE id = '%s'", params.MerchantID)
	}

	if params.Name != "" {
		prefix := "WHERE"
		if nwhere > 0 {
			prefix = "AND"
		}
		nwhere += 1
		query += fmt.Sprintf(" %s name ILIKE '%%%s%%'", prefix, params.Name)
	}

	if params.MerchantCategory == "SmallRestaurant" || params.MerchantCategory == "MediumRestaurant" || params.MerchantCategory == "LargeRestaurant" || params.MerchantCategory == "MerchandiseRestaurant" || params.MerchantCategory == "BoothKiosk" || params.MerchantCategory == "ConvenienceStore" {
		prefix := "WHERE"
		if nwhere > 0 {
			prefix = "AND"
		}
		nwhere += 1
		query += fmt.Sprintf(" %s merchant_category = '%s'", prefix, params.MerchantCategory)
	}

	if params.CreatedAt == "asc" || params.CreatedAt == "desc" {
		query += fmt.Sprintf(" ORDER BY created_at %s", params.CreatedAt)
	}

	if params.Limit != 0 {
		query += fmt.Sprintf(" LIMIT %d", params.Limit)
	} else {
		query += " LIMIT 5"
	}
	if params.Offset != 0 {
		query += fmt.Sprintf(" OFFSET %d", params.Offset)
	} else {
		query += " OFFSET 0"
	}

	// log.Println(query)

	err := r.db.Select(&merchants, query)
	if err != nil {
		log.Println(err)
		return merchants, nil //localError.ErrInternalServer("Failed to find merchants", err)
	}

	return merchants, nil
}

// Count and check if any merchant id provided in parameter is exists.
// It will return error if the merchant ID amount given in parameter is not same with the database count
func (u *merchantRepository) CheckMerchantIDs(IDs []string) ([]Merchant, *localError.GlobalError) {
	// Construct the query
	q := `
		SELECT *
		FROM merchants
		WHERE id in (?)
	`

	// Fill the quety with the place holder
	// q = fmt.Sprintf(q, generatePlaceholder(len(IDs)))

	// Create variable to store data
	var merchants []Merchant

	// Generate Query
	query, args, err := sqlx.In(q, IDs)
	if err != nil {
		return nil, localError.ErrInternalServer(err.Error(), err)
	}

	query = u.db.Rebind(query)

	err = u.db.Select(&merchants, query, args...)
	if err != nil {
		return nil, localError.ErrInternalServer(err.Error(), err)
	}

	return merchants, nil
}

func (u *merchantRepository) CheckItemIDs(IDs []string) ([]Item, *localError.GlobalError) {
	// Construct the query
	q := `
		SELECT *
		FROM items
		WHERE id in (?)
	`

	// Fill the quety with the place holder
	// q = fmt.Sprintf(q, generatePlaceholder(len(IDs)))

	// Create variable to store data
	var items []Item

	// Generate Query
	query, args, err := sqlx.In(q, IDs)
	if err != nil {
		return nil, localError.ErrInternalServer(err.Error(), err)
	}

	query = u.db.Rebind(query)

	err = u.db.Select(&items, query, args...)
	if err != nil {
		return nil, localError.ErrInternalServer(err.Error(), err)
	}

	return items, nil
}

// func generatePlaceholder(length int) string {
// 	var (
// 		placeholderTmpl []string
// 		placeholder     string
// 	)

// 	for i := 0; i < length; i++ {
// 		placeholderTmpl = append(placeholderTmpl, fmt.Sprintf("$%d", i+1))
// 	}

// 	placeholder = strings.Join(placeholderTmpl, ",")

// 	return placeholder
// }
