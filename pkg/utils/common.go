package utils

import "strconv"

func StringToInt64(str string) (int64, error) {
	i64, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0, err
	}
	return i64, nil
}

func Int64ToString(i64 int64) string {
	return strconv.FormatInt(i64, 10)
}
