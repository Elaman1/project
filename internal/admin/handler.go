package admin

import (
	"fmt"
	"myproject/config"
	"net/http"
)

func PanelHandler(w http.ResponseWriter, r *http.Request, ctxApp config.CtxApp) {
	const op = "panelHandler"

	fmt.Println("admin panel")
}
