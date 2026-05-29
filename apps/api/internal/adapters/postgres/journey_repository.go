package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tooseriuz/tsr-pg/apps/api/internal/domain"
)

type JourneyRepository struct {
	pool *pgxpool.Pool
}

func NewJourneyRepository(pool *pgxpool.Pool) JourneyRepository {
	return JourneyRepository{pool: pool}
}

func (r JourneyRepository) List(ctx context.Context) ([]domain.Journey, error) {
	rows, err := r.pool.Query(ctx, `
		select name, timestamp, location, thumbnail
		from journeys
		order by timestamp desc, id desc
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	journeys := make([]domain.Journey, 0)
	for rows.Next() {
		var journey domain.Journey
		if err := rows.Scan(&journey.Name, &journey.Timestamp, &journey.Location, &journey.Thumbnail); err != nil {
			return nil, err
		}

		journeys = append(journeys, journey)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return journeys, nil
}
