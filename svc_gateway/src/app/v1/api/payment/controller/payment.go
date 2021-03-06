package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sofyan48/svc_gateway/src/app/v1/api/payment/entity"
	"github.com/sofyan48/svc_gateway/src/app/v1/api/payment/service"
	"github.com/sofyan48/svc_gateway/src/app/v1/utility/rest"
)

// PaymentController ...
type PaymentController struct {
	Service service.PaymentServiceInterface
}

// PaymentControllerHandler ...
func PaymentControllerHandler() *PaymentController {
	return &PaymentController{
		Service: service.PaymentServiceHandler(),
	}
}

// PaymentControllerInterface ...
type PaymentControllerInterface interface {
	PaymentCreate(context *gin.Context)
	UpdatePaidPayment(context *gin.Context)
	GetPaymentData(context *gin.Context)
	DeletePayment(context *gin.Context)
}

// PaymentCreate ...
func (ctrl *PaymentController) PaymentCreate(context *gin.Context) {
	payload := &entity.PaymentRequest{}
	context.ShouldBind(payload)
	result, err := ctrl.Service.PaymentCreateService(payload)
	if err != nil {
		rest.ResponseMessages(context, http.StatusInternalServerError, err.Error())
		return
	}
	rest.ResponseData(context, http.StatusOK, result)
	return
}

// GetPaymentData ...
func (ctrl *PaymentController) GetPaymentData(context *gin.Context) {
	uuid := context.Param("uuid")
	result, err := ctrl.Service.PaymentGetStatus(uuid)
	if err != nil {
		rest.ResponseMessages(context, http.StatusInternalServerError, err.Error())
		return
	}
	rest.ResponseData(context, http.StatusOK, result)
	return
}

// UpdatePaidPayment ...
func (ctrl *PaymentController) UpdatePaidPayment(context *gin.Context) {
	uuid := context.Param("uuid")
	payload := &entity.PaymentPaidRequest{}
	context.ShouldBind(payload)
	result, err := ctrl.Service.PaymentUpdateOrder(uuid, payload)
	if err != nil {
		rest.ResponseMessages(context, http.StatusInternalServerError, err.Error())
		return
	}
	rest.ResponseData(context, http.StatusOK, result)
}

// DeletePayment ...
func (ctrl *PaymentController) DeletePayment(context *gin.Context) {
	rest.ResponseMessages(context, http.StatusOK, "OK")
}
