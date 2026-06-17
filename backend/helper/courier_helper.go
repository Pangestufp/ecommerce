package helper

import (
	"backend/dto"
	"regexp"
	"strings"
)

var allowedServices = map[string][]string{
	"jne": {
		"OKE",
		"REG",
		"YES",
		"CTC",
		"CTCYES",
	},

	"sicepat": {
		"REG",
		"BEST",
		"HALU",
	},

	"jnt": {
		"EZ",
	},

	"anteraja": {
		"REG",
		"ND",
		"SD",
	},

	"pos": {
		"Pos Reguler",
		"Pos Nextday",
	},

	"tiki": {
		"REG",
		"ONS",
		"ECO",
	},

	"rex": {
		"REG",
		"EXP",
	},

	"ninja": {
		"STANDARD",
	},

	"ide": {
		"STD",
	},

	"wahana": {},
	"lion":   {},
	"sap":    {},
	"ncs":    {},
	"rpx":    {},
	"dse":    {},
}

var serviceYearPattern = regexp.MustCompile(`\d{4}`)

func IsAllowedCourier(code string, service string) bool {

	code = strings.ToLower(
		strings.TrimSpace(code),
	)

	service = strings.TrimSpace(
		service,
	)

	if services, ok := allowedServices[code]; ok {

		for _, s := range services {

			if strings.EqualFold(
				s,
				service,
			) {
				return true
			}
		}

		return false
	}

	if service == "" {
		return false
	}

	if serviceYearPattern.MatchString(service) {
		return false
	}

	upper := strings.ToUpper(
		service,
	)

	if strings.Contains(
		upper,
		"_REV",
	) {
		return false
	}

	return true
}

func GetShippingGroup(service string) string {
	switch {

	case strings.EqualFold(
		service,
		"SD",
	):
		return "same_day"

	case strings.EqualFold(
		service,
		"SAME DAY",
	):
		return "same_day"

	case strings.EqualFold(
		service,
		"ND",
	):
		return "next_day"

	case strings.EqualFold(
		service,
		"YES",
	):
		return "next_day"

	default:
		return "regular"
	}
}

func BuildDisplayName(name string, service string) string {

	if service == "" {
		return name
	}

	return name + " " + service
}

func MarkRecommended(items []dto.ShippingOption) {
	if len(items) == 0 {
		return
	}

	lowest := items[0].Cost

	for i := range items {
		if items[i].Cost <= lowest+5000 {
			items[i].IsRecommended = true
		}
	}
}
