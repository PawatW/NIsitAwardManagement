package exceptions

import "errors"

var ErrUserNotFound = errors.New("user not found")
var ErrEmailAlreadyExists = errors.New("email already exists")
var ErrNisitIDAlreadyExists = errors.New("nisit id already exists")
var ErrUserInactive = errors.New("user account is deactivated")
var ErrCampusNotFound = errors.New("campus not found")
var ErrFacultyNotFound = errors.New("faculty not found")
var ErrDepartmentNotFound = errors.New("department not found")
var ErrInvalidCampusID = errors.New("invalid campus id format")
var ErrInvalidFacultyID = errors.New("invalid faculty id format")
var ErrInvalidDepartmentID = errors.New("invalid department id format")
var ErrFacultyNotInCampus = errors.New("faculty does not belong to campus")
var ErrDepartmentNotInFaculty = errors.New("department does not belong to faculty")
