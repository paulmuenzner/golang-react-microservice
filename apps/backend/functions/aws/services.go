package s3service

import (
	"bytes"

	"context"

	"log"

	"github.com/aws/aws-sdk-go-v2/aws"

	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// Interface für S3-Operationen (Ermöglicht Mocking für Tests)

type S3Client interface {
	Upload(ctx context.Context, bucket, key string, data []byte) error
	Delete(ctx context.Context, bucket, key string) error
	Rename(ctx context.Context, bucket, oldKey, newKey string) error
	MoveFolder(ctx context.Context, bucket, sourcePrefix, destinationPrefix string) error
	ListFiles(ctx context.Context, bucket, prefix string) ([]string, error)
	EmptyBucket(ctx context.Context, bucket string) error
}

// Implementierung der S3-Funktionen

type s3Service struct {
	client   *s3.Client
	uploader *manager.Uploader
}

// Konstruktor für S3-Service

func NewS3Service(cfg aws.Config) S3Client {
	client := s3.NewFromConfig(cfg)
	uploader := manager.NewUploader(client)
	return &s3Service{
		client:   client,
		uploader: uploader,
	}
}

// Hochladen einer Datei in S3

func (s *s3Service) Upload(ctx context.Context, bucket, key string, data []byte) error {
	_, err := s.uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: &bucket,
		Key:    &key,
		Body:   bytes.NewReader(data),
	})
	return err
}

// Löschen einer Datei aus S3

func (s *s3Service) Delete(ctx context.Context, bucket, key string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	return err
}

// Umbenennen einer Datei durch Kopieren und Löschen

func (s *s3Service) Rename(ctx context.Context, bucket, oldKey, newKey string) error {
	_, err := s.client.CopyObject(ctx, &s3.CopyObjectInput{
		Bucket:     &bucket,
		CopySource: aws.String(bucket + "/" + oldKey),
		Key:        &newKey,
	})

	if err != nil {
		return err
	}
	return s.Delete(ctx, bucket, oldKey)

}

// Verschieben eines Ordners innerhalb von S3

func (s *s3Service) MoveFolder(ctx context.Context, bucket, sourcePrefix, destinationPrefix string) error {
	files, err := s.ListFiles(ctx, bucket, sourcePrefix)
	if err != nil {
		return err
	}

	for _, file := range files {
		newPath := destinationPrefix + file[len(sourcePrefix):]
		if err := s.Rename(ctx, bucket, file, newPath); err != nil {
			log.Printf("Fehler beim Verschieben von %s nach %s: %v", file, newPath, err)
			return err
		}
	}
	return nil
}

// Auflisten aller Dateien mit bestimmtem Prefix

func (s *s3Service) ListFiles(ctx context.Context, bucket, prefix string) ([]string, error) {
	var files []string
	paginator := s3.NewListObjectsV2Paginator(s.client, &s3.ListObjectsV2Input{
		Bucket: &bucket,
		Prefix: &prefix,
	})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, obj := range page.Contents {
			files += append(files, *obj.Key)
		}
	}
	return files, nil
}

// Löschen aller Dateien im Bucket

func (s *s3Service) EmptyBucket(ctx context.Context, bucket string) error {
	files, err := s.ListFiles(ctx, bucket, "")
	if err != nil {
		return err
	}

	for _, file := range files {
		if err := s.Delete(ctx, bucket, file); err != nil {
			log.Printf("Fehler beim Löschen von %s: %v", file, err)
			return err
		}
	}
	return nil
}
