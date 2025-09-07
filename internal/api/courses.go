package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"

	"mathtermind-go/internal/db"
	apperrors "mathtermind-go/internal/errors"
)

// ListCoursesHandler handles GET /api/v1/courses
func ListCoursesHandler(pool *pgxpool.Pool) apperrors.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		limit := 20
		offset := 0
		if v := r.URL.Query().Get("limit"); v != "" {
			lv, err := strconv.Atoi(v)
			if err != nil || lv <= 0 || lv > 100 {
				return apperrors.Errorf(apperrors.ErrCodeValidation, "invalid limit parameter").WithDetails(map[string]any{"limit": v})
			}
			limit = lv
		}
		if v := r.URL.Query().Get("offset"); v != "" {
			ov, err := strconv.Atoi(v)
			if err != nil || ov < 0 {
				return apperrors.Errorf(apperrors.ErrCodeValidation, "invalid offset parameter").WithDetails(map[string]any{"offset": v})
			}
			offset = ov
		}

		courses, err := db.ListCourses(r.Context(), pool, limit, offset)
		if err != nil {
			return apperrors.Wrap(err, apperrors.ErrCodeDBQuery, "failed to list courses")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		return json.NewEncoder(w).Encode(map[string]any{
			"items":  courses,
			"limit":  limit,
			"offset": offset,
		})
	}
}
