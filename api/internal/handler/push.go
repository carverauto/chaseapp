package handler

import (
	"net/http"
)

// SubscribePush registers a device for push notifications.
// POST /api/v1/push/subscribe
func SubscribePush(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement push subscription
	Error(w, http.StatusNotImplemented, "Not implemented")
}

// UnsubscribePush removes a device from push notifications.
// POST /api/v1/push/unsubscribe
func UnsubscribePush(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement push unsubscription
	Error(w, http.StatusNotImplemented, "Not implemented")
}

// GetSafariPushPackage generates a Safari push notification package.
// GET /api/v1/push/safari-package
func GetSafariPushPackage(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement Safari push package generation
	Error(w, http.StatusNotImplemented, "Not implemented")
}
