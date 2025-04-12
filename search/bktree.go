package search

import "unicode/utf8"

type BKNode struct {
	Word     string
	Children map[int]*BKNode
}

func (n *BKNode) Add(word string) {
	distance := levenshteinDistance(n.Word, word)
	if child, exists := n.Children[distance]; exists {
		child.Add(word)
	} else {
		n.Children[distance] = &BKNode{Word: word, Children: make(map[int]*BKNode)}
	}
}

func (n *BKNode) Search(word string, maxDistance int) []string {
	results := []string{}
	distance := levenshteinDistance(n.Word, word)
	if distance <= maxDistance {
		results = append(results, n.Word)
	}
	for d, child := range n.Children {
		if distance-maxDistance <= d && d <= distance+maxDistance {
			results = append(results, child.Search(word, maxDistance)...)
		}
	}
	return results
}

func levenshteinDistance(a, b string) int {
	la := utf8.RuneCountInString(a)
	lb := utf8.RuneCountInString(b)
	if la == 0 {
		return lb
	}
	if lb == 0 {
		return la
	}

	matrix := make([][]int, la+1)
	for i := range matrix {
		matrix[i] = make([]int, lb+1)
	}

	for i := 0; i <= la; i++ {
		matrix[i][0] = i
	}
	for j := 0; j <= lb; j++ {
		matrix[0][j] = j
	}

	for i := 1; i <= la; i++ {
		for j := 1; j <= lb; j++ {
			cost := 0
			if []rune(a)[i-1] != []rune(b)[j-1] {
				cost = 1
			}
			matrix[i][j] = minimum(
				matrix[i-1][j]+1,
				matrix[i][j-1]+1,
				matrix[i-1][j-1]+cost,
			)
		}
	}
	return matrix[la][lb]
}

func minimum(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}
