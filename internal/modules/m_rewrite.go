package modules

import (
	"bytes"
	"fmt"
	"github.com/buger/jsonparser"
	str "github.com/profzone/eden-framework/pkg/strings"
	"github.com/valyala/fasthttp"
)

const (
	dollar  byte = '$'
	lParent byte = '('
	rParent byte = ')'
	dot     byte = '.'
)

var (
	context = []byte("context")
	origin  = []byte("origin")
	path    = []byte("path")
	query   = []byte("query")
	cookie  = []byte("cookie")
	header  = []byte("header")
	body    = []byte("body")
)

type rewriteExpr struct {
	origin []byte
	parsed []byte
	tokens []token
	buffer *bytes.Buffer
	err    error
}

func newRewriteExpr(pattern string) rewriteExpr {
	r := rewriteExpr{
		origin: []byte(pattern),
		tokens: make([]token, 0),
		buffer: bytes.NewBuffer([]byte{}),
	}
	r.scan()
	return r
}

func (r *rewriteExpr) scan() {
	var prevToken byte
	var tokenStart = -1
	for index, c := range r.origin {
		switch c {
		case dollar:
			if prevToken != rParent {
				// constToken
				r.tokens = append(r.tokens, &constToken{value: r.origin[tokenStart+1 : index]})
				tokenStart = index
			}
		case lParent:
			if prevToken != dollar {
				r.err = fmt.Errorf("[Col %d]syntax error: \"(\" must followed by \"$\"", index)
				return
			}
			tokenStart = index
		case rParent:
			if tokenStart == 0 {
				r.err = fmt.Errorf("[Col %d]syntax error: can't find previours \"(\"", index)
				return
			}
			r.tokens = append(r.tokens, newParamToken(r.origin[tokenStart+1:index]))
			tokenStart = index
		}

		if index == len(r.origin)-1 && tokenStart != index {
			r.tokens = append(r.tokens, &constToken{value: r.origin[tokenStart+1:]})
		}
		prevToken = c
	}
}

func (r *rewriteExpr) apply(req *fasthttp.Request) error {
	for _, token := range r.tokens {
		err := token.apply(r.buffer, req)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *rewriteExpr) uri() string {
	return r.buffer.String()
}

func (r *rewriteExpr) Error() error {
	return r.err
}

type token interface {
	apply(buf *bytes.Buffer, req *fasthttp.Request) error
}

func newParamToken(value []byte) token {
	sep := bytes.Split(value, []byte{dot})
	if bytes.Compare(sep[0], origin) == 0 {
		return newOriginToken(sep[1:])
	} else if bytes.Compare(sep[0], context) == 0 {

	}
	return nil
}

func newOriginToken(value [][]byte) token {
	if bytes.Compare(value[0], query) == 0 {
		return &originQueryToken{key: value[1]}
	} else if bytes.Compare(value[0], body) == 0 {
		return &originBodyToken{path: str.BytesToStrings(value[1:])}
	} else if bytes.Compare(value[0], cookie) == 0 {

	} else if bytes.Compare(value[0], header) == 0 {

	} else if bytes.Compare(value[0], path) == 0 {

	}

	return nil
}

type constToken struct {
	value []byte
}

func (r *constToken) apply(buf *bytes.Buffer, req *fasthttp.Request) error {
	_, err := buf.Write(r.value)
	return err
}

type originQueryToken struct {
	key []byte
}

func (r *originQueryToken) apply(buf *bytes.Buffer, req *fasthttp.Request) error {
	_, err := buf.Write(req.URI().QueryArgs().Peek(string(r.key)))
	return err
}

type originBodyToken struct {
	path []string
}

func (r *originBodyToken) apply(buf *bytes.Buffer, req *fasthttp.Request) error {
	value, _, _, err := jsonparser.Get(req.Body(), r.path...)
	if err != nil {
		return err
	}

	_, err = buf.Write(value)
	return err
}
