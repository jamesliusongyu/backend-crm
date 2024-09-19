package config

import (
	"log"
	"strings"
)

// type string string

const (
	CUSTOMERA     string = "customerA"
	CUSTOMERB     string = "customerB"
	BLUSHIPPING   string = "bluShipping"
	COLUMBUS_TEST string = "columbusTest"
	UNKNOWN       string = "unknown"
)

func MakeMapping(email string) string {
	log.Println(email)
	mapping := map[string]string{
		"@customera.com":   CUSTOMERA,
		"@customerb.com":   CUSTOMERB,
		"@blushipping.com": BLUSHIPPING,
		"@outlook.com":     COLUMBUS_TEST,
	}

	for domain, tenant := range mapping {
		if strings.HasSuffix(email, domain) {
			return tenant
		}
	}
	return UNKNOWN
}
