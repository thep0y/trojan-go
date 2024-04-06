package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/rs/zerolog/log"

	"github.com/thep0y/trojan-go/common"
	"github.com/thep0y/trojan-go/config"
	"github.com/thep0y/trojan-go/statistic"
	"github.com/thep0y/trojan-go/statistic/memory"
)

const Name = "MYSQL"

type Authenticator struct {
	*memory.Authenticator
	db             *sql.DB
	updateDuration time.Duration
	ctx            context.Context
}

func (a *Authenticator) updater() {
	for {
		for _, user := range a.ListUsers() {
			// swap upload and download for users
			hash := user.Hash()
			sent, recv := user.ResetTraffic()

			s, err := a.db.Exec(
				"UPDATE `users` SET `upload`=`upload`+?, `download`=`download`+? WHERE `password`=?;",
				recv,
				sent,
				hash,
			)
			if err != nil {
				log.Error().
					Err(err).
					Msg("failed to update data to user table")
				continue
			}
			if r, err := s.RowsAffected(); err != nil {
				if r == 0 {
					a.DelUser(hash)
				}
			}
		}
		log.Info().Msg("buffered data has been written into the database")

		// update memory
		rows, err := a.db.Query("SELECT password,quota,download,upload FROM users")
		if err != nil || rows.Err() != nil {
			log.Error().Err(err).Msg("failed to pull data from the database")
			time.Sleep(a.updateDuration)
			continue
		}
		for rows.Next() {
			var hash string
			var quota, download, upload int64
			err := rows.Scan(&hash, &quota, &download, &upload)
			if err != nil {
				log.Error().Err(err).Msg("failed to obtain data from the query result")
				break
			}
			if download+upload < quota || quota < 0 {
				a.AddUser(hash)
			} else {
				a.DelUser(hash)
			}
		}

		select {
		case <-time.After(a.updateDuration):
		case <-a.ctx.Done():
			log.Debug().Msg("MySQL daemon exiting...")
			return
		}
	}
}

func connectDatabase(
	driverName, username, password, ip string,
	port int,
	dbName string,
) (*sql.DB, error) {
	path := strings.Join(
		[]string{
			username,
			":",
			password,
			"@tcp(",
			ip,
			":",
			fmt.Sprintf("%d", port),
			")/",
			dbName,
			"?charset=utf8",
		},
		"",
	)
	return sql.Open(driverName, path)
}

func NewAuthenticator(ctx context.Context) (statistic.Authenticator, error) {
	cfg := config.FromContext(ctx, Name).(*Config)
	db, err := connectDatabase(
		"mysql",
		cfg.MySQL.Username,
		cfg.MySQL.Password,
		cfg.MySQL.ServerHost,
		cfg.MySQL.ServerPort,
		cfg.MySQL.Database,
	)
	if err != nil {
		return nil, common.NewError("Failed to connect to database server").Base(err)
	}
	memoryAuth, err := memory.NewAuthenticator(ctx)
	if err != nil {
		return nil, err
	}
	a := &Authenticator{
		db:             db,
		ctx:            ctx,
		updateDuration: time.Duration(cfg.MySQL.CheckRate) * time.Second,
		Authenticator:  memoryAuth.(*memory.Authenticator),
	}
	go a.updater()
	log.Debug().Msg("mysql authenticator created")
	return a, nil
}

func init() {
	statistic.RegisterAuthenticatorCreator(Name, NewAuthenticator)
}
