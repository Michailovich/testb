package gql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"testb/models"
)

type Resolver struct {
	DB *sql.DB
}

// ============================
//               Query
// ============================

func (r *Resolver) Mains(ctx context.Context, args struct{ IncludeDeleted bool }) ([]*mainResolver, error) {
	sqlQ := `SELECT id,title,sub_id,sub_obj,created_at,updated_at,deleted_at FROM main`
	if !args.IncludeDeleted {
		sqlQ += " WHERE deleted_at IS NULL"
	}
	rows, err := r.DB.QueryContext(ctx, sqlQ)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*mainResolver
	for rows.Next() {
		var m models.Main
		var dt sql.NullTime
		if err := rows.Scan(&m.ID, &m.Title, &m.SubID, &m.SubObj, &m.CreatedAt, &m.UpdatedAt, &dt); err != nil {
			return nil, err
		}
		if dt.Valid {
			m.DeletedAt = &dt.Time
		}
		out = append(out, &mainResolver{m: &m, db: r.DB})
	}
	return out, nil
}

func (r *Resolver) Main(ctx context.Context, args struct {
	ID             int32
	IncludeDeleted bool
}) (*mainResolver, error) {
	sqlQ := `SELECT id,title,sub_id,sub_obj,created_at,updated_at,deleted_at FROM main WHERE id=$1`
	if !args.IncludeDeleted {
		sqlQ += " AND deleted_at IS NULL"
	}
	var m models.Main
	var dt sql.NullTime
	err := r.DB.QueryRowContext(ctx, sqlQ, args.ID).Scan(
		&m.ID, &m.Title, &m.SubID, &m.SubObj, &m.CreatedAt, &m.UpdatedAt, &dt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if dt.Valid {
		m.DeletedAt = &dt.Time
	}
	return &mainResolver{m: &m, db: r.DB}, nil
}

// ============================
//             Mutations
// ============================

type ToolInput struct {
	Title       string
	Description *string
}
type TableInput struct{ Name string }
type ChairInput struct {
	Name string
	Type string
}

type MainInput struct {
	Title string
	Tool  *ToolInput
	Table *TableInput
	Chair *ChairInput
}

type MainUpdateInput struct {
	Title     *string
	DeletedAt *string
}

func (r *Resolver) CreateMain(ctx context.Context, args struct{ Input MainInput }) (*mainResolver, error) {
	in := args.Input

	// ровно один из трех
	cnt := 0
	if in.Tool != nil {
		cnt++
	}
	if in.Table != nil {
		cnt++
	}
	if in.Chair != nil {
		cnt++
	}
	if cnt != 1 {
		return nil, errors.New("provide exactly one of tool/table/chair")
	}

	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// создаем заглушку main
	var mainID int
	err = tx.QueryRowContext(ctx,
		`INSERT INTO main (title,sub_id,sub_obj) VALUES($1,0,'') RETURNING id`,
		in.Title,
	).Scan(&mainID)
	if err != nil {
		return nil, err
	}

	// вставляем в нужную дочернюю таблицу
	var subObj string
	var subID int
	if in.Tool != nil {
		err = tx.QueryRowContext(ctx,
			`INSERT INTO tools (title,description,main_id) VALUES($1,$2,$3) RETURNING id`,
			in.Tool.Title, in.Tool.Description, mainID,
		).Scan(&subID)
		subObj = "TOOL"
	}
	if in.Table != nil {
		err = tx.QueryRowContext(ctx,
			`INSERT INTO "tables" (name,main_id) VALUES($1,$2) RETURNING id`,
			in.Table.Name, mainID,
		).Scan(&subID)
		subObj = "TABLE"
	}
	if in.Chair != nil {
		if in.Chair.Type != "ABC" && in.Chair.Type != "CDE" {
			return nil, fmt.Errorf("invalid chair.type: %s", in.Chair.Type)
		}
		err = tx.QueryRowContext(ctx,
			`INSERT INTO chairs (name,type,main_id) VALUES($1,$2,$3) RETURNING id`,
			in.Chair.Name, in.Chair.Type, mainID,
		).Scan(&subID)
		subObj = "CHAIR"
	}
	if err != nil {
		return nil, err
	}

	// апдейтим main
	_, err = tx.ExecContext(ctx,
		`UPDATE main SET sub_id=$1,sub_obj=$2,updated_at=NOW() WHERE id=$3`,
		subID, subObj, mainID,
	)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	// возвращаем свежий main
	return r.Main(ctx, struct {
		ID             int32
		IncludeDeleted bool
	}{ID: int32(mainID), IncludeDeleted: true})
}

func (r *Resolver) UpdateMain(ctx context.Context, args struct {
	ID    int32
	Input MainUpdateInput
}) (*mainResolver, error) {
	in := args.Input
	set := []string{"updated_at=NOW()"}
	vars := []interface{}{}
	idx := 1

	if in.Title != nil {
		idx++
		set = append(set, fmt.Sprintf("title=$%d", idx))
		vars = append(vars, *in.Title)
	}
	if in.DeletedAt != nil {
		idx++
		dt, err := time.Parse(time.RFC3339, *in.DeletedAt)
		if err != nil {
			return nil, fmt.Errorf("bad deletedAt: %w", err)
		}
		set = append(set, fmt.Sprintf("deleted_at=$%d", idx))
		vars = append(vars, dt)
	}

	sqlQ := fmt.Sprintf(
		"UPDATE main SET %s WHERE id=$1 RETURNING id,title,sub_id,sub_obj,created_at,updated_at,deleted_at",
		strings.Join(set, ","),
	)
	vars = append([]interface{}{args.ID}, vars...)

	row := r.DB.QueryRowContext(ctx, sqlQ, vars...)
	var m models.Main
	var dt sql.NullTime
	if err := row.Scan(
		&m.ID, &m.Title, &m.SubID, &m.SubObj,
		&m.CreatedAt, &m.UpdatedAt, &dt,
	); err != nil {
		return nil, err
	}
	if dt.Valid {
		m.DeletedAt = &dt.Time
	}
	return &mainResolver{m: &m, db: r.DB}, nil
}

func (r *Resolver) DeleteMain(ctx context.Context, args struct{ ID int32 }) (bool, error) {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return false, err
	}
	defer tx.Rollback()

	res, err := tx.ExecContext(ctx,
		`UPDATE main SET deleted_at=NOW() WHERE id=$1 AND deleted_at IS NULL`,
		args.ID,
	)
	if err != nil {
		return false, err
	}
	ra, _ := res.RowsAffected()

	// soft‐delete и в дочерних
	_, _ = tx.ExecContext(ctx, `UPDATE tools    SET deleted_at=NOW() WHERE main_id=$1`, args.ID)
	_, _ = tx.ExecContext(ctx, `UPDATE "tables" SET deleted_at=NOW() WHERE main_id=$1`, args.ID)
	_, _ = tx.ExecContext(ctx, `UPDATE chairs   SET deleted_at=NOW() WHERE main_id=$1`, args.ID)

	if err := tx.Commit(); err != nil {
		return false, err
	}
	return ra > 0, nil
}

