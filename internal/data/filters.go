package data

import "strconv"

type CursorFilters struct {
	Limit  int
	Cursor string
}

type Metadata struct {
	NextCursor string `json:"nextCursor,omitempty"`
}

func calculateMetadata(count int, limit int, notes []*Note) Metadata {
	if count == 0 {
		return Metadata{}
	}

	if count < limit {
		return Metadata{}
	}

	lastNote := notes[len(notes)-1]
	return Metadata{NextCursor: strconv.FormatInt(lastNote.ID, 10)}
}
