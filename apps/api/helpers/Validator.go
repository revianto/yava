package helpers

import (
	"reflect"
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
)

// CustomValidators menyimpan semua custom validator yang terdaftar
var CustomValidators = map[string]validator.Func{}

// RegisterCustomValidator mendaftarkan custom validator
func RegisterCustomValidator(tag string, fn validator.Func) {
	CustomValidators[tag] = fn
}

// ApplyCustomValidators mengaplikasikan semua custom validator ke instance validator
func ApplyCustomValidators(v *validator.Validate) {
	for tag, fn := range CustomValidators {
		v.RegisterValidation(tag, fn)
	}
}

// LoadDefaultCustomValidators memuat validator custom bawaan
func LoadDefaultCustomValidators() {
	// isDate - validasi format tanggal YYYY-MM-DD
	RegisterCustomValidator("isDate", func(fl validator.FieldLevel) bool {
		value := fl.Field().String()
		if value == "" {
			return true // biarkan required yang handle jika wajib
		}

		// Format yang valid: YYYY-MM-DD
		_, err := time.Parse("2006-01-02", value)
		return err == nil
	})

	// isDateTime - validasi format datetime YYYY-MM-DD HH:MM:SS
	RegisterCustomValidator("isDateTime", func(fl validator.FieldLevel) bool {
		value := fl.Field().String()
		if value == "" {
			return true
		}

		// Try multiple datetime formats
		formats := []string{
			time.RFC3339,
			time.RFC3339Nano,
		}

		for _, format := range formats {
			if _, err := time.Parse(format, value); err == nil {
				return true
			}
		}
		return false
	})

	// isTime - validasi format waktu HH:MM:SS
	RegisterCustomValidator("isTime", func(fl validator.FieldLevel) bool {
		value := fl.Field().String()
		if value == "" {
			return true
		}

		formats := []string{
			"15:04:05",
			"15:04",
		}

		for _, format := range formats {
			if _, err := time.Parse(format, value); err == nil {
				return true
			}
		}
		return false
	})

	// isPhone - validasi format nomor telepon Indonesia
	RegisterCustomValidator("isPhone", func(fl validator.FieldLevel) bool {
		value := fl.Field().String()
		if value == "" {
			return true
		}

		// Format: 08xxx, +62xxx, 62xxx
		pattern := `^(\+62|62|0)[0-9]{9,13}$`
		matched, _ := regexp.MatchString(pattern, value)
		return matched
	})

	// isNIK - validasi NIK Indonesia (16 digit)
	RegisterCustomValidator("isNIK", func(fl validator.FieldLevel) bool {
		value := fl.Field().String()
		if value == "" {
			return true
		}

		pattern := `^[0-9]{16}$`
		matched, _ := regexp.MatchString(pattern, value)
		return matched
	})

	// isNPWP - validasi NPWP Indonesia
	RegisterCustomValidator("isNPWP", func(fl validator.FieldLevel) bool {
		value := fl.Field().String()
		if value == "" {
			return true
		}

		// Format: XX.XXX.XXX.X-XXX.XXX atau 15 digit angka
		pattern := `^[0-9]{15}$|^[0-9]{2}\.[0-9]{3}\.[0-9]{3}\.[0-9]-[0-9]{3}\.[0-9]{3}$`
		matched, _ := regexp.MatchString(pattern, value)
		return matched
	})

	// isSlug - validasi format slug (lowercase, dash separated)
	RegisterCustomValidator("isSlug", func(fl validator.FieldLevel) bool {
		value := fl.Field().String()
		if value == "" {
			return true
		}

		pattern := `^[a-z0-9]+(-[a-z0-9]+)*$`
		matched, _ := regexp.MatchString(pattern, value)
		return matched
	})

	// notEmpty - berbeda dengan required, ini validasi string tidak kosong setelah trim
	RegisterCustomValidator("notEmpty", func(fl validator.FieldLevel) bool {
		if fl.Field().Kind() != reflect.String {
			return true
		}
		value := fl.Field().String()
		return len(value) > 0 && value != " "
	})

	// isPositive - validasi angka positif (> 0)
	RegisterCustomValidator("isPositive", func(fl validator.FieldLevel) bool {
		switch fl.Field().Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return fl.Field().Int() > 0
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return fl.Field().Uint() > 0
		case reflect.Float32, reflect.Float64:
			return fl.Field().Float() > 0
		}
		return false
	})

	// isNonNegative - validasi angka >= 0
	RegisterCustomValidator("isNonNegative", func(fl validator.FieldLevel) bool {
		switch fl.Field().Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return fl.Field().Int() >= 0
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return true // uint selalu >= 0
		case reflect.Float32, reflect.Float64:
			return fl.Field().Float() >= 0
		}
		return false
	})

	// isTimezone - validasi IANA timezone (e.g., "Asia/Jakarta", "Asia/Makassar")
	RegisterCustomValidator("isTimezone", func(fl validator.FieldLevel) bool {
		value := fl.Field().String()
		if value == "" {
			return true
		}

		_, err := time.LoadLocation(value)
		return err == nil
	})

	// isTimeWithZone - validasi waktu dengan timezone offset (e.g., "15:04:05+07:00")
	RegisterCustomValidator("isTimeWithZone", func(fl validator.FieldLevel) bool {
		value := fl.Field().String()
		if value == "" {
			return true
		}

		formats := []string{
			"15:04:05Z07:00", // 14:30:00+07:00
			"15:04:05-07:00", // 14:30:00-07:00
			"15:04:05Z",      // 14:30:00Z (UTC)
			"15:04Z07:00",    // 14:30+07:00
			"15:04-07:00",    // 14:30-07:00
		}

		for _, format := range formats {
			if _, err := time.Parse(format, value); err == nil {
				return true
			}
		}
		return false
	})
}
