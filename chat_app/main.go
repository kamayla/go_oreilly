package main

import (
	"flag"
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
)

var avatars Avatar = TryAvatars{
	UseFileSystemAvatar,
	UseAuthAvatar,
	UseGravatar,
}

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
	http.Handle("/upload", &templateHandler{filename: "upload.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.HandleFunc("/uploader", uploaderHandler)
	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:   "auth",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})
		w.Header().Set("Location", "/chat")
		w.WriteHeader(http.StatusTemporaryRedirect)
	})

	http.Handle("/avatars/", http.StripPrefix("/avatars/", http.FileServer(http.Dir("./avatars"))))

	go r.run()

	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal(err)
	}
}
