package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/goccy/go-json"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/revianto/yava/api/app/controllers"
	"github.com/revianto/yava/api/app/middlewares"
	"github.com/revianto/yava/api/exceptions"
	"github.com/revianto/yava/api/helpers"
	"github.com/revianto/yava/api/routes"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {

	// Set Timezone
	tz := helpers.GetEnv("APP_TIMEZONE")
	if tz == "" {
		tz = "UTC"
	}
	loc, err := time.LoadLocation(tz)
	if err != nil {
		log.Fatal(err)
	}
	time.Local = loc

	// connect DB
	ConnDB()

	// Load translations
	helpers.LoadTranslations()

	// Load custom validators
	helpers.LoadDefaultCustomValidators()

	appName := helpers.GetEnv("APP_NAME")
	if appName == "" {
		appName = "Go Framework"
	}
	app := fiber.New(fiber.Config{
		AppName:           appName,
		CaseSensitive:     true,
		StrictRouting:     true,
		EnablePrintRoutes: true,
		JSONEncoder:       json.Marshal,
		JSONDecoder:       json.Unmarshal,
		ErrorHandler:      ErrorHandler,
	})

	app.Use(recover.New())

	// CORS Middleware
	allowedOrigins := helpers.GetEnv("CORS_ORIGINS")
	if allowedOrigins == "" {
		allowedOrigins = "http://localhost:3000,http://localhost:5173"
	}

	app.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, X-Requested-With",
		AllowMethods:     "GET, POST, HEAD, PUT, DELETE, PATCH, OPTIONS",
		AllowCredentials: true,
	}))

	app.Use(middlewares.SetContexFiber)

	// protected static files
	app.Use("/template", middlewares.CheckClientToken())
	app.Static("/template", "./public/template")
	app.Use("/image", middlewares.CheckClientToken())
	app.Static("/image", "./public/img", fiber.Static{
		Compress:      true,           // Gzip compression
		CacheDuration: 24 * time.Hour, // Browser caching
	})

	app.Get("/", middlewares.SetContexFiber, func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	admin := app.Group("/api/admin/:locale")
	admin.Use(middlewares.ValidateLocale)
	{
		routes.Api(admin)
	}

	// 404 Handler
	app.Use(func(c *fiber.Ctx) error {
		return exceptions.ResponseErrorException(c, exceptions.PageNotFoundErrorException(c, "Page Not Found")) // => 404 "Not Found"
	})

	log.Fatal(app.Listen(helpers.GetEnv("PORT")))
}

func ErrorHandler(c *fiber.Ctx, err error) error {
	var e *fiber.Error
	if errors.As(err, &e) {
		switch e.Code {
		case fiber.StatusBadRequest: // 400
			return exceptions.ResponseErrorException(c, exceptions.ErrorException(c, 400, e.Message))
		case fiber.StatusUnauthorized: // 401
			return exceptions.ResponseErrorException(c, exceptions.AuthException(c, 401, "Unauthorized"))
		case fiber.StatusForbidden: // 403
			return exceptions.ResponseErrorException(c, exceptions.AuthException(c, 403, "Forbidden"))
		case fiber.StatusNotFound: // 404
			return exceptions.ResponseErrorException(c, exceptions.PageNotFoundErrorException(c, "Page Not Found"))
		case fiber.StatusMethodNotAllowed: // 405
			return exceptions.ResponseErrorException(c, exceptions.ErrorException(c, 405, "Method Not Allowed"))
		case fiber.StatusUnprocessableEntity: // 422
			return exceptions.ResponseErrorException(c, exceptions.ValidateException(c, 422, e.Message))
		case fiber.StatusTooManyRequests: // 429
			return exceptions.ResponseErrorException(c, exceptions.ThrottleException(c, 429, "Too Many Requests"))
		case fiber.StatusServiceUnavailable: // 503
			return exceptions.ResponseErrorException(c, exceptions.ErrorException(c, 503, "Service Unavailable"))
		}
	}

	// Get caller info for debugging
	_, file, line, _ := runtime.Caller(1)
	log.Printf("Internal Error at %s:%d - %s", file, line, err.Error())
	return exceptions.ResponseErrorException(c, exceptions.InternalErrorException(c, "Internal Server Error"))
}


func ConnDB() {
	sslmode := helpers.GetEnv("DB_SSLMODE")
	if sslmode == "" {
		sslmode = "disable"
	}
	dbTz := helpers.GetEnv("APP_TIMEZONE")
	if dbTz == "" {
		dbTz = "UTC"
	}
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s prefer_simple_protocol=true",
		helpers.GetEnv("DB_HOST"),
		helpers.GetEnv("DB_PORT"),
		helpers.GetEnv("DB_USERNAME"),
		helpers.GetEnv("DB_PASSWORD"),
		helpers.GetEnv("DB_DATABASE"),
		sslmode,
		dbTz,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Buat Callback Function
	enforceTx := func(db *gorm.DB) {
		// Cek tipe ConnPool (Connection Pool)
		// Jika ConnPool adalah *sql.Tx, berarti sedang dalam transaction.
		// Jika *sql.DB, berarti koneksi biasa (non-transaction).
		switch db.Statement.ConnPool.(type) {
		case *sql.Tx:
			// Aman, ini transaction
			return
		case *gorm.PreparedStmtTX:
			// Aman (Transaction dengan Prepared Statement aktif)
			return
		default:
			// Bahaya, ini bukan transaction!
			db.AddError(errors.New("UNSAFE OPERATION: Create/Update/Delete must be within a transaction!"))
		}
	}

	// Daftarkan Callback ke event Create, Update, dan Delete
	// GORM akan menjalankan fungsi ini SEBELUM query dieksekusi
	db.Callback().Create().Before("gorm:create").Register("enforce_tx", enforceTx)
	db.Callback().Update().Before("gorm:update").Register("enforce_tx", enforceTx)
	db.Callback().Delete().Before("gorm:delete").Register("enforce_tx", enforceTx)

	controllers.DB = db
}
