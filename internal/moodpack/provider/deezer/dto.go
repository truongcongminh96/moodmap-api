package deezer

type searchResponse struct {
	Data []struct {
		Title  string `json:"title"`
		Link   string `json:"link"`
		Artist struct {
			Name string `json:"name"`
		} `json:"artist"`
	} `json:"data"`
}
