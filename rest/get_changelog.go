package rest

import (
	"net/http"

	"github.com/arelate/gaugin/gaugin_middleware"
	"github.com/arelate/gaugin/view_models"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
)

func GetChangelog(w http.ResponseWriter, r *http.Request) {

	// GET /changelog?id

	id := r.URL.Query().Get("id")

	idRedux, err := getRedux(
		http.DefaultClient,
		id,
		false,
		vangogh_local_data.ChangelogProperty)

	if err != nil {
		http.Error(w, nod.ErrorStr("error getting redux"), http.StatusInternalServerError)
		return
	}

	gaugin_middleware.DefaultHeaders(nil, w)

	cvm := view_models.NewChangelog(idRedux[id])

	if err := tmpl.ExecuteTemplate(w, "changelog-page", cvm); err != nil {
		http.Error(w, nod.ErrorStr("template exec error"), http.StatusInternalServerError)
		return
	}
}