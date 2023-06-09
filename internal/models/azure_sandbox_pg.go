package models

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

// Make sure we conform to the SandboxData interface
var _ SandboxData = (*AzureSandboxPostgres)(nil)

type AzureSandboxPostgres struct {
	dbPool *pgxpool.Pool
}

func NewAzureSandboxesPostgres(dbPool *pgxpool.Pool) *AzureSandboxPostgres {

	return &AzureSandboxPostgres{
		dbPool: dbPool,
	}

}

func (s *AzureSandboxPostgres) Insert(name string, expireTime time.Time) (string, error) {
	id := ""

	err := s.dbPool.QueryRow(context.Background(), "SELECT public.insert_sandbox($1, $2)", name, expireTime).Scan(&id)

	return id, err
}

func (s *AzureSandboxPostgres) Delete(id string) (bool, error) {
	ok := false

	err := s.dbPool.QueryRow(context.Background(), "SELECT public.delete_sandbox($1)", id).Scan(&ok)

	return ok, err
}

func (s *AzureSandboxPostgres) GetAll(limit int, offset int) ([]SandboxDetails, error) {

	rows, err := s.dbPool.Query(context.Background(), "SELECT * FROM public.get_sandbox_all($1, $2)", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sandboxes := make([]SandboxDetails, 0)

	for rows.Next() {
		var sandbox SandboxDetails

		err := rows.Scan(&sandbox.UUID, &sandbox.Name, &sandbox.CreatedAt, &sandbox.UpdatedAt, &sandbox.ExpiresAt, &sandbox.Status)
		if err != nil {
			return nil, err
		}

		sandboxes = append(sandboxes, sandbox)
	}

	return sandboxes, nil
}

func (s *AzureSandboxPostgres) GetByName(name string) ([]SandboxDetails, error) {

	rows, err := s.dbPool.Query(context.Background(), "SELECT * FROM public.get_sandbox_by_name($1)", name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sandboxes := make([]SandboxDetails, 0)

	for rows.Next() {
		var sandbox SandboxDetails

		err := rows.Scan(&sandbox.UUID, &sandbox.Name, &sandbox.CreatedAt, &sandbox.UpdatedAt, &sandbox.ExpiresAt, &sandbox.Status)
		if err != nil {
			return nil, err
		}

		sandboxes = append(sandboxes, sandbox)
	}

	return sandboxes, nil
}

func (s *AzureSandboxPostgres) GetByID(id string) (SandboxDetails, error) {

	sandbox := SandboxDetails{}

	err := s.dbPool.QueryRow(context.Background(), "SELECT * FROM public.get_sandbox_by_id($1)", id).Scan(
		&sandbox.UUID,
		&sandbox.Name,
		&sandbox.CreatedAt,
		&sandbox.UpdatedAt,
		&sandbox.ExpiresAt,
		&sandbox.Status)

	return sandbox, err
}

func (s *AzureSandboxPostgres) UpdateExpiration(id string, expiresAt time.Time) (bool, error) {
	ok := false

	err := s.dbPool.QueryRow(context.Background(), "SELECT public.update_sandbox_expires_at($1, $2)", id, expiresAt).Scan(&ok)

	return ok, err
}

func (s *AzureSandboxPostgres) UpdateStatus(id string, status string) (bool, error) {
	ok := false

	err := s.dbPool.QueryRow(context.Background(), "SELECT public.update_sandbox_status($1, $2)", id, status).Scan(&ok)

	return ok, err
}
