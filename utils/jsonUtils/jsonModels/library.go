package jsonModels

type Game struct {
	IGDBID              int      `json:"igdb_id"`
	Name                string   `json:"name"`
	Description         string   `json:"description"`
	PlayTime            int      `json:"playTime"`
	Achievments         []int    `json:"achievments"`
	Executable          string   `json:"executable"`
	Running             bool     `json:"running"`
	Favorite            bool     `json:"favorite"`
	CoverURL            string   `json:"cover_url"`
	ArtworkUrlList      []string   `json:"artwork_url_list"`
	ScreenshotUrlList   []string   `json:"screenshot_url_list"`
	CoverFilename       string   `json:"cover_filename"`
	ArtworkFilenames    []string `json:"artwork_filenames"`
	ScreenshotFilenames []string `json:"screenshot_filenames"`
}

type Library struct {
	Games map[int]Game `json:"games"`
}

func (l Library) GetLibrary() Library {
	return l
}

func (l *Library) DefaultValues() {
	l.Games = map[int]Game{}
}
