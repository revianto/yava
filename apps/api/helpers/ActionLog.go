package helpers

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// WriteLog inserts a tr_action_logs row asynchronously (fire-and-forget).
// Uses raw Exec to bypass the enforceTx GORM callback (only registered on
// gorm:create/update/delete, not gorm:raw). Log failures never fail the caller.
func WriteLog(db *gorm.DB, c *fiber.Ctx, actionText string) {
	userId := GetUserId(c)
	if userId == 0 {
		return
	}
	go func() {
		err := db.Session(&gorm.Session{NewDB: true}).
			Exec(`INSERT INTO tr_action_logs (user_id, action_text, created_at)
				  VALUES (?, ?, NOW())`, userId, actionText).Error
		if err != nil {
			log.Printf("[WriteLog] %v", err)
		}
	}()
}
