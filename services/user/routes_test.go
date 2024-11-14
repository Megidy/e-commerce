package user

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	user "github.com/Megidy/e-commerce/frontend/templates/user"
	"github.com/Megidy/e-commerce/services/auth"
	"github.com/Megidy/e-commerce/types"
	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
)

type mockTemplate struct {
	Case string
}

func (m *mockTemplate) GetDataFromForm(c *gin.Context, key string) string {
	switch m.Case {
	case "Login:Pass:CorrectPayload":
		if key == "email" {
			return "PastTestEmail@gmail.com"
		}
		if key == "password" {
			return "PastTestPassword"
		}
	case "Login:Fail:NoEmailFound":
		if key == "email" {
			return "Fail:NoEmailFound,TestEmail@gmail.com"
		}
		if key == "password" {
			return "FailTestPassword"
		}
	case "Login:Fail:InvalidPassword":
		if key == "email" {
			return "Fail:IncorrectPassword,TestEmail@gmail.com"
		}
		if key == "password" {
			return "FailTestPassword"
		}
	}

	return ""
}

type mockUserStore struct {
}

func (m *mockUserStore) CreateUser(user *types.User) error {
	return nil
}
func (m *mockUserStore) AlreadyExists(user *types.User) (bool, error) {
	if user.Email == "Fail:NoEmailFound,TestEmail@gmail.com" {
		return false, nil
	}

	return true, nil
}
func (m *mockUserStore) GetUserByEmail(email string) (types.User, error) {
	if email == "PastTestEmail@gmail.com" {
		pass, _ := auth.HashPassword("PastTestPassword")
		return types.User{
			ID:       "test",
			Name:     "testName",
			LastName: "testLastName",
			Email:    email,
			Password: string(pass),
			Created:  "11.11.1111",
			Role:     "user",
		}, nil

	} else if email == "Fail:IncorrectPassword,TestEmail@gmail.com" {
		pass, _ := auth.HashPassword("another password found in db")
		return types.User{
			ID:       "test",
			Name:     "testName",
			LastName: "testLastName",
			Email:    email,
			Password: string(pass),
			Created:  "11.11.1111",
			Role:     "user",
		}, nil

	}
	return types.User{}, nil
}

func (m *mockUserStore) GetUserById(id string) (types.User, error) {
	return types.User{}, nil
}

func TestLoadLogInTemplate(t *testing.T) {
	r, w := io.Pipe()
	expectedLabelName := "Email"
	go func() {
		_ = user.Login(false, "").Render(context.Background(), w)
		_ = w.Close()
	}()
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		t.Fatalf("error while creating document : %v", err)
	}
	if actualLabelName := doc.Find("label").First().Text(); actualLabelName != expectedLabelName {
		t.Errorf("expected label name %q, got %q ", expectedLabelName, actualLabelName)
	}
}

func TestLogIn(t *testing.T) {

	userStore := &mockUserStore{}

	t.Run("Should Pass if payload is correct", func(t *testing.T) {
		templates := &mockTemplate{Case: "Login:Pass:CorrectPayload"}
		handler := NewHandler(templates, userStore)
		req, err := http.NewRequest(http.MethodPost, "/login/enter", nil)
		if err != nil {
			t.Errorf("error while creating new request : %v", err)
		}
		rr := httptest.NewRecorder()
		router := gin.Default()
		router.POST("/login/enter", func(c *gin.Context) {
			handler.LogIn(c)

		})
		router.ServeHTTP(rr, req)

	})
	t.Run("Should Fail if email is not found in DB ", func(t *testing.T) {
		template := &mockTemplate{Case: "Login:Fail:NoEmailFound"}
		handler := NewHandler(template, userStore)
		req, err := http.NewRequest(http.MethodPost, "/login/enter", nil)
		if err != nil {
			t.Errorf("error while creating new request : %v", err)
		}
		rr := httptest.NewRecorder()
		router := gin.Default()
		router.POST("/login/enter", func(c *gin.Context) {
			handler.LogIn(c)
			if c.Writer.Header().Get("hasEmail") != "false" {
				t.Errorf("expected email to be not found")
			}

		})
		router.ServeHTTP(rr, req)
	})
	t.Run("Should fail if password of payload doesn`t match with native password", func(t *testing.T) {
		template := &mockTemplate{Case: "Login:Fail:InvalidPassword"}
		handler := NewHandler(template, userStore)
		req, err := http.NewRequest(http.MethodPost, "/login/enter", nil)
		if err != nil {
			t.Errorf("error while creating new request :%v", err)
		}
		rr := httptest.NewRecorder()
		router := gin.Default()
		router.POST("/login/enter", func(c *gin.Context) {
			handler.LogIn(c)
			if c.Writer.Header().Get("correctPassword") != "false" {
				t.Errorf("expected password to be incorrect")
			}
		})
		router.ServeHTTP(rr, req)
	})
}