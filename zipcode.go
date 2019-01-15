// zipcode 将邮编范围整理为单个邮编序列
package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
)

var buf bytes.Buffer        // 输出缓冲
var separator = "-"         // 邮编起止分隔符
var cotnumber = 0           // 国家列编号
var colnumber = []int{1, 2} // 起止邮编列编号

func main() {
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		line := input.Text()
		secs := strings.Split(line, ",")
		if len(colnumber) > 1 {
			secs[colnumber[0]] = secs[colnumber[0]] + "-" + secs[colnumber[1]]
			secs[colnumber[1]] = ""
		}
		var before, after, zipstr, country string
		var zipnum1, zipnum2 int
		var err1, err2 error
		for i, s := range secs {
			if i == colnumber[0] {
				zipsecs := strings.Split(s, separator)
				zipstr, zipnum1, err1 = extract(zipsecs[0])
				if err1 != nil {
					zipstr = zipsecs[0]
				}
				_, zipnum2, err2 = extract(zipsecs[1])
				continue
			}
			if i == cotnumber {
				country = s
				continue
			}
			if i < colnumber[0] {
				before += s + ","
				continue
			}
			after += s + ","
		}
		if err1 != nil || err2 != nil {
			line := fmt.Sprintf("%s%s_%s,%s\n", before, country, zipstr, after)
			buf.WriteString(line)
			continue
		}
		for i := zipnum1; i <= zipnum2; i++ {
			line := fmt.Sprintf("%s%s_%s%d,%s\n", before, country, zipstr, i, after)
			buf.WriteString(line)
		}
	}
	if err := input.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "zipcode: %v\n", err)
		os.Exit(1)
	}
	fmt.Print(buf.String())
}

func extract(str string) (string, int, error) {
	var s, e = 0, 0
	var mks = true
	for i, r := range str {
		if unicode.IsNumber(r) {
			if mks && r != '0' {
				s = i
				mks = false
			}
			if !mks {
				e = i + 1
			}
		} else {
			if !mks {
				break
			}
		}
	}
	num, err := strconv.Atoi(string(str[s:e]))
	return string(str[:s]), num, err
}
