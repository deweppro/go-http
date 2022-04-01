package routes

import (
	"reflect"
	"testing"
)

func TestUnit_split(t *testing.T) {
	type args struct {
		uri string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{name: "Case1", args: args{uri: ""}, want: []string{}},
		{name: "Case2", args: args{uri: "/a/b/"}, want: []string{"a", "b"}},
		{name: "Case3", args: args{uri: "/a/////b/"}, want: []string{"a", "b"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := split(tt.args.uri); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("split() = %v, want %v", got, tt.want)
			}
		})
	}
}
