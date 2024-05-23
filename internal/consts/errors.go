package consts

const (
	EmptyAuthorError        = "empty author error"
	WrongPageError          = "wrong page number error"
	WrongPageSizeError      = "wrong page size error"
	WrongIdError            = "wrong id error"
	TooLongContentError     = "content is too long"
	CommentsNotAllowedError = "comments not allowed for this post"

	CreatingPostError = "error with creating new post"
	GettingPostError  = "error with getting post"
	PostNotFountError = "error with getting post"

	CreatingCommentError = "error with creating new comment"
	GettingCommentError  = "error with getting comments"
	GettingRepliesError  = "error with getting replies ti comment"
)

const (
	InternalErrorType = "Internal Server Error"
	BadRequestType    = "Bad Request"
	NotFoundType      = "Not Found Error"
)

const (
	WrongLimitOffsetError   = "limit and offset must be not negative"
	ThereIsNoObserversError = "there is no connection to the observer"
)
