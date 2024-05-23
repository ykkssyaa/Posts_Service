package service

import (
	"database/sql"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/ykkssyaa/Posts_Service/internal/consts"
	mock_gateway "github.com/ykkssyaa/Posts_Service/internal/gateway/mock"
	"github.com/ykkssyaa/Posts_Service/internal/models"
	lg "github.com/ykkssyaa/Posts_Service/pkg/logger"
	"github.com/ykkssyaa/Posts_Service/pkg/pagination"
	"testing"
)

func IsEqualPost(want, got models.Post) bool {

	if want.ID != got.ID || want.CommentsAllowed != got.CommentsAllowed ||
		want.Content != got.Content || want.Author != got.Author {
		return false
	}

	return true
}

func TestPostsService_CreatePost(t *testing.T) {

	defaultPost := models.Post{
		ID:              1,
		Author:          "Author",
		Content:         "Cnt",
		CommentsAllowed: true,
		Name:            "Name",
	}

	tests := []struct {
		name     string
		post     models.Post
		want     models.Post
		wantErr  bool
		repoRes  models.Post
		repoErr  error
		repoSkip bool
	}{
		{
			name:     "Positive",
			post:     defaultPost,
			want:     defaultPost,
			wantErr:  false,
			repoErr:  nil,
			repoRes:  defaultPost,
			repoSkip: false,
		},
		{
			name:     "Error from repo",
			post:     defaultPost,
			want:     models.Post{},
			wantErr:  true,
			repoErr:  errors.New("some error"),
			repoRes:  models.Post{},
			repoSkip: false,
		},
		{
			name:     "Wrong author",
			post:     models.Post{Author: ""},
			want:     models.Post{},
			wantErr:  true,
			repoErr:  nil,
			repoRes:  models.Post{},
			repoSkip: true,
		},
		{
			name:     "Wrong content",
			post:     models.Post{Author: "Au1", Content: string(make([]byte, consts.MaxContentLength+1))},
			want:     models.Post{},
			wantErr:  true,
			repoErr:  nil,
			repoRes:  models.Post{},
			repoSkip: true,
		},
	}

	logger := lg.InitLogger()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctl := gomock.NewController(t)
			defer ctl.Finish()

			repo := mock_gateway.NewMockPosts(ctl)

			if !tt.repoSkip {
				repo.EXPECT().CreatePost(tt.post).Return(tt.repoRes, tt.repoErr)
			}

			p := NewPostsService(repo, logger)

			got, err := p.CreatePost(tt.post)

			if (err != nil) != tt.wantErr {
				t.Errorf("CreatePost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !IsEqualPost(got, tt.want) {
				t.Errorf("CreatePost() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPostsService_GetAllPosts(t *testing.T) {

	data := []models.Post{
		{
			ID:              1,
			Name:            "N1",
			Content:         "C1",
			Author:          "A1",
			CommentsAllowed: true,
		},
		{
			ID:              2,
			Name:            "N2",
			Content:         "C2",
			Author:          "A2",
			CommentsAllowed: true,
		},
		{
			ID:              3,
			Name:            "N3",
			Content:         "C3",
			Author:          "A3",
			CommentsAllowed: false,
		},
	}

	type args struct {
		page     *int
		pageSize *int
	}
	tests := []struct {
		name     string
		want     []models.Post
		wantErr  bool
		repoRes  []models.Post
		repoErr  error
		repoSkip bool
		args     args
	}{
		{
			name:     "Positive",
			args:     args{page: nil, pageSize: nil},
			want:     data,
			wantErr:  false,
			repoRes:  data,
			repoErr:  nil,
			repoSkip: false,
		},
		{
			name:     "Error from repo",
			args:     args{page: nil, pageSize: nil},
			want:     nil,
			wantErr:  true,
			repoRes:  nil,
			repoErr:  errors.New("some error"),
			repoSkip: false,
		},
		{
			name:     "Empty result from repo",
			args:     args{page: nil, pageSize: nil},
			want:     nil,
			wantErr:  false,
			repoRes:  nil,
			repoErr:  nil,
			repoSkip: false,
		},
		{
			name:     "With not null page number",
			args:     args{page: ptr(1), pageSize: nil},
			want:     data,
			wantErr:  false,
			repoRes:  data,
			repoErr:  nil,
			repoSkip: false,
		},
		{
			name:     "With not null page size",
			args:     args{page: nil, pageSize: ptr(1)},
			want:     data,
			wantErr:  false,
			repoRes:  data,
			repoErr:  nil,
			repoSkip: false,
		},
		{
			name:     "With not null page size and number",
			args:     args{page: ptr(1), pageSize: ptr(3)},
			want:     data,
			wantErr:  false,
			repoRes:  data,
			repoErr:  nil,
			repoSkip: false,
		},
		{
			name:     "With wrong page size",
			args:     args{page: ptr(1), pageSize: ptr(-3)},
			want:     nil,
			wantErr:  true,
			repoRes:  nil,
			repoErr:  nil,
			repoSkip: true,
		},
		{
			name:     "With wrong page number",
			args:     args{page: ptr(-1), pageSize: ptr(3)},
			want:     nil,
			wantErr:  true,
			repoRes:  nil,
			repoErr:  nil,
			repoSkip: true,
		},
		{
			name:     "With wrong page number and size",
			args:     args{page: ptr(-1), pageSize: ptr(-3)},
			want:     nil,
			wantErr:  true,
			repoRes:  nil,
			repoErr:  nil,
			repoSkip: true,
		},
	}

	logger := lg.InitLogger()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctl := gomock.NewController(t)
			defer ctl.Finish()

			repo := mock_gateway.NewMockPosts(ctl)

			if !tt.repoSkip {

				offset, limit := pagination.GetOffsetAndLimit(tt.args.page, tt.args.pageSize)
				repo.EXPECT().GetAllPosts(limit, offset).Return(tt.repoRes, tt.repoErr)
			}

			p := NewPostsService(repo, logger)
			got, err := p.GetAllPosts(tt.args.page, tt.args.pageSize)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllPosts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(got) != len(tt.want) {
				t.Errorf("GetAllPosts() got = %v len= %d, want %v len= %d", got, len(got), tt.want, len(tt.want))
				return
			}
			for i := 0; i < len(got); i++ {
				if !IsEqualPost(got[i], tt.want[i]) {
					t.Errorf("GetAllPosts() got = %v, want %v", got[i], tt.want[i])
					return
				}
			}

		})
	}
}

