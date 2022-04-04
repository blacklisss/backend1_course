package handlers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"gb/backend1_course/internal/usecases/app/repos/linkrepo"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type Handlers struct {
	ln *linkrepo.Links
}

func NewHandlers(ln *linkrepo.Links) *Handlers {
	r := &Handlers{
		ln: ln,
	}
	return r
}

// Link defines model for Link.
type Link struct {
	ID        uuid.UUID   `json:"id,omitempty"`
	AdminLink string      `json:"adminlink,omitempty"`
	Count     uint64      `json:"count,omitempty"`
	Hash      string      `json:"hash,omitempty"`
	IPStat    []*LinkStat `json:"ipStat,omitempty"`
	Link      string      `json:"link,omitempty"`
}

// LinkStat defines model for LinkStat.
type LinkStat struct {
	LinkID      *uuid.UUID `json:"id,omitempty"`
	IP          string     `json:"ip,omitempty"`
	RequestTime time.Time  `json:"requestTime,omitempty"`
}

func (rt *Handlers) CreateLink(ctx context.Context, l string) (Link, error) {
	nbu, err := rt.ln.Create(ctx, l)
	if err != nil {
		return Link{}, fmt.Errorf("error when creating: %w", err)
	}

	return Link{
		ID:        nbu.ID,
		Hash:      nbu.Hash,
		AdminLink: nbu.AdminLink,
		Link:      nbu.Link,
	}, nil
}

var ErrLinkNotFound = errors.New("link not found")

// read?uid=...
func (rt *Handlers) GetLink(ctx context.Context, hash string, r *http.Request) (Link, error) {
	nbu, err := rt.ln.Read(ctx, hash, r)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Link{}, ErrLinkNotFound
		}
		return Link{}, fmt.Errorf("error when reading: %w", err)
	}

	return Link{
		ID:        nbu.ID,
		Hash:      nbu.Hash,
		AdminLink: nbu.AdminLink,
		Link:      nbu.Link,
	}, nil
}

func (rt *Handlers) DeleteLink(ctx context.Context, hash string, r *http.Request) (Link, error) {
	nbu, err := rt.ln.Delete(ctx, hash, r)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Link{}, ErrLinkNotFound
		}
		return Link{}, fmt.Errorf("error when reading: %w", err)
	}

	return Link{
		ID:        nbu.ID,
		Hash:      nbu.Hash,
		AdminLink: nbu.AdminLink,
		Link:      nbu.Link,
	}, nil
}

func (rt *Handlers) GetStat(ctx context.Context, hash string) (*Link, error) {
	nbu, nds, err := rt.ln.Stat(ctx, hash)
	if err != nil {
		return nil, err
	}

	n := &Link{
		ID:        nbu.ID,
		Link:      nbu.Link,
		AdminLink: nbu.AdminLink,
		Count:     nbu.Count,
		Hash:      nbu.Hash,
	}

	for _, s := range nds {
		tmp := &LinkStat{
			LinkID:      &s.LinkID,
			IP:          s.IPAddr,
			RequestTime: s.RequestTime,
		}
		n.IPStat = append(n.IPStat, tmp)
	}
	return n, nil
}
