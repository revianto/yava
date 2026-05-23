package helpers

import (
	"errors"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/revianto/yava/api/exceptions"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func init() {
	// Set timezone Jakarta untuk konsistensi test
	loc, _ := time.LoadLocation("Asia/Jakarta")
	time.Local = loc
}

func TestConv_Default(t *testing.T) {
	assert.Equal(t, "default", Conv(nil).Default("default").String())
	assert.Equal(t, "value", Conv("value").Default("default").String())
}

func TestType_String(t *testing.T) {
	assert.Equal(t, "test", Conv("test").String())
	assert.Equal(t, "123", Conv(123).String())
	assert.Equal(t, "true", Conv(true).String())
	valStr := "pointer"
	assert.Equal(t, "pointer", Conv(&valStr).String())

	now := time.Now()
	assert.Equal(t, now.Format("2006-01-02 15:04:05"), Conv(now).String())
}

func TestType_Int(t *testing.T) {
	assert.Equal(t, 123, Conv("123").Int())
	assert.Equal(t, 123, Conv(123).Int())
	assert.Equal(t, 1, Conv(true).Int())
	assert.Equal(t, 0, Conv(false).Int())
	assert.Equal(t, 123, Conv(123.45).Int())
}

func TestType_Float(t *testing.T) {
	assert.Equal(t, 123.45, Conv(123.45).Float())
	assert.Equal(t, 123.45, Conv("123.45").Float())
}

func TestType_Bool(t *testing.T) {
	assert.True(t, Conv(true).Bool())
	assert.True(t, Conv("true").Bool())
	assert.True(t, Conv("1").Bool())
	assert.True(t, Conv(1).Bool())

	assert.False(t, Conv(false).Bool())
	assert.False(t, Conv("false").Bool())
	assert.False(t, Conv("0").Bool())
	assert.False(t, Conv(0).Bool())
}

func TestType_Time(t *testing.T) {
	// Test format ISO8601 dengan Timezone (Format RFC3339)
	ts := "2026-02-01T14:00:00+07:00"
	parsed := Conv(ts).Time()
	assert.Equal(t, 2026, parsed.Year())
	assert.Equal(t, time.Month(2), parsed.Month())
	assert.Equal(t, 1, parsed.Day())
	assert.Equal(t, 14, parsed.Hour())

	// Test RAW MySQL Format (Tanpa Timezone) -> Harus dianggap Local (WIB)
	tsRaw := "2026-02-01 14:00:00"
	parsedRaw := Conv(tsRaw).Time()
	assert.Equal(t, 14, parsedRaw.Hour())
	_, offset := parsedRaw.Zone()
	assert.Equal(t, 7*3600, offset, "Harus offset +07:00 (WIB)")
}

func TestType_TimeFormat(t *testing.T) {
	// Harusnya output RFC3339
	now := time.Date(2026, 2, 1, 14, 0, 0, 0, time.Local)
	expected := now.Format(time.RFC3339)
	assert.Equal(t, expected, Conv(now).TimeFormat())
}

func TestType_Decimal(t *testing.T) {
	d := decimal.NewFromFloat(123.45)
	assert.Equal(t, d.String(), Conv(123.45).Decimal().String())
	assert.Equal(t, d.String(), Conv("123.45").Decimal().String())
}

func TestType_MapAny(t *testing.T) {
	jsonStr := `{"key": "value"}`
	m := Conv(jsonStr).MapAny()
	assert.NotNil(t, m)
	assert.Equal(t, "value", m["key"])

	fiberMap := fiber.Map{"key": "value"}
	m2 := Conv(fiberMap).MapAny()
	assert.Equal(t, "value", m2["key"])
}

func TestType_SliceMapAny(t *testing.T) {
	jsonStr := `[{"key": "value"}]`
	s := Conv(jsonStr).SliceMapAny()
	assert.Len(t, s, 1)
	assert.Equal(t, "value", s[0]["key"])
}

func TestType_IsEmpty(t *testing.T) {
	assert.True(t, Conv("").IsEmpty())
	assert.True(t, Conv(nil).IsEmpty())
	assert.True(t, Conv(0).IsEmpty())
	assert.False(t, Conv("a").IsEmpty())
	assert.False(t, Conv(1).IsEmpty())
}

func TestType_IsNil(t *testing.T) {
	assert.True(t, Conv(nil).IsNil())
	assert.False(t, Conv("").IsNil())
}

func TestType_StartEndDateTime(t *testing.T) {
	dateStr := "2026-02-01"

	// StartOfDay
	start := Conv(dateStr).StartDateTime()
	assert.Equal(t, "2026-02-01 00:00:00", start)

	// EndOfDay
	end := Conv(dateStr).EndDateTime()
	assert.Equal(t, "2026-02-01 23:59:59", end)

	// StartOfDay - Time
	startTime := Conv(dateStr).StartDateTimeTime()
	assert.Equal(t, 2026, startTime.Year())
	assert.Equal(t, 0, startTime.Hour())
	assert.Equal(t, 0, startTime.Minute())
	assert.Equal(t, 0, startTime.Second())

	// EndOfDay - Time
	endTime := Conv(dateStr).EndDateTimeTime()
	assert.Equal(t, 2026, endTime.Year())
	assert.Equal(t, 23, endTime.Hour())
	assert.Equal(t, 59, endTime.Minute())
	assert.Equal(t, 59, endTime.Second())
}

func TestType_Error(t *testing.T) {
	err := errors.New("test error")
	assert.Equal(t, err, Conv(err).Error())
	assert.Equal(t, "test string", Conv("test string").Error().Error())

	appErr := exceptions.AppError{Error: "app error"}
	assert.Equal(t, "app error", Conv(appErr).Error().Error())
}

func TestType_GetMapValueAsString(t *testing.T) {
	m := map[string]interface{}{
		"name": "Olsera",
		"age":  10,
	}

	assert.Equal(t, "Olsera", Conv(m).GetMapValueAsString("name", "Default"))
	assert.Equal(t, "10", Conv(m).GetMapValueAsString("age", "0"))
	assert.Equal(t, "Default", Conv(m).GetMapValueAsString("missing", "Default"))
}
