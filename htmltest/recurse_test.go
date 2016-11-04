package htmltest

import "testing"

func TestRecurseDirectory(t *testing.T) {
	options := map[string]interface{}{
		"DirectoryPath": "fixtures/links/rootLink",
	}
	SetOptions(options)
	documents := RecurseDirectory("")

	// There should be three documents in this directory, did we find them all?
	t_assertEqual(t, len(documents), 3)
}
