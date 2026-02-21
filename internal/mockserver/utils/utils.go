package utils

import (
	cryptoRand "crypto/rand"
	"encoding/binary"
	"fmt"
	"math/rand"
	"regexp"
)

func ValidateInternalID(id string) bool {
	if len(id) != 18 {
		return false
	}
	matched, err := regexp.MatchString(`^[a-zA-Z0-9]+$`, id)
	if err != nil {
		return false
	}
	return matched
}

func ValidateCampaignID(id string) bool {
	return ValidateInternalID(id)
}

func ValidateUUID(id string) bool {
	uuidRegex := `^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$`
	matched, err := regexp.MatchString(uuidRegex, id)
	if err != nil {
		return false
	}
	return matched
}

func newRand() *rand.Rand {
	var b [8]byte
	_, _ = cryptoRand.Read(b[:]) // crypto-safe seed
	seed := int64(binary.LittleEndian.Uint64(b[:]))
	return rand.New(rand.NewSource(seed))
}

func GenerateRandomCampaignID() string {
	r := newRand()
	return "701PI0000" + fmt.Sprintf("%05d", r.Intn(100000))
}

func NewUUID() string {
	r := newRand()

	u := make([]byte, 16)
	r.Read(u)

	u[6] = (u[6] & 0x0F) | 0x40
	u[8] = (u[8] & 0x3F) | 0x80

	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		u[0:4],
		u[4:6],
		u[6:8],
		u[8:10],
		u[10:16],
	)
}

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateRandomMessageId() string {
	r := newRand()
	result := make([]byte, 64)
	for i := range result {
		result[i] = letters[r.Intn(len(letters))]
	}
	return string(result)
}
