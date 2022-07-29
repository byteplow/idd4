package kratos

import (
	"context"
	"crypto/rand"
	"math/big"

	"github.com/byteplow/idd4/internal/container"
	kratos "github.com/ory/kratos-client-go"
)

const (
	javascriptISOString = "2006-01-02T15:04:05.999Z"
)

func CreateInvite() (string, error) {
	identity, err := createIdentity()
	if err != nil {
		return "", err
	}

	link, err := createRecoveryLink(identity.Id)
	if err != nil {
		return "", err
	}

	return link, nil
}

func createRecoveryLink(identityId string) (string, error) {
	expiresIn := "256h"

	link, _, err := container.KratosAdminClient.AdminCreateSelfServiceRecoveryLink(context.Background()).AdminCreateSelfServiceRecoveryLinkBody(kratos.AdminCreateSelfServiceRecoveryLinkBody{
		ExpiresIn:  &expiresIn,
		IdentityId: identityId,
	}).Execute()

	if err != nil {
		return "", nil
	}

	return link.RecoveryLink, nil
}

func createIdentity() (*kratos.Identity, error) {
	id, err := generateUsername()
	if err != nil {
		return nil, err
	}

	identity, _, err := container.KratosAdminClient.AdminCreateIdentity(context.Background()).AdminCreateIdentityBody(kratos.AdminCreateIdentityBody{
		SchemaId: "default",
		Traits: map[string]interface{}{
			"username": id,
		},
	}).Execute()
	if err != nil {
		return nil, err
	}

	return identity, nil
}

func generateRandomString(n int) (string, error) {
	const letters = "0123456789abcdefghijklmnopqrstuvwxyz"
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		ret[i] = letters[num.Int64()]
	}

	return string(ret), nil
}

func generateUsername() (string, error) {
	str, err := generateRandomString(12)
	if err != nil {
		return "nil", err
	}

	return str[:6] + "-" + str[6:], nil
}