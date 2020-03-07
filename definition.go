package palapi

// TODO: Correct that
type PartOfSpeech string

type WordSpecialDefinition struct {
	Case       string
	Definition string
}

type WordDefinition struct {
	Definition   string
	PartOfSpeech PartOfSpeech
}
