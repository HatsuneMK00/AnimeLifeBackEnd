package response

import "time"

type AnimeRecord struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
	Name      string    `json:"name"`
	NameJp    string    `json:"name_jp"`
	Cover     string    `json:"cover"`
	BangumiId int       `json:"bangumi_id"`
	Rating    int       `json:"rating"`
	Comment   string    `json:"comment"`
	RecordAt  time.Time `json:"record_at"`
}
