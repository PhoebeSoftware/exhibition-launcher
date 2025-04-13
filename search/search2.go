package search

import (
	"exhibition-launcher/library"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"strings"
)

type FuzzyManager struct {
	LibraryManager *library.LibraryManager
	GamesMap       map[string]int
}

func getKeys(m map[string]int) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func (fuzzyManager *FuzzyManager) IndexFuzzy() {
	fuzzyManager.GamesMap = make(map[string]int)
	for id, game := range fuzzyManager.LibraryManager.Library.Games {
		name := strings.ToLower(game.Name)
		// Remove weird characters
		name = strings.Map(func(r rune) rune {
			if strings.ContainsRune(",!?.$#@%^&*()-+=[]/';:<>", r) {
				return -1
			}
			return r
		}, name)
		fuzzyManager.GamesMap[name] = id
	}
}

func (fuzzyManager *FuzzyManager) SearchByName(name string) []int {
	name = strings.ToLower(name)
	// Remove weird characters
	name = strings.Map(func(r rune) rune {
		if strings.ContainsRune(",!?.$#@%^&*()-+=[]/';:<>", r) {
			return -1
		}
		return r
	}, name)

	gameNames := getKeys(fuzzyManager.GamesMap)

	searchResults := fuzzy.Find(name, gameNames)
	var results []int
	for _, searchResult := range searchResults {
		results = append(results, fuzzyManager.GamesMap[searchResult])
	}

	return results
}











