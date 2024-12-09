package controller

import (
	"net/http"
	"time"

	"github.com/fajarherdian22/credit_bank/exception"
	"github.com/fajarherdian22/credit_bank/helper"
	"github.com/fajarherdian22/credit_bank/model/web"
	"github.com/fajarherdian22/credit_bank/repository"
	"github.com/fajarherdian22/credit_bank/service"
	"github.com/fajarherdian22/credit_bank/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CustomerController struct {
	customerService *service.CustomerServiceImpl
}

func NewCustomerController(customerService *service.CustomerServiceImpl) *CustomerController {
	return &CustomerController{customerService: customerService}
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

func (controller *CustomerController) GetCustomerId(c *gin.Context) {
	var req struct {
		ID string `json:"id" binding:"required,len=36"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	payload, err := controller.customerService.GetCustomer(c.Request.Context(), req.ID)
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

func (controller *CustomerController) CreateCustomersUser(c *gin.Context) {
	var req CreateCustomersRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.ErrorResponse(err))
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
