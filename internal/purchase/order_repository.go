package purchase

import (
	localError "belimang/pkg/error"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

type IOrderRepository interface {
	CreateEstimation(entity *OrderEstimation) (string, *localError.GlobalError)
	CreateOrderMerchant(orderEstimationID string, entity []OrderEstimationDetail) *localError.GlobalError
	PlaceOrder(orderEstimationID string) (string, *localError.GlobalError)
}

type orderRepository struct {
	db *sqlx.DB
}

// PlaceOrder implements IOrderRepository.
func (repo *orderRepository) PlaceOrder(orderEstimationID string) (string, *localError.GlobalError) {
	// Search for estimation ID
	searchQuery := "SELECT id from order_estimation where id = $1"

	estErr := repo.db.QueryRow(searchQuery, orderEstimationID).Err()
	if estErr != nil {
		return "", localError.ErrNotFound("Estimation data not found", estErr)
	}

	// Order ID
	var id string

	// Construct query
	q := "INSERT INTO orders (order_estimation_id) values ($1) returning id"

	err := repo.db.QueryRowx(q, orderEstimationID).Scan(&id)
	if err != nil {
		return "", localError.ErrInternalServer(err.Error(), err)
	}

	return id, nil
}

// CreateOrderMerchant implements IOrderRepository.
func (repo *orderRepository) CreateOrderMerchant(orderEstimationID string, entity []OrderEstimationDetail) *localError.GlobalError {
	log.Println(orderEstimationID)
	// Construct insert query & param
	q := "INSERT INTO order_estimation_items (order_estimation_id,item_id,quantity) VALUES "
	var insertParam []any

	// Loop to get the full data to be stored
	for i, data := range entity {
		pos := i * 3

		// Generate placeholder
		q += fmt.Sprintf("($%d,$%d,$%d),", pos+1, pos+2, pos+3)

		// Generate binding value
		insertParam = append(insertParam, orderEstimationID, data.ItemID, data.Quantity)
	}

	q = q[:len(q)-1] // Hilangkan ","

	log.Println(q, insertParam)
	_, err := repo.db.Exec(q, insertParam...)
	if err != nil {
		return localError.ErrInternalServer(err.Error(), err)
	}

	return nil
}

// CreateEstimation store the user estimated order time and price
func (repo *orderRepository) CreateEstimation(entity *OrderEstimation) (string, *localError.GlobalError) {
	var id string

	// Insert Query
	q := `INSERT INTO order_estimation 
			(user_id,user_location_lat,user_location_long,total_price,estimated_delivery_time) 
			values 
				($1,$2,$3,$4,$5)
			RETURNING id
		`

	err := repo.db.QueryRowx(
		q,
		entity.UserID,
		entity.UserLat,
		entity.UserLong,
		entity.Price,
		entity.EstimatedTime,
	).Scan(&id)

	if err != nil {
		return "", localError.ErrInternalServer(err.Error(), err)
	}

	return id, nil
}

func NewOrderRepository(db *sqlx.DB) IOrderRepository {
	return &orderRepository{
		db: db,
	}
}
