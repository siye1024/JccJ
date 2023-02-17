package api

import (
	"context"
	"dousheng/db"
	"dousheng/rpcserver/user/kitex_gen/user"
	"encoding/base64"
	"fmt"
	"github.com/cloudwego/kitex/pkg/kerrors"
	"golang.org/x/crypto/argon2"
	"math/rand"
)

type CreateUserOp struct {
	ctx context.Context
}

// NewCreateUserOp new CreateUserOp
func NewCreateUserOp(ctx context.Context) *CreateUserOp {
	return &CreateUserOp{ctx: ctx}
}

// CreateUser create user info
func (s *CreateUserOp) CreateUser(req *user.DouyinUserRegisterRequest, arg2Params *Argon2Params) error {
	users, err := db.QueryUser(s.ctx, req.Username)
	if err != nil {
		return err
	}
	if len(users) != 0 {
		err := kerrors.NewBizStatusError(10002, "Username Already Exist")
		return err
	}

	// encode the plaintext password and save it to the db
	password, err := encodedPassword(req.Password, arg2Params)
	if err != nil {
		return err
	}

	return db.CreateUser(s.ctx, []*db.User{{
		UserName: req.Username,
		Password: password,
	}})
}

// encode the plaintext password to the hash password
func encodedPassword(password string, arg2Params *Argon2Params) (encodedHash string, err error) {
	salt, err := generateRandomBytes(arg2Params.SaltLength)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, arg2Params.Iterations, arg2Params.Memory, arg2Params.Parallelism, arg2Params.KeyLength)
	// Baset4 encode the salt and hashed password
	base64Salt := base64.RawStdEncoding.EncodeToString(salt)
	base64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// return hash according to the standard encoded hash representation
	encodedHash = fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, arg2Params.Memory, arg2Params.Iterations, arg2Params.Parallelism, base64Salt, base64Hash)

	return encodedHash, nil
}

// generateRandomBytes returns a random bytes.
func generateRandomBytes(saltLength uint32) ([]byte, error) {
	buf := make([]byte, saltLength)
	_, err := rand.Read(buf)
	if err != nil {
		return nil, err
	}

	return buf, nil
}
