package converters

import "strconv"

func ConvertStringToFloat64(value string) (float64, error) {
	return strconv.ParseFloat(value, 64)
}

func ConvertStringToInt64(value string) (int64, error) {
	return strconv.ParseInt(value, 10, 64)
}
