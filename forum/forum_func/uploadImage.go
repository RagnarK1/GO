package forum_func

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func uploadImg(w http.ResponseWriter, r *http.Request) (string, string, error) {
	// retrieve file from posted form-data
	file, handler, err := r.FormFile("uploadImg")
	if err != nil {
		return "false", "", errors.New("")
	}
	defer file.Close()

	// only allow files 20MB or under
	maxFileSize := int64(20971520)
	if handler.Size > maxFileSize {
		return "", "File too large to upload. Max size is 20mb", errors.New("file size exceeds 20mb")
	}
	// only allow jpeg, png and gif fileextentions
	ext := filepath.Ext(handler.Filename)
	if ext != ".jpeg" && ext != ".png" && ext != ".gif" && ext != ".jpg" {
		return "", "File extention is not .jpeg, .png or .gif", errors.New("file extention not supported")
	}

	// check if upload_images dir exists and create if it doesn't
	path := "common/uploaded_images"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			fmt.Println(err)
			return "false", "Error creating dir: " + path, err
		}
	}
	// create temp file on server
	filename := strings.Split(handler.Filename, ".")[0]
	tempFile, err := ioutil.TempFile(path, filename+"-*"+ext)
	if err != nil {
		return "", "Error creating temp file on the server", err
	}
	defer tempFile.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return "", "Error reading bytes into file", err
	}
	tempFile.Write(fileBytes)
	return strings.Split(tempFile.Name(), "/")[2], "", nil
}

func maybeDeleteImg(id, remove, img string) string {
	if img != "false" {
		// delete previous img
		deleteImg(id)
	}
	if remove == "on" {
		deleteImg(id)
		if img == "false" {
			img = "remove"
		}
	}
	return img
}

// takes topic id and deletes the image associated with that post
func deleteImg(id string) {
	var deleteImg string
	sqlScript := `select topic_img from topic_header where topic_id = ` + id
	value := SelectSQLScript(DatabaseName, sqlScript)
	var topic_img sql.NullString
	for value.Next() {
		value.Scan(&topic_img)
		deleteImg = topic_img.String
	}
	os.Remove("common/uploaded_images/" + deleteImg)
}
