package palapi


type SourceID string

type Source struct {
	ID SourceID
	Name string

	URL string
	Metadata map[string]string
}
