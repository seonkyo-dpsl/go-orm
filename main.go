package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type UpcomingFeatureWaitlist struct {
	ID     string `gorm:"primaryKey"`
	Target string `gorm:"type:varchar(255)"`
	Email  string `gorm:"type:varchar(255)"`
}

func (UpcomingFeatureWaitlist) TableName() string {
	return "upcoming_feature_waitlist"
}

func main() {
	app := fiber.New()

	dsn := "host=localhost user=seonkyo password=test123 dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Seoul"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("데이터베이스 연결 실패:", err)
	}

	app.Get("/upcoming-feature-waitlist", func(c *fiber.Ctx) error {
		upcomingFeatureWaitlists := getAllUpcomingFeatureWaitlists(db)
		return c.JSON(upcomingFeatureWaitlists)
	})

	app.Get("/upcoming-feature-waitlist/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		upcomingFeatureWaitlist := getUpcomingFeatureWaitlistByID(db, id)
		return c.JSON(upcomingFeatureWaitlist)
	})

	app.Post("/upcoming-feature-waitlist", func(c *fiber.Ctx) error {
		upcomingFeatureWaitlist := new(UpcomingFeatureWaitlist)
		if err := c.BodyParser(upcomingFeatureWaitlist); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		createUpcomingFeatureWaitlist(db, upcomingFeatureWaitlist.Target, upcomingFeatureWaitlist.Email)
		return c.Status(201).JSON(upcomingFeatureWaitlist)
	})

	app.Put("/upcoming-feature-waitlist/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		upcomingFeatureWaitlist := new(UpcomingFeatureWaitlist)
		if err := c.BodyParser(upcomingFeatureWaitlist); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		updateUpcomingFeatureWaitlist(db, id, upcomingFeatureWaitlist.Target, upcomingFeatureWaitlist.Email)
		return c.Status(200).JSON(upcomingFeatureWaitlist)
	})

	app.Delete("/upcoming-feature-waitlist/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		deleteUpcomingFeatureWaitlist(db, id)
		return c.Status(204).SendString("Deleted")
	})

	// getUpcomingFeatureWaitlistByID(db, "019300407141-ef58-2aac-0ecc-509727c2")
	// createUpcomingFeatureWaitlist(db, "integration", "seonkyo.jeong@deepsales.com")
	// updateUpcomingFeatureWaitlist(db, "c2ca914f-3284-4110-8d09-7a4659cecfdc", "seonkyo.jeong-modified@deepsales.com")
	// deleteUpcomingFeatureWaitlist(db, "c2ca914f-3284-4110-8d09-7a4659cecfdc")

	log.Fatal(app.Listen(":3000"))
}

func getAllUpcomingFeatureWaitlists(db *gorm.DB) []UpcomingFeatureWaitlist {
	var upcomingFeatureWaitlists []UpcomingFeatureWaitlist
	db.Find(&upcomingFeatureWaitlists)
	return upcomingFeatureWaitlists
}

func getUpcomingFeatureWaitlistByID(db *gorm.DB, id string) *UpcomingFeatureWaitlist {
	var upcomingFeatureWaitlist UpcomingFeatureWaitlist
	db.First(&upcomingFeatureWaitlist, "id = ?", id)

	if upcomingFeatureWaitlist.ID == "" {
		return nil
	}

	return &upcomingFeatureWaitlist
}

func createUpcomingFeatureWaitlist(db *gorm.DB, target string, email string) {
	upcomingFeatureWaitlist := UpcomingFeatureWaitlist{ID: uuid.New().String(), Target: target, Email: email}
	result := db.Create(&upcomingFeatureWaitlist)
	if result.Error != nil {
		log.Println("upcomingFeatureWaitlist 생성 실패:", result.Error)
	} else {
		fmt.Printf("upcomingFeatureWaitlist 생성 성공: ID = %d\n", upcomingFeatureWaitlist.ID)
	}
}

func updateUpcomingFeatureWaitlist(db *gorm.DB, id string, target string, email string) {
	var upcomingFeatureWaitlist UpcomingFeatureWaitlist
	db.First(&upcomingFeatureWaitlist, "id = ?", id)
	upcomingFeatureWaitlist.Email = email
	upcomingFeatureWaitlist.Target = target
	db.Save(&upcomingFeatureWaitlist)
	fmt.Printf("upcomingFeatureWaitlist ID %d 정보 수정 완료\n", id)
}

func deleteUpcomingFeatureWaitlist(db *gorm.DB, id string) {
	db.Delete(&UpcomingFeatureWaitlist{}, "id = ?", id)
	fmt.Printf("upcomingFeatureWaitlist ID %d 삭제 완료\n", id)
}
