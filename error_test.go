package server

import (
	"testing"

	"github.com/pkg/errors"
)

func TestIsErrorNotFound(t *testing.T) {
	type args struct {
		err error
	}
	type want struct {
		isNotFound bool
	}
	cases := map[string]struct {
		args
		want
	}{
		"ExpectTrue": {
			args: args{
				err: errNotFound{},
			},
			want: want{
				isNotFound: true,
			},
		},
		"ExpectFalse": {
			args: args{
				err: errors.New("some other error"),
			},
			want: want{
				isNotFound: false,
			},
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			res := IsErrorNotFound(tc.err)
			if res != tc.want.isNotFound {
				t.Errorf("Expected %v but got %v", tc.want.isNotFound, res)
			}
		})
	}
}
