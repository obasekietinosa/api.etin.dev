package data

type CursorFilters struct {
	Limit  int
	Cursor string
}

type Metadata struct {
	NextCursor string `json:"nextCursor,omitempty"`
}
