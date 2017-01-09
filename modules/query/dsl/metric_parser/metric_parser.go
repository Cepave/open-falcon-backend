package metric_parser

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"
)

var g = &grammar{
	rules: []*rule{
		{
			name: "MetricFilter",
			pos:  position{line: 5, col: 1, offset: 27},
			expr: &actionExpr{
				pos: position{line: 5, col: 16, offset: 42},
				run: (*parser).callonMetricFilter1,
				expr: &seqExpr{
					pos: position{line: 5, col: 16, offset: 42},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 5, col: 16, offset: 42},
							label: "e",
							expr: &ruleRefExpr{
								pos:  position{line: 5, col: 18, offset: 44},
								name: "Expr",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 5, col: 23, offset: 49},
							name: "EOF",
						},
					},
				},
			},
		},
		{
			name: "Expr",
			pos:  position{line: 9, col: 1, offset: 73},
			expr: &actionExpr{
				pos: position{line: 9, col: 8, offset: 80},
				run: (*parser).callonExpr1,
				expr: &seqExpr{
					pos: position{line: 9, col: 8, offset: 80},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 9, col: 8, offset: 80},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 9, col: 10, offset: 82},
							label: "t",
							expr: &ruleRefExpr{
								pos:  position{line: 9, col: 12, offset: 84},
								name: "OrTerm",
							},
						},
						&labeledExpr{
							pos:   position{line: 9, col: 19, offset: 91},
							label: "orTerms",
							expr: &zeroOrMoreExpr{
								pos: position{line: 9, col: 27, offset: 99},
								expr: &ruleRefExpr{
									pos:  position{line: 9, col: 27, offset: 99},
									name: "OrClause",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 9, col: 37, offset: 109},
							name: "_",
						},
					},
				},
			},
		},
		{
			name: "OrClause",
			pos:  position{line: 13, col: 1, offset: 165},
			expr: &choiceExpr{
				pos: position{line: 13, col: 12, offset: 176},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 13, col: 12, offset: 176},
						run: (*parser).callonOrClause2,
						expr: &seqExpr{
							pos: position{line: 13, col: 12, offset: 176},
							exprs: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 13, col: 12, offset: 176},
									name: "_",
								},
								&litMatcher{
									pos:        position{line: 13, col: 14, offset: 178},
									val:        "or",
									ignoreCase: false,
								},
								&oneOrMoreExpr{
									pos: position{line: 13, col: 19, offset: 183},
									expr: &ruleRefExpr{
										pos:  position{line: 13, col: 19, offset: 183},
										name: "EMPTY_CHAR",
									},
								},
								&labeledExpr{
									pos:   position{line: 13, col: 31, offset: 195},
									label: "t",
									expr: &ruleRefExpr{
										pos:  position{line: 13, col: 33, offset: 197},
										name: "OrTerm",
									},
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 15, col: 5, offset: 225},
						run: (*parser).callonOrClause10,
						expr: &seqExpr{
							pos: position{line: 15, col: 5, offset: 225},
							exprs: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 15, col: 5, offset: 225},
									name: "_",
								},
								&litMatcher{
									pos:        position{line: 15, col: 7, offset: 227},
									val:        "or",
									ignoreCase: false,
								},
								&notExpr{
									pos: position{line: 15, col: 12, offset: 232},
									expr: &seqExpr{
										pos: position{line: 15, col: 14, offset: 234},
										exprs: []interface{}{
											&oneOrMoreExpr{
												pos: position{line: 15, col: 14, offset: 234},
												expr: &ruleRefExpr{
													pos:  position{line: 15, col: 14, offset: 234},
													name: "EMPTY_CHAR",
												},
											},
											&ruleRefExpr{
												pos:  position{line: 15, col: 26, offset: 246},
												name: "OrTerm",
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "OrTerm",
			pos:  position{line: 19, col: 1, offset: 313},
			expr: &actionExpr{
				pos: position{line: 19, col: 10, offset: 322},
				run: (*parser).callonOrTerm1,
				expr: &seqExpr{
					pos: position{line: 19, col: 10, offset: 322},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 19, col: 10, offset: 322},
							label: "t",
							expr: &ruleRefExpr{
								pos:  position{line: 19, col: 12, offset: 324},
								name: "Term",
							},
						},
						&labeledExpr{
							pos:   position{line: 19, col: 17, offset: 329},
							label: "andTerms",
							expr: &zeroOrMoreExpr{
								pos: position{line: 19, col: 26, offset: 338},
								expr: &ruleRefExpr{
									pos:  position{line: 19, col: 26, offset: 338},
									name: "AndClause",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "AndClause",
			pos:  position{line: 23, col: 1, offset: 405},
			expr: &choiceExpr{
				pos: position{line: 23, col: 13, offset: 417},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 23, col: 13, offset: 417},
						run: (*parser).callonAndClause2,
						expr: &seqExpr{
							pos: position{line: 23, col: 13, offset: 417},
							exprs: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 23, col: 13, offset: 417},
									name: "_",
								},
								&litMatcher{
									pos:        position{line: 23, col: 15, offset: 419},
									val:        "and",
									ignoreCase: false,
								},
								&oneOrMoreExpr{
									pos: position{line: 23, col: 21, offset: 425},
									expr: &ruleRefExpr{
										pos:  position{line: 23, col: 21, offset: 425},
										name: "EMPTY_CHAR",
									},
								},
								&labeledExpr{
									pos:   position{line: 23, col: 33, offset: 437},
									label: "t",
									expr: &ruleRefExpr{
										pos:  position{line: 23, col: 35, offset: 439},
										name: "Term",
									},
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 25, col: 5, offset: 465},
						run: (*parser).callonAndClause10,
						expr: &seqExpr{
							pos: position{line: 25, col: 5, offset: 465},
							exprs: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 25, col: 5, offset: 465},
									name: "_",
								},
								&litMatcher{
									pos:        position{line: 25, col: 7, offset: 467},
									val:        "and",
									ignoreCase: false,
								},
								&notExpr{
									pos: position{line: 25, col: 13, offset: 473},
									expr: &seqExpr{
										pos: position{line: 25, col: 15, offset: 475},
										exprs: []interface{}{
											&oneOrMoreExpr{
												pos: position{line: 25, col: 15, offset: 475},
												expr: &ruleRefExpr{
													pos:  position{line: 25, col: 15, offset: 475},
													name: "EMPTY_CHAR",
												},
											},
											&ruleRefExpr{
												pos:  position{line: 25, col: 27, offset: 487},
												name: "OrTerm",
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Term",
			pos:  position{line: 29, col: 1, offset: 555},
			expr: &actionExpr{
				pos: position{line: 29, col: 8, offset: 562},
				run: (*parser).callonTerm1,
				expr: &labeledExpr{
					pos:   position{line: 29, col: 8, offset: 562},
					label: "t",
					expr: &choiceExpr{
						pos: position{line: 29, col: 11, offset: 565},
						alternatives: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 29, col: 11, offset: 565},
								name: "QuoteExpr",
							},
							&ruleRefExpr{
								pos:  position{line: 29, col: 23, offset: 577},
								name: "FactorOp",
							},
						},
					},
				},
			},
		},
		{
			name: "QuoteExpr",
			pos:  position{line: 33, col: 1, offset: 607},
			expr: &actionExpr{
				pos: position{line: 33, col: 13, offset: 619},
				run: (*parser).callonQuoteExpr1,
				expr: &seqExpr{
					pos: position{line: 33, col: 13, offset: 619},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 33, col: 13, offset: 619},
							val:        "(",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 33, col: 17, offset: 623},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 33, col: 19, offset: 625},
							label: "e",
							expr: &ruleRefExpr{
								pos:  position{line: 33, col: 21, offset: 627},
								name: "Expr",
							},
						},
						&labeledExpr{
							pos:   position{line: 33, col: 26, offset: 632},
							label: "rightParen",
							expr: &zeroOrOneExpr{
								pos: position{line: 33, col: 37, offset: 643},
								expr: &litMatcher{
									pos:        position{line: 33, col: 37, offset: 643},
									val:        ")",
									ignoreCase: false,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "FactorOp",
			pos:  position{line: 40, col: 1, offset: 772},
			expr: &choiceExpr{
				pos: position{line: 40, col: 12, offset: 783},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 40, col: 12, offset: 783},
						run: (*parser).callonFactorOp2,
						expr: &seqExpr{
							pos: position{line: 40, col: 12, offset: 783},
							exprs: []interface{}{
								&labeledExpr{
									pos:   position{line: 40, col: 12, offset: 783},
									label: "left",
									expr: &ruleRefExpr{
										pos:  position{line: 40, col: 17, offset: 788},
										name: "Factor",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 40, col: 24, offset: 795},
									name: "_",
								},
								&labeledExpr{
									pos:   position{line: 40, col: 26, offset: 797},
									label: "op",
									expr: &ruleRefExpr{
										pos:  position{line: 40, col: 29, offset: 800},
										name: "Op",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 40, col: 32, offset: 803},
									name: "_",
								},
								&labeledExpr{
									pos:   position{line: 40, col: 34, offset: 805},
									label: "right",
									expr: &ruleRefExpr{
										pos:  position{line: 40, col: 40, offset: 811},
										name: "Factor",
									},
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 42, col: 5, offset: 885},
						run: (*parser).callonFactorOp12,
						expr: &seqExpr{
							pos: position{line: 42, col: 5, offset: 885},
							exprs: []interface{}{
								&labeledExpr{
									pos:   position{line: 42, col: 5, offset: 885},
									label: "left",
									expr: &ruleRefExpr{
										pos:  position{line: 42, col: 10, offset: 890},
										name: "Factor",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 42, col: 17, offset: 897},
									name: "_",
								},
								&ruleRefExpr{
									pos:  position{line: 42, col: 19, offset: 899},
									name: "Op",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Op",
			pos:  position{line: 46, col: 1, offset: 969},
			expr: &actionExpr{
				pos: position{line: 46, col: 6, offset: 974},
				run: (*parser).callonOp1,
				expr: &ruleRefExpr{
					pos:  position{line: 46, col: 6, offset: 974},
					name: "VIABLE_CHARS",
				},
			},
		},
		{
			name: "Factor",
			pos:  position{line: 57, col: 1, offset: 1142},
			expr: &choiceExpr{
				pos: position{line: 57, col: 10, offset: 1151},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 57, col: 10, offset: 1151},
						run: (*parser).callonFactor2,
						expr: &seqExpr{
							pos: position{line: 57, col: 10, offset: 1151},
							exprs: []interface{}{
								&labeledExpr{
									pos:   position{line: 57, col: 10, offset: 1151},
									label: "v",
									expr: &choiceExpr{
										pos: position{line: 57, col: 13, offset: 1154},
										alternatives: []interface{}{
											&ruleRefExpr{
												pos:  position{line: 57, col: 13, offset: 1154},
												name: "Metric",
											},
											&ruleRefExpr{
												pos:  position{line: 57, col: 22, offset: 1163},
												name: "NUMBER",
											},
										},
									},
								},
								&andExpr{
									pos: position{line: 57, col: 30, offset: 1171},
									expr: &ruleRefExpr{
										pos:  position{line: 57, col: 31, offset: 1172},
										name: "END_WORD",
									},
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 59, col: 5, offset: 1202},
						run: (*parser).callonFactor10,
						expr: &ruleRefExpr{
							pos:  position{line: 59, col: 5, offset: 1202},
							name: "VIABLE_CHARS",
						},
					},
				},
			},
		},
		{
			name: "Metric",
			pos:  position{line: 63, col: 1, offset: 1279},
			expr: &actionExpr{
				pos: position{line: 63, col: 10, offset: 1288},
				run: (*parser).callonMetric1,
				expr: &seqExpr{
					pos: position{line: 63, col: 10, offset: 1288},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 63, col: 10, offset: 1288},
							val:        "$",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 63, col: 14, offset: 1292},
							label: "metric",
							expr: &choiceExpr{
								pos: position{line: 63, col: 22, offset: 1300},
								alternatives: []interface{}{
									&litMatcher{
										pos:        position{line: 63, col: 22, offset: 1300},
										val:        "max",
										ignoreCase: false,
									},
									&litMatcher{
										pos:        position{line: 63, col: 30, offset: 1308},
										val:        "min",
										ignoreCase: false,
									},
									&litMatcher{
										pos:        position{line: 63, col: 38, offset: 1316},
										val:        "avg",
										ignoreCase: false,
									},
									&litMatcher{
										pos:        position{line: 63, col: 46, offset: 1324},
										val:        "med",
										ignoreCase: false,
									},
									&litMatcher{
										pos:        position{line: 63, col: 54, offset: 1332},
										val:        "mdev",
										ignoreCase: false,
									},
									&litMatcher{
										pos:        position{line: 63, col: 63, offset: 1341},
										val:        "loss",
										ignoreCase: false,
									},
									&litMatcher{
										pos:        position{line: 63, col: 72, offset: 1350},
										val:        "count",
										ignoreCase: false,
									},
									&litMatcher{
										pos:        position{line: 63, col: 82, offset: 1360},
										val:        "pck_sent",
										ignoreCase: false,
									},
									&litMatcher{
										pos:        position{line: 63, col: 95, offset: 1373},
										val:        "pck_received",
										ignoreCase: false,
									},
									&litMatcher{
										pos:        position{line: 63, col: 112, offset: 1390},
										val:        "num_agent",
										ignoreCase: false,
									},
									&litMatcher{
										pos:        position{line: 63, col: 126, offset: 1404},
										val:        "num_target",
										ignoreCase: false,
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "NUMBER",
			pos:  position{line: 68, col: 1, offset: 1499},
			expr: &actionExpr{
				pos: position{line: 68, col: 10, offset: 1508},
				run: (*parser).callonNUMBER1,
				expr: &seqExpr{
					pos: position{line: 68, col: 10, offset: 1508},
					exprs: []interface{}{
						&oneOrMoreExpr{
							pos: position{line: 68, col: 10, offset: 1508},
							expr: &charClassMatcher{
								pos:        position{line: 68, col: 10, offset: 1508},
								val:        "[0-9]",
								ranges:     []rune{'0', '9'},
								ignoreCase: false,
								inverted:   false,
							},
						},
						&zeroOrOneExpr{
							pos: position{line: 68, col: 17, offset: 1515},
							expr: &seqExpr{
								pos: position{line: 68, col: 18, offset: 1516},
								exprs: []interface{}{
									&litMatcher{
										pos:        position{line: 68, col: 18, offset: 1516},
										val:        ".",
										ignoreCase: false,
									},
									&oneOrMoreExpr{
										pos: position{line: 68, col: 22, offset: 1520},
										expr: &charClassMatcher{
											pos:        position{line: 68, col: 22, offset: 1520},
											val:        "[0-9]",
											ranges:     []rune{'0', '9'},
											ignoreCase: false,
											inverted:   false,
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "VIABLE_CHARS",
			pos:  position{line: 72, col: 1, offset: 1562},
			expr: &actionExpr{
				pos: position{line: 72, col: 16, offset: 1577},
				run: (*parser).callonVIABLE_CHARS1,
				expr: &oneOrMoreExpr{
					pos: position{line: 72, col: 16, offset: 1577},
					expr: &charClassMatcher{
						pos:        position{line: 72, col: 16, offset: 1577},
						val:        "[^ \\t\\n\\r]",
						chars:      []rune{' ', '\t', '\n', '\r'},
						ignoreCase: false,
						inverted:   true,
					},
				},
			},
		},
		{
			name: "_",
			pos:  position{line: 76, col: 1, offset: 1622},
			expr: &zeroOrMoreExpr{
				pos: position{line: 76, col: 5, offset: 1626},
				expr: &ruleRefExpr{
					pos:  position{line: 76, col: 5, offset: 1626},
					name: "EMPTY_CHAR",
				},
			},
		},
		{
			name: "END_WORD",
			pos:  position{line: 77, col: 1, offset: 1638},
			expr: &choiceExpr{
				pos: position{line: 77, col: 12, offset: 1649},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 77, col: 12, offset: 1649},
						name: "EOF",
					},
					&oneOrMoreExpr{
						pos: position{line: 77, col: 18, offset: 1655},
						expr: &ruleRefExpr{
							pos:  position{line: 77, col: 18, offset: 1655},
							name: "EMPTY_CHAR",
						},
					},
					&litMatcher{
						pos:        position{line: 77, col: 32, offset: 1669},
						val:        ")",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "EMPTY_CHAR",
			pos:  position{line: 78, col: 1, offset: 1673},
			expr: &charClassMatcher{
				pos:        position{line: 78, col: 14, offset: 1686},
				val:        "[ \\t\\n\\r]",
				chars:      []rune{' ', '\t', '\n', '\r'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "EOF",
			pos:  position{line: 79, col: 1, offset: 1696},
			expr: &notExpr{
				pos: position{line: 79, col: 7, offset: 1702},
				expr: &anyMatcher{
					line: 79, col: 8, offset: 1703,
				},
			},
		},
	},
}

func (c *current) onMetricFilter1(e interface{}) (interface{}, error) {
	return e, nil
}

func (p *parser) callonMetricFilter1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onMetricFilter1(stack["e"])
}

func (c *current) onExpr1(t, orTerms interface{}) (interface{}, error) {
	return newBoolFilterImpl(true, t, orTerms), nil
}

func (p *parser) callonExpr1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onExpr1(stack["t"], stack["orTerms"])
}

func (c *current) onOrClause2(t interface{}) (interface{}, error) {
	return t, nil
}

func (p *parser) callonOrClause2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOrClause2(stack["t"])
}

func (c *current) onOrClause10() (interface{}, error) {
	panic(fmt.Errorf("OR what? \"%v\"", string(c.text)))
}

func (p *parser) callonOrClause10() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOrClause10()
}

func (c *current) onOrTerm1(t, andTerms interface{}) (interface{}, error) {
	return newBoolFilterImpl(false, t, andTerms), nil
}

func (p *parser) callonOrTerm1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOrTerm1(stack["t"], stack["andTerms"])
}

func (c *current) onAndClause2(t interface{}) (interface{}, error) {
	return t, nil
}

func (p *parser) callonAndClause2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onAndClause2(stack["t"])
}

func (c *current) onAndClause10() (interface{}, error) {
	panic(fmt.Errorf("AND what? \"%v\"", string(c.text)))
}

func (p *parser) callonAndClause10() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onAndClause10()
}

func (c *current) onTerm1(t interface{}) (interface{}, error) {
	return t, nil
}

func (p *parser) callonTerm1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onTerm1(stack["t"])
}

func (c *current) onQuoteExpr1(e, rightParen interface{}) (interface{}, error) {
	if rightParen == nil {
		panic(fmt.Errorf("Need right parenthese ')' for: \"%v\"", string(c.text)))
	}

	return e, nil
}

func (p *parser) callonQuoteExpr1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onQuoteExpr1(stack["e"], stack["rightParen"])
}

func (c *current) onFactorOp2(left, op, right interface{}) (interface{}, error) {
	return newFilterImpl(left, string(op.([]byte)), right), nil
}

func (p *parser) callonFactorOp2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onFactorOp2(stack["left"], stack["op"], stack["right"])
}

func (c *current) onFactorOp12(left interface{}) (interface{}, error) {
	panic(fmt.Errorf("Need right factor: [%s]", string(c.text)))
}

func (p *parser) callonFactorOp12() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onFactorOp12(stack["left"])
}

func (c *current) onOp1() (interface{}, error) {
	op := string(c.text)

	switch op {
	case ">=", "<=", "==", "!=", ">", "<":
		return c.text, nil
	}

	panic(fmt.Errorf("Unknown operator: [%s]", op))
}

func (p *parser) callonOp1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOp1()
}

func (c *current) onFactor2(v interface{}) (interface{}, error) {
	return v, nil
}

func (p *parser) callonFactor2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onFactor2(stack["v"])
}

func (c *current) onFactor10() (interface{}, error) {
	panic(fmt.Errorf("Unknown factor: [%s]", string(c.text)))
}

func (p *parser) callonFactor10() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onFactor10()
}

func (c *current) onMetric1(metric interface{}) (interface{}, error) {
	metricName := string(metric.([]byte))
	return mapOfMetric[metricName], nil
}

func (p *parser) callonMetric1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onMetric1(stack["metric"])
}

func (c *current) onNUMBER1() (interface{}, error) {
	return string(c.text), nil
}

func (p *parser) callonNUMBER1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onNUMBER1()
}

func (c *current) onVIABLE_CHARS1() (interface{}, error) {
	return string(c.text), nil
}

func (p *parser) callonVIABLE_CHARS1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onVIABLE_CHARS1()
}

var (
	// errNoRule is returned when the grammar to parse has no rule.
	errNoRule = errors.New("grammar has no rule")

	// errInvalidEncoding is returned when the source is not properly
	// utf8-encoded.
	errInvalidEncoding = errors.New("invalid encoding")

	// errNoMatch is returned if no match could be found.
	errNoMatch = errors.New("no match found")
)

// Option is a function that can set an option on the parser. It returns
// the previous setting as an Option.
type Option func(*parser) Option

// Debug creates an Option to set the debug flag to b. When set to true,
// debugging information is printed to stdout while parsing.
//
// The default is false.
func Debug(b bool) Option {
	return func(p *parser) Option {
		old := p.debug
		p.debug = b
		return Debug(old)
	}
}

// Memoize creates an Option to set the memoize flag to b. When set to true,
// the parser will cache all results so each expression is evaluated only
// once. This guarantees linear parsing time even for pathological cases,
// at the expense of more memory and slower times for typical cases.
//
// The default is false.
func Memoize(b bool) Option {
	return func(p *parser) Option {
		old := p.memoize
		p.memoize = b
		return Memoize(old)
	}
}

// Recover creates an Option to set the recover flag to b. When set to
// true, this causes the parser to recover from panics and convert it
// to an error. Setting it to false can be useful while debugging to
// access the full stack trace.
//
// The default is true.
func Recover(b bool) Option {
	return func(p *parser) Option {
		old := p.recover
		p.recover = b
		return Recover(old)
	}
}

// ParseFile parses the file identified by filename.
func ParseFile(filename string, opts ...Option) (interface{}, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ParseReader(filename, f, opts...)
}

// ParseReader parses the data from r using filename as information in the
// error messages.
func ParseReader(filename string, r io.Reader, opts ...Option) (interface{}, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return Parse(filename, b, opts...)
}

// Parse parses the data from b using filename as information in the
// error messages.
func Parse(filename string, b []byte, opts ...Option) (interface{}, error) {
	return newParser(filename, b, opts...).parse(g)
}

// position records a position in the text.
type position struct {
	line, col, offset int
}

func (p position) String() string {
	return fmt.Sprintf("%d:%d [%d]", p.line, p.col, p.offset)
}

// savepoint stores all state required to go back to this point in the
// parser.
type savepoint struct {
	position
	rn rune
	w  int
}

type current struct {
	pos  position // start position of the match
	text []byte   // raw text of the match
}

// the AST types...

type grammar struct {
	pos   position
	rules []*rule
}

type rule struct {
	pos         position
	name        string
	displayName string
	expr        interface{}
}

type choiceExpr struct {
	pos          position
	alternatives []interface{}
}

type actionExpr struct {
	pos  position
	expr interface{}
	run  func(*parser) (interface{}, error)
}

type seqExpr struct {
	pos   position
	exprs []interface{}
}

type labeledExpr struct {
	pos   position
	label string
	expr  interface{}
}

type expr struct {
	pos  position
	expr interface{}
}

type andExpr expr
type notExpr expr
type zeroOrOneExpr expr
type zeroOrMoreExpr expr
type oneOrMoreExpr expr

type ruleRefExpr struct {
	pos  position
	name string
}

type andCodeExpr struct {
	pos position
	run func(*parser) (bool, error)
}

type notCodeExpr struct {
	pos position
	run func(*parser) (bool, error)
}

type litMatcher struct {
	pos        position
	val        string
	ignoreCase bool
}

type charClassMatcher struct {
	pos        position
	val        string
	chars      []rune
	ranges     []rune
	classes    []*unicode.RangeTable
	ignoreCase bool
	inverted   bool
}

type anyMatcher position

// errList cumulates the errors found by the parser.
type errList []error

func (e *errList) add(err error) {
	*e = append(*e, err)
}

func (e errList) err() error {
	if len(e) == 0 {
		return nil
	}
	e.dedupe()
	return e
}

func (e *errList) dedupe() {
	var cleaned []error
	set := make(map[string]bool)
	for _, err := range *e {
		if msg := err.Error(); !set[msg] {
			set[msg] = true
			cleaned = append(cleaned, err)
		}
	}
	*e = cleaned
}

func (e errList) Error() string {
	switch len(e) {
	case 0:
		return ""
	case 1:
		return e[0].Error()
	default:
		var buf bytes.Buffer

		for i, err := range e {
			if i > 0 {
				buf.WriteRune('\n')
			}
			buf.WriteString(err.Error())
		}
		return buf.String()
	}
}

// parserError wraps an error with a prefix indicating the rule in which
// the error occurred. The original error is stored in the Inner field.
type parserError struct {
	Inner  error
	pos    position
	prefix string
}

// Error returns the error message.
func (p *parserError) Error() string {
	return p.prefix + ": " + p.Inner.Error()
}

// newParser creates a parser with the specified input source and options.
func newParser(filename string, b []byte, opts ...Option) *parser {
	p := &parser{
		filename: filename,
		errs:     new(errList),
		data:     b,
		pt:       savepoint{position: position{line: 1}},
		recover:  true,
	}
	p.setOptions(opts)
	return p
}

// setOptions applies the options to the parser.
func (p *parser) setOptions(opts []Option) {
	for _, opt := range opts {
		opt(p)
	}
}

type resultTuple struct {
	v   interface{}
	b   bool
	end savepoint
}

type parser struct {
	filename string
	pt       savepoint
	cur      current

	data []byte
	errs *errList

	recover bool
	debug   bool
	depth   int

	memoize bool
	// memoization table for the packrat algorithm:
	// map[offset in source] map[expression or rule] {value, match}
	memo map[int]map[interface{}]resultTuple

	// rules table, maps the rule identifier to the rule node
	rules map[string]*rule
	// variables stack, map of label to value
	vstack []map[string]interface{}
	// rule stack, allows identification of the current rule in errors
	rstack []*rule

	// stats
	exprCnt int
}

// push a variable set on the vstack.
func (p *parser) pushV() {
	if cap(p.vstack) == len(p.vstack) {
		// create new empty slot in the stack
		p.vstack = append(p.vstack, nil)
	} else {
		// slice to 1 more
		p.vstack = p.vstack[:len(p.vstack)+1]
	}

	// get the last args set
	m := p.vstack[len(p.vstack)-1]
	if m != nil && len(m) == 0 {
		// empty map, all good
		return
	}

	m = make(map[string]interface{})
	p.vstack[len(p.vstack)-1] = m
}

// pop a variable set from the vstack.
func (p *parser) popV() {
	// if the map is not empty, clear it
	m := p.vstack[len(p.vstack)-1]
	if len(m) > 0 {
		// GC that map
		p.vstack[len(p.vstack)-1] = nil
	}
	p.vstack = p.vstack[:len(p.vstack)-1]
}

func (p *parser) print(prefix, s string) string {
	if !p.debug {
		return s
	}

	fmt.Printf("%s %d:%d:%d: %s [%#U]\n",
		prefix, p.pt.line, p.pt.col, p.pt.offset, s, p.pt.rn)
	return s
}

func (p *parser) in(s string) string {
	p.depth++
	return p.print(strings.Repeat(" ", p.depth)+">", s)
}

func (p *parser) out(s string) string {
	p.depth--
	return p.print(strings.Repeat(" ", p.depth)+"<", s)
}

func (p *parser) addErr(err error) {
	p.addErrAt(err, p.pt.position)
}

func (p *parser) addErrAt(err error, pos position) {
	var buf bytes.Buffer
	if p.filename != "" {
		buf.WriteString(p.filename)
	}
	if buf.Len() > 0 {
		buf.WriteString(":")
	}
	buf.WriteString(fmt.Sprintf("%d:%d (%d)", pos.line, pos.col, pos.offset))
	if len(p.rstack) > 0 {
		if buf.Len() > 0 {
			buf.WriteString(": ")
		}
		rule := p.rstack[len(p.rstack)-1]
		if rule.displayName != "" {
			buf.WriteString("rule " + rule.displayName)
		} else {
			buf.WriteString("rule " + rule.name)
		}
	}
	pe := &parserError{Inner: err, prefix: buf.String()}
	p.errs.add(pe)
}

// read advances the parser to the next rune.
func (p *parser) read() {
	p.pt.offset += p.pt.w
	rn, n := utf8.DecodeRune(p.data[p.pt.offset:])
	p.pt.rn = rn
	p.pt.w = n
	p.pt.col++
	if rn == '\n' {
		p.pt.line++
		p.pt.col = 0
	}

	if rn == utf8.RuneError {
		if n > 0 {
			p.addErr(errInvalidEncoding)
		}
	}
}

// restore parser position to the savepoint pt.
func (p *parser) restore(pt savepoint) {
	if p.debug {
		defer p.out(p.in("restore"))
	}
	if pt.offset == p.pt.offset {
		return
	}
	p.pt = pt
}

// get the slice of bytes from the savepoint start to the current position.
func (p *parser) sliceFrom(start savepoint) []byte {
	return p.data[start.position.offset:p.pt.position.offset]
}

func (p *parser) getMemoized(node interface{}) (resultTuple, bool) {
	if len(p.memo) == 0 {
		return resultTuple{}, false
	}
	m := p.memo[p.pt.offset]
	if len(m) == 0 {
		return resultTuple{}, false
	}
	res, ok := m[node]
	return res, ok
}

func (p *parser) setMemoized(pt savepoint, node interface{}, tuple resultTuple) {
	if p.memo == nil {
		p.memo = make(map[int]map[interface{}]resultTuple)
	}
	m := p.memo[pt.offset]
	if m == nil {
		m = make(map[interface{}]resultTuple)
		p.memo[pt.offset] = m
	}
	m[node] = tuple
}

func (p *parser) buildRulesTable(g *grammar) {
	p.rules = make(map[string]*rule, len(g.rules))
	for _, r := range g.rules {
		p.rules[r.name] = r
	}
}

func (p *parser) parse(g *grammar) (val interface{}, err error) {
	if len(g.rules) == 0 {
		p.addErr(errNoRule)
		return nil, p.errs.err()
	}

	// TODO : not super critical but this could be generated
	p.buildRulesTable(g)

	if p.recover {
		// panic can be used in action code to stop parsing immediately
		// and return the panic as an error.
		defer func() {
			if e := recover(); e != nil {
				if p.debug {
					defer p.out(p.in("panic handler"))
				}
				val = nil
				switch e := e.(type) {
				case error:
					p.addErr(e)
				default:
					p.addErr(fmt.Errorf("%v", e))
				}
				err = p.errs.err()
			}
		}()
	}

	// start rule is rule [0]
	p.read() // advance to first rune
	val, ok := p.parseRule(g.rules[0])
	if !ok {
		if len(*p.errs) == 0 {
			// make sure this doesn't go out silently
			p.addErr(errNoMatch)
		}
		return nil, p.errs.err()
	}
	return val, p.errs.err()
}

func (p *parser) parseRule(rule *rule) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseRule " + rule.name))
	}

	if p.memoize {
		res, ok := p.getMemoized(rule)
		if ok {
			p.restore(res.end)
			return res.v, res.b
		}
	}

	start := p.pt
	p.rstack = append(p.rstack, rule)
	p.pushV()
	val, ok := p.parseExpr(rule.expr)
	p.popV()
	p.rstack = p.rstack[:len(p.rstack)-1]
	if ok && p.debug {
		p.print(strings.Repeat(" ", p.depth)+"MATCH", string(p.sliceFrom(start)))
	}

	if p.memoize {
		p.setMemoized(start, rule, resultTuple{val, ok, p.pt})
	}
	return val, ok
}

func (p *parser) parseExpr(expr interface{}) (interface{}, bool) {
	var pt savepoint
	var ok bool

	if p.memoize {
		res, ok := p.getMemoized(expr)
		if ok {
			p.restore(res.end)
			return res.v, res.b
		}
		pt = p.pt
	}

	p.exprCnt++
	var val interface{}
	switch expr := expr.(type) {
	case *actionExpr:
		val, ok = p.parseActionExpr(expr)
	case *andCodeExpr:
		val, ok = p.parseAndCodeExpr(expr)
	case *andExpr:
		val, ok = p.parseAndExpr(expr)
	case *anyMatcher:
		val, ok = p.parseAnyMatcher(expr)
	case *charClassMatcher:
		val, ok = p.parseCharClassMatcher(expr)
	case *choiceExpr:
		val, ok = p.parseChoiceExpr(expr)
	case *labeledExpr:
		val, ok = p.parseLabeledExpr(expr)
	case *litMatcher:
		val, ok = p.parseLitMatcher(expr)
	case *notCodeExpr:
		val, ok = p.parseNotCodeExpr(expr)
	case *notExpr:
		val, ok = p.parseNotExpr(expr)
	case *oneOrMoreExpr:
		val, ok = p.parseOneOrMoreExpr(expr)
	case *ruleRefExpr:
		val, ok = p.parseRuleRefExpr(expr)
	case *seqExpr:
		val, ok = p.parseSeqExpr(expr)
	case *zeroOrMoreExpr:
		val, ok = p.parseZeroOrMoreExpr(expr)
	case *zeroOrOneExpr:
		val, ok = p.parseZeroOrOneExpr(expr)
	default:
		panic(fmt.Sprintf("unknown expression type %T", expr))
	}
	if p.memoize {
		p.setMemoized(pt, expr, resultTuple{val, ok, p.pt})
	}
	return val, ok
}

func (p *parser) parseActionExpr(act *actionExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseActionExpr"))
	}

	start := p.pt
	val, ok := p.parseExpr(act.expr)
	if ok {
		p.cur.pos = start.position
		p.cur.text = p.sliceFrom(start)
		actVal, err := act.run(p)
		if err != nil {
			p.addErrAt(err, start.position)
		}
		val = actVal
	}
	if ok && p.debug {
		p.print(strings.Repeat(" ", p.depth)+"MATCH", string(p.sliceFrom(start)))
	}
	return val, ok
}

func (p *parser) parseAndCodeExpr(and *andCodeExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAndCodeExpr"))
	}

	ok, err := and.run(p)
	if err != nil {
		p.addErr(err)
	}
	return nil, ok
}

