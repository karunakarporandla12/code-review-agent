package app

import "testing"

func TestTechnologiesIncludeMoreSupportedStacks(t *testing.T) {
	if !validFrontendTechnology(normalizeGeneratedTechnology("vue js")) {
		t.Fatal("expected vue to be treated as a supported frontend technology")
	}

	if !validFrontendTechnology(normalizeGeneratedTechnology("next.js")) {
		t.Fatal("expected nextjs to be treated as a supported frontend technology")
	}

	if !validBackendTechnology(normalizeGeneratedTechnology("nodejs")) {
		t.Fatal("expected nodejs to be treated as a supported backend technology")
	}

	if !validBackendTechnology(normalizeGeneratedTechnology("python")) {
		t.Fatal("expected python to be treated as a supported backend technology")
	}

	if !validDatabaseTechnology(normalizeGeneratedTechnology("sqlite")) {
		t.Fatal("expected sqlite to be treated as a supported database technology")
	}
}

func TestPromptTechnologyNormalizationMatchesCommonAliases(t *testing.T) {
	cases := map[string]string{
		"react js":        "react",
		"ReactJS":         "react",
		"Next.js":         "nextjs",
		"Spring Boot":     "java",
		"Go Lang":         "golang",
		"node express":    "node",
		"Postgres SQL":    "postgresql",
		"mongo db":        "mongodb",
		"sqlite database": "sqlite",
	}

	for input, want := range cases {
		if got := normalizeGeneratedTechnology(input); got != want {
			t.Fatalf("normalizeGeneratedTechnology(%q) = %q, want %q", input, got, want)
		}
	}
}
