<h1>2fa-auth</h1>
<p>
</p>

> 2fa-auth is a one-time passcode (time-based) generator, ideal for use in two-factor authentication, that supports Google Authenticator and other two-factor devices.

## Install

```sh
go get github.com/tsirysndr/2fa-auth
```

## Usage

Let's say you have a user that wants to enable two-factor authentication, and you intend to do two-factor authentication using an app like Google Authenticator, Duo Security, Authy, etc. This is a three-step process:

1. Generate a secret
2. Show a QR code for the user to scan in
3. Authenticate the token for the first time

### Generating a key

This will generate a secret key of length 16, which will be the secret key for the user.

```go
import (
  fmt
  auth "github.com/tsirysndr/2fa-auth"
)

...

g, _ := auth.GenerateSecret(&auth.Options{
  Length:     16,
  Symbols:    true,
  OtpauthURL: true,
  Name:       name,
  Issuer:     issuer,
})
secret, _ := json.Marshal(g)
fmt.Println(string(secret))
```

### Verifying the token

After the user scans the QR code, ask the user to enter in the token that they see in their app. Then, verify it against the secret.

```go
import (
  fmt
  auth "github.com/tsirysndr/2fa-auth"
)

...

fmt.Println(auth.VerifyOTP(secret, code))
```

## Author

👤 **Tsiry Sandratraina**

* Website: https://tsiry-sandratraina.netlify.com
* Github: [@tsirysndr](https://github.com/tsirysndr)

## Show your support

Give a ⭐️ if this project helped you!
