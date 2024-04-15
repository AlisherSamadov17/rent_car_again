package storage

import (
	"context"
	"rent-car/api/models"
	"time"
)

type IStorage interface {
	CloseDB()
	Car() ICarStorage
	Customer() ICustomerStorage
	Order() IOrderStorage
	Redis() IRedisStorage
}

type ICarStorage interface {
	Create(context.Context,models.CreateCar) (string, error)
	GetByID(ctx context.Context,id string) (models.GetByIDCar, error)
	GetAvaibleCars(ctx context.Context,req models.GetAllCarsRequest) (models.GetAllCarsResponse, error)
	GetAll(context.Context,models.GetAllCarsRequest) (models.GetAllCarsRealResponse, error)
	Update(ctx context.Context,car models.UpdateCar) (string, error)
	Delete(ctx context.Context,id string) error
}

type ICustomerStorage interface {
	Create(context.Context,models.Customer) (string, error)
	GetByID(ctx context.Context,id string) (models.Customer, error)
	GetAllCustomer(ctx context.Context,req models.GetAllCustomersRequest) (models.GetAllCustomersResponse, error)
	UpdateCustomer(context.Context,models.Customer) (string, error)
	UpdateCustomerPassword(context.Context,models.PasswordOfCustomer)(string, error)
	GetPasswordforLogin(ctx context.Context, phone string) (string, error)
	Delete(ctx context.Context,id string) error
	GetByLogin(context.Context, string) (models.GetAllCustomer, error)
	GetGmail(ctx context.Context, gmail string) (string, error)
	CustomerRegisterCreate(ctx context.Context, customer models.LoginCustomer) (string, error)
}

type IOrderStorage interface {
	Create(context.Context,models.CreateOrder) (string, error)
	GetByID(ctx context.Context,id string) (models.OrderAll, error)
	GetAll(ctx context.Context,request models.GetAllOrdersRequest) (models.GetAllOrdersResponse, error)
	Update(context.Context,models.UpdateOrder) (string, error)
	Delete(ctx context.Context,id string) error
	UpdateOrderStatus(context.Context,models.GetOrder) (string, error)
}

type IRedisStorage interface {
	SetX(ctx context.Context, key string, value interface{}, duration time.Duration) error
	Get(ctx context.Context, key string) interface{}
	Del(ctx context.Context, key string) error
}