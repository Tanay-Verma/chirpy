package main

import (
	"fmt"
	"testing"
)

var cases = []struct {
	Og_msg       string
	Filtered_msg string
}{
	{
		Og_msg:       "This is a kerfuffle opinion I need to share with the world",
		Filtered_msg: "This is a **** opinion I need to share with the world",
	},
	{
		Og_msg:       "I had something interesting for breakfast",
		Filtered_msg: "I had something interesting for breakfast",
	},
}

func TestProfanityFilter(t *testing.T) {
	restrictedWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	for _, test := range cases {
		t.Run(
			fmt.Sprintf(
				"%q \nget converted to \n%q\n",
				test.Og_msg,
				test.Filtered_msg,
			),
			func(t *testing.T) {
				got := getCleanedBody(test.Og_msg, restrictedWords)
				if got != test.Filtered_msg {
					t.Errorf("got: \n%q want: \n%q\n", got, test.Filtered_msg)
				}
			},
		)
	}
}
