package utils

import (
	"reflect"
	"testing"
)

func TestRandomNamed(t *testing.T) {
	var names []*string

	for i := 0; i < 10_000; i++ {
		n := GetRandomName()
		names = append(names, n)
	}

	for idx, n := range names {
		for _, m := range names[idx+1:] {
			if e := reflect.DeepEqual(*n, *m); e {
				t.Errorf("%v should not equal %v", *n, *m)
			}
		}
	}

}
