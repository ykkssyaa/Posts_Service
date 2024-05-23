package in_memory

import (
	"github.com/ykkssyaa/Posts_Service/internal/models"
	"testing"
)

func isEqual(want, got models.Comment) bool {
	if want.Post != got.Post || want.ID != got.ID || want.Author != got.Author || want.Content != got.Content {
		return false
	}

	if (want.ReplyTo == nil) != (got.ReplyTo == nil) {
		return false
	}

	if (want.ReplyTo != nil) && (got.ReplyTo != nil) && *got.ReplyTo != *want.ReplyTo {
		return false
	}

	return true
}

func TestCommentsInMemory_CreateComment(t *testing.T) {

	type args struct {
		comment models.Comment
	}
	tests := []struct {
		name    string
		args    args
		want    models.Comment
		wantErr bool
	}{
		{
			name: "Positive Creating",
			args: args{comment: models.Comment{
				Post:    1,
				Content: "Test",
				Author:  "Au",
				ReplyTo: nil,
			}},
			want: models.Comment{
				Post:    1,
				Content: "Test",
				Author:  "Au",
				ReplyTo: nil,
				ID:      1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCommentsInMemory(10)

			got, err := c.CreateComment(tt.args.comment)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateComment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !isEqual(got, tt.want) {
				t.Errorf("CreateComment() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCommentsInMemory_GetCommentsByPost(t *testing.T) {

	data := []*models.Comment{
		{
			ID:      1,
			Post:    1,
			Content: "Test1",
			Author:  "Au",
			ReplyTo: nil,
		},
		{
			ID:      2,
			Post:    1,
			Content: "Test2",
			Author:  "Au2",
			ReplyTo: nil,
		},
		{
			ID:      3,
			Post:    1,
			Content: "Test3",
			Author:  "Au",
			ReplyTo: nil,
		},
		{
			ID:      4,
			Post:    1,
			Content: "Test4",
			Author:  "Au",
			ReplyTo: ptr(1),
		},
		{
			ID:      5,
			Post:    1,
			Content: "Test5",
			Author:  "Au",
			ReplyTo: ptr(1),
		},
		{
			ID:      6,
			Post:    2,
			Content: "Test6",
			Author:  "Au",
			ReplyTo: nil,
		},
		{
			ID:      7,
			Post:    2,
			Content: "Test6",
			Author:  "Au",
			ReplyTo: ptr(6),
		},
		{
			ID:      8,
			Post:    2,
			Content: "Test6",
			Author:  "Au",
			ReplyTo: ptr(7),
		},
	}

	type args struct {
		postId int
		limit  int
		offset int
	}
	tests := []struct {
		name    string
		args    args
		data    []*models.Comment
		want    []*models.Comment
		wantErr bool
	}{
		{
			name:    "get all comments 1",
			args:    args{limit: -1, offset: 0, postId: 1},
			data:    data,
			want:    data[:3],
			wantErr: false,
		},
		{
			name:    "get all comments 2",
			args:    args{limit: -1, offset: 0, postId: 2},
			data:    data,
			want:    data[5:6],
			wantErr: false,
		},
		{
			name:    "get comments with offset",
			args:    args{limit: -1, offset: 1, postId: 1},
			data:    data,
			want:    data[1:3],
			wantErr: false,
		},
		{
			name:    "get comments with big offset",
			args:    args{limit: -1, offset: 10, postId: 1},
			data:    data,
			want:    nil,
			wantErr: false,
		},
		{
			name:    "get comments with limit",
			args:    args{limit: 2, offset: 0, postId: 1},
			data:    data,
			want:    data[0:2],
			wantErr: false,
		},
		{
			name:    "get comments with limit 1",
			args:    args{limit: 2, offset: 0, postId: 1},
			data:    data,
			want:    data[0:2],
			wantErr: false,
		},
		{
			name:    "get comments with limit 0",
			args:    args{limit: 0, offset: 0, postId: 1},
			data:    data,
			want:    nil,
			wantErr: false,
		},
		{
			name:    "get all comments with big limit",
			args:    args{limit: 100, offset: 0, postId: 1},
			data:    data,
			want:    data[:3],
			wantErr: false,
		},
		{
			name:    "get all comments with negative limit",
			args:    args{limit: -10, offset: 0, postId: 1},
			data:    data,
			want:    nil,
			wantErr: true,
		},
		{
			name:    "get all comments with negative offset",
			args:    args{limit: 10, offset: -10, postId: 1},
			data:    data,
			want:    nil,
			wantErr: true,
		},
		{
			name:    "get all comments with negative offset and limit",
			args:    args{limit: -10, offset: -10, postId: 1},
			data:    data,
			want:    nil,
			wantErr: true,
		},
		{
			name:    "get all comments with offset and big limit",
			args:    args{limit: 10, offset: 1, postId: 1},
			data:    data,
			want:    data[1:3],
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCommentsInMemory(10)

			for _, comment := range tt.data {
				_, _ = c.CreateComment(*comment)
			}

			got, err := c.GetCommentsByPost(tt.args.postId, tt.args.limit, tt.args.offset)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetCommentsByPost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != len(tt.want) {
				t.Errorf("GetCommentsByPost() got = %v len= %d, want %v len= %d", got, len(got), tt.want, len(tt.want))
				return
			}
			for i := 0; i < len(got); i++ {
				if !isEqual(*got[i], *tt.want[i]) {
					t.Errorf("GetCommentsByPost() got = %v, want %v", got[i], tt.want[i])
					return
				}
			}
		})
	}
}

func TestCommentsInMemory_GetRepliesOfComment(t *testing.T) {

	data := []*models.Comment{
		{
			ID:      1,
			Post:    1,
			Content: "Test1",
			Author:  "Au",
			ReplyTo: nil,
		},
		{
			ID:      2,
			Post:    1,
			Content: "Test2",
			Author:  "Au2",
			ReplyTo: nil,
		},
		{
			ID:      3,
			Post:    1,
			Content: "Test3",
			Author:  "Au",
			ReplyTo: nil,
		},
		{
			ID:      4,
			Post:    1,
			Content: "Test4",
			Author:  "Au",
			ReplyTo: ptr(1),
		},
		{
			ID:      5,
			Post:    1,
			Content: "Test5",
			Author:  "Au",
			ReplyTo: ptr(1),
		},
		{
			ID:      6,
			Post:    2,
			Content: "Test6",
			Author:  "Au",
			ReplyTo: nil,
		},
		{
			ID:      7,
			Post:    2,
			Content: "Test6",
			Author:  "Au",
			ReplyTo: ptr(6),
		},
		{
			ID:      8,
			Post:    2,
			Content: "Test6",
			Author:  "Au",
			ReplyTo: ptr(7),
		},
	}

	type args struct {
		commentId int
	}
	tests := []struct {
		name    string
		args    args
		want    []*models.Comment
		data    []*models.Comment
		wantErr bool
	}{
		{
			name:    "Getting replies 1",
			args:    args{commentId: 1},
			want:    data[3:5],
			data:    data,
			wantErr: false,
		},
		{
			name:    "Getting replies 2",
			args:    args{commentId: 6},
			want:    data[6:7],
			data:    data,
			wantErr: false,
		},
		{
			name:    "Getting replies 3",
			args:    args{commentId: 7},
			want:    data[7:8],
			data:    data,
			wantErr: false,
		},
		{
			name:    "Getting nil replies",
			args:    args{commentId: 10},
			want:    nil,
			data:    data,
			wantErr: false,
		},
		{
			name:    "Getting replies with empty storage",
			args:    args{commentId: 10},
			want:    nil,
			data:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCommentsInMemory(10)

			for _, comment := range tt.data {
				_, _ = c.CreateComment(*comment)
			}

			got, err := c.GetRepliesOfComment(tt.args.commentId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRepliesOfComment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(got) != len(tt.want) {
				t.Errorf("GetCommentsByPost() got = %v len= %d, want %v len= %d", got, len(got), tt.want, len(tt.want))
				return
			}
			for i := 0; i < len(got); i++ {
				if !isEqual(*got[i], *tt.want[i]) {
					t.Errorf("GetCommentsByPost() got = %v, want %v", got[i], tt.want[i])
					return
				}
			}
		})
	}
}

func TestNewCommentsInMemory(t *testing.T) {
	type args struct {
		count int
	}
	tests := []struct {
		name    string
		args    args
		wantCap int
	}{
		{
			name:    "positive creating",
			args:    args{count: 10},
			wantCap: 10,
		},
		{
			name:    "creating big slice",
			args:    args{count: 1000},
			wantCap: 1000,
		},
		{
			name:    "creating small slice",
			args:    args{count: 1},
			wantCap: 1,
		},
		{
			name:    "creating small slice",
			args:    args{count: 0},
			wantCap: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			c := NewCommentsInMemory(tt.args.count)
			if c.comments == nil || cap(c.comments) != tt.wantCap {
				t.Errorf("NewCommentsInMemory() error creating want: %v got %v", tt.wantCap, cap(c.comments))
			}

		})
	}
}

func ptr(i int) *int {
	p := i
	return &p
}
