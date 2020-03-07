package palapi

import (
	"sort"
	"time"

	"github.com/asdine/storm/v3"
	log "github.com/sirupsen/logrus"
)

func (m *Manager) ReportWord(word string) (*Word, error) {
	w, err := m.Persistence.GetWord(word)
	if err == nil {
		if !time.Now().After(w.LastUpdate.Add(m.MaxAntiquityOfWord)) {
			return w, nil
		}
	}

	if err != storm.ErrNotFound {
		return nil, err
	}

	// word not found or needs update

	// providers sorted
	sort.SliceStable(m.Providers, func(i, j int) bool {
		return m.Providers[j].Source().Relevancy < m.Providers[i].Source().Relevancy
	})

	for i, provider := range m.Providers {
		source := provider.Source()
		log.WithFields(log.Fields{"source": source.Name, "relevancy": source.Relevancy}).Info(" walking")


		provider.AvailableFeatures()
	}



}