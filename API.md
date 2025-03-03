# Language Learning Portal API Documentation

## Groups

### GET /api/groups
Returns a paginated list of word groups.

```json
{
  "items": [
    {
      "id": 1,
      "name": "Basic Greetings",
      "word_count": 4
    }
  ],
  "pagination": {
    "current_page": 1,
    "total_pages": 1,
    "total_items": 1,
    "items_per_page": 100
  }
}
```

### GET /api/groups/:id
Returns details of a specific group.

```json
{
  "id": 1,
  "name": "Basic Greetings",
  "word_count": 4,
  "words": [
    {
      "id": 1,
      "arabic": "مرحبا",
      "roman": "marhaban",
      "english": "hello"
    }
  ]
}
```

### GET /api/groups/:id/words
Returns a paginated list of words in a group.

```json
{
  "items": [
    {
      "id": 1,
      "arabic": "مرحبا",
      "roman": "marhaban",
      "english": "hello"
    }
  ],
  "pagination": {
    "current_page": 1,
    "total_pages": 1,
    "total_items": 4,
    "items_per_page": 100
  }
}
```

### GET /api/groups/:id/study_sessions
Returns a paginated list of study sessions for a group.

```json
{
  "items": [
    {
      "id": 1,
      "activity_name": "Vocabulary Quiz",
      "group_name": "Basic Greetings",
      "start_time": "2025-02-20T08:55:56Z",
      "review_items_count": 3
    }
  ],
  "pagination": {
    "current_page": 1,
    "total_pages": 1,
    "total_items": 1,
    "items_per_page": 100
  }
}
```

## Words

### GET /api/words
Returns a paginated list of all words.

```json
{
  "items": [
    {
      "id": 1,
      "arabic": "مرحبا",
      "roman": "marhaban",
      "english": "hello"
    }
  ],
  "pagination": {
    "current_page": 1,
    "total_pages": 1,
    "total_items": 4,
    "items_per_page": 100
  }
}
```

### GET /api/words/:id
Returns details of a specific word.

```json
{
  "id": 1,
  "arabic": "مرحبا",
  "roman": "marhaban",
  "english": "hello",
  "groups": [
    {
      "id": 1,
      "name": "Basic Greetings"
    }
  ]
}
```

## Study Activities

### GET /api/study_activities
Returns a paginated list of study activities.

```json
{
  "items": [
    {
      "id": 1,
      "name": "Vocabulary Quiz",
      "thumbnail_url": "/images/vocab-quiz.png",
      "description": "Test your vocabulary knowledge",
      "launch_url": "/activities/vocab-quiz",
      "stats": {
        "total_sessions": 1,
        "total_words_reviewed": 3,
        "accuracy_rate": 66.67
      }
    }
  ],
  "pagination": {
    "current_page": 1,
    "total_pages": 1,
    "total_items": 1,
    "items_per_page": 100
  }
}
```

### GET /api/study_activities/:id
Returns details of a specific study activity.

```json
{
  "id": 1,
  "name": "Vocabulary Quiz",
  "thumbnail_url": "/images/vocab-quiz.png",
  "description": "Test your vocabulary knowledge",
  "launch_url": "/activities/vocab-quiz",
  "stats": {
    "total_sessions": 1,
    "total_words_reviewed": 3,
    "accuracy_rate": 66.67
  },
  "recent_sessions": [
    {
      "id": 1,
      "group_name": "Basic Greetings",
      "start_time": "2025-02-20T08:55:56Z",
      "stats": {
        "total_words": 3,
        "correct_count": 2,
        "wrong_count": 1
      }
    }
  ]
}
```

### POST /api/study_activities
Creates a new study activity.

Request:
```json
{
  "name": "Flashcards",
  "thumbnail_url": "/images/flashcards.png",
  "description": "Learn with interactive flashcards",
  "launch_url": "/activities/flashcards"
}
```

Response:
```json
{
  "id": 2,
  "name": "Flashcards",
  "thumbnail_url": "/images/flashcards.png",
  "description": "Learn with interactive flashcards",
  "launch_url": "/activities/flashcards"
}
```

## Study Sessions

### GET /api/study_sessions
Returns a paginated list of study sessions.

```json
{
  "items": [
    {
      "id": 1,
      "activity_name": "Vocabulary Quiz",
      "group_name": "Basic Greetings",
      "start_time": "2025-02-20T08:55:56Z",
      "stats": {
        "total_words": 3,
        "correct_count": 2,
        "wrong_count": 1
      }
    }
  ],
  "pagination": {
    "current_page": 1,
    "total_pages": 1,
    "total_items": 1,
    "items_per_page": 100
  }
}
```

### GET /api/study_sessions/:id/words
Returns a paginated list of words reviewed in a study session.

```json
{
  "items": [
    {
      "arabic": "مرحبا",
      "roman": "marhaban",
      "english": "hello",
      "is_correct": true,
      "reviewed_at": "2025-02-20T08:55:56Z"
    }
  ],
  "pagination": {
    "current_page": 1,
    "total_pages": 1,
    "total_items": 3,
    "items_per_page": 100
  }
}
```

### POST /api/study_sessions/:id/words/:word_id/review
Records a word review in a study session.

Request:
```json
{
  "is_correct": true
}
```

Response:
```json
{
  "id": 4,
  "word_id": 1,
  "session_id": 1,
  "is_correct": true,
  "reviewed_at": "2025-02-26T11:03:23Z"
}
```

## Dashboard

### GET /api/dashboard/last_study_session
Returns details of the last study session.

```json
{
  "id": 1,
  "activity_name": "Vocabulary Quiz",
  "group_name": "Basic Greetings",
  "start_time": "2025-02-20T08:55:56Z",
  "stats": {
    "total_words": 3,
    "correct_count": 2,
    "wrong_count": 1
  }
}
```

### GET /api/dashboard/study_progress
Returns study progress statistics.

```json
{
  "daily_stats": [
    {
      "date": "2025-02-20",
      "correct_count": 2,
      "wrong_count": 1
    }
  ],
  "total_stats": {
    "total_words_studied": 3,
    "total_correct": 2,
    "total_wrong": 1,
    "accuracy_rate": 66.67
  }
}
```

### GET /api/dashboard/quick-stats
Returns quick overview statistics.

```json
{
  "total_words_available": 4,
  "words_studied": 3,
  "study_sessions_completed": 1,
  "last_study_session": {
    "activity_name": "Vocabulary Quiz",
    "group_name": "Basic Greetings",
    "correct_count": 2,
    "wrong_count": 1
  }
}
```

## Reset Endpoints

### POST /api/reset_history
Deletes all study history (sessions and reviews) but keeps words and groups.

```json
{
  "message": "Study history has been reset successfully"
}
```

### POST /api/full_reset
Deletes all data and resets the database to initial state.

```json
{
  "message": "Database has been fully reset successfully"
}
```

## Error Responses

All endpoints return appropriate error responses in this format:

```json
{
  "error": "Error message here",
  "code": "ERROR_CODE_HERE"
}
```

Common error codes:
- `INVALID_REQUEST`
- `NOT_FOUND`
- `INTERNAL_ERROR`
