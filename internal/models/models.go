package models

import (
	"time"

	"github.com/google/uuid"
)

// Base model that includes common fields for all models
type Base struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// User represents a user in the system
type User struct {
	Base
	Username          string     `json:"username" db:"username"`
	Email             string     `json:"email" db:"email"`
	PasswordHash      string     `json:"-" db:"password_hash"`
	FirstName         *string    `json:"first_name,omitempty" db:"first_name"`
	LastName          *string    `json:"last_name,omitempty" db:"last_name"`
	AvatarURL         *string    `json:"avatar_url,omitempty" db:"avatar_url"`
	ProfileData       JSONB      `json:"profile_data,omitempty" db:"profile_data"`
	AgeGroup          string     `json:"age_group" db:"age_group"`
	Points            int        `json:"points" db:"points"`
	ExperienceLevel   int       `json:"experience_level" db:"experience_level"`
	TotalStudyTimeMin int        `json:"total_study_time_min" db:"total_study_time_min"`

	// Relationships
	Settings        []UserSetting         `json:"settings,omitempty"`
	Progress        []Progress            `json:"progress,omitempty"`
	ContentProgress []UserContentProgress `json:"content_progress,omitempty"`
	Notifications   []UserNotification    `json:"notifications,omitempty"`
}

// UserSetting contains user preferences and settings
type UserSetting struct {
	Base
	UserID                        uuid.UUID `json:"user_id" db:"user_id"`
	Theme                         string    `json:"theme" db:"theme"`
	NotificationDailyReminder     bool      `json:"notification_daily_reminder" db:"notification_daily_reminder"`
	NotificationAchievementAlerts bool      `json:"notification_achievement_alerts" db:"notification_achievement_alerts"`
	NotificationStudyTime         string    `json:"notification_study_time" db:"notification_study_time"`
	AccessibilityFontSize         string    `json:"accessibility_font_size" db:"accessibility_font_size"`
	AccessibilityHighContrast     bool      `json:"accessibility_high_contrast" db:"accessibility_high_contrast"`
	StudyDailyGoalMin             int       `json:"study_daily_goal_min" db:"study_daily_goal_min"`
	StudyPreferredSubject         string    `json:"study_preferred_subject" db:"study_preferred_subject"`
}

// UserNotification represents a notification sent to a user
type UserNotification struct {
	Base
	UserID    uuid.UUID  `json:"user_id" db:"user_id"`
	Type      string     `json:"type" db:"type"`
	Title     string     `json:"title" db:"title"`
	Message   string     `json:"message" db:"message"`
	IsRead    bool       `json:"is_read" db:"is_read"`
	RelatedID *uuid.UUID `json:"related_id,omitempty" db:"related_id"`
}

