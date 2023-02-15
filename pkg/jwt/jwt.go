package jwt

import (
	"github.com/cloudwego/kitex/pkg/kerrors"
	"github.com/golang-jwt/jwt"
)

// JWT signing Key
type JWT struct {
	SigningKey []byte
}

// private claims, share information between parties that agree on using them
// CustomClaims Structured version of Claims Section, as referenced at https://tools.ietf.org/html/rfc7519#section-4.1 See examples for how to use this with your own claim types
type CustomClaims struct {
	Id   int64
	Time int64
	jwt.StandardClaims
}

func NewJWT(SigningKey []byte) *JWT {
	return &JWT{
		SigningKey,
	}
}

// CreateToken creates a new token
func (j *JWT) CreateToken(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// zap.S().Debugf(token.SigningString())
	return token.SignedString(j.SigningKey)

}

// ParseToken parses the token.
func (j *JWT) ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, kerrors.NewBizStatusError(10003, "JWT ERROR:That's not even a token")
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, kerrors.NewBizStatusError(10004, "JWT ERROR:Token expired")
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, kerrors.NewBizStatusError(10005, "JWT ERROR:Token is not active yet")
			} else {
				return nil, kerrors.NewBizStatusError(10006, "JWT ERROR:Couldn't handle this token")
			}

		}
	}
	// verify the token claims
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, kerrors.NewBizStatusError(10006, "JWT ERROR:Couldn't handle this token")
}
