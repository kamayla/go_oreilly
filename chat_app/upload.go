package main

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func uploaderHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.FormValue("userid")

	multiWriter := io.MultiWriter(w, os.Stdout)
	file, header, err := r.FormFile("avatarFile")

	if err != nil {
		io.WriteString(multiWriter, err.Error())
		return
	}

	defer file.Close()

	data, err := io.ReadAll(file)

	if err != nil {
		io.WriteString(multiWriter, err.Error())
		return
	}

	filename := filepath.Join("avatars", userID+filepath.Ext(header.Filename))
	err = os.WriteFile(filename, data, 0777)

	if err != nil {
		io.WriteString(multiWriter, err.Error())
		return
	}
	io.WriteString(multiWriter, "成功")
}
