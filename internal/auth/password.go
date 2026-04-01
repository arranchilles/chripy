package auth

import "github.com/alexedwards/argon2id"

func HashPassword(password string) (string, error) {
	hashArgs := argon2id.Params{
		Iterations:  1,
		Parallelism: 1,
		Memory:      1024,
		SaltLength:  16,
		KeyLength:   16,
	}

	hashedPassword, err := argon2id.CreateHash(password, &hashArgs)
	if err != nil {
		return "", err
	}
	return hashedPassword, nil
}

func CheckPasswordHash(password string, hash string) (bool, error) {
	isPassword, err := argon2id.ComparePasswordAndHash(password, hash)

	if err != nil {
		return false, err
	}

	return isPassword, nil
}
