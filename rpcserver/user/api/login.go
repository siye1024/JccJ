package api

import (
	"context"
	"crypto/subtle"
	"dousheng/db"
	"dousheng/rpcserver/kitex_gen/user"
	"encoding/base64"
	"fmt"
	"github.com/cloudwego/kitex/pkg/kerrors"
	"golang.org/x/crypto/argon2"
	"strings"
)

type Argon2Params struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

type CheckUserOp struct {
	ctx context.Context
}

// NewCheckUserOp new CheckUser Operation
func NewCheckUserOp(ctx context.Context) *CheckUserOp {
	return &CheckUserOp{
		ctx: ctx,
	}
}

// CheckUser check user info, such as username, password and so on
func (s *CheckUserOp) CheckUser(req *user.DouyinUserRegisterRequest) (int64, error) {
	username := req.Username
	users, err := db.QueryUser(s.ctx, username)
	if err != nil {
		return 0, err
	}
	if len(users) == 0 {
		err := kerrors.NewBizStatusError(10007, "Invalid Username")
		return 0, err
	}

	userLogin := users[0]
	pwdMatch, err := cmpPasswordAndHash(req.Password, userLogin.Password)
	if err != nil {
		return 0, err
	}

	if !pwdMatch {
		err := kerrors.NewBizStatusError(10008, "Invalid Username of Password")
		return 0, err
	}

	return int64(userLogin.ID), nil
}

// cmpPasswordAndHash compares the password and hash of the input password
func cmpPasswordAndHash(password, encodedHash string) (match bool, err error) {
	// get the parameters, salt and decoded the password hash
	arg2Params, salt, hash, err := decodeHash(encodedHash)
	if err != nil {
		return false, err
	}

	// compute the hash of the input password with the same parameters
	inputHash := argon2.IDKey([]byte(password), salt, arg2Params.Iterations, arg2Params.Memory, arg2Params.Parallelism, arg2Params.KeyLength)

	// Check the contents of the hashed passwords are indentical
	if subtle.ConstantTimeCompare(hash, inputHash) == 1 {
		return true, nil
	}
	return false, nil
}

func decodeHash(encodedHash string) (arg2Params *Argon2Params, salt, hash []byte, err error) {
	values := strings.Split(encodedHash, "$")
	if len(values) != 6 {
		err := kerrors.NewBizStatusError(10010, "Encoded Hash isn't in the correct format")
		return nil, nil, nil, err
	}

	var version int
	_, err = fmt.Sscanf(values[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, err
	}
	if version != argon2.Version {
		err := kerrors.NewBizStatusError(10010, "Encoded Hash isn't in the correct version")
		return nil, nil, nil, err
	}

	arg2Params = &Argon2Params{}
	if _, err := fmt.Sscanf(values[3], "m=%d,t=%d,p=%d", &arg2Params.Memory, &arg2Params.Iterations, &arg2Params.Parallelism); err != nil {
		return nil, nil, nil, err
	}

	salt, err = base64.RawStdEncoding.Strict().DecodeString(values[4])
	if err != nil {
		return nil, nil, nil, err
	}
	arg2Params.SaltLength = uint32(len(salt))

	hash, err = base64.RawStdEncoding.Strict().DecodeString(values[5])
	if err != nil {
		return nil, nil, nil, err
	}
	arg2Params.KeyLength = uint32(len(hash))

	return arg2Params, salt, hash, nil
}
