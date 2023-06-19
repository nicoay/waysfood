package handlers

import (
	"fmt"
	"log"
	dto "mytask/dto/result"
	transactiondto "mytask/dto/transaction"
	"mytask/models"
	repositories "mytask/repository"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-playground/validator"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"gopkg.in/gomail.v2"
)

type handlerTransaction struct {
	TransactionRepo repositories.TransactionsRepo
	OrderRepo       repositories.OrderRepo
	CartRepo        repositories.CartRepo
}

func HandlerTransaction(TransactionRepo repositories.TransactionsRepo, OrderRepo repositories.OrderRepo,
	CartRepo repositories.CartRepo) *handlerTransaction {
	return &handlerTransaction{TransactionRepo, OrderRepo, CartRepo}
}

func (h *handlerTransaction) FindTransaction(c echo.Context) error {
	Transactions, err := h.TransactionRepo.FindTransaction()
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Code: http.StatusOK, Data: Transactions})
}

func (h *handlerTransaction) GetTransaction(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	Transaction, err := h.TransactionRepo.GetTransaction(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Code: http.StatusOK, Data: Transaction})
}

func (h *handlerTransaction) GetUserTransaction(c echo.Context) error {
	userLogin := c.Get("userLogin")
	Id := userLogin.(jwt.MapClaims)["id"].(float64)

	Transaction, err := h.TransactionRepo.GetUserTransaction(int(Id))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Code: http.StatusOK, Data: Transaction})
}

func (h *handlerTransaction) GetPartnerTransaction(c echo.Context) error {
	userLogin := c.Get("userLogin")
	Id := userLogin.(jwt.MapClaims)["id"].(float64)

	Transaction, err := h.TransactionRepo.GetPartnerTransaction(int(Id))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Code: http.StatusOK, Data: Transaction})
}

func (h *handlerTransaction) CreateTransaction(c echo.Context) error {

	request := new(transactiondto.CreateTransactionReq)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
	}

	userLogin := c.Get("userLogin")
	buyerId := userLogin.(jwt.MapClaims)["id"].(float64)

	validation := validator.New()
	err := validation.Struct(request)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
	}

	var transactionIsMatch = false
	var transactionId int

	for !transactionIsMatch {
		transactionId = int(time.Now().Unix())
		transactionData, _ := h.TransactionRepo.GetTransaction(transactionId)
		if transactionData.ID == 0 {
			transactionIsMatch = true
		}
	}

	// Membuat struct Transaction dari data permintaan
	transaction := models.Transaction{
		ID:         transactionId,
		BuyerID:    int(buyerId),
		SellerID:   request.SellerID,
		TotalPrice: request.TotalPrice,
		Status:     "pending",
	}

	dataTransaction, err := h.TransactionRepo.CreateTransaction(transaction)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()})
	}

	orders, err := h.OrderRepo.GetOrderbyUser(int(buyerId))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()})
	}

	for _, order := range orders {
		cart := models.Carts{
			Qty:           order.Qty,
			ProductID:     order.ProductID,
			TransactionID: dataTransaction.ID,
		}
		h.CartRepo.CreateCarts(cart)
	}
	fmt.Println(transaction.TotalPrice, "haikal kontol")
	var s = snap.Client{}
	s.New(os.Getenv("SERVER_KEY"), midtrans.Sandbox)

	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  strconv.Itoa(dataTransaction.ID),
			GrossAmt: transaction.TotalPrice,
		},
		CreditCard: &snap.CreditCardDetails{
			Secure: true,
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: dataTransaction.Buyer.Fullname,
			Email: dataTransaction.Buyer.Email,
		},
	}

	snapResp, _ := s.CreateTransaction(req)

	return c.JSON(http.StatusOK, dto.SuccessResult{Code: http.StatusOK, Data: snapResp})
}

func (h *handlerTransaction) Notification(c echo.Context) error {
	var notificationPayload map[string]interface{}

	if err := c.Bind(&notificationPayload); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
	}

	transactionStatus := notificationPayload["transaction_status"].(string)
	fraudStatus := notificationPayload["fraud_status"].(string)
	orderId := notificationPayload["order_id"].(string)

	order_id, _ := strconv.Atoi(orderId)
	transaction, _ := h.TransactionRepo.GetTransaction(order_id)

	fmt.Print("ini payloadnya", notificationPayload)

	if transactionStatus == "capture" {
		if fraudStatus == "challenge" {
			// TODO set transaction status on your database to 'challenge'
			// e.g: 'Payment status challenged. Please take action on your Merchant Administration Portal
			h.TransactionRepo.UpdateTransaction("PENDING", order_id)
		} else if fraudStatus == "accept" {
			// TODO set transaction status on your database to 'success'
			SendMail("SUCCESS", transaction)
			h.TransactionRepo.UpdateTransaction("SUCCESS", order_id)
		}
	} else if transactionStatus == "settlement" {
		// TODO set transaction status on your databaase to 'SUCCESS'
		SendMail("SUCCESS", transaction)
		h.TransactionRepo.UpdateTransaction("SUCCESS", order_id)
	} else if transactionStatus == "deny" {
		// TODO you can ignore 'deny', because most of the time it allows payment retries
		// and later can become success
		h.TransactionRepo.UpdateTransaction("FAILED", order_id)
	} else if transactionStatus == "CANCEL" || transactionStatus == "expire" {
		// TODO set transaction status on your databaase to 'failure'
		h.TransactionRepo.UpdateTransaction("FAILED", order_id)
	} else if transactionStatus == "PENDING" {
		// TODO set transaction status on your databaase to 'PENDING' / waiting payment
		h.TransactionRepo.UpdateTransaction("PENDING", order_id)
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Code: http.StatusOK, Data: notificationPayload})
}

func SendMail(status string, transaction models.Transaction) {

	if status != transaction.Status && (status == "SUCCEESS") {
		var CONFIG_SMTP_HOST = "smtp.gmail.com"
		var CONFIG_SMTP_PORT = 587
		var CONFIG_SENDER_NAME = "DEWETOUR <demo.misaeltimpolas04@gmail.com>"
		var CONFIG_AUTH_EMAIL = os.Getenv("EMAIL_SYSTEM")
		var CONFIG_AUTH_PASSWORD = os.Getenv("PASSWORD_SYSTEM")

		var trip = "TRIP BOOKING"
		var price = strconv.Itoa(int(transaction.TotalPrice))

		mailer := gomail.NewMessage()
		mailer.SetHeader("From", CONFIG_SENDER_NAME)
		mailer.SetHeader("To", transaction.Buyer.Email)
		mailer.SetHeader("Subject", "Transaction Status")
		mailer.SetBody("text/html", fmt.Sprintf(`<!DOCTYPE html>
	  <html lang="en">
		<head>
		<meta charset="UTF-8" />
		<meta http-equiv="X-UA-Compatible" content="IE=edge" />
		<meta name="viewport" content="width=device-width, initial-scale=1.0" />
		<title>Document</title>
		<style>
		  h1 {
		  color: brown;
		  }
		</style>
		</head>
		<body>
		<h2>Product payment :</h2>
		<ul style="list-style-type:none;">
		  <li>Name : %s</li>
		  <li>Total payment: Rp.%s</li>
		  <li>Status : <b>%s</b></li>
		</ul>
		</body>
	  </html>`, trip, price, status))

		dialer := gomail.NewDialer(
			CONFIG_SMTP_HOST,
			CONFIG_SMTP_PORT,
			CONFIG_AUTH_EMAIL,
			CONFIG_AUTH_PASSWORD,
		)

		err := dialer.DialAndSend(mailer)
		if err != nil {
			log.Fatal(err.Error())
		}

		log.Println("Mail sent! to " + transaction.Buyer.Email)
	}
}