func (p *parser) parseAndExpr(and *andExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAndExpr"))
	}

	pt := p.pt
	p.pushV()
	_, ok := p.parseExpr(and.expr)
	p.popV()
	p.restore(pt)
	return nil, ok
}

func (p *parser) parseAnyMatcher(any *anyMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAnyMatcher"))
	}

	if p.pt.rn != utf8.RuneError {
		start := p.pt
		p.read()
		return p.sliceFrom(start), true
	}
	return nil, false
}

func (p *parser) parseCharClassMatcher(chr *charClassMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseCharClassMatcher"))
	}

	cur := p.pt.rn
	// can't match EOF
	if cur == utf8.RuneError {
		return nil, false
	}
	start := p.pt
	if chr.ignoreCase {
		cur = unicode.ToLower(cur)
	}

	// try to match in the list of available chars
	for _, rn := range chr.chars {
		if rn == cur {
			if chr.inverted {
				return nil, false
			}
			p.read()
			return p.sliceFrom(start), true
		}
	}

	// try to match in the list of ranges
	for i := 0; i < len(chr.ranges); i += 2 {
		if cur >= chr.ranges[i] && cur <= chr.ranges[i+1] {
			if chr.inverted {
				return nil, false
			}
			p.read()
			return p.sliceFrom(start), true
		}
	}

	// try to match in the list of Unicode classes
	for _, cl := range chr.classes {
		if unicode.Is(cl, cur) {
			if chr.inverted {
				return nil, false
			}
			p.read()
			return p.sliceFrom(start), true
		}
	}

	if chr.inverted {
		p.read()
		return p.sliceFrom(start), true
	}
	return nil, false
}

func (p *parser) parseChoiceExpr(ch *choiceExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseChoiceExpr"))
	}

	for _, alt := range ch.alternatives {
		p.pushV()
		val, ok := p.parseExpr(alt)
		p.popV()
		if ok {
			return val, ok
		}
	}
	return nil, false
}

func (p *parser) parseLabeledExpr(lab *labeledExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseLabeledExpr"))
	}

	p.pushV()
	val, ok := p.parseExpr(lab.expr)
	p.popV()
	if ok && lab.label != "" {
		m := p.vstack[len(p.vstack)-1]
		m[lab.label] = val
	}
	return val, ok
}

func (p *parser) parseLitMatcher(lit *litMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseLitMatcher"))
	}

	start := p.pt
	for _, want := range lit.val {
		cur := p.pt.rn
		if lit.ignoreCase {
			cur = unicode.ToLower(cur)
		}
		if cur != want {
			p.restore(start)
			return nil, false
		}
		p.read()
	}
	return p.sliceFrom(start), true
}

func (p *parser) parseNotCodeExpr(not *notCodeExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseNotCodeExpr"))
	}

	ok, err := not.run(p)
	if err != nil {
		p.addErr(err)
	}
	return nil, !ok
}

func (p *parser) parseNotExpr(not *notExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseNotExpr"))
	}

	pt := p.pt
	p.pushV()
	_, ok := p.parseExpr(not.expr)
	p.popV()
	p.restore(pt)
	return nil, !ok
}

