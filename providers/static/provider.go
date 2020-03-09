package static

import (
	"gopkg.in/thedevsaddam/gojsonq.v2"
	"os"
	"strconv"
	"strings"

	"github.com/minskylab/palapi"
	"github.com/pkg/errors"
)

type SerialType string

const JSON SerialType = "json"
const YAML SerialType = "yaml"
const XML SerialType = "xml"
const CSV SerialType = "csv"

type WordQuery func(word string, feature palapi.Feature) string

type Provider struct {
	filepath string
	query WordQuery
	fileType SerialType
	features []palapi.Feature
	status palapi.ProviderStatus

	file *os.File

	jsonTape *gojsonq.JSONQ
	yamlTape *gojsonq.JSONQ
	xmlTape *gojsonq.JSONQ
	csvTape *gojsonq.JSONQ
}

func NewProvider(filepath string, features []palapi.Feature, fileType SerialType, query WordQuery) (*Provider, error) {
	file, err := os.OpenFile(filepath, os.O_RDONLY, os.ModeAppend)
	if err != nil {
		return nil, errors.Wrap(err, "invalid file, it cannot open")
	}

	return &Provider{
		query:    query,
		filepath: filepath,
		fileType: fileType,
		features: features,
		status:   palapi.IDLE,
		file: file,
		jsonTape: gojsonq.New().Reader(file),
		yamlTape: gojsonq.New(gojsonq.SetDecoder(&yamlDecoder{})).Reader(file),
		xmlTape: gojsonq.New(gojsonq.SetDecoder(&xmlDecoder{})).Reader(file),
		csvTape: gojsonq.New(gojsonq.SetDecoder(&csvDecoder{})).Reader(file),
	}, nil
}

func (p *Provider) Source() palapi.Source {
	return palapi.Source{
		ID:        "static",
		Name:      "Static Files",
		Relevancy: 1,
		URL:       "",
		Metadata: map[string]string{
			"type": string(p.fileType),
		},
	}
}

func (p *Provider) AvailableFeatures() []palapi.Feature {
	return p.features
}

func (p *Provider) FindWord(word string) (*palapi.Report, error) {
	var tape *gojsonq.JSONQ


	switch p.fileType {
	case JSON:
		tape = p.jsonTape
	case YAML:
		tape = p.yamlTape
	case XML:
		tape = p.xmlTape
	case CSV:
		tape = p.csvTape
	default:
		return nil, ErrInvalidStaticFileType
	}

	report := palapi.Report{}
	for _, feature := range p.features {
		query := p.query(word, feature)
		result, err := tape.FindR(query)
		if err != nil {
			return nil, errors.Wrap(err, "query to your JSON file failed")
		}

		switch feature {
		case palapi.Definitions:
			definitions, err := result.StringSlice()
			if err != nil {
				return nil, errors.Wrap(err, "definitions extraction failed")
			}
			defs := make([]palapi.WordDefinition, 0)
			for _, d := range definitions {
				df := d
				partOfSpeech := ""
				parts := strings.Split(d, "|")
				if len(parts) > 1 {
					df = strings.Trim(parts[0], " /\",'")
					partOfSpeech = strings.Trim(parts[1], " /\",'")
				}
				defs = append(defs, palapi.WordDefinition{
					Definition:   df,
					PartOfSpeech: palapi.PartOfSpeech(partOfSpeech),
				})
			}
			report.Definitions = &defs

		case palapi.Synonyms:
			synonyms, err := result.StringSlice()
			if err != nil {
				return nil, errors.Wrap(err, "synonyms extraction failed")
			}
			report.Synonyms = &synonyms

		case palapi.Antonyms:
			antonyms, err := result.StringSlice()
			if err != nil {
				return nil, errors.Wrap(err, "antonyms extraction failed")
			}
			report.Synonyms = &antonyms

		case palapi.Examples:
			antonyms, err := result.StringSlice()
			if err != nil {
				return nil, errors.Wrap(err, "antonyms extraction failed")
			}
			report.Synonyms = &antonyms

		case palapi.Frequency:
			freq, err := result.String()
			if err != nil {
				return nil, errors.Wrap(err, "frequency extraction failed")
			}

			// expected format: "0.23|24.212|1.54"
			parts := strings.Split(freq, "|")
			if len(parts) != 3 {
				return nil, errors.New("you need 3 parts for a word's frequency: \"Zipf|PerMillion|Diversity\"")
			}

			zipf, err := strconv.ParseFloat(parts[0], 64)
			if err != nil {
				return nil, errors.Wrap(err, "cannot parse the zipfs freq")
			}

			perMillion, err := strconv.ParseFloat(parts[1], 64)
			if err != nil {
				return nil, errors.Wrap(err, "cannot parse the perMillion freq")
			}

			diversity, err := strconv.ParseFloat(parts[2], 64)
			if err != nil {
				return nil, errors.Wrap(err, "cannot parse the diversity freq")
			}


			report.Frequency = &palapi.WordFrequency{
				Zipf:       zipf,
				PerMillion: perMillion,
				Diversity:  diversity,
			}
		}
	}

	return &report, nil
}

func (p *Provider) Status() palapi.ProviderStatus {
	return p.status
}
