package seller

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

func TestDoesSellerEmailExist(t *testing.T) {
	db, err := store.SetupTestDB("../../.env")
	assert.NoError(t, err)

	addDummyAccounts(db)

	testEmail := "test@gmail.com"
	testEmail2 := "test2@gmail.com"
	testEmail3 := "Test@gmail.com"
	testEmail4 := ""

	//Test 1: Positive result when email exists in the database
	res := doesSellerEmailExist(db, testEmail)
	assert.Equal(t, true, res)

	//Test 2: Positive result when email exists in the database
	res = doesSellerEmailExist(db, testEmail2)
	assert.Equal(t, true, res)

	//Test 3: Negative result when email is similar to one in the email but
	//is not the same
	res = doesSellerEmailExist(db, testEmail3)
	assert.Equal(t, false, res)

	//Test 4: Negative result when test email is an empty string
	res = doesSellerEmailExist(db, testEmail4)
	assert.NotEqual(t, true, res)

	store.CloseDB(db)
}

func TestDoesSellerNameExist(t *testing.T) {
	db, err := store.SetupTestDB("../../.env")
	assert.NoError(t, err)

	addDummyAccounts(db)

	testEmail := "Test1"
	testEmail2 := "Test2"
	testEmail3 := "test3"
	testEmail4 := ""

	//Test 1: Positive result when seller_name exists in the database
	res := doesSellerNameExist(db, testEmail)
	assert.Equal(t, true, res)

	//Test 2: Positive result when seller_name exists in the database
	res = doesSellerNameExist(db, testEmail2)
	assert.Equal(t, true, res)

	//Test 3: Negative result when seller_name is similar to one in the database but
	//is not the same
	res = doesSellerNameExist(db, testEmail3)
	assert.Equal(t, false, res)

	//Test 4: Negative result when test seller_name is an empty string
	res = doesSellerNameExist(db, testEmail4)
	assert.NotEqual(t, true, res)

	store.CloseDB(db)
}

func TestDoesSellerExist(t *testing.T) {
	db, err := store.SetupTestDB("../../.env")
	assert.NoError(t, err)

	sellerIds := addDummyAccounts(db)

	// Test 1: Positive result when seller_id exists in the database
	res := DoesSellerExist(db, sellerIds[0])
	assert.Equal(t, true, res)

	// Test 2: Positive result when seller_id exists in the database
	res = DoesSellerExist(db, sellerIds[1])
	assert.Equal(t, true, res)

	// Test 3: Negative result when seller_id is similar to one in the database but
	// is not the same
	res = DoesSellerExist(db, sellerIds[2]+"0")
	assert.Equal(t, false, res)

	// Test 4: Negative result when test seller_id is an empty string
	res = DoesSellerExist(db, "")
	assert.NotEqual(t, true, res)

	store.CloseDB(db)
}

func TestSellerLogin(t *testing.T) {
	db, err := store.SetupTestDB("../../.env")
	addDummyAccounts(db)

	testLogin1 := data.UserLoginData{Email: "test@gmail.com", Password: "Test1234"}
	testLogin2 := data.UserLoginData{Email: "test2@gmail.com", Password: "Test1234"}
	testLogin3 := data.UserLoginData{Email: "test@gmail.com", Password: "Test12345"}
	testLogin4 := data.UserLoginData{Email: "test8@gmail.com", Password: "Test1234"}
	testLogin5 := data.UserLoginData{Email: "", Password: "Test1234"}
	testLogin6 := data.UserLoginData{Email: "test@gmail.com", Password: ""}

	//Test 1: Positive Test where username and password are both correct
	res, err := SellerLogin(db, testLogin1)
	assert.Empty(t, err)
	assert.Equal(t, res.Email, testLogin1.Email)

	//Test 2: Positive Test where username and password are both correct
	res, err = SellerLogin(db, testLogin2)
	assert.Empty(t, err)
	assert.Equal(t, res.Email, testLogin2.Email)

	//Test 3: Negative Test where username is correct but password is incorrect
	res, err = SellerLogin(db, testLogin3)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "Incorrect user email or password!")

	//Test 4: Negative Test where username is incorrect
	res, err = SellerLogin(db, testLogin4)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "Incorrect user email or password!")

	//Test 5: Negative Test where username is empty string
	res, err = SellerLogin(db, testLogin5)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "Incorrect user email or password!")

	//Test 6: Negative Test where password is empty string
	res, err = SellerLogin(db, testLogin6)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "Incorrect user email or password!")

	store.CloseDB(db)
}

func TestSellerSignUp(t *testing.T) {
	db, err := store.SetupTestDB("../../.env")

	testSignup1 := data.SellerSignUpData{Email: "test@gmail.com", Password: "Test1234", SellerName: "Test1"}
	testSignup2 := data.SellerSignUpData{Email: "test2@gmail.com", Password: "Test1234", SellerName: "Test2"}
	testSignup3 := data.SellerSignUpData{Email: "test@gmail.com", Password: "Test1234", SellerName: "Test3"}
	testSignup4 := data.SellerSignUpData{Email: "test@gmail9.com", Password: "Test1234", SellerName: "Test2"}

	//Test 1: Positive Test case, where signup is successful
	res, err := SellerSignUp(db, testSignup1)
	assert.Empty(t, err)
	assert.Equal(t, res.Email, testSignup1.Email)

	//Test 2: Positive Test case, where signup is successful
	res, err = SellerSignUp(db, testSignup2)
	assert.Empty(t, err)
	assert.Equal(t, res.Email, testSignup2.Email)

	//Test 3: Negative Test case, where email already exists
	res, err = SellerSignUp(db, testSignup3)
	assert.Error(t, err)

	//Test 3: Negative Test case, where seller_name already exists
	res, err = SellerSignUp(db, testSignup4)
	assert.Error(t, err)

	store.CloseDB(db)
}

func addDummyAccounts(db *sql.DB) []string {
	var dummyAccounts []data.SellerSignUpData = []data.SellerSignUpData{{Email: "test@gmail.com", Password: "Test1234", SellerName: "Test1"},
		{Email: "test2@gmail.com", Password: "Test1234", SellerName: "Test2"}, {Email: "test3@gmail.com", Password: "Test1234", SellerName: "Test3"}}
	var sellerIds []string
	for i := 0; i < len(dummyAccounts); i++ {
		var sellerId string
		query := `INSERT INTO sellers(email, password, seller_name) VALUES ($1,$2,$3) RETURNING seller_id;`
		hashedPwd, _ := utils.HashAndSalt([]byte(dummyAccounts[i].Password))
		db.QueryRowContext(context.Background(), query, dummyAccounts[i].Email, hashedPwd, dummyAccounts[i].SellerName).Scan(&sellerId)
		sellerIds = append(sellerIds, sellerId)
	}

	return sellerIds
}
