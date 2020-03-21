package controllers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/NiranjanShetty8/bookmarkapp/model"
	"github.com/NiranjanShetty8/bookmarkapp/security"
	"github.com/NiranjanShetty8/bookmarkapp/services"
	"github.com/NiranjanShetty8/bookmarkapp/web"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

type CategoryController struct {
	CategoryService *services.CategoryService
}

func (cc *CategoryController) RegisterRoutes(router *mux.Router) {
	subRouter := router.PathPrefix("/api/bookmarkapp/user/{userid}").Subrouter()
	subRouter.Use(security.AuthMiddleWareFunc)
	subRouter.HandleFunc("/category", cc.GetAllCategories).Methods("GET")
	subRouter.HandleFunc("/category/{categoryid}", cc.GetCategoryByID).Methods("GET")
	subRouter.HandleFunc("/category/name/{categoryname}", cc.GetCategoryByName).Methods("GET")
	subRouter.HandleFunc("/category", cc.AddCategory).Methods("POST")
	subRouter.HandleFunc("/category/{categoryid}", cc.UpdateCategory).Methods("PUT")
	subRouter.HandleFunc("/category/{categoryid}", cc.DeleteCategory).Methods("DELETE")
}

func (cc *CategoryController) GetAllCategories(w http.ResponseWriter, r *http.Request) {
	uid, err := ParseID(&w, r, "userid")
	if err != nil {
		return
	}
	categories := []model.Category{}
	err = cc.CategoryService.GetAllCategories(uid, &categories)
	if err != nil {
		web.RespondError(&w, err)
		return
	}
	web.RespondJSON(&w, http.StatusOK, categories)
}

func (cc *CategoryController) GetCategoryByID(w http.ResponseWriter, r *http.Request) {
	userID, err := ParseID(&w, r, "userid")
	if err != nil {
		return
	}
	category := model.Category{}
	categoryID, err := ParseID(&w, r, "categoryid")
	if err != nil {
		return
	}
	err = cc.CategoryService.GetCategory(userID, categoryID, &category)
	if err != nil {
		web.RespondError(&w, err)
	}
}

func (cc *CategoryController) GetCategoryByName(w http.ResponseWriter, r *http.Request) {
	_, err := ParseID(&w, r, "userid")
	if err != nil {
		return
	}
	category := model.Category{}
	err = cc.CategoryService.GetCategoryByName(mux.Vars(r)["categoryname"], &category)
	if err != nil {
		web.RespondError(&w, err)
	}
}

func (cc *CategoryController) AddCategory(w http.ResponseWriter, r *http.Request) {
	userID, err := ParseID(&w, r, "userid")
	if err != nil {
		return
	}
	category := model.Category{}
	err = web.UnmarshalJSON(r, category)
	if err != nil {
		web.RespondError(&w, web.NewValidationError("error", map[string]string{"error": err.Error()}))
		return
	}
	err = validateCategory(&category)
	if err != nil {
		web.RespondError(&w, web.NewValidationError("error", map[string]string{"error": err.Error()}))
		return
	}
	category.UserID = userID
	err = cc.CategoryService.AddCategory(&category)
	if err != nil {
		web.RespondError(&w, err)
	}
}

func (cc *CategoryController) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	userID, err := ParseID(&w, r, "userid")
	if err != nil {
		return
	}
	category := model.Category{}
	err = web.UnmarshalJSON(r, category)
	if err != nil {
		web.RespondError(&w, web.NewValidationError("error", map[string]string{"error": err.Error()}))
		return
	}
	err = validateCategory(&category)
	if err != nil {
		web.RespondError(&w, web.NewValidationError("error", map[string]string{"error": err.Error()}))
		return
	}
	catergoryID, err := ParseID(&w, r, "categoryid")
	if err != nil {
		return
	}
	category.ID = catergoryID
	category.UserID = userID
	err = cc.CategoryService.UpdateCategory(&category)
	if err != nil {
		web.RespondError(&w, err)
	}
}

func (cc *CategoryController) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	userID, err := ParseID(&w, r, "userid")
	if err != nil {
		return
	}
	catergoryID, err := ParseID(&w, r, "categoryid")
	if err != nil {
		return
	}
	err = cc.CategoryService.DeleteCategory(userID, catergoryID)
	if err != nil {
		web.RespondError(&w, err)
	}
}

func NewCategoryController(cs *services.CategoryService) *CategoryController {
	return &CategoryController{
		CategoryService: cs,
	}
}

func ParseID(w *http.ResponseWriter, r *http.Request, mapKey string) (uuid.UUID, error) {
	id := mux.Vars(r)[mapKey]
	uid, err := uuid.FromString(id)
	if err != nil {
		web.RespondError(w, web.NewValidationError(mapKey,
			map[string]string{"Error": "Invalid " + mapKey}))
		return uuid.Nil, err
	}
	return uid, err
}

func validateCategory(category *model.Category) error {
	name := strings.TrimSpace(category.Name)
	if name == "" {
		return errors.New("Name is required")
	}
	return nil
}
