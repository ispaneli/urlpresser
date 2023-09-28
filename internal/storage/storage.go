package storage

import (
	"encoding/json"
	"math/rand"
	"os"
	"sync"
)

const (
	charset   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	keyLength = 6
)

type Store struct {
	shortURLs       map[string]string
	originalURLs    map[string]string
	fileStoragePath string
	sync.Mutex
}

func (s *Store) GetShortURL(originalURL string) (shortURL string) {
	s.Lock()
	defer s.Unlock()

	shortURL, ok := s.shortURLs[originalURL]
	if !ok {
		shortURL = s.generateShortURL()
		s.shortURLs[originalURL] = shortURL
		s.originalURLs[shortURL] = originalURL
		s.saveStorageToFile()
	}

	return shortURL
}

func (s *Store) GetOriginURL(shortURL string) (originURL string, exist bool) {
	s.Lock()
	defer s.Unlock()

	originURL, exist = s.originalURLs[shortURL]
	return
}

func (s *Store) generateShortURL() (shortURL string) {
	for {
		shortKey := make([]byte, keyLength)
		for i := range shortKey {
			shortKey[i] = charset[rand.Intn(len(charset))]
		}
		shortURL = string(shortKey)

		if _, ok := s.originalURLs[shortURL]; !ok {
			break
		}
	}

	return
}

func NewStorage(fileStoragePath string) (*Store, error) {
	store := &Store{
		originalURLs:    make(map[string]string),
		shortURLs:       make(map[string]string),
		fileStoragePath: fileStoragePath,
	}

	if fileStoragePath != "" {
		if _, err := os.Stat(fileStoragePath); err == nil {
			data, err := os.ReadFile(fileStoragePath)
			if err != nil {
				return nil, err
			}

			var storedData []map[string]string
			if err := json.Unmarshal(data, &storedData); err != nil {
				return nil, err
			}

			for _, item := range storedData {
				store.originalURLs[item["short_url"]] = item["original_url"]
				store.shortURLs[item["original_url"]] = item["short_url"]
			}
		} else {
			if _, err := os.Create(fileStoragePath); err != nil {
				return nil, err
			}
		}
	}

	return store, nil

}

func (s *Store) saveStorageToFile() error {
	if s.fileStoragePath == "" {
		return nil
	}

	data := make([]map[string]string, 0, len(s.originalURLs))
	for shortURL, originalURL := range s.originalURLs {
		data = append(data, map[string]string{"short_url": shortURL, "original_url": originalURL})
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = os.WriteFile(s.fileStoragePath, jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}
