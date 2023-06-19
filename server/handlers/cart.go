package handlers

import (
	"fmt"
	dto "mytask/dto/result"
	"mytask/models"
	repositories "mytask/repository"
	"net/http"
	"strconv"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

type handlerCart struct {
	CartRepo repositories.CartRepo
}

func HandlerCart(CartRepo repositories.CartRepo) *handlerCart {
	return &handlerCart{CartRepo}
}

func (h *handlerCart) FindCarts(c echo.Context) error {
	carts, err := h.CartRepo.FindCarts()
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Code: http.StatusOK, Data: carts})
}

func (h *handlerCart) GetCarts(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	carts, err := h.CartRepo.GetCarts(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Code: http.StatusOK, Data: carts})
}

func (h *handlerCart) CreateCarts(c echo.Context) error {
	request := new(models.Carts)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
	}

	validation := validator.New()
	err := validation.Struct(request)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
	}

	carts := models.Carts{}
	fmt.Println(carts, "ini cart")

	data, err := h.CartRepo.CreateCarts(carts)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Code: http.StatusOK, Data: data})
}
