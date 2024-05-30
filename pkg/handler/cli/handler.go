package cli

import (
	"flag"
	"news-aggregator/pkg/service"
)

type Handler struct {
	services *service.ArticleService
}

func NewHandler(services *service.ArticleService) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitCommands() {
	// Init CLI arguments
	help := flag.Bool("help", false, "Show all available arguments and their descriptions.")
	sources := flag.String("sources", "", "Select the desired news sources to get the news from. Supported sources: abcnews, bbc, washingtontimes, nbc, usatoday")
	keywords := flag.String("keywords", "", "Specify the keywords to filter the news by.")
	dateStart := flag.String("date-start", "", "Specify the start date to filter the news by (format: YYYY-MM-DD).")
	dateEnd := flag.String("date-end", "", "Specify the end date to filter the news by (format: YYYY-MM-DD).")

	flag.Parse()

	if *help {
		flag.Usage()
		return
	}

	h.Execute(*sources, *keywords, *dateStart, *dateEnd)
}

func (h *Handler) Execute(sources, keywords, dateStart, dateEnd string) {
	articles := h.loadData()
	filteredArticles := h.filterArticles(articles, sources, keywords, dateStart, dateEnd)
	h.printArticles(filteredArticles)
}
