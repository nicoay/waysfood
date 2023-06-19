package handlers

import (
	dto "mytask/dto/result"
	"mytask/models"
	repositories "mytask/repository"
	"net/http"
	"strconv"

	"github.com/go-playground/validator"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

type handlerOrder struct {
	OrderRepo repositories.OrderRepo
}

func HandlerOrder(OrderRepo repositories.OrderRepo) *handlerOrder {
	return &handlerOrder{OrderRepo}
}

func (h *handlerOrder) FindOrder(c echo.Context) error {
	Order, err := h.OrderRepo.FindOrder()
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Code: http.StatusOK, Data: Order})
}

func (h *handlerOrder) GetOrderByUserSeller(c echo.Context) error {
	SID, _ := strconv.Atoi(c.Param("id"))
	userLogin := c.Get("userLogin")
	BID := userLogin.(jwt.MapClaims)["id"].(float64)
	Order, err := h.OrderRepo.GetOrderByUserSeller(int(BID), SID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Code: http.StatusOK, Data: Order})
}
func (h *handlerOrder) GetOrderByUserProduct(c echo.Context) error {
	idp, _ := strconv.Atoi(c.Param("ids"))
	userLogin := c.Get("userLogin")
	BID := userLogin.(jwt.MapClaims)["id"].(float64)
	Order, err := h.OrderRepo.GetOrderByUserProduct(int(BID), idp)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Code: http.StatusOK, Data: Order})
}

func (h *handlerOrder) GetOrderbyUser(c echo.Context) error {
	userLogin := c.Get("userLogin")
	buyerId := userLogin.(jwt.MapClaims)["id"].(float64)
	Orders, err := h.OrderRepo.GetOrderbyUser(int(buyerId))

	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Code: http.StatusOK, Data: Orders})
}

func (h *handlerOrder) GetOrder(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	Order, err := h.OrderRepo.GetOrder(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Code: http.StatusOK, Data: Order})
}

func (h *handlerOrder) CreateOrder(c echo.Context) error {
	// userLogin := c.Get("userLogin")
	// userID := userLogin.(jwt.MapClaims)["id"].(float64)
	request := new(models.Order)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
	}

	validation := validator.New()
	err := validation.Struct(request)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
	}

	// data form pattern submit to pattern entity db user
	Order := models.Order{
		Qty:       request.Qty,
		BuyerID:   request.BuyerID,
		SellerID:  request.SellerID,
		ProductID: request.ProductID,
	}

	OrderBuyer, _ := h.OrderRepo.GetOrderByUserProduct(request.BuyerID, request.ProductID)

	if OrderBuyer.ID == 0 {
		data, err := h.OrderRepo.CreateOrder(Order)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()})
		}
		return c.JSON(http.StatusOK, dto.SuccessResult{Code: http.StatusOK, Data: data})
	} else {
		OrderBuyer.Qty = request.Qty
		UpdateOrder, _ := h.OrderRepo.UpdateOrder(OrderBuyer)
		return c.JSON(http.StatusOK, dto.SuccessResult{Code: http.StatusOK, Data: UpdateOrder})
	}
}

func (h *handlerOrder) DeleteAllOrder(c echo.Context) error {

	userLogin := c.Get("userLogin")
	buyerId := userLogin.(jwt.MapClaims)["id"].(float64)
	data, err := h.OrderRepo.DeleteAllOrder(int(buyerId))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Code: http.StatusOK, Data: data})
}

func (h *handlerOrder) DeleteOrder(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	Order, err := h.OrderRepo.GetOrder(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
	}

	data, err := h.OrderRepo.DeleteOrder(Order, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Code: http.StatusOK, Data: data})
}
