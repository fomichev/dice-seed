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
var BASE = 6
var M_REGEXP = "[^1-6]+"

// Bits of entropy we want.
var ENT = 256

var ENTROPY_PER_ROLL = math.Log2(float64(BASE))

var ROLLS = int(math.Floor(float64(ENT)/ENTROPY_PER_ROLL) + 1)

func rollsToSeed(s string, map6to0 bool) (*big.Int, error) {
	seed := big.NewInt(0)
	base := big.NewInt(int64(BASE))
	for _, c := range s {
		i, err := strconv.Atoi(string(c))
		if err != nil {
			return nil, err
		}

		if map6to0 && i == 6 {
			i = 0
		}

		if i < 0 || i >= BASE {
			return nil, fmt.Errorf("invalid input %c", c)
		}

		// seed = seed * base + i
		seed.Mul(seed, base)
		seed.Add(seed, big.NewInt(int64(i)))
	}

	return seed, nil
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

	seed, err := rollsToSeed(s, true)
	if err != nil {
		log.Fatal(err)
	}

	entropy := make([]byte, ENT/8)
	for i, b := range seed.Bytes() {
		if i >= ENT/8 {
			break
		}
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
