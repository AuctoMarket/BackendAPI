package buyer

import (
	"BackendAPI/data"
	"BackendAPI/store"
	"BackendAPI/utils"
	"context"
	"database/sql"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestDoesBuyerEmailExist(t *testing.T) {

	db, err := store.SetupTestDB()
	assert.NoError(t, err)

	addDummyAccounts(db)

	testEmail := "test@gmail.com"
	testEmail2 := "test2@gmail.com"
	testEmail3 := "Test@gmail.com"
	testEmail4 := ""

	//Test 1: Positive result when email exists in the database
	res := doesBuyerEmailExist(db, testEmail)
	assert.Equal(t, true, res)

	//Test 2: Positive result when email exists in the database
	res = doesBuyerEmailExist(db, testEmail2)
	assert.Equal(t, true, res)

	//Test 3: Negative result when email is similar to one in the email but
	//is not the same
	res = doesBuyerEmailExist(db, testEmail3)
	assert.Equal(t, false, res)

	//Test 4: Negative result when test email is an empty string
	res = doesBuyerEmailExist(db, testEmail4)
	assert.NotEqual(t, true, res)

	store.CloseDB(db)
}

func TestBuyerLogin(t *testing.T) {
	db, err := store.SetupTestDB()
	addDummyAccounts(db)

	testLogin1 := data.BuyerLoginData{Email: "test@gmail.com", Password: "Test1234"}
	testLogin2 := data.BuyerLoginData{Email: "test2@gmail.com", Password: "Test1234"}
	testLogin3 := data.BuyerLoginData{Email: "test@gmail.com", Password: "Test12345"}
	testLogin4 := data.BuyerLoginData{Email: "test8@gmail.com", Password: "Test1234"}
	testLogin5 := data.BuyerLoginData{Email: "", Password: "Test1234"}
	testLogin6 := data.BuyerLoginData{Email: "test@gmail.com", Password: ""}

	//Test 1: Positive Test where username and password are both correct
	res, err := BuyerLogin(db, testLogin1)
	assert.Empty(t, err)
	assert.Equal(t, res.Email, testLogin1.Email)

	//Test 2: Positive Test where username and password are both correct
	res, err = BuyerLogin(db, testLogin2)
	assert.Empty(t, err)
	assert.Equal(t, res.Email, testLogin2.Email)

	//Test 3: Negative Test where username is correct but password is incorrect
	res, err = BuyerLogin(db, testLogin3)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "Incorrect user email or password!")

	//Test 4: Negative Test where username is incorrect
	res, err = BuyerLogin(db, testLogin4)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "Incorrect user email or password!")

	//Test 5: Negative Test where username is empty string
	res, err = BuyerLogin(db, testLogin5)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "Incorrect user email or password!")

	//Test 6: Negative Test where password is empty string
	res, err = BuyerLogin(db, testLogin6)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "Incorrect user email or password!")

	store.CloseDB(db)
}

func TestBuyerSignUp(t *testing.T) {
	db, err := store.SetupTestDB()

	testSignup1 := data.BuyerSignUpData{Email: "test@gmail.com", Password: "Test1234"}
	testSignup2 := data.BuyerSignUpData{Email: "test2@gmail.com", Password: "Test1234"}
	testSignup3 := data.BuyerSignUpData{Email: "test@gmail.com", Password: "Test1234"}

	//Test 1: Positive Test case, where signup is successful
	res, err := BuyerSignUp(db, testSignup1)
	assert.Empty(t, err)
	assert.Equal(t, res.Email, testSignup1.Email)

	//Test 2: Positive Test case, where signup is successful
	res, err = BuyerSignUp(db, testSignup2)
	assert.Empty(t, err)
	assert.Equal(t, res.Email, testSignup2.Email)

	//Test 3: Negative Test case, where email already exists
	res, err = BuyerSignUp(db, testSignup3)
	assert.Error(t, err)

	store.CloseDB(db)
}

func addDummyAccounts(db *sql.DB) {
	var dummyAccounts []data.BuyerSignUpData = []data.BuyerSignUpData{{Email: "test@gmail.com", Password: "Test1234"},
		{Email: "test2@gmail.com", Password: "Test1234"}, {Email: "test3@gmail.com", Password: "Test1234"}}

	for i := 0; i < len(dummyAccounts); i++ {
		query := `INSERT INTO buyers(email, password) VALUES ($1,$2);`
		hashedPwd, _ := utils.HashAndSalt([]byte(dummyAccounts[i].Password))
		db.ExecContext(context.Background(), query, dummyAccounts[i].Email, hashedPwd)
	}
}
