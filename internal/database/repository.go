package database

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/go-sql-driver/mysql"
)

type Repository interface {
	GetUrlDetail(url string) (UrlDetail, error)
	SetLongUrl(url, longUrl string) (int64, error)
	RecordMeta(urlId int, ip, location, deviceType string) error
	Init() error
}

func NewRepository() Repository {
	return new(MysqlRepo)
}

type MysqlRepo struct {
	db          *sql.DB
	initialized bool
}

// Init Initializes the db variable and establishes a connection to the database
func (m *MysqlRepo) Init() error {
	cfg := mysql.Config{
		User:   os.Getenv("DB_USER"),
		Passwd: os.Getenv("DB_PASSWD"),
		Net:    "tcp",
		Addr:   os.Getenv("DB_ADDR"),
		DBName: os.Getenv("DB_NAME"),
	}
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return err
	}
	if err = db.Ping(); err != nil {
		return err
	}
	m.db = db
	m.initialized = true
	return nil
}

// GetUrlDetail retrives the details of the given shortened url
func (m *MysqlRepo) GetUrlDetail(url string) (UrlDetail, error) {
	if !m.initialized {
		return UrlDetail{}, errors.New("GetUrlDetail: repository is not initialized. Did you forget to call Init()?")
	}
	var u UrlDetail
	row := m.db.QueryRow("SELECT id, url, long_url, visit_count from url_maps where url=?", url)

	err := row.Scan(&u.ID, &u.Url, &u.LongUrl, &u.VisitCount)
	if err != nil {
		if err == sql.ErrNoRows {
			return u, fmt.Errorf("GetUrlDetail: no record with url %s", url)
		}
		return u, fmt.Errorf("GetUrlDetail: retrieving %s failed with %v", url, err)
	}
	return u, err
}

// SetLongUrl creates a new record with the given url and longUrl
func (m *MysqlRepo) SetLongUrl(url, longUrl string) (int64, error) {
	if !m.initialized {
		return 0, errors.New("SetLongUrl: repository is not initialized. Did you forget to call Init()?")
	}
	r, err := m.db.Exec("INSERT INTO url_maps (url, long_url) values (?, ?)", url, longUrl)

	if err != nil {
		return 0, err
	}
	id, err := r.LastInsertId()
	return id, err
}

// RecordMeta creates a new record with the given metadata
func (m *MysqlRepo) RecordMeta(urlId int, ip, location, deviceType string) error {
	if !m.initialized {
		return errors.New("RecordMeta: repository is not initialized. Did you forget to call Init()?")
	}
	_, err := m.db.Exec("INSERT INTO url_meta (url_id, ip, location, device_type) values (?, ?, ?, ?)", urlId, ip, location, deviceType)
	if err != nil {
		return err
	}
	_, err = m.db.Exec("UPDATE url_maps SET visit_count=visit_count+1 WHERE id=?", urlId)
	return err
}
