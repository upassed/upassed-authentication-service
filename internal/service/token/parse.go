package token

import "github.com/golang-jwt/jwt"

func (service *serviceImpl) parseToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrInvalidKeyType
		}

		return []byte(service.cfg.Jwt.Secret), nil
	})
}
