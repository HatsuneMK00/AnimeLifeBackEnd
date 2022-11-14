package request

type AnimeRecordRequest struct {
	AnimeName   string `json:"animeName"`
	AnimeRating int    `json:"animeRating"`
}

type AnimeRecordUpdateRequest struct {
	AnimeRating int `json:"animeRating"`
	AnimeId     int `json:"animeId"`
	BangumiId   int `json:"bangumiId"`
}

type AnimeRecordDeleteRequest struct {
	AnimeId int `json:"animeId"`
}
