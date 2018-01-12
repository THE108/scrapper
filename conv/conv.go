package conv

import (
	"fmt"
	"strconv"
	"strings"
)

func ToMoney(s string) (float64, error) {
	s = strings.ToLower(strings.TrimSpace(s))

	if s == "free" {
		return 0., nil
	}

	s = strings.TrimSpace(strings.TrimPrefix(s, "sgd"))

	return strconv.ParseFloat(s, 64)
}

func ToDeliveryType(s string) (string, error) {
	s = strings.ToLower(s)

	if strings.Contains(s, "liveup") {
		return "liveup", nil
	}

	s = strings.Trim(s, ": ")

	switch s {
	case "standard", "express", "economy":
		return s, nil
	}

	return "", fmt.Errorf("unknown delivery type: %s", s)
}
