package document

import (
	"fmt"
	"strings"
)

type Sentences []Sentence

type Sentence struct {
	RawParse                     string       `json:"parse"`
	Dependencies                 []Dependency `json:"basicDependencies"`
	EnhancedPlusPlusDependencies []Dependency `json:"enhancedPlusPlusDependencies"`
	Tokens                       Tokens       `json:"tokens"`
}

type Parse struct {
	Index int
	Text  string
	Pos   string

	Token *Token

	Parent   *Parse
	Children []*Parse
}

func (s *Sentence) Parse() (*Parse, error) {
	tokens := make(map[int]*Token, len(s.Tokens)+1)
	for _, t := range s.Tokens {
		tokens[t.Index] = &t
	}

	replacer := strings.NewReplacer("\n", "", "\t", "")
	rawParse := replacer.Replace(strings.Trim(s.RawParse, " \r\n\t"))

	decoder := parseDecoder{
		rawParse:   rawParse,
		parseIndex: 0,
		tokenIndex: s.Tokens[0].Index,
		tokens:     tokens,
	}

	parse, _, err := decoder.Decode(0)
	return parse, err
}

type parseDecoder struct {
	rawParse   string
	tokens     map[int]*Token
	parseIndex int
	tokenIndex int
}

func (d *parseDecoder) Decode(start int) (parse *Parse, end int, err error) {
	if d.rawParse[start] != '(' {
		return nil, -1, fmt.Errorf("invalid start position. [%s]", d.rawParse)
	}

	p := start + strings.Index(d.rawParse[start:], " ")
	parse = &Parse{
		Index: d.parseIndex,
		Pos:   d.rawParse[start+1 : p],
	}
	d.parseIndex++

	for current := p + 1; current < len(d.rawParse); current++ {
		if d.rawParse[current] == '(' {
			child, end, err := d.Decode(current)
			if err != nil {
				return nil, -1, err
			}

			child.Parent = parse
			parse.Children = append(parse.Children, child)
			current = end

		} else if d.rawParse[current] == ')' {
			if len(parse.Children) == 0 {
				parse.Token = d.tokens[d.tokenIndex]
				parse.Text = parse.Token.OriginalText
				d.tokenIndex++

				if parse.Token == nil {
					return nil, -1, fmt.Errorf("Couldn't detect the token")
				}

			} else {
				// Skip all end of closing parenthesis
				//  (NNP noun))))))))))))
				//               ^
			}
			return parse, current, nil
		}
	}
	return nil, -1, fmt.Errorf("Leached the end of parse")
}
