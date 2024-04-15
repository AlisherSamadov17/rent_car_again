package service

import (
	"context"
	"fmt"
	"rent-car/api/models"
	"rent-car/config"
	"rent-car/pkg"
	"rent-car/pkg/jwt"
	"rent-car/pkg/logger"
	"rent-car/pkg/logger/password"
	"rent-car/pkg/smtp"
	"rent-car/storage"
	"time"
)

type authService struct {
	storage storage.IStorage
	log     logger.ILogger
	redis   storage.IRedisStorage

}

func NewAuthService(storage storage.IStorage, log logger.ILogger,redis storage.IRedisStorage) authService {
	return authService{
		storage: storage,
		log:     log,
		redis: redis,
	}
}

func (a authService) CustomerLogin(ctx context.Context, loginRequest models.CustomerLoginRequest) (models.CustomerLoginResponse, error) {
	fmt.Println(" loginRequest.Login: ", loginRequest.Login)
	customer, err := a.storage.Customer().GetByLogin(ctx, loginRequest.Login)
	if err != nil {
		a.log.Error("error while getting customer credentials by login", logger.Error(err))
		return models.CustomerLoginResponse{}, err
	}

	if err = password.CompareHashAndPassword(customer.Password, loginRequest.Password); err != nil {
		a.log.Error("error while comparing password", logger.Error(err))
		return models.CustomerLoginResponse{}, err
	}

	m := make(map[interface{}]interface{})

	m["user_id"] = customer.Id
	m["user_role"] = config.CUSTOMER_ROLE

    accessToken,refreshToken,err :=jwt.GenJWT(m)
	if err != nil {
		a.log.Error("error while generating tokens for customer login",logger.Error(err))
	}

	return models.CustomerLoginResponse{
		AccessToken: accessToken,
		RefreshToken: refreshToken,
	},nil
}


func (a authService) CustomerRegister(ctx context.Context, loginRequest models.CustomerRegisterRequest) error {
	_,err := a.storage.Customer().GetGmail(ctx,loginRequest.Mail)
	if err != nil {
		a.log.Error("gmail already registered",logger.Error(err))
		return err
	}
	otpCode := pkg.GenerateOTP()

	msg := fmt.Sprintf("Your otp code is: %v, for registering RENT_CAR. Don't give it to anyone", otpCode)
	
	fmt.Printf("Your otp code is: %v, for registering RENT_CAR. Don't give it to anyone", otpCode)
	
	fmt.Println(loginRequest.Mail,"----------",otpCode)

	err = a.redis.SetX(ctx, loginRequest.Mail, otpCode, time.Minute*2)
	if err != nil {
		a.log.Error("error while setting otpCode to redis customer register", logger.Error(err))
		return err
	}
	
    err = smtp.SendMail(loginRequest.Mail, msg)
	if err != nil {
		a.log.Error("error while sending otp code to customer register", logger.Error(err))
		return err
	}
   return nil
}
