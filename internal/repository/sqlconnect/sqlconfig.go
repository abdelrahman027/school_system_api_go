package sqlconnect

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func ConnectDB(dbName string) (*sql.DB, error) {
	// godotenv.Load()
	db_user := os.Getenv("DB_USER")
	db_password := os.Getenv("DB_PASSWORD")
	db_host := os.Getenv("DB_HOST")
	db_port := os.Getenv("DB_PORT")

	fmt.Println("trying to connect to MariaDB")
	// connectionString := "root:Anapple27@tcp(127.0.0.0:3306)/" + dbName
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", db_user, db_password, db_host, db_port, dbName)
	fmt.Println(connectionString)
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return nil, err
	}
	fmt.Println("Connected  to MariaDB")
	return db, nil
}
