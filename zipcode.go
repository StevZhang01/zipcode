// zipcode 将邮编范围整理为单个邮编序列
package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"unicode"
)

var buf bytes.Buffer        // 输出缓冲
var separator = "-"         // 邮编起止分隔符
var cotnumber = 2           // 国家列编号
var colnumber = []int{3, 4} // 起止邮编列编号
var title = true            //是否有表头

func main() {
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		// 跳过空行
		line := input.Text()
		if line == "" {
			continue
		}
		// 解析表头
		if title {
			buf.WriteString(line + "\n")
			title = false
			continue
		}
		secs := strings.Split(line, ",")
		if len(colnumber) > 1 {
			secs[colnumber[0]] = secs[colnumber[0]] + "-" + secs[colnumber[1]]
			secs[colnumber[1]] = ""
		}
		var before, after, zipstr, zipnum1, zipnum2, country string
		for i, s := range secs {
			if i == colnumber[0] {
				zipsecs := strings.Split(s, separator)
				zipstr, zipnum1 = extract(zipsecs[0])
				_, zipnum2 = extract(zipsecs[1])
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
		if zipnum1 == "" || zipnum2 == "" {
			line := fmt.Sprintf("%s,%s_%s, %s\n", before, country, zipstr, after)
			buf.WriteString(line)
			continue
		}
		l := len(zipnum1)
		if len(zipnum2) > len(zipnum1) {
			l = len(zipnum2)
		}
		num1, _ := strconv.Atoi(zipnum1)
		num2, _ := strconv.Atoi(zipnum2)
		l1 := len(fmt.Sprintf("%d", num1))
		l2 := len(fmt.Sprintf("%d", num2))
		if l1 > l2 {
			l -= l1
		} else {
			l -= l2
		}
		pres := prefixes(num1, num2)
		var pre0 string
		for i := 0; i < l; i++ {
			pre0 += "0"
		}
		for _, pre := range pres {
			pre = zipstr + pre0 + pre
			line := fmt.Sprintf("%s,%s_%s,%s\n", before, country, pre, after)
			buf.WriteString(line)
		}

	}
	if err := input.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "zipcode: %v\n", err)
		os.Exit(1)
	}
	fmt.Print(buf.String())
}

func extract(str string) (string, string) {
	var s, e = 0, 0
	var mark = true
	for i, r := range str {
		if unicode.IsNumber(r) {
			if mark {
				s = i
				mark = false
			}
			if !mark {
				e = i + 1
			}
			continue
		}
		if !mark {
			break
		}
	}
	if e == 0 {
		return str, ""
	}
	return string(str[:s]), string(str[s:e])
}

func prefixes(x, y int) []string {
	var r []string
	if x > y {
		x, y = y, x
	}
	l := len(fmt.Sprintf("%d", y))
	n := 0
	mark := true
	for {
		k := int(math.Pow(10, float64(n)))
		if mark && x != 0 && x%(10*k) == 0 {
			n++
			continue
		}
		if mark && x == 0 {
			mark = false
			n = l - 1
			continue
		}
		if y-x < k-1 {
			mark = false
			n--
			continue
		}
		f := fmt.Sprintf("%%0%dd", l-n)
		r = append(r, fmt.Sprintf(f, x/k))
		if y-x == k-1 {
			break
		}
		x += k
	}
	return r
}
