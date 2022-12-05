package request

type AnimeRecordRequest struct {
	AnimeName   string `json:"animeName"`
	AnimeRating int    `json:"animeRating"`
	Commment    string `json:"comment"`
}

type AnimeRecordUpdateRequest struct {
	AnimeRating int    `json:"animeRating"`
	AnimeId     int    `json:"animeId"`
	BangumiId   int    `json:"bangumiId"`
	Comment     string `json:"comment"`
}

type AnimeRecordDeleteRequest struct {
	AnimeId int `json:"animeId"`
}
