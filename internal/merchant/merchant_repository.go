package merchant

import (
	"database/sql"
	"errors"
	"fmt"
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
	FindAllMerchants(params GetMerchantQueryParams) ([]Merchant, *localError.GlobalError)
	FindNearbyMerchants(location Location, params GetMerchantQueryParams) ([]MerchantWithItemQueryResult, *localError.GlobalError)
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
		return merchants, nil//localError.ErrInternalServer("Failed to find merchants", err)
	}

	return merchants, nil
}

func (r *merchantRepository) FindNearbyMerchants(location Location, params GetMerchantQueryParams) ([]MerchantWithItemQueryResult, *localError.GlobalError) {
	merchants := []MerchantWithItemQueryResult{}

	subquery := `
	select
	m.id as merchant_id,
	m.name as merchant_name,
	merchant_category,
	m.image_url as merchant_image_url,
	location_lat,
	location_long,
	m.created_at as merchant_created_at,
	i.id as item_id,
	i.name as item_name,
	product_category,
	price,
	i.image_url as item_image_url,
	i.created_at as item_created_at 
	from merchants m inner join items i on m.id = i.merchant_id` 
	nwhere := 0

	if params.MerchantID != "" {
		nwhere += 1
		subquery += fmt.Sprintf(" WHERE m.id = '%s'", params.MerchantID)
	}

	if params.Name != "" {
		prefix := "WHERE"
		if nwhere > 0 {
			prefix = "AND"
		}
		nwhere += 1
		subquery += fmt.Sprintf(" %s (m.name ILIKE '%%%s%%' OR i.name ILIKE '%%%s%%')", prefix, params.Name, params.Name)
	}

	if params.MerchantCategory == "SmallRestaurant" || params.MerchantCategory == "MediumRestaurant" || params.MerchantCategory == "LargeRestaurant" || params.MerchantCategory == "MerchandiseRestaurant" || params.MerchantCategory == "BoothKiosk" || params.MerchantCategory == "ConvenienceStore" {
		prefix := "WHERE"
		if nwhere > 0 {
			prefix = "AND"
		}
		nwhere += 1
		subquery += fmt.Sprintf(" %s merchant_category = '%s'", prefix, params.MerchantCategory)
	}

	query := fmt.Sprintf(`
	WITH myconstant (rad) as (
		values (pi()/180)
	 ),
	 mydata as (%s)
	 select
	 2 * 6371 * asin( |/( sin((%f*rad-location_lat*rad)/2::decimal)^2 + (sin((%f*rad-location_long*rad)/2::decimal)^2) * cos(location_lat*rad) * cos(%f*rad) ) ) as distance,
	 mydata.* from mydata, myconstant
	 order by distance asc;`, subquery, location.Lat, location.Long, location.Lat)

	// log.Println(query)

	err := r.db.Select(&merchants, query)
	if err != nil {
		log.Println(err)
		return merchants, nil
	}

	return merchants, nil
}