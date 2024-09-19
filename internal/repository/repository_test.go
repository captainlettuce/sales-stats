package repository

import (
	"github.com/Masterminds/squirrel"
	"github.com/captainlettuce/sales-stats/pkg/types"
	"reflect"
	"testing"
	"time"
)

func Test_whereOrWhereIn(t *testing.T) {
	type args struct {
		q      squirrel.SelectBuilder
		col    string
		filter []int
	}
	type testCase struct {
		name string
		args args
		want squirrel.SelectBuilder
	}
	tests := []testCase{
		{
			name: "singular input generates where",
			args: args{
				col:    "col",
				filter: []int{1},
			},
			want: squirrel.SelectBuilder{}.Where(squirrel.Eq{"col": 1}),
		},
		{
			name: "multiple input generates where in (...)",
			args: args{
				col:    "multiCol",
				filter: []int{1, 2},
			},
			want: squirrel.SelectBuilder{}.Where(squirrel.Eq{"multiCol": []int{1, 2}}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := whereOrWhereIn(tt.args.q, tt.args.col, tt.args.filter); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("whereOrWhereIn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_applyDateRange(t *testing.T) {
	time1, err := time.Parse(time.DateOnly, "2024-01-01")
	if err != nil {
		t.Fatal(err)
	}
	time2, err := time.Parse(time.DateOnly, "2024-01-02")
	if err != nil {
		t.Fatal(err)
	}
	date1 := types.Date{time1}
	date2 := types.Date{time2}

	// The column is assumed to be named 'col' so use that in the want-clause if the column name is needed
	type args struct {
		r *types.DateRange
	}
	tests := []struct {
		name string
		args args
		want squirrel.SelectBuilder
	}{
		{
			name: "basic regression test",
			args: args{
				r: &types.DateRange{Before: &date2, After: &date1},
			},
			want: squirrel.SelectBuilder{}.Where("col::DATE BETWEEN ?::DATE AND ?::DATE", date1.Format(time.DateOnly), date2.Format(time.DateOnly)),
		},
		{
			name: "after",
			args: args{
				r: &types.DateRange{After: &date1},
			},
			want: squirrel.SelectBuilder{}.Where("col::DATE > ?::DATE", date1.Format(time.DateOnly)),
		},
		{
			name: "before",
			args: args{
				r: &types.DateRange{Before: &date1},
			},
			want: squirrel.SelectBuilder{}.Where("col::DATE < ?::DATE", date1.Format(time.DateOnly)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := applyDateRange(squirrel.SelectBuilder{}, "col", tt.args.r); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("applyDateRange() = %v, want %v", got, tt.want)
			}
		})
	}
}

func ref[T any](v T) *T {
	return &v
}

func Test_applyIntRange(t *testing.T) {

	// The column is assumed to be named 'col' so use that in the want-clause if the column name is needed
	type args struct {
		r types.IntegerRange
	}
	tests := []struct {
		name string
		args args
		want squirrel.SelectBuilder
	}{
		{
			name: "over and under",
			args: args{
				r: types.IntegerRange{
					Over:  ref(int64(1)),
					Under: ref(int64(2)),
				},
			},
			want: squirrel.SelectBuilder{}.Where("col > ?", ref(int64(1))).Where("col < ?", ref(int64(2))),
		},
		{
			name: "over",
			args: args{
				r: types.IntegerRange{
					Over: ref(int64(1)),
				},
			},
			want: squirrel.SelectBuilder{}.Where("col > ?", ref(int64(1))),
		},
		{
			name: "under",
			args: args{
				r: types.IntegerRange{
					Under: ref(int64(1)),
				},
			},
			want: squirrel.SelectBuilder{}.Where("col < ?", ref(int64(1))),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := applyIntRange(squirrel.SelectBuilder{}, "col", tt.args.r); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("applyIntRange() = %v, want %v", got, tt.want)
			}
		})
	}
}
