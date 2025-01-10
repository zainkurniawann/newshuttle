package utils

import (
	"context"
	"errors"
	"log"

	"firebase.google.com/go/v4"
    "firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

var FirebaseApp *firebase.App

func InitFirebase() {
	opt := option.WithCredentialsFile("./service-account.json") // Ganti dengan path ke file JSON service account
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing Firebase app: %v", err)
	}
	FirebaseApp = app
}

// Send a notification to a device
func SendNotification(userUUID, title, status string) error {
    deviceToken, err := getDeviceToken(userUUID)
    if err != nil {
        return err
    }

    client, err := FirebaseApp.Messaging(context.Background())
    if err != nil {
        return errors.New("fcm: failed to get Firebase Messaging client")
    }

    var body string
    switch status {
    case "home":
        body = "Your student is at home."
    case "waiting_to_be_taken_to_school":
        body = "The school driver is on the way to pick your children up."
    case "going_to_school":
        body = "Your children are on the way to school."
    case "at_school":
        body = "Your children have arrived at school."
    case "waiting_to_be_taken_to_home":
        body = "The school driver is on the way to take your children home."
    case "going_to_home":
        body = "Your children are on the way home."
    default:
        return errors.New("fcm: invalid status")
    }

    message := &messaging.Message{
        Notification: &messaging.Notification{
            Title: title,
            Body:  body,
        },
        Token: deviceToken,
    }

    _, err = client.Send(context.Background(), message)
    if err != nil {
        return errors.New("fcm: failed to send message")
    }
    return nil
}

// Get device token from the database
func getDeviceToken(userUUID string) (string, error) {
    var deviceToken string
    query := "SELECT device_token FROM fcm_tokens WHERE user_uuid = $1"
    log.Println("Query yang dijalankan:", query, "dengan userUUID:", userUUID) // Logging query sebelum eksekusi
    err := db.QueryRow(query, userUUID).Scan(&deviceToken)
    if err != nil {
        log.Println("Error saat mengambil token:", err) // Logging error jika terjadi
        return "", errors.New("fcm: failed to get device token")
    }
    log.Println("Device token yang didapat:", deviceToken) // Logging hasil device token
    return deviceToken, nil
}