package models

import (
	"gopkg.in/mgo.v2/bson"
	"gost/dbmodels"
	"time"
)

// Struct representing an user account. This is a database dbmodels
type ApplicationUser struct {
	Id bson.ObjectId `json:"id"`

	Email                   string    `json:"email"`
	Password                string    `json:"password"`
	AccountType             int       `json:"accountType"`
	ResetPasswordToken      string    `json:"resetPasswordTokenToken"`
	ResetPasswordExpireDate time.Time `json:"resetPasswordExpireDate"`
}

func (user *ApplicationUser) PopConstrains() {
	// Nothing to do here for now
}

func (user *ApplicationUser) Expand(dbUser *dbmodels.ApplicationUser) {
	user.Id = dbUser.Id
	user.Email = dbUser.Email
	user.Password = dbUser.Password
	user.AccountType = dbUser.AccountType
	user.ResetPasswordToken = dbUser.ResetPasswordToken
	user.ResetPasswordExpireDate = dbUser.ResetPasswordExpireDate

	user.PopConstrains()
}

func (user *ApplicationUser) Collapse() *dbmodels.ApplicationUser {
	dbUser := dbmodels.ApplicationUser{
		Id:                      user.Id,
		Email:                   user.Email,
		Password:                user.Password,
		AccountType:             user.AccountType,
		ResetPasswordToken:      user.ResetPasswordToken,
		ResetPasswordExpireDate: user.ResetPasswordExpireDate,
	}

	return &dbUser
}
