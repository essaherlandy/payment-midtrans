package services

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

var serverKey = "SB-Mid-server-UrlF5ajBqWoVjTJyimvVHdXc"

type Request struct {
	Email string `json:"email"`
}

func CreateTransaction(c *fiber.Ctx) error {
	var req Request

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid request",
		})
	}

	var snapClient snap.Client
	snapClient.New(serverKey, midtrans.Sandbox)

	orderID := fmt.Sprintf("ORDER-%d", time.Now().Unix())

	snapReq := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  orderID,
			GrossAmt: 150000,
		},
		CustomerDetail: &midtrans.CustomerDetails{
			Email: req.Email,
		},
	}

	snapResp, err := snapClient.CreateTransaction(snapReq)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"token":    snapResp.Token,
		"order_id": orderID,
	})
}

func MidtransCallback(c *fiber.Ctx) error {
	var payload map[string]interface{}

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid payload",
		})
	}

	orderID := payload["order_id"].(string)
	status := payload["transaction_status"].(string)
	gross := payload["gross_amount"].(string)
	signature := payload["signature_key"].(string)

	// 🔐 VALIDASI SIGNATURE (WAJIB)
	raw := orderID + gross + status + serverKey
	hash := sha512.Sum512([]byte(raw))
	expected := hex.EncodeToString(hash[:])

	if signature != expected {
		return c.Status(403).JSON(fiber.Map{
			"error": "invalid signature",
		})
	}

	// 💳 HANDLE STATUS PAYMENT
	switch status {

	case "settlement":
		fmt.Println("✅ PAID:", orderID)
		// TODO: update DB -> paid

	case "pending":
		fmt.Println("⏳ PENDING:", orderID)
		// TODO: update DB -> pending

	case "expire":
		fmt.Println("❌ EXPIRED:", orderID)
		// TODO: update DB -> expired

	case "cancel":
		fmt.Println("🚫 CANCELLED:", orderID)
		// TODO: update DB -> cancelled
	}

	return c.JSON(fiber.Map{
		"message": "ok",
	})
}
