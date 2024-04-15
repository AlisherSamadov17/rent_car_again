package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"rent-car/api/models"
	"rent-car/config"
	"rent-car/pkg"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type carRepo struct {
	db *pgxpool.Pool	
}

func NewCar(db *pgxpool.Pool) carRepo {
	return carRepo{
		db: db,
	}
}

func (c *carRepo) Create(ctx context.Context,car models.CreateCar) (string, error) {

	id := uuid.New()

	query := ` INSERT INTO cars(
		id,
		name,
		brand,
		model,
		hourse_power,
		colour,
		engine_cap,
		year)
		VALUES($1,$2,$3,$4,$5,$6,$7,$8) 
	`

	ctx,cancel:= context.WithTimeout(ctx,config.TimewithContex)
	defer cancel()

	_,err := c.db.Exec(ctx,query,
		id.String(),
		car.Name, car.Brand,
		car.Model, car.HoursePower,
		car.Colour, car.EngineCap,car.Year)

	if err != nil {
		fmt.Println(err.Error())
	}

	return id.String(), nil
}

func (c *carRepo) Update(ctx context.Context,car models.UpdateCar) (string, error) {

	query := `UPDATE cars set
			name=$1,
			year=$2,
			brand=$3,
			model=$4,
			hourse_power=$5,
			colour=$6,
			engine_cap=$7,
			updated_at=CURRENT_TIMESTAMP WHERE id = $8 AND deleted_at = 0`
	ctx,cancel:= context.WithTimeout(ctx,config.TimewithContex)
	defer cancel()

	    hoursepower := car.HoursePower
		enginecap :=car.EngineCap

    
	_, err := c.db.Exec(ctx,query,
		car.Name,car.Year,car.Brand,
		car.Model, hoursepower,
		car.Colour, enginecap, car.Id)

	if err != nil {
		fmt.Println(err.Error())
	}

	return car.Id, nil
}

func (c *carRepo) GetAll(ctx context.Context,req models.GetAllCarsRequest) (models.GetAllCarsRealResponse, error) {
	var (
		resp   = models.GetAllCarsRealResponse{}
		filter = ""
	)
	offset := (req.Page - 1) * req.Limit

	if req.Search != "" {
		filter += fmt.Sprintf(` and name ILIKE  '%%%v%%' `, req.Search)
	}

	filter += fmt.Sprintf("OFFSET %v LIMIT %v", offset, req.Limit)
	fmt.Println("filter:", filter)

	ctx,cancel:= context.WithTimeout(ctx,config.TimewithContex)
	defer cancel()

	rows, err := c.db.Query(ctx,`select count(id) OVER(),id,name,brand,model,hourse_power,engine_cap,colour,created_at,updated_at,year FROM cars WHERE deleted_at = 0` + filter + ``)
	if err != nil {
		return resp, err
	}
	for rows.Next() {
		var (
			car      = models.CarReal{}
			updateAt     sql.NullString
			createdAt    sql.NullString
		)
        value := car.HoursePower
		value2 :=car.EngineCap
		// hoursepower := pkg.EncodeToStringFor(value)
		// enginecap := pkg.EncodeToStringFor(int(value2))

		if err := rows.Scan(
			&resp.Count,
			&car.Id,
			&car.Name,
			&car.Brand,
			&car.Model,
			&car.Year,
			&value,
			&value2,
			&car.Colour,			
			&updateAt,
			&createdAt); err != nil {
			return resp, err
		}
        
         
		car.UpdatedAt = pkg.NullStringToString(updateAt)
		car.CreatedAt = pkg.NullStringToString(createdAt)
		resp.Cars = append(resp.Cars, car)
	}
	return resp, nil
}

func (c *carRepo) GetAvaibleCars(ctx context.Context,req models.GetAllCarsRequest) (models.GetAllCarsResponse, error) {
	var (
		resp   = models.GetAllCarsResponse{}
		filter = ""
	)
	offset := (req.Page - 1) * req.Limit

	if req.Search != "" {
		filter += fmt.Sprintf(` and name ILIKE  '%%%v%%' `, req.Search)
	}

	filter += fmt.Sprintf(" OFFSET %v LIMIT %v", offset, req.Limit)
	fmt.Println("filter:", filter)

	ctx,cancel:= context.WithTimeout(ctx,config.TimewithContex)
	defer cancel()
	rows,err := c.db.Query(ctx,`
	SELECT name,model 
	FROM cars
	WHERE id NOT IN (SELECT car_id FROM orders WHERE car_id IS NOT NULL)`+ filter + ``)
	if err != nil {
		return resp,err
	}
	for rows.Next() {
	     var (
			car = models.Car{}
		 )

		if err := rows.Scan(

			&car.Name,
			&car.Model,);err != nil{
				return resp,err
			}

		resp.Cars = append(resp.Cars, car)	
	}
	if err = rows.Err();err != nil {
		return resp,err
	   }
	   countQuery := `SELECT COUNT(*) AS count_of_available_cars
	   FROM cars
	   WHERE id NOT IN (SELECT car_id FROM orders WHERE car_id IS NOT NULL)`

	   err = c.db.QueryRow(ctx,countQuery).Scan(&resp.Count)
		 if err != nil{
			return resp,err
		 }
	   return resp,nil
}

func (c *carRepo) GetByID(ctx context.Context,id string) (models.GetByIDCar, error) {
	car := models.GetByIDCar{}

	var (
		updateAt     sql.NullString
		createdAt    sql.NullString
	)
	ctx,cancel:= context.WithTimeout(ctx,config.TimewithContex)
	defer cancel()

	if err := c.db.QueryRow(ctx,`select id,name,brand,model,hourse_power,colour,engine_cap,year,created_at,updated_at from cars where id = $1`, id).Scan(
		&car.Id,
		&car.Name,
		&car.Brand,
		&car.Model,
		&car.HoursePower,
		&car.Colour,
		&car.EngineCap,
		&car.Year,
		&updateAt,
		&createdAt,
	); err != nil {
		return car, err
	}
	car.UpdatedAt = pkg.NullStringToString(updateAt)
	car.CreatedAt = pkg.NullStringToString(createdAt)
	return car, nil
}

func (c *carRepo) Delete(ctx context.Context,id string) error {

	query := `delete from cars WHERE id = $1`

	ctx,cancel:= context.WithTimeout(ctx,config.TimewithContex)
	defer cancel()

	_, err := c.db.Exec(ctx,query, id)
	if err != nil {
		return err
	}

	return nil
}