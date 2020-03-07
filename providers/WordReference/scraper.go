package wordreference

import (
	"path"
	"strings"
	"time"

	"github.com/gocolly/colly"
	log "github.com/sirupsen/logrus"
)

func (p *Provider) extractByWord(word string) ([]string, []string, time.Duration, error) {
	c := colly.NewCollector()
	synonyms := make([]string, 0)
	antonyms := make([]string, 0)
	t1 := time.Now()
	c.OnHTML("div#article .trans.clickable ul li", func(e *colly.HTMLElement) {
		// synonyms
		text := strings.TrimSpace(e.Text)
		for _, s := range strings.Split(text, ", ") {
			synonyms = append(synonyms, strings.TrimSpace(s))
		}
	})

	c.OnHTML("div#article .trans.clickable ul ul", func(e *colly.HTMLElement) {
		// antonyms
		text := strings.TrimSpace(e.Text)
		for _, s := range strings.Split(text, ", ") {
			antonyms = append(antonyms, strings.TrimSpace(s))
		}
	})

	log.WithFields(log.Fields{
		"synonyms": synonyms,
		"antonyms": antonyms,
	}).Info(word)

	url := path.Join(p.baseURL, strings.Trim(word, " ,.\\/'"))
	if err := c.Visit(url); err != nil {
		return nil, nil, 0, err
	}

	return synonyms, antonyms, time.Since(t1), nil

}
