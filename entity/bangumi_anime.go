package entity

type BangumiAnime struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	NameCn string `json:"name_cn"`
	Images struct {
		Large  string `json:"large"`
		Common string `json:"common"`
		Medium string `json:"medium"`
		Small  string `json:"small"`
		Grid   string `json:"grid"`
	} `json:"images"`
}

type BangumiResponse struct {
	Results int            `json:"results"`
	List    []BangumiAnime `json:"list"`
}
