package route

import (
	"encoding/json"
	"github.com/MidSmer/cloud-app/db"
	"io"
	"io/ioutil"
	"net/http"
)

type createArticleRequest struct {
	Content string `json:"content"`
}

type createArticleDataResult struct {
	Key string `json:"key"`
}

type createArticleResult struct {
	Code    string                  `json:"code"`
	Message string                  `json:"message"`
	Data    createArticleDataResult `json:"data"`
	Success bool                    `json:"success"`
}

func CreateArticle(w http.ResponseWriter, r *http.Request) {
	var (
		req createArticleRequest
		res createArticleResult
	)

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		Logger.Panic(err)

		res = createArticleResult{
			Code:    "fail",
			Message: "body read fail",
			Success: false,
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(res); err != nil {
			Logger.Panic(err)
		}
		return
	}
	if err := r.Body.Close(); err != nil {
		Logger.Panic(err)

		res = createArticleResult{
			Code:    "fail",
			Message: "body read fail",
			Success: false,
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(res); err != nil {
			Logger.Panic(err)
		}
		return
	}
	if err := json.Unmarshal(body, &req); err != nil {
		Logger.Panic(err)

		res = createArticleResult{
			Code:    "fail",
			Message: "body parsing fail",
			Success: false,
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(res); err != nil {
			Logger.Panic(err)
		}
		return
	}

	key, err := db.CreateArticle(req.Content)
	if err != nil {
		Logger.Error("create article fail!", err)

		res = createArticleResult{
			Code:    "fail",
			Message: "create fail",
			Success: false,
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(res); err != nil {
			Logger.Panic(err)
		}
		return
	}

	res = createArticleResult{
		Code:    "success",
		Message: "create success",
		Data: createArticleDataResult{
			Key: key,
		},
		Success: true,
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		Logger.Panic(err)
	}
}

type updateArticleRequest struct {
	Key     string `json:"key"`
	Content string `json:"content"`
}

type updateArticleDataResult struct {
	Key string `json:"key"`
}

type updateArticleResult struct {
	Code    string                  `json:"code"`
	Message string                  `json:"message"`
	Data    updateArticleDataResult `json:"data"`
	Success bool                    `json:"success"`
}

func UpdateArticle(w http.ResponseWriter, r *http.Request) {
	var (
		req updateArticleRequest
		res updateArticleResult
	)

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		Logger.Panic(err)

		res = updateArticleResult{
			Code:    "fail",
			Message: "body read fail",
			Success: false,
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(res); err != nil {
			Logger.Panic(err)
		}
		return
	}
	if err := r.Body.Close(); err != nil {
		Logger.Panic(err)

		res = updateArticleResult{
			Code:    "fail",
			Message: "body read fail",
			Success: false,
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(res); err != nil {
			Logger.Panic(err)
		}
		return
	}
	if err := json.Unmarshal(body, &req); err != nil {
		Logger.Panic(err)

		res = updateArticleResult{
			Code:    "fail",
			Message: "body parsing fail",
			Success: false,
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(res); err != nil {
			Logger.Panic(err)
		}
		return
	}

	key, err := db.UpdateArticle(req.Key, req.Content)
	if err != nil {
		Logger.Error("update article fail!", err)

		res = updateArticleResult{
			Code:    "fail",
			Message: "update fail",
			Success: false,
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(res); err != nil {
			Logger.Panic(err)
		}
		return
	}

	res = updateArticleResult{
		Code:    "success",
		Message: "update success",
		Data: updateArticleDataResult{
			Key: key,
		},
		Success: true,
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		Logger.Panic(err)
	}
}

type fetchArticleRequest struct {
	Key string `json:"key"`
}

type fetchArticleDataResult struct {
	Content string `json:"content"`
}

type fetchArticleResult struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Data    fetchArticleDataResult `json:"data"`
	Success bool                   `json:"success"`
}

func FetchArticle(w http.ResponseWriter, r *http.Request) {
	var (
		req fetchArticleRequest
		res fetchArticleResult
	)

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		Logger.Panic(err)

		res = fetchArticleResult{
			Code:    "fail",
			Message: "body read fail",
			Success: false,
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(res); err != nil {
			Logger.Panic(err)
		}
		return
	}
	if err := r.Body.Close(); err != nil {
		Logger.Panic(err)

		res = fetchArticleResult{
			Code:    "fail",
			Message: "body read fail",
			Success: false,
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(res); err != nil {
			Logger.Panic(err)
		}
		return
	}
	if err := json.Unmarshal(body, &req); err != nil {
		Logger.Panic(err)

		res = fetchArticleResult{
			Code:    "fail",
			Message: "body parsing fail",
			Success: false,
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(res); err != nil {
			Logger.Panic(err)
		}
		return
	}

	content, err := db.GetArticleForKey(req.Key)
	if err != nil {
		Logger.Error("fetch article fail!", err)

		res = fetchArticleResult{
			Code:    "fail",
			Message: "fetch fail",
			Success: false,
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(res); err != nil {
			Logger.Panic(err)
		}
		return
	}

	res = fetchArticleResult{
		Code:    "success",
		Message: "fetch success",
		Data: fetchArticleDataResult{
			Content: content,
		},
		Success: true,
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		Logger.Panic(err)
	}
}
