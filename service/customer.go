package service

import (
	"context"
	"errors"
	"rent-car/api/models"
	"rent-car/pkg/logger"
	"rent-car/storage"
)


type customerService struct {
	storage storage.IStorage
	logger logger.ILogger
	redis  storage.IRedisStorage
}

func NewCustomerService(storage storage.IStorage,logger logger.ILogger,redis storage.IRedisStorage) customerService {
	return customerService{
		storage: storage,
		logger: logger,
		redis: redis,
	}
}

func (cs customerService) Create(ctx context.Context, customer models.Customer) (string,error) {
	pkey,err := cs.storage.Customer().Create(ctx,customer)
	if err != nil {
		cs.logger.Error("ERROR in service layer while creating customer", logger.Error(err))
		return "", err
	}
	return pkey,nil
}

func (cs customerService) Update(ctx context.Context, customer models.Customer) (string,error) {
	pkey, err := cs.storage.Customer().UpdateCustomer(ctx,customer)
	if err != nil {
		cs.logger.Error("ERROR in service layer while updating customer", logger.Error(err))
		return "",err
	}

	err = cs.redis.Del(ctx,customer.Id)
	if err != nil {
		cs.logger.Error("error while setting otpCode to redis customer update", logger.Error(err))
		return "error redis update",err
	}

	return pkey,nil
}

func (cs customerService) GetByIDCustomer(ctx context.Context, id string) (models.Customer,error) {
	customer, err := cs.storage.Customer().GetByID(ctx,id)
	if err != nil {
		cs.logger.Error("ERROR in service layer while getting by id customer", logger.Error(err))
		return models.Customer{},err
	}
	return customer,nil
}

func (cs customerService) Delete(ctx context.Context, id string) (error) {
	err := cs.storage.Customer().Delete(ctx,id)
	if err != nil {
		cs.logger.Error("ERROR in service layer while deleting customer", logger.Error(err))
		return err
	}

	err = cs.redis.Del(ctx, id)
	if err != nil {
		cs.logger.Error("error while setting otpCode to redis customer deleted", logger.Error(err))
		return err
	}
	
	return nil
}

func (cs customerService) GetCustomerAll(ctx context.Context,customer models.GetAllCustomersRequest) (models.GetAllCustomersResponse, error) {
	customers, err := cs.storage.Customer().GetAllCustomer(ctx,customer)
	if err != nil {
		cs.logger.Error("ERROR in service layer while getting all customers", logger.Error(err))
		return customers,err
	}
	return customers,nil
}

func (cs customerService) UpdatePassword(ctx context.Context, customer models.PasswordOfCustomer) (string, error) {

	pKey, err := cs.storage.Customer().UpdateCustomerPassword(ctx, customer)
	if err != nil {
		cs.logger.Error("ERROR in service layer while updating Customer", logger.Error(err))
		return "", err
	}

	return pKey, nil
}

func (u customerService) GetPasswordforLogin(ctx context.Context, phone string) (string, error) {
	pKey, err := u.storage.Customer().GetPasswordforLogin(ctx, phone)
	if err != nil {
		u.logger.Error("ERROR in service layer while getbyID Customer",logger.Error(err))
		return "Error", err
	}

	return pKey, nil
}


func (u customerService) CustomerRegisterCreate(ctx context.Context, customer models.LoginCustomer) (string, error) {
	OTPCODE := u.storage.Redis().Get(ctx, customer.Gmail)
    OTPCODEStr, ok := OTPCODE.(string)
    if !ok {
        u.logger.Error("error in service layer while creating customer", logger.Error(errors.New("failed to convert OTP code to string")))
        return "the code did not match", errors.New("failed to convert OTP code to string")
    }

	if OTPCODEStr != customer.GmailCode {
        u.logger.Error("error in service layer while creating customer", logger.Error(errors.New("the code you entered is not the same as the code sent to your gmail address")))
        return "the code did not match", errors.New("the code you entered is not the same as the code sent to your gmail address")
    }

	pKey, err := u.storage.Customer().CustomerRegisterCreate(ctx, customer)
    if err != nil {
        u.logger.Error("error in service layer while creating customer", logger.Error(err))
        return "", err
    }
	return pKey,nil
}
