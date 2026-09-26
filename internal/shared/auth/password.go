package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/callmhejerry/sms/internal/shared/apierror"
	"golang.org/x/crypto/argon2"
)

const (
	argonTime    = 1
	argonMemory  = 64 * 1024
	argonThreads = 4
	argonKeyLen  = 32
	saltLength   = 16
)

func HashPassword(password string) (string, error) {
	salt := make([]byte, saltLength)

	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, argonTime, argonMemory, argonThreads, argonKeyLen)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// Format: argon2id$v=19$m=65536,t=1,p=4$salt$hash

	encoded := fmt.Sprintf("argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, argonMemory, argonTime, argonThreads, b64Salt, b64Hash)

	return encoded, nil
}

func CheckPassword(password, encodedHash string) (bool, *apierror.AppError) {
	parts := strings.Split(encodedHash, "$")

	if len(parts) != 5 {
		return false, apierror.Internal(nil, "Invalid hash format")
	}

	var memory, time uint32
	var threads uint8

	// parts[0] = "argon2id"
	// parts[1] = "v=19"
	// parts[2] = "m=65536,t=1,p=4"
	// parts[3] = salt
	// parts[4] = hash

	if _, err := fmt.Sscanf(parts[2], "m=%d,t=%d,p=%d", &memory, &time, &threads); err != nil {
		return false, apierror.Internal(err, "Something went wrong")
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[3])

	if err != nil {
		return false, apierror.Internal(err, "Something went wrong")
	}

	decodedhash, err := base64.RawStdEncoding.DecodeString(parts[4])

	if err != nil {
		return false, nil
	}

	hash := argon2.IDKey([]byte(password), salt, time, memory, threads, uint32(len(decodedhash)))

	if subtle.ConstantTimeCompare(hash, decodedhash) == 1 {
		return true, nil
	}
	return false, nil
}
