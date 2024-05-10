package transport

// func (t *Transport) Oauth2ConnectForm(w http.ResponseWriter, r *http.Request) {
// 	t.render(w, http.StatusOK, "auth.gohtml", t.newTemplateData(r))
// }
//
// func (t *Transport) Oauth2Connect(w http.ResponseWriter, r *http.Request) {
// 	url, err := t.App.Oauth2ConnectProvider(chi.URLParam(r, "provider"))
// 	if err != nil {
// 		t.serverError(w, err)
// 		return
// 	}
// 	t.App.Log.Debug("url to redirect: " + url)
//
// 	w.Header().Set("location", url)
// 	w.WriteHeader(http.StatusTemporaryRedirect)
// }
//
// func (t *Transport) Oauth2Success(w http.ResponseWriter, r *http.Request) {
// 	if err := r.ParseForm(); err != nil {
// 		t.serverError(w, err)
// 	}
//
// 	code := r.Form.Get("code")
// 	prov := chi.URLParam(r, "provider")
// 	if code == "" || prov == "" {
// 		t.clientError(w, http.StatusBadRequest)
// 		return
// 	}
//
// 	err := t.App.Oauth2Success(prov, code)
// 	if err != nil {
// 		t.serverError(w, err)
// 		return
// 	}
//
// 	t.App.Log.Info("successfully authenticated with: " + prov)
// 	w.WriteHeader(http.StatusOK)
// }
//
// func (t *Transport) AddToken(rw http.ResponseWriter, r *http.Request) {
// 	msg, err := t.App.AuthAddAPIToken(r.Form.Get("token"), chi.URLParam(r, "provider"))
// 	if err != nil {
// 		t.serverError(rw, err)
// 		return
// 	}
// 	t.App.Log.Info("valid token for: " + msg)
// 	rw.WriteHeader(http.StatusCreated)
// }
