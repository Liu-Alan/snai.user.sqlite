package main

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	"io"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func RegUser(c *gin.Context) {
	account := strings.TrimSpace(c.Query("account"))
	password := strings.TrimSpace(c.Query("password"))

	message := "success"
	status := 200

	if account == "" || password == "" {
		message = "Account or password cannot be empty"
		status = 400
	} else {
		db, err := sql.Open("sqlite3", "./db/user.db")
		defer db.Close()
		if err != nil {
			message = "Service error"
			status = 400
		} else {
			qsql := "Select account From user where account=?"
			count := 0
			rows, err := db.Query(qsql, account)
			if err != nil {
				message = "Service error"
				status = 400
			} else {
				for rows.Next() {
					count++
				}
				if count > 0 {
					message = "Account already exists"
					status = 400
				} else {
					fmt.Println(password)
					h := md5.New()
					io.WriteString(h, password)
					pwdmd5 := hex.EncodeToString(h.Sum(nil))
					stmt, err := db.Prepare("Insert Into user(account,password) values(?,?)")
					if err != nil {
						message = "Service error"
						status = 400
					}

					_, err = stmt.Exec(account, pwdmd5)
					if err != nil {
						message = "Service error"
						status = 400
					}

					message = "success"
					status = 200
				}
			}
		}
	}
	c.JSON(200, gin.H{
		"status":  status,
		"message": message,
	})
}

func main() {
	router := gin.Default()
	router.GET("/reguser", RegUser)
	router.Run(":8080")
}
