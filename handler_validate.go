package main

import (
	"encoding/json"
	"net/http"
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

func HandlerValidateChirp(w http.ResponseWriter, r *http.Request){
	type parameters struct{
		Body string `json:"body"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil{
		respondWithError(w, 400, "Something went wrong")
	}
	if len(params.Body) > 140{
		respondWithError(w, 400, "Chirp is too long")
	}else{
		type respParameters struct{
			CleanedBody string `json:"cleaned_body"`
		}
		respondWithJson(w, 200, respParameters{
			CleanedBody: cleanProfanity(params.Body),},
		)
	}
}
