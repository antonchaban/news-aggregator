package cli

import (
	"flag"
	"fmt"
	"github.com/antonchaban/news-aggregator/pkg/filter"
	"github.com/antonchaban/news-aggregator/pkg/service"
	"os"
)

type Handler interface {
	execute(f filter.Filters, sortOrder string) error
	initCommands() error
}

// cliHandler is a struct that holds a reference to the articleService.
type cliHandler struct {
	service service.ArticleService
}

func NewHandler(asvc service.ArticleService) (Handler, error) {
	h := &cliHandler{service: asvc}
	err := h.initCommands()
	if err != nil {
		return nil, err
	}
	return h, nil
}

// initCommands initializes the command-line interface (CLI) commands and flags.
// It parses the flags and executes the appropriate command based on the provided flags.
func (h *cliHandler) initCommands() error {
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
		return nil
	}

	err := h.execute(filter.Filters{
		Source:    *sources,
		Keyword:   *keywords,
		StartDate: *dateStart,
		EndDate:   *dateEnd,
	}, *sortOrder)
	if err != nil {
		return err
	}
	return nil
}

// execute loads the articles, filters them based on the provided sources, keywords, and date range,
// and then prints the filtered articles.
func (h *cliHandler) execute(f filter.Filters, sortOrder string) error {
	err := h.service.LoadDataFromFiles()
	if err != nil {
		return err
	}

	filteredArticles, err := h.filterArticles(f)
	if err != nil {
		return err
	}
	sortedArticles := h.sortArticles(filteredArticles, sortOrder)
	h.printArticles(sortedArticles, f)
	return nil
}
