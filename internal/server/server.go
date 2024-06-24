package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/antonchaban/news-aggregator/pkg/backuper"
	"github.com/antonchaban/news-aggregator/pkg/handler/web"
	"github.com/antonchaban/news-aggregator/pkg/model"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
	certFile   string
	keyFile    string
}

func (s *Server) Run(port string, handler http.Handler, artHandler web.Handler) error {
	s.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}
	articles, err := backuper.NewLoader(artHandler.ArticleService()).LoadAllFromFile()
	if err != nil {
		return err
	}
	err = artHandler.ArticleService().SaveAll(articles)

	if err := s.httpServer.ListenAndServeTLS(s.certFile, s.keyFile); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logrus.Fatalf("Could not listen on %s: %v\n", port, err)
		return err
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context, articles []model.Article) error {
	fmt.Println("Shutting down the server...")
	err := backuper.NewSaver(articles).SaveAllToFile()
	if err != nil {
		return err
	}
	return s.httpServer.Shutdown(ctx)
}

func NewServer(certFile, keyFile string) *Server {
	return &Server{
		certFile: certFile,
		keyFile:  keyFile,
	}
}
