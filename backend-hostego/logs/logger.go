package logs

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

var crashLogger *log.Logger
var accessLogger *log.Logger

func InitLogger() {
	// Ensure logs directory exists
	if err := os.MkdirAll("logs", 0755); err != nil {
		log.Fatal("❌ Failed to create logs directory:", err)
	}

	// Main application log
	appLogFile, err := os.OpenFile("logs/backend.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("❌ Failed to open main log file:", err)
	}

	// Crash-specific log with more details
	crashLogFile, err := os.OpenFile("logs/crashes.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("❌ Failed to open crash log file:", err)
	}

	// Access log for requests
	accessLogFile, err := os.OpenFile("logs/access.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("❌ Failed to open access log file:", err)
	}

	// Setup multiple writers
	appMulti := io.MultiWriter(os.Stdout, appLogFile)
	crashMulti := io.MultiWriter(os.Stdout, crashLogFile, appLogFile)
	accessMulti := io.MultiWriter(accessLogFile)

	// Configure main logger
	log.SetOutput(appMulti)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// Configure crash logger
	crashLogger = log.New(crashMulti, "CRASH: ", log.Ldate|log.Ltime|log.Lshortfile)

	// Configure access logger
	accessLogger = log.New(accessMulti, "", log.Ldate|log.Ltime)

	log.Println("✅ Enhanced logging system initialized")
}

// LogCrash logs detailed crash information
func LogCrash(r interface{}, context string, additionalData map[string]interface{}) {
	if crashLogger == nil {
		log.Printf("🚨 CRASH (logger not initialized): %v in %s", r, context)
		return
	}

	// Get caller information
	_, file, line, ok := runtime.Caller(2)
	if ok {
		file = filepath.Base(file)
	}

	crashLogger.Printf("🚨 CRITICAL CRASH DETECTED 🚨")
	crashLogger.Printf("Time: %s", time.Now().Format("2006-01-02 15:04:05"))
	crashLogger.Printf("Context: %s", context)
	crashLogger.Printf("Location: %s:%d", file, line)
	crashLogger.Printf("Panic Reason: %v", r)

	if additionalData != nil {
		crashLogger.Printf("Additional Data:")
		for key, value := range additionalData {
			crashLogger.Printf("  %s: %v", key, value)
		}
	}

	crashLogger.Printf("Stack Trace:")
	// Get stack trace for the crash
	stack := make([]byte, 4096)
	n := runtime.Stack(stack, false)
	if n > 0 {
		crashLogger.Printf("%s", stack[:n])
	}
	crashLogger.Printf("═══════════════════════════════════════════════════")
}

// LogAccess logs request access information
func LogAccess(method, path, userAgent, ip string, status int, duration time.Duration) {
	if accessLogger == nil {
		return
	}

	accessLogger.Printf("%s %s | Status: %d | Duration: %v | IP: %s | UA: %s",
		method, path, status, duration, ip, userAgent)
}

// LogError logs application errors with context
func LogError(err error, context string, additionalData map[string]interface{}) {
	errorMsg := fmt.Sprintf("❌ ERROR in %s: %v", context, err)

	for key, value := range additionalData {
		errorMsg += fmt.Sprintf(" | %s: %v", key, value)
	}

	log.Println(errorMsg)
}

// LogWarning logs warnings with context
func LogWarning(message string, context string, additionalData map[string]interface{}) {
	warningMsg := fmt.Sprintf("⚠️ WARNING in %s: %s", context, message)

	for key, value := range additionalData {
		warningMsg += fmt.Sprintf(" | %s: %v", key, value)
	}

	log.Println(warningMsg)
}
