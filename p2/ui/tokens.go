package ui

type tType int

const (
	Keyword tType = iota
)

type Token struct {
	index   int
	raw     []rune
	Type    tType
	lenght  int
	currPos int
}

func NewToken(i int, tip tType) *Token {
	return &Token{
		index: i,
		Type:  tip,
	}
}

type Line struct {
	index     int
	tokens    []*Token
	currToken *Token
}

func (t *Token) GetString() string {
	str := ""
	for v := range t.raw {
		str = str + string(t.raw[v])
	}
	return str
}
func (t *Token) SetType(tip tType) {
	t.Type = tip
}

func (t *Token) GetType() tType {
	return t.Type
}
func (t *Token) GetPos() int {
	return t.currPos
}

func (t *Token) Forward() int {
	if t.currPos == t.lenght-1 {
		return t.currPos
	} else {
		t.currPos++
		return t.currPos
	}
}
