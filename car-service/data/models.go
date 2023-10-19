package data

import (
	"context"
	"database/sql"
	"time"
)

const dbTimeout = time.Second * 3

var db *sql.DB

// New is the function used to create an instance of the data package. It returns the type
// Model, which embeds all the types we want to be available to our application.
func New(dbPool *sql.DB) Models {
	db = dbPool

	return Models{
		CarRequest: CarRequest{},
		Car:        Car{},
	}
}

// Models is the type for this package. Note that any model that is included as a member
// in this type is available to us throughout the application, anywhere that the
// app variable is used, provided that the model is also added in the New function.
type Models struct {
	CarRequest CarRequest
	Car        Car
}

type CarRequest struct {
	ID        int       `json:"id"`
	UserId    int       `json:"user_id"`
	UserName  string    `json:"user_name"`
	CarType   string    `json:"car_type"`
	CarId     int       `json:"car_id"`
	City      string    `json:"city"`
	Address   string    `json:"address"`
	Active    bool      `json:"active"`
	Rating    int       `json:"rating,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Car struct {
	ID        int       `json:"id"`
	UserId    int       `json:"user_id"`
	City      int       `json:"city"`
	CarType   string    `json:"car_type"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// GetAllCarRequestByCity returns active car requests by city and car type
func (c *Car) GetAllCarRequestByCity(city, carType string) ([]*CarRequest, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
        SELECT id, user_id, user_name, car_type, car_id, city, address, active, rating, created_at, updated_at
        FROM car_requests
        WHERE active = true AND city = $1 AND car_type = $2
    `

	rows, err := db.QueryContext(ctx, query, city, carType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var carRequests []*CarRequest

	for rows.Next() {
		var carRequest CarRequest
		err := rows.Scan(
			&carRequest.ID,
			&carRequest.UserId,
			&carRequest.UserName,
			&carRequest.CarType,
			&carRequest.CarId,
			&carRequest.City,
			&carRequest.Address,
			&carRequest.Active,
			&carRequest.Rating,
			&carRequest.CreatedAt,
			&carRequest.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		carRequests = append(carRequests, &carRequest)
	}

	return carRequests, nil
}

func (u *Car) InsertCar(car Car) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var newID int
	stmt := `insert into users (user_id, city, car_type, created_at, updated_at)
		values ($1, $2, $3, $4, $5, $6, $7) returning id`

	err := db.QueryRowContext(ctx, stmt,
		car.UserId,
		car.City,
		car.CarType,
		time.Now(),
		time.Now(),
	).Scan(&newID)

	if err != nil {
		return 0, err
	}

	return newID, nil
}
