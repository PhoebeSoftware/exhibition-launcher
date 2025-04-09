package library

import (
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

func (l *LibraryManager) CacheImageToDisk(gameName string, cachingPath string, uri string) (string, error) {
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
	pathToFile := filepath.Join(cachingPath, gameName, fileName)

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

	// Change the path so it is relative to the frontend
	// In ./frontend
	// for example: ./cache/{game}/image.jpg
	relativePath := filepath.Join("..", "cache", gameName, fileName)
	// Encode the path so weird characters like ', ", ? dont blow things up but exclude /'s for paths
	encodedPath := encodePathSegments(relativePath)

	fmt.Println("Succesfully cached image:", fileName)

	return encodedPath, nil
}

func (l *LibraryManager) CacheAllImagesAndChangePaths(game *jsonModels.Game, gameData igdb.ApiGame) error {
	pathToCache := filepath.Join("./frontend/src/cache")

	var (
		err error
	)

	game.CoverURL, err = l.CacheImageToDisk(gameData.Name, pathToCache, gameData.CoverURL)
	if err != nil {
		return err
	}

	for i, uri := range game.ArtworkUrlList {
		game.ArtworkUrlList[i], err = l.CacheImageToDisk(gameData.Name, pathToCache, uri)
		if err != nil {
			return err
		}
	}

	for i, uri := range game.ScreenshotUrlList {
		game.ScreenshotUrlList[i], err = l.CacheImageToDisk(gameData.Name, pathToCache, uri)
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
