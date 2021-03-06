/*
Copyright 2016 Google Inc. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package parser

import (
	"testing"
)

type lexTest struct {
	name      string
	input     string
	tokens    Tokens
	errString string
}

var (
	tEOF = token{kind: tokenEndOfFile}
)

var lexTests = []lexTest{
	{"empty", "", Tokens{}, ""},
	{"whitespace", "  \t\n\r\r\n", Tokens{}, ""},

	{"brace L", "{", Tokens{{kind: tokenBraceL, data: "{"}}, ""},
	{"brace R", "}", Tokens{{kind: tokenBraceR, data: "}"}}, ""},
	{"bracket L", "[", Tokens{{kind: tokenBracketL, data: "["}}, ""},
	{"bracket R", "]", Tokens{{kind: tokenBracketR, data: "]"}}, ""},
	{"colon", ":", Tokens{{kind: tokenOperator, data: ":"}}, ""},
	{"colon2", "::", Tokens{{kind: tokenOperator, data: "::"}}, ""},
	{"colon3", ":::", Tokens{{kind: tokenOperator, data: ":::"}}, ""},
	{"arrow right", "->", Tokens{{kind: tokenOperator, data: "->"}}, ""},
	{"less than minus", "<-", Tokens{{kind: tokenOperator, data: "<"},
		{kind: tokenOperator, data: "-"}}, ""},
	{"comma", ",", Tokens{{kind: tokenComma, data: ","}}, ""},
	{"dollar", "$", Tokens{{kind: tokenDollar, data: "$"}}, ""},
	{"dot", ".", Tokens{{kind: tokenDot, data: "."}}, ""},
	{"paren L", "(", Tokens{{kind: tokenParenL, data: "("}}, ""},
	{"paren R", ")", Tokens{{kind: tokenParenR, data: ")"}}, ""},
	{"semicolon", ";", Tokens{{kind: tokenSemicolon, data: ";"}}, ""},

	{"not 1", "!", Tokens{{kind: tokenOperator, data: "!"}}, ""},
	{"not 2", "! ", Tokens{{kind: tokenOperator, data: "!"}}, ""},
	{"not equal", "!=", Tokens{{kind: tokenOperator, data: "!="}}, ""},
	{"tilde", "~", Tokens{{kind: tokenOperator, data: "~"}}, ""},
	{"plus", "+", Tokens{{kind: tokenOperator, data: "+"}}, ""},
	{"minus", "-", Tokens{{kind: tokenOperator, data: "-"}}, ""},

	{"number 0", "0", Tokens{{kind: tokenNumber, data: "0"}}, ""},
	{"number 1", "1", Tokens{{kind: tokenNumber, data: "1"}}, ""},
	{"number 1.0", "1.0", Tokens{{kind: tokenNumber, data: "1.0"}}, ""},
	{"number 0.10", "0.10", Tokens{{kind: tokenNumber, data: "0.10"}}, ""},
	{"number 0e100", "0e100", Tokens{{kind: tokenNumber, data: "0e100"}}, ""},
	{"number 1e100", "1e100", Tokens{{kind: tokenNumber, data: "1e100"}}, ""},
	{"number 1.1e100", "1.1e100", Tokens{{kind: tokenNumber, data: "1.1e100"}}, ""},
	{"number 1.1e-100", "1.1e-100", Tokens{{kind: tokenNumber, data: "1.1e-100"}}, ""},
	{"number 1.1e+100", "1.1e+100", Tokens{{kind: tokenNumber, data: "1.1e+100"}}, ""},
	{"number 0100", "0100", Tokens{
		{kind: tokenNumber, data: "0"},
		{kind: tokenNumber, data: "100"},
	}, ""},
	{"number 10+10", "10+10", Tokens{
		{kind: tokenNumber, data: "10"},
		{kind: tokenOperator, data: "+"},
		{kind: tokenNumber, data: "10"},
	}, ""},
	{"number 1.+3", "1.+3", Tokens{}, "number 1.+3:1:3 Couldn't lex number, junk after decimal point: '+'"},
	{"number 1e!", "1e!", Tokens{}, "number 1e!:1:3 Couldn't lex number, junk after 'E': '!'"},
	{"number 1e+!", "1e+!", Tokens{}, "number 1e+!:1:4 Couldn't lex number, junk after exponent sign: '!'"},

	{"double string \"hi\"", "\"hi\"", Tokens{{kind: tokenStringDouble, data: "hi"}}, ""},
	{"double string \"hi nl\"", "\"hi\n\"", Tokens{{kind: tokenStringDouble, data: "hi\n"}}, ""},
	{"double string \"hi\\\"\"", "\"hi\\\"\"", Tokens{{kind: tokenStringDouble, data: "hi\\\""}}, ""},
	{"double string \"hi\\nl\"", "\"hi\\\n\"", Tokens{{kind: tokenStringDouble, data: "hi\\\n"}}, ""},
	{"double string \"hi", "\"hi", Tokens{}, "double string \"hi:1:1 Unterminated String"},

	{"single string 'hi'", "'hi'", Tokens{{kind: tokenStringSingle, data: "hi"}}, ""},
	{"single string 'hi nl'", "'hi\n'", Tokens{{kind: tokenStringSingle, data: "hi\n"}}, ""},
	{"single string 'hi\\''", "'hi\\''", Tokens{{kind: tokenStringSingle, data: "hi\\'"}}, ""},
	{"single string 'hi\\nl'", "'hi\\\n'", Tokens{{kind: tokenStringSingle, data: "hi\\\n"}}, ""},
	{"single string 'hi", "'hi", Tokens{}, "single string 'hi:1:1 Unterminated String"},

	{"assert", "assert", Tokens{{kind: tokenAssert, data: "assert"}}, ""},
	{"else", "else", Tokens{{kind: tokenElse, data: "else"}}, ""},
	{"error", "error", Tokens{{kind: tokenError, data: "error"}}, ""},
	{"false", "false", Tokens{{kind: tokenFalse, data: "false"}}, ""},
	{"for", "for", Tokens{{kind: tokenFor, data: "for"}}, ""},
	{"function", "function", Tokens{{kind: tokenFunction, data: "function"}}, ""},
	{"if", "if", Tokens{{kind: tokenIf, data: "if"}}, ""},
	{"import", "import", Tokens{{kind: tokenImport, data: "import"}}, ""},
	{"importstr", "importstr", Tokens{{kind: tokenImportStr, data: "importstr"}}, ""},
	{"in", "in", Tokens{{kind: tokenIn, data: "in"}}, ""},
	{"local", "local", Tokens{{kind: tokenLocal, data: "local"}}, ""},
	{"null", "null", Tokens{{kind: tokenNullLit, data: "null"}}, ""},
	{"self", "self", Tokens{{kind: tokenSelf, data: "self"}}, ""},
	{"super", "super", Tokens{{kind: tokenSuper, data: "super"}}, ""},
	{"tailstrict", "tailstrict", Tokens{{kind: tokenTailStrict, data: "tailstrict"}}, ""},
	{"then", "then", Tokens{{kind: tokenThen, data: "then"}}, ""},
	{"true", "true", Tokens{{kind: tokenTrue, data: "true"}}, ""},

	{"identifier", "foobar123", Tokens{{kind: tokenIdentifier, data: "foobar123"}}, ""},
	{"identifier", "foo bar123", Tokens{{kind: tokenIdentifier, data: "foo"}, {kind: tokenIdentifier, data: "bar123"}}, ""},

	{"c++ comment", "// hi", Tokens{}, ""},                                                                     // This test doesn't look at fodder (yet?)
	{"hash comment", "# hi", Tokens{}, ""},                                                                     // This test doesn't look at fodder (yet?)
	{"c comment", "/* hi */", Tokens{}, ""},                                                                    // This test doesn't look at fodder (yet?)
	{"c comment no term", "/* hi", Tokens{}, "c comment no term:1:1 Multi-line comment has no terminating */"}, // This test doesn't look at fodder (yet?)

	{
		"block string spaces",
		`|||
  test
    more
  |||
    foo
|||`,
		Tokens{
			{
				kind:                  tokenStringBlock,
				data:                  "test\n  more\n|||\n  foo\n",
				stringBlockIndent:     "  ",
				stringBlockTermIndent: "",
			},
		},
		"",
	},
	{
		"block string tabs",
		`|||
	test
	  more
	|||
	  foo
|||`,
		Tokens{
			{
				kind:                  tokenStringBlock,
				data:                  "test\n  more\n|||\n  foo\n",
				stringBlockIndent:     "\t",
				stringBlockTermIndent: "",
			},
		},
		"",
	},
	{
		"block string mixed",
		`|||
	  	test
	  	  more
	  	|||
	  	  foo
|||`,
		Tokens{
			{
				kind:                  tokenStringBlock,
				data:                  "test\n  more\n|||\n  foo\n",
				stringBlockIndent:     "\t  \t",
				stringBlockTermIndent: "",
			},
		},
		"",
	},
	{
		"block string blanks",
		`|||

  test


    more
  |||
    foo
|||`,
		Tokens{
			{
				kind:                  tokenStringBlock,
				data:                  "\ntest\n\n\n  more\n|||\n  foo\n",
				stringBlockIndent:     "  ",
				stringBlockTermIndent: "",
			},
		},
		"",
	},
	{
		"block string bad indent",
		`|||
  test
 foo
|||`,
		Tokens{},
		"block string bad indent:1:1 Text block not terminated with |||",
	},
	{
		"block string eof",
		`|||
  test`,
		Tokens{},
		"block string eof:1:1 Unexpected EOF",
	},
	{
		"block string not term",
		`|||
  test
`,
		Tokens{},
		"block string not term:1:1 Text block not terminated with |||",
	},
	{
		"block string no ws",
		`|||
test
|||`,
		Tokens{},
		"block string no ws:1:1 Text block's first line must start with whitespace",
	},

	{"verbatim_string1", `@""`, Tokens{{kind: tokenVerbatimStringDouble, data: ""}}, ""},
	{"verbatim_string2", `@''`, Tokens{{kind: tokenVerbatimStringSingle, data: ""}}, ""},
	{"verbatim_string3", `@""""`, Tokens{{kind: tokenVerbatimStringDouble, data: `"`}}, ""},
	{"verbatim_string4", `@''''`, Tokens{{kind: tokenVerbatimStringSingle, data: "'"}}, ""},
	{"verbatim_string5", `@"\n"`, Tokens{{kind: tokenVerbatimStringDouble, data: "\\n"}}, ""},
	{"verbatim_string6", `@"''"`, Tokens{{kind: tokenVerbatimStringDouble, data: "''"}}, ""},

	{"verbatim_string_unterminated", `@"blah blah`, Tokens{}, "verbatim_string_unterminated:1:1 Unterminated String"},
	{"verbatim_string_junk", `@blah blah`, Tokens{}, "verbatim_string_junk:1:1 Couldn't lex verbatim string, junk after '@': 98"},

	{"op *", "*", Tokens{{kind: tokenOperator, data: "*"}}, ""},
	{"op /", "/", Tokens{{kind: tokenOperator, data: "/"}}, ""},
	{"op %", "%", Tokens{{kind: tokenOperator, data: "%"}}, ""},
	{"op &", "&", Tokens{{kind: tokenOperator, data: "&"}}, ""},
	{"op |", "|", Tokens{{kind: tokenOperator, data: "|"}}, ""},
	{"op ^", "^", Tokens{{kind: tokenOperator, data: "^"}}, ""},
	{"op =", "=", Tokens{{kind: tokenOperator, data: "="}}, ""},
	{"op <", "<", Tokens{{kind: tokenOperator, data: "<"}}, ""},
	{"op >", ">", Tokens{{kind: tokenOperator, data: ">"}}, ""},
	{"op >==|", ">==|", Tokens{{kind: tokenOperator, data: ">==|"}}, ""},

	{"junk", "💩", Tokens{}, "junk:1:1 Could not lex the character '\\U0001f4a9'"},
}

func tokensEqual(ts1, ts2 Tokens) bool {
	if len(ts1) != len(ts2) {
		return false
	}
	for i := range ts1 {
		t1, t2 := ts1[i], ts2[i]
		if t1.kind != t2.kind {
			return false
		}
		if t1.data != t2.data {
			return false
		}
		if t1.stringBlockIndent != t2.stringBlockIndent {
			return false
		}
		if t1.stringBlockTermIndent != t2.stringBlockTermIndent {
			return false
		}
	}
	return true
}

func TestLex(t *testing.T) {
	for _, test := range lexTests {
		// Copy the test tokens and append an EOF token
		testTokens := append(Tokens(nil), test.tokens...)
		testTokens = append(testTokens, tEOF)
		tokens, err := Lex(test.name, test.input)
		var errString string
		if err != nil {
			errString = err.Error()
		}
		if errString != test.errString {
			t.Errorf("%s: error result does not match. got\n\t%+v\nexpected\n\t%+v",
				test.name, errString, test.errString)
		}
		if err == nil && !tokensEqual(tokens, testTokens) {
			t.Errorf("%s: got\n\t%+v\nexpected\n\t%+v", test.name, tokens, testTokens)
		}
	}
}

// TODO: test fodder, test position reporting
