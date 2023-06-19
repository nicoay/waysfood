package handlers

import (
	"context"
	"fmt"
	productdto "mytask/dto/product"
	dto "mytask/dto/result"
	"mytask/models"
	repositories "mytask/repository"
	"net/http"
	"os"
	"strconv"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/go-playground/validator"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

type handlerProduct struct {
	ProductRepo repositories.ProductRepo
	UserRepo    repositories.UseRepository
}

func HandlerProduct(ProductRepo repositories.ProductRepo, UserRepo repositories.UseRepository) *handlerProduct {
	return &handlerProduct{ProductRepo, UserRepo}
}

func (h *handlerProduct) FindProduct(c echo.Context) error {
	Product, err := h.ProductRepo.FindProduct()
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Code: http.StatusOK, Data: Product})
}

func (h *handlerProduct) FindProductPartner(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	// userLogin := c.Get("userLogin")
	// userID := userLogin.(jwt.MapClaims)["id"].(float64)

	Products, err := h.ProductRepo.FindProductPartner(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Code: http.StatusOK, Data: Products})
}

func (h *handlerProduct) GetProduct(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	var Product models.Product
	Product, err := h.ProductRepo.GetProduct(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Code: http.StatusOK, Data: Product})
}

func (h *handlerProduct) CreateProduct(c echo.Context) error {
	dataFile := c.Get("dataFile").(string)

	userLogin := c.Get("userLogin")
	sellerId := userLogin.(jwt.MapClaims)["id"].(float64)

	Price, _ := strconv.ParseInt(c.FormValue("price"), 10, 64)
	// Qty, _ := strconv.Atoi(c.FormValue("qty"))
	// UserId, _ := strconv.Atoi(c.FormValue("user_id"))
	request := productdto.CreateProductReq{
		Title:  c.FormValue("title"),
		Image:  dataFile,
		Price:  Price,
		UserId: int(sellerId),
	}

	validation := validator.New()
	err := validation.Struct(request)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
	}
	var ctx = context.Background()
	var CLOUD_NAME = os.Getenv("CLOUD_NAME")
	var API_KEY = os.Getenv("API_KEY")
	var API_SECRET = os.Getenv("API_SECRET")

	// Add your Cloudinary credentials ...
	cld, _ := cloudinary.NewFromParams(CLOUD_NAME, API_KEY, API_SECRET)

	// Upload file to Cloudinary ...
	resp, err := cld.Upload.Upload(ctx, dataFile, uploader.UploadParams{Folder: "WAYSFOOD"})
	// data form pattern submit to pattern entity db Product

	if err != nil {
		fmt.Println(err.Error())
	}
	Products := models.Product{
		Title:  request.Title,
		Image:  resp.SecureURL,
		Price:  request.Price,
		UserID: request.UserId,
	}

	data, err := h.ProductRepo.CreateProduct(Products)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Code: http.StatusOK, Data: data})
}