func TestPostsService_GetPostById(t *testing.T) {

	defaultPost := models.Post{
		ID:              1,
		Author:          "Author",
		Content:         "Cnt",
		CommentsAllowed: true,
		Name:            "Name",
	}

	logger := lg.InitLogger()

	tests := []struct {
		name     string
		postId   int
		want     models.Post
		wantErr  bool
		repoRes  models.Post
		repoErr  error
		repoSkip bool
	}{
		{
			name:     "Positive",
			postId:   1,
			want:     defaultPost,
			wantErr:  false,
			repoRes:  defaultPost,
			repoErr:  nil,
			repoSkip: false,
		},
		{
			name:     "Error from repo",
			postId:   1,
			want:     models.Post{},
			wantErr:  true,
			repoRes:  models.Post{},
			repoErr:  errors.New("some error"),
			repoSkip: false,
		},
		{
			name:     "Sql no errors from repo",
			postId:   1,
			want:     models.Post{},
			wantErr:  true,
			repoRes:  models.Post{},
			repoErr:  sql.ErrNoRows,
			repoSkip: false,
		},
		{
			name:     "Wrong Id",
			postId:   -1,
			want:     models.Post{},
			wantErr:  true,
			repoRes:  models.Post{},
			repoErr:  nil,
			repoSkip: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()

			repo := mock_gateway.NewMockPosts(ctl)

			if !tt.repoSkip {
				repo.EXPECT().GetPostById(tt.postId).Return(tt.repoRes, tt.repoErr)
			}

			p := NewPostsService(repo, logger)

			got, err := p.GetPostById(tt.postId)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetPostById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !IsEqualPost(got, tt.want) {
				t.Errorf("GetPostById() got = %v, want %v", got, tt.want)
			}
		})
	}
}
