package handlers

import (
	"context"
	"fmt"
	dto "mytask/dto/result"
	usersdto "mytask/dto/user"
	"mytask/models"
	repositories "mytask/repository"
	"os"

	"net/http"
	"strconv"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/go-playground/validator"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

type HandlerUser struct {
	UserRepository repositories.UseRepository
}

func UserHandler(UserRepository repositories.UseRepository) *HandlerUser {
	return &HandlerUser{UserRepository}
}

// handler

func (h *HandlerUser) FindUsers(c echo.Context) error {
	users, err := h.UserRepository.FindUser()
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Code: http.StatusOK, Data: users})
}

func (h *HandlerUser) GetUser(c echo.Context) error {
	// id, _ := strconv.Atoi(c.Param("id"))
	userLogin := c.Get("userLogin")
	userID := userLogin.(jwt.MapClaims)["id"].(float64)

	user, err := h.UserRepository.GetUser(int(userID))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Code: http.StatusOK, Data: user})
}

func (h *HandlerUser) FindPartner(c echo.Context) error {
	user, err := h.UserRepository.FindPartner()
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
	}
	return c.JSON(http.StatusOK, dto.SuccessResult{Code: http.StatusOK, Data: user})
}

func (h *HandlerUser) CreateUser(c echo.Context) error {
	dataFile := c.Get("dataFile").(string)
	request := new(usersdto.CreateUserRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
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

	if err != nil {
		fmt.Println(err.Error())
	}

	user := models.User{
		Fullname: request.Fullname,
		Email:    request.Email,
		Password: request.Password,
		Gender:   request.Gender,
		Phone:    request.Phone,
		Location: request.Location,
		Role:     request.Role,
		Image:    resp.SecureURL,
	}

	data, err := h.UserRepository.CreateUser(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()})
	}
	return c.JSON(http.StatusOK, dto.SuccessResult{Code: http.StatusOK, Data: data})
}

func (h *HandlerUser) UpdateUser(c echo.Context) error {
	dataFile := c.Get("dataFile").(string)
	fmt.Println("this is data file", dataFile)
	userLogin := c.Get("userLogin")
	userID := userLogin.(jwt.MapClaims)["id"].(float64)

	var ctx = context.Background()
	var CLOUD_NAME = os.Getenv("CLOUD_NAME")
	var API_KEY = os.Getenv("API_KEY")
	var API_SECRET = os.Getenv("API_SECRET")

	// Add your Cloudinary credentials ...
	cld, _ := cloudinary.NewFromParams(CLOUD_NAME, API_KEY, API_SECRET)

	// Upload file to Cloudinary ...
	resp, err := cld.Upload.Upload(ctx, dataFile, uploader.UploadParams{Folder: "WAYSFOOD"})

	if err != nil {
		fmt.Println(err.Error())
	}

	request := usersdto.UpdateUserRequest{
		Fullname: c.FormValue("fullname"),
		Email:    c.FormValue("email"),
		Location: c.FormValue("location"),
		Image:    resp.SecureURL,
		Phone:    c.FormValue("phone"),
	}

	profile, err := h.UserRepository.GetUser(int(userID))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
	}

	if request.Email != "" {
		profile.Email = request.Email
	}
	if request.Phone != "" {
		profile.Phone = request.Phone
	}
	if request.Fullname != "" {
		profile.Fullname = request.Fullname
	}
	if request.Location != "" {
		profile.Location = request.Location
	}

	if request.Image != "" {
		profile.Image = resp.SecureURL
	}

	data, err := h.UserRepository.UpdateUser(profile)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Code: http.StatusOK, Data: data})
}

func (h *HandlerUser) DeleteUser(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	user, err := h.UserRepository.GetUser(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
	}

	data, err := h.UserRepository.DeleteUser(user, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Code: http.StatusOK, Data: data})
}

// convert models.user to userdto.UserResponse . n make it flexible
