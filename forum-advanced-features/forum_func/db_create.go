// package forum_func

// import (
// 	"database/sql"
// 	"fmt"
// 	"os"
// 	"strconv"
// )

// func CreateDatabase(databaseName string) {
// 	database, _ := sql.Open("sqlite3", "file:test.s3db?_auth&_auth_user=admin&_auth_pass=admin&_auth_crypt=sha1")
// 	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS people (id INTEGER PRIMARY KEY, firstname TEXT, lastname TEXT)")
// 	statement.Exec()
// 	statement, _ = database.Prepare("INSERT INTO people (firstname, lastname) VALUES (?, ?)")
// 	statement.Exec("John", "Doe")
// 	rows, _ := database.Query("SELECT id, firstname, lastname FROM people")
// 	var id int
// 	var firstname string
// 	var lastname string
// 	for rows.Next() {
// 		rows.Scan(&id, &firstname, &lastname)
// 		fmt.Println(strconv.Itoa(id) + ": " + firstname + " " + lastname)
// 	}
// 	os.Exit(0)
// }

package forum_func

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3" // This is SQL driver
)

func CreateDatabase(databaseName string) {
	// file, err := os.Create(databaseName) // Create SQLite file
	file, err := sql.Open("sqlite3", databaseName)
	if err != nil {
		log.Fatal(err.Error())
	}
	file.Close()
	log.Println(databaseName, "created")

	createTables(databaseName) // Create Database Tables
	addFirstRecords(databaseName)
	fmt.Println("Database created!")

}

func createTables(databaseName string) {
	createUserSQL := `CREATE TABLE users (
		[user_email] TEXT NOT NULL PRIMARY KEY,
		[user_name] TEXT NOT NULL,
		[user_pw] TEXT NOT NULL,
		[user_type] TEXT NOT NULL DEFAULT 'USER',
		[creation_date] timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
		[status] TEXT DEFAULT 'ACTIVE' NOT NULL,
		[status_date] timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
		[login_time] timestamp DEFAULT CURRENT_TIMESTAMP,
		[login_time_prev] timestamp DEFAULT CURRENT_TIMESTAMP
	);` // SQL Statement for Create Table
	ExecSQLScript(databaseName, createUserSQL)

	createLoginMethodSQL := `CREATE TABLE login_method (
		[user_email] TEXT NOT NULL,
		[method] TEXT NOT NULL,
		[creation_date] timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
		[status] TEXT DEFAULT 'ACTIVE' NOT NULL,
		[status_date] timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
		FOREIGN KEY(user_email) REFERENCES users(user_email)
	);` // SQL Statement for Create Table
	ExecSQLScript(databaseName, createLoginMethodSQL)

	createMainTopicSQL := `CREATE TABLE topic_header (
		[topic_id] INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		[user_email] TEXT NOT NULL,
		[topic_name] TEXT NOT NULL,
		[topic_text] TEXT NOT NULL,
		[topic_img] TEXT DEFAULT 'false' NOT NULL,
		[creation_date] timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
		[status] TEXT DEFAULT 'PENDING' NOT NULL,
		[status_text] TEXT,
		[status_date] timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
		FOREIGN KEY(user_email) REFERENCES users(user_email)
	  );` // SQL Statement for Create Table
	ExecSQLScript(databaseName, createMainTopicSQL)

	createMainTopicCategoriesSQL := `CREATE TABLE topic_categories (
		[topic_id] INTEGER NOT NULL,
		[category_name] TEXT NOT NULL,
		[creation_date] timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
		[status] TEXT DEFAULT 'ACTIVE' NOT NULL,
		[status_date] timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
		FOREIGN KEY(topic_id) REFERENCES topic_header(topic_id)
	)`
	ExecSQLScript(databaseName, createMainTopicCategoriesSQL)

	createMainTopicContentSQL := `CREATE TABLE topic_content (
		[content_id] INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		[topic_id] NOT NULL,
		[user_email] TEXT NOT NULL,
		[content_name] TEXT NOT NULL,
		[content_text] TEXT NOT NULL,
		[creation_date] timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
		[status] TEXT DEFAULT 'ACTIVE' NOT NULL,
		[status_date] timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
		FOREIGN KEY(user_email) REFERENCES users(user_email),
		FOREIGN KEY(topic_id) REFERENCES topic_header(topic_id)
	  );` // SQL Statement for Create Table
	ExecSQLScript(databaseName, createMainTopicContentSQL)

	createContentEmotionSQL := `CREATE TABLE topic_emotion (
		[emotion_id] INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		[topic_id] NOT NULL,
		[content_id] INTEGER NULL,
		[user_email] TEXT NOT NULL,
		[emotion] TEXT NOT NULL,
		[emotion_text] TEXT NOT NULL,
		[creation_date] timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
		[status] TEXT DEFAULT 'ACTIVE' NOT NULL,
		[status_date] timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
		FOREIGN KEY(topic_id) REFERENCES topic_header(topic_id),
		FOREIGN KEY(user_email) REFERENCES users(user_email),
		FOREIGN KEY(content_id) REFERENCES topic_content(content_id)
	  );` // SQL Statement for Create Table
	ExecSQLScript(databaseName, createContentEmotionSQL)

	createSetupSQL := `CREATE TABLE setup (
		[setup_type] TEXT NOT NULL,
		[setup_value] TEXT NOT NULL
	  );` // SQL Statement for Create Table
	ExecSQLScript(databaseName, createSetupSQL)

}

