package rest

import (
	"fmt"
	"github.com/mercadolibre/golang-restclient/rest"
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	fmt.Println("about to start test cases...")
	rest.StartMockupServer()
	os.Exit(m.Run())
}

func TestLoginUserTimeoutFromApi(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		URL:          "http://localhost:8080/users/login",
		HTTPMethod:   http.MethodPost,
		ReqBody:      `{"email":"email@gmail.com","password":"password"}`,
		RespHTTPCode: -1,
		RespBody:     `{}`,
	})
	repository := usersRepository{}

	user, err := repository.LoginUser("email@gmail.com", "password")

	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status)
	assert.EqualValues(t, "invalid restclient response when trying to login user", err.Message)
}

func TestLoginUserInvalidErrorInterface(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		URL:          "http://localhost:8080/users/login",
		HTTPMethod:   http.MethodPost,
		ReqBody:      `{"email":"email@gmail.com","password":"password"}`,
		RespHTTPCode: http.StatusNotFound,
		RespBody:     `{"message": "invalid login credentials", "status": "404", "error": "not_found"}`,
	})
	repository := usersRepository{}

	user, err := repository.LoginUser("email@gmail.com", "password")

	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status)
	assert.EqualValues(t, "invalid error interface when trying to login user", err.Message)
}

func TestLoginUserInvalidLoginCredentials(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		URL:          "http://localhost:8080/users/login",
		HTTPMethod:   http.MethodPost,
		ReqBody:      `{"email":"email@gmail.com","password":"password"}`,
		RespHTTPCode: http.StatusNotFound,
		RespBody:     `{"message": "invalid login credentials", "status": 404, "error": "not_found"}`,
	})
	repository := usersRepository{}

	user, err := repository.LoginUser("email@gmail.com", "password")

	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusNotFound, err.Status)
	assert.EqualValues(t, "invalid error interface when trying to login user", err.Message)
}

func TestLoginUserInvalidUserJsonResponse(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		URL:          "http://localhost:8080/users/login",
		HTTPMethod:   http.MethodPost,
		ReqBody:      `{"email":"email@gmail.com","password":"password"}`,
		RespHTTPCode: http.StatusOK,
		RespBody:     `{"id": "1", "first_name": "Fede", "last_name": "Leon", "email": "email@gmail.com"}`,
	})
	repository := usersRepository{}

	user, err := repository.LoginUser("email@gmail.com", "password")

	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status)
	assert.EqualValues(t, "error when trying to unmarshall users response", err.Message)
}

func TestLoginUserNoError(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		URL:          "http://localhost:8080/users/login",
		HTTPMethod:   http.MethodPost,
		ReqBody:      `{"email":"email@gmail.com","password":"password"}`,
		RespHTTPCode: http.StatusOK,
		RespBody:     `{"id": 1, "first_name": "Fede", "last_name": "Leon", "email": "email@gmail.com"}`,
	})
	repository := usersRepository{}

	user, err := repository.LoginUser("email@gmail.com", "password")

	assert.NotNil(t, user)
	assert.Nil(t, err)
	assert.EqualValues(t, http.StatusOK, err.Status)
	assert.EqualValues(t, 1, user.Id)
	assert.EqualValues(t, "Fede", user.FirstName)
	assert.EqualValues(t, "Leon", user.LastName)
	assert.EqualValues(t, "email@gmail.com", user.Email)
}
