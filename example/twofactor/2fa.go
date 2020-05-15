package twofactor

import (
	"fmt"

	"github.com/hokaccha/go-prettyjson"
	auth "github.com/tsirysndr/2fa-auth"
)

func GenerateOtpURL(name, issuer string, secretLength int, symbols bool) {
	g, _ := auth.GenerateSecret(&auth.Options{
		Length:     secretLength,
		Symbols:    symbols,
		OtpauthURL: true,
		Name:       name,
		Issuer:     issuer,
	})
	secret, _ := prettyjson.Marshal(g)
	fmt.Println(string(secret))
}

func VerifyOTP(secret, code string) {
	authenticated, _ := auth.VerifyOTP(secret, code)
	if authenticated {
		fmt.Println("Authentication Success!")
		return
	}
	fmt.Println("Authentication Failed!")
}
