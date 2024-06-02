package api

import (
	"net/http"

	"nhooyr.io/websocket"

	"github.com/mykysha/gogames/snake/pkg/log"
)

type API struct {
	logger     log.Logger
	mux        *http.ServeMux
	template   *Template
	page       *IndexPage
	screenChan chan string
}

func NewAPI(logger log.Logger, screenChan chan string) *API {
	api := &API{
		logger:     logger,
		mux:        http.NewServeMux(),
		template:   newTemplate(),
		page:       newIndexPage(),
		screenChan: screenChan,
	}

	api.registerEndpoints()

	return api
}

func (a *API) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	a.mux.ServeHTTP(w, req)
}

func (a *API) registerEndpoints() {
	a.mux.HandleFunc("/css/{file}", a.cssFileHandler)
	a.mux.HandleFunc("/", a.indexPageHandler)
	a.mux.HandleFunc("/game", a.wsHandler)
}

func (a *API) cssFileHandler(w http.ResponseWriter, req *http.Request) {
	a.logger.Info("GET /css/" + req.PathValue("file"))
	http.ServeFile(w, req, "css/"+req.PathValue("file"))
}

func (a *API) indexPageHandler(w http.ResponseWriter, _ *http.Request) {
	a.logger.Info("/")

	if err := a.template.Render(w, "index", a.page); err != nil {
		a.logger.Error("failed to render index page", "error", err)
	}
}

func (a *API) wsHandler(w http.ResponseWriter, req *http.Request) {
	a.logger.Info("/game")

	wsConn, err := websocket.Accept(w, req, nil)
	if err != nil {
		a.logger.Error("failed to accept websocket conn", "error", err)

		return
	}

	defer wsConn.Close(websocket.StatusNormalClosure, "idk unexpected or something")

	for screen := range a.screenChan {
		if err := a.page.UpdateScreen(screen); err != nil {
			a.logger.Error("failed to update screen", "error", err)

			continue
		}

		writer, err := wsConn.Writer(req.Context(), websocket.MessageText)
		if err != nil {
			a.logger.Error("failed to write to websocket conn", "error", err)

			return
		}

		a.logger.Info("Sending screen to client")

		if err := a.template.Render(writer, "screen", a.page); err != nil {
			a.logger.Error("failed to render screen", "error", err)
		}

		writer.Close()
	}
}
