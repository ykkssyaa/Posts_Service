package pagination

func GetOffsetAndLimit(page, pageSize *int) (offset, limit int) {
	if page == nil || pageSize == nil {
		limit = -1
		offset = 0
	} else {
		offset = (*page - 1) * *pageSize
		limit = *pageSize
	}
	return
}
