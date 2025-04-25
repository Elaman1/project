package admin

import (
	"encoding/json"
	"myproject/config"
	"net/http"
)

func PanelHandler(w http.ResponseWriter, r *http.Request, ctxApp config.CtxApp) {
	const op = "panelHandler"

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "admin panel"})
}
