package format

import (
	"fmt"
	"strconv"
	"strings"
)

func CentsToDollarsI(val int) string {
	return CentsToDollars(int64(val))
}
func BpsToPercentI(val int) string {
	return BpsToPercent(int64(val))
}
func CentsToDollars(val int64) string {

	vals := strings.Split(strconv.Itoa(int(val)), "")
	if len(vals) == 1 {
		vals = append([]string{"0", "0"}, vals...)
	} else if len(vals) == 2 {
		vals = append([]string{"0"}, vals...)
	}

	return fmt.Sprintf("%s.%s", strings.Join(vals[:len(vals)-2], ""), strings.Join(vals[len(vals)-2:], ""))
}
func BpsToPercent(val int64) string {
	vals := strings.Split(strconv.Itoa(int(val)), "")
	if len(vals) <= 2 {
		return fmt.Sprintf("%d", val)
	}

	return fmt.Sprintf("%s.%s", strings.Join(vals[:len(vals)-2], ""), strings.Join(vals[len(vals)-2:], ""))
}
func Strptr(string2 string) *string {
	return &string2
}
