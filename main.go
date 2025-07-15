package main

import (
	"bytes"
	"fmt"
	"net"
	"net/http"
	"os"
	"text/template"
)

const (
	ansi_red    = "\x1b[31m"
	ansi_yellow = "\x1b[33m"
	ansi_blue   = "\x1b[34m"
	ansi_reset  = "\x1b[0m"
)

const tmpl = `<!DOCTYPE html>
<html>
	<head>
	  <meta charset="UTF-8" />
	  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
	  <style>
		body {
		  margin: 0; font-family: sans-serif;
		  display: flex; justify-content: center; align-items: center;
		  min-height: 100vh; background: #f5f5f5;
		}
		.box {
		  background: #fff; padding: 2rem; border-radius: 10px;
		  box-shadow: 0 4px 12px rgba(0,0,0,0.1); text-align: center;
		  width: 90%; max-width: 360px;
		}
		h1 { margin-bottom: 1rem; font-size: 1.4rem; }
		.file-input {
		  display: inline-block; margin: 1rem 0;
		}
		input[type="file"]::file-selector-button {
		  background: #007bff; color: #fff; border: none; padding: 0.4rem 1rem;
		  border-radius: 5px; cursor: pointer; font-size: 0.9rem;
		}
		input[type="file"]::file-selector-button:hover {
		  background: #0056b3;
		}
		input[type="submit"] {
		  background: #28a745; color: white; border: none;
		  padding: 0.5rem 1.2rem; border-radius: 6px; cursor: pointer;
		  font-size: 1rem;
		}
		.upload-status {
		  margin-top: 1.5rem; padding: 0.75rem; background: #e6ffed;
		  color: #207a39; border: 1px solid #b2e5c1; border-radius: 6px;
		  font-size: 0.95rem; word-wrap: break-word;
		}
	  </style>
	</head>
	<body>
	  <div class="box">
		<h1>Upload File</h1>

		<form method="POST" enctype="multipart/form-data">
		  <input class="file-input" type="file" name="file" id="file" required />
		  <br>
		  <input type="submit" value="Submit" />
		</form>

		{{if .FILE_UPLOAD}}
		  <div class="upload-status">
			File "<strong>{{.FILE_UPLOAD}}</strong>" uploaded successfully.
		  </div>
		{{end}}
	  </div>
	</body>
</html>`

type handler struct{ reciver bool }

func red(s string) string {
	return ansi_red + s + ansi_reset + "\n"
}

func blue(s string) string {
	return ansi_blue + s + ansi_reset + "\n"
}

func handlePost(r *http.Request) string {
	f, fh, err := r.FormFile("file")
	if err != nil {
		fmt.Print(red(err.Error()))
		return ""
	}

	buf := make([]byte, fh.Size)
	_, err = f.Read(buf)
	if err != nil {
		fmt.Print(red(err.Error()))
		return ""
	}

	fmt.Printf(blue("[FILE IN] [%s] [%d]"), fh.Filename, fh.Size)
	os.WriteFile(fh.Filename, buf, 0664)

	return fh.Filename
}

func fillTemplate(tmpl string, data map[string]any) []byte {
	t := template.Must(template.New("tmpl").Parse(tmpl))
	var buf bytes.Buffer
	t.Execute(&buf, data)
	return buf.Bytes()
}

func (h handler) serveReciver(w http.ResponseWriter, r *http.Request) {
	fileName := ""
	if r.Method == "POST" {
		fileName = handlePost(r)
	}

	data := map[string]any{
		"FILE_UPLOAD": fileName,
	}

	w.Write(fillTemplate(tmpl, data))
}

func (h handler) serveTransceiver(w http.ResponseWriter, r *http.Request) {
	fs := http.FileServer(http.Dir("."))
	fs.ServeHTTP(w, r)
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Printf(blue("[REQUEST IN] [%s] [%s %s] %s"), r.UserAgent(), r.Method, r.URL.Path, r.RemoteAddr)

	redir := func() {
		w.Header().Set("Location", "/")
		w.WriteHeader(http.StatusFound)
	}

	if r.URL.Path == "/r" {
		h.reciver = true
		redir()
		return
	}

	if r.URL.Path == "/t" {
		h.reciver = false
		redir()
		return
	}

	if h.reciver {
		h.serveReciver(w, r)
		return
	}
	h.serveTransceiver(w, r)
}

func main() {
	port := 8000
	h := handler{reciver: true}

	addr := fmt.Sprintf("0.0.0.0:%d", port)
	l, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Println(red(err.Error()))
		return
	}
	fmt.Printf(blue("Connected to %s"), addr)

	http.Serve(l, &h)
}
