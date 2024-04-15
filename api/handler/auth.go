package handler

import (
	"fmt"
	"net/http"
	_ "rent-car/api/docs"
	"rent-car/api/models"
	"rent-car/pkg/check"
	"rent-car/pkg/logger/password"

	"github.com/gin-gonic/gin"
)

// CustomerLogin godoc
// @Router       /customer/login [POST]
// @Summary      Customer login
// @Description  Customer login
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        login body models.CustomerLoginRequest true "login"
// @Success      201  {object}  models.CustomerLoginResponse
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h *Handler) CustomerLogin(c *gin.Context)  {
	loginReq := models.CustomerLoginRequest{}

	if err := c.ShouldBindJSON(&loginReq);err != nil {
		handlerResponseLog(c, h.Log, "error while binding body", http.StatusBadRequest, err)
		return
	}
	fmt.Println("loginReq:",loginReq)

	loginResp,err := h.Services.Auth().CustomerLogin(c.Request.Context(),loginReq)
	if err != nil {
		handlerResponseLog(c,h.Log,"unauthorized",http.StatusUnauthorized,err)
		return
	}
	handlerResponseLog(c,h.Log,"succes",http.StatusOK,loginResp)
}


// CustomerRegister godoc
// @Router       /customer/register [POST]
// @Summary      Customer register
// @Description  Customer register
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        register body models.CustomerRegisterRequest true "register"
// @Success      201  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h *Handler) CustomerRegister(c *gin.Context)  {
	loginReq := models.CustomerRegisterRequest{}

	if err := c.ShouldBindJSON(&loginReq); err != nil {
		handlerResponseLog(c, h.Log, "error while reading body", http.StatusBadRequest, err)
		return
	}
	
	if err := check.CheckEmail(loginReq.Mail); err != nil {
		handlerResponseLog(c,h.Log,"Email address does not exist or is undeliverable "+loginReq.Mail, http.StatusBadRequest, err.Error())
		return
	}

	err := h.Services.Auth().CustomerRegister(c.Request.Context(), loginReq)
	if err != nil {
		handlerResponseLog(c, h.Log, "", http.StatusInternalServerError,err)
		return
	}
	handlerResponseLog(c, h.Log, "Otp sent successfull", http.StatusOK, "okey")
}


// CustomerCreateRegister godoc
// @Router       /customer/auth/create [POST]
// @Summary      Customer password check and create 
// @Description  Customer password check and create
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        login body models.LoginCustomer true "login"
// @Success      201  {object}  models.LoginCustomer
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h *Handler) CustomerRegisterCreate(c *gin.Context) {
	loginReq := models.LoginCustomer{}

	if err := c.ShouldBindJSON(&loginReq); err != nil {
		handlerResponseLog(c, h.Log, "error while binding body", http.StatusBadRequest, err)
		return
	}
	fmt.Println("loginReq: ", loginReq)

	if err := check.CheckEmail(loginReq.Gmail); err != nil {
		handlerResponseLog(c,h.Log,"Email address does not exist or is undeliverable "+loginReq.Gmail, http.StatusBadRequest, err.Error())
		return
	}
	if err := check.ValidatePassword(loginReq.Password); err != nil {
		handlerResponseLog(c,h.Log,"error while validating  password, password: "+loginReq.Password, http.StatusBadRequest, err.Error())
		return
	}

	HashPassword,err:=password.HashPassword(loginReq.Password)
	if err!=nil{
		handlerResponseLog(c, h.Log, "password hashed error", http.StatusUnauthorized, err)
	}

	loginReq.Password=HashPassword

	loginResp, err := h.Services.Customer().CustomerRegisterCreate(c.Request.Context(), loginReq)
	if err != nil {
		handlerResponseLog(c, h.Log, "erorororor", http.StatusUnauthorized, err)
		return
	}

	handlerResponseLog(c, h.Log, "Succes", http.StatusOK, loginResp)

}