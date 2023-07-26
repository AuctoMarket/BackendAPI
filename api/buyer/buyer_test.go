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

	addDummyAccounts(db)

	testEmail := "test@gmail.com"
	testEmail2 := "test2@gmail.com"
	testEmail3 := "Test@gmail.com"
	testEmail4 := ""

	//Test 1: Positive result when email exists in the database
	res, err := doesBuyerEmailExist(db, testEmail)
	assert.NoError(t, err)
	assert.Equal(t, res, true)

	//Test 2: Positive result when email exists in the database
	res, err = doesBuyerEmailExist(db, testEmail2)
	assert.NoError(t, err)
	assert.Equal(t, res, true)

	//Test 3: Negative result when email is similar to one in the email but
	//is not the same
	res, err = doesBuyerEmailExist(db, testEmail3)
	assert.NoError(t, err)
	assert.NotEqual(t, res, true)

	//Test 4: Negative result when test email is an empty string
	res, err = doesBuyerEmailExist(db, testEmail4)
	assert.NoError(t, err)
	assert.NotEqual(t, res, true)

	store.CleaupTestDB(db)
}

func TestBuyerLogin(t *testing.T) {
	db, err := store.SetupTestDB()

	addDummyAccounts(db)

	testLogin1 := data.LoginData{Email: "test@gmail.com", Password: "Test1234"}
	testLogin2 := data.LoginData{Email: "test2@gmail.com", Password: "Test1234"}
	testLogin3 := data.LoginData{Email: "test@gmail.com", Password: "Test12345"}
	testLogin4 := data.LoginData{Email: "test8@gmail.com", Password: "Test1234"}
	testLogin5 := data.LoginData{Email: "", Password: "Test1234"}
	testLogin6 := data.LoginData{Email: "test@gmail.com", Password: ""}

	//Test 1: Positive Test where username and password are both correct
	res, err := BuyerLogin(db, testLogin1)
	assert.NoError(t, err)
	assert.Equal(t, res.Email, testLogin1.Email)

	//Test 2: Positive Test where username and password are both correct
	res, err = BuyerLogin(db, testLogin2)
	assert.NoError(t, err)
	assert.Equal(t, res.Email, testLogin2.Email)

	//Test 3: Negative Test where username is correct but password is incorrect
	res, err = BuyerLogin(db, testLogin3)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "Incorrect username or password")

	//Test 4: Negative Test where username is incorrect
	res, err = BuyerLogin(db, testLogin4)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "Incorrect username or password")

	//Test 5: Negative Test where username is empty string
	res, err = BuyerLogin(db, testLogin5)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "Incorrect username or password")

	//Test 6: Negative Test where password is empty string
	res, err = BuyerLogin(db, testLogin6)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "Incorrect username or password")

	store.CleaupTestDB(db)
}

func TestBuyerSignUp(t *testing.T) {
	db, err := store.SetupTestDB()

	testSignup1 := data.SignUpData{Email: "test@gmail.com", Password: "Test1234"}
	testSignup2 := data.SignUpData{Email: "test2@gmail.com", Password: "Test1234"}
	testSignup3 := data.SignUpData{Email: "test@gmail.com", Password: "Test1234"}

	//Test 1: Positive Test case, where signup is successful
	res, err := BuyerSignUp(db, testSignup1)
	assert.NoError(t, err)
	assert.Equal(t, res.Email, testSignup1.Email)

	//Test 2: Positive Test case, where signup is successful
	res, err = BuyerSignUp(db, testSignup2)
	assert.NoError(t, err)
	assert.Equal(t, res.Email, testSignup2.Email)

	//Test 3: Negative Test case, where email already exists
	res, err = BuyerSignUp(db, testSignup3)
	assert.Error(t, err)

	store.CleaupTestDB(db)
}

func addDummyAccounts(db *sql.DB) {
	var dummyAccounts []data.SignUpData = []data.SignUpData{{Email: "test@gmail.com", Password: "Test1234"},
		{Email: "test2@gmail.com", Password: "Test1234"}, {Email: "test3@gmail.com", Password: "Test1234"}}

	for i := 0; i < len(dummyAccounts); i++ {
		query := `INSERT INTO buyers(email, password) VALUES ($1,$2);`
		hashedPwd, _ := utils.HashAndSalt([]byte(dummyAccounts[i].Password))
		db.ExecContext(context.Background(), query, dummyAccounts[i].Email, hashedPwd)
	}
}
