package jwts

import "testing"

func TestParseToken(t *testing.T) {
	tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjM0NTQ0NjMsInRva2VuIjoiMTAwNSJ9.JhFiitSHJqelNmQ7DO7hi9hoLTISlg4IqDfJLM_W4Uo"
	ParseToken(tokenString, "msproject")
}
