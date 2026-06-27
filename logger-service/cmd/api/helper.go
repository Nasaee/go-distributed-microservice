package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type JsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func (app *Application) readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1048576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)
	if err != nil {
		return err
	}

	/*
		หลังจากที่ระบบแปลงข้อมูล JSON ชุดแรกเข้าสู่ตัวแปร data ไปเรียบร้อยแล้วในบรรทัดก่อนหน้า คำสั่งนี้จะทำการ Decode ข้อมูลถัดไปที่เหลืออยู่ใน Request Body ทันที โดยใช้ตัวรับเป็น Struct ว่าง (&struct{}{}) ซึ่งเราไม่ได้ต้องการนำข้อมูลนี้ไปใช้งานจริง แต่มีไว้เพื่อตรวจจับว่า "ยังมีข้อมูลหลงเหลืออยู่ใน Stream อีกหรือไม่"
	*/
	err = dec.Decode(&struct{}{})
	/*
		หาก Request Body ไม่มีข้อมูลใดๆ หลงเหลืออยู่แล้ว (เป็นสิ่งที่เราคาดหวังสำหรับการส่ง JSON เดี่ยว) ตัว Decode จะชนกับจุดสิ้นสุดของไฟล์และส่งข้อผิดพลาด io.EOF (End of File) กลับมา
		แต่ถ้าค่า err ที่ได้รับ ไม่ใช่ io.EOF (เช่น ได้ค่าเป็น nil ซึ่งหมายความว่ามีข้อมูล JSON ชุดที่สองส่งพ่วงมาด้วย หรือได้ Error รูปแบบอื่น เช่น Syntax Error) แสดงว่า Request Body มีความผิดปกติเกิดขึ้น
	*/
	if err != io.EOF {
		return errors.New("body must have only a single json value")
	}

	return nil
}

func (app *Application) writeJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	/*
		pattern headers[0] นี้เป็นที่นิยมใน Go เพราะ simple และ caller ไม่มีเหตุผลต้องส่งมากกว่า 1 ตัว
	*/
	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		return err
	}

	return nil
}

func (app *Application) errorJSON(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload JsonResponse
	payload.Error = true
	payload.Message = err.Error()

	return app.writeJSON(w, statusCode, payload)
}
