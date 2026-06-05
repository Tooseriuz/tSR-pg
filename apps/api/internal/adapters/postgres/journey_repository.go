package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tooseriuz/tsr-pg/apps/api/internal/domain"
	"github.com/tooseriuz/tsr-pg/apps/api/internal/service"
)

type JourneyRepository struct {
	pool *pgxpool.Pool
}

func NewJourneyRepository(pool *pgxpool.Pool) JourneyRepository {
	return JourneyRepository{pool: pool}
}

func (r JourneyRepository) List(ctx context.Context) ([]domain.Journey, error) {
	rows, err := r.pool.Query(ctx, `
		select id, name, timestamp, location, thumbnail
		from journeys
		order by highlight desc, timestamp desc, id desc
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	journeys := make([]domain.Journey, 0)
	for rows.Next() {
		var journey domain.Journey
		if err := rows.Scan(&journey.ID, &journey.Name, &journey.Timestamp, &journey.Location, &journey.Thumbnail); err != nil {
			return nil, err
		}

		journeys = append(journeys, journey)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return journeys, nil
}

func (r JourneyRepository) Get(ctx context.Context, id int64) (domain.JourneyContent, error) {
	var journey domain.JourneyContent
	err := r.pool.QueryRow(ctx, `
		select j.name, j.timestamp, j.created_at, jc.content
		from journeys j
		join journey_contents jc on jc.journey_id = j.id
		where j.id = $1
	`, id).Scan(&journey.Name, &journey.Timestamp, &journey.CreatedAt, &journey.Content)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.JourneyContent{}, service.ErrJourneyNotFound
	}
	if err != nil {
		return domain.JourneyContent{}, err
	}

	return journey, nil
}

func (r JourneyRepository) Create(ctx context.Context, journey domain.CreateJourney) (int64, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	var id int64
	err = tx.QueryRow(ctx, `
		insert into journeys (name, timestamp, location, thumbnail, highlight)
		values ($1, $2, $3, $4, $5)
		returning id
	`, journey.Name, journey.Timestamp, journey.Location, journey.Thumbnail, journey.Highlight).Scan(&id)
	if err != nil {
		return 0, err
	}

	_, err = tx.Exec(ctx, `
		insert into journey_contents (journey_id, content)
		values ($1, $2)
	`, id, journey.Content)
	if err != nil {
		return 0, err
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, err
	}

	return id, nil
}

func (r JourneyRepository) CreateImage(ctx context.Context, path string) (string, error) {
	var id string
	err := r.pool.QueryRow(ctx, `
		insert into journey_images (path)
		values ($1)
		returning id
	`, path).Scan(&id)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (r JourneyRepository) GetImage(ctx context.Context, id string) (domain.JourneyImage, error) {
	var image domain.JourneyImage
	err := r.pool.QueryRow(ctx, `
		select id, path, created_at
		from journey_images
		where id = $1
	`, id).Scan(&image.ID, &image.Path, &image.CreatedAt)
	if err != nil {
		return domain.JourneyImage{}, err
	}

	return image, nil
}
