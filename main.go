package main

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type Node struct {
	val   string
	left  *Node
	right *Node
}

const LEFTBRACKET uint8 = '('
const RIGHTBRACKET uint8 = ')'
const AND uint8 = '&'
const OR uint8 = '|'

func BuildTree(expression string) (n *Node, e error) {
	root := &Node{}

	s1 := []*Node{}
	s2 := []uint8{}
	s3 := []int{}
	for i := 0; i < len(expression); {
		if len(s3) > 0 && expression[i] != LEFTBRACKET && expression[i] != RIGHTBRACKET {
			i++
			continue
		} else if expression[i] == LEFTBRACKET {
			s3 = append(s3, i)
			i++
			continue
		} else if expression[i] == RIGHTBRACKET {
			if len(s3) == 0 {
				return nil, errors.New("invalid expression,please check the number of brackets")
			}
			var start int
			start, s3 = s3[len(s3)-1], s3[:len(s3)-1]
			if len(s3) == 0 {
				if len(s1) == 2 {
					val, err := BuildTree(expression[start+1 : i])
					if err != nil {
						return nil, err
					}
					if s2[0] == AND || (s2[0] == s2[1]) {
						var v1, v2 *Node
						v2, s1 = s1[len(s1)-1], s1[:len(s1)-1]
						v1, s1 = s1[len(s1)-1], s1[:len(s1)-1]
						switch s2[0] {
						case OR:
							root := &Node{}
							root.val = string(OR)
							root.left = v1
							root.right = v2
							s1 = append(s1, root)
						case AND:
							root := &Node{}
							root.val = string(AND)
							root.left = v1
							root.right = v2
							s1 = append(s1, root)
						}
						s2 = s2[1:]
						s1 = append(s1, val)
					} else {
						var v1 *Node
						v1, s1 = s1[len(s1)-1], s1[:len(s1)-1]
						s2 = s2[:len(s2)-1]
						root := &Node{}
						root.val = string(AND)
						root.left = v1
						root.right = val
						s1 = append(s1, root)
					}
				} else {
					val, err := BuildTree(expression[start+1 : i])
					if err != nil {
						return nil, err
					}
					s1 = append(s1, val)
				}
			}
			i++
			continue
		} else if (expression[i] == AND) || (expression[i] == OR) {
			s2 = append(s2, expression[i])
			i++
			continue
		} else {
			j := i + 1
			for j < len(expression) && expression[j] != AND && expression[j] != OR && expression[j] != LEFTBRACKET && expression[j] != RIGHTBRACKET {
				j++
			}
			var v *Node
			if IsValid(expression[i:j]) {
				v = &Node{val: expression[i:j]}
			} else {
				return nil, errors.New(fmt.Sprintf("bad label expression:%s", expression[i:j]))
			}
			if len(s1) == 2 {
				if s2[0] == AND || s2[0] == s2[1] {
					var v2 *Node
					var v1 *Node
					v2, s1 = s1[len(s1)-1], s1[:len(s1)-1]
					v1, s1 = s1[len(s1)-1], s1[:len(s1)-1]
					switch s2[0] {
					case AND:
						root := &Node{}
						root.val = string(AND)
						root.left = v1
						root.right = v2
						s1 = append(s1, root)
					case OR:
						root := &Node{}
						root.val = string(OR)
						root.left = v1
						root.right = v2
						s1 = append(s1, root)
					}
					s2 = s2[1:]
					s1 = append(s1, v)
				} else {
					var v1 *Node
					v1, s1 = s1[len(s1)-1], s1[:len(s1)-1]
					s2 = s2[:len(s2)-1]
					root := &Node{}
					root.val = string(AND)
					root.left = v1
					root.right = v
					s1 = append(s1, root)
				}
			} else {
				s1 = append(s1, v)
			}
			i = j
		}
	}
	if len(s3) != 0 {
		return nil, errors.New("invalid expression,please check the number of brackets")
	}
	if len(s2) == 0 {
		if len(s1) == 0 {
			return nil, errors.New("invalid expression(null expression)")
		} else {
			return s1[0], nil
		}
	} else if len(s2) > 1 {
		return nil, errors.New("wrong expression,check the operator please")
	} else if s2[0] == AND {
		if len(s1) != 2 {
			return nil, errors.New("wrong expression,missing operand")
		}
		root.val = string(AND)
		root.left = s1[0]
		root.right = s1[1]
		return root, nil
	} else {
		if len(s1) != 2 {
			return nil, errors.New("wrong expression,missing operand")
		}
		root.val = string(OR)
		root.left = s1[0]
		root.right = s1[1]
		return root, nil
	}
}

func IsValid(s string) bool {
	res, _ := regexp.Match(`^[^!,^=,^ ,^\t]+!{0,1}=[^!,^=,^ ,^\t]+$`, []byte(s))
	return res
}

func PostorderTraversal(root *Node, list *[]string) {
	var l string
	var r string
	if root.left != nil {
		PostorderTraversal(root.left, list)
	}
	if root.right != nil {
		PostorderTraversal(root.right, list)
	}
	if l != "" {
		*list = append(*list, l)
	}
	if r != "" {
		*list = append(*list, r)
	}
	*list = append(*list, root.val)
}

func Converse2ReversePolishNotation(root *Node) string {
	const space uint8 = ' '
	list := &[]string{}
	PostorderTraversal(root, list)
	return strings.Join(*list, " ")
}

func Evaluation(labelMap map[string]string, expression string) bool {
	var stack []bool
	exp := strings.Split(expression, ` `)
	for i := 0; i < len(exp); i++ {
		if exp[i] == "&" {
			var v1, v2 bool
			v2, stack = stack[len(stack)-1], stack[:len(stack)-1]
			v1, stack = stack[len(stack)-1], stack[:len(stack)-1]
			stack = append(stack, v1 && v2)
		} else if exp[i] == "|" {
			var v1, v2 bool
			v2, stack = stack[len(stack)-1], stack[:len(stack)-1]
			v1, stack = stack[len(stack)-1], stack[:len(stack)-1]
			stack = append(stack, v1 || v2)
		} else {
			if strings.Contains(exp[i], "!") {
				val := strings.Split(exp[i], "!=")
				if _, ok := labelMap[val[0]]; ok {
					if labelMap[val[0]] == val[1] {
						stack = append(stack, false)
					} else {
						stack = append(stack, true)
					}
				} else {
					return false
				}
			} else {
				val := strings.Split(exp[i], "=")
				if _, ok := labelMap[val[0]]; ok {
					if labelMap[val[0]] == val[1] {
						stack = append(stack, true)
					} else {
						stack = append(stack, false)
					}
				} else {
					return false
				}
			}
		}
	}
	return stack[len(stack)-1]
}

func main() {
	expression := "(idc=beijing|app=mac)&env=product"
	tree, err := BuildTree(expression)
	if err != nil {

	}
	reverse := Converse2ReversePolishNotation(tree)
	fmt.Println(tree)
	fmt.Println(reverse)

	lables := make(map[string]string)
	lables["idc"] = "beijing"
	lables["env"] = "product"
	lables["app"] = "x"

	notation := Evaluation(lables, reverse)

	fmt.Println(notation)
}
