package user

import (
	"context"
	"io"
	"testing"

	user "github.com/Megidy/e-commerce/frontend/templates/user"
	"github.com/PuerkitoBio/goquery"
)

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
