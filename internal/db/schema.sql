-- Database schema for Mathtermind
-- This file contains all the table definitions and relationships

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    avatar_url VARCHAR(255),
    profile_data JSONB,
    age_group VARCHAR(50) NOT NULL,
    points INTEGER NOT NULL DEFAULT 0,
    experience_level INTEGER NOT NULL DEFAULT 1,
    total_study_time_min INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- User settings
CREATE TABLE user_settings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    theme VARCHAR(20) NOT NULL DEFAULT 'light',
    notification_daily_reminder BOOLEAN NOT NULL DEFAULT true,
    notification_achievement_alerts BOOLEAN NOT NULL DEFAULT true,
    notification_study_time VARCHAR(5) NOT NULL DEFAULT '09:00',
    accessibility_font_size VARCHAR(10) NOT NULL DEFAULT 'medium',
    accessibility_high_contrast BOOLEAN NOT NULL DEFAULT false,
    study_daily_goal_min INTEGER NOT NULL DEFAULT 30,
    study_preferred_subject VARCHAR(50) NOT NULL DEFAULT 'MATH',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Application settings
CREATE TABLE settings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    key VARCHAR(255) NOT NULL UNIQUE,
    value TEXT NOT NULL,
    description TEXT,
    is_protected BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- User notifications
CREATE TABLE user_notifications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL,
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    is_read BOOLEAN NOT NULL DEFAULT false,
    related_id UUID,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Tags for content categorization
CREATE TABLE tags (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL UNIQUE,
    category VARCHAR(50) NOT NULL DEFAULT 'TOPIC',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Courses
CREATE TABLE courses (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    topic VARCHAR(50) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    duration_min INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Course tags (many-to-many relationship)
CREATE TABLE course_tags (
    course_id UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    tag_id UUID NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY (course_id, tag_id)
);

-- Lessons
CREATE TABLE lessons (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    course_id UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    lesson_order INTEGER NOT NULL,
    estimated_time_min INTEGER NOT NULL,
    points_reward INTEGER NOT NULL DEFAULT 10,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Base content table (polymorphic)
CREATE TABLE content (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    lesson_id UUID NOT NULL REFERENCES lessons(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    "order" INTEGER NOT NULL,
    content_type VARCHAR(20) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Theory content
CREATE TABLE theory_content (
    id UUID PRIMARY KEY REFERENCES content(id) ON DELETE CASCADE,
    text_content TEXT NOT NULL,
    examples JSONB,
    references JSONB
);

-- Exercise content
CREATE TABLE exercise_content (
    id UUID PRIMARY KEY REFERENCES content(id) ON DELETE CASCADE,
    problems JSONB NOT NULL,
    estimated_time_min INTEGER
);

-- Assessment content
CREATE TABLE assessment_content (
    id UUID PRIMARY KEY REFERENCES content(id) ON DELETE CASCADE,
    questions JSONB NOT NULL,
    time_limit_min INTEGER,
    passing_score FLOAT NOT NULL DEFAULT 70.0,
    attempts_allowed INTEGER NOT NULL DEFAULT 3
);

-- Interactive content
CREATE TABLE interactive_content (
    id UUID PRIMARY KEY REFERENCES content(id) ON DELETE CASCADE,
    interactive_type VARCHAR(50) NOT NULL,
    content_data JSONB NOT NULL,
    config JSONB
);

-- Resource content
CREATE TABLE resource_content (
    id UUID PRIMARY KEY REFERENCES content(id) ON DELETE CASCADE,
    resource_type VARCHAR(50) NOT NULL,
    url VARCHAR(1024) NOT NULL,
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    resource_metadata JSONB
);

-- User progress tracking
CREATE TABLE progress (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    course_id UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    current_lesson_id UUID REFERENCES lessons(id) ON DELETE CASCADE,
    total_points_earned INTEGER NOT NULL DEFAULT 0,
    time_spent_min INTEGER NOT NULL DEFAULT 0,
    progress_percentage FLOAT NOT NULL DEFAULT 0.0,
    progress_data JSONB NOT NULL DEFAULT '{}'::jsonb,
    last_accessed TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    is_completed BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- User content progress
CREATE TABLE user_content_progress (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content_id UUID NOT NULL REFERENCES content(id) ON DELETE CASCADE,
    is_completed BOOLEAN NOT NULL DEFAULT false,
    score FLOAT,
    attempts INTEGER NOT NULL DEFAULT 0,
    time_spent INTEGER NOT NULL DEFAULT 0,
    last_interaction TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, content_id)
);

-- Content state for resuming
CREATE TABLE content_states (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    progress_id UUID NOT NULL REFERENCES progress(id) ON DELETE CASCADE,
    content_id UUID NOT NULL REFERENCES content(id) ON DELETE CASCADE,
    state_type VARCHAR(50) NOT NULL,
    numeric_value FLOAT,
    json_value JSONB,
    text_value TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, content_id, state_type)
);

-- Completed lessons
CREATE TABLE completed_lessons (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    lesson_id UUID NOT NULL REFERENCES lessons(id) ON DELETE CASCADE,
    course_id UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    completed_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    score FLOAT,
    time_spent INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Completed courses
CREATE TABLE completed_courses (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    course_id UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    completed_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    final_score FLOAT,
    total_time_spent INTEGER NOT NULL,
    completed_lessons_count INTEGER NOT NULL,
    achievements_earned UUID[],
    certificate_id UUID,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, course_id)
);

-- User answers
CREATE TABLE user_answers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content_id UUID NOT NULL REFERENCES content(id) ON DELETE CASCADE,
    question_id VARCHAR(100) NOT NULL,
    answer_data JSONB NOT NULL,
    is_correct BOOLEAN NOT NULL,
    points_earned INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for better query performance
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_user_settings_user_id ON user_settings(user_id);
CREATE INDEX idx_user_notifications_user_id ON user_notifications(user_id);
CREATE INDEX idx_user_notifications_type ON user_notifications(type);
CREATE INDEX idx_courses_topic ON courses(topic);
CREATE INDEX idx_lessons_course_id ON lessons(course_id);
CREATE INDEX idx_content_lesson_id ON content(lesson_id);
CREATE INDEX idx_content_content_type ON content(content_type);
CREATE INDEX idx_progress_user_course ON progress(user_id, course_id);
CREATE INDEX idx_user_content_progress_user_content ON user_content_progress(user_id, content_id);
CREATE INDEX idx_content_states_user_content ON content_states(user_id, content_id);
CREATE INDEX idx_completed_lessons_user_course ON completed_lessons(user_id, course_id);
CREATE INDEX idx_completed_courses_user ON completed_courses(user_id);
CREATE INDEX idx_user_answers_user_content ON user_answers(user_id, content_id);
