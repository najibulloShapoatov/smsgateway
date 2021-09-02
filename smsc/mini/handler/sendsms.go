package handler

import (
	"context"
	"net/http"
	"smsc/mini/middleware"
	"smsc/pkg/models"
	"smsc/pkg/utils"
	"time"
)

func Sent(w http.ResponseWriter, r *http.Request) {

	customerID, err := middleware.Authentication(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	aplhanumeric, err := (&models.Alphanumeric{CustomerID: customerID, Name: r.URL.Query().Get("from")}).Get(r.Context())
	if err != nil {
		http.Error(w, "alphanumeric not found ", http.StatusBadRequest)
		return
	}
	phone := r.URL.Query().Get("to")
	if len(phone) <= 0 {
		http.Error(w, "phone not valid", http.StatusBadRequest)
	}
	text := r.URL.Query().Get("text")

	sentsms, err := (&models.SentSms{
		Alphanumeric: aplhanumeric.Name,
		Phone:        phone,
		Content:      text,
	}).Insert(r.Context())

	go func() {
		sentsms.Status = utils.SentSms(sentsms.Alphanumeric, sentsms.Phone, sentsms.Content)
		sentsms.UpdatedAt = time.Now()
		sentsms.Update(context.Background())
	}()

	http.Error(w, http.StatusText(http.StatusAccepted), http.StatusAccepted)
	return

}
