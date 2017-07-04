package accounts

import (
	"fmt"
	"net/http"
	"time"

	"github.com/hypnoglow/pascont/account"
	"github.com/hypnoglow/pascont/kit"
	"github.com/hypnoglow/pascont/kit/form"
	"github.com/hypnoglow/pascont/kit/schema"
)

// PostAccounts is a handler for:
// POST /accounts
func (c RestController) PostAccounts(w http.ResponseWriter, req *http.Request) {
	var accForm postAccountForm
	form.PopulateFormFromJSON(req.Body, &accForm)
	if !accForm.Validate() {
		kit.RespondWithFormErrors(w, http.StatusBadRequest, accForm.ValidationErrors())
		return
	}

	// It must be not possible to create multiple accounts with same name.
	exists, err := c.accountRepo.Exists(accForm.Name)
	if err != nil {
		c.errorLogger.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if exists {
		kit.RespondWithError(w, http.StatusConflict, schema.NewError(
			fmt.Sprintf("Account with name %s already exists", accForm.Name),
			"name",
			accForm.Name,
		))
		return
	}

	passwordHash, err := c.hasher.GenerateHashFromPassword([]byte(accForm.Password))
	if err != nil {
		c.errorLogger.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	app := account.NewApplication(accForm.Name, passwordHash, time.Now())
	acc, err := c.accountRepo.Accept(app)
	if err != nil {
		c.errorLogger.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	kit.RespondJSON(w, http.StatusCreated, schema.NewResultBody(
		postAccountSchema{
			ID:        acc.ID,
			Name:      acc.Name,
			CreatedAt: acc.CreatedAt,
			UpdatedAt: acc.UpdatedAt,
		},
		nil,
	))
}

const (
	// nameMinLen is min len of the account name in bytes
	nameMinLen = 4

	// nameMaxLen is max len of the account name in bytes.
	nameMaxLen = 64

	// passwordMinLen is min len of the account password in bytes.
	passwordMinLen = 8
)

type postAccountForm struct {
	form.BaseForm
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (f *postAccountForm) Validate() bool {
	if len(f.Name) < nameMinLen {
		f.AddError(fmt.Sprintf("Name must be at least %d bytes", nameMinLen), "name", f.Name)
	}
	if len(f.Name) > nameMaxLen {
		f.AddError(fmt.Sprintf("Name must be at most %d bytes", nameMaxLen), "name", f.Name)
	}

	if len(f.Password) < passwordMinLen {
		f.AddError(fmt.Sprintf("Password must be at least %d bytes", passwordMinLen), "password", nil)
	}

	return len(f.ValidationErrors()) == 0
}

type postAccountSchema struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
