package swagger

import (
	"fmt"
	"net/http"

	"github.com/swaggo/swag"
)

const scalarHTML = `<!doctype html>
<html>
<head>
  <title>ms-catalog API</title>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
</head>
<body>
  <script
    id="api-reference"
    data-url="/swagger/doc.json"
    data-configuration='{"theme":"purple","layout":"modern"}'
  ></script>
  <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
</body>
</html>`

func Spec(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, _ = fmt.Fprint(w, swaggerDoc())
}

func swaggerDoc() string {
	s, err := swag.ReadDoc()
	if err != nil {
		return `{"error":"swagger spec not available"}`
	}
	return s
}

func ScalarUI(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = fmt.Fprint(w, scalarHTML)
}
