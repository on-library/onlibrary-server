package reviews

import (
	"fmt"
	"net/http"
	authModel "onlibrary/auth/models"
	bookModel "onlibrary/books/models"
	"onlibrary/common"
	"onlibrary/database"
	reviewModel "onlibrary/reviews/models"

	"github.com/golang-jwt/jwt"
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

	FindReviewRequest struct {
		BookID	uuid.UUID	`json:"book_id"`
	}
)

func (controller ReviewController) Routes() []common.Route {
	return []common.Route {
		{
			Method: echo.POST,
			Path: "/review/add",
			Handler: controller.AddReview,
			Middleware: []echo.MiddlewareFunc{common.JwtMiddleware()},
		},
		{
			Method: echo.DELETE,
			Path: "/review/delete",
			Handler: controller.DeleteReview,
		},
		{
			Method: echo.POST,
			Path: "/review/findbybook",
			Handler: controller.FindReviewByBook,
		},
	}
}

func (controller ReviewController) AddReview(c echo.Context) error {
	db := database.GetInstance()
	params := new(AddReviewRequest)

	if err:=c.Bind(params); err!= nil {
		return c.JSON(http.StatusBadRequest, err)
	}


	user := c.Get("user").(*jwt.Token)
	// fmt.Print(user)
	claims := user.Claims.(*common.JwtCustomClaims)

	var book bookModel.Book
	var userProfile authModel.Auth

	if err := db.Select("id,username","name").First(&userProfile, "id = ?",claims.ID); err.Error != nil {
		return c.JSON(http.StatusBadRequest,echo.Map{
			"message":"User nout found",
			"status":"error",
		})
	}


	newId := uuid.NewV1()

	var review = reviewModel.Review{ID: newId,Comment: params.Comment,Rating: params.Rating,BookRefer: params.BookId}

	if err:= db.First(&book, "book_id = ?", params.BookId); err.Error != nil {
		var r = struct {
			common.GeneralResponseJSON
		}{
			GeneralResponseJSON: common.GeneralResponseJSON{Message: "Book not found"},
		}
		fmt.Println(err.Error)
		return c.JSON(http.StatusBadRequest, r)
	}

	
	db.Model(&userProfile).Association("Reviews").Append(&review)
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

	if err := db.Where("review_id = ?", params.ID).Delete(&review); err!=nil{
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

func (controller ReviewController) FindReviewByBook(c echo.Context) error {
	db:=database.GetInstance()
	params := new(FindReviewRequest)

	type ReviewsWithUser struct {
		reviewModel.Review
		User authModel.Auth		`json:"user"`
	}

	if err:=c.Bind(params);err!= nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	var reviews []reviewModel.Review
	var reviewsWithUser []ReviewsWithUser
	var user authModel.Auth


	if err:= db.Model(&reviews).Where("book_refer = ?",params.BookID).Order("created_at desc").Find(&reviews);err.Error !=nil{
		return c.JSON(http.StatusOK, echo.Map{
			"message":"",
			"status":"error",
		})
	}

	for i:=0; i<len(reviews);i++{

		if err := db.Select("username", "name").First(&user, "id = ?", reviews[i].AuthReviewRefer);err.Error!=nil{
			return c.JSON(http.StatusBadRequest, echo.Map{"message":"user id not found","status":"error"})
		}

		reviewsWithUser = append(reviewsWithUser, ReviewsWithUser{Review: reviews[i],User:user } )
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message":"",
		"status":"success",
		"data":reviewsWithUser,
	})
}