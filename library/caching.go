package library

import (
	"exhibition-launcher/igdb"
	"exhibition-launcher/utils"
	"exhibition-launcher/utils/jsonUtils/jsonModels"
	"fmt"
	"github.com/google/uuid"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

func StartImageServer() {
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir(GetImageCachePath()))))
	go func() {
		err := http.ListenAndServe(":34115", nil)
		if err != nil {
			fmt.Println(err)
			return
		}
	}()
	fmt.Println("Image server running at localhost:34115")
}

func GetImageCachePath() string {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		// fallback
		cacheDir = os.TempDir()
	}
	return filepath.Join(cacheDir, "Exhibition-Launcher", "images")
}

func (l *LibraryManager) CacheImageToDisk(gameName string, uri string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return "", err
	}

	resp, err := l.Client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	fileName := url.QueryEscape(gameName + "-" + uuid.New().String() + ".jpg")
	pathToFile := filepath.Join(GetImageCachePath(), fileName)

	if err = os.MkdirAll(filepath.Dir(pathToFile), 0755); err != nil {
		return "", err
	}

	file, err := os.OpenFile(pathToFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
	if err != nil {
		return "", fmt.Errorf("error opening file while caching: %w", err)
	}

	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return "", err
	}

	// Encode the path so weird characters like ', ", ? dont blow things up but exclude /'s for paths
	fmt.Println("Succesfully cached image:", pathToFile)
	return fileName, nil
}

func (l *LibraryManager) GetCoverURL(coverFileName string, coverURL string) string {
	if coverFileName != "" && l.Settings.UseCaching {
		imageURL, _ := l.GetImageURL(coverFileName)
		return imageURL
	}
	return coverURL
}

func (l *LibraryManager) GetAllImageURLs(filenames []string, urls []string) []string {
	var listOfImages []string
	if len(filenames) > 0 {
		for _, filename := range filenames {
			imageURL, err := l.GetImageURL(filename)
			if err != nil {
				listOfImages = urls
				break
			}
			listOfImages = append(listOfImages, imageURL)
		}
	} else {
		listOfImages = urls
	}
	return listOfImages
}

func (l *LibraryManager) GetImageURL(fileName string) (string, error) {
	ok := utils.FileExists(filepath.Join(GetImageCachePath(), fileName))
	if !ok || !l.Settings.UseCaching {
		return "", fmt.Errorf("file not found or caching is turned off defaulting back to https")
	}
	return "http://localhost:34115/images/" + url.QueryEscape(fileName), nil
}

func (l *LibraryManager) CacheAllImages(game *jsonModels.Game, gameData igdb.ApiGame) error {
	var (
		err error
	)

	game.CoverFilename, err = l.CacheImageToDisk(gameData.Name, gameData.CoverURL)
	if err != nil {
		return err
	}

	err = l.CacheArtworks(game, gameData)
	if err != nil {
		return err
	}
	err = l.CacheScreenshots(game, gameData)
	if err != nil {
		return err
	}

	return nil
}

func (l *LibraryManager) CacheArtworks(game *jsonModels.Game, gameData igdb.ApiGame) error {
	for _, uri := range gameData.ArtworkUrlList {
		fileName, err := l.CacheImageToDisk(gameData.Name, uri)
		game.ArtworkFilenames = append(game.ArtworkFilenames, fileName)
		if err != nil {
			return err
		}
	}
	return nil
}
func (l *LibraryManager) CacheScreenshots(game *jsonModels.Game, gameData igdb.ApiGame) error {
	for _, uri := range gameData.ScreenshotUrlList {
		fileName, err := l.CacheImageToDisk(gameData.Name, uri)
		game.ScreenshotFilenames = append(game.ScreenshotFilenames, fileName)
		if err != nil {
			return err
		}
	}
	return nil
}

// CheckForCache checks if cache exists if it does not it caches the images
func (l *LibraryManager) CheckForCache() {
	for _, game := range l.Library.Games {
		if game.ArtworkFilenames != nil && game.ScreenshotFilenames != nil && game.CoverFilename != "" {
			continue
		}
		fmt.Println("No cache found trying to refetch...")

		gameData, err := l.APIManager.GetGameData(game.IGDBID)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if gameData.Name == "" {
			fmt.Println("Failed to get game data", err)
			continue
		}

		if game.ScreenshotFilenames == nil {
			err = l.CacheScreenshots(&game, gameData)
			if err != nil {
				fmt.Println(err)
				continue
			}
		}

		if game.ArtworkFilenames == nil {
			err = l.CacheArtworks(&game, gameData)
			if err != nil {
				fmt.Println(err)
				continue
			}
		}

		if game.CoverFilename == "" {
			game.CoverFilename, err = l.CacheImageToDisk(gameData.Name, gameData.CoverURL)
			if err != nil {
				fmt.Println(err)
				continue
			}
		}

		l.Library.Games[game.IGDBID] = game

		err = l.JsonManager.Save()
		if err != nil {
			fmt.Println("Failed to save data", err)
			continue
		}

	}
}
