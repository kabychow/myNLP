package mynlp

import (
	"bufio"
	"github.com/RadhiFadlillah/go-sastrawi"
	"github.com/reiver/go-porterstemmer"
	"io"
	"log"
	"os"
	"path"
	"runtime"
	"strings"
	"unicode"
)

var stopwords = getStopwords()
var dictionaries = getDictionaries()
var enStem = porterstemmer.StemString
var idStem = sastrawi.NewStemmer(sastrawi.DefaultDictionary).Stem
var maxLen int

type Nlp struct {
	Raw     string
	Tk      Tokens
	TkStops Tokens
}

func NLP(text string) *Nlp {
	text = strings.TrimSpace(text)
	result := &Nlp{Raw: text}
	startLen := maxLen
	for len([]rune(text)) > 0 {
		if len([]rune(text)) < startLen {
			startLen = len([]rune(text))
		}
		word := string([]rune(text)[0:startLen])
		for {
			if len([]rune(word)) == 1 {
				break
			}
			_, ok := dictionaries[word]
			if ok || filter(word) {
				break
			}
			word = string([]rune(word)[0 : len([]rune(word))-1])
		}
		if strings.TrimSpace(word) != "" {
			tmp := enStem(idStem(word))
			result.TkStops = append(result.TkStops, tmp)
			if _, found := stopwords[word]; !found {
				result.Tk = append(result.Tk, tmp)
			}
		}
		text = string([]rune(text)[len([]rune(word)):])
	}
	return result
}

func getStopwords() map[string]bool {
	tmp := map[string]bool{}
	fi, err := os.Open(fromCurrentPath("/data/stopwords.txt"))
	if err != nil {
		log.Fatal(err)
	}
	defer fi.Close()
	br := bufio.NewReader(fi)
	for a, _, c := br.ReadLine(); c != io.EOF; a, _, c = br.ReadLine() {
		tmp[string(a)] = true
	}
	tmp[" "] = true
	return tmp
}

func getDictionaries() map[string]bool {
	dict := map[string]bool{}
	fi, err := os.Open(fromCurrentPath("/data/dictionary.txt"))
	if err != nil {
		log.Fatal(err)
	}
	defer fi.Close()
	br := bufio.NewReader(fi)
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		token := string(a)
		currLen := len([]rune(token))
		if currLen > maxLen {
			maxLen = currLen
		}
		if len([]rune(token)) >= 2 {
			dict[token] = true
		}
	}
	return dict
}

func fromCurrentPath(append string) string {
	var _, filename, _, _ = runtime.Caller(0)
	return path.Dir(filename) + append
}

func filter(text string) bool {
	for _, r := range []rune(text) {
		if unicode.Is(unicode.Scripts["Han"], r) || unicode.IsSpace(r) || unicode.IsControl(r) || unicode.IsPunct(r) {
			return false
		}
	}
	return true
}
