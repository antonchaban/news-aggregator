package server

import (
	"context"
	"crypto/tls"
	"github.com/antonchaban/news-aggregator/pkg/handler/web"
	"github.com/antonchaban/news-aggregator/pkg/model"
	"github.com/antonchaban/news-aggregator/pkg/service"
	"github.com/antonchaban/news-aggregator/pkg/storage/inmemory"
	"net/http"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestNewServer(t *testing.T) {
	type args struct {
		certFile string
		keyFile  string
	}
	tests := []struct {
		name string
		args args
		want *Server
	}{
		{
			name: "valid certificate and key files",
			args: args{
				certFile: ".server.crt",
				keyFile:  ".server.key",
			},
			want: &Server{
				certFile: ".server.crt",
				keyFile:  ".server.key",
			},
		},
		{
			name: "empty certificate and key files",
			args: args{
				certFile: "",
				keyFile:  "",
			},
			want: &Server{
				certFile: "",
				keyFile:  "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewServer(tt.args.certFile, tt.args.keyFile); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewServer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_Run(t *testing.T) {
	// Set up in-memory databases and services
	db := inmemory.New()
	srcDb := inmemory.NewSrc()
	articleService := service.New(db)
	sourceService := service.NewSourceService(db, srcDb)
	handler := web.NewHandler(articleService, sourceService)

	// Set up environment variables
	os.Setenv("CERT_FILE", "server.crt")
	os.Setenv("KEY_FILE", "server.key")
	os.Setenv("PORT", "443")
	os.Setenv("SAVES_DIR", "testdata") // Assuming a temporary directory for testing

	tests := []struct {
		name     string
		certFile string
		keyFile  string
		port     string
		handler  http.Handler
		wantErr  bool
	}{
		{
			name:     "run",
			certFile: os.Getenv("CERT_FILE"),
			keyFile:  os.Getenv("KEY_FILE"),
			port:     os.Getenv("PORT"),
			handler:  handler.InitRoutes(),
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				certFile: tt.certFile,
				keyFile:  tt.keyFile,
			}

			serverErr := make(chan error)
			go func() {
				serverErr <- s.Run(tt.port, tt.handler, *handler)
			}()

			time.Sleep(2 * time.Second)

			httpClient := &http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
				},
			}

			resp, err := httpClient.Get("https://localhost:" + tt.port + "/articles")
			if err != nil {
				if !tt.wantErr {
					t.Errorf("Unexpected error: %v", err)
				}
			} else {
				resp.Body.Close()
				if tt.wantErr {
					t.Errorf("Expected error but got none")
				}
			}

			err = s.Shutdown(context.Background(), []model.Article{}, []model.Source{})
			if err != nil {
				t.Fatalf("Failed to shutdown server: %v", err)
			}

			err = <-serverErr
			if (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServer_Shutdown(t *testing.T) {
	db := inmemory.New()
	srcDb := inmemory.NewSrc()
	articleService := service.New(db)
	sourceService := service.NewSourceService(db, srcDb)
	handler := web.NewHandler(articleService, sourceService)
	type args struct {
		ctx      context.Context
		articles []model.Article
		sources  []model.Source
	}

	os.Setenv("CERT_FILE", "server.crt")
	os.Setenv("KEY_FILE", "server.key")
	os.Setenv("PORT", "443")
	os.Setenv("SAVES_DIR", "testdata")

	server := &http.Server{
		Addr:           ":" + os.Getenv("PORT"),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	tests := []struct {
		name    string
		fields  Server
		args    args
		wantErr bool
	}{
		{
			name: "shutdown",
			fields: Server{
				httpServer: server,
				certFile:   os.Getenv("CERT_FILE"),
				keyFile:    os.Getenv("KEY_FILE"),
			},
			args: args{
				ctx: context.Background(),
				articles: []model.Article{
					{Id: 1, Title: "Test Article 1"},
				},
				sources: []model.Source{
					{Id: 1, Name: "Test Source 1"},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				httpServer: tt.fields.httpServer,
				certFile:   tt.fields.certFile,
				keyFile:    tt.fields.keyFile,
			}

			serverErr := make(chan error)
			go func() {
				serverErr <- s.Run(os.Getenv("PORT"), handler.InitRoutes(), *handler)
			}()

			// Allow some time for the server to start
			time.Sleep(2 * time.Second)

			if err := s.Shutdown(tt.args.ctx, tt.args.articles, tt.args.sources); (err != nil) != tt.wantErr {
				t.Errorf("Shutdown() error = %v, wantErr %v", err, tt.wantErr)
			}

			err := <-serverErr
			if (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
