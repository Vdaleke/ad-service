package ads

import "time"

type Ad struct {
	ID        int64
	Title     string
	Text      string
	AuthorID  int64
	Published bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
