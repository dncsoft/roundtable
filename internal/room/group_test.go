package room

import "testing"

func roomSummary(t *testing.T, h *Hub, id string) RoomSummary {
	t.Helper()
	for _, s := range h.Rooms() {
		if s.ID == id {
			return s
		}
	}
	t.Fatalf("방 요약 없음: %s", id)
	return RoomSummary{}
}

func TestRoomGroupPatchAndSummary(t *testing.T) {
	h, r := newTestHub(t, &fakeProvider{id: "codex"})
	if s := roomSummary(t, h, r.ID); s.Group != "" || s.Workspace != r.Workspace {
		t.Fatalf("초기 요약: %+v", s)
	}
	if _, err := h.UpdateRoom(r.ID, RoomPatch{Group: ref("  설계 검토 ")}); err != nil {
		t.Fatal(err)
	}
	if s := roomSummary(t, h, r.ID); s.Group != "설계 검토" {
		t.Fatalf("그룹 지정: %q", s.Group)
	}
	// Other patches leave the group alone.
	if _, err := h.UpdateRoom(r.ID, RoomPatch{Title: ref("새 이름")}); err != nil {
		t.Fatal(err)
	}
	if s := roomSummary(t, h, r.ID); s.Group != "설계 검토" {
		t.Fatalf("다른 변경 뒤 그룹: %q", s.Group)
	}
	// A blank group returns the room to its folder group.
	if _, err := h.UpdateRoom(r.ID, RoomPatch{Group: ref("   ")}); err != nil {
		t.Fatal(err)
	}
	if s := roomSummary(t, h, r.ID); s.Group != "" {
		t.Fatalf("그룹 해제: %q", s.Group)
	}
}
