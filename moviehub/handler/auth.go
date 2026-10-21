import (
    "encoding/json"
    "errors"
    "net/http"
    "strings"
    "unicode/utf8"

    "moviehub/auth"
    "moviehub/mw"
    "moviehub/store"
)

// ... Register และ Login ของเดิมอยู่ตรงนี้ ไม่แตะ ...

// Me (ด่านที่ 4) ตอบข้อมูลของคนที่ถือตั๋ว เลขสมาชิกมาจากตั๋วที่ RequireLogin ตรวจแล้วเท่านั้น
func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
    id, ok := mw.MemberID(r.Context())
    if !ok {
        writeError(w, http.StatusUnauthorized, "NO_TOKEN", "กรุณาเข้าสู่ระบบก่อน")
        return
    }
    m, err := h.Store.MemberByID(r.Context(), id)
    if errors.Is(err, store.ErrNotFound) {
        writeError(w, http.StatusUnauthorized, "MEMBER_GONE", "ไม่พบบัญชีนี้แล้ว กรุณาเข้าสู่ระบบใหม่")
        return
    }
    if err != nil {
        internal(w, err)
        return
    }
    writeJSON(w, http.StatusOK, m)
}