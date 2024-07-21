package swagger

// ResponseUserRegistered indicates that the user was registered successfully.
// @Description User registered successfully
// @Response
type ResponseUserRegistered struct {
	Message string `json:"message"`
}

// ResponseUserLoggedIn indicates that the user was logged in successfully.
// @Description User logged in successfully
// @Response
type ResponseUserLoggedIn struct {
	Message string `json:"message"`
}

// ResponseUserLoggedOut indicates that the user was logged out successfully.
// @Description User logged out successfully
// @Response
type ResponseUserLoggedOut struct {
	Message string `json:"message"`
}

// ResponseTokenVerified indicates that the token is valid.
// @Description Token valid
// @Response
type ResponseTokenVerified struct {
	Message string `json:"message"`
}

// ResponseBadRequest indicates that the request was invalid.
// @Description Bad request
// @Response
type ResponseBadRequest struct {
	Error string `json:"error"`
}

// ResponseInternalServerError indicates an internal server error.
// @Description Internal server error
// @Response
type ResponseInternalServerError struct {
	Error string `json:"error"`
}

// ResponseUnauthorized indicates that the user is not authorized.
// @Description Unauthorized
// @Response
type ResponseUnauthorized struct {
	Error string `json:"error"`
}

// ResponseNoteRetrievedRestricted represents the response for a retrieved note.
// @Description Note retrieved successfully
// @Response
type ResponseNoteRetrievedRestricted struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

// ResponseNoteDeleted indicates that the note was deleted successfully.
// @Description Note deleted successfully
// @Response
type ResponseNoteDeleted struct {
	Message string `json:"message"`
}

// ResponseNoteNotFound indicates that the note was not found.
// @Description Note not found
// @Response
type ResponseNoteNotFound struct {
	Error string `json:"error"`
}

// ResponseNoteExpired indicates that the note has expired.
// @Description Note expired
// @Response
type ResponseNoteExpiredOrReachedMaxViews struct {
	Error string `json:"error"`
}

// ResponeNotFound indicates that the resource was not found.
// @Description Not found
// @Response
type ResponseNotFound struct {
	Error string `json:"error"`
}
