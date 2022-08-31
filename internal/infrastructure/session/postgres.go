// Copyright (c) 2022 Isaque Veras
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package session

import (
	"database/sql"

	"github.com/Masterminds/squirrel"

	"github.com/isaqueveras/power-sso/pkg/database/postgres"
	"github.com/isaqueveras/power-sso/pkg/oops"
)

// pgSession is the implementation
// of transaction for the session repository
type pgSession struct {
	DB *postgres.DBTransaction
}

// create add session of the user in database
func (pg *pgSession) create(userID, clientIP, userAgent *string) (sessionID *string, err error) {
	if err = pg.DB.Builder.
		Insert("sessions").
		Columns("user_id", "expires_at", "ip", "user_agent").
		Values(userID, squirrel.Expr("NOW() + '15 minutes'"), clientIP, userAgent).
		Suffix(`RETURNING "id"`).
		Scan(&sessionID); err != nil {
		return nil, oops.Err(err)
	}

	if _, err = pg.DB.Builder.
		Update("users").
		Set("number_failed_attempts", 0).
		Set("last_failure_date", nil).
		Where("id = ?", userID).
		Exec(); err != nil && err != sql.ErrNoRows {
		return nil, oops.Err(err)
	}

	return
}

func (pg *pgSession) delete(sessionID *string) (err error) {
	if _, err = pg.DB.Builder.
		Update("sessions").
		Set("deleted_at", squirrel.Expr("NOW()")).
		Where("id = ? AND deleted_at IS NULL", sessionID).
		Exec(); err != nil {
		return oops.Err(err)
	}

	return
}