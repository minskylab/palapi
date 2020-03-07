package palapi

type SourceID string

type Source struct {
	ID        SourceID
	Name      string
	Relevancy float64
	URL       string
	Metadata  map[string]string
}
