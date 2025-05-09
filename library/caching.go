package library

import (
	"exhibition-launcher/utils"
	"exhibition-launcher/utils/json_utils/json_models"
	"fmt"
	"github.com/google/uuid"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

func StartImageServer() {
	mux := http.DefaultServeMux
	mux.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir(GetImageCachePath()))))
	go func() {
		err := http.ListenAndServe(":34115", mux)
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
	imageURL, err := l.GetImageURL(coverFileName)
	if err != nil {
		if l.Settings.CacheImagesToDisk {
			go l.CheckForCache()
		}
		return coverURL
	}
	return imageURL
}

func (l *LibraryManager) GetAllImageURLs(filenames []string, urls []string) []string {
	var listOfImages []string
	if len(filenames) > 0 {
		for _, filename := range filenames {
			imageURL, err := l.GetImageURL(filename)
			if err != nil {
				if l.Settings.CacheImagesToDisk {
					go l.CheckForCache()
				}
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
	if !ok || !l.Settings.CacheImagesToDisk {
		return "", fmt.Errorf("file not found or caching is turned off defaulting back to https")
	}
	return "http://localhost:34115/images/" + url.QueryEscape(fileName), nil
}

func (l *LibraryManager) CacheAllImages(game *json_models.Game) error {
	var (
		err error
	)

	if game.CoverURL != "" {
		game.CoverFilename, err = l.CacheImageToDisk(game.Name, game.CoverURL)
		if err != nil {
			return err
		}
	}

	err = l.CacheArtworks(game)
	if err != nil {
		return err
	}
	err = l.CacheScreenshots(game)
	if err != nil {
		return err
	}

	return nil
}

func (l *LibraryManager) CacheArtworks(game *json_models.Game) error {
	for _, uri := range game.ArtworkUrlList {
		fileName, err := l.CacheImageToDisk(game.Name, uri)
		game.ArtworkFilenames = append(game.ArtworkFilenames, fileName)
		if err != nil {
			return err
		}
	}
	return nil
}
func (l *LibraryManager) CacheScreenshots(game *json_models.Game) error {
	for _, uri := range game.ScreenshotUrlList {
		fileName, err := l.CacheImageToDisk(game.Name, uri)
		game.ScreenshotFilenames = append(game.ScreenshotFilenames, fileName)
		if err != nil {
			return err
		}
	}
	return nil
}

var isCaching bool

// CheckForCache checks if cache exists if it does not it caches the images
func (l *LibraryManager) CheckForCache() {
	// isCaching to make sure multiple cache checks aren't running
	if isCaching {
		return
	}
	isCaching = true
	defer func() {
		isCaching = false
	}()
	cachePath := GetImageCachePath()
	for _, game := range l.Library.Games {

		// Check for single missing files
		for _, filename := range game.ArtworkFilenames {
			if utils.FileExists(cachePath, filename) {
				continue
			}
			game.ArtworkFilenames = nil
		}

		// Check for single missing files
		for _, filename := range game.ScreenshotFilenames {
			if utils.FileExists(cachePath, filename) {
				continue
			}
			game.ScreenshotFilenames = nil
		}

		// Check for entire lists missing
		if game.ArtworkFilenames != nil &&
			game.ScreenshotFilenames != nil &&
			game.CoverFilename != "" &&
			utils.FileExists(cachePath, game.CoverFilename) {
			continue
		}

		gameData, err := l.ProxyClient.GetMetadataByIGDBID(game.IGDBID)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if gameData.Name == "" {
			fmt.Println("Failed to get game data")
			continue
		}
		fmt.Println("No cache found trying to refetch...")

		if len(game.ScreenshotFilenames) <= 0 || game.ScreenshotFilenames == nil {
			err = l.CacheScreenshots(&game)
			if err != nil {
				fmt.Println(err)
				continue
			}
		}

		if len(game.ArtworkFilenames) <= 0 || game.ArtworkFilenames == nil {
			err = l.CacheArtworks(&game)
			if err != nil {
				fmt.Println(err)
				continue
			}
		}

		if game.CoverFilename == "" || !utils.FileExists(GetImageCachePath(), game.CoverFilename) {
			game.CoverFilename, err = l.CacheImageToDisk(gameData.Name, gameData.GetCoverURL())
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
