package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"math"
	"math/big"
	"regexp"
	"strconv"
	"strings"

	"github.com/tyler-smith/go-bip39"
)

// https://en.bitcoin.it/wiki/Seed_phrase

// Number of sides on the dice.
var M = 6
var M_REGEXP = "[^1-6]+"

// Bits of entropy we want.
var ENT = 256

var ENTROPY_PER_ROLL = math.Log2(float64(M))
var ROLLS = int(math.Ceil(float64(ENT) / ENTROPY_PER_ROLL))

func rollsToSeed(s string, inc int) (*big.Int, error) {
	x := big.NewInt(0)
	m := big.NewInt(int64(M))
	for _, c := range s {
		i, err := strconv.Atoi(string(c))
		if err != nil {
			return nil, err
		}

		i -= inc

		if i < 0 || i >= M {
			return nil, fmt.Errorf("invalid input %c", c)
		}

		x.Mul(x, m)
		x.Add(x, big.NewInt(int64(i)))
	}

	// throw away bits higher than ENT
	max := big.NewInt(0).Exp(big.NewInt(2), big.NewInt(int64(ENT)), nil)
	max.Sub(max, big.NewInt(1))
	x.And(x, max)

	return x, nil
}

func main() {
	reg, err := regexp.Compile(M_REGEXP)
	if err != nil {
		log.Fatal(err)
	}

	s := ""
	for len(s) < ROLLS {
		fmt.Printf("Roll the dice %3d times, press ENTER to see remaining attempts. ", ROLLS-len(s))
		fmt.Printf("Enter string of 1-6: ")

		var input string
		_, err := fmt.Scanf("%s", &input)
		if err != nil {
			log.Fatal(err)
		}

		input = reg.ReplaceAllString(input, "")
		s += input
	}

	seed, err := rollsToSeed(s, 1)
	if err != nil {
		log.Fatal(err)
	}

	entropy := make([]byte, ENT/8)
	for i, b := range seed.Bytes() {
		entropy[i] = b
	}

	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Dice:", s)
	fmt.Println("Entropy:", hex.EncodeToString(entropy))
	fmt.Println("Mnemonic:")

	for i, w := range strings.Split(mnemonic, " ") {
		fmt.Printf("%2d. %s\n", i+1, w)
	}
}
