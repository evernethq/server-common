package pagination

func GetPageOffset(pageNum, pageSize int64) int {
	return int((pageNum - 1) * pageSize)
}
