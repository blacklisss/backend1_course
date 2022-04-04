package linkrepo

import (
	"context"
	"fmt"
	"gb/backend1_course/internal/entities/linkentity"
	"net/http"
)

type LinkStore interface {
	Create(ctx context.Context, l string) (*linkentity.Link, error)
	Read(ctx context.Context, hash string, r *http.Request) (*linkentity.Link, error)
	Delete(ctx context.Context, hash string, r *http.Request) error
	Stat(ctx context.Context, hash string) (*linkentity.Link, []*linkentity.IPStat, error)
}

type Links struct {
	lstore LinkStore
}

func NewLinks(lstore LinkStore) *Links {
	return &Links{
		lstore: lstore,
	}
}

func (ls *Links) Create(ctx context.Context, l string) (*linkentity.Link, error) {
	link, err := ls.lstore.Create(ctx, l)
	if err != nil {
		return nil, fmt.Errorf("create link error: %w", err)
	}
	return link, nil
}

func (ls *Links) Read(ctx context.Context, hash string, r *http.Request) (*linkentity.Link, error) {
	l, err := ls.lstore.Read(ctx, hash, r)
	if err != nil {
		return nil, fmt.Errorf("read link error: %w", err)
	}
	return l, nil
}

func (ls *Links) Delete(ctx context.Context, hash string, r *http.Request) (*linkentity.Link, error) {
	l, err := ls.lstore.Read(ctx, hash, r)
	if err != nil {
		return nil, fmt.Errorf("delete link error: %w", err)
	}
	return l, ls.lstore.Delete(ctx, hash, r)
}

func (ls *Links) Stat(ctx context.Context, hash string) (*linkentity.Link, []*linkentity.IPStat, error) {
	l, s, err := ls.lstore.Stat(ctx, hash)
	if err != nil {
		return nil, nil, fmt.Errorf("read link error: %w", err)
	}
	return l, s, nil
}
