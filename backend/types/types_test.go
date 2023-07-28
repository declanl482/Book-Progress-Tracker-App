package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTypes(t *testing.T) {

	// (1) Test request inputs for a User type.
	t.Run("TestUser", func(t *testing.T) {
		users := []struct {
			user    User
			isValid bool
		}{
			{
				// Valid user.
				user: User{
					Username: "foo",
					Email:    "foo@bar.com",
					Password: "foo",
					Books:    []Book{},
				},
				isValid: true,
			},
			{
				// Missing username.
				user: User{
					Username: "",
					Email:    "foo@bar.com",
					Password: "foo",
					Books:    []Book{},
				},
				isValid: false,
			},
			{
				// Missing email.
				user: User{
					Username: "foo",
					Email:    "",
					Password: "foo",
					Books:    []Book{},
				},
				isValid: false,
			},
			{
				// Missing password.
				user: User{
					Username: "foo",
					Email:    "foo@bar.com",
					Password: "",
					Books:    []Book{},
				},
				isValid: false,
			},
			{
				// Invalid email variant 1.
				user: User{
					Username: "foo",
					Email:    "foobar.com",
					Password: "foo",
					Books:    []Book{},
				},
				isValid: false,
			},
			{
				// Invalid email variant 2.
				user: User{
					Username: "foo",
					Email:    "foo@barcom",
					Password: "foo",
					Books:    []Book{},
				},
				isValid: false,
			},
		}

		for _, userData := range users {
			err := userData.user.ValidateUser()
			if userData.isValid {
				// Expect no errors (err == nil).
				assert.NoError(t, err, "Expected no errors for valid user case, got: %v.", err)
			} else {
				// Expect errors (err != nil).
				assert.Error(t, err, "Expected errors for invalid user case, got: nil.")
			}
		}
	})

	// (2) Test request inputs for user credentials.
	t.Run("TestCredentials", func(t *testing.T) {
		credentials := []struct {
			creds   Credentials
			isValid bool
		}{
			{
				// Valid user credentials.
				creds: Credentials{
					Email:    "foo@bar.com",
					Password: "foo",
				},
				isValid: true,
			},
			{
				// Missing email.
				creds: Credentials{
					Email:    "",
					Password: "foo",
				},
				isValid: false,
			},
			{
				// Missing password.
				creds: Credentials{
					Email:    "foo@bar.com",
					Password: "",
				},
				isValid: false,
			},
			{
				// Invalid email variant 1.
				creds: Credentials{
					Email:    "foobar.com",
					Password: "foo",
				},
				isValid: false,
			},
			{
				// Invalid email variant 1.
				creds: Credentials{
					Email:    "foo@barcom",
					Password: "foo",
				},
				isValid: false,
			},
		}

		for _, credData := range credentials {
			err := credData.creds.ValidateCredentials()
			if credData.isValid {
				// We expect no errors (err == nil).
				assert.NoError(t, err, "Expected no error for valid credentials case, got: %v.", err)
			} else {
				// Credentials are not valid, we expect errors (err != nil).
				assert.Error(t, err, "Expected error for invalid credentials case, got: nil.")
			}
		}

	})

	// (3) Test request inputs for a Book type.
	t.Run("TestBook", func(t *testing.T) {
		books := []struct {
			book    Book
			isValid bool
		}{
			{
				// Valid book with edition.
				book: Book{
					Title:      "Sample Book 1",
					Edition:    4,
					Author:     "John Doe",
					PagesCount: 300,
					PagesRead:  150,
				},
				isValid: true,
			},
			{
				// Valid book without edition.
				book: Book{
					Title:      "Sample Book 2",
					Author:     "John Doe",
					PagesCount: 300,
					PagesRead:  150,
				},
				isValid: true,
			},
			{
				// Missing title.
				book: Book{
					Title:      "",
					Author:     "John Doe",
					PagesCount: 300,
					PagesRead:  150,
				},
				isValid: false,
			},
			{
				// Missing author.
				book: Book{
					Title:      "Sample Book 4",
					Author:     "",
					PagesCount: 300,
					PagesRead:  150,
				},
				isValid: false,
			},
			{
				// Invalid pages count.
				book: Book{
					Title:      "Sample Book 5",
					Author:     "John Doe",
					PagesCount: 0,
					PagesRead:  0,
				},
				isValid: false,
			},
			{
				// Invalid pages read.
				book: Book{
					Title:      "Sample Book 6",
					Author:     "John Doe",
					PagesCount: 300,
					PagesRead:  -1,
				},
				isValid: false,
			},
			{
				// pages read > pages count.
				book: Book{
					Title:      "Sample Book 7",
					Author:     "John Doe",
					PagesCount: 300,
					PagesRead:  600,
				},
				isValid: false,
			},
		}

		for _, bookData := range books {
			err := bookData.book.ValidateBook()
			if bookData.isValid {
				// We expect no errors (err == nil).
				assert.NoError(t, err, "Expected no error for valid book case, got: %v.", err)
			} else {
				// Book is not valid, we expect errors (err != nil).
				assert.Error(t, err, "Expected error for invalid book case, got: nil.")
			}
		}
	})
}
