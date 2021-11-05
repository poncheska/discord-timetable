package utils

import (
	"log"
	"testing"
)

func TestGenerateRandomId(t *testing.T) {
	s := GenerateRandomId(6)
	log.Println(s)
	s = GenerateRandomId(6)
	log.Println(s)
	s = GenerateRandomId(6)
	log.Println(s)
	s = GenerateRandomId(6)
	log.Println(s)
}

func TestShuffle(t *testing.T) {
	in := map[string]string{
		"1234":"1Bob",
		"2345":"2Bib",
		"3456":"3Lol",
		"4567":"4Kek",
	}
	log.Println(Shuffle(in))
	log.Println(Shuffle(in))
	log.Println(Shuffle(in))
	log.Println(Shuffle(in))
	log.Println(Shuffle(in))
}