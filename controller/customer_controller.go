package controller

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/fajarherdian22/credit_bank/exception"
	"github.com/fajarherdian22/credit_bank/helper"
	"github.com/fajarherdian22/credit_bank/model/web"
	"github.com/fajarherdian22/credit_bank/repository"
	"github.com/fajarherdian22/credit_bank/service"
	"github.com/fajarherdian22/credit_bank/token"
	"github.com/fajarherdian22/credit_bank/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CustomerController struct {
	customerService *service.CustomerServiceImpl
	tokenMaker      token.Maker
}

func NewCustomerController(customerService *service.CustomerServiceImpl, tokenMaker token.Maker) *CustomerController {
	return &CustomerController{
		customerService: customerService,
		tokenMaker:      tokenMaker,
	}
}

type CreateCustomersRequest struct {
	Nik          string  `json:"nik" binding:"required,len=16"`
	Password     string  `json:"password" binding:"required,min=6"`
	Email        string  `json:"email" binding:"required,email"`
	FullName     string  `json:"full_name" binding:"required"`
	LegalName    string  `json:"legal_name" binding:"required"`
	TempatLahir  string  `json:"tempat_lahir" binding:"required"`
	TanggalLahir string  `json:"tanggal_lahir" binding:"required"`
	Gaji         float64 `json:"gaji" binding:"required,numeric,gt=0"`
	PhotoKtp     string  `json:"photo_ktp" binding:"required"`
	FotoSelfie   string  `json:"foto_selfie" binding:"required"`
}

type CustomerResponse struct {
	ID       string `json:"id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
}

func NewUserResponse(customer repository.Customer) CustomerResponse {
	return CustomerResponse{
		ID:       customer.ID,
		FullName: customer.FullName,
		Email:    customer.Email,
	}
}

type LoginUserResponse struct {
	SessionID             string           `json:"session_id"`
	AccessToken           string           `json:"access_token"`
	AccessTokenExpiresAt  time.Time        `json:"access_token_expires_at"`
	RefreshToken          string           `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time        `json:"refresh_token_expires_at"`
	Customer              CustomerResponse `json:"customer"`
}

func createCustomersPayload(req CreateCustomersRequest, pw string, tgl_lahir time.Time) repository.CreateCustomersParams {
	return repository.CreateCustomersParams{
		ID:             uuid.NewString(),
		Nik:            req.Nik,
		HashedPassword: pw,
		Email:          req.Email,
		FullName:       req.FullName,
		LegalName:      req.LegalName,
		TempatLahir:    req.TempatLahir,
		TanggalLahir:   tgl_lahir,
		Gaji:           req.Gaji,
		PhotoKtp:       req.PhotoKtp,
		FotoSelfie:     req.FotoSelfie,
	}
}

func (controller *CustomerController) GetCustomersInfo(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	payload, err := controller.customerService.GetCustomer(c.Request.Context(), req.Email)
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

func (controller *CustomerController) LoginCustomers(c *gin.Context) {
	type CreateLoginReq struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}

	var req CreateLoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	customer, err := controller.customerService.GetCustomer(c, req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, helper.ErrorResponse(err))
			return
		}
		c.JSON(http.StatusInternalServerError, helper.ErrorResponse(err))
		return
	}

	err = util.CheckPassword(req.Password, customer.HashedPassword)
	if err != nil {
		c.JSON(http.StatusUnauthorized, helper.ErrorResponse(err))
		return
	}

	accessToken, accessPayload, err := controller.tokenMaker.CreateToken(customer.Email, customer.ID, 15*time.Minute)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helper.ErrorResponse(err))
		return
	}

	refreshToken, refreshPayload, err := controller.tokenMaker.CreateToken(customer.Email, customer.ID, 24*time.Hour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helper.ErrorResponse(err))
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

	session, err := controller.customerService.CreateSession(c, arg)

	if err != nil {
		c.JSON(http.StatusInternalServerError, helper.ErrorResponse(err))
		return
	}

	rsp := LoginUserResponse{
		SessionID:             session.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
		Customer:              NewUserResponse(customer),
	}
	c.JSON(http.StatusOK, rsp)

}

func (controller *CustomerController) CreateCustomersUser(c *gin.Context) {
	var req CreateCustomersRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, helper.ErrorResponse(err))
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)

	if err != nil {
		c.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return
	}

	tanggalLahir, err := util.ValidateDate(req.TanggalLahir)

	if err != nil {
		c.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	arg := createCustomersPayload(req, hashedPassword, tanggalLahir)

	payload, err := controller.customerService.CreateCustomers(c.Request.Context(), arg)
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
