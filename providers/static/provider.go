package static

import (
	"github.com/minskylab/palapi"
)

type Provider struct {
	fileType string
	features []palapi.Feature
	status palapi.ProviderStatus
}

func (p *Provider) Source() palapi.Source {
	return palapi.Source{
		ID:        "static",
		Name:      "Static Files",
		Relevancy: 1,
		URL:       "",
		Metadata: map[string]string{
			"type": p.fileType,
		},
	}
}

func (p *Provider) AvailableFeatures() []palapi.Feature {
	return p.features
}

func (p *Provider) FindWord(word string) (*palapi.Report, error) {
	switch p.fileType {
	case "json":
	case "yaml":
	case "xml":
	case "csv":
	default:
		return nil, ErrInvalidStaticFileType
	}
}

func (p *Provider) Status() palapi.ProviderStatus {
	return p.status
}
