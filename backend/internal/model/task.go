package model

import "time"

type TaskStatus string

const (
	StatusTodo          TaskStatus = "todo"
	StatusInProgress    TaskStatus = "in_progress"
	StatusStagingReview TaskStatus = "staging_review"
	StatusDone          TaskStatus = "done"
)

type Priority string

const (
	PriorityP1 Priority = "p1"
	PriorityP2 Priority = "p2"
	PriorityP3 Priority = "p3"
	PriorityP4 Priority = "p4"
)

type MemberRole string

const (
	RoleOwner  MemberRole = "owner"
	RoleAdmin  MemberRole = "admin"
	RoleMember MemberRole = "member"
)

type NotificationType string

const (
	NotifReminder   NotificationType = "reminder"
	NotifAlert      NotificationType = "alert"
	NotifMention    NotificationType = "mention"
	NotifAssignment NotificationType = "assignment"
)

type ActivityAction string

const (
	ActionTaskCreated    ActivityAction = "task_created"
	ActionTaskUpdated    ActivityAction = "task_updated"
	ActionTaskDeleted    ActivityAction = "task_deleted"
	ActionTaskAssigned   ActivityAction = "task_assigned"
	ActionStatusChanged  ActivityAction = "status_changed"
	ActionCommentAdded   ActivityAction = "comment_added"
	ActionCommentDeleted ActivityAction = "comment_deleted"
	ActionMemberAdded    ActivityAction = "member_added"
	ActionMemberRemoved  ActivityAction = "member_removed"
	ActionSprintCreated  ActivityAction = "sprint_created"
	ActionSprintUpdated  ActivityAction = "sprint_updated"
	ActionLabelCreated   ActivityAction = "label_created"
	ActionLabelDeleted   ActivityAction = "label_deleted"
)

type Password struct {
	PlainText *string
	Hash      []byte
}

type User struct {
	ID           string    `json:"id" bson:"_id"`
	Email        string    `json:"email" bson:"email"`
	Name         string    `json:"name" bson:"name"`
	AvatarURL    string    `json:"avatar_url,omitempty" bson:"avatar_url,omitempty"`
	PasswordHash Password  `json:"-" bson:"password_hash"`
	CreatedAt    time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" bson:"updated_at"`
}

var AnonymousUser = &User{}

type Token struct {
	PlainText string    `json:"token" bson:"token"`
	Hash      []byte    `json:"-" bson:"hash"`
	UserID    string    `json:"-" bson:"user_id"`
	Expiry    time.Time `json:"expiry" bson:"expiry"`
	Scope     string    `json:"-" bson:"scope"`
}

type Project struct {
	ID          string    `json:"id" bson:"_id"`
	Name        string    `json:"name" bson:"name"`
	Description string    `json:"description,omitempty" bson:"description,omitempty"`
	OwnerID     string    `json:"owner_id" bson:"owner_id"`
	MemberCount int       `json:"member_count" bson:"member_count"`
	TaskCount   int       `json:"task_count" bson:"task_count"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" bson:"updated_at"`
}

type ProjectMember struct {
	UserID    string     `json:"user_id" bson:"user_id"`
	Name      string     `json:"name" bson:"name"`
	Email     string     `json:"email" bson:"email"`
	AvatarURL string     `json:"avatar_url,omitempty" bson:"avatar_url,omitempty"`
	Role      MemberRole `json:"role" bson:"role"`
	JoinedAt  time.Time  `json:"joined_at" bson:"joined_at"`
}

type Task struct {
	ID           string        `json:"id" bson:"_id"`
	ProjectID    string        `json:"project_id" bson:"project_id"`
	Title        string        `json:"title" bson:"title"`
	Description  string        `json:"description,omitempty" bson:"description,omitempty"`
	Status       TaskStatus    `json:"status" bson:"status"`
	Priority     Priority      `json:"priority" bson:"priority"`
	DueDate      *time.Time    `json:"due_date,omitempty" bson:"due_date,omitempty"`
	IsPastDue    bool          `json:"is_past_due" bson:"-"`
	AssigneeID   *string       `json:"assignee_id,omitempty" bson:"assignee_id,omitempty"`
	Assignee     *User         `json:"assignee,omitempty" bson:"-"`
	ReporterID   string        `json:"reporter_id" bson:"reporter_id"`
	Reporter     *User         `json:"reporter,omitempty" bson:"-"`
	LabelIDs     []string      `json:"label_ids,omitempty" bson:"label_ids,omitempty"`
	Labels       []Label       `json:"labels,omitempty" bson:"-"`
	SprintID     *string       `json:"sprint_id,omitempty" bson:"sprint_id,omitempty"`
	CommentCount int           `json:"comment_count" bson:"comment_count"`
	TimeTracking *TimeTracking `json:"time_tracking,omitempty" bson:"time_tracking,omitempty"`
	CreatedAt    time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at" bson:"updated_at"`
}

type TimeTracking struct {
	EstimatedHours  float64 `json:"estimated_hours" bson:"estimated_hours"`
	LoggedHours     float64 `json:"logged_hours" bson:"logged_hours"`
	IsOverdue       bool    `json:"is_overdue" bson:"-"`
	OverdueDuration *string `json:"overdue_duration,omitempty" bson:"-"`
}

// stored in a separate mongo collection
type Comment struct {
	ID        string    `json:"id" bson:"_id"`
	TaskID    string    `json:"task_id" bson:"task_id"`
	UserID    string    `json:"user_id" bson:"user_id"`
	User      *User     `json:"user,omitempty" bson:"user"`
	Content   string    `json:"content" bson:"content"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

type Label struct {
	ID        string    `json:"id" bson:"_id"`
	ProjectID string    `json:"project_id" bson:"project_id"`
	Name      string    `json:"name" bson:"name"`
	Color     string    `json:"color,omitempty" bson:"color,omitempty"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}

type ActivityEntry struct {
	ID        string                 `json:"id" bson:"_id"`
	ProjectID string                 `json:"project_id" bson:"project_id"`
	TaskID    *string                `json:"task_id,omitempty" bson:"task_id,omitempty"`
	UserID    string                 `json:"user_id" bson:"user_id"`
	User      *User                  `json:"user,omitempty" bson:"-"`
	Action    ActivityAction         `json:"action" bson:"action"`
	Details   map[string]interface{} `json:"details,omitempty" bson:"details,omitempty"`
	CreatedAt time.Time              `json:"created_at" bson:"created_at"`
}

type Sprint struct {
	ID        string    `json:"id" bson:"_id"`
	ProjectID string    `json:"project_id" bson:"project_id"`
	Name      string    `json:"name" bson:"name"`
	Label     string    `json:"label,omitempty" bson:"label,omitempty"`
	StartDate time.Time `json:"start_date" bson:"start_date"`
	EndDate   time.Time `json:"end_date" bson:"end_date"`
	IsActive  bool      `json:"is_active" bson:"is_active"`
	TaskCount int       `json:"task_count" bson:"task_count"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

type Notification struct {
	ID            string           `json:"id" bson:"_id"`
	UserID        string           `json:"user_id" bson:"user_id"`
	Type          NotificationType `json:"type" bson:"type"`
	Title         string           `json:"title" bson:"title"`
	Message       string           `json:"message" bson:"message"`
	IsRead        bool             `json:"is_read" bson:"is_read"`
	ReferenceType *string          `json:"reference_type,omitempty" bson:"reference_type,omitempty"`
	ReferenceID   *string          `json:"reference_id,omitempty" bson:"reference_id,omitempty"`
	CreatedAt     time.Time        `json:"created_at" bson:"created_at"`
}
