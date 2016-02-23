package gibson

import (
	"net/http/httptest"
	"testing"
)

func TestEmpty(t *testing.T) {
	app := New()
	ts := httptest.NewServer(app)
	defer ts.Close()

}
