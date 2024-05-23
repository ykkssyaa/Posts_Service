package service

import (
	"database/sql"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/ykkssyaa/Posts_Service/internal/consts"
	mock_gateway "github.com/ykkssyaa/Posts_Service/internal/gateway/mock"
	"github.com/ykkssyaa/Posts_Service/internal/models"
	mock_service "github.com/ykkssyaa/Posts_Service/internal/service/mock"
	lg "github.com/ykkssyaa/Posts_Service/pkg/logger"
	"github.com/ykkssyaa/Posts_Service/pkg/pagination"
	"github.com/ykkssyaa/Posts_Service/pkg/responce_errors"
	"testing"
)

func IsEqualComment(want, got models.Comment) bool {
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

func TestCommentsService_CreateComment(t *testing.T) {

	type FuncRes struct {
		com  models.Comment
		err  error
		skip bool
	}

	defaultComment := models.Comment{
		ID:      1,
		Post:    1,
		Author:  "Au1",
		Content: "Text1",
		ReplyTo: nil,
	}

	defaultPost := models.Post{
		ID:              1,
		Author:          "Author2",
		Content:         "Cnt2",
		CommentsAllowed: true,
		Name:            "Name",
	}

	postWithoutAllowedComments := models.Post{
		ID:              1,
		Author:          "Author3",
		Content:         "Cnt3",
		CommentsAllowed: false,
		Name:            "Name3",
	}

	type args struct {
		comment models.Comment
	}
	tests := []struct {
		name       string
		args       args
		want       FuncRes
		repoRes    FuncRes
		getterRes  models.Post
		getterErr  error
		skipGetter bool
	}{
		{
			name:      "Positive creation",
			args:      args{comment: defaultComment},
			want:      FuncRes{com: defaultComment, err: nil},
			repoRes:   FuncRes{com: defaultComment, err: nil},
			getterRes: defaultPost,
			getterErr: nil,
		},
		{
			name:      "Comments Not Allowed",
			args:      args{comment: defaultComment},
			want:      FuncRes{com: models.Comment{}, err: responce_errors.ResponseError{}},
			repoRes:   FuncRes{com: defaultComment, err: nil, skip: true},
			getterRes: postWithoutAllowedComments,
			getterErr: nil,
		},
		{
			name:      "Post Not Found",
			args:      args{comment: defaultComment},
			want:      FuncRes{com: models.Comment{}, err: responce_errors.ResponseError{}},
			repoRes:   FuncRes{com: models.Comment{}, err: nil, skip: true},
			getterRes: models.Post{},
			getterErr: sql.ErrNoRows,
		},
		{
			name:      "Error with creating",
			args:      args{comment: defaultComment},
			want:      FuncRes{com: models.Comment{}, err: responce_errors.ResponseError{}},
			repoRes:   FuncRes{com: models.Comment{}, err: errors.New("some error")},
			getterRes: defaultPost,
			getterErr: nil,
		},
		{
			name:       "Wrong Author",
			args:       args{comment: models.Comment{Author: ""}},
			want:       FuncRes{com: models.Comment{}, err: responce_errors.ResponseError{}},
			repoRes:    FuncRes{com: models.Comment{}, err: nil, skip: true},
			getterRes:  defaultPost,
			getterErr:  nil,
			skipGetter: true,
		},
		{
			name:       "Wrong Content",
			args:       args{comment: models.Comment{Author: "A", Content: string(make([]byte, consts.MaxContentLength+1))}},
			want:       FuncRes{com: models.Comment{}, err: responce_errors.ResponseError{}},
			repoRes:    FuncRes{com: models.Comment{}, err: nil, skip: true},
			getterRes:  defaultPost,
			getterErr:  nil,
			skipGetter: true,
		},
		{
			name:       "Wrong Post Id",
			args:       args{comment: models.Comment{Author: "A", Post: -1}},
			want:       FuncRes{com: models.Comment{}, err: responce_errors.ResponseError{}},
			repoRes:    FuncRes{com: models.Comment{}, err: nil, skip: true},
			getterRes:  defaultPost,
			getterErr:  nil,
			skipGetter: true,
		},
	}

	logger := lg.InitLogger()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctl := gomock.NewController(t)
			defer ctl.Finish()

			repo := mock_gateway.NewMockComments(ctl)
			getter := mock_service.NewMockPostGetter(ctl)

			if !tt.repoRes.skip {
				repo.EXPECT().CreateComment(tt.args.comment).Return(tt.repoRes.com, tt.repoRes.err)
			}

			if !tt.skipGetter {
				getter.EXPECT().GetPostById(tt.args.comment.Post).Return(tt.getterRes, tt.getterErr)
			}

			c := NewCommentsService(repo, logger, getter)

			got, err := c.CreateComment(tt.args.comment)

			if (err != nil) != (tt.want.err != nil) {
				t.Errorf("CreateComment() error = %v, wantErr %v", err, tt.want.err)
				return
			}

			if !IsEqualComment(got, tt.want.com) {
				t.Errorf("CreateComment() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCommentsService_GetCommentsByPost(t *testing.T) {

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
			Post:    2,
			Content: "Test6",
			Author:  "Au",
			ReplyTo: nil,
		},
		{
			ID:      5,
			Post:    2,
			Content: "Test8",
			Author:  "Audff",
			ReplyTo: nil,
		},
	}

	type args struct {
		postId   int
		page     *int
		pageSize *int
	}
	tests := []struct {
		name     string
		args     args
		want     []*models.Comment
		wantErr  bool
		repoRes  []*models.Comment
		repoErr  error
		repoSkip bool
	}{
		{
			name:     "positive getting",
			args:     args{postId: 1, pageSize: nil, page: nil},
			want:     data[:3],
			wantErr:  false,
			repoRes:  data[:3],
			repoErr:  nil,
			repoSkip: false,
		},
		{
			name:     "error from repo",
			args:     args{postId: 1, pageSize: nil, page: nil},
			want:     nil,
			wantErr:  true,
			repoRes:  nil,
			repoErr:  errors.New("some error"),
			repoSkip: false,
		},
		{
			name:     "empty repo result",
			args:     args{postId: 1, pageSize: nil, page: nil},
			want:     nil,
			wantErr:  false,
			repoRes:  nil,
			repoErr:  nil,
			repoSkip: false,
		},
		{
			name:     "wrong postId",
			args:     args{postId: -1, pageSize: nil, page: nil},
			want:     nil,
			wantErr:  true,
			repoRes:  nil,
			repoErr:  nil,
			repoSkip: true,
		},
		{
			name:     "wrong pageSize",
			args:     args{postId: 1, pageSize: ptr(-1), page: nil},
			want:     nil,
			wantErr:  true,
			repoRes:  nil,
			repoErr:  nil,
			repoSkip: true,
		},
		{
			name:     "wrong page number",
			args:     args{postId: 1, pageSize: nil, page: ptr(-10)},
			want:     nil,
			wantErr:  true,
			repoRes:  nil,
			repoErr:  nil,
			repoSkip: true,
		},
		{
			name:     "with nil pageSize",
			args:     args{postId: 1, pageSize: nil, page: ptr(10)},
			want:     data[:1],
			wantErr:  false,
			repoRes:  data[:1],
			repoErr:  nil,
			repoSkip: false,
		},
		{
			name:     "with nil page number",
			args:     args{postId: 1, pageSize: ptr(1), page: nil},
			want:     data[:1],
			wantErr:  false,
			repoRes:  data[:1],
			repoErr:  nil,
			repoSkip: false,
		},
		{
			name:     "with not nil page number and pageSize",
			args:     args{postId: 1, pageSize: ptr(1), page: ptr(1)},
			want:     data[:1],
			wantErr:  false,
			repoRes:  data[:1],
			repoErr:  nil,
			repoSkip: false,
		},
	}

	logger := lg.InitLogger()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctl := gomock.NewController(t)
			defer ctl.Finish()

			repo := mock_gateway.NewMockComments(ctl)
			getter := mock_service.NewMockPostGetter(ctl)

			if !tt.repoSkip {
				offset, limit := pagination.GetOffsetAndLimit(tt.args.page, tt.args.pageSize)
				repo.EXPECT().GetCommentsByPost(tt.args.postId, limit, offset).
					Return(tt.repoRes, tt.repoErr)
			}

			c := NewCommentsService(repo, logger, getter)

			got, err := c.GetCommentsByPost(tt.args.postId, tt.args.page, tt.args.pageSize)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCommentsByPost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(got) != len(tt.want) {
				t.Errorf("GetCommentsByPost() got = %v len= %d, want %v len= %d", got, len(got), tt.want, len(tt.want))
				return
			}
			for i := 0; i < len(got); i++ {
				if !IsEqualComment(*got[i], *tt.want[i]) {
					t.Errorf("GetCommentsByPost() got = %v, want %v", got[i], tt.want[i])
					return
				}
			}

		})
	}
}

