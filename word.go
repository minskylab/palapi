package palapi

import "time"

type WordID string

// Word is an struct to wrap the basic features of a word on palapi.
// The ID of the word is the self word.
// synonyms and antonyms are ordered by near to the word, the near is provided by the source.
type Word struct {
	ID          WordID
	LastUpdate  time.Time
	Sources     []SourceID
	Definitions []WordDefinition
	Synonyms    map[int64]WordID
	Antonyms    map[int64]WordID
	Examples    []Sentence
	Frequency   WordFrequency
}
