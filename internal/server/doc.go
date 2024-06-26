// This package contains the server implementation.
// It is responsible for running the HTTPs server and handling the server lifecycle.
// On it's start, it initializes the sources, loads the articles from the backup file, and saves them to the database.
// This package is used by the main package to start the server and handle the server lifecycle.
//
// Usage example:
//
//	server := server.NewServer("path/to/certFile", "path/to/keyFile")
//	err := server.Run("8080", yourHttpHandler, yourWebHandler)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// On shutdown:
//	err = server.Shutdown(context.Background(), articles)
//	if err != nil {
//	    log.Fatal(err)
//	}
package server
