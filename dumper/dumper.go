// Package dumper はPerlのData::Dumperに相当するデバッグ用ダンプユーティリティです。
package dumper

import (
	"fmt"
	"reflect"
	"strings"
	"unsafe"
)

const indentStr = "  "

// Dump は渡した値を人間が読みやすい形式の文字列で返します。
// 複数の値を渡すと $var1, $var2, ... として連結されます。
func Dump(vars ...any) string {
	var sb strings.Builder
	for i, v := range vars {
		fmt.Fprintf(&sb, "$var%d = ", i+1)
		dumpValue(&sb, reflect.ValueOf(v), 0, make(map[uintptr]bool))
		sb.WriteString("\n")
	}
	return sb.String()
}

// DumpNamed は変数名付きでダンプ文字列を返します。
func DumpNamed(name string, v any) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "$%s = ", name)
	dumpValue(&sb, reflect.ValueOf(v), 0, make(map[uintptr]bool))
	sb.WriteString("\n")
	return sb.String()
}

// Dd はダンプ結果を標準出力に書き出します（デバッグ用）。
func Dd(vars ...any) {
	fmt.Print(Dump(vars...))
}

// DdNamed は変数名付きでダンプ結果を標準出力に書き出します。
func DdNamed(name string, v any) {
	fmt.Print(DumpNamed(name, v))
}

func dumpValue(sb *strings.Builder, v reflect.Value, depth int, visited map[uintptr]bool) {
	if !v.IsValid() {
		sb.WriteString("nil")
		return
	}

	switch v.Kind() {
	case reflect.Ptr, reflect.Interface:
		if v.IsNil() {
			sb.WriteString("nil")
			return
		}
		if v.Kind() == reflect.Ptr {
			ptr := v.Pointer()
			if visited[ptr] {
				fmt.Fprintf(sb, "*(%p) /* 循環参照 */", unsafe.Pointer(v.Pointer()))
				return
			}
			visited[ptr] = true
			defer delete(visited, ptr)
			sb.WriteString("&")
			dumpValue(sb, v.Elem(), depth, visited)
		} else {
			dumpValue(sb, v.Elem(), depth, visited)
		}

	case reflect.Struct:
		t := v.Type()
		fmt.Fprintf(sb, "%s{\n", t.Name())
		for i := range t.NumField() {
			f := t.Field(i)
			indent(sb, depth+1)
			fmt.Fprintf(sb, "%s: ", f.Name)
			fv := v.Field(i)
			// エクスポートされていないフィールドは型情報のみ表示
			if !f.IsExported() {
				fmt.Fprintf(sb, "<%s>", f.Type)
			} else {
				dumpValue(sb, fv, depth+1, visited)
			}
			sb.WriteString(",\n")
		}
		indent(sb, depth)
		sb.WriteString("}")

	case reflect.Map:
		if v.IsNil() {
			fmt.Fprintf(sb, "map[%s]%s(nil)", v.Type().Key(), v.Type().Elem())
			return
		}
		fmt.Fprintf(sb, "map[%s]%s{\n", v.Type().Key(), v.Type().Elem())
		for _, key := range v.MapKeys() {
			indent(sb, depth+1)
			dumpValue(sb, key, depth+1, visited)
			sb.WriteString(": ")
			dumpValue(sb, v.MapIndex(key), depth+1, visited)
			sb.WriteString(",\n")
		}
		indent(sb, depth)
		sb.WriteString("}")

	case reflect.Slice:
		if v.IsNil() {
			fmt.Fprintf(sb, "[]%s(nil)", v.Type().Elem())
			return
		}
		fallthrough
	case reflect.Array:
		if v.Len() == 0 {
			fmt.Fprintf(sb, "[]%s{}", v.Type().Elem())
			return
		}
		fmt.Fprintf(sb, "[]%s{\n", v.Type().Elem())
		for i := range v.Len() {
			indent(sb, depth+1)
			dumpValue(sb, v.Index(i), depth+1, visited)
			sb.WriteString(",\n")
		}
		indent(sb, depth)
		sb.WriteString("}")

	case reflect.String:
		fmt.Fprintf(sb, "%q", v.String())

	case reflect.Bool:
		fmt.Fprintf(sb, "%t", v.Bool())

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		fmt.Fprintf(sb, "%d", v.Int())

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		fmt.Fprintf(sb, "%d", v.Uint())

	case reflect.Float32, reflect.Float64:
		fmt.Fprintf(sb, "%g", v.Float())

	case reflect.Complex64, reflect.Complex128:
		fmt.Fprintf(sb, "%g", v.Complex())

	case reflect.Chan:
		if v.IsNil() {
			fmt.Fprintf(sb, "(%s)(nil)", v.Type())
		} else {
			fmt.Fprintf(sb, "(%s)(%p)", v.Type(), unsafe.Pointer(v.Pointer()))
		}

	case reflect.Func:
		if v.IsNil() {
			fmt.Fprintf(sb, "(%s)(nil)", v.Type())
		} else {
			fmt.Fprintf(sb, "(%s)(%p)", v.Type(), unsafe.Pointer(v.Pointer()))
		}

	default:
		fmt.Fprintf(sb, "%#v", v.Interface())
	}
}

func indent(sb *strings.Builder, depth int) {
	for range depth {
		sb.WriteString(indentStr)
	}
}
