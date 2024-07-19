type User struct {
    ID        uint           `json:"id"`
    Username  string         `json:"username"`
    Email     string         `json:"email"`
    Password  string         `json:"password"`
    FirstName string         `json:"first_name"`
    LastName  string         `json:"last_name"`
    RoleID    uint            `json:"role_id"` // Make sure this tag matches the JSON key
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `json:"deleted_at"`
}
