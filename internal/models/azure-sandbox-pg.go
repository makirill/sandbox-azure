package models

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type AzureSandboxPostgres struct {
	db_pool *pgxpool.Pool
}

func InitAzureSandboxesPostgres(connString string) (*AzureSandboxPostgres, error) {
	db_pool, err := pgxpool.Connect(context.Background(), connString)
	if err != nil {
		return nil, err
	}

	// TODO: Where to close it??

	return &AzureSandboxPostgres{
		db_pool: db_pool,
	}, nil

}

func (s *AzureSandboxPostgres) Insert(name string) (string, error) {
	id := ""

	err := s.db_pool.QueryRow(context.Background(), "SELECT public.insert_sandbox($1)", name).Scan(&id)

	return id, err
}

func (s *AzureSandboxPostgres) Delete(id string) (bool, error) {
	ok := false

	err := s.db_pool.QueryRow(context.Background(), "SELECT public.delete_sandbox($1)", id).Scan(&ok)

	return ok, err
}

func (s *AzureSandboxPostgres) GetAll() ([]SandboxDetails, error) {

	rows, err := s.db_pool.Query(context.Background(), "SELECT * FROM public.get_sandbox_all()")
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

	rows, err := s.db_pool.Query(context.Background(), "SELECT * FROM public.get_sandbox_by_name($1)", name)
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

	err := s.db_pool.QueryRow(context.Background(), "SELECT * FROM public.get_sandbox_by_id($1)", id).Scan(
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

	err := s.db_pool.QueryRow(context.Background(), "SELECT public.update_sandbox_expires_at($1, $2)", id, expiresAt).Scan(&ok)

	return ok, err
}

func (s *AzureSandboxPostgres) UpdateStatus(id string, status string) (bool, error) {
	ok := false

	err := s.db_pool.QueryRow(context.Background(), "SELECT public.update_sandbox_status($1, $2)", id, status).Scan(&ok)

	return ok, err
}