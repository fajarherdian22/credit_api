package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/fajarherdian22/credit_bank/exception"
	"github.com/fajarherdian22/credit_bank/helper"
	"github.com/fajarherdian22/credit_bank/repository"
	"github.com/fajarherdian22/credit_bank/service"
	"github.com/fajarherdian22/credit_bank/token"
	"github.com/fajarherdian22/credit_bank/util"
	"github.com/fajarherdian22/credit_bank/web"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type CustomerController struct {
	CustomerService *service.CustomerServiceImpl
	tokenMaker      token.Maker
	Validate        *validator.Validate
}

func NewCustomerController(CustomerService *service.CustomerServiceImpl, tokenMaker token.Maker, validate *validator.Validate) *CustomerController {
	return &CustomerController{
		CustomerService: CustomerService,
		tokenMaker:      tokenMaker,
		Validate:        validate,
	}
}

func (controller *CustomerController) LoginCustomers(c *gin.Context) {
	type CreateLoginReq struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}

	var req CreateLoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		exception.ErrorHandler(c, err)
		return
	}

	customer, err := controller.CustomerService.GetCustomer(c, req.Email)
	if err != nil {
		exception.ErrorHandler(c, err)
		return
	}

	err = util.CheckPassword(req.Password, customer.HashedPassword)
	if err != nil {
		exception.ErrorHandler(c, exception.NewNotAuthError("invalid password"))
		return
	}

	accessToken, accessPayload, err := controller.tokenMaker.CreateToken(customer.Email, customer.ID, 15*time.Minute)
	if err != nil {
		exception.ErrorHandler(c, exception.NewInternalError("failed to create access token"))
		return
	}

	refreshToken, refreshPayload, err := controller.tokenMaker.CreateToken(customer.Email, customer.ID, 24*time.Hour)
	if err != nil {
		exception.ErrorHandler(c, exception.NewInternalError("failed to refresh access token"))
		return
	}

	arg := repository.CreateSessionParams{
		ID:           refreshPayload.ID,
		Email:        customer.Email,
		CustomerID:   customer.ID,
		RefreshToken: refreshToken,
		UserAgent:    c.Request.UserAgent(),
		ClientIp:     c.ClientIP(),
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
	}

	session, err := controller.CustomerService.CreateSession(c, arg)

	if err != nil {
		exception.ErrorHandler(c, exception.NewInternalError("failed to create session"))
		return
	}

	rsp := web.LoginUserResponse{
		SessionID:             session.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
		Customer:              web.NewCustomerResponse(customer),
	}

	WebResponse := web.WebResponse{
		Code:   200,
		Data:   rsp,
		Status: "OK",
	}

	helper.HandleEncodeWriteJson(c, WebResponse)
}

func (controller *CustomerController) CreateCustomersUser(c *gin.Context) {

	var req web.CreateCustomersRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		exception.ErrorHandler(c, err)
		return
	}

	if err := controller.Validate.Struct(req); err != nil {
		exception.ErrorHandler(c, err)
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)

	if err != nil {
		exception.ErrorHandler(c, exception.NewInternalError("failed to hash password"))
		return
	}

	tanggalLahir, err := util.ValidateDate(req.TanggalLahir)
	if err != nil {
		exception.ErrorHandler(c, exception.NewBadRequestError(err.Error()))
		return
	}

	arg := web.CreateCustomersPayload(req, hashedPassword, tanggalLahir)

	payload, err := controller.CustomerService.CreateCustomers(c.Request.Context(), arg)
	if err != nil {
		exception.ErrorHandler(c, err)
		return
	}

	WebResponse := web.WebResponse{
		Code:   200,
		Data:   payload,
		Status: "OK",
	}

	helper.HandleEncodeWriteJson(c, WebResponse)
}

func (controller *CustomerController) RenewAccessToken(c *gin.Context) {
	type RenewAccessTokenRequest struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	type RenewAccessTokenResponse struct {
		AccessToken          string    `json:"access_token"`
		AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
	}

	var req RenewAccessTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		exception.ErrorHandler(c, exception.NewBadRequestError(err.Error()))
		return
	}

	refreshPayload, err := controller.tokenMaker.VerifiyToken(req.RefreshToken)
	if err != nil {
		exception.ErrorHandler(c, exception.NewNotAuthError(err.Error()))
		return
	}
	session, err := controller.CustomerService.GetSession(c, refreshPayload.ID)
	if err != nil {
		exception.ErrorHandler(c, err)
		return
	}

	if session.IsBlocked {
		err := fmt.Errorf("Blocked")
		exception.ErrorHandler(c, exception.NewNotAuthError(err.Error()))
		return
	}

	if session.Email != refreshPayload.Email {
		err := fmt.Errorf("incorrect session email")
		exception.ErrorHandler(c, exception.NewNotAuthError(err.Error()))
		return
	}

	if session.RefreshToken != req.RefreshToken {
		err := fmt.Errorf("missmatch session token")
		exception.ErrorHandler(c, exception.NewNotAuthError(err.Error()))
		return
	}

	if time.Now().After(session.ExpiresAt) {
		err := fmt.Errorf("expired session")
		exception.ErrorHandler(c, exception.NewNotAuthError(err.Error()))
		return
	}

	accessToken, AccessPayload, err := controller.tokenMaker.CreateToken(session.Email, session.CustomerID, 15*time.Minute)
	if err != nil {
		exception.ErrorHandler(c, exception.NewInternalError(err.Error()))
		return
	}

	rsp := RenewAccessTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: AccessPayload.ExpiredAt,
	}
	c.JSON(http.StatusOK, rsp)

}
