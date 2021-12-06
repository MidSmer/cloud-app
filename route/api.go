package route

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strings"
)

func API(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	apiName := vars["name"]

	switch apiName {
	case "upload":
		uploadAPI(w, r)
	}
}

func multipartUpload(destURL string, f io.Reader, fields map[string]string) (*http.Response, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fw, err := writer.CreateFormFile(fields["file"], fields["filename"])
	if err != nil {
		return nil, fmt.Errorf("CreateFormFile %v", err)
	}

	_, err = io.Copy(fw, f)
	if err != nil {
		return nil, fmt.Errorf("copying fileWriter %v", err)
	}

	for k, v := range fields {
		if k == "file" {
			continue
		}
		_ = writer.WriteField(k, v)
	}

	err = writer.Close() // close writer before POST request
	if err != nil {
		return nil, fmt.Errorf("writerClose: %v", err)
	}

	resp, err := http.Post(destURL, writer.FormDataContentType(), body)
	if err != nil {
		return nil, err
	}

	return resp, nil

	// req, err := http.NewRequest("POST", destURL, body)
	// if err != nil {
	//  return nil, err
	// }

	// req.Header.Set("Content-Type", writer.FormDataContentType())

	// if req.Close && req.Body != nil {
	//  defer req.Body.Close()
	// }

	// return http.DefaultClient.Do(req)
}

func uploadAPI(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(2 << 10); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Something bad happened!"))
		return
	}

	target := r.PostFormValue("target")
	targetUri := ""
	fields := make(map[string]string)
	var resultParse func(*bytes.Buffer, []byte) error

	switch target {
	case "smSite":
		targetUri = "https://sm.ms/api/v2/upload"
		fields = map[string]string{
			"file":   "smfile",
			"format": "json",
		}
		resultParse = func(w *bytes.Buffer, r []byte) error {
			type smDataRes struct {
				Filename string `json:"filename"`
				Height   int    `json:"height"`
				Width    int    `json:"width"`
				Size     int    `json:"size"`
				Url      string `json:"url"`
			}
			type smRes struct {
				Code    string    `json:"code"`
				Message string    `json:"message"`
				Data    smDataRes `json:"data"`
				Success bool      `json:"success"`
			}

			var res smRes
			if err := json.Unmarshal(r, &res); err != nil {
				return err
			}

			if err := json.NewEncoder(w).Encode(res); err != nil {
				return err
			}

			return nil
		}
	}
	if targetUri == "" {
		Logger.Error("web upload api handler fail! not target uri")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		Logger.Error("web upload api handler fail! not file ", err)
		return
	}
	defer file.Close()

	nameParts := strings.Split(header.Filename, ".")
	filename := nameParts[1]
	fields["filename"] = filename

	res, err := multipartUpload(targetUri, file, fields)
	if err != nil {
		Logger.Error("web upload api handler fail! target upload fail ", err)
		return
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		Logger.Error("web upload api handler fail! response read fail ", err)
		return
	}

	resAPI := &bytes.Buffer{}
	if err := resultParse(resAPI, body); err != nil {
		Logger.Error("web upload api handler fail! result parse fail ", err)
		return
	}

	w.Write(resAPI.Bytes())
}
