package request

type AnimeRecordRequest struct {
	AnimeName   string `json:"animeName"`
	AnimeRating int    `json:"animeRating"`
}
