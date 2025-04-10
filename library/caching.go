package library

import (
	"encoding/base64"
	"exhibition-launcher/igdb"
	"exhibition-launcher/utils/jsonUtils/jsonModels"
	"fmt"
	"github.com/google/uuid"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func getImageCachePath() string {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		// fallback
		cacheDir = os.TempDir()
	}
	return filepath.Join(cacheDir, "exhibtion-launcher", "images")
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
	fileName := gameName + "-" + uuid.New().String() + ".jpg"
	pathToFile := filepath.Join(getImageCachePath(), fileName)

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

func (l *LibraryManager) GetImageBase64(filename string) (string, error) {
	path := filepath.Join(getImageCachePath(), filename)
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	base64Data := base64.StdEncoding.EncodeToString(data)
	return "data:image/png;base64," + base64Data, nil
}

func (l *LibraryManager) CacheAllImagesAndChangePaths(game *jsonModels.Game, gameData igdb.ApiGame) error {
	var (
		err error
	)

	game.CoverURL, err = l.CacheImageToDisk(gameData.Name, gameData.CoverURL)
	if err != nil {
		return err
	}

	for i, uri := range game.ArtworkUrlList {
		game.ArtworkUrlList[i], err = l.CacheImageToDisk(gameData.Name , uri)
		if err != nil {
			return err
		}
	}

	for i, uri := range game.ScreenshotUrlList {
		game.ScreenshotUrlList[i], err = l.CacheImageToDisk(gameData.Name, uri)
		if err != nil {
			return err
		}
	}
	return nil
}

func encodePathSegments(path string) string {
	// Split the path into segments using "/" as a delimiter.
	segments := strings.Split(path, "/")

	// Encode each segment.
	for i, segment := range segments {
		segments[i] = url.PathEscape(segment)
	}

	// Join the encoded segments using "/" as the delimiter.
	return strings.Join(segments, "/")
}
