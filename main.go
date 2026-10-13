package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, addr := range addrs {
		ipNet, ok := addr.(*net.IPNet)
		if !ok || ipNet.IP.IsLoopback() {
			continue
		}
		ip4 := ipNet.IP.To4()
		if ip4 != nil && ip4[0] == 192 && ip4[1] == 168 {
			return ip4.String()
		}
	}
	return ""
}

func main() {
	portFlag := flag.String("p", "1725", "Port to listen on")
	localFlag := flag.Bool("local", false, "Bind to localhost (127.0.0.1) instead of 0.0.0.0")
	certFlag := flag.String("cert", "", "TLS certificate file (enables HTTPS)")
	keyFlag := flag.String("key", "", "TLS key file (enables HTTPS)")
	logFlag := flag.Bool("log", false, "Enable access logging to stdout")
	liteFlag := flag.Bool("lite", false, "Force lightweight UI (no JavaScript)")
	persistDLFlag := flag.Bool("persist-dl", false, "Persist download counts to .dlcounts.json")
	userFlag := flag.String("user", envOr("JOTASRV_USER", "admin"), "Login username (env JOTASRV_USER)")
	passFlag := flag.String("pass", envOr("JOTASRV_PASS", "1234"), "Login password (env JOTASRV_PASS)")
	publicFlag := flag.Bool("public", false, "Serve ./public_html (read-only, no login) at /public/")
	publicDirFlag := flag.String("public-dir", "", "Serve this directory (read-only, no login) at /public/; implies -public")
	noAuthFlag := flag.Bool("no-auth", false, "Disable login")
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
	if st, statErr := os.Stat(rootDir); statErr != nil || !st.IsDir() {
		log.Fatalf("Root directory %s does not exist or is not a directory", rootDir)
	}
	if (*certFlag == "") != (*keyFlag == "") {
		log.Fatalf("-cert and -key must be used together")
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

	// Browsers attach cached basic-auth credentials to cross-site form posts,
	// so reject cross-origin state-changing requests (CSRF).
	var handler http.Handler = http.NewCrossOriginProtection().Handler(http.HandlerFunc(handleRequest))
	if !*noAuthFlag {
		handler = newAuthHandler(handler, *userFlag, *passFlag)
	}
	publicDir := ""
	if *publicFlag || *publicDirFlag != "" {
		src := *publicDirFlag
		if src == "" {
			src = "public_html"
		}
		publicDir, err = filepath.Abs(src)
		if err != nil {
			log.Fatalf("Failed to resolve public directory %s: %v", src, err)
		}
		if st, statErr := os.Stat(publicDir); statErr != nil || !st.IsDir() {
			log.Fatalf("-public needs a directory to serve (%s); use -public-dir to set one", publicDir)
		}
		handler = newPublicHandler(handler, publicDir)
	}
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
	if !*localFlag {
		if localIP := getLocalIP(); localIP != "" {
			scheme := "http"
			if *certFlag != "" && *keyFlag != "" {
				scheme = "https"
			}
			fmt.Printf("Local network: %s://%s:%s\n", scheme, localIP, *portFlag)
		}
	}
	fmt.Printf("Serving directory: %s\n", rootDir)
	if publicDir != "" {
		fmt.Printf("Public content (read-only, no login): %s -> /public/\n", publicDir)
	}
	if *certFlag != "" && *keyFlag != "" {
		fmt.Println("TLS enabled")
	}
	if *noAuthFlag {
		fmt.Println("Login disabled")
	} else {
		fmt.Printf("Login enabled (user: %s)\n", *userFlag)
		if *userFlag == "admin" && *passFlag == "1234" {
			fmt.Println("WARNING: using the default credentials admin:1234, change them with -user/-pass")
		}
	}
	if liteMode {
		fmt.Println("Lightweight mode (no JavaScript)")
	}
	fmt.Println("Press Ctrl+C to stop.")

	<-ctx.Done()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Shutdown error: %v", err)
	}

	if *persistDLFlag {
		saveDownloadCounts()
	}
	fmt.Println("Server stopped.")
}
