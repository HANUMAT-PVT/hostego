package logs

import (
	"io"
	"log"
	"os"
)

func InitLogger() {
	file, err := os.OpenFile("logs/backend.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("❌ Failed to open log file:", err)
	}

	multi := io.MultiWriter(os.Stdout, file)
	log.SetOutput(multi)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}
