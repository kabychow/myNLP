package mytokenizer

import (
	"bufio"
	"io"
	"log"
	"os"
	"strings"
)

type DictRecord struct {
	TF    string
	Token string
	POS   string
}

type Dict struct {
	Records map[string]DictRecord
	maxLen  int
}

func NewDict(dictPath string) *Dict {
	dict := &Dict{Records: make(map[string]DictRecord)}
	fi, err := os.Open(dictPath)
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

		res := strings.Split(string(a), " ")

		var TF, pos string

		token := res[0]

		if len(res) > 1 {
			TF = res[1]
			pos = res[2]
		}

		currLen := len([]rune(token))
		if currLen > dict.maxLen {
			dict.maxLen = currLen
		}

		if len([]rune(token)) >= 2 {
			dict.Records[token] = DictRecord{
				TF:  TF,
				POS: pos,
			}
		}
	}
	return dict
}

