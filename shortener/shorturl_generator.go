package shortener

import (
	"crypto/sha256"
	"fmt"
	"math/big"
	"os"

	"github.com/itchyny/base58-go"
)

func sha2560f(input string) []byte {
	algorithm := sha256.New()
	algorithm.Write([]byte(input))
	return algorithm.Sum(nil)
}

func base58Encoded(bytes []byte) string {
	encoding := base58.BitcoinEncoding
	encoded, err := encoding.Encode(bytes)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	return string(encoded)
}

func GenerateShortLink(originalUrl string, userId string) string {
	urlHashByte := sha2560f(originalUrl + userId)
	generatedNumber := new(big.Int).SetBytes(urlHashByte).Uint64()
	urlEncoded := base58Encoded([]byte(fmt.Sprintf("%d", generatedNumber)))
	return urlEncoded[:8]
}
