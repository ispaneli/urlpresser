package storage

import (
	"math/rand"
	"sync"
)

const (
	charset   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	keyLength = 6
)

type Store struct {
	shortURLs    map[string]string
	originalURLs map[string]string
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

func NewStorage() *Store {
	return &Store{
		originalURLs: make(map[string]string),
		shortURLs:    make(map[string]string),
	}
}
