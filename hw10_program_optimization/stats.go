package hw10_program_optimization //nolint:golint,stylecheck

import (
	"bufio"
	"io"
	"log"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

type User struct {
	Email string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	return countDomains(r, domain)
}

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func countDomains(r io.Reader, domain string) (DomainStat, error) {
	res := make(DomainStat)
	var user *User
	reader := bufio.NewReader(r)
	for {
		line, _, _ := reader.ReadLine()
		if line == nil {
			break
		}
		if err := json.Unmarshal(line, &user); err != nil {
			log.Println("E!", err)

			continue
		}
		matched := strings.Contains(user.Email, "."+domain)
		if matched {
			domain := strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])
			res[domain]++
			user.Email = ""
		}
	}

	return res, nil
}
