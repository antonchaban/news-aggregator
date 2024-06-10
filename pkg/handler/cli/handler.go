package cli

import (
	"flag"
	"fmt"
	"log"
	"news-aggregator/pkg/service"
	"os"
)

// Handler is a struct that holds a reference to the ArticleService.
type Handler struct {
	Service service.Article
}

func NewHandler(services service.Article) *Handler {
	return &Handler{Service: services}
}

// InitCommands initializes the command-line interface (CLI) commands and flags.
// It parses the flags and executes the appropriate command based on the provided flags.
func (h *Handler) InitCommands() {
	helpDesc := "Show all available arguments and their descriptions."
	sourcesDesc := "Select the desired news sources to get the news from. Supported sources: abcnews, bbc, washingtontimes, nbc, usatoday"
	keywordsDesc := "Specify the keywords to filter the news by."
	dateStartDesc := "Specify the start date to filter the news by (format: YYYY-MM-DD)."
	dateEndDesc := "Specify the end date to filter the news by (format: YYYY-MM-DD)."
	sortOrderDesc := "Specify the sort order for the news by date (ASC or DESC)."

	help := flag.Bool("help", false, helpDesc)
	sources := flag.String("sources", "", sourcesDesc)
	keywords := flag.String("keywords", "", keywordsDesc)
	dateStart := flag.String("date-start", "", dateStartDesc)
	dateEnd := flag.String("date-end", "", dateEndDesc)
	sortOrder := flag.String("sort-order", "DESC", sortOrderDesc)

	flag.Usage = func() {
		fmt.Printf("Usage of %s:\n", os.Args[0])
		fmt.Printf("  -help\n\t%s\n", helpDesc)
		fmt.Printf("  -sources string\n\t%s\n", sourcesDesc)
		fmt.Printf("  -keywords string\n\t%s\n", keywordsDesc)
		fmt.Printf("  -date-start string\n\t%s\n", dateStartDesc)
		fmt.Printf("  -date-end string\n\t%s\n", dateEndDesc)
		fmt.Printf("  -sort-order string\n\t%s\n", sortOrderDesc)
	}

	flag.Parse()

	if *help {
		flag.Usage()
		return
	}

	h.Execute(*sources, *keywords, *dateStart, *dateEnd, *sortOrder)
}

// Execute loads the articles, filters them based on the provided sources, keywords, and date range,
// and then prints the filtered articles.
func (h *Handler) Execute(sources, keywords, dateStart, dateEnd, sortOrder string) {
	err := h.loadData()
	if err != nil {
		log.Fatalf("Error loading data: %v", err)
	}

	filteredArticles := h.filterArticles(sources, keywords, dateStart, dateEnd)
	sortedArticles := h.sortArticles(filteredArticles, sortOrder)
	h.printArticles(sortedArticles, sources, keywords, dateStart, dateEnd)
}