// ============================
//           Resolvers
// ============================

type mainResolver struct {
	m  *models.Main
	db *sql.DB
}

func (r *mainResolver) ID() int32         { return r.m.ID }
func (r *mainResolver) Title() string     { return r.m.Title }
func (r *mainResolver) SubObj() string    { return r.m.SubObj }
func (r *mainResolver) SubId() int32      { return r.m.SubID }
func (r *mainResolver) CreatedAt() string { return r.m.CreatedAt.Format(time.RFC3339) }
func (r *mainResolver) UpdatedAt() string { return r.m.UpdatedAt.Format(time.RFC3339) }
func (r *mainResolver) DeletedAt() *string {
	if r.m.DeletedAt == nil {
		return nil
	}
	s := r.m.DeletedAt.Format(time.RFC3339)
	return &s
}

func (r *mainResolver) Tool(ctx context.Context) (*toolResolver, error) {
	if r.m.SubObj != "TOOL" {
		return nil, nil
	}
	var t models.Tool
	var dt sql.NullTime
	err := r.db.QueryRowContext(ctx,
		`SELECT id,title,description,main_id,created_at,updated_at,deleted_at
     FROM tools WHERE id=$1`, r.m.SubID,
	).Scan(&t.ID, &t.Title, &t.Description, &t.MainID, &t.CreatedAt, &t.UpdatedAt, &dt)
	if err != nil {
		return nil, err
	}
	if dt.Valid {
		t.DeletedAt = &dt.Time
	}
	return &toolResolver{t: &t}, nil
}

