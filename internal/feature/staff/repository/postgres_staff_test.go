package repository

import (
	"Test/internal/feature/staff/entity"
	"Test/internal/pkg/provider"
	"context"
	"log"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
)

var dataProvider = func() provider.DataProvider {
	pool, err := setupTestDB("")
	if err != nil {
		log.Fatalf("setup testing db error: %v", err)
		return nil
	}

	return provider.NewPostgresProvider(pool)
}()
var testStaffRepo = NewStaffRepo(dataProvider)

func setupTestDB(DSN string) (*pgxpool.Pool, error) {
	dbConfig, err := pgxpool.ParseConfig(DSN)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), dbConfig)
	if err != nil {
		return nil, err
	}

	return pool, nil
}

func TestCreateStaff(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tests := []*entity.Staff{
		{
			ID:             uuid.New(),
			UserID:         uuid.New(),
			OrganizationID: uuid.New(),
			RoleID:         uuid.New(),
			Position:       "",
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
		{
			ID:             uuid.New(),
			UserID:         uuid.New(),
			OrganizationID: uuid.New(),
			RoleID:         uuid.New(),
			Position:       "",
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
		{
			ID:             uuid.New(),
			UserID:         uuid.New(),
			OrganizationID: uuid.New(),
			RoleID:         uuid.New(),
			Position:       "",
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
	}

	for _, tt := range tests {
		err := testStaffRepo.Create(ctx, tt)
		if err != nil {
			t.Errorf("create staff error: %v", err)
			return
		}
		staff, err := testStaffRepo.Get(ctx, tt.ID)
		if err != nil {
			t.Errorf("get staff error: %v", err)
			return
		}
		assert.Equal(t, staff.ID, tt.ID)
		assert.Equal(t, staff.Position, tt.Position)
		assert.Equal(t, staff.RoleID, tt.RoleID)
		assert.Equal(t, staff.UserID, tt.UserID)
		assert.Equal(t, staff.OrganizationID, tt.OrganizationID)
	}

}
