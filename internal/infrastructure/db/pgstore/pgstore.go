package pgstore

import (
	"context"
	"database/sql"
	"fmt"
	"gb/backend1_course/internal/entities/linkentity"
	"gb/backend1_course/internal/infrastructure/utils"
	"gb/backend1_course/internal/usecases/app/repos/linkrepo"
	"net/http"
	"time"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v4/stdlib" // Postgresql driver
)

var _ linkrepo.LinkStore = &Links{}

type DBPgLink struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
	Hash      string
	AdminLink string
	Link      string
	Count     uint64
}

type DBPgLinkStat struct {
	LinkID      uuid.UUID
	IPAddr      string
	RequestTime time.Time
}

type Links struct {
	db *sql.DB
}

func NewLinks(dsn string) (*Links, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}
	if err != nil {
		db.Close()
		return nil, err
	}
	ls := &Links{
		db: db,
	}
	return ls, nil
}

func (ls *Links) Close() {
	ls.db.Close()
}

func (ls *Links) Create(ctx context.Context, l string) (*linkentity.Link, error) {
	id, _ := uuid.NewUUID()
	hash := utils.GenerateRandomString(10)
	adminLink := utils.GenerateRandomString(12)

	dbu := &DBPgLink{
		ID:        id,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Hash:      hash,
		AdminLink: adminLink,
		Link:      l,
		Count:     0,
	}

	_, err := ls.db.ExecContext(ctx, `INSERT INTO links
	(id, created_at, updated_at, deleted_at, hash, admin_link, link, count)
	values ($1, $2, $3, $4, $5, $6, $7, $8)`,
		dbu.ID,
		dbu.CreatedAt,
		dbu.UpdatedAt,
		nil,
		dbu.Hash,
		dbu.AdminLink,
		dbu.Link,
		dbu.Count,
	)
	if err != nil {
		return nil, err
	}

	ret := &linkentity.Link{
		ID:        dbu.ID,
		Link:      dbu.Link,
		Hash:      dbu.Hash,
		AdminLink: dbu.AdminLink,
		Count:     dbu.Count,
	}

	return ret, nil
}

func (ls *Links) Delete(ctx context.Context, hash string, r *http.Request) error {
	_, err := ls.db.ExecContext(ctx, `UPDATE links SET deleted_at = $2 WHERE hash = $1`,
		hash, time.Now(),
	)
	return err
}

func (ls *Links) Read(ctx context.Context, hash string, r *http.Request) (*linkentity.Link, error) {
	dbu := &DBPgLink{}
	err := ls.db.QueryRowContext(ctx,
		`SELECT id, created_at, updated_at, deleted_at, hash, admin_link, link, count
	FROM links WHERE hash = $1`, hash).Scan(
		&dbu.ID,
		&dbu.CreatedAt,
		&dbu.UpdatedAt,
		&dbu.DeletedAt,
		&dbu.Hash,
		&dbu.AdminLink,
		&dbu.Link,
		&dbu.Count,
	)
	if err != nil {
		return nil, err
	}

	if dbu.DeletedAt != nil {
		return &linkentity.Link{
			ID:        dbu.ID,
			Hash:      dbu.Hash,
			AdminLink: dbu.AdminLink,
			Link:      dbu.Link,
			Count:     dbu.Count,
		}, fmt.Errorf("ссылка удалена")
	}

	ip := r.RemoteAddr

	_, err = ls.db.ExecContext(ctx, `INSERT INTO link_stats
	(link_id, ip)
	values ($1, $2)`,
		dbu.ID,
		ip,
	)
	if err != nil {
		return nil, err
	}

	_, err = ls.db.ExecContext(ctx, `UPDATE links SET count = count + 1 WHERE id=$1`,
		dbu.ID,
	)
	if err != nil {
		return nil, err
	}
	return &linkentity.Link{
		ID:        dbu.ID,
		Hash:      dbu.Hash,
		AdminLink: dbu.AdminLink,
		Link:      dbu.Link,
		Count:     dbu.Count,
	}, nil
}

func (ls *Links) Stat(ctx context.Context, hash string) (*linkentity.Link, []*linkentity.IPStat, error) {
	var stats []*linkentity.IPStat

	dbu := &DBPgLink{}
	stat := &DBPgLinkStat{}

	err := ls.db.QueryRowContext(ctx,
		`SELECT id, created_at, updated_at, deleted_at, hash, admin_link, link, count
	FROM links WHERE admin_link = $1`, hash).Scan(
		&dbu.ID,
		&dbu.CreatedAt,
		&dbu.UpdatedAt,
		&dbu.DeletedAt,
		&dbu.Hash,
		&dbu.AdminLink,
		&dbu.Link,
		&dbu.Count,
	)
	if err != nil {
		return nil, nil, err
	}

	rows, err := ls.db.QueryContext(ctx,
		`SELECT link_id, ip, created_at FROM link_stats WHERE link_id = $1`, dbu.ID)
	if err != nil {
		return nil, nil, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(
			&stat.LinkID,
			&stat.IPAddr,
			&stat.RequestTime,
		)
		if err != nil {
			return nil, nil, err
		}

		tmpStat := &linkentity.IPStat{
			LinkID:      stat.LinkID,
			IPAddr:      stat.IPAddr,
			RequestTime: stat.RequestTime,
		}

		stats = append(stats, tmpStat)
	}

	return &linkentity.Link{
		ID:        dbu.ID,
		Hash:      dbu.Hash,
		AdminLink: dbu.AdminLink,
		Link:      dbu.Link,
		Count:     dbu.Count,
	}, stats, nil
}
