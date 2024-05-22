package service

import (
	"github.com/ykkssyaa/Posts_Service/internal/consts"
	"github.com/ykkssyaa/Posts_Service/internal/gateway"
	"github.com/ykkssyaa/Posts_Service/internal/models"
	"github.com/ykkssyaa/Posts_Service/pkg/logger"
	re "github.com/ykkssyaa/Posts_Service/pkg/responce_errors"
)

type CommentsService struct {
	repo   gateway.Comments
	logger *logger.Logger
}

func NewCommentsService(repo gateway.Comments, logger *logger.Logger) *CommentsService {
	return &CommentsService{repo: repo, logger: logger}
}

func (c CommentsService) CreateComment(comment models.Comment) (models.Comment, error) {
	if len(comment.Author) == 0 {
		c.logger.Err.Println(consts.EmptyAuthorError)
		return models.Comment{}, re.ResponseError{
			Message: consts.EmptyAuthorError,
			Type:    consts.BadRequestType,
		}
	}

	if len(comment.Content) >= consts.MaxContentLength {
		c.logger.Err.Println(consts.TooLongContentError, len(comment.Content))
		return models.Comment{}, re.ResponseError{
			Message: consts.TooLongContentError,
			Type:    consts.BadRequestType,
		}
	}

	newComment, err := c.repo.CreateComment(comment)
	if err != nil {
		c.logger.Err.Println(consts.CreatingCommentError, err.Error())
		return models.Comment{}, re.ResponseError{
			Message: consts.CreatingCommentError,
			Type:    consts.InternalErrorType,
		}
	}

	return newComment, nil
}

func (c CommentsService) GetCommentsByPost(postId int) ([]*models.Comment, error) {

	if postId <= 0 {
		c.logger.Err.Println(consts.WrongIdError, postId)
		return nil, re.ResponseError{
			Message: consts.WrongIdError,
			Type:    consts.BadRequestType,
		}
	}

	comments, err := c.repo.GetCommentsByPost(postId)
	if err != nil {
		c.logger.Err.Println(consts.GettingCommentError, postId, err.Error())
		return nil, re.ResponseError{
			Message: consts.GettingCommentError,
			Type:    consts.InternalErrorType,
		}
	}

	return comments, nil
}
