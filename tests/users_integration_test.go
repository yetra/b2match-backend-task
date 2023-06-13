package tests

import (
	"b2match/backend/database"
	"b2match/backend/models"
	"b2match/backend/routes"
	"os"
	"reflect"

	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var router *gin.Engine

var testCompany models.Company
var testUsers []models.User

func requestJSON(content interface{}) io.Reader {
	jsonbytes, err := json.Marshal(content)
	if err != nil {
		panic(err)
	}

	return io.NopCloser(bytes.NewBuffer(jsonbytes))
}

func setUpTestData() {
	testCompany = models.Company{Name: "TestCompany", Location: "Somewhere"}

	if err := database.DB.Create(&testCompany).Error; err != nil {
		panic("couldn't set up test data")
	}

	testUsers = []models.User{
		{FirstName: "First1", LastName: "Last1", EMail: "test1@test.com", Password: "testpass", CompanyID: 1},
		{FirstName: "First2", LastName: "Last2", EMail: "test2@test.com", Password: "testpass", CompanyID: 1},
	}

	if err := database.DB.Create(&testUsers).Error; err != nil {
		panic("couldn't set up test data 1")
	}
}

func tearDownTestData() {
	if err := database.DB.Delete(&testUsers).Error; err != nil {
		panic("couldn't tear down test data")
	}
	if err := database.DB.Delete(&testCompany).Error; err != nil {
		panic("couldn't tear down test data")
	}
}

func TestMain(m *testing.M) {
	database.SetUpDB("test.db", &gorm.Config{})
	setUpTestData()

	router = routes.CreateRouter()

	exitCode := m.Run()

	tearDownTestData()

	os.Exit(exitCode)
}

func TestGetUsers(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users", nil)

	router.ServeHTTP(w, req)

	assert.EqualValues(t, http.StatusOK, w.Code)

	var users []models.User
	if err := json.Unmarshal(w.Body.Bytes(), &users); err != nil {
		assert.Fail(t, err.Error())
	}

	assert.EqualValues(t, 2, len(users))
	assert.True(t, reflect.DeepEqual(testUsers, users))
}

func TestGetUserByID(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/1", nil)

	router.ServeHTTP(w, req)

	assert.EqualValues(t, http.StatusOK, w.Code)

	var user models.User
	if err := json.Unmarshal(w.Body.Bytes(), &user); err != nil {
		assert.Fail(t, err.Error())
	}

	assert.True(t, reflect.DeepEqual(testUsers[0], user))
}

func TestGetUserByIDInvalid(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/100", nil)

	router.ServeHTTP(w, req)

	assert.EqualValues(t, http.StatusNotFound, w.Code)
}

func TestCreateUser(t *testing.T) {
	content := map[string]interface{}{
		"first_name": "TestFirstName",
		"last_name":  "TestLastName",
		"e_mail":     "test@test.com",
		"password":   "testpass",
		"company_id": 1,
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/users", requestJSON(content))

	router.ServeHTTP(w, req)

	assert.EqualValues(t, http.StatusCreated, w.Code)

	var createdUserBody models.User
	if err := json.Unmarshal(w.Body.Bytes(), &createdUserBody); err != nil {
		assert.Fail(t, err.Error())
	}

	var createdUser models.User
	if err := database.DB.Find(&createdUser, createdUserBody.ID).Error; err != nil {
		assert.Fail(t, err.Error())
	}

	assert.EqualValues(t, content["first_name"], createdUser.FirstName)
	assert.EqualValues(t, content["last_name"], createdUser.LastName)
	assert.EqualValues(t, content["e_mail"], createdUser.EMail)
	assert.EqualValues(t, content["password"], createdUser.Password)
	assert.EqualValues(t, content["company_id"], createdUser.CompanyID)

	database.DB.Delete(&createdUser)
}

func TestCreateUserMissingFields(t *testing.T) {
	invalidContents := []map[string]interface{}{
		{
			"last_name":  "TestLastName",
			"e_mail":     "testinvalid@test.com",
			"password":   "testpass",
			"company_id": 1,
		},
		{
			"first_name": "TestFirstName",
			"e_mail":     "testinvalid@test.com",
			"password":   "testpass",
			"company_id": 1,
		},
		{
			"first_name": "TestFirstName",
			"last_name":  "TestLastName",
			"password":   "testpass",
			"company_id": 1,
		},
		{
			"first_name": "TestFirstName",
			"last_name":  "TestLastName",
			"e_mail":     "testinvalid@test.com",
			"company_id": 1,
		},
		{
			"first_name": "TestFirstName",
			"last_name":  "TestLastName",
			"e_mail":     "testinvalid@test.com",
			"password":   "testpass",
		},
	}

	for _, content := range invalidContents {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/users", requestJSON(content))

		router.ServeHTTP(w, req)

		assert.EqualValues(t, http.StatusBadRequest, w.Code)
	}
}

func TestCreateUserInvalidCompanyID(t *testing.T) {
	invalidContent := map[string]interface{}{
		"first_name": "TestFirstName",
		"last_name":  "TestLastName",
		"e_mail":     "testinvalid@test.com",
		"password":   "testpass",
		"company_id": 100,
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/users", requestJSON(invalidContent))

	router.ServeHTTP(w, req)

	assert.EqualValues(t, http.StatusNotFound, w.Code)
}

func TestUpdateUser(t *testing.T) {
	content := map[string]interface{}{
		"first_name": "UpdatedFirstName",
		"last_name":  "UpdatedLastName",
		"location":   "New Location",
		"about":      "New About...",
		"password":   "newpass",
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/users/1", requestJSON(content))

	router.ServeHTTP(w, req)

	assert.EqualValues(t, http.StatusOK, w.Code)

	var updatedUserBody models.User
	if err := json.Unmarshal(w.Body.Bytes(), &updatedUserBody); err != nil {
		assert.Fail(t, err.Error())
	}

	var updatedUser models.User
	if err := database.DB.Find(&updatedUser, 1).Error; err != nil {
		assert.Fail(t, err.Error())
	}

	assert.EqualValues(t, content["first_name"], updatedUser.FirstName)
	assert.EqualValues(t, content["last_name"], updatedUser.LastName)
	assert.EqualValues(t, content["location"], updatedUser.Location)
	assert.EqualValues(t, content["about"], updatedUser.About)
	assert.EqualValues(t, content["password"], updatedUser.Password)
}

func TestUpdateUserInvalidID(t *testing.T) {
	content := map[string]interface{}{
		"first_name": "UpdatedFirstName",
		"last_name":  "UpdatedLastName",
		"location":   "New Location",
		"about":      "New About...",
		"password":   "newpass",
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/users/100", requestJSON(content))

	router.ServeHTTP(w, req)

	assert.EqualValues(t, http.StatusNotFound, w.Code)
}

func TestDeleteUser(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/users/1", nil)

	router.ServeHTTP(w, req)

	assert.EqualValues(t, http.StatusNoContent, w.Code)
	assert.Error(t, database.DB.First(&models.User{}, 1).Error)
}

func TestDeleteUserInvalidID(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/users/100", nil)

	router.ServeHTTP(w, req)

	assert.EqualValues(t, http.StatusNotFound, w.Code)
}
