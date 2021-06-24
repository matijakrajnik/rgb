package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"rgb/internal/conf"
	"rgb/internal/store"
	"strconv"
	"time"

	"github.com/cristalhq/jwt/v3"
	"github.com/rs/zerolog/log"
)

var (
	jwtSigner   jwt.Signer
	jwtVerifier jwt.Verifier
)

func jwtSetup(conf conf.Config) {
	var err error
	key := []byte(conf.JwtSecret)

	jwtSigner, err = jwt.NewSignerHS(jwt.HS256, key)
	if err != nil {
		log.Panic().Err(err).Msg("Error creating JWT signer")
	}

	jwtVerifier, err = jwt.NewVerifierHS(jwt.HS256, key)
	if err != nil {
		log.Panic().Err(err).Msg("Error creating JWT verifier")
	}
}

func generateJWT(user *store.User) string {
	claims := &jwt.RegisteredClaims{
		ID:        fmt.Sprint(user.ID),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
	}
	builder := jwt.NewBuilder(jwtSigner)
	token, err := builder.Build(claims)
	if err != nil {
		log.Panic().Err(err).Msg("Error building JWT")
	}
	return token.String()
}

func verifyJWT(tokenStr string) (int, error) {
	token, err := jwt.Parse([]byte(tokenStr))
	if err != nil {
		log.Error().Err(err).Str("tokenStr", tokenStr).Msg("Error parsing JWT")
		return 0, err
	}

	if err := jwtVerifier.Verify(token.Payload(), token.Signature()); err != nil {
		log.Error().Err(err).Msg("Error verifying token")
		return 0, err
	}

	var claims jwt.StandardClaims
	if err := json.Unmarshal(token.RawClaims(), &claims); err != nil {
		log.Error().Err(err).Msg("Error unmarshalling JWT claims")
		return 0, err
	}

	if notExpired := claims.IsValidAt(time.Now()); !notExpired {
		return 0, errors.New("Token expired.")
	}

	id, err := strconv.Atoi(claims.ID)
	if err != nil {
		log.Error().Err(err).Str("claims.ID", claims.ID).Msg("Error converting claims ID to number")
		return 0, errors.New("ID in token is not valid")
	}
	return id, err
}
