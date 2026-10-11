package substackbrowser

import (
	"strings"
	"testing"
)

func TestSetTitleAndSubtitleIncludesSpanishPlaceholders(t *testing.T) {
	js, err := SetTitleAndSubtitle("Título de prueba", "Subtítulo corto")
	if err != nil {
		t.Fatal(err)
	}
	for _, needle := range []string{
		`placeholder="Título"`,
		`placeholder="Añade un subtítulo..."`,
		"looksLikeSubtitleHint",
		"title_set",
	} {
		if !strings.Contains(js, needle) {
			t.Fatalf("title JS missing %q", needle)
		}
	}
}

func TestScheduleAfterContinueJSIncludesContinuar(t *testing.T) {
	js, err := ScheduleAfterContinueJS(ScheduleAfterContinueOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(js, "Continuar") {
		t.Fatal("schedule JS should try Continuar for Spanish publications")
	}
	if !strings.Contains(js, "añadir etiquetas") {
		t.Fatal("schedule JS should detect Spanish publish settings copy")
	}
}
