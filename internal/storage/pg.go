package storage

import (
    "fmt"

    "github.com/you/url-shortener/internal/model"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

type Store struct {
    db *gorm.DB
}

func NewPostgres(dsn string) (*Store, error) {
    g, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        return nil, err
    }
    // Авто-миграция схемы
    if err := g.AutoMigrate(&model.URL{}); err != nil {
        return nil, fmt.Errorf("migrate: %w", err)
    }
    return &Store{db: g}, nil
}

func (s *Store) Create(u *model.URL) error {
    return s.db.Create(u).Error
}

func (s *Store) FindByCode(code string) (*model.URL, error) {
    var u model.URL
    err := s.db.First(&u, "code = ?", code).Error
    return &u, err
}

func (s *Store) IncrementClicks(code string) error {
    return s.db.Model(&model.URL{}).
        Where("code = ?", code).
        UpdateColumn("clicks", gorm.Expr("clicks + ?", 1)).
        Error
}