package document

type Tokens []Token

func (ts Tokens) Index(i int) *Token {
	for _, t := range ts {
		if t.Index == i {
			return &t
		}
	}
	return nil
}

type Token struct {
	Word         string `json:"word"`
	OriginalText string `json:"originalText"`
	Lemma        string `json:"lemma"`
	Pos          string `json:"pos"`
	Ner          string `json:"ner"`

	CharacterOffsetBegin int `json:"characterOffsetBegin"`
	CharacterOffsetEnd   int `json:"characterOffsetEnd"`

	Before string `json:"before"`
	Index  int    `json:"index"`
	After  string `json:"after"`
}
