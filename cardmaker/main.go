package main

import (
	"japvocrus/internal/anki"
	"japvocrus/internal/dict"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	audioFolder := os.Getenv("TTS_OUTPUT_FOLDER")
	ttsEnabled := os.Getenv("CREATE_TTS") != "false"
	output := "out.apkg"

	if audioFolder == "" {
		panic("TTS_OUTPUT_FOLDER not set")
	}

	translations := []dict.Translation{}

	if err := anki.GenerateApkg(translations, audioFolder, output, ttsEnabled); err != nil {
		panic(err)
	}
}
