package jwt

import (
	"fmt"
	"log"

	"github.com/nats-io/jwt/v2"
	"github.com/nats-io/nkeys"
)

func GenerateToken() {
	log.SetFlags(0)

	operatorKP, _ := nkeys.CreateOperator()

	operatorPub, _ := operatorKP.PublicKey()
	fmt.Printf("operator pubkey: %s\n", operatorPub)

	operatorSeed, _ := operatorKP.Seed()
	fmt.Printf("operator seed: %s\n\n", string(operatorSeed))

	accountKP, _ := nkeys.CreateAccount()

	accountPub, _ := accountKP.PublicKey()
	fmt.Printf("account pubkey: %s\n", accountPub)

	accountSeed, _ := accountKP.Seed()
	fmt.Printf("account seed: %s\n", string(accountSeed))

	accountClaims := jwt.NewAccountClaims(accountPub)
	accountClaims.Name = "my-account"

	accountClaims.Limits.JetStreamLimits.DiskStorage = -1
	accountClaims.Limits.JetStreamLimits.MemoryStorage = -1

	fmt.Printf("account claims: %s\n", accountClaims)

	accountJWT, _ := accountClaims.Encode(operatorKP)
	fmt.Printf("account jwt: %s\n\n", accountJWT)

	userKP, _ := nkeys.CreateUser()

	userPub, _ := userKP.PublicKey()
	fmt.Printf("user pubkey: %s\n", userPub)

	userSeed, _ := userKP.Seed()
	fmt.Printf("user seed: %s\n", string(userSeed))

	userClaims := jwt.NewUserClaims(userPub)
	userClaims.Name = "my-user"

	userClaims.Limits.Data = 1024 * 1024 * 1024

	userClaims.Permissions.Pub.Allow.Add("foo.>", "bar.>")
	userClaims.Permissions.Sub.Allow.Add("_INBOX.>")

	fmt.Printf("userclaims: %s\n", userClaims)

	userJWT, _ := userClaims.Encode(accountKP)
	fmt.Printf("user jwt: %s\n", userJWT)

	creds, _ := jwt.FormatUserConfig(userJWT, userSeed)
	fmt.Printf("creds file: %s\n", creds)
}
