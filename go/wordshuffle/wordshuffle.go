package wordshuffle

import (
	"net/http"
	"fmt"
)

func HandlePata(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	word1 := r.FormValue("word1")
	word2 := r.FormValue("word2")
	newWord := makeNewWords([]rune(word1), []rune(word2))
	fmt.Fprintf(w, "%s", newWord)
	http.ServeFile(w, r, "view/layout.html")
}

func makeNewWords(word1 []rune, word2 []rune) string {
	var newWord []rune
	if len(word1) > len(word2) {
		for i := 0; i < len(word2); i++ {
			newWord = append(append(newWord, word1[i:i+1]...), word2[i:i+1]...)
		}
		newWord = append(newWord, word1[len(word2):]...)
	} else {
		for i := 0; i < len(word1); i++ {
			newWord = append(append(newWord, word1[i:i+1]...), word2[i:i+1]...)
		}
		newWord = append(newWord, word2[len(word1):]...)
	}
	return string(newWord)
}
