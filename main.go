package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/FactomProject/btcutil/base58"
)

var alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
var assetNames = []string{"PNT", "USD", "EUR", "JPY", "GBP", "CAD", "CHF", "INR", "SGD", "CNY", "HKD", "XAU", "XAG", "XPD", "XPT", "XBT", "ETH", "LTC", "XBC", "FCT"}

// - pre: a string of Base58 characters that the output should start with
// - payload: the amount of bytes in the payload
func findFirst(pre string, payload int) []byte {
	for i := 0; i < len(pre); i++ {
		if !strings.ContainsAny(pre[i:i+1], alphabet) {
			fmt.Println(pre, "contains non Base58 characters:", pre[i:i+1])
			return nil
		}
	}

	min, _ := hex.DecodeString(strings.Repeat("00", payload)) // left and right boundaries
	max, _ := hex.DecodeString(strings.Repeat("FF", payload))

	prefix := make([]byte, len(pre))
	for i := 0; i < len(pre); i++ {
		for j := 0; j < 256; j++ {
			prefix[i] = byte(j)
			lh := base58.Encode(append((prefix), min...))

			if matchall(lh, pre, i) {
				prefix = cos(prefix, i) // if we find the pre-prefix, we already surpassed the minimum
				break
			}

			// nothing found in range, carry over and reset loop
			if j == 255 && i > 0 {
				prefix = co(prefix, i)
				j = 0
			}
		}
	}
	prefix = co(prefix, len(prefix)-1) // readd the last subtract
	fmt.Printf("%s(%d) = %x, %s to %s\n", pre, payload, prefix, base58.Encode(append(prefix, min...)), base58.Encode(append(prefix, max...)))
	//fmt.Println(base58.Encode(append(prefix, min...)))
	return prefix
}

func main() {
	fmt.Println("Human Readable Unique Addresses")
	for _, p := range assetNames {
		findFirst("p"+p, 32+4) // pubkey + checksum
	}

	fmt.Println()
	fmt.Println("Addresses with currency encoded inside payload")

	prefix := findFirst("PEG", 32+4+5)
	adr := make([]byte, 32) // empty RCD

	for _, p := range assetNames {
		fmt.Println(pegAddr(fmt.Sprintf("p%s_", p), prefix, adr))
	}
}

func pegAddr(name string, prefix, adr []byte) string {
	adr = append([]byte(name), adr...)
	adr = append(prefix, adr...)
	hash := sha256.Sum256(adr)
	csum := sha256.Sum256(hash[:])
	return name + base58.Encode(append(adr, csum[:4]...))
}
