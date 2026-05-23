package helpers

import (
	"encoding/json"
	"errors"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/revianto/yava/api/exceptions"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// Daftar format waktu yang didukung untuk parsing
var timeFormats = []string{
	"2006-01-02 15:04:05",
	"2006-01-02",
	"15:04:05",
	"2006-01-02 15:04",
	"2006-01-02T15:04:05Z",
	"2006-01-02T15:04:05+07:00",
	"2006-01-02T15:04:05-07:00",
	"2006-01-02T15:04:05.000Z",
	"2006-01-02T15:04:05.000+07:00",
	"2006-01-02T15:04:05.000-07:00",
	"2006-01-02T15:04:05.000000Z",
	"2006-01-02T15:04:05.000000+07:00",
	"2006-01-02T15:04:05.000000-07:00",
	"2006-01-02T15:04:05.000000000Z",
	"2006-01-02T15:04:05.000000000+07:00",
	"2006-01-02T15:04:05.000000000-07:00",
	"2006-01-02T15:04:05.000000000+0700",
	"2006-01-02T15:04:05.000000000-0700",
	"2006-01-02T15:04:05.000000000+07",
	"2006-01-02T15:04:05.000000000-07",
	time.RFC3339,
	time.RFC3339Nano,
	"2006-01-02T15:04Z07:00",
	"2006-01-02T15:04-07:00",
	"2006-01-02T15Z07:00",
	"2006-01-02T15-07:00",
	"2006-01-02T15Z",
	"2006-01-02T15",
	"2006-01-02T15:04:05.000Z07:00",
	"2006-01-02T15:04:05.000",
	"2006-01-02T15:04:05",
	"2006-01-02T15:04",
	"2006-01",
	"2006",
	"2006-01-02 15:04:05.000Z07:00",
	"2006-01-02 15:04:05.000-07:00",
	"2006-01-02 15:04:05.000Z",
	"2006-01-02 15:04:05.000",
	"2006-01-02 15",
	"2006-01-02 15:04:05.000000000Z07:00",
	"2006-01-02 15:04:05.000000000-07:00",
}

type Type struct {
	value interface{}
}

func Conv(d interface{}) *Type {
	return &Type{
		value: d,
	}
}

func (c *Type) Default(d interface{}) *Type {
	if c.value == nil {
		c.value = d
	}
	return c
}

func (c Type) Error() error {
	if c.value == nil {
		return nil
	}

	switch data := c.value.(type) {
	case error:
		return data
	case *runtime.TypeAssertionError:
		return data
	case exceptions.AppError:
		return errors.New(Conv(data.Error).String())
	case string:
		return errors.New(data)
	case *string:
		if data == nil {
			return nil
		}
		return errors.New(*data)
	case int:
		return errors.New(Conv(data).String())
	case int64:
		return errors.New(Conv(data).String())
	case float32:
		return errors.New(Conv(data).String())
	case float64:
		return errors.New(Conv(data).String())
	case decimal.Decimal:
		return errors.New(data.String())
	case time.Time:
		return errors.New(data.Format("2006-01-02 15:04:05"))
	default:
		panic("Type is " + reflect.TypeOf(data).String())
	}
}

func (c Type) String() string {
	if c.value == nil {
		return ""
	}

	switch data := c.value.(type) {
	case string:
		return data
	case []string:
		result, _ := json.Marshal(data)
		return string(result)
	case exceptions.AppError:
		return Conv(data.Error).Default("Please try again later").String()
	case *string:
		if data == nil {
			return ""
		}
		return *data
	case *strings.Reader:
		byteData, _ := data.ReadByte()
		return string(byteData)
	case time.Time:
		return data.Format("2006-01-02 15:04:05")
	case *time.Time:
		if data == nil {
			return ""
		}
		return data.Format("2006-01-02 15:04:05")
	case int:
		return strconv.Itoa(data)
	case error:
		return data.Error()
	case int8:
		return strconv.Itoa(int(data))
	case int32:
		return strconv.Itoa(int(data))
	case uint32:
		return strconv.FormatUint(uint64(data), 10)
	case int64:
		return strconv.FormatInt(data, 10)
	case uint64:
		return strconv.FormatUint(data, 10)
	case float32:
		return strconv.FormatFloat(float64(data), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(data, 'f', -1, 64)
	case []uint8:
		return string(data)
	case bool:
		return strconv.FormatBool(data)
	case decimal.Decimal:
		return data.String()
	case map[string]interface{}:
		result, _ := json.Marshal(data)
		return string(result)
	case []map[string]interface{}:
		result, _ := json.Marshal(data)
		return string(result)
	case []interface{}:
		result, _ := json.Marshal(data)
		return string(result)
	default:
		panic("Type is " + reflect.TypeOf(data).String())
	}
}

func (c Type) Int() int {
	if c.value == nil {
		return 0
	}

	switch data := c.value.(type) {
	case string:
		if data == "" {
			return 0
		}
		i, e := strconv.ParseInt(data, 10, 64)
		if e != nil {
			return int(Conv(data).Decimal().IntPart())
		}
		return int(i)
	case bool:
		if data {
			return 1
		}
		return 0
	case int:
		return data
	case int8:
		return int(data)
	case int32:
		return int(data)
	case int64:
		return int(data)
	case uint:
		return int(data)
	case uint8:
		return int(data)
	case uint32:
		return int(data)
	case uint64:
		return int(data)
	case float32:
		return int(data)
	case float64:
		return int(data)
	case decimal.Decimal:
		return int(data.IntPart())
	default:
		panic("Type is " + reflect.TypeOf(data).String())
	}
}

func (c Type) Int64() int64 {
	return int64(c.Int())
}

func (c Type) Float() float64 {
	if c.value == nil {
		return 0
	}

	switch data := c.value.(type) {
	case string:
		if data == "" {
			return 0
		}
		i, _ := strconv.ParseFloat(data, 64)
		return i
	case int:
		return float64(data)
	case int8:
		return float64(data)
	case int32:
		return float64(data)
	case int64:
		return float64(data)
	case uint:
		return float64(data)
	case uint8:
		return float64(data)
	case uint32:
		return float64(data)
	case uint64:
		return float64(data)
	case float32:
		return float64(data)
	case decimal.Decimal:
		return data.InexactFloat64()
	case float64:
		return data
	default:
		panic("Type is " + reflect.TypeOf(data).String())
	}
}

func (c Type) Decimal(quotes ...bool) decimal.Decimal {
	// Catatan: Parameter quotes diabaikan untuk menghindari race condition pada global state
	// decimal.MarshalJSONWithoutQuotes. Konfigurasi ini harus dilakukan di main level.

	if c.value == nil {
		return decimal.Zero
	}

	switch data := c.value.(type) {
	case []uint8:
		num := decimal.Zero
		json.Unmarshal(data, &num)
		return num
	case decimal.Decimal:
		return data
	case int:
		return decimal.NewFromInt(int64(data))
	case int32:
		return decimal.NewFromInt32(data)
	case int64:
		return decimal.NewFromInt(data)
	case uint:
		return decimal.NewFromInt(int64(data))
	case uint32:
		return decimal.NewFromInt(int64(data))
	case uint64:
		return decimal.NewFromInt(int64(data))
	case float32:
		return decimal.NewFromFloat32(data)
	case float64:
		return decimal.NewFromFloat(data)
	case string:
		if data == "" {
			return decimal.Zero
		}
		num, _ := decimal.NewFromString(data)
		return num
	default:
		panic("Type is " + reflect.TypeOf(data).String())
	}
}

func (c Type) Byte() []byte {
	if c.value == nil {
		return nil
	}

	switch data := c.value.(type) {
	case []byte:
		return data
	case fiber.Map:
		byteData, _ := json.Marshal(data)
		return byteData
	case string:
		return []byte(data)
	case []string:
		mapByte, _ := json.Marshal(data)
		return mapByte
	case []interface{}:
		mapByte, _ := json.Marshal(data)
		return mapByte
	case []map[string]interface{}:
		mapByte, _ := json.Marshal(data)
		return mapByte
	case *[]map[string]interface{}:
		if data == nil {
			return nil
		}
		mapByte, _ := json.Marshal(data)
		return mapByte
	case map[string]interface{}:
		mapByte, _ := json.Marshal(data)
		return mapByte
	case *map[string]interface{}:
		if data == nil {
			return nil
		}
		mapByte, _ := json.Marshal(data)
		return mapByte
	default:
		kind := reflect.ValueOf(data).Kind()
		if kind == reflect.Struct || kind == reflect.Slice || kind == reflect.Ptr {
			byteData, _ := json.Marshal(data)
			return byteData
		}
		panic("Type is " + reflect.TypeOf(data).String())
	}
}

func (c Type) Time() time.Time {
	if c.value == nil {
		return time.Time{}
	}

	switch data := c.value.(type) {
	case time.Time:
		return data
	case *time.Time:
		if data == nil {
			return time.Time{}
		}
		return *data
	case int64:
		return time.Unix(data, 0)
	case int:
		return time.Unix(int64(data), 0)
	case gorm.DeletedAt:
		return data.Time
	case string:
		if data == "" {
			return time.Time{}
		}
		// Coba semua format yang didukung
		for _, format := range timeFormats {
			if result, err := time.ParseInLocation(format, data, time.Local); err == nil {
				return result
			}
		}
		// Jika tidak ada format yang cocok, kembalikan zero time
		return time.Time{}
	default:
		return time.Time{}
	}
}

func (c Type) TimeFormat() string {
	return c.Time().Format(time.RFC3339)
}

func (c Type) Bool() bool {
	if c.value == nil {
		return false
	}

	switch data := c.value.(type) {
	case bool:
		return data
	case *bool:
		if data == nil {
			return false
		}
		return *data
	case string:
		result, err := strconv.ParseBool(data)
		if err != nil {
			// Coba parsing sebagai angka
			if num, numErr := strconv.Atoi(data); numErr == nil {
				return num != 0
			}
			return false
		}
		return result
	case int:
		return data != 0
	case int8:
		return data != 0
	case int16:
		return data != 0
	case int32:
		return data != 0
	case int64:
		return data != 0
	case uint:
		return data != 0
	case uint8:
		return data != 0
	case uint16:
		return data != 0
	case uint32:
		return data != 0
	case uint64:
		return data != 0
	case float32:
		return data != 0
	case float64:
		return data != 0
	case decimal.Decimal:
		return !data.IsZero()
	default:
		panic("Param to bool undefined, Type is " + reflect.TypeOf(data).String())
	}
}

func (c Type) GetMapValueAsString(key string, defaultValue string) string {
	mapData := c.MapAny()
	if mapData == nil {
		return defaultValue
	}
	anyData, exists := mapData[key]
	if !exists || anyData == nil {
		return defaultValue
	}
	return Conv(anyData).Default(defaultValue).String()
}

// IsNil checks if the underlying value is nil
func (c Type) IsNil() bool {
	return c.value == nil
}

// IsEmpty checks if the underlying value is empty (nil, empty string, zero, etc.)
func (c Type) IsEmpty() bool {
	if c.value == nil {
		return true
	}

	switch data := c.value.(type) {
	case string:
		return data == ""
	case []byte:
		return len(data) == 0
	case int:
		return data == 0
	case int64:
		return data == 0
	case float64:
		return data == 0
	case bool:
		return !data
	case decimal.Decimal:
		return data.IsZero()
	case time.Time:
		return data.IsZero()
	default:
		v := reflect.ValueOf(data)
		switch v.Kind() {
		case reflect.Slice, reflect.Map, reflect.Array:
			return v.Len() == 0
		case reflect.Ptr, reflect.Interface:
			return v.IsNil()
		}
		return false
	}
}

// StartDateTimeTime mengembalikan time.Time awal hari (00:00:00)
func (c Type) StartDateTimeTime() time.Time {
	t := c.Time()
	if t.IsZero() {
		return time.Time{}
	}
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// StartDateTime mengkonversi tanggal ke datetime string dengan waktu 00:00:00
func (c Type) StartDateTime() string {
	t := c.StartDateTimeTime()
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}

// EndDateTimeTime mengembalikan time.Time akhir hari (23:59:59)
func (c Type) EndDateTimeTime() time.Time {
	t := c.Time()
	if t.IsZero() {
		return time.Time{}
	}
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, t.Location())
}

// EndDateTime mengkonversi tanggal ke datetime string dengan waktu 23:59:59
func (c Type) EndDateTime() string {
	t := c.EndDateTimeTime()
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}

// unmarshalJSON mencoba parsing data ke target struct/map/slice
func (c Type) unmarshalJSON(data interface{}, target interface{}) bool {
	var bytes []byte
	switch v := data.(type) {
	case string:
		if v == "" {
			return false
		}
		bytes = []byte(v)
	case []byte:
		if len(v) == 0 {
			return false
		}
		bytes = v
	default:
		// Coba marshal dulu
		if b, err := json.Marshal(data); err == nil {
			bytes = b
		} else {
			return false
		}
	}

	return json.Unmarshal(bytes, target) == nil
}

// SliceMap mengkonversi any ke []map[string]any
func (c Type) SliceMapAny() []map[string]any {
	if c.value == nil {
		return nil
	}

	switch data := c.value.(type) {
	case []map[string]any:
		return data
	case []fiber.Map:
		result := make([]map[string]any, len(data))
		for i, m := range data {
			result[i] = map[string]any(m)
		}
		return result
	case []interface{}:
		result := make([]map[string]any, 0, len(data))
		for _, item := range data {
			if m, ok := item.(map[string]interface{}); ok {
				result = append(result, m)
			} else if m, ok := item.(map[string]any); ok {
				result = append(result, m)
			} else if m, ok := item.(fiber.Map); ok {
				result = append(result, map[string]any(m))
			}
		}
		return result
	default:
		var result []map[string]any
		if c.unmarshalJSON(data, &result) {
			return result
		}
		return nil
	}
}

// MapAny mengkonversi any ke map[string]any
func (c Type) MapAny() map[string]any {
	if c.value == nil {
		return nil
	}

	switch data := c.value.(type) {
	case map[string]any:
		return data
	case fiber.Map:
		return map[string]any(data)
	default:
		var result map[string]any
		if c.unmarshalJSON(data, &result) {
			return result
		}
		return nil
	}
}
