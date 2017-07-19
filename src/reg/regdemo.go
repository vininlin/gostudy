package main

import (
	"regexp"
	"fmt"
)

func main()  {
	text := "Hello 世界！123 Go."

	reg := regexp.MustCompile(`[a-z]+`)
	fmt.Printf("1.%q\n", reg.FindAllString(text, -1))

	reg = regexp.MustCompile(`[^a-z]+`)
	fmt.Printf("2.%q\n", reg.FindAllString(text, -1))

	reg = regexp.MustCompile(`[\w]+`)
	fmt.Printf("3.%q\n", reg.FindAllString(text, -1))

	reg = regexp.MustCompile(`[^\w\s]+`)
	fmt.Printf("4.%q\n", reg.FindAllString(text, -1))

	reg = regexp.MustCompile(`[[:upper:]]]+`)
	fmt.Printf("5.%q\n", reg.FindAllString(text, -1))

	reg = regexp.MustCompile(`[[:^ascii:]]+`)
	fmt.Printf("6.%q\n", reg.FindAllString(text, -1))

	reg = regexp.MustCompile(`[\pP]+`)
	fmt.Printf("7.%q\n", reg.FindAllString(text, -1))

	reg = regexp.MustCompile(`[\PP]+`)
	fmt.Printf("8.%q\n", reg.FindAllString(text, -1))

	reg = regexp.MustCompile(`[\p{Han}]+`)
	fmt.Printf("9。%q\n", reg.FindAllString(text, -1))

	reg = regexp.MustCompile(`[\P{Han}]+`)
	fmt.Printf("10.%q\n", reg.FindAllString(text, -1))

	reg = regexp.MustCompile(`Hello|Go`)
	fmt.Printf("11.%q\n", reg.FindAllString(text, -1))

	// 查找行首以 H 开头，以空格结尾的字符串
	reg = regexp.MustCompile(`^H.*\s`)
	fmt.Printf("12.%q\n", reg.FindAllString(text, -1))

	//查找行首以 H 开头，以空白结尾的字符串（非贪婪模式）
	reg = regexp.MustCompile(`(?U)^H.*\s`)
	fmt.Printf("13.%q\n", reg.FindAllString(text, -1))

	// 查找以 hello 开头（忽略大小写），以 Go 结尾的字符串
	reg = regexp.MustCompile(`(?i:^hello).*Go`)
	fmt.Printf("14.%q\n", reg.FindAllString(text, -1))

	//查找 Go.
	reg = regexp.MustCompile(`\QGo.\E`)
	fmt.Printf("15.%q\n", reg.FindAllString(text, -1))

	// 查找从行首开始，以空格结尾的字符串（非贪婪模式）
	reg = regexp.MustCompile(`(?U)^.*`)
	fmt.Printf("16.%q\n", reg.FindAllString(text, -1))

	// 查找以空格开头，到行尾结束，中间不包含空格字符串
	reg = regexp.MustCompile(` [^ ]*$`)
	fmt.Printf("17.%q\n", reg.FindAllString(text, -1))

	// 查找“单词边界”之间的字符串
	reg = regexp.MustCompile(`(?U)\b.+\b`)
	fmt.Printf("18.%q\n", reg.FindAllString(text, -1))

	// 查找连续 1 次到 4 次的非空格字符，并以 o 结尾的字符串
	reg = regexp.MustCompile(`[^ ]{1,4}o`)
	fmt.Printf("19.%q\n", reg.FindAllString(text, -1))

	// 查找 Hello 或 Go
	reg = regexp.MustCompile(`(?:Hell|G)o`)
	fmt.Printf("20.%q\n", reg.FindAllString(text, -1))

	// 查找 Hello 或 Go，替换为 Hellooo、Gooo
	reg = regexp.MustCompile(`(?P<n>Hell|G)o`)
	fmt.Printf("21.%q\n", reg.ReplaceAllString(text, "${n}ooo"))

	// 交换 Hello 和 Go
	reg = regexp.MustCompile(`(Hello)(.*)(Go)`)
	fmt.Printf("22.%q\n", reg.ReplaceAllString(text, "$3$2$1"))

	// 特殊字符的查找
	reg = regexp.MustCompile(`[\f\t\n\r\v\123\x7F\x{10FFFF}\\\^\$\.\*\+\?\{\}\(\)\[\]|]`)
	fmt.Printf("23.%q\n", reg.ReplaceAllString("\f\t\n\r\v\123\x7F\U0010FFFF\\^$.*+?{}()[]|", "-"))
}
