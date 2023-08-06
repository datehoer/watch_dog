package main

import (
	"crypto/sha1"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"strings"
)

func ensureCacheDirExists() {
	// Check if the cache directory exists
	if _, err := os.Stat("./cache"); os.IsNotExist(err) {
		// If not, create it
		err := os.Mkdir("./cache", 0755)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func getQueryKey(params map[string]string) string {
	// Combine the parameters into a string
	paramStr := fmt.Sprintf("%v", params)
	// Create a hash of the parameters string
	hasher := sha1.New()
	hasher.Write([]byte(paramStr))
	return hex.EncodeToString(hasher.Sum(nil))
}

func main() {
	ensureCacheDirExists()
	r := gin.Default()

	r.GET("/data-entry", func(c *gin.Context) {
		params := map[string]string{
			"db":      c.Query("db"),
			"table":   c.Query("table"),
			"limit":   c.Query("limit"),
			"orderby": c.Query("orderby"),
		}

		queryKey := getQueryKey(params)
		cacheFilePath := fmt.Sprintf("./cache/%s.json", queryKey)

		// Check if update flag is present
		updateFlag := c.Query("update") != ""

		var entries []map[string]string

		if !updateFlag {
			// Try to read from cache
			if cachedData, err := os.ReadFile(cacheFilePath); err == nil {
				if err := json.Unmarshal(cachedData, &entries); err == nil {
					c.JSON(200, entries)
					return
				}
			}
		}

		// Perform the query
		tableNames := strings.Split(params["table"], ",")
		orderBys := strings.Split(params["orderby"], ",")
		db, err := sql.Open("mysql", params["db"])
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		for i, tableName := range tableNames {
			orderBy := orderBys[i]
			query := fmt.Sprintf("SELECT * FROM %s ORDER BY %s DESC LIMIT %s", tableName, orderBy, params["limit"])
			rows, err := db.Query(query)
			if err != nil {
				log.Fatal(err)
			}
			columns, err := rows.Columns()
			if err != nil {
				log.Fatal(err)
			}
			values := make([]sql.RawBytes, len(columns))
			scanArgs := make([]interface{}, len(values))
			for i := range values {
				scanArgs[i] = &values[i]
			}
			for rows.Next() {
				err := rows.Scan(scanArgs...)
				if err != nil {
					log.Fatal(err)
				}
				entry := make(map[string]string)
				for i, value := range values {
					entry[columns[i]] = string(value)
				}
				entries = append(entries, entry)
			}
			rows.Close() // Close rows after reading
		}

		// Save to cache
		encodedData, err := json.Marshal(entries)
		if err != nil {
			log.Fatal(err)
		}
		if err := os.WriteFile(cacheFilePath, encodedData, 0644); err != nil {
			log.Fatal(err)
		}

		c.JSON(200, entries)
	})

	r.Run()
}

//http://localhost:8080/data-entry?db=root:Qq125638/*-@tcp(124.221.222.201:3306)/spider_show&table=book&limit=100&orderby=update_time?