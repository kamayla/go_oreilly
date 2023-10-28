package main

import (
	"flag"
	"fmt"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(
			template.ParseFiles(
				filepath.Join("templates", t.filename),
			),
		)
	})

	data := map[string]any{
		"Host": r.Host,
	}

	if authCoolie, err := r.Cookie("auth"); err == nil {
		fmt.Println(authCoolie.Value)
		data["UserData"] = objx.MustFromBase64(authCoolie.Value)
	}

	t.templ.Execute(w, data)
}

func main() {
	addr := flag.String("addr", ":8080", "アプリケーションのアドレス")
	flag.Parse() // フラグを解析
	gomniauth.SetSecurityKey("hoge")
	gomniauth.WithProviders(
		google.New(
			"462491467947-ebliib85g10b0k220opeucvd6fnunpk0.apps.googleusercontent.com",
			"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",
			"http://localhost:8080/auth/callback/google",
		),
	)
	r := newRoom()

	http.Handle("/chat", MustAuth(
		&templateHandler{
			filename: "chat.html",
		},
	))

	http.Handle("/room", r)

	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)

	go r.run()

	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal(err)
	}
}
