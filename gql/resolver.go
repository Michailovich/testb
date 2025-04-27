package gql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"testb/models"
)

type Resolver struct {
	DB *sql.DB
}

// MainResolver resolves Main type fields
type MainResolver struct {
	m  *models.Main
	db *sql.DB
}

func (r *MainResolver) ID() int32 {
	return int32(r.m.ID)
}

func (r *MainResolver) Title() string {
	return r.m.Title
}

func (r *MainResolver) SubID() *int32 {
	if r.m.SubID != nil {
		val := int32(*r.m.SubID)
		return &val
	}
	return nil
}

func (r *MainResolver) SubObj() *string {
	return r.m.SubObj
}

func (r *MainResolver) CreatedAt() string {
	return r.m.CreatedAt.Format(time.RFC3339)
}

func (r *MainResolver) UpdatedAt() string {
	return r.m.UpdatedAt.Format(time.RFC3339)
}

func (r *MainResolver) DeletedAt() *string {
	if r.m.DeletedAt != nil {
		s := r.m.DeletedAt.Format(time.RFC3339)
		return &s
	}
	return nil
}

func (r *MainResolver) Tools(ctx context.Context) ([]*ToolResolver, error) {
	if len(r.m.Tools) > 0 {
		var tools []*ToolResolver
		for _, t := range r.m.Tools {
			tools = append(tools, &ToolResolver{t: &t})
		}
		return tools, nil
	}

	query := `SELECT id, title, description, main_id, created_at, updated_at, deleted_at 
	          FROM tools WHERE main_id = $1 AND deleted_at IS NULL`
	rows, err := r.db.QueryContext(ctx, query, r.m.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tools []*ToolResolver
	for rows.Next() {
		var t models.Tool
		var description sql.NullString
		var deletedAt sql.NullTime

		err := rows.Scan(
			&t.ID,
			&t.Title,
			&description,
			&t.MainID,
			&t.CreatedAt,
			&t.UpdatedAt,
			&deletedAt,
		)
		if err != nil {
			return nil, err
		}

		t.Description = description.String
		if deletedAt.Valid {
			t.DeletedAt = &deletedAt.Time
		}

		tools = append(tools, &ToolResolver{t: &t})
	}
	return tools, nil
}

func (r *MainResolver) Tables(ctx context.Context) ([]*TableResolver, error) {
	if len(r.m.Tables) > 0 {
		var tables []*TableResolver
		for _, t := range r.m.Tables {
			tables = append(tables, &TableResolver{t: &t})
		}
		return tables, nil
	}

	query := `SELECT id, name, main_id, created_at, updated_at, deleted_at 
	          FROM tables WHERE main_id = $1 AND deleted_at IS NULL`
	rows, err := r.db.QueryContext(ctx, query, r.m.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []*TableResolver
	for rows.Next() {
		var t models.Table
		var deletedAt sql.NullTime

		err := rows.Scan(
			&t.ID,
			&t.Name,
			&t.MainID,
			&t.CreatedAt,
			&t.UpdatedAt,
			&deletedAt,
		)
		if err != nil {
			return nil, err
		}

		if deletedAt.Valid {
			t.DeletedAt = &deletedAt.Time
		}

		tables = append(tables, &TableResolver{t: &t})
	}
	return tables, nil
}

func (r *MainResolver) Chairs(ctx context.Context) ([]*ChairResolver, error) {
	if len(r.m.Chairs) > 0 {
		var chairs []*ChairResolver
		for _, c := range r.m.Chairs {
			chairs = append(chairs, &ChairResolver{c: &c})
		}
		return chairs, nil
	}

	query := `SELECT id, name, type, main_id, created_at, updated_at, deleted_at 
	          FROM chairs WHERE main_id = $1 AND deleted_at IS NULL`
	rows, err := r.db.QueryContext(ctx, query, r.m.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chairs []*ChairResolver
	for rows.Next() {
		var c models.Chair
		var deletedAt sql.NullTime

		err := rows.Scan(
			&c.ID,
			&c.Name,
			&c.Type,
			&c.MainID,
			&c.CreatedAt,
			&c.UpdatedAt,
			&deletedAt,
		)
		if err != nil {
			return nil, err
		}

		if deletedAt.Valid {
			c.DeletedAt = &deletedAt.Time
		}

		chairs = append(chairs, &ChairResolver{c: &c})
	}
	return chairs, nil
}

// ToolResolver resolves Tool type fields
type ToolResolver struct {
	t *models.Tool
}

func (r *ToolResolver) ID() int32 {
	return int32(r.t.ID)
}

func (r *ToolResolver) Title() string {
	return r.t.Title
}

func (r *ToolResolver) Description() *string {
	if r.t.Description == "" {
		return nil
	}
	return &r.t.Description
}

func (r *ToolResolver) MainID() int32 {
	return int32(r.t.MainID)
}

func (r *ToolResolver) CreatedAt() string {
	return r.t.CreatedAt.Format(time.RFC3339)
}

func (r *ToolResolver) UpdatedAt() string {
	return r.t.UpdatedAt.Format(time.RFC3339)
}

func (r *ToolResolver) DeletedAt() *string {
	if r.t.DeletedAt != nil {
		s := r.t.DeletedAt.Format(time.RFC3339)
		return &s
	}
	return nil
}

// TableResolver resolves Table type fields
type TableResolver struct {
	t *models.Table
}

func (r *TableResolver) ID() int32 {
	return int32(r.t.ID)
}

func (r *TableResolver) Name() string {
	return r.t.Name
}

func (r *TableResolver) MainID() int32 {
	return int32(r.t.MainID)
}

func (r *TableResolver) CreatedAt() string {
	return r.t.CreatedAt.Format(time.RFC3339)
}

func (r *TableResolver) UpdatedAt() string {
	return r.t.UpdatedAt.Format(time.RFC3339)
}

func (r *TableResolver) DeletedAt() *string {
	if r.t.DeletedAt != nil {
		s := r.t.DeletedAt.Format(time.RFC3339)
		return &s
	}
	return nil
}

// ChairResolver resolves Chair type fields
type ChairResolver struct {
	c *models.Chair
}

func (r *ChairResolver) ID() int32 {
	return int32(r.c.ID)
}

func (r *ChairResolver) Name() string {
	return r.c.Name
}

func (r *ChairResolver) Type() string {
	return r.c.Type
}

func (r *ChairResolver) MainID() int32 {
	return int32(r.c.MainID)
}

func (r *ChairResolver) CreatedAt() string {
	return r.c.CreatedAt.Format(time.RFC3339)
}

func (r *ChairResolver) UpdatedAt() string {
	return r.c.UpdatedAt.Format(time.RFC3339)
}

func (r *ChairResolver) DeletedAt() *string {
	if r.c.DeletedAt != nil {
		s := r.c.DeletedAt.Format(time.RFC3339)
		return &s
	}
	return nil
}

// Query resolvers
func (r *Resolver) Mains(ctx context.Context, args struct {
	IncludeDeleted bool `graphql:"includeDeleted"`
}) ([]*MainResolver, error) {

	query := `SELECT id, title, sub_id, sub_obj, created_at, updated_at, deleted_at FROM main`
	if !args.IncludeDeleted {
		query += " WHERE deleted_at IS NULL"
	}

	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var mains []*MainResolver
	for rows.Next() {
		var m models.Main
		var subID sql.NullInt64
		var subObj sql.NullString
		var deletedAt sql.NullTime

		err := rows.Scan(
			&m.ID,
			&m.Title,
			&subID,
			&subObj,
			&m.CreatedAt,
			&m.UpdatedAt,
			&deletedAt,
		)
		if err != nil {
			return nil, err
		}

		if subID.Valid {
			val := int32(subID.Int64)
			m.SubID = &val
		}
		if subObj.Valid {
			m.SubObj = &subObj.String
		}
		if deletedAt.Valid {
			m.DeletedAt = &deletedAt.Time
		}

		mains = append(mains, &MainResolver{m: &m, db: r.DB})
	}

	return mains, nil
}

func (r *Resolver) Main(ctx context.Context, args struct {
	ID             int32
	IncludeDeleted bool `graphql:"includeDeleted"`
}) (*MainResolver, error) {
	query := `SELECT id, title, sub_id, sub_obj, created_at, updated_at, deleted_at FROM main WHERE id = $1`
	if !args.IncludeDeleted {
		query += " AND deleted_at IS NULL"
	}

	var m models.Main
	var subID sql.NullInt64
	var subObj sql.NullString
	var deletedAt sql.NullTime

	err := r.DB.QueryRowContext(ctx, query, args.ID).Scan(
		&m.ID,
		&m.Title,
		&subID,
		&subObj,
		&m.CreatedAt,
		&m.UpdatedAt,
		&deletedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("main with id %d not found", args.ID)
		}
		return nil, err
	}

	if subID.Valid {
		val := int32(subID.Int64)
		m.SubID = &val
	}
	if subObj.Valid {
		m.SubObj = &subObj.String
	}
	if deletedAt.Valid {
		m.DeletedAt = &deletedAt.Time
	}

	return &MainResolver{m: &m, db: r.DB}, nil
}

// Mutation resolvers
func (r *Resolver) CreateMain(ctx context.Context, args struct{ Input models.MainInput }) (*MainResolver, error) {
	var m models.Main
	now := time.Now()

	query := `INSERT INTO main (title, sub_id, sub_obj, created_at, updated_at) 
              VALUES ($1, $2, $3, $4, $5) 
              RETURNING id, title, sub_id, sub_obj, created_at, updated_at, deleted_at`

	var subID interface{} = nil
	if args.Input.SubID != nil {

		subIDValue := int64(*args.Input.SubID)
		subID = subIDValue
	}

	var subObj interface{} = nil
	if args.Input.SubObj != nil {
		subObj = *args.Input.SubObj
	}

	var returnedSubID sql.NullInt64
	var returnedSubObj sql.NullString
	var deletedAt sql.NullTime

	err := r.DB.QueryRowContext(ctx, query,
		args.Input.Title,
		subID,
		subObj,
		now,
		now,
	).Scan(
		&m.ID,
		&m.Title,
		&returnedSubID,
		&returnedSubObj,
		&m.CreatedAt,
		&m.UpdatedAt,
		&deletedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create main: %w", err)
	}

	if returnedSubID.Valid {
		subIDValue := int32(returnedSubID.Int64)
		m.SubID = &subIDValue
	}

	if returnedSubObj.Valid {
		m.SubObj = &returnedSubObj.String
	}

	if deletedAt.Valid {
		m.DeletedAt = &deletedAt.Time
	}

	return &MainResolver{m: &m, db: r.DB}, nil
}

func (r *Resolver) UpdateMain(ctx context.Context, args struct {
	ID    int32
	Input models.MainUpdateInput
}) (*MainResolver, error) {
	query := "UPDATE main SET updated_at = NOW()"
	params := []interface{}{}
	paramCount := 0

	if args.Input.Title != nil {
		paramCount++
		query += fmt.Sprintf(", title = $%d", paramCount)
		params = append(params, *args.Input.Title)
	}

	if args.Input.SubID != nil {
		paramCount++
		query += fmt.Sprintf(", sub_id = $%d", paramCount)
		params = append(params, *args.Input.SubID)
	}

	if args.Input.SubObj != nil {
		paramCount++
		query += fmt.Sprintf(", sub_obj = $%d", paramCount)
		params = append(params, *args.Input.SubObj)
	}

	if args.Input.DeletedAt != nil {
		paramCount++
		query += fmt.Sprintf(", deleted_at = $%d", paramCount)

		if *args.Input.DeletedAt == "" {
			params = append(params, nil)
		} else {
			deletedAt, err := time.Parse(time.RFC3339, *args.Input.DeletedAt)
			if err != nil {
				return nil, fmt.Errorf("invalid deleted_at format: %v", err)
			}
			params = append(params, deletedAt)
		}
	}

	paramCount++
	query += fmt.Sprintf(" WHERE id = $%d RETURNING id, title, sub_id, sub_obj, created_at, updated_at, deleted_at", paramCount)
	params = append(params, args.ID)

	var m models.Main
	var subID sql.NullInt64
	var subObj sql.NullString
	var deletedAt sql.NullTime

	err := r.DB.QueryRowContext(ctx, query, params...).Scan(
		&m.ID,
		&m.Title,
		&subID,
		&subObj,
		&m.CreatedAt,
		&m.UpdatedAt,
		&deletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("main with id %d not found", args.ID)
		}
		return nil, fmt.Errorf("update failed: %v", err)
	}

	if subID.Valid {
		val := int32(subID.Int64)
		m.SubID = &val
	}
	if subObj.Valid {
		m.SubObj = &subObj.String
	}
	if deletedAt.Valid {
		m.DeletedAt = &deletedAt.Time
	}

	return &MainResolver{m: &m, db: r.DB}, nil
}

func (r *Resolver) DeleteMain(ctx context.Context, args struct{ ID int32 }) (bool, error) {
	result, err := r.DB.ExecContext(ctx,
		"UPDATE main SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL",
		args.ID)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected > 0, nil
}
