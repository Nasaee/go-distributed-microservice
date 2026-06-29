package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
	Log    LogPayload  `json:"log,omitempty"`
	Mail   MailPayload `json:"mail,omitempty"`
}

type MailPayload struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Application) broker(w http.ResponseWriter, r *http.Request) {
	payload := JsonResponse{
		Error:   false,
		Message: "Broker Service is running",
	}

	err := app.writeJSON(w, http.StatusOK, payload)
	if err != nil {
		http.Error(w, "internal server eroror", http.StatusInternalServerError)
		return
	}
}

func (app *Application) handleSubmission(w http.ResponseWriter, r *http.Request) {
	var RequestPayload RequestPayload

	err := app.readJSON(w, r, &RequestPayload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	switch RequestPayload.Action {
	case "auth":
		app.authenticate(w, RequestPayload.Auth)
	case "log":
		app.log(w, RequestPayload.Log)
	case "mail":
		app.mail(w, RequestPayload.Mail)
	default:
		app.errorJSON(w, errors.New("invalid action"))
	}
}

func (app *Application) authenticate(w http.ResponseWriter, p AuthPayload) {
	// 1. แปลงข้อมูล AuthPayload (อีเมลและรหัสผ่าน) เป็น JSON format พร้อมจัดรูปแบบให้สวยงาม
	jsonData, _ := json.MarshalIndent(p, "", "\t")

	// 2. สร้าง HTTP Request (POST) เพื่อเตรียมส่งข้อมูลไปยัง Authentication Microservice
	request, err := http.NewRequest("POST", "http://authentication:8081/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// 3. เรียกใช้งาน HTTP Client เพื่อทำการยิง Request ไปยัง Microservice
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer response.Body.Close()

	// 4. ตรวจสอบสถานะการตอบกลับจาก Authentication Microservice
	if response.StatusCode == http.StatusUnauthorized {
		app.errorJSON(w, errors.New("invalid credentials"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("error calling auth service"))
		return
	}

	var serviceResponse JsonResponse

	// 5. อ่านและถอดรหัส (Decode) JSON ผลลัพธ์ที่ได้มาจาก Authentication Service
	err = json.NewDecoder(response.Body).Decode(&serviceResponse)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// 6. ตรวจสอบเงื่อนไขข้อผิดพลาดเชิงตรรกะที่ส่งกลับมาจากเซอร์วิส
	if serviceResponse.Error {

		app.errorJSON(w, err, http.StatusUnauthorized)
		return
	}

	// 7. เมื่อการยืนยันตัวตนสำเร็จ ส่งผลลัพธ์กลับไปยัง Client ของ Broker Service
	var payload JsonResponse
	payload.Error = false
	payload.Message = "Authorized"
	payload.Data = serviceResponse.Data

	err = app.writeJSON(w, http.StatusAccepted, payload)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (app *Application) log(w http.ResponseWriter, p LogPayload) {
	jsonData, _ := json.MarshalIndent(p, "", "\t")

	request, err := http.NewRequest("POST", "http://logger:8082/log", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, err)
		return
	}

	var payload JsonResponse
	payload.Error = false
	payload.Message = "Logged"

	err = app.writeJSON(w, http.StatusAccepted, payload)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (app *Application) mail(w http.ResponseWriter, p MailPayload) {
	jsonData, _ := json.MarshalIndent(p, "", "\t")

	request, err := http.NewRequest("POST", "http://mailer:8080/send", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resposne, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer resposne.Body.Close()

	if resposne.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("error calling mailer service"))
		return
	}

	var payload JsonResponse
	payload.Error = false
	payload.Message = "Message sent to " + p.To

	app.writeJSON(w, http.StatusAccepted, payload)
}
