package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

func main() {
	portFlag := flag.String("p", "1725", "Port to listen on")
	localFlag := flag.Bool("local", false, "Bind to localhost (127.0.0.1) instead of 0.0.0.0")
	certFlag := flag.String("cert", "", "TLS certificate file (enables HTTPS)")
	keyFlag := flag.String("key", "", "TLS key file (enables HTTPS)")
	logFlag := flag.Bool("log", false, "Enable access logging to stdout")
	liteFlag := flag.Bool("lite", false, "Force lightweight UI (no JavaScript)")
	persistDLFlag := flag.Bool("persist-dl", false, "Persist download counts to .dlcounts.json")
	versionFlag := flag.Bool("version", false, "Print version and exit")
	flag.Parse()

	if *versionFlag {
		fmt.Printf("jotasrv v%s\n", Version)
		os.Exit(0)
	}

	var err error
	rootDir, err = os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current directory: %v", err)
	}

	if flag.NArg() > 0 {
		rootDir = flag.Arg(0)
	}

	rootDir, err = filepath.Abs(rootDir)
	if err != nil {
		log.Fatalf("Failed to resolve root directory: %v", err)
	}
	liteMode = *liteFlag

	if *persistDLFlag {
		loadDownloadCounts()
	}

	host := "0.0.0.0"
	if *localFlag {
		host = "127.0.0.1"
	}
	addr := fmt.Sprintf("%s:%s", host, *portFlag)

	var handler http.Handler = http.HandlerFunc(handleRequest)
	if *logFlag {
		handler = newLoggingHandler(handler)
	}

	server := &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: 10 * time.Second,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigCh
		log.Println("Shutting down...")
		cancel()
	}()

	go func() {
		if *certFlag != "" && *keyFlag != "" {
			if err := server.ListenAndServeTLS(*certFlag, *keyFlag); err != http.ErrServerClosed {
				log.Fatalf("Server failed: %v", err)
			}
		} else {
			if err := server.ListenAndServe(); err != http.ErrServerClosed {
				log.Fatalf("Server failed: %v", err)
			}
		}
	}()

	fmt.Printf("Starting server at http://%s\n", addr)
	fmt.Printf("Serving directory: %s\n", rootDir)
	if *certFlag != "" && *keyFlag != "" {
		fmt.Println("TLS enabled")
	}
	if liteMode {
		fmt.Println("Lightweight mode (no JavaScript)")
	}
	fmt.Println("Press Ctrl+C to stop.")

	<-ctx.Done()

	if *persistDLFlag {
		saveDownloadCounts()
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Shutdown error: %v", err)
	}
	fmt.Println("Server stopped.")
}
