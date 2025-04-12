package search

import (
	"exhibition-launcher/library"
	"math"
	"strings"
)

type SearchManager struct {
	LibraryManager *library.LibraryManager
	SortedIGDBIDs  []int
	IndexedGames   map[string]int
	BKNode         *BKNode
}

var illegalCharacters = []string{
	",",
	"!",
	"?",
	".",
	"$",
	"#",
	"@",
	"%",
	"^",
	"&",
	"*",
	"(",
	")",
	"-",
	"+",
	"=",
	"[",
	"]",
	"/",
	"'",
	";",
	":",
	"<",
	">",
}

func (searchManager *SearchManager) IndexGames() {
	searchManager.SortedIGDBIDs = searchManager.LibraryManager.GetSortedIDs()
	searchManager.IndexedGames = map[string]int{}
	var firstGameName string
	for _, id := range searchManager.SortedIGDBIDs {
		name := strings.ToLower(searchManager.LibraryManager.Library.Games[id].Name)
		for _, character := range illegalCharacters {
			name = strings.Trim(name, character)
		}
		if firstGameName == "" {
			firstGameName = name
		}
		searchManager.IndexedGames[name] = id
		splitNames := strings.Split(name, " ")

		for _, splitName := range splitNames {
			searchManager.IndexedGames[splitName] = id
		}

	}

	root := &BKNode{
		Word:     firstGameName,
		Children: make(map[int]*BKNode),
	}

	for name, _ := range searchManager.IndexedGames {
		if name == firstGameName {
			continue
		}
		root.Add(name)
	}
	searchManager.BKNode = root
}

func (searchManager *SearchManager) SearchByName(name string) []int {
	L := 0
	R := len(searchManager.SortedIGDBIDs) - 1
	arr := searchManager.SortedIGDBIDs
	games := searchManager.IndexedGames
	var result []int

	searchResults := searchManager.BKNode.Search(name, 4)
	if len(searchResults) <= 0 {
		return result
	}

	for _, searchResult := range searchResults {
		for L < R {
			middle := int(math.Floor(float64((L + R) / 2)))
			if arr[middle] == games[searchResult] {
				result = append(result, games[searchResult])
				break
			} else if games[searchResult] > arr[middle] {
				L = middle + 1
			} else {
				R = middle - 1
			}
		}
	}

	return result
}
