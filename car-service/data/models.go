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
	ID        int           `json:"id"`
	UserId    int           `json:"user_id"`
	UserName  string        `json:"user_name"`
	CarType   string        `json:"car_type"`
	CarId     sql.NullInt64 `json:"car_id"`
	City      string        `json:"city"`
	Address   string        `json:"address"`
	Active    bool          `json:"active"`
	Rating    int           `json:"rating"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

type Car struct {
	ID        int       `json:"id"`
	UserId    int       `json:"user_id"`
	CarName   string    `json:"car_name"`
	City      string    `json:"city"`
	CarType   string    `json:"car_type"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// GetAllCarRequestByCity returns active car requests by city and car type
func (c *CarRequest) GetAllCarRequestByCity(city, carType string, active bool) ([]*CarRequest, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var rows *sql.Rows
	var err error

	if len(city) > 0 && len(carType) > 0 {
		query := `
			SELECT id, user_id, user_name, car_type, car_id, city, address, active, rating, created_at, updated_at
			FROM car_requests
			WHERE city = $1 AND car_type = $2 AND active = $3
    	`
		rows, err = db.QueryContext(ctx, query, city, carType, active)
	} else if len(city) > 0 && len(carType) == 0 {
		query := `
			SELECT id, user_id, user_name, car_type, car_id, city, address, active, rating, created_at, updated_at
			FROM car_requests
			WHERE city = $1 AND active = $2
    	`
		rows, err = db.QueryContext(ctx, query, city, active)
	} else if len(city) == 0 && len(carType) > 0 {
		query := `
			SELECT id, user_id, user_name, car_type, car_id, city, address, active, rating, created_at, updated_at
			FROM car_requests
			WHERE car_type = $1 AND active = $2
    	`
		rows, err = db.QueryContext(ctx, query, carType, active)
	} else {
		query := `
			SELECT id, user_id, user_name, car_type, car_id, city, address, active, rating, created_at, updated_at
			FROM car_requests
			WHERE active = $1
    	`
		rows, err = db.QueryContext(ctx, query, active)
	}

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

// GetAllCarRequestByCity returns active car requests by user_id
func (c *CarRequest) GetCarRequestByUser(userId int) ([]*CarRequest, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var rows *sql.Rows
	var err error

	query := `
		SELECT id, user_id, user_name, car_type, car_id, city, address, active, rating, created_at, updated_at
		FROM car_requests
		WHERE user_id = $1 AND active = true
	`
	rows, err = db.QueryContext(ctx, query, userId)

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
	stmt := `insert into cars (user_id, city, car_name, car_type, created_at, updated_at)
		values ($1, $2, $3, $4, $5, $6) returning id`

	err := db.QueryRowContext(ctx, stmt,
		car.UserId,
		car.City,
		car.CarName,
		car.CarType,
		time.Now(),
		time.Now(),
	).Scan(&newID)

	if err != nil {
		return 0, err
	}

	return newID, nil
}

func (cr *CarRequest) InsertCarRequest(carRequest CarRequest) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var newID int
	stmt := `INSERT INTO car_requests (user_id, user_name, car_type, car_id, city, address, active, rating, created_at, updated_at)
             VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`

	err := db.QueryRowContext(ctx, stmt,
		carRequest.UserId,
		carRequest.UserName,
		carRequest.CarType,
		nil,
		carRequest.City,
		carRequest.Address,
		true,
		0,
		time.Now(),
		time.Now(),
	).Scan(&newID)

	if err != nil {
		return 0, err
	}

	return newID, nil
}

// GetAllCars returns cars by user ID
func (c *Car) GetAllCars(userId int, active bool) ([]*Car, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		SELECT id, user_id, car_name, city, car_type, created_at, updated_at
		FROM cars
		WHERE user_id = $1 AND active = $2
	`

	rows, err := db.QueryContext(ctx, query, userId, active)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cars []*Car

	for rows.Next() {
		var car Car
		err := rows.Scan(
			&car.ID,
			&car.UserId,
			&car.CarName,
			&car.City,
			&car.CarType,
			&car.CreatedAt,
			&car.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		cars = append(cars, &car)
	}

	return cars, nil
}

// GetCarByID retrieves a car by its ID.
func (c *Car) GetCarByID(id int) (*Car, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
        SELECT id, user_id, car_name, city, car_type, created_at, updated_at
        FROM cars
        WHERE id = $1
    `

	var car Car
	err := db.QueryRowContext(ctx, query, id).Scan(
		&car.ID,
		&car.UserId,
		&car.CarName,
		&car.City,
		&car.CarType,
		&car.CreatedAt,
		&car.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &car, nil
}

// Update updates a car's information in the database.
func (c *Car) Update() error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `
        UPDATE cars
        SET
            user_id = $1,
            car_name = $2,
            city = $3,
            car_type = $4,
            active = $5,
            updated_at = $6
        WHERE id = $7
    `

	_, err := db.ExecContext(ctx, stmt,
		c.UserId,
		c.CarName,
		c.City,
		c.CarType,
		c.Active,
		time.Now(),
		c.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

// GetCarRequestByID retrieves a car request by its ID.
func (cr *CarRequest) GetCarRequestByID(id int) (*CarRequest, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
        SELECT id, user_id, user_name, car_type, car_id, city, address, active, rating, created_at, updated_at
        FROM car_requests
        WHERE id = $1
    `

	var carRequest CarRequest
	err := db.QueryRowContext(ctx, query, id).Scan(
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

	return &carRequest, nil
}

// Update updates a car request's information in the database.
func (cr *CarRequest) Update() error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `
        UPDATE car_requests
        SET
            user_id = $1,
            user_name = $2,
            car_type = $3,
            car_id = $4,
            city = $5,
            address = $6,
            active = $7,
            rating = $8,
            updated_at = $9
        WHERE id = $10
    `

	_, err := db.ExecContext(ctx, stmt,
		cr.UserId,
		cr.UserName,
		cr.CarType,
		cr.CarId,
		cr.City,
		cr.Address,
		cr.Active,
		cr.Rating,
		time.Now(),
		cr.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

// DeleteCar deletes a car from the database based on its ID.
func (c *Car) DeleteCar() error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
        DELETE FROM cars
        WHERE id = $1
    `

	_, err := db.ExecContext(ctx, query, c.ID)
	if err != nil {
		return err
	}

	return nil
}
