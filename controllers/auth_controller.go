// Package controllers Auth controllers
package controllers

import (
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"

	"github.com/feserr/pheme-backend/models"
)

// Register godoc
// @Summary      Register a user
// @Description  add a user
// @Tags         authentication
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.User
// @Failure      400  {object}  models.Message
// @Router       /register [post]
func Register(c *fiber.Ctx) error {
	body := models.UserParamsNew{}

	if err := c.BodyParser(&body); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Invalid JSON body",
		})
	}

	if err := validate.Struct(body); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Wrong JSON params",
		})
	}

	user, err := models.RegisterUser(body.UserName, body.Email, body.Avatar, body.Password)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Email already used",
		})
	}

	return c.JSON(user)
}

// Login godoc
// @Summary      Login to a user
// @Description  user login
// @Tags         authentication
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.Message
// @Failure      404  {object}  models.Message
// @Failure      400  {object}  models.Message
// @Failure      500  {object}  models.Message
// @Router       /login [post]
func Login(c *fiber.Ctx) error {
	body := models.UserParamsLogin{}

	if err := c.BodyParser(&body); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Invalid JSON body",
		})
	}

	if err := validate.Struct(body); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Wrong JSON params",
		})
	}

	user, err := models.GetUserByEmail(body.Email)
	if err != nil {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "User not found",
		})
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(body.Password)); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Incorrect password",
		})
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(user.ID)),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	})

	token, err := claims.SignedString([]byte(SecretKey))

	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Could not login",
		})
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(user)
}

// User godoc
// @Summary      Retrieve the user
// @Description  get the user cookie
// @Tags         authentication
// @Produce      json
// @Success      200  {object}  models.User
// @Failure      401  {object}  models.Message
// @Router       /user [get]
func User(c *fiber.Ctx) error {
	user, err := models.GetUser(c, SecretKey)
	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "Unauthenticated",
		})
	}

	return c.JSON(user)
}

// Logout godoc
// @Summary      Logout the user
// @Description  logout and remove the user cookie
// @Tags         authentication
// @Produce      json
// @Success      200  {object}  models.Message
// @Router       /logout [post]
func Logout(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "Success",
	})
}

// Delete godoc
// @Summary      Delete the user
// @Description  delete the user
// @Tags         authentication
// @Produce      json
// @Success      200  {object}  models.Message
// @Failure      400  {object}  models.Message
// @Failure      401  {object}  models.Message
// @Router       /user [delete]
func Delete(c *fiber.Ctx) error {
	user, err := models.GetUser(c, SecretKey)
	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "Unauthenticated",
		})
	}

	err = models.DeleteByID(user.ID)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "User don't exist",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Success",
	})
}