func TestCommentsService_GetRepliesOfComment(t *testing.T) {

	data := []*models.Comment{
		{
			ID:      2,
			Post:    1,
			Content: "Test2",
			Author:  "Au2",
			ReplyTo: ptr(1),
		},
		{
			ID:      3,
			Post:    1,
			Content: "Test3",
			Author:  "Au",
			ReplyTo: ptr(1),
		},
		{
			ID:      4,
			Post:    1,
			Content: "Test6",
			Author:  "Au",
			ReplyTo: ptr(1),
		},
		{
			ID:      5,
			Post:    1,
			Content: "Test8",
			Author:  "Audff",
			ReplyTo: ptr(1),
		},
	}

	type args struct {
		commentId int
	}
	tests := []struct {
		name     string
		args     args
		want     []*models.Comment
		wantErr  bool
		repoRes  []*models.Comment
		repoErr  error
		repoSkip bool
	}{
		{
			name:     "Positive Getting",
			args:     args{commentId: 1},
			want:     data,
			wantErr:  false,
			repoRes:  data,
			repoSkip: false,
			repoErr:  nil,
		},
		{
			name:     "Wrong Comment Id",
			args:     args{commentId: -1},
			want:     nil,
			wantErr:  true,
			repoRes:  nil,
			repoSkip: true,
			repoErr:  nil,
		},
		{
			name:     "Error from repo",
			args:     args{commentId: 1},
			want:     nil,
			wantErr:  true,
			repoRes:  nil,
			repoSkip: false,
			repoErr:  errors.New("some error"),
		},
		{
			name:     "Nil result from repo",
			args:     args{commentId: 1},
			want:     nil,
			wantErr:  false,
			repoRes:  nil,
			repoSkip: false,
			repoErr:  nil,
		},
	}

	logger := lg.InitLogger()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctl := gomock.NewController(t)
			defer ctl.Finish()

			repo := mock_gateway.NewMockComments(ctl)
			getter := mock_service.NewMockPostGetter(ctl)

			if !tt.repoSkip {
				repo.EXPECT().GetRepliesOfComment(tt.args.commentId).Return(tt.repoRes, tt.repoErr)
			}

			c := NewCommentsService(repo, logger, getter)

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
				if !IsEqualComment(*got[i], *tt.want[i]) {
					t.Errorf("GetCommentsByPost() got = %v, want %v", got[i], tt.want[i])
					return
				}
			}
		})
	}
}

func ptr(i int) *int {
	p := i
	return &p
}
