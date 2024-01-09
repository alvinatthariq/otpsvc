package controllers

import (
	"github.com/alvinatthariq/otpsvc/domain"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type controller struct {
	gorm   *gorm.DB
	router *mux.Router
	domain domain.DomainItf
}

func Init(gorm *gorm.DB, router *mux.Router, domain domain.DomainItf) {
	c := &controller{
		gorm:   gorm,
		router: router,
		domain: domain,
	}

	c.Serve()
}

func (c *controller) Serve() {
	// otp
	c.router.HandleFunc("/v1/otp/request", c.RequestOTP).Methods("POST")
	c.router.HandleFunc("/v1/otp/validate", c.ValidateOTP).Methods("POST")
}
