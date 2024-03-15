package main

import (
	"fmt"
)

func main() {

	s := `xxx)version/**/()/**/1=1 or 1=1 or name="(xx x)" and 1=sleep(1) or age=(1 or 1=1) or 1=1/* 1 or 1=1)`
	fmt.Println(l(r(s, ')')))
}

func ReplaceCharAt(s string, index int, c rune) string {
	if index < 0 || index >= len(s) {
		return s
	}
	// 创建新的字符串：原始字符串的前半部分 + 新字符 + 原始字符串的后半部分
	return s[:index] + string(c) + s[index+1:]
}

func RemoveCharAt(s string, index int) string {
	// 检查索引是否有效
	if index < 0 || index >= len(s) {
		return s
	}

	// 将字符串转换为rune切片，以便正确处理多字节字符
	runes := []rune(s)

	// 移除指定位置的字符
	runes = append(runes[:index], runes[index+1:]...)

	// 将rune切片转换回字符串
	return string(runes)
}

func l(s string) string {
	for i := 0; i < len(s); i++ {
		ch := s[i]
		if ch == '(' {
			if i == len(s)-1 {
				s = RemoveCharAt(s, i)
				break
			}
			found := false
			for j := i + 1; j <= len(s)-1; j++ {
				ch1 := s[j]
				//// 如果发现双引号， 成对双引号内的"("不算
				if ch1 == '"' || ch1 == '\'' {
					found = true
				}
				// 找到说明闭合了
				if !found && ch1 == ')' {
					break
				}
				if ch1 == '(' {
					s = RemoveCharAt(s, i)
				}
			}
		}
	}
	return s
}

func r(s string, sym rune) string {
	for i := len(s) - 1; i >= 0; i-- {
		ch := s[i]
		if ch == uint8(sym) {
			found := false
			if i == 0 {
				s = RemoveCharAt(s, i)
				break
			}
			for j := i - 1; j > -1; j-- {
				ch1 := s[j]
				// 如果发现双引号， 成对双引号内的"("不算
				if ch1 == '"' || ch1 == '\'' {
					found = true
				}

				// 找到说明闭合了
				if !found && ch1 == '(' {
					// 检查中间是否为空
					fmt.Println("=", string(s[j]))
					break
				}

				if ch1 == uint8(sym) {
					s = RemoveCharAt(s, i) //"T" + s[i+1:]
					//return r(s)
				}

				if j == 0 && !found {
					s = RemoveCharAt(s, i)
				}
			}
		}
	}
	return s
}
