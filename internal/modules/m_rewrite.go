package modules

import (
	"bytes"
	"fmt"
	"github.com/buger/jsonparser"
	str "github.com/profzone/eden-framework/pkg/strings"
	"github.com/valyala/fasthttp"
	"longhorn/proxy/pkg/route"
)

const (
	rewriteTokenDollar  byte = '$'
	rewriteTokenLParent byte = '('
	rewriteTokenRParent byte = ')'
	rewriteTokenDot     byte = '.'
)

var (
	rewriteExprContext = []byte("context")
	rewriteExprOrigin  = []byte("origin")
	rewriteExprPath    = []byte("path")
	rewriteExprQuery   = []byte("query")
	rewriteExprCookie  = []byte("cookie")
	rewriteExprHeader  = []byte("header")
	rewriteExprBody    = []byte("body")
)

type rewriteExpr struct {
	req    *fasthttp.Request
	origin []byte
	params map[string][]byte
	tokens []token
	buffer *bytes.Buffer
	err    error
}

func newRewriteExpr(req *fasthttp.Request, pattern string, params route.Params) rewriteExpr {
	r := rewriteExpr{
		req:    req,
		origin: []byte(pattern),
		params: make(map[string][]byte),
		tokens: make([]token, 0),
		buffer: bytes.NewBuffer([]byte{}),
	}

	for _, p := range params {
		r.params[p.Key] = []byte(p.Value)
	}
	r.scan()
	return r
}

func (r *rewriteExpr) scan() {
	var prevToken byte
	var tokenStart = -1
	for index, c := range r.origin {
		switch c {
		case rewriteTokenDollar:
			if prevToken != rewriteTokenRParent {
				// constToken
				r.tokens = append(r.tokens, &constToken{value: r.origin[tokenStart+1 : index]})
				tokenStart = index
			}
		case rewriteTokenLParent:
			if prevToken != rewriteTokenDollar {
				r.err = fmt.Errorf("[Col %d]syntax error: \"(\" must followed by \"$\"", index)
				return
			}
			tokenStart = index
		case rewriteTokenRParent:
			if tokenStart == 0 {
				r.err = fmt.Errorf("[Col %d]syntax error: can't find previours \"(\"", index)
				return
			}
			token, err := newParamToken(r.origin[tokenStart+1 : index])
			if err != nil {
				r.err = err
				return
			}
			r.tokens = append(r.tokens, token)
			tokenStart = index
		}

		if index == len(r.origin)-1 && tokenStart != index {
			r.tokens = append(r.tokens, &constToken{value: r.origin[tokenStart+1:]})
		}
		prevToken = c
	}
}

func (r *rewriteExpr) apply() error {
	for _, token := range r.tokens {
		err := token.apply(r.buffer, r)
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
	apply(buf *bytes.Buffer, expr *rewriteExpr) error
}

func newParamToken(value []byte) (token, error) {
	sep := bytes.Split(value, []byte{rewriteTokenDot})
	if bytes.Compare(sep[0], rewriteExprOrigin) == 0 {
		return newOriginToken(sep[1:])
	} else if bytes.Compare(sep[0], rewriteExprContext) == 0 {

	}
	return nil, fmt.Errorf("syntax error: not support expression: %s", string(sep[0]))
}

func newOriginToken(value [][]byte) (token, error) {
	if bytes.Compare(value[0], rewriteExprQuery) == 0 {
		return originQueryToken{key: string(value[1])}, nil
	} else if bytes.Compare(value[0], rewriteExprBody) == 0 {
		return originBodyToken{path: str.BytesToStrings(value[1:])}, nil
	} else if bytes.Compare(value[0], rewriteExprCookie) == 0 {
		return originCookieToken{key: string(value[1])}, nil
	} else if bytes.Compare(value[0], rewriteExprHeader) == 0 {
		return originHeaderToken{key: string(value[1])}, nil
	} else if bytes.Compare(value[0], rewriteExprPath) == 0 {
		return originPathToken{key: string(value[1])}, nil
	}

	return nil, fmt.Errorf("syntax error: not support origin expression: %s", string(value[0]))
}

type constToken struct {
	value []byte
}

func (r *constToken) apply(buf *bytes.Buffer, _ *rewriteExpr) error {
	_, err := buf.Write(r.value)
	return err
}

type originQueryToken struct {
	key string
}

func (r originQueryToken) apply(buf *bytes.Buffer, expr *rewriteExpr) error {
	_, err := buf.Write(expr.req.URI().QueryArgs().Peek(r.key))
	return err
}

type originBodyToken struct {
	path []string
}

func (r originBodyToken) apply(buf *bytes.Buffer, expr *rewriteExpr) error {
	value, _, _, err := jsonparser.Get(expr.req.Body(), r.path...)
	if err != nil {
		return err
	}

	_, err = buf.Write(value)
	return err
}

type originCookieToken struct {
	key string
}

func (o originCookieToken) apply(buf *bytes.Buffer, expr *rewriteExpr) error {
	_, err := buf.Write(expr.req.Header.Cookie(o.key))
	return err
}

type originHeaderToken struct {
	key string
}

func (o originHeaderToken) apply(buf *bytes.Buffer, expr *rewriteExpr) error {
	_, err := buf.Write(expr.req.Header.Peek(o.key))
	return err
}

type originPathToken struct {
	key string
}

func (o originPathToken) apply(buf *bytes.Buffer, expr *rewriteExpr) error {
	_, err := buf.Write(expr.params[o.key])
	return err
}
