package dumper_test

import (
	"fmt"
	"testing"

	"github.com/itoken417/goutils/util/dumper"
)

type Address struct {
	Zip  string
	City string
}

type Person struct {
	Name    string
	Age     int
	Scores  []float64
	Tags    map[string]string
	Address Address
	Friend  *Person
}

func TestDumpPrimitives(t *testing.T) {
	fmt.Print(dumper.Dump(42, "hello", true, 3.14))
}

func TestDumpStruct(t *testing.T) {
	p := Person{
		Name:   "田中",
		Age:    30,
		Scores: []float64{95.5, 80.0, 72.3},
		Tags:   map[string]string{"role": "admin", "lang": "go"},
		Address: Address{
			Zip:  "100-0001",
			City: "東京都千代田区",
		},
	}
	fmt.Print(dumper.DumpNamed("person", p))
}

func TestDumpPointerAndCycle(t *testing.T) {
	a := &Person{Name: "Alice", Age: 25}
	b := &Person{Name: "Bob", Age: 28, Friend: a}
	a.Friend = b // 循環参照
	fmt.Print(dumper.DumpNamed("alice", a))
}

func TestDumpNilAndEmpty(t *testing.T) {
	fmt.Print(dumper.Dump(nil, []int(nil), []int{}, map[string]int(nil)))
}

func TestDd(t *testing.T) {
	dumper.Dd([]int{1, 2, 3}, map[string]bool{"ok": true})
}
