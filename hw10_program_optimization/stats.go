package hw10programoptimization

import (
	"encoding/json"
	"io"
	"strings"
)

type User struct {
	Email string `json:"email"`
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat)
	decoder := json.NewDecoder(r)

	for {
		var user User

		if err := decoder.Decode(&user); err == io.EOF {
			break
		} else if err != nil {
			continue
		}

		email := strings.ToLower(user.Email)
		if idx := strings.LastIndex(email, "@"); idx != -1 {
			domainPart := email[idx+1:]

			if strings.HasSuffix(domainPart, domain) {
				result[domainPart]++
			}
		}
	}

	return result, nil
}
