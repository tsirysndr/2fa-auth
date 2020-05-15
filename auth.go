package auth

import (
	"crypto/rand"
	"encoding/base32"
	"encoding/hex"
	"log"
	"math"
	"net/url"
	"strings"

	g "github.com/tsirysndr/dgoogauth"
)

type GeneratedSecret struct {
	Ascii        string `json:"ascii,omitempty"`
	Hex          string `json:"hex,omitempty"`
	Base32       string `json:"base32,omitempty"`
	QrCodeAscii  string `json:"qr_code_ascii,omitempty"`
	QrCodeHex    string `json:"qr_code_hex,omitempty"`
	QrCodeBase32 string `json:"qr_code_base32,omitempty"`
	GoogleAuthQR string `json:"google_auth_qr,omitempty"`
	OtpAuthQR    string `json:"otp_auth_qr,omitempty"`
	OtpType      string `json:"otp_type,omitempty"`
}

type Options struct {
	Length       int
	Name         string
	QrCodes      bool
	GoogleAuthQr bool
	OtpauthURL   bool
	Symbols      bool
	Issuer       string
}

type OtpauthOptions struct {
	Secret    string
	Label     string
	Type      string
	Counter   int
	Issuer    string
	Algorithm string
	Digits    int
	Period    int
	Encoding  string
}

func GenerateSecret(opt *Options) (*GeneratedSecret, error) {
	key := GenerateSecretASCII(opt.Length, opt.Symbols)
	secretKey := GeneratedSecret{
		Ascii:  key,
		Hex:    hex.EncodeToString([]byte(key)),
		Base32: strings.ReplaceAll(base32.StdEncoding.EncodeToString([]byte(key)), "=", ""),
	}

	if opt.QrCodes {
		secretKey.QrCodeAscii = "https://chart.googleapis.com/chart?chs=166x166&chld=L|0&cht=qr&chl=" + url.QueryEscape(secretKey.Ascii)
		secretKey.QrCodeHex = "https://chart.googleapis.com/chart?chs=166x166&chld=L|0&cht=qr&chl=" + url.QueryEscape(secretKey.Hex)
		secretKey.QrCodeBase32 = "https://chart.googleapis.com/chart?chs=166x166&chld=L|0&cht=qr&chl=" + url.QueryEscape(secretKey.Base32)
	}

	if opt.OtpauthURL {
		secretKey.OtpAuthQR, _ = GenerateOtpauthURL(&OtpauthOptions{
			Secret:    secretKey.Base32,
			Label:     opt.Name,
			Issuer:    opt.Issuer,
			Counter:   0,
			Algorithm: "sha1",
		})
	}
	return &secretKey, nil
}

func GenerateSecretASCII(length int, symbols bool) string {
	bytes := make([]byte, length)
	rand.Read(bytes)
	set := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXTZabcdefghiklmnopqrstuvwxyz"
	if symbols {
		set += "!@#$%^&*()<>?/[]{},.:;"
	}
	output := ""
	for _, b := range bytes {
		output += strings.Split(set, "")[int(math.Floor(float64(b)/255.0*float64((len(set)-1))))]
	}
	return output
}

func GenerateOtpauthURL(opt *OtpauthOptions) (string, error) {
	otpType := "totp"

	if opt.Type != "" {
		otpType = opt.Type
	}

	if opt.Type == "hotp" && opt.Counter == 0 {
		log.Fatal("Missing counter value for HOTP'")
	}

	u := url.URL{
		Scheme: "otpauth",
		Host:   otpType,
		Path:   url.QueryEscape(opt.Label),
	}

	q := u.Query()
	q.Set("secret", opt.Secret)
	q.Set("issuer", opt.Issuer)
	u.RawQuery = q.Encode()
	u.Path = opt.Label
	return u.String(), nil
}

func VerifyOTP(secret, code string) (bool, error) {
	otpconf := &g.OTPConfig{
		Secret:      secret,
		WindowSize:  3,
		HotpCounter: 0,
	}
	return otpconf.Authenticate(code)
}
