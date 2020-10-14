package main

import (
	"email/models"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"strings"
	"time"

	"github.com/jasonlvhit/gocron"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	godotenv.Load(".env")
	// migrator()
	// insert()
	// task()

	gocron.Every(30).Second().Do(task)
	<-gocron.Start()
}

func connectDB() *gorm.DB {

	user := os.Getenv("USER_DATABASE")
	pass := os.Getenv("PASS_DATABASE")
	port := os.Getenv("PORT_DATABASE")
	host := os.Getenv("HOST_DATABASE")
	name := os.Getenv("NAME_DATABASE")

	conn := user + ":" + pass + "@tcp(" + host + ":" + port + ")/" + name + "?charset=utf8mb4&parseTime=True&loc=Local"

	db, errConn := gorm.Open(mysql.Open(conn), &gorm.Config{})
	if errConn != nil {
		panic("Failed to connect database")
	}

	return db
}

func sendMail(to []string, subject string, message string) error {
	auth := smtp.PlainAuth("", os.Getenv("SMTP_USER"), os.Getenv("SMTP_PASS"), os.Getenv("SMTP_HOST"))

	message = "From: " + os.Getenv("SMTP_USER") + "\n" +
		"To: " + strings.Join(to, ",") + "\n" +
		"Subject: " + subject + "\n" +
		message

	smtpAddr := os.Getenv("SMTP_HOST") + ":" + os.Getenv("SMTP_PORT")
	log.Println("Sending...")
	log.Println("Message", message)
	return smtp.SendMail(smtpAddr, auth, os.Getenv("SMTP_USER"), to, []byte(message))
}

func migrator() {
	var check bool
	db := connectDB()
	check = db.Migrator().HasTable(&models.Email{})
	if !check {
		db.Migrator().CreateTable(&models.Email{})
		fmt.Println("Create table emails")
	}
}

func insert() {
	t, _ := time.Parse(time.RFC3339, "2020-10-14T16:00:00+07:00")
	email := models.Email{
		To:      "yovilatte@gmail.com",
		Subject: "Halo Test Email",
		Message: "Halo. Ini adalah test email",
		SendAt:  t,
		IsSent:  false,
	}

	db := connectDB()

	db.Create(&email)
}

func task() {
	log.Println("I do task")
	var emails []models.Email

	db := connectDB()

	db.Where("is_sent = ?", false).Find(&emails)
	nowNano := time.Now().Unix()

	for _, value := range emails {
		var to []string
		to = append(to, value.To)
		log.Println(value.SendAt.Unix(), nowNano, value.SendAt.Unix() <= nowNano)
		if value.SendAt.Unix() <= nowNano {
			log.Println("I send email")
			err := sendMail(to, value.Subject, value.Message)
			if err != nil {
				log.Println("Send failed")
			} else {
				log.Println("Send success")
				value.IsSent = true
				db.Save(&value)
			}
		}

	}
}
