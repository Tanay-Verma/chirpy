package main

import "strings"

const replacementWord = "****"

func getCleanedBody(message string, restrictedWords map[string]struct{}) string {
	msgTokens := strings.Split(message, " ")
	for i, word := range msgTokens {
		if _, ok := restrictedWords[strings.ToLower(word)]; ok {
			msgTokens[i] = replacementWord
		}
	}

	return strings.Join(msgTokens, " ")
}
