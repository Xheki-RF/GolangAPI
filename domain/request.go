// nolint
package domain

type Email struct {
	Email string `validate:"required,email"`
}

type UpdateEmail struct {
	OldEmail string `validate:"required,email" json:"old_email"`
	NewEmail string `validate:"required,email,nefield=OldEmail" json:"new_email"`
}

type DeleteEmail struct {
	Email string `validate:"required,email"`
}
