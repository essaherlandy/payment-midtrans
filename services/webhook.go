package services

import (
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func Webhook(c *fiber.Ctx) error {
	var payload map[string]interface{}

	if err := json.Unmarshal(c.Body(), &payload); err != nil {
		return c.Status(400).SendString("invalid payload")
	}

	orderID := payload["order_id"].(string)
	status := payload["transaction_status"].(string)
	gross := payload["gross_amount"].(string)
	signature := payload["signature_key"].(string)

	// VALIDASI SIGNATURE (WAJIB)
	raw := orderID + gross + status + serverKey
	hash := sha512.Sum512([]byte(raw))
	expected := hex.EncodeToString(hash[:])

	if signature != expected {
		return c.Status(403).SendString("invalid signature")
	}

	// HANDLE STATUS
	switch status {
	case "settlement":
		fmt.Println("PAID:", orderID)
	case "pending":
		fmt.Println("PENDING:", orderID)
	case "expire":
		fmt.Println("EXPIRED:", orderID)
	case "cancel":
		fmt.Println("CANCELLED:", orderID)
	}

	return c.SendString("OK")
}
