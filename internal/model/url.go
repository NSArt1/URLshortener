package model

import "time"

type URL struct {
    ID          uint      `gorm:"primaryKey"`
    Code        string    `gorm:"uniqueIndex;size:8"`
    OriginalURL string    `gorm:"not null"`
    Clicks      uint      `gorm:"default:0"`
    CreatedAt   time.Time
    UpdatedAt   time.Time
}