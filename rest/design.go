package rest

import (
	"github.com/gofiber/fiber"
	log "github.com/sirupsen/logrus"
)

func (s *Service) describe() {
	apiGroup := s.app.Group("/api/v1")

	apiGroup.Get("/all", func(c *fiber.Ctx) {
		log.WithField("ip", c.IP()).Info("new all request")
		word, err := s.parentManager.ReportWord(c.Query("word"))
		if err != nil {
			log.WithField("ip", c.IP()).Error(err)
			_ = c.JSON(map[string]error{"error": err})
			return
		}

		if err := c.JSON(word); err != nil {
			log.WithField("ip", c.IP()).Error(err)
			_ = c.JSON(map[string]error{"error": err})
			return
		}
	})

	apiGroup.Get("/definitions", func(c *fiber.Ctx) {
		log.WithField("ip", c.IP()).Info("new definitions request")
		word, err := s.parentManager.ReportWord(c.Query("word"))
		if err != nil {
			log.WithField("ip", c.IP()).Error(err)
			_ = c.JSON(map[string]error{"error": err})
			return
		}

		if err := c.JSON(word.Definitions); err != nil {
			log.WithField("ip", c.IP()).Error(err)
			_ = c.JSON(map[string]error{"error": err})
			return
		}
	})

	apiGroup.Get("/synonyms", func(c *fiber.Ctx) {
		log.WithField("ip", c.IP()).Info("new synonyms request")
		word, err := s.parentManager.ReportWord(c.Query("word"))
		if err != nil {
			log.WithField("ip", c.IP()).Error(err)
			_ = c.JSON(map[string]error{"error": err})
			return
		}

		if err := c.JSON(word.Synonyms); err != nil {
			log.WithField("ip", c.IP()).Error(err)
			_ = c.JSON(map[string]error{"error": err})
			return
		}
	})

	apiGroup.Get("/antonyms", func(c *fiber.Ctx) {
		log.WithField("ip", c.IP()).Info("new antonyms request")
		word, err := s.parentManager.ReportWord(c.Query("word"))
		if err != nil {
			log.WithField("ip", c.IP()).Error(err)
			_ = c.JSON(map[string]error{"error": err})
			return
		}

		if err := c.JSON(word.Antonyms); err != nil {
			log.WithField("ip", c.IP()).Error(err)
			_ = c.JSON(map[string]error{"error": err})
			return
		}
	})

	apiGroup.Get("/examples", func(c *fiber.Ctx) {
		log.WithField("ip", c.IP()).Info("new examples request")
		word, err := s.parentManager.ReportWord(c.Query("word"))
		if err != nil {
			log.WithField("ip", c.IP()).Error(err)
			_ = c.JSON(map[string]error{"error": err})
			return
		}

		if err := c.JSON(word.Examples); err != nil {
			log.WithField("ip", c.IP()).Error(err)
			_ = c.JSON(map[string]error{"error": err})
			return
		}
	})


}