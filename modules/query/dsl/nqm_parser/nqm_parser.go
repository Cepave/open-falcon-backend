package nqm_parser

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
			name: "Query",
			pos:  position{line: 5, col: 1, offset: 24},
			expr: &actionExpr{
				pos: position{line: 5, col: 9, offset: 32},
				run: (*parser).callonQuery1,
				expr: &seqExpr{
					pos: position{line: 5, col: 9, offset: 32},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 5, col: 9, offset: 32},
							label: "allParams",
							expr: &zeroOrMoreExpr{
								pos: position{line: 5, col: 19, offset: 42},
								expr: &ruleRefExpr{
									pos:  position{line: 5, col: 19, offset: 42},
									name: "QueryParam",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 5, col: 31, offset: 54},
							name: "EOF",
						},
					},
				},
			},
		},
		{
			name: "QueryParam",
			pos:  position{line: 9, col: 1, offset: 97},
			expr: &choiceExpr{
				pos: position{line: 9, col: 14, offset: 110},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 9, col: 14, offset: 110},
						run: (*parser).callonQueryParam2,
						expr: &seqExpr{
							pos: position{line: 9, col: 14, offset: 110},
							exprs: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 9, col: 14, offset: 110},
									name: "_",
								},
								&labeledExpr{
									pos:   position{line: 9, col: 16, offset: 112},
									label: "param",
									expr: &choiceExpr{
										pos: position{line: 9, col: 23, offset: 119},
										alternatives: []interface{}{
											&ruleRefExpr{
												pos:  position{line: 9, col: 23, offset: 119},
												name: "TimeFilter",
											},
											&ruleRefExpr{
												pos:  position{line: 9, col: 36, offset: 132},
												name: "NodeFilter",
											},
										},
									},
								},
								&ruleRefExpr{
									pos:  position{line: 9, col: 48, offset: 144},
									name: "_",
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 11, col: 5, offset: 171},
						run: (*parser).callonQueryParam10,
						expr: &seqExpr{
							pos: position{line: 11, col: 5, offset: 171},
							exprs: []interface{}{
								&labeledExpr{
									pos:   position{line: 11, col: 5, offset: 171},
									label: "paramName",
									expr: &ruleRefExpr{
										pos:  position{line: 11, col: 15, offset: 181},
										name: "ParamName",
									},
								},
								&zeroOrOneExpr{
									pos: position{line: 11, col: 25, offset: 191},
									expr: &seqExpr{
										pos: position{line: 11, col: 26, offset: 192},
										exprs: []interface{}{
											&litMatcher{
												pos:        position{line: 11, col: 26, offset: 192},
												val:        "=",
												ignoreCase: false,
											},
											&zeroOrOneExpr{
												pos: position{line: 11, col: 30, offset: 196},
												expr: &ruleRefExpr{
													pos:  position{line: 11, col: 30, offset: 196},
													name: "ParamValue",
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
		},
		{
			name: "NodeFilter",
			pos:  position{line: 15, col: 1, offset: 291},
			expr: &choiceExpr{
				pos: position{line: 15, col: 14, offset: 304},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 15, col: 14, offset: 304},
						run: (*parser).callonNodeFilter2,
						expr: &seqExpr{
							pos: position{line: 15, col: 14, offset: 304},
							exprs: []interface{}{
								&labeledExpr{
									pos:   position{line: 15, col: 14, offset: 304},
									label: "nodeProperty",
									expr: &ruleRefExpr{
										pos:  position{line: 15, col: 27, offset: 317},
										name: "NodeProperty",
									},
								},
								&litMatcher{
									pos:        position{line: 15, col: 40, offset: 330},
									val:        "=",
									ignoreCase: false,
								},
								&labeledExpr{
									pos:   position{line: 15, col: 44, offset: 334},
									label: "propertyValue",
									expr: &ruleRefExpr{
										pos:  position{line: 15, col: 58, offset: 348},
										name: "MultiLiteralString",
									},
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 17, col: 5, offset: 426},
						run: (*parser).callonNodeFilter9,
						expr: &seqExpr{
							pos: position{line: 17, col: 5, offset: 426},
							exprs: []interface{}{
								&labeledExpr{
									pos:   position{line: 17, col: 5, offset: 426},
									label: "paramName",
									expr: &ruleRefExpr{
										pos:  position{line: 17, col: 15, offset: 436},
										name: "NodeProperty",
									},
								},
								&labeledExpr{
									pos:   position{line: 17, col: 28, offset: 449},
									label: "assignedValue",
									expr: &zeroOrOneExpr{
										pos: position{line: 17, col: 42, offset: 463},
										expr: &seqExpr{
											pos: position{line: 17, col: 43, offset: 464},
											exprs: []interface{}{
												&litMatcher{
													pos:        position{line: 17, col: 43, offset: 464},
													val:        "=",
													ignoreCase: false,
												},
												&ruleRefExpr{
													pos:  position{line: 17, col: 47, offset: 468},
													name: "ParamValue",
												},
											},
										},
									},
								},
								&ruleRefExpr{
									pos:  position{line: 17, col: 60, offset: 481},
									name: "END_WORD",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "TimeFilter",
			pos:  position{line: 21, col: 1, offset: 575},
			expr: &choiceExpr{
				pos: position{line: 21, col: 14, offset: 588},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 21, col: 14, offset: 588},
						run: (*parser).callonTimeFilter2,
						expr: &seqExpr{
							pos: position{line: 21, col: 14, offset: 588},
							exprs: []interface{}{
								&labeledExpr{
									pos:   position{line: 21, col: 14, offset: 588},
									label: "timeParamName",
									expr: &ruleRefExpr{
										pos:  position{line: 21, col: 28, offset: 602},
										name: "TimeParamName",
									},
								},
								&litMatcher{
									pos:        position{line: 21, col: 42, offset: 616},
									val:        "=",
									ignoreCase: false,
								},
								&labeledExpr{
									pos:   position{line: 21, col: 46, offset: 620},
									label: "timeValue",
									expr: &choiceExpr{
										pos: position{line: 21, col: 57, offset: 631},
										alternatives: []interface{}{
											&ruleRefExpr{
												pos:  position{line: 21, col: 57, offset: 631},
												name: "ISO_8601",
											},
											&ruleRefExpr{
												pos:  position{line: 21, col: 68, offset: 642},
												name: "UNIX_TIME",
											},
										},
									},
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 23, col: 5, offset: 709},
						run: (*parser).callonTimeFilter11,
						expr: &seqExpr{
							pos: position{line: 23, col: 5, offset: 709},
							exprs: []interface{}{
								&labeledExpr{
									pos:   position{line: 23, col: 5, offset: 709},
									label: "paramName",
									expr: &ruleRefExpr{
										pos:  position{line: 23, col: 15, offset: 719},
										name: "TimeParamName",
									},
								},
								&labeledExpr{
									pos:   position{line: 23, col: 29, offset: 733},
									label: "assignedValue",
									expr: &zeroOrOneExpr{
										pos: position{line: 23, col: 43, offset: 747},
										expr: &seqExpr{
											pos: position{line: 23, col: 44, offset: 748},
											exprs: []interface{}{
												&litMatcher{
													pos:        position{line: 23, col: 44, offset: 748},
													val:        "=",
													ignoreCase: false,
												},
												&zeroOrOneExpr{
													pos: position{line: 23, col: 48, offset: 752},
													expr: &ruleRefExpr{
														pos:  position{line: 23, col: 48, offset: 752},
														name: "ParamValue",
													},
												},
											},
										},
									},
								},
								&ruleRefExpr{
									pos:  position{line: 23, col: 62, offset: 766},
									name: "END_WORD",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "NodeProperty",
			pos:  position{line: 27, col: 1, offset: 860},
			expr: &actionExpr{
				pos: position{line: 27, col: 16, offset: 875},
				run: (*parser).callonNodeProperty1,
				expr: &seqExpr{
					pos: position{line: 27, col: 16, offset: 875},
					exprs: []interface{}{
						&choiceExpr{
							pos: position{line: 27, col: 17, offset: 876},
							alternatives: []interface{}{
								&litMatcher{
									pos:        position{line: 27, col: 17, offset: 876},
									val:        "agent",
									ignoreCase: false,
								},
								&litMatcher{
									pos:        position{line: 27, col: 27, offset: 886},
									val:        "target",
									ignoreCase: false,
								},
							},
						},
						&litMatcher{
							pos:        position{line: 27, col: 37, offset: 896},
							val:        ".",
							ignoreCase: false,
						},
						&choiceExpr{
							pos: position{line: 27, col: 42, offset: 901},
							alternatives: []interface{}{
								&litMatcher{
									pos:        position{line: 27, col: 42, offset: 901},
									val:        "isp",
									ignoreCase: false,
								},
								&litMatcher{
									pos:        position{line: 27, col: 50, offset: 909},
									val:        "province",
									ignoreCase: false,
								},
								&litMatcher{
									pos:        position{line: 27, col: 63, offset: 922},
									val:        "city",
									ignoreCase: false,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "TimeParamName",
			pos:  position{line: 31, col: 1, offset: 963},
			expr: &actionExpr{
				pos: position{line: 31, col: 17, offset: 979},
				run: (*parser).callonTimeParamName1,
				expr: &choiceExpr{
					pos: position{line: 31, col: 18, offset: 980},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 31, col: 18, offset: 980},
							val:        "starttime",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 31, col: 32, offset: 994},
							val:        "endtime",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "ParamName",
			pos:  position{line: 35, col: 1, offset: 1038},
			expr: &actionExpr{
				pos: position{line: 35, col: 13, offset: 1050},
				run: (*parser).callonParamName1,
				expr: &oneOrMoreExpr{
					pos: position{line: 35, col: 13, offset: 1050},
					expr: &charClassMatcher{
						pos:        position{line: 35, col: 13, offset: 1050},
						val:        "[^ =\\t\\n\\r]",
						chars:      []rune{' ', '=', '\t', '\n', '\r'},
						ignoreCase: false,
						inverted:   true,
					},
				},
			},
		},
		{
			name: "ParamValue",
			pos:  position{line: 38, col: 1, offset: 1095},
			expr: &actionExpr{
				pos: position{line: 38, col: 14, offset: 1108},
				run: (*parser).callonParamValue1,
				expr: &oneOrMoreExpr{
					pos: position{line: 38, col: 14, offset: 1108},
					expr: &charClassMatcher{
						pos:        position{line: 38, col: 14, offset: 1108},
						val:        "[^ \\t\\n\\r]",
						chars:      []rune{' ', '\t', '\n', '\r'},
						ignoreCase: false,
						inverted:   true,
					},
				},
			},
		},
		{
			name: "MultiLiteralString",
			pos:  position{line: 42, col: 1, offset: 1153},
			expr: &choiceExpr{
				pos: position{line: 42, col: 22, offset: 1174},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 42, col: 22, offset: 1174},
						run: (*parser).callonMultiLiteralString2,
						expr: &seqExpr{
							pos: position{line: 42, col: 22, offset: 1174},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 42, col: 22, offset: 1174},
									val:        "%",
									ignoreCase: false,
								},
								&labeledExpr{
									pos:   position{line: 42, col: 26, offset: 1178},
									label: "autoCondition",
									expr: &ruleRefExpr{
										pos:  position{line: 42, col: 40, offset: 1192},
										name: "LiteralStringForAutoCond",
									},
								},
								&litMatcher{
									pos:        position{line: 42, col: 65, offset: 1217},
									val:        "%",
									ignoreCase: false,
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 44, col: 5, offset: 1269},
						run: (*parser).callonMultiLiteralString8,
						expr: &seqExpr{
							pos: position{line: 44, col: 5, offset: 1269},
							exprs: []interface{}{
								&labeledExpr{
									pos:   position{line: 44, col: 5, offset: 1269},
									label: "first",
									expr: &ruleRefExpr{
										pos:  position{line: 44, col: 11, offset: 1275},
										name: "LiteralString",
									},
								},
								&labeledExpr{
									pos:   position{line: 44, col: 25, offset: 1289},
									label: "rest",
									expr: &zeroOrMoreExpr{
										pos: position{line: 44, col: 30, offset: 1294},
										expr: &ruleRefExpr{
											pos:  position{line: 44, col: 31, offset: 1295},
											name: "RestLiteralString",
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
			name: "LiteralStringForAutoCond",
			pos:  position{line: 48, col: 1, offset: 1368},
			expr: &actionExpr{
				pos: position{line: 48, col: 28, offset: 1395},
				run: (*parser).callonLiteralStringForAutoCond1,
				expr: &oneOrMoreExpr{
					pos: position{line: 48, col: 28, offset: 1395},
					expr: &charClassMatcher{
						pos:        position{line: 48, col: 28, offset: 1395},
						val:        "[^ \\t\\n\\r,=%]",
						chars:      []rune{' ', '\t', '\n', '\r', ',', '=', '%'},
						ignoreCase: false,
						inverted:   true,
					},
				},
			},
		},
		{
			name: "LiteralString",
			pos:  position{line: 51, col: 1, offset: 1442},
			expr: &actionExpr{
				pos: position{line: 51, col: 17, offset: 1458},
				run: (*parser).callonLiteralString1,
				expr: &oneOrMoreExpr{
					pos: position{line: 51, col: 17, offset: 1458},
					expr: &charClassMatcher{
						pos:        position{line: 51, col: 17, offset: 1458},
						val:        "[^ \\t\\n\\r,=]",
						chars:      []rune{' ', '\t', '\n', '\r', ',', '='},
						ignoreCase: false,
						inverted:   true,
					},
				},
			},
		},
		{
			name: "RestLiteralString",
			pos:  position{line: 55, col: 1, offset: 1505},
			expr: &choiceExpr{
				pos: position{line: 55, col: 21, offset: 1525},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 55, col: 21, offset: 1525},
						run: (*parser).callonRestLiteralString2,
						expr: &seqExpr{
							pos: position{line: 55, col: 21, offset: 1525},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 55, col: 21, offset: 1525},
									val:        ",",
									ignoreCase: false,
								},
								&labeledExpr{
									pos:   position{line: 55, col: 25, offset: 1529},
									label: "sValue",
									expr: &ruleRefExpr{
										pos:  position{line: 55, col: 32, offset: 1536},
										name: "LiteralString",
									},
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 57, col: 5, offset: 1576},
						run: (*parser).callonRestLiteralString7,
						expr: &seqExpr{
							pos: position{line: 57, col: 5, offset: 1576},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 57, col: 5, offset: 1576},
									val:        ",",
									ignoreCase: false,
								},
								&labeledExpr{
									pos:   position{line: 57, col: 9, offset: 1580},
									label: "errorLiteralValue",
									expr: &ruleRefExpr{
										pos:  position{line: 57, col: 27, offset: 1598},
										name: "ParamValue",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "ISO_8601",
			pos:  position{line: 61, col: 1, offset: 1689},
			expr: &actionExpr{
				pos: position{line: 61, col: 12, offset: 1700},
				run: (*parser).callonISO_86011,
				expr: &seqExpr{
					pos: position{line: 61, col: 12, offset: 1700},
					exprs: []interface{}{
						&oneOrMoreExpr{
							pos: position{line: 61, col: 12, offset: 1700},
							expr: &charClassMatcher{
								pos:        position{line: 61, col: 12, offset: 1700},
								val:        "[0-9]",
								ranges:     []rune{'0', '9'},
								ignoreCase: false,
								inverted:   false,
							},
						},
						&litMatcher{
							pos:        position{line: 61, col: 19, offset: 1707},
							val:        "-",
							ignoreCase: false,
						},
						&oneOrMoreExpr{
							pos: position{line: 61, col: 23, offset: 1711},
							expr: &charClassMatcher{
								pos:        position{line: 61, col: 23, offset: 1711},
								val:        "[0-9]",
								ranges:     []rune{'0', '9'},
								ignoreCase: false,
								inverted:   false,
							},
						},
						&litMatcher{
							pos:        position{line: 61, col: 30, offset: 1718},
							val:        "-",
							ignoreCase: false,
						},
						&oneOrMoreExpr{
							pos: position{line: 61, col: 34, offset: 1722},
							expr: &charClassMatcher{
								pos:        position{line: 61, col: 34, offset: 1722},
								val:        "[0-9]",
								ranges:     []rune{'0', '9'},
								ignoreCase: false,
								inverted:   false,
							},
						},
						&zeroOrOneExpr{
							pos: position{line: 61, col: 41, offset: 1729},
							expr: &seqExpr{
								pos: position{line: 61, col: 42, offset: 1730},
								exprs: []interface{}{
									&litMatcher{
										pos:        position{line: 61, col: 42, offset: 1730},
										val:        "T",
										ignoreCase: false,
									},
									&oneOrMoreExpr{
										pos: position{line: 61, col: 46, offset: 1734},
										expr: &charClassMatcher{
											pos:        position{line: 61, col: 46, offset: 1734},
											val:        "[0-9]",
											ranges:     []rune{'0', '9'},
											ignoreCase: false,
											inverted:   false,
										},
									},
									&zeroOrOneExpr{
										pos: position{line: 61, col: 53, offset: 1741},
										expr: &seqExpr{
											pos: position{line: 61, col: 54, offset: 1742},
											exprs: []interface{}{
												&litMatcher{
													pos:        position{line: 61, col: 54, offset: 1742},
													val:        ":",
													ignoreCase: false,
												},
												&oneOrMoreExpr{
													pos: position{line: 61, col: 58, offset: 1746},
													expr: &charClassMatcher{
														pos:        position{line: 61, col: 58, offset: 1746},
														val:        "[0-9]",
														ranges:     []rune{'0', '9'},
														ignoreCase: false,
														inverted:   false,
													},
												},
												&zeroOrOneExpr{
													pos: position{line: 61, col: 65, offset: 1753},
													expr: &seqExpr{
														pos: position{line: 61, col: 66, offset: 1754},
														exprs: []interface{}{
															&zeroOrOneExpr{
																pos: position{line: 61, col: 66, offset: 1754},
																expr: &charClassMatcher{
																	pos:        position{line: 61, col: 66, offset: 1754},
																	val:        "[Z+-]",
																	chars:      []rune{'Z', '+', '-'},
																	ignoreCase: false,
																	inverted:   false,
																},
															},
															&oneOrMoreExpr{
																pos: position{line: 61, col: 73, offset: 1761},
																expr: &charClassMatcher{
																	pos:        position{line: 61, col: 73, offset: 1761},
																	val:        "[0-9]",
																	ranges:     []rune{'0', '9'},
																	ignoreCase: false,
																	inverted:   false,
																},
															},
															&zeroOrOneExpr{
																pos: position{line: 61, col: 80, offset: 1768},
																expr: &seqExpr{
																	pos: position{line: 61, col: 81, offset: 1769},
																	exprs: []interface{}{
																		&litMatcher{
																			pos:        position{line: 61, col: 81, offset: 1769},
																			val:        ":",
																			ignoreCase: false,
																		},
																		&oneOrMoreExpr{
																			pos: position{line: 61, col: 85, offset: 1773},
																			expr: &charClassMatcher{
																				pos:        position{line: 61, col: 85, offset: 1773},
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
			name: "UNIX_TIME",
			pos:  position{line: 65, col: 1, offset: 1817},
			expr: &actionExpr{
				pos: position{line: 65, col: 13, offset: 1829},
				run: (*parser).callonUNIX_TIME1,
				expr: &oneOrMoreExpr{
					pos: position{line: 65, col: 13, offset: 1829},
					expr: &charClassMatcher{
						pos:        position{line: 65, col: 13, offset: 1829},
						val:        "[0-9]",
						ranges:     []rune{'0', '9'},
						ignoreCase: false,
						inverted:   false,
					},
				},
			},
		},
		{
			name: "_",
			pos:  position{line: 69, col: 1, offset: 1866},
			expr: &zeroOrMoreExpr{
				pos: position{line: 69, col: 5, offset: 1870},
				expr: &ruleRefExpr{
					pos:  position{line: 69, col: 5, offset: 1870},
					name: "EMPTY_CHAR",
				},
			},
		},
		{
			name: "END_WORD",
			pos:  position{line: 70, col: 1, offset: 1882},
			expr: &choiceExpr{
				pos: position{line: 70, col: 12, offset: 1893},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 70, col: 12, offset: 1893},
						name: "EOF",
					},
					&oneOrMoreExpr{
						pos: position{line: 70, col: 18, offset: 1899},
						expr: &ruleRefExpr{
							pos:  position{line: 70, col: 18, offset: 1899},
							name: "EMPTY_CHAR",
						},
					},
				},
			},
		},
		{
			name: "EMPTY_CHAR",
			pos:  position{line: 71, col: 1, offset: 1911},
			expr: &charClassMatcher{
				pos:        position{line: 71, col: 14, offset: 1924},
				val:        "[ \\t\\n\\r]",
				chars:      []rune{' ', '\t', '\n', '\r'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "EOF",
			pos:  position{line: 72, col: 1, offset: 1934},
			expr: &notExpr{
				pos: position{line: 72, col: 7, offset: 1940},
				expr: &anyMatcher{
					line: 72, col: 8, offset: 1941,
				},
			},
		},
	},
}

func (c *current) onQuery1(allParams interface{}) (interface{}, error) {
	return toSetters(allParams), nil
}

func (p *parser) callonQuery1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onQuery1(stack["allParams"])
}

func (c *current) onQueryParam2(param interface{}) (interface{}, error) {
	return param, nil
}

func (p *parser) callonQueryParam2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onQueryParam2(stack["param"])
}

func (c *current) onQueryParam10(paramName interface{}) (interface{}, error) {
	return empty_param_setter, fmt.Errorf("Unknown parameter: %q", paramName)
}

func (p *parser) callonQueryParam10() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onQueryParam10(stack["paramName"])
}

func (c *current) onNodeFilter2(nodeProperty, propertyValue interface{}) (interface{}, error) {
	return buildSetterFunc(nodeProperty, propertyValue)
}

func (p *parser) callonNodeFilter2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onNodeFilter2(stack["nodeProperty"], stack["propertyValue"])
}

func (c *current) onNodeFilter9(paramName, assignedValue interface{}) (interface{}, error) {
	return empty_param_setter, buildErrorForInvalidParam(paramName, assignedValue)
}

func (p *parser) callonNodeFilter9() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onNodeFilter9(stack["paramName"], stack["assignedValue"])
}

func (c *current) onTimeFilter2(timeParamName, timeValue interface{}) (interface{}, error) {
	return buildSetterFunc(timeParamName, timeValue)
}

func (p *parser) callonTimeFilter2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onTimeFilter2(stack["timeParamName"], stack["timeValue"])
}

func (c *current) onTimeFilter11(paramName, assignedValue interface{}) (interface{}, error) {
	return empty_param_setter, buildErrorForInvalidParam(paramName, assignedValue)
}

func (p *parser) callonTimeFilter11() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onTimeFilter11(stack["paramName"], stack["assignedValue"])
}

func (c *current) onNodeProperty1() (interface{}, error) {
	return string(c.text), nil
}

func (p *parser) callonNodeProperty1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onNodeProperty1()
}

func (c *current) onTimeParamName1() (interface{}, error) {
	return string(c.text), nil
}

func (p *parser) callonTimeParamName1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onTimeParamName1()
}

func (c *current) onParamName1() (interface{}, error) {
	return string(c.text), nil
}

func (p *parser) callonParamName1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onParamName1()
}

func (c *current) onParamValue1() (interface{}, error) {
	return string(c.text), nil
}

func (p *parser) callonParamValue1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onParamValue1()
}

func (c *current) onMultiLiteralString2(autoCondition interface{}) (interface{}, error) {
	return parseAutoCondition(autoCondition)
}

func (p *parser) callonMultiLiteralString2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onMultiLiteralString2(stack["autoCondition"])
}

func (c *current) onMultiLiteralString8(first, rest interface{}) (interface{}, error) {
	return combineStringLiterals(first, rest), nil
}

func (p *parser) callonMultiLiteralString8() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onMultiLiteralString8(stack["first"], stack["rest"])
}

func (c *current) onLiteralStringForAutoCond1() (interface{}, error) {
	return string(c.text), nil
}

func (p *parser) callonLiteralStringForAutoCond1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onLiteralStringForAutoCond1()
}

func (c *current) onLiteralString1() (interface{}, error) {
	return string(c.text), nil
}

func (p *parser) callonLiteralString1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onLiteralString1()
}

func (c *current) onRestLiteralString2(sValue interface{}) (interface{}, error) {
	return sValue, nil
}

func (p *parser) callonRestLiteralString2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onRestLiteralString2(stack["sValue"])
}

func (c *current) onRestLiteralString7(errorLiteralValue interface{}) (interface{}, error) {
	return "", fmt.Errorf("Illegal literal value: \"%v\"", errorLiteralValue)
}

func (p *parser) callonRestLiteralString7() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onRestLiteralString7(stack["errorLiteralValue"])
}

func (c *current) onISO_86011() (interface{}, error) {
	return parseIso8601(c)
}

func (p *parser) callonISO_86011() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onISO_86011()
}

func (c *current) onUNIX_TIME1() (interface{}, error) {
	return parseUnixTime(c)
}

func (p *parser) callonUNIX_TIME1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onUNIX_TIME1()
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
