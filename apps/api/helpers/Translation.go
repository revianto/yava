package helpers

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Translations menyimpan semua terjemahan
// Format key: "locale/filename.key" contoh: "id/validation.required"
var Translations = map[string]interface{}{}

// LangPath path ke folder lang (default: "lang")
var LangPath = "lang"

// Trans mengambil terjemahan berdasarkan locale dan key
// locale: kode bahasa (id, en, dll)
// key: key terjemahan format "filename.key" (contoh: "validation.required", "error.not_found")
// attributes: map untuk placeholder replacement
//
// Placeholder yang didukung:
//   - :name -> akan diganti dengan attributes["name"]
//   - {0}, {1} -> akan diganti dengan attributes berdasarkan urutan
//
// Contoh penggunaan:
//
//	Trans("id", "validation.required", map[string]interface{}{"attribute": "email"})
//	-> "Field email wajib diisi"
//
//	Trans("id", "error.not_found", nil)
//	-> "Data tidak ditemukan"
func Trans(locale string, key string, attributes map[string]interface{}) string {
	fullKey := locale + "/" + key

	if val, ok := Translations[fullKey]; ok {
		dataLang := fmt.Sprintf("%v", val)

		if attributes != nil {
			i := 0
			for k, v := range attributes {
				// Replace :key format (Laravel style)
				dataLang = strings.ReplaceAll(dataLang, ":"+k, fmt.Sprintf("%v", v))
				// Replace {0}, {1} format (indexed)
				dataLang = strings.ReplaceAll(dataLang, "{"+fmt.Sprintf("%v", i)+"}", fmt.Sprintf("%v", v))
				i++
			}
		}

		return strings.TrimSpace(dataLang)
	}

	// Fallback: return key jika tidak ditemukan
	return key
}

// TransWithFallback sama seperti Trans tapi dengan fallback ke bahasa default
func TransWithFallback(locale string, key string, attributes map[string]interface{}, fallbackLocale string) string {
	result := Trans(locale, key, attributes)

	// Jika hasil sama dengan key (tidak ditemukan), coba fallback
	if result == key && fallbackLocale != locale {
		result = Trans(fallbackLocale, key, attributes)
	}

	return result
}

// LoadTranslations memuat semua file translation dari folder lang
// Struktur folder: lang/{locale}/{filename}.json
// Contoh: lang/id/validation.json, lang/en/error.json
func LoadTranslations() error {
	return LoadTranslationsFrom(LangPath)
}

// LoadTranslationsFrom memuat translations dari path tertentu
func LoadTranslationsFrom(basePath string) error {
	// Baca semua folder locale di dalam lang/
	locales, err := os.ReadDir(basePath)
	if err != nil {
		return fmt.Errorf("failed to read lang directory: %w", err)
	}

	for _, localeDir := range locales {
		if !localeDir.IsDir() {
			continue
		}

		locale := localeDir.Name()
		localePath := filepath.Join(basePath, locale)

		// Baca semua file JSON di folder locale
		files, err := os.ReadDir(localePath)
		if err != nil {
			continue
		}

		for _, file := range files {
			if file.IsDir() || !strings.HasSuffix(file.Name(), ".json") {
				continue
			}

			// Load file JSON
			filePath := filepath.Join(localePath, file.Name())
			if err := loadTranslationFile(locale, filePath); err != nil {
				return fmt.Errorf("failed to load %s: %w", filePath, err)
			}
		}
	}

	return nil
}

// loadTranslationFile memuat satu file translation JSON
func loadTranslationFile(locale, filePath string) error {
	// Baca file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	// Parse JSON
	var translations map[string]interface{}
	if err := json.Unmarshal(data, &translations); err != nil {
		return fmt.Errorf("invalid JSON in %s: %w", filePath, err)
	}

	// Ambil nama file tanpa extension sebagai prefix
	baseName := filepath.Base(filePath)
	prefix := strings.TrimSuffix(baseName, filepath.Ext(baseName))

	// Register setiap translation dengan format: locale/prefix.key
	for key, value := range translations {
		fullKey := locale + "/" + prefix + "." + key
		Translations[fullKey] = value
	}

	return nil
}

// RegisterTranslation untuk mendaftarkan terjemahan secara programmatic
// locale: kode bahasa
// key: key terjemahan (format: "filename.key")
// value: nilai terjemahan
func RegisterTranslation(locale, key, value string) {
	Translations[locale+"/"+key] = value
}

// RegisterTranslations untuk mendaftarkan banyak terjemahan sekaligus
// locale: kode bahasa
// prefix: prefix untuk key (biasanya nama file, contoh: "validation")
// translations: map[key]value
func RegisterTranslations(locale, prefix string, translations map[string]string) {
	for key, value := range translations {
		fullKey := locale + "/" + prefix + "." + key
		Translations[fullKey] = value
	}
}

// TransError mengambil pesan error dari errors.json berdasarkan locale dan field (error tag).
// Fallback: locale lain → "id" → key mentah.
//
// Contoh:
//
//	TransError("id", "required", map[string]interface{}{"attribute": "nama"})
//	→ "Field nama wajib diisi"
//
//	TransError("id", "min", map[string]interface{}{"attribute": "password", "param": "6"})
//	→ "Field password minimal 6 karakter"
func TransError(locale, field string, attributes map[string]any) string {
	key := "errors." + field

	msg := Trans(locale, key, attributes)
	if msg != key {
		return msg
	}

	// Fallback ke "id"
	if locale != "id" {
		msg = Trans("id", key, attributes)
		if msg != key {
			return msg
		}
	}

	// Fallback ke default error message
	defaultMsg := Trans(locale, "errors.default", attributes)
	if defaultMsg != "errors.default" {
		return defaultMsg
	}

	return field
}

// Err shorthand untuk TransError dengan satu attribute field.
//
// Contoh:
//
//	Err("id", "required", "nama")
//	→ "Field nama wajib diisi"
func Err(locale, field, attribute string) string {
	return TransError(locale, field, map[string]any{"attribute": attribute})
}

// GetValidationMessage mengambil pesan validasi berdasarkan locale dan tag
func GetValidationMessage(locale, tag, field, param string) string {
	key := "validations." + tag

	// Coba ambil dengan Trans
	msg := Trans(locale, key, map[string]interface{}{
		"attribute": field,
		"param":     param,
	})

	// Jika tidak ditemukan, coba fallback ke English
	if msg == key {
		msg = Trans("en", key, map[string]interface{}{
			"attribute": field,
			"param":     param,
		})
	}

	// Jika masih tidak ditemukan, gunakan default
	if msg == key {
		defaultMsg := Trans(locale, "validations.default", map[string]interface{}{
			"attribute": field,
		})
		if defaultMsg != "validations.default" {
			msg = defaultMsg
		}
	}

	return msg
}

// ReloadTranslations reload semua translations (berguna untuk hot reload)
func ReloadTranslations() error {
	// Clear existing translations
	Translations = make(map[string]interface{})
	return LoadTranslations()
}

// GetAllTranslations return semua translations (untuk debugging)
func GetAllTranslations() map[string]interface{} {
	return Translations
}

// HasTranslation cek apakah translation key ada
func HasTranslation(locale, key string) bool {
	_, ok := Translations[locale+"/"+key]
	return ok
}
