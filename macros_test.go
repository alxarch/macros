package macros

import (
	"testing"
)

func TestReplaceMacros(t *testing.T) {
	a := func(dst []byte, tag string) ([]byte, error) {
		switch tag {
		case "FOO":
			dst = append(dst, "foo"...)
		case "BAR":
			dst = append(dst, "bar"...)
		default:
			dst = append(dst, "baz"...)
		}
		return dst, nil
	}
	tpl := Quick("${", "}")
	w, err := tpl.Replace(nil, "${FOO} ${BAR} ${BAZ}", a)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if string(w) != "foo bar baz" {
		t.Errorf("Invalid replace: %s != %s", string(w), "foo bar baz")
	}
}

// func BenchmarkQuickTemplateExecute(b *testing.B) {
// 	b.ReportAllocs()
// 	a := func(dst []byte, tag string) ([]byte, error) {
// 		switch tag {
// 		case "FOO":
// 			dst = append(dst, "foo"...)
// 		case "BAR":
// 			dst = append(dst, "bar"...)
// 		default:
// 			dst = append(dst, "baz"...)
// 		}
// 		return dst, nil
// 	}
// 	w := bytes.NewBuffer(make([]byte, 0, 64))
// 	tpl := Quick("${", "}", make([]byte, 0, 64))
// 	var err error
// 	for i := 0; i < b.N; i++ {
// 		w.Reset()
// 		_, err = tpl.Execute(w, "${FOO} ${BAR} ${BAZ}", a)
// 		if err != nil {
// 			b.Error(err)
// 			return
// 		}
// 	}

// }

func BenchmarkQuickTemplateReplace(b *testing.B) {
	b.ReportAllocs()
	a := func(dst []byte, tag string) ([]byte, error) {
		switch tag {
		case "FOO":
			dst = append(dst, "foo"...)
		case "BAR":
			dst = append(dst, "bar"...)
		default:
			dst = append(dst, "baz"...)
		}
		return dst, nil
	}
	w := make([]byte, 0, 64)
	tpl := Quick("${", "}")
	var err error
	for i := 0; i < b.N; i++ {
		_, err = tpl.Replace(w[:0], "${FOO} ${BAR} ${BAZ}", a)
		if err != nil {
			b.Error(err)
			return
		}
	}

}
