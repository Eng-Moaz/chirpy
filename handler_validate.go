package main

import (
	"strings"
)

func BadWordsSet() map[string]struct{}{
	BadWords := make(map[string]struct{})
	BadWords["kerfuffle"] = struct{}{}
	BadWords["sharbert"] = struct{}{}
	BadWords["fornax"] = struct{}{}
	return BadWords
}

func cleanProfanity(chirp string) string{
	BadWords := BadWordsSet()
	finalString := make([]string, 0)
	cleaned := strings.SplitSeq(chirp, " ")
	for word := range cleaned{
		if _, exists := BadWords[strings.ToLower(word)]; exists{
			finalString = append(finalString, "****")
		}else{
			finalString = append(finalString, word)
		}	
	}
	return strings.Join(finalString, " ")
}