// Course represents a learning course
type Course struct {
	Base
	Topic       string `json:"topic" db:"topic"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	DurationMin int    `json:"duration_min" db:"duration_min"`

	// Relationships
	Lessons []Lesson `json:"lessons,omitempty"`
	Tags    []Tag    `json:"tags,omitempty" gorm:"many2many:course_tags;"`
}

// Tag represents a category or label for courses
type Tag struct {
	Base
	Name     string    `json:"name" db:"name"`
	Category string    `json:"category" db:"category"`
	Courses  []*Course `json:"courses,omitempty" gorm:"many2many:course_tags;"`
}

// CourseTag represents the many-to-many relationship between courses and tags
type CourseTag struct {
	CourseID uuid.UUID `json:"course_id" db:"course_id"`
	TagID    uuid.UUID `json:"tag_id" db:"tag_id"`
}

// Lesson represents a lesson within a course
type Lesson struct {
	Base
	CourseID         uuid.UUID `json:"course_id" db:"course_id"`
	Title            string    `json:"title" db:"title"`
	LessonOrder      int       `json:"lesson_order" db:"lesson_order"`
	EstimatedTimeMin int       `json:"estimated_time_min" db:"estimated_time_min"`
	PointsReward     int       `json:"points_reward" db:"points_reward"`

	// Relationships
	Course   *Course   `json:"course,omitempty"`
	Contents []Content `json:"contents,omitempty"`
}

// ContentType represents the type of content (theory, exercise, etc.)
type ContentType string

const (
	ContentTypeTheory      ContentType = "theory"
	ContentTypeExercise    ContentType = "exercise"
	ContentTypeAssessment  ContentType = "assessment"
	ContentTypeInteractive ContentType = "interactive"
	ContentTypeResource    ContentType = "resource"
)

// Content represents a piece of learning content
type Content struct {
	Base
	LessonID    uuid.UUID   `json:"lesson_id" db:"lesson_id"`
	Title       string      `json:"title" db:"title"`
	Description *string     `json:"description,omitempty" db:"description"`
	Order       int         `json:"order" db:"order"`
	ContentType ContentType `json:"content_type" db:"content_type"`

	// Relationships
	Lesson *Lesson `json:"-"`

	// Content type specific fields (only one should be set)
	Theory      *TheoryContent      `json:"theory,omitempty"`
	Exercise    *ExerciseContent    `json:"exercise,omitempty"`
	Assessment  *AssessmentContent  `json:"assessment,omitempty"`
	Interactive *InteractiveContent `json:"interactive,omitempty"`
	Resource    *ResourceContent    `json:"resource,omitempty"`
}

// TheoryContent contains theory content details
type TheoryContent struct {
	ID          uuid.UUID `json:"id" db:"id"`
	TextContent string    `json:"text_content" db:"text_content"`
	Examples    JSONB     `json:"examples,omitempty" db:"examples"`
	References  JSONB     `json:"references,omitempty" db:"references"`
}

// ExerciseContent contains exercise content details
type ExerciseContent struct {
	ID            uuid.UUID `json:"id" db:"id"`
	Problems      JSONB     `json:"problems" db:"problems"`
	EstimatedTime *int      `json:"estimated_time,omitempty" db:"estimated_time"`
}

// AssessmentContent contains assessment content details
type AssessmentContent struct {
	ID              uuid.UUID `json:"id" db:"id"`
	Questions       JSONB     `json:"questions" db:"questions"`
	TimeLimit       *int      `json:"time_limit,omitempty" db:"time_limit"`
	PassingScore    float64   `json:"passing_score" db:"passing_score"`
	AttemptsAllowed int       `json:"attempts_allowed" db:"attempts_allowed"`
}

// InteractiveContent contains interactive content details
type InteractiveContent struct {
	ID              uuid.UUID `json:"id" db:"id"`
	InteractiveType string    `json:"interactive_type" db:"interactive_type"`
	ContentData     JSONB     `json:"content_data" db:"content_data"`
	Config          JSONB     `json:"config,omitempty" db:"config"`
}

// ResourceContent contains resource content details
type ResourceContent struct {
	ID               uuid.UUID  `json:"id" db:"id"`
	ResourceType     string     `json:"resource_type" db:"resource_type"`
	URL              string     `json:"url" db:"url"`
	CreatedBy        *uuid.UUID `json:"created_by,omitempty" db:"created_by"`
	ResourceMetadata JSONB     `json:"resource_metadata,omitempty" db:"resource_metadata"`
}

// Progress tracks user progress in a course
type Progress struct {
	Base
	UserID             uuid.UUID `json:"user_id" db:"user_id"`
	CourseID           uuid.UUID `json:"course_id" db:"course_id"`
	CurrentLessonID    *uuid.UUID `json:"current_lesson_id,omitempty" db:"current_lesson_id"`
	TotalPointsEarned  int       `json:"total_points_earned" db:"total_points_earned"`
	TimeSpentMin       int       `json:"time_spent_min" db:"time_spent_min"`
	ProgressPercentage float64   `json:"progress_percentage" db:"progress_percentage"`
	ProgressData       JSONB     `json:"progress_data" db:"progress_data"`
	LastAccessed       time.Time `json:"last_accessed" db:"last_accessed"`
	IsCompleted        bool      `json:"is_completed" db:"is_completed"`

	// Relationships
	User          *User     `json:"user,omitempty"`
	Course        *Course   `json:"course,omitempty"`
	CurrentLesson *Lesson   `json:"current_lesson,omitempty"`
	ContentStates []ContentState `json:"content_states,omitempty"`
}

// UserContentProgress tracks user progress for a specific content item
type UserContentProgress struct {
	Base
	UserID          uuid.UUID `json:"user_id" db:"user_id"`
	ContentID       uuid.UUID `json:"content_id" db:"content_id"`
	IsCompleted     bool      `json:"is_completed" db:"is_completed"`
	Score           *float64  `json:"score,omitempty" db:"score"`
	Attempts        int       `json:"attempts" db:"attempts"`
	TimeSpentMin    int       `json:"time_spent_min" db:"time_spent_min"`
	LastInteraction time.Time `json:"last_interaction" db:"last_interaction"`

	// Relationships
	User    *User    `json:"user,omitempty"`
	Content *Content `json:"content,omitempty"`
}

// ContentState stores detailed state information for resuming content
type ContentState struct {
	Base
	UserID       uuid.UUID `json:"user_id" db:"user_id"`
	ProgressID   uuid.UUID `json:"progress_id" db:"progress_id"`
	ContentID    uuid.UUID `json:"content_id" db:"content_id"`
	StateType    string    `json:"state_type" db:"state_type"`
	NumericValue *float64 `json:"numeric_value,omitempty" db:"numeric_value"`
	JSONValue    JSONB     `json:"json_value,omitempty" db:"json_value"`
	TextValue    *string   `json:"text_value,omitempty" db:"text_value"`

	// Relationships
	User     *User    `json:"user,omitempty"`
	Progress *Progress `json:"progress,omitempty"`
	Content  *Content  `json:"content,omitempty"`
}

// CompletedLesson records when a user completes a lesson
type CompletedLesson struct {
	Base
	UserID       uuid.UUID `json:"user_id" db:"user_id"`
	LessonID     uuid.UUID `json:"lesson_id" db:"lesson_id"`
	CourseID     uuid.UUID `json:"course_id" db:"course_id"`
	Score        *float64  `json:"score,omitempty" db:"score"`
	TimeSpentMin int       `json:"time_spent_min" db:"time_spent_min"`

	// Relationships
	User   *User   `json:"user,omitempty"`
	Lesson *Lesson `json:"lesson,omitempty"`
	Course *Course `json:"course,omitempty"`
}

// CompletedCourse records when a user completes a course
type CompletedCourse struct {
	Base
	UserID                uuid.UUID `json:"user_id" db:"user_id"`
	CourseID              uuid.UUID `json:"course_id" db:"course_id"`
	FinalScore            *float64  `json:"final_score,omitempty" db:"final_score"`
	TotalTimeSpentMin     int       `json:"total_time_spent_min" db:"total_time_spent_min"`
	CompletedLessonsCount int       `json:"completed_lessons_count" db:"completed_lessons_count"`
	AchievementsEarned    []string  `json:"achievements_earned,omitempty" db:"achievements_earned"`
	CertificateID         *string   `json:"certificate_id,omitempty" db:"certificate_id"`

	// Relationships
	User   *User   `json:"user,omitempty"`
	Course *Course `json:"course,omitempty"`
}

// UserAnswer records a user's answer to a question
type UserAnswer struct {
	Base
	UserID       uuid.UUID `json:"user_id" db:"user_id"`
	ContentID    uuid.UUID `json:"content_id" db:"content_id"`
	QuestionID   string    `json:"question_id" db:"question_id"`
	AnswerData   JSONB     `json:"answer_data" db:"answer_data"`
	IsCorrect    bool      `json:"is_correct" db:"is_correct"`
	PointsEarned int       `json:"points_earned" db:"points_earned"`

	// Relationships
	User    *User    `json:"user,omitempty"`
	Content *Content `json:"content,omitempty"`
}

// JSONB is a wrapper around map[string]interface{} for JSONB database fields
type JSONB map[string]any