func (r *mainResolver) Table(ctx context.Context) (*tableResolver, error) {
	if r.m.SubObj != "TABLE" {
		return nil, nil
	}
	var t models.Table
	var dt sql.NullTime
	err := r.db.QueryRowContext(ctx,
		`SELECT id,name,main_id,created_at,updated_at,deleted_at
     FROM "tables" WHERE id=$1`, r.m.SubID,
	).Scan(&t.ID, &t.Name, &t.MainID, &t.CreatedAt, &t.UpdatedAt, &dt)
	if err != nil {
		return nil, err
	}
	if dt.Valid {
		t.DeletedAt = &dt.Time
	}
	return &tableResolver{t: &t}, nil
}

func (r *mainResolver) Chair(ctx context.Context) (*chairResolver, error) {
	if r.m.SubObj != "CHAIR" {
		return nil, nil
	}
	var c models.Chair
	var dt sql.NullTime
	err := r.db.QueryRowContext(ctx,
		`SELECT id,name,type,main_id,created_at,updated_at,deleted_at
     FROM chairs WHERE id=$1`, r.m.SubID,
	).Scan(&c.ID, &c.Name, &c.Type, &c.MainID, &c.CreatedAt, &c.UpdatedAt, &dt)
	if err != nil {
		return nil, err
	}
	if dt.Valid {
		c.DeletedAt = &dt.Time
	}
	return &chairResolver{c: &c}, nil
}

// child‐resolvers

type toolResolver struct{ t *models.Tool }

func (r *toolResolver) ID() int32     { return r.t.ID }
func (r *toolResolver) Title() string { return r.t.Title }
func (r *toolResolver) Description() *string {
	return r.t.Description
}
func (r *toolResolver) MainId() int32     { return r.t.MainID }
func (r *toolResolver) CreatedAt() string { return r.t.CreatedAt.Format(time.RFC3339) }
func (r *toolResolver) UpdatedAt() string { return r.t.UpdatedAt.Format(time.RFC3339) }
func (r *toolResolver) DeletedAt() *string {
	if r.t.DeletedAt == nil {
		return nil
	}
	s := r.t.DeletedAt.Format(time.RFC3339)
	return &s
}

type tableResolver struct{ t *models.Table }

func (r *tableResolver) ID() int32         { return r.t.ID }
func (r *tableResolver) Name() string      { return r.t.Name }
func (r *tableResolver) MainId() int32     { return r.t.MainID }
func (r *tableResolver) CreatedAt() string { return r.t.CreatedAt.Format(time.RFC3339) }
func (r *tableResolver) UpdatedAt() string { return r.t.UpdatedAt.Format(time.RFC3339) }
func (r *tableResolver) DeletedAt() *string {
	if r.t.DeletedAt == nil {
		return nil
	}
	s := r.t.DeletedAt.Format(time.RFC3339)
	return &s
}

type chairResolver struct{ c *models.Chair }

func (r *chairResolver) ID() int32         { return r.c.ID }
func (r *chairResolver) Name() string      { return r.c.Name }
func (r *chairResolver) Type() string      { return r.c.Type }
func (r *chairResolver) MainId() int32     { return r.c.MainID }
func (r *chairResolver) CreatedAt() string { return r.c.CreatedAt.Format(time.RFC3339) }
func (r *chairResolver) UpdatedAt() string { return r.c.UpdatedAt.Format(time.RFC3339) }
func (r *chairResolver) DeletedAt() *string {
	if r.c.DeletedAt == nil {
		return nil
	}
	s := r.c.DeletedAt.Format(time.RFC3339)
	return &s
}
