package regex

import (
	//"fmt"
	"regexp"
)

var all_space_reg = regexp.MustCompile(` 　`)
var space_reg = regexp.MustCompile(` `)
var zspace_reg = regexp.MustCompile(`　`)

func RmAllSpace(str string) string {
	result := all_space_reg.ReplaceAllString(str, "")
	return result
}

func RmSpace(str string) string {
	result := space_reg.ReplaceAllString(str, "")
	return result
}
func RmZenSpace(str string) string {
	result := zspace_reg.ReplaceAllString(str, "")
	return result
}

func RegC(rule, str string) bool {
	r := regexp.MustCompile(rule)
	return r.MatchString(str)
}

func RegM(rule, str string) []string {
	r := regexp.MustCompile(rule)
	result := r.FindAllString(str, -1)
	//	fmt.Println(result)
	return result
}

func RegSM(rule, str string) [][]string {
	r := regexp.MustCompile(rule)
	result := r.FindAllStringSubmatch(str, -1)
	// fmt.Println(result)
	return result
}

func RegR(rule, repl string, str string) string {
	r := regexp.MustCompile(rule)
	result := r.ReplaceAllString(str, repl)
	return result
}
