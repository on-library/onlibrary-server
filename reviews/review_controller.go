package reviews

import (
	"net/http"
	bookModel "onlibrary/books/models"
	"onlibrary/common"
	"onlibrary/database"
	reviewModel "onlibrary/reviews/models"

	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
)

type(
	ReviewController struct {

	}

	AddReviewRequest struct {
		BookId		uuid.UUID		`json:"bookId"`
		Comment		string		`json:"comment"`
		Rating		uint		`json:"rating"`
	}

	DeleteReviewRequest struct {
		ID		uuid.UUID	`json:"id"`
	}
)

func (controller ReviewController) Routes() []common.Route {
	return []common.Route {
		{
			Method: echo.POST,
			Path: "/review/add",
			Handler: controller.AddReview,
		},
		{
			Method: echo.DELETE,
			Path: "/review/delete",
			Handler: controller.DeleteReview,
		},
	}
}

func (controller ReviewController) AddReview(c echo.Context) error {
	db := database.GetInstance()
	params := new(AddReviewRequest)

	if err:=c.Bind(params); err!= nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	var book bookModel.Book

	newId := uuid.NewV1()

	var review = reviewModel.Review{ID: newId,Comment: params.Comment,Rating: params.Rating,BookRefer: params.BookId}

	if err:= db.First(&book, "id = ?", params.BookId); err != nil {
		var r = struct {
			common.GeneralResponseJSON
		}{
			GeneralResponseJSON: common.GeneralResponseJSON{Message: "Book not found"},
		}
	
		return c.JSON(http.StatusBadRequest, r)
	}

	
	db.Model(&book).Association("Reviews").Append(&review)

	
	var r = struct {
		common.GeneralResponseJSON
		Data  reviewModel.Review `json:"data"`
	}{
		GeneralResponseJSON: common.GeneralResponseJSON{Message: "success"},
		Data: review,
	}


	return c.JSON(http.StatusOK, r)
}

func (controller ReviewController) DeleteReview(c echo.Context) error {
	db:=database.GetInstance()
	params := new(DeleteReviewRequest)

	if err:=c.Bind(params);err!= nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	var review reviewModel.Review

	if err := db.Where("id = ?", params.ID).Delete(&review); err!=nil{
		var r = struct {
			common.GeneralResponseJSON
		}{
			GeneralResponseJSON: common.GeneralResponseJSON{Message: "Review not found"},
		}
		return c.JSON(http.StatusBadRequest,r)
	}

	var r = struct {
		common.GeneralResponseJSON
		Data uuid.UUID `json:"data"`
	}{
		GeneralResponseJSON: common.GeneralResponseJSON{Message: "success"},
		Data: params.ID,
	}

	return c.JSON(http.StatusOK, r)
}