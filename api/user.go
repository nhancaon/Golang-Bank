package api

import (
	"log"
	"net/http"
	db "simple-bank/db/sqlc"
	"simple-bank/util"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
)

type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type createUserResponse struct {
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashedPassword,
		FullName:       req.FullName,
		Email:          req.Email,
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			log.Println(pgErr.Code)
			switch pgErr.Code {
			case "23505": // unique_violation
				ctx.JSON(http.StatusForbidden, errorResponse(pgErr))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := createUserResponse{
		Username: user.Username,
		FullName: user.FullName,
		Email:    user.Email,
		PasswordChangedAt: func() time.Time {
			if user.PasswordChangedAt.Valid {
				return user.PasswordChangedAt.Time
			}
			return time.Time{} // hoặc time.Now() nếu muốn
		}(),
		CreatedAt: func() time.Time {
			if user.CreatedAt.Valid {
				return user.CreatedAt.Time
			}
			return time.Time{}
		}(),
	}

	ctx.JSON(http.StatusOK, rsp)
}
