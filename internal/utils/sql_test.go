package utils

import (
	"database/sql"
	"reflect"
	"testing"
	"time"
)

func TestToNullString(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want sql.NullString
	}{
		{
			name: "empty string",
			args: args{
				s: "",
			},
			want: sql.NullString{Valid: false},
		},
		{
			name: "non-empty string",
			args: args{
				s: "test",
			},
			want: sql.NullString{String: "test", Valid: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToNullString(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToNullString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToNullTime(t *testing.T) {
	timestamp := time.Now()
	type args struct {
		t time.Time
	}
	tests := []struct {
		name string
		args args
		want sql.NullTime
	}{
		{
			name: "zero time",
			args: args{
				t: time.Time{},
			},
			want: sql.NullTime{Valid: false},
		},
		{
			name: "non-zero time",
			args: args{
				t: timestamp,
			},
			want: sql.NullTime{Time: timestamp, Valid: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToNullTime(tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToNullTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToNullInt32(t *testing.T) {
	type args struct {
		i int32
	}
	tests := []struct {
		name string
		args args
		want sql.NullInt32
	}{
		{
			name: "zero int",
			args: args{
				i: 0,
			},
			want: sql.NullInt32{Valid: false},
		},
		{
			name: "non-zero int",
			args: args{
				i: 42,
			},
			want: sql.NullInt32{Int32: 42, Valid: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToNullInt32(tt.args.i); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToNullInt32() = %v, want %v", got, tt.want)
			}
		})
	}
}
