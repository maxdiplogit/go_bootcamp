package structs

import (
    "time"
)

// === Reusable building blocks ===

type Identifier struct {
    ID   int64  `json:"id"`
    UUID string `json:"uuid"`
}

type Timestamps struct {
    CreatedAt time.Time  `json:"created_at"`
    UpdatedAt time.Time  `json:"updated_at"`
    DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type AuditTrail struct {
    CreatedBy int64  `json:"created_by"`
    UpdatedBy int64  `json:"updated_by"`
    Version   int    `json:"version"`
}

// === Domain types compose these blocks ===

type User struct {
    UserIdentifier *Identifier
    Timestamps
    AuditTrail
    Email    string `json:"email"`
    FullName string `json:"full_name"`
}

type Order struct {
    OrderIdentifier Identifier
    Timestamps
    AuditTrail
    UserID int64   `json:"user_id"`
    Total  float64 `json:"total"`
}

