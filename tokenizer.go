package tokenizer

import (
	"bufio"
	"io"
	"os"
	"path"
	"runtime"
	"strings"
	"unicode"
)

type Tokenizer struct {
	dict      *Dict
	dictPath  string
	stopWords map[string]bool
}

func NewTokenizer(stopwordsPath string) *Tokenizer {
	tk := &Tokenizer{}
	var _, filename, _, _ = runtime.Caller(0)
	tk.dict = NewDict(path.Dir(filename) + "data/dictionary.txt")
	tk.loadStopwords(stopwordsPath)
	return tk
}

func (tk *Tokenizer) tokenize(text string) []string {
	var result []string
	startLen := tk.dict.maxLen
	text = strings.Trim(text, " ")
	for len([]rune(text)) > 0 {
		if len([]rune(text)) < startLen {
			startLen = len([]rune(text))
		}
		word := string([]rune(text)[0:startLen])
		for {
			if len([]rune(word)) == 1 {
				break
			}
			_, ok := tk.dict.Records[word]
			if ok || filter(word) {
				break
			}
			word = string([]rune(word)[0 : len([]rune(word))-1])
		}
		if _, found := tk.stopWords[word]; !found {
			result = append(result, word)
		}
		text = string([]rune(text)[len([]rune(word)):])
	}
	return result
}

func filter(text string) bool {
	for _, r := range []rune(text) {
		if unicode.Is(unicode.Scripts["Han"], r) || unicode.IsSpace(r) || unicode.IsControl(r) || unicode.IsPunct(r) {
			return false
		}
	}
	return true
}

func (tk *Tokenizer) loadStopwords(path string) {
	tk.stopWords = map[string]bool{}
	fi, _ := os.Open(path)
	defer fi.Close()
	br := bufio.NewReader(fi)
	for a, _, c := br.ReadLine(); c != io.EOF; a, _, c = br.ReadLine() {
		tk.stopWords[string(a)] = true
	}
	tk.stopWords[" "] = true
}

