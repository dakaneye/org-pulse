package github

import "testing"

func TestFilterRepos(t *testing.T) {
	repos := []Repo{
		{Name: "active", Owner: "org", IsArchived: false, IsFork: false},
		{Name: "archived", Owner: "org", IsArchived: true, IsFork: false},
		{Name: "fork", Owner: "org", IsArchived: false, IsFork: true},
		{Name: "other", Owner: "org", IsArchived: false, IsFork: false},
	}

	t.Run("no filters", func(t *testing.T) {
		got := FilterRepos(repos, nil, nil)
		if len(got) != 2 {
			t.Fatalf("got %d repos, want 2 (active + other)", len(got))
		}
	})

	t.Run("include filter", func(t *testing.T) {
		got := FilterRepos(repos, []string{"active"}, nil)
		if len(got) != 1 || got[0].Name != "active" {
			t.Errorf("got %v, want [active]", got)
		}
	})

	t.Run("exclude filter", func(t *testing.T) {
		got := FilterRepos(repos, nil, []string{"other"})
		if len(got) != 1 || got[0].Name != "active" {
			t.Errorf("got %v, want [active]", got)
		}
	})

	t.Run("include archived is still excluded", func(t *testing.T) {
		got := FilterRepos(repos, []string{"archived"}, nil)
		if len(got) != 0 {
			t.Errorf("got %v, want empty (archived repos excluded)", got)
		}
	})

	t.Run("empty input", func(t *testing.T) {
		got := FilterRepos(nil, nil, nil)
		if len(got) != 0 {
			t.Errorf("got %d, want 0", len(got))
		}
	})
}
