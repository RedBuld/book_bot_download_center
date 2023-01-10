package database

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB_Session struct {
	params          *DB_Params
	logger          *log.Logger
	pool            *pgxpool.Pool
	config          *pgxpool.Config
	done            chan bool
	notifyConnClose chan bool
	isReady         bool
}

type DB_Params struct {
	Server             string
	MaxConnectAttempts int
}

const (
	reconnectDelay   = 2 * time.Second
	healthCheckDelay = 2 * time.Second
)

var (
	errAlreadyClosed = errors.New("already closed: not connected to the server")
)

func NewDB(params *DB_Params) *DB_Session {
	session := DB_Session{
		params:          params,
		logger:          log.New(os.Stdout, "", log.LstdFlags),
		done:            make(chan bool),
		notifyConnClose: make(chan bool),
	}

	config, err := pgxpool.ParseConfig(session.params.Server)
	if err != nil {
		panic(err)
	}

	session.logger.Println("DB config valid!")
	session.config = config

	session.logger.Println("DB starting connection")
	go session.handleReconnect()

	return &session
}

func (session *DB_Session) handleReconnect() {
	for {
		session.isReady = false
		session.logger.Println("DB attempting to connect")

		err := session.connect()

		if err != nil {
			session.logger.Printf("DB Error: %+v\n", err)

			select {
			case <-session.done:
				return
			case <-time.After(reconnectDelay):
			}
			continue
		}

		select {
		case <-session.done:
			return
		case <-session.notifyConnClose:
			session.logger.Println("DB connection closed. Reconnecting...")
		}
	}
}

func (session *DB_Session) connect() error {
	pool, err := pgxpool.NewWithConfig(context.Background(), session.config)
	if err != nil {
		return err
	}
	session.pool = pool

	err = session.ping()
	if err != nil {
		return err
	}

	go func() {
		ticker := time.NewTicker(healthCheckDelay)
		defer ticker.Stop()
		for {
			<-ticker.C
			err := session.ping()
			if err != nil {
				session.notifyConnClose <- true
				break
			}
		}
	}()

	session.isReady = true
	session.logger.Println("DB connected!")
	session.logger.Println("DB setup!")

	return nil
}

func (session *DB_Session) ping() error {
	err := session.pool.Ping(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func (session *DB_Session) Close() error {
	session.logger.Println("Stopping DB")
	if !session.isReady {
		return errAlreadyClosed
	}
	session.pool.Close()
	close(session.done)
	close(session.notifyConnClose)
	session.isReady = false
	return nil
}
