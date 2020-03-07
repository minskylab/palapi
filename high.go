package palapi

import (
	"sort"
	"time"

	"github.com/asdine/storm/v3"
	log "github.com/sirupsen/logrus"
)

func (m *Manager) reportWord(word string, deepest int64) (*Word, error) {
	if deepest > m.DeepMaxExploration {
		return nil, nil
	}

	log.WithField("deep", deepest).Infof("reporting: %s", word)

	w, err := m.Persistence.GetWord(word)
	if err == nil && !time.Now().After(w.LastUpdate.Add(m.MaxAntiquityOfWord)) {
		return w, nil
	}

	if err != storm.ErrNotFound {
		return nil, err
	}

	// word not found or needs update
	// providers sorted
	sort.SliceStable(m.Providers, func(i, j int) bool {
		return m.Providers[j].Source().Relevancy < m.Providers[i].Source().Relevancy
	})

	var syntheticWord Word

	definitions := make([]WordDefinition, 0)
	synonyms := make([]WordID, 0)
	antonyms := make([]WordID, 0)
	examples := make([]Sentence, 0)
	frequency := WordFrequency{}

	for _, provider := range m.Providers {
		source := provider.Source()

		log.WithFields(log.Fields{"source": source.Name, "relevancy": source.Relevancy}).Info("walking")
		report := provider.FindWord(word)

		log.WithFields(log.Fields{
			"word": report.Word,
			"extract_duration": report.ExtractionDuration,
			"query_duration": report.QueryDuration,
		}).Info("report")

		for _, feature := range  provider.AvailableFeatures() {
			switch feature {
			case Definitions:
				for _, def := range *report.Definitions {
					definitions = append(definitions, def)
				}
			case Synonyms:
				for _, syn := range *report.Synonyms {
					word, err := m.reportWord(syn, deepest+1)
					if err != nil {
						return nil, err
					}
					if word != nil {
						synonyms = append(synonyms, word.ID)
					}
				}
			case Antonyms:
				for _, ant := range *report.Antonyms {
					word, err := m.reportWord(ant, deepest+1)
					if err != nil {
						return nil, err
					}
					if word != nil {
						antonyms = append(antonyms, word.ID)
					}
				}
			case Frequency:
				frequency = *report.Frequency
			case Examples:
				for _, e := range *report.Examples {
					examples = append(examples, Sentence(e))
				}
			default:
				return nil, ErrFeatureNotExist
			}
		}

		mapSynonyms := map[int64]WordID{}
		for i, s := range synonyms {
			mapSynonyms[int64(i)] = s
		}

		mapAntonyms := map[int64]WordID{}
		for i, s := range antonyms {
			mapAntonyms[int64(i)] = s
		}

		syntheticWord = Word{
			ID:          WordID(word),
			LastUpdate:  time.Now(),
			Source:      provider.Source().ID,
			Definitions: definitions,
			Synonyms:    mapSynonyms,
			Antonyms:    mapAntonyms,
			Examples:    examples,
			Frequency:   frequency,
		}
	}

	return m.Persistence.SaveWord(syntheticWord)
}