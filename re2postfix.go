// // rewrite infix regular expressionsinto equivalent postfix expressions
// // https://en.wikipedia.org/wiki/Shunting-yard_algorithm
// package main

// import "strings"

// // 操作符优先级
// var operatorPriority = map[byte]int{
// 	'*': 4,
// 	'?': 4,
// 	'+': 4,
// 	'.': 3,
// 	'|': 2,
// 	'(': 1, // TODO: 左右括号算操作符吗？没用的话去掉
// }

// var operator_stack *bytestack
// var postfix_result strings.Builder

// func re2postfix(re string) string {
// 	operator_stack = new_bytestack()
// 	postfix_result.Reset()

// 	should_add_concat := false

// 	for i := 0; i < len(re); i++ {
// 		ch := re[i]

// 		if ch == '*' || ch == '?' || ch == '+' {
// 			should_add_concat = true
// 			push_operator(ch)
// 			continue
// 		}

// 		if ch == '|' {
// 			should_add_concat = false
// 			push_operator(ch)
// 			continue
// 		}

// 		if ch == '(' {
// 			if should_add_concat {
// 				push_operator('.')
// 			}
// 			operator_stack.in(ch)
// 			should_add_concat = false
// 			continue
// 		}

// 		if ch == ')' {
// 			var operator byte

// 			for !operator_stack.is_empty() {
// 				operator = operator_stack.out()
// 				if operator == '(' {
// 					break
// 				}
// 				postfix_result.WriteByte(operator)
// 			}

// 			if operator != '(' {
// 				panic("unmatched ')'")
// 			}

// 			should_add_concat = true
// 			continue
// 		}

// 		if should_add_concat {
// 			push_operator('.')
// 		}

// 		postfix_result.WriteByte(ch)
// 		should_add_concat = true
// 	}

// 	for !operator_stack.is_empty() {
// 		operator := operator_stack.out()
// 		if operator == '(' {
// 			panic("unmatched '('")
// 		}
// 		postfix_result.WriteByte(operator)
// 	}

// 	return postfix_result.String()
// }

// func push_operator(operator byte) {
// 	current_priority := operatorPriority[operator]

// 	for !operator_stack.is_empty() {
// 		top := operator_stack.out()
// 		top_priority := operatorPriority[top]

// 		if top_priority >= current_priority {
// 			postfix_result.WriteByte(top)
// 		} else {
// 			operator_stack.in(top)
// 			break
// 		}
// 	}

// 	operator_stack.in(operator)
// }
