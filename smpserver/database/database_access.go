/*
 * smp_server.
 * Copyright (C) 2021  Sergio Cruz
 *
 * This program is free software; you can redistribute it and/or
 * modify it under the terms of the GNU General Public License
 * as published by the Free Software Foundation; either version 2
 * of the License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 *  along with this program; if not, write to the Free Software
 *  Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.
 */

package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	logger "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

var connection *sql.DB

func GetConnection() *sql.DB {
	if connection != nil {
		return connection
	}
	var err error
	connection, err = sql.Open("sqlite3", "./data/smp.db")
	if err != nil {
		panic(err)
	}
	return connection
}

func Close() {
	if connection != nil {
		err := connection.Close()
		if err != nil {
			logger.Error("Error closinc database connection", err)
		}
		err = nil
	}
}

func Init() {
	conn := GetConnection()
	var err error
	_, err = conn.Exec("CREATE TABLE users (name varchar(100) primary key, pwd varchar(500))")
	if err != nil {
		logger.Warn(err.Error())
	}
}

func StoreUser(user string, password string) error {
	phash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// insert statement
	conn := GetConnection()
	sqlStatement := "INSERT INTO users (name, pwd) VALUES ($1, $2)"
	_, err = conn.Exec(sqlStatement, user, phash)
	if err != nil {
		return err
	}

	return nil
}

func VerifyUser(user string, password string) error {
	conn := GetConnection()
	var phash sql.NullString
	err := conn.QueryRow("SELECT pwd FROM users WHERE name = ?", user).Scan(&phash)
	if err != nil {
		return err
	}
	return bcrypt.CompareHashAndPassword([]byte(phash.String), []byte(password))
}
