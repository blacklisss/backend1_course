package linkentity

import (
	"time"

	"github.com/google/uuid"
)

type IPStat struct {
	LinkID      uuid.UUID
	IPAddr      string
	RequestTime time.Time
}

type Link struct {
	ID        uuid.UUID
	Hash      string
	AdminLink string
	Link      string
	Count     uint64
}
