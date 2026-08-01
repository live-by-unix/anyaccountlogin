package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/live-by-unix/anyaccountlogin/pkg/auth"
	"github.com/live-by-unix/anyaccountlogin/pkg/device"
	"github.com/live-by-unix/anyaccountlogin/internal/version"
)

var (
	authManager *auth.AuthManager
)

func main() {
	log.Printf("AnyAccountLogin Daemon v%s starting...", version.GetVersion())
	log.Printf("Platform: %s/%s", runtime.GOOS, runtime.GOARCH)

	// Initialize auth manager
	authManager = auth.NewAuthManager()

	// Setup signal handling for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start HTTP server for local communication
	go startHTTPServer(ctx)

	// Start platform-specific integration
	go startPlatformIntegration(ctx)

	// Wait for shutdown signal
	<-sigChan
	log.Println("Shutting down daemon...")
}

func startHTTPServer(ctx context.Context) {
	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "OK")
	})

	// Version endpoint
	mux.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, version.GetVersion())
	})

	// Validate endpoint
	mux.HandleFunc("/validate", validateHandler)

	// Device info endpoint
	mux.HandleFunc("/device", deviceInfoHandler)

	// Start server on Unix socket or localhost
	server := &http.Server{
		Addr:    "127.0.0.1:58432",
		Handler: mux,
	}

	// Try to use Unix socket on Unix-like systems
	if runtime.GOOS != "windows" {
		socketPath := "/var/run/anyaccountlogin.sock"
		os.Remove(socketPath) // Remove if exists

		listener, err := net.Listen("unix", socketPath)
		if err == nil {
			os.Chmod(socketPath, 0660)
			log.Printf("HTTP server listening on Unix socket: %s", socketPath)
			go func() {
				if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
					log.Printf("HTTP server error: %v", err)
				}
			}()
			return
		}
		log.Printf("Failed to create Unix socket, falling back to TCP: %v", err)
	}

	// Fallback to TCP
	log.Printf("HTTP server listening on %s", server.Addr)
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("HTTP server error: %v", err)
		}
	}()

	// Graceful shutdown
	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	server.Shutdown(shutdownCtx)
}

func validateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	flashDrivePath := r.FormValue("flash_drive")
	password := r.FormValue("password")
	user := r.FormValue("user")

	if flashDrivePath == "" || password == "" {
		http.Error(w, "Missing required parameters", http.StatusBadRequest)
		return
	}

	log.Printf("Validation request for user: %s, flash drive: %s", user, flashDrivePath)

	if err := authManager.ValidateFlashDrive(flashDrivePath, password); err != nil {
		log.Printf("Validation failed: %v", err)
		http.Error(w, "Authentication failed", http.StatusUnauthorized)
		return
	}

	log.Printf("Validation successful for user: %s", user)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "OK")
}

func deviceInfoHandler(w http.ResponseWriter, r *http.Request) {
	deviceInfo, err := device.GetDeviceInfo()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get device info: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"system_uuid":"%s","tpm_id":"%s","platform":"%s"}`,
		deviceInfo.SystemUUID, deviceInfo.TPMID, deviceInfo.Platform)
}

func startPlatformIntegration(ctx context.Context) {
	switch runtime.GOOS {
	case "linux":
		startLinuxIntegration(ctx)
	case "darwin":
		startMacOSIntegration(ctx)
	case "windows":
		startWindowsIntegration(ctx)
	default:
		log.Printf("No platform integration for %s", runtime.GOOS)
	}
}

func startLinuxIntegration(ctx context.Context) {
	log.Println("Starting Linux platform integration")
	// Monitor for PAM requests, flash drive insertion/removal, etc.
	// This would integrate with the PAM module and systemd
}

func startMacOSIntegration(ctx context.Context) {
	log.Println("Starting macOS platform integration")
	// Monitor for loginwindow requests, flash drive insertion/removal, etc.
	// This would integrate with the loginwindow plugin and launchd
}

func startWindowsIntegration(ctx context.Context) {
	log.Println("Starting Windows platform integration")
	// Monitor for Credential Provider requests, flash drive insertion/removal, etc.
	// This would integrate with the Credential Provider and Windows Service
}
