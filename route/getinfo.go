package route

import (
	"github.com/MidSmer/cloud-app/external/v2ray.com/core"
	"net/http"
	"strconv"
)

func GetInfo(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
	id := v.Get("id")
	accountManager := core.GetAccountManagerInstance()

	err := accountManager.Add(id)

	w.Header().Set("Content-Type", "text/html")

	w.Write([]byte(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Info</title>
</head>
<body>
`))

	if err == nil {
		w.Write([]byte(`
<div>success</div>
`))
	} else {
		w.Write([]byte(`
<div>user:
`))
		w.Write([]byte(strconv.Itoa(len(accountManager.Get()))))
		w.Write([]byte(`
</div>
`))
	}

	w.Write([]byte(`
</body>
</html>
`))
}
