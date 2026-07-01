package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)


func TestID(t *testing.T){
	testID := uuid.New()
	tokenSecret := "testingToken"
	testDuration := time.Minute * 15	
	secretString, err := MakeJWT(testID, tokenSecret, testDuration) 
	if err != nil{
		t.Errorf("Error in making a secretString")
	}
	gotId, err := ValidateJWT(secretString, tokenSecret)
	if err != nil{
		t.Errorf("Error in validating: %v", err)
	}
	if gotId != testID{
		t.Errorf("Validation failed expected: %v, got: %v", testID, gotId)
	}
}
