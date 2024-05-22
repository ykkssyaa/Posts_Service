package responce_errors

type ResponseError struct {
	Message string
	Type    string
}

func (r ResponseError) Error() string {
	return r.Message
}

func (r ResponseError) Extensions() map[string]interface{} {
	return map[string]interface{}{
		"message": r.Message,
		"type":    r.Type,
	}
}
