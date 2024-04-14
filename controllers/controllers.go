package controllers

import (
	"avito-app/database"
	"avito-app/middlewares"
	"avito-app/models"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func ParseToken(tokenString string) (claims *models.Claims, err error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(middlewares.JWT), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*models.Claims)
	if !ok {
		return nil, err
	}
	return claims, nil
}

func GetUserBanner(c *fiber.Ctx) error {
	tagID := c.Query("tag_id")
	featureID := c.Query("feature_id")
	if tagID == "" || featureID == "" {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"error": "Bad Request",
		})
	}
	token := c.Get("token")
	claims, _ := ParseToken(token)
	var banner models.Banner
	if claims.Role == "admin" {
		database.DB.Raw("SELECT * FROM banners WHERE array_position(tag_ids , @value) IS NOT NULL AND feature_id = @value2 ORDER BY updated_at DESC",
			sql.Named("value", tagID), sql.Named("value2", featureID)).Find(&banner)
	} else {
		database.DB.Raw("SELECT * FROM banners WHERE array_position(tag_ids , @value) IS NOT NULL AND feature_id = @value2 AND is_active = true ORDER BY updated_at DESC",
			sql.Named("value", tagID), sql.Named("value2", featureID)).Find(&banner)
	}
	if banner.BannerID == 0 {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"error": "Not Found",
		})
	}
	return c.JSON(banner.Content)
}

func GetBanner(c *fiber.Ctx) error {
	var banners []models.Banner
	query := database.DB.Where("")
	tagID := c.Query("tag_id")
	featureID := c.Query("feature_id")
	limit := c.Query("limit")
	offset := c.Query("offset")
	if featureID != "" {
		ID, _ := strconv.Atoi(featureID)
		query.Where("feature_id = ?", ID)
	}
	if tagID != "" {
		ID, _ := strconv.Atoi(tagID)
		query.Where(fmt.Sprintf("array_position(tag_ids , %d) IS NOT NULL", ID))
	}
	if limit != "" {
		lmt, _ := strconv.Atoi(limit)
		query.Limit(lmt)
		if offset != "" {
			offst, _ := strconv.Atoi(offset)
			query.Offset(offst)
		}
	}
	query.Find(&banners)
	c.Status(fiber.StatusOK)
	return c.JSON(banners)
}

func PostBanner(c *fiber.Ctx) error {
	data := new(models.Banner)
	if err := c.BodyParser(&data); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"error": "Bad Request",
		})
	}
	database.DB.Create(data)
	c.Status(fiber.StatusCreated)
	return c.JSON(fiber.Map{
		"banner_id": data.BannerID,
	})
}

func PatchBanner(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"error": "Bad Request",
		})
	}
	data := new(models.Banner)
	if err := c.BodyParser(&data); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"error": "Bad Request",
		})
	}
	var banner models.Banner
	result := database.DB.Where("banner_id = ?", id).First(&banner)
	if result.RowsAffected == 0 {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"error": "Not Found",
		})
	}
	database.DB.Model(banner).Updates(data)
	c.Status(fiber.StatusOK)
	return c.SendString("")
}

func DeleteBanner(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"error": "Bad Request",
		})
	}
	var banner = models.Banner{BannerID: uint(id)}
	result := database.DB.First(&banner)
	if result.RowsAffected == 0 {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"error": "Not Found",
		})
	}
	database.DB.Delete(banner)
	c.Status(fiber.StatusNoContent)
	return c.SendString("")
}

func PostRegister(c *fiber.Ctx) error {
	var data map[string]string
	if err := c.BodyParser(&data); err != nil || data["password"] == "" || data["login"] == "" {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"error": "Bad Request",
		})
	}
	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)
	user := models.User{
		Login:    data["login"],
		Role:     "user",
		Password: password,
	}
	result := database.DB.Create(&user)
	if result.RowsAffected == 0 {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"error": result.Error.Error(),
		})
	}
	c.Status(fiber.StatusCreated)
	return c.JSON(fiber.Map{
		"user_id": user.ID,
	})
}

func PostLogin(c *fiber.Ctx) error {
	var data map[string]string
	if err := c.BodyParser(&data); err != nil || data["password"] == "" || data["login"] == "" {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"error": "Bad Request",
		})
	}
	var user models.User
	result := database.DB.Where("login = ?", data["login"]).First(&user)
	if result.RowsAffected == 0 {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"error": result.Error.Error(),
		})
	}
	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"])); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"error": "Bad Request",
		})
	}
	claims := &models.Claims{
		Role: user.Role,
		StandardClaims: jwt.StandardClaims{
			Issuer:    strconv.Itoa(int(user.ID)),
			Subject:   user.Login,
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(middlewares.JWT))
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	c.Status(fiber.StatusOK)
	return c.JSON(fiber.Map{
		"token": tokenString,
	})
}
