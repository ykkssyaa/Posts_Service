package pagination

import "testing"

func TestGetOffsetAndLimit(t *testing.T) {
	type args struct {
		page     *int
		pageSize *int
	}
	tests := []struct {
		name       string
		args       args
		wantOffset int
		wantLimit  int
	}{
		{
			name:       "Both nil",
			args:       args{nil, nil},
			wantOffset: 0,
			wantLimit:  -1,
		},
		{
			name:       "Nil page",
			args:       args{nil, ptr(10)},
			wantOffset: 0,
			wantLimit:  -1,
		},
		{
			name:       "Nil pageSize",
			args:       args{ptr(2), nil},
			wantOffset: 0,
			wantLimit:  -1,
		},
		{
			name:       "First page",
			args:       args{ptr(1), ptr(10)},
			wantOffset: 0,
			wantLimit:  10,
		},
		{
			name:       "Second page",
			args:       args{ptr(2), ptr(10)},
			wantOffset: 10,
			wantLimit:  10,
		},
		{
			name:       "Third page",
			args:       args{ptr(3), ptr(5)},
			wantOffset: 10,
			wantLimit:  5,
		},
		{
			name:       "Zero page",
			args:       args{ptr(0), ptr(10)},
			wantOffset: 0,
			wantLimit:  -1,
		},
		{
			name:       "Negative page",
			args:       args{ptr(-1), ptr(10)},
			wantOffset: 0,
			wantLimit:  -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOffset, gotLimit := GetOffsetAndLimit(tt.args.page, tt.args.pageSize)
			if gotOffset != tt.wantOffset {
				t.Errorf("GetOffsetAndLimit() gotOffset = %v, want %v", gotOffset, tt.wantOffset)
			}
			if gotLimit != tt.wantLimit {
				t.Errorf("GetOffsetAndLimit() gotLimit = %v, want %v", gotLimit, tt.wantLimit)
			}
		})
	}
}

func ptr(i int) *int {
	return &i
}