func (p *parser) parseOneOrMoreExpr(expr *oneOrMoreExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseOneOrMoreExpr"))
	}

	var vals []interface{}

	for {
		p.pushV()
		val, ok := p.parseExpr(expr.expr)
		p.popV()
		if !ok {
			if len(vals) == 0 {
				// did not match once, no match
				return nil, false
			}
			return vals, true
		}
		vals = append(vals, val)
	}
}

func (p *parser) parseRuleRefExpr(ref *ruleRefExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseRuleRefExpr " + ref.name))
	}

	if ref.name == "" {
		panic(fmt.Sprintf("%s: invalid rule: missing name", ref.pos))
	}

	rule := p.rules[ref.name]
	if rule == nil {
		p.addErr(fmt.Errorf("undefined rule: %s", ref.name))
		return nil, false
	}
	return p.parseRule(rule)
}

func (p *parser) parseSeqExpr(seq *seqExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseSeqExpr"))
	}

	var vals []interface{}

	pt := p.pt
	for _, expr := range seq.exprs {
		val, ok := p.parseExpr(expr)
		if !ok {
			p.restore(pt)
			return nil, false
		}
		vals = append(vals, val)
	}
	return vals, true
}

func (p *parser) parseZeroOrMoreExpr(expr *zeroOrMoreExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseZeroOrMoreExpr"))
	}

	var vals []interface{}

	for {
		p.pushV()
		val, ok := p.parseExpr(expr.expr)
		p.popV()
		if !ok {
			return vals, true
		}
		vals = append(vals, val)
	}
}

func (p *parser) parseZeroOrOneExpr(expr *zeroOrOneExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseZeroOrOneExpr"))
	}

	p.pushV()
	val, _ := p.parseExpr(expr.expr)
	p.popV()
	// whether it matched or not, consider it a match
	return val, true
}

func rangeTable(class string) *unicode.RangeTable {
	if rt, ok := unicode.Categories[class]; ok {
		return rt
	}
	if rt, ok := unicode.Properties[class]; ok {
		return rt
	}
	if rt, ok := unicode.Scripts[class]; ok {
		return rt
	}

	// cannot happen
	panic(fmt.Sprintf("invalid Unicode class: %s", class))
}
