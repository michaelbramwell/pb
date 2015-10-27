package paulabramwell

import (
	"appengine"
	"appengine/datastore"
	"encoding/json"
	"fmt"
	htemplate "html/template"
	"io"
	"net/http"
	ttemplate "text/template"
	"time"
)

type ReponseMsg struct {
	Name string
	Body string
	Key  string
}

type ContentItem struct {
	PageNameValue      string
	PageTitleValue     string
	PageMetaDescrValue string
	PageHeaderValue    string
	PageBodyValue      []byte
	PageFooterValue    string
	PathName           string
	DateModified       time.Time
}

type ContentItemWithKey struct {
	ContentItem
	PageBodyStr string
	Key         string
}

var htmlTempl = ttemplate.Must(ttemplate.ParseFiles("tmpl/root.html"))
var htmlTemplEdit = htemplate.Must(htemplate.ParseFiles("tmpl/edit.html"))

func init() {
	http.HandleFunc("/", root)
	http.HandleFunc("/edit", edit)
	http.HandleFunc("/edit/post", post)
	http.HandleFunc("/add", add)
}

func root(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	q := datastore.NewQuery("ContentItem").
		Filter("PathName=", "/edit")

	var contentItemWithKey ContentItemWithKey
	var contentItem ContentItem

	for t := q.Run(c); ; {
		key, err := t.Next(&contentItem)

		if err == datastore.Done {
			//serveError(c, w, err)
			break
		}

		if err != nil {
			serveError(c, w, err)
			return
		}

		contentItemWithKey.ContentItem = contentItem
		contentItemWithKey.Key = key.Encode()
		contentItemWithKey.PageBodyStr = string(contentItem.PageBodyValue)

	}

	w.Header().Set("Content-Type", "text/html")
	htmlTempl.Execute(w, contentItemWithKey)

}

func add(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	htmlTemplEdit.Execute(w, "")
}

func edit(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	q := datastore.NewQuery("ContentItem").
		Filter("PathName=", "/edit")

	var contentItemWithKey ContentItemWithKey
	var contentItem ContentItem

	for t := q.Run(c); ; {
		key, err := t.Next(&contentItem)

		if err == datastore.Done {
			break
		}

		if err != nil {
			serveError(c, w, err)
			return
		}

		contentItemWithKey.ContentItem = contentItem
		contentItemWithKey.Key = key.Encode()
		contentItemWithKey.PageBodyStr = string(contentItem.PageBodyValue)

	}

	w.Header().Set("Content-Type", "text/html")
	htmlTemplEdit.Execute(w, contentItemWithKey)

}

func post(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	ci := ContentItem{
		PageNameValue:      r.FormValue("pageNameValue"),
		PageTitleValue:     r.FormValue("pageTitleValue"),
		PageMetaDescrValue: r.FormValue("pageMetaDescrValue"),
		PageHeaderValue:    r.FormValue("pageHeaderValue"),
		PageBodyValue:      []byte(r.FormValue("pageBodyValue")),
		PageFooterValue:    r.FormValue("pageFooterValue"),
		PathName:           r.FormValue("pathname"),
		DateModified:       time.Now(),
	}

	if r.FormValue("key") == "undefined" || r.FormValue("key") == "" {

		key, err := datastore.Put(c, datastore.NewIncompleteKey(c, "ContentItem", nil), &ci)
		keyStr := key.Encode()

		if err != nil {
			serveError(c, w, err)
			return
		}

		rm := ReponseMsg{"Success", "true", keyStr}
		jrm, _ := json.Marshal(rm)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, string(jrm))

		return
	}

	key, err := datastore.DecodeKey(r.FormValue("key"))

	if err != nil {
		serveError(c, w, err)
		return
	}

	if _, err := datastore.Put(c, key, &ci); err != nil {
		serveError(c, w, err)
		return
	}

	keyStr := r.FormValue("key")
	rm := ReponseMsg{"Success", "true", keyStr}
	jrm, _ := json.Marshal(rm)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(jrm))

	return
}

func serveError(c appengine.Context, w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "text/plain")
	io.WriteString(w, "Internal Server Error - "+err.Error()+"\n")
	c.Errorf("%v", err)
}
