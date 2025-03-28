package entity

type WordOfWisdom struct {
	Chapter int    `json:"chapter"`
	Line    int    `json:"line"`
	Text    string `json:"text"`
}
