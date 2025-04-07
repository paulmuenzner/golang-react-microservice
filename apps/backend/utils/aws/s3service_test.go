package s3service_test

import (
	"context"
	"testing"

	"meinprojekt/s3service"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/stretchr/testify/assert"
)

const (
	testBucket  = "mein-test-bucket"
	testPrefix  = "testordner/"
	testFile    = testPrefix + "testdatei.txt"
	renamedFile = testPrefix + "umbenannt.txt"
)

// Setup function to initialize S3 service
func setupS3(t *testing.T) s3service.S3Client {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		t.Fatalf("Fehler beim Laden der AWS-Konfiguration: %v", err)
	}
	return s3service.NewS3Service(cfg)
}

func TestS3Operations(t *testing.T) {
	s3 := setupS3(t)
	ctx := context.TODO()

	// 1️⃣ Datei hochladen
	err := s3.Upload(ctx, testBucket, testFile, []byte("Testinhalt für S3"))
	assert.NoError(t, err, "Upload fehlgeschlagen")

	// 2️⃣ Datei umbenennen
	err = s3.Rename(ctx, testBucket, testFile, renamedFile)
	assert.NoError(t, err, "Umbenennen fehlgeschlagen")

	// 3️⃣ Dateien auflisten
	files, err := s3.ListFiles(ctx, testBucket, testPrefix)
	assert.NoError(t, err, "Fehler beim Auflisten der Dateien")
	assert.Contains(t, files, renamedFile, "Datei nicht im Ordner gefunden")

	// 4️⃣ Datei löschen
	err = s3.Delete(ctx, testBucket, renamedFile)
	assert.NoError(t, err, "Löschen fehlgeschlagen")

	// 5️⃣ Ordner verschieben
	err = s3.MoveFolder(ctx, testBucket, testPrefix, "neuerordner/")
	assert.NoError(t, err, "Ordner verschieben fehlgeschlagen")
}
