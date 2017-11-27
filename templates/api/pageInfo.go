package api

type pageInfoResolver struct {
	startID     string
	endID       string
	hasNext     bool
	hasPrevious bool
}

func (r *pageInfoResolver) StartCursor() *string {
	cursor := encodeCursor(r.startID)
	return &cursor
}

func (r *pageInfoResolver) EndCursor() *string {
	cursor := encodeCursor(r.endID)
	return &cursor
}

func (r *pageInfoResolver) HasNextPage() bool {
	return r.hasNext
}

func (r *pageInfoResolver) HasPreviousPage() bool {
	return r.hasPrevious
}
