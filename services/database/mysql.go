package database

import (
	"database/sql"
	"errors"
	"log"
	"math"
	"os"
	"strings"

	"github.com/JackPairce/MicroService/services/hashing"
	t "github.com/JackPairce/MicroService/services/types"
	"github.com/go-sql-driver/mysql"
)

type DB struct {
	DB *sql.DB
}

func (s *DB) Connect() error {
	// Capture connection properties.
	cfg := mysql.Config{
		User:                 os.Getenv("DBUSER"),
		Passwd:               os.Getenv("DBPASS"),
		Net:                  "tcp",
		Addr:                 os.Getenv("DBHOST"),
		DBName:               os.Getenv("DBNAME"),
		AllowNativePasswords: true,
	}

	DB, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return err
	}
	pingErr := DB.Ping()
	if pingErr != nil {
		return pingErr
	}
	s.DB = DB
	return nil
}

func (s *DB) CheckUserExistence(user *t.User) error {
	stmt, err := s.DB.Prepare("SELECT COUNT(ID) FROM Users WHERE Name = ?;")
	if err != nil {
		return err
	}
	defer stmt.Close()

	var count int
	err = stmt.QueryRow(user.Name).Scan(&count)
	if err == sql.ErrNoRows {
		return errors.New("user not found")
	} else if err != nil {
		return err
	}
	return nil
}

func (s *DB) CheckUserPassword(user *t.User) (int64, error) {
	stmt, err := s.DB.Prepare("SELECT ID,Password FROM Users WHERE Name =?;")
	if err != nil {
		return int64(math.NaN()), err
	}
	defer stmt.Close()
	var ID int
	var Password string
	err = stmt.QueryRow(user.Name).Scan(&ID, &Password)
	if err == sql.ErrNoRows {
		return int64(math.NaN()), errors.New("user not found")
	} else if err != nil {
		return int64(math.NaN()), err
	}
	if hashing.CheckPasswordHash(user.Password, Password) {
		return int64(ID), nil
	}
	return int64(math.NaN()), errors.New("wrong password")
}

func (s *DB) AddUser(user *t.User) (int64, error) {
	stmt, err := s.DB.Prepare("INSERT INTO Users (Name, Password, PeerAddress) VALUES (?, ?, ?);")
	if err != nil {
		return int64(math.NaN()), err
	}
	defer stmt.Close()
	HashedPass, err := hashing.HashPassword(user.Name)
	if err != nil {
		return int64(math.NaN()), err
	}
	res, err := stmt.Exec(user.Name, HashedPass, user.Peeraddress)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return int64(math.NaN()), errors.New("user already exists")
		}
		return int64(math.NaN()), err
	}
	ID, err := res.LastInsertId()
	if err != nil {
		return int64(math.NaN()), err
	}
	return ID, nil
}

func (s *DB) GetUserPeerAdress(id int64) (string, error) {
	stmt, err := s.DB.Prepare("SELECT PeerAddress FROM Users WHERE Id =?;")
	if err != nil {
		return "", err
	}
	defer stmt.Close()
	var PeerAddress string
	err = stmt.QueryRow(id).Scan(&PeerAddress)
	if err != nil {
		return "", err
	}

	return PeerAddress, nil
}

func (s *DB) UserLogin(id int64, peeraddress string) error {
	stmt, err := s.DB.Prepare("UPDATE Users SET PeerAddress = ?, IsActive = 1 WHERE Id =?;")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(peeraddress, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *DB) UserLogout(id int64) error {
	stmt, err := s.DB.Prepare("UPDATE Users SET PeerAddress = NULL,IsActive = 0 WHERE ID =?;")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}
	return nil
}

func (s *DB) AddFile(file *t.File) error {
	stmt, err := s.DB.Prepare("INSERT INTO Files VALUES (?,?,?,?,?);")
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, err := stmt.Exec(file.Id, file.Filename, file.Size, file.Type, file.Ownerid)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	log.Println("adding", n, "file")
	return nil
}

func (s *DB) SearchFile(id int, searchTerm string, exact bool) (*[]*t.File, error) {
	stmt, err := s.DB.Prepare("SELECT DISTINCT ID, FileName, FileSize, FileType, UserID FROM Files WHERE UserID != ? AND FileName LIKE ?;")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	if !exact {
		searchTerm = "%" + searchTerm + "%"
	}
	rows, err := stmt.Query(id, searchTerm)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var files []*t.File
	for rows.Next() {
		var file t.File
		err = rows.Scan(&file.Id, &file.Filename, &file.Size, &file.Type, &file.Ownerid)
		if err != nil {
			return nil, err
		}
		file.Name = strings.Split(file.Filename, ".")[0]
		files = append(files, &file)
	}
	return &files, nil
}

func (s *DB) GetAllFiles(UserID int) (*[]*t.File, error) {
	stmt, err := s.DB.Prepare("SELECT ID FROM Files WHERE UserID = ?;")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(UserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var files []*t.File
	for rows.Next() {
		var file t.File
		err = rows.Scan(&file.Id)
		if err != nil {
			return nil, err
		}
		files = append(files, &file)
	}
	return &files, nil
}

func (s *DB) DeleteFile(id int) error {
	stmt, err := s.DB.Prepare("DELETE FROM Files WHERE ID =?;")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}
	return nil
}
