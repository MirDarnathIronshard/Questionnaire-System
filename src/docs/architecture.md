# Detailed System Architecture Explanation

## 1. System Overview
The voting system is designed as a multi-user platform that allows for creating and managing questionnaires. The key features include:

- **Questionnaire Management**: Users can create, edit, and manage surveys with different question types
- **Access Control**: Different permission levels for different user roles
- **Response Collection**: Secure collection and storage of survey responses
- **Analytics**: Real-time analysis and visualization of survey results

## 2. Core Modules Explained

### 2.1 Authentication Module
This module handles all user identity and access management:

- **User Registration**:
    - Collects email and national ID
    - Validates national ID format using Iranian national ID rules
    - Sends verification email
    - Creates user account after verification

- **JWT Authentication**:
    - Generates JWT tokens upon login
    - Includes user role and permissions in token
    - Manages token expiration and refresh
    - Validates tokens on protected routes

- **RBAC (Role-Based Access Control)**:
    - Uses Casbin for role management
    - Defines roles like admin, questionnaire owner, participant
    - Maps permissions to roles
    - Enforces access rules

### 2.2 Questionnaire Management Module
Handles the core functionality of survey creation and management:

- **Question Types Support**:
    - Multiple choice
    - Single choice
    - Text response
    - Numerical response
    - File upload

- **Ordering Options**:
    - Sequential: Questions appear in fixed order
    - Random: Questions are shuffled for each participant
    - Conditional: Questions appear based on previous answers

- **Time Management**:
    - Schedule start and end dates
    - Set time limits for completion
    - Allow/restrict response editing
    - Track response timestamps

### 2.3 Access Control Module
Manages permissions at both system and questionnaire levels:

- **System Level**:
    - Global roles (admin, user, etc.)
    - System-wide permissions
    - User group management

- **Questionnaire Level**:
    - Owner permissions
    - Participant permissions
    - Viewer permissions
    - Response anonymity settings

### 2.4 Analytics Module
Provides comprehensive analysis of questionnaire responses:

- **Real-time Monitoring**:
    - Live response tracking
    - Participation rates
    - Completion rates
    - Response distribution

- **Statistical Analysis**:
    - Response summaries
    - Cross-tabulation
    - Trend analysis
    - Custom calculations

- **Export Options**:
    - CSV export
    - Excel reports
    - PDF summaries
    - Raw data access

### 2.5 Notification Module
Handles all system communications:

- **Email Notifications**:
    - Registration verification
    - Password reset
    - Questionnaire invitations
    - Response confirmations

- **System Notifications**:
    - Real-time updates
    - Status changes
    - Deadline reminders
    - System alerts

## 3. Technical Implementation

### Database Structure
```sql
-- Key tables and relationships
Users (
    id SERIAL PRIMARY KEY,
    email VARCHAR UNIQUE,
    national_id VARCHAR UNIQUE,
    password_hash VARCHAR,
    role VARCHAR,
    created_at TIMESTAMP
)

Questionnaires (
    id SERIAL PRIMARY KEY,
    title VARCHAR,
    description TEXT,
    owner_id INTEGER REFERENCES Users(id),
    start_time TIMESTAMP,
    end_time TIMESTAMP,
    status VARCHAR,
    settings JSONB
)

Questions (
    id SERIAL PRIMARY KEY,
    questionnaire_id INTEGER REFERENCES Questionnaires(id),
    type VARCHAR,
    content TEXT,
    settings JSONB,
    order INTEGER
)

Responses (
    id SERIAL PRIMARY KEY,
    question_id INTEGER REFERENCES Questions(id),
    user_id INTEGER REFERENCES Users(id),
    content TEXT,
    created_at TIMESTAMP
)
```

### API Structure
```plaintext
/api
  /auth
    POST /register
    POST /login
    POST /verify-email
    POST /reset-password
  
  /questionnaires
    GET /
    POST /
    GET /{id}
    PUT /{id}
    DELETE /{id}
    GET /{id}/responses
    
  /responses
    POST /
    GET /{id}
    PUT /{id}
    
  /analytics
    GET /questionnaire/{id}/summary
    GET /questionnaire/{id}/export
    GET /questionnaire/{id}/live
```

### Security Implementation
1. **Authentication Flow**:
   ```plaintext
   1. User submits credentials
   2. System validates credentials
   3. Generates JWT token with:
      - User ID
      - Role
      - Permissions
      - Expiration time
   4. Returns token to client
   5. Client includes token in Authorization header
   ```

2. **Permission Validation Flow**:
   ```plaintext
   1. Request arrives with JWT token
   2. Middleware extracts user info
   3. Casbin checks permission rules
   4. If allowed: proceed
   5. If denied: return 403
   ```

### Real-time Features Implementation
Using WebSocket for:
- Live response tracking
- Real-time analytics updates
- System notifications
- Chat functionality

```javascript
// WebSocket events
{
  "type" - "RESPONSE_RECEIVED",
  "data" - {
    "questionnaireId": "123",
    "responseCount": 45,
    "completionRate": 78.5
  }
}
```

### Caching Strategy
Using Redis for:
- Session management
- Response rate limiting
- Analytics caching
- Temporary data storage

```plaintext
Cache Keys:
- user:{id}:session
- questionnaire:{id}:stats
- user:{id}:responses
- system:active_questionnaires
```

This detailed explanation covers the main components of the system. Would you like me to elaborate on any specific aspect or provide more technical details about any particular module?