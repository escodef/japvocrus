package main

import (
	"japvocrus/internal/anki"
	"japvocrus/internal/csvreader"
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

	base := os.Getenv("JARDIC_URL")
	CSVFile := os.Getenv("CSV_FILE")
	audioFolder := os.Getenv("TTS_OUTPUT_FOLDER")
	ttsEnabled := os.Getenv("CREATE_TTS") != "false"
	output := "out.apkg"

	if audioFolder == "" {
		panic("TTS_OUTPUT_FOLDER not set")
	}

	client, err := dict.NewJardicClient(base)
	if err != nil {
		log.Fatal(err)
	}

	csvr, err := csvreader.Load(CSVFile)
	if err != nil {
		log.Fatal(err)
	}

	var translations []dict.Translation

	for i := range 2 {
		translation, err := client.GetTranslation(csvr[i].Word)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(translation)

		translations = append(translations, *translation)
	}

	if err := anki.GenerateApkg(translations, audioFolder, output, ttsEnabled); err != nil {
		panic(err)
	}
}