func addFirstRecords(databaseName string) {
	userInsert := `insert into users (user_email, user_name, user_pw, user_type) values('admin@admin.admin','Admin','d121150ba9e51fff', 'ADMIN')` // SQL Statement for Create Table
	ExecSQLScript(databaseName, userInsert)
	mainTopicInsert := `insert into topic_header (user_email, topic_name, topic_text, status) values('admin@admin.admin','Welcome','Welcome to our FORUM Project', 'ACTIVE')` // SQL Statement for Create Table
	ExecSQLScript(databaseName, mainTopicInsert)
	topicContent1Insert := `insert into topic_content (topic_id, user_email, content_name, content_text) select topic_id, user_email, 'Have you ever!', 'Have you ever tried to create a FORUM page?' from topic_header` // SQL Statement for Create Table
	ExecSQLScript(databaseName, topicContent1Insert)
	topicContent1Like1 := `insert into topic_emotion (topic_id, content_id, user_email, emotion, emotion_text) select topic_id, null, user_email, 'Like', 'Like' from topic_header` // SQL Statement for Create Table
	ExecSQLScript(databaseName, topicContent1Like1)
	topicContent1Like2 := `insert into topic_emotion (topic_id, content_id, user_email, emotion, emotion_text) select topic_id, content_id, user_email, 'Like', 'Like' from topic_content` // SQL Statement for Create Table
	ExecSQLScript(databaseName, topicContent1Like2)
	mainTopic2Insert := `insert into topic_header (user_email, topic_name, topic_text) values('admin@admin.admin','No content in this topic','Only theme created but no text inside')` // SQL Statement for Create Table
	ExecSQLScript(databaseName, mainTopic2Insert)
	topicContent1DisLike1 := `insert into topic_emotion (topic_id, content_id, user_email, emotion, emotion_text) select topic_id, null, user_email, 'Dislike', 'Dislike' from topic_header where topic_id = 2` // SQL Statement for Create Table
	ExecSQLScript(databaseName, topicContent1DisLike1)
	SetupType1 := `insert into setup (setup_type, setup_value) select 'CATEGORY','GO reloaded' union all select 'CATEGORY','ASCII art' union all select 'CATEGORY','Groupie tracker' union all select 'CATEGORY','Lem-in' union all select 'CATEGORY','Forum'` // SQL Statement for Create Table
	ExecSQLScript(databaseName, SetupType1)
	SetupType2 := `insert into setup (setup_type, setup_value) select 'USER_TYPE', 'USER' union all select 'USER_TYPE', 'MODERATOR' union all select 'USER_TYPE', 'ADMIN'` // SQL Statement for Create Table
	ExecSQLScript(databaseName, SetupType2)
	SetupType3 := `insert into setup (setup_type, setup_value) select 'STATUS', 'ACTIVE' union all select 'STATUS', 'PENDING' union all select 'STATUS', 'DECLINED' union all select 'STATUS', 'DELETED'` // SQL Statement for Create Table
	ExecSQLScript(databaseName, SetupType3)
	topicCategories := `insert into topic_categories (topic_id, category_name) values(1, 'Lem-in'),(1, 'ASCII art')`
	ExecSQLScript(databaseName, topicCategories)
}
