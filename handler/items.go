package handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	db "wallet/database"
	model "wallet/models"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

var itemIDKey = "itemID"

func items(router chi.Router) {
	router.Get("/", fetchData)
	router.Post("/", createData)
	router.Route("/{itemId}", func(router chi.Router) {
		router.Use(ItemContext)
		router.Get("/", fetchDataByID)
		router.Put("/", updateData)
		router.Delete("/", deleteData)
	})
}

func ItemContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		itemId := chi.URLParam(r, "itemId")
		if itemId == "" {
			render.Render(w, r, ErrorRenderer(fmt.Errorf("item ID is required")))
			return
		}
		id, err := strconv.Atoi(itemId)
		if err != nil {
			render.Render(w, r, ErrorRenderer(fmt.Errorf("invalid item ID")))
		}
		ctx := context.WithValue(r.Context(), itemIDKey, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func createData(w http.ResponseWriter, r *http.Request) {
	item := &model.AccountData{}
	if err := render.Bind(r, item); err != nil {
		render.Render(w, r, ErrBadRequest)
		return
	}
	if err := dbInstance.CreateData(item); err != nil {
		render.Render(w, r, ErrorRenderer(err))
		return
	}
	if err := render.Render(w, r, item); err != nil {
		render.Render(w, r, ServerErrorRenderer(err))
		return
	}
}

func fetchData(w http.ResponseWriter, r *http.Request) {
	items, err := dbInstance.FetchData()
	if err != nil {
		render.Render(w, r, ServerErrorRenderer(err))
		return
	}
	if err := render.Render(w, r, items); err != nil {
		render.Render(w, r, ErrorRenderer(err))
	}
}

func fetchDataByID(w http.ResponseWriter, r *http.Request) {
	itemID := r.Context().Value(itemIDKey).(int)
	item, err := dbInstance.FetchDataByID(itemID)
	if err != nil {
		if err == db.NoMatch {
			render.Render(w, r, ErrNotFound)
		} else {
			render.Render(w, r, ErrorRenderer(err))
		}
		return
	}
	if err := render.Render(w, r, &item); err != nil {
		render.Render(w, r, ServerErrorRenderer(err))
		return
	}
}

func deleteData(w http.ResponseWriter, r *http.Request) {
	itemId := r.Context().Value(itemIDKey).(int)
	err := dbInstance.DeleteData(itemId)
	if err != nil {
		if err == db.NoMatch {
			render.Render(w, r, ErrNotFound)
		} else {
			render.Render(w, r, ServerErrorRenderer(err))
		}
		return
	}
}

func updateData(w http.ResponseWriter, r *http.Request) {
	itemId := r.Context().Value(itemIDKey).(int)
	itemData := model.AccountData{}
	if err := render.Bind(r, &itemData); err != nil {
		render.Render(w, r, ErrBadRequest)
		return
	}
	item, err := dbInstance.UpdateData(itemId, itemData)
	if err != nil {
		if err == db.NoMatch {
			render.Render(w, r, ErrNotFound)
		} else {
			render.Render(w, r, ServerErrorRenderer(err))
		}
		return
	}
	if err := render.Render(w, r, &item); err != nil {
		render.Render(w, r, ServerErrorRenderer(err))
		return
	}
}
