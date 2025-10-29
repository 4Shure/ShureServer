package utils

import (
	"reflect"
	"strings"
	"time"
)

func FormatEpoch(millis int64) string {
	return time.UnixMilli(millis).
		UTC().
		Format(time.RFC3339)
}

func NowUTC() int64 {
	return time.Now().
		UTC().
		UnixMilli()
}

func FromEpoch(rfc string) (int64, error) {
	t, err := time.Parse(time.RFC3339, rfc)
	if err != nil {
		return 0, err
	}
	return t.UnixMilli(), nil
}

// IsHourExact checks if the given epoch milliseconds represents
// an exact hour (e.g., 14:00:00.000).
func IsHourExact(millis int64) bool {
	// An exact hour is perfectly divisible by the number
	// of milliseconds in one hour.
	const millisInHour = 3600000
	return millis%millisInHour == 0
}

func Sanitize(o any) {
	v := reflect.ValueOf(o)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		panic("sanitize: expected pointer to struct")
	}

	v = v.Elem()
	if v.Kind() != reflect.Struct {
		panic("sanitize: expected struct")
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		switch field.Kind() {
		case reflect.String:
			field.SetString(sanitizeString(field.String()))

		case reflect.Slice:
			if field.Type().Elem().Kind() == reflect.String {
				for j := 0; j < field.Len(); j++ {
					field.Index(j).SetString(sanitizeString(field.Index(j).String()))
				}
			}
		}
	}
}

func sanitizeString(s string) string {
	return strings.TrimSpace(s)
}
