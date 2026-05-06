package runstate

import "testing"

func TestNextForStatusPreservesCommandReasonAndAddsGateMetadata(t *testing.T) {
	next := NextForStatus(StatusIntakeComplete, "run-001")
	if next == nil {
		t.Fatal("expected next action")
	}
	if next.Command == "" || next.Reason == "" {
		t.Fatalf("command and reason must remain populated: %+v", next)
	}
	if next.Stage != "refine" || next.GateKind != "semantic" {
		t.Fatalf("unexpected gate metadata: %+v", next)
	}
	if next.ContextCommand != "specops context run-001" {
		t.Fatalf("context command = %q", next.ContextCommand)
	}
	if next.NoteCommand != "specops note run-001 --stage refine --text <file-or-inline>" {
		t.Fatalf("note command = %q", next.NoteCommand)
	}
	if !next.HumanInputRecommended || len(next.SuggestedQuestions) != 3 {
		t.Fatalf("unexpected operator guidance: %+v", next)
	}
}
