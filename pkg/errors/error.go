package errors

var (
	ErrInconsistentIDs = New("inconsistent IDs")
	ErrAlreadyExists   = New("already exists")
	ErrNotFound        = New("not found")
)

var (
	// ErrUnauthorizedAccess indicates missing or invalid credentials provided
	// when accessing a protected resource.
	ErrUnauthorizedAccess = New("missing or invalid credentials provided")

	// ErrMalformedEntity indicates malformed entity specification (e.g.
	// invalid owner or ID).
	ErrMalformedEntity = New("malformed entity specification")

	// ErrCreateUUID indicates error in creating uuid for entity creation
	ErrCreateUUID = New("uuid creation failed")

	// ErrCreateEntity indicates error in creating entity or entities
	ErrCreateEntity = New("create entity failed")

	// ErrUpdateEntity indicates error in updating entity or entities
	ErrUpdateEntity = New("update entity failed")

	// ErrAuthorization indicates a failure occurred while authorizing the entity.
	ErrAuthorization = New("failed to perform authorization over the entity")

	// ErrViewEntity indicates error in viewing entity or entities
	ErrViewEntity = New("view entity failed")

	// ErrRemoveEntity indicates error in removing entity
	ErrRemoveEntity = New("remove entity failed")

	// ErrConnect indicates error in adding connection
	ErrConnect = New("add connection failed")

	// ErrDisconnect indicates error in removing connection
	ErrDisconnect = New("remove connection failed")

	// ErrFailedToRetrieveThings failed to retrieve things.
	ErrFailedToRetrieveThings = New("failed to retrieve group members")

	// ErrWrongPassword indicates error in wrong password
	ErrWrongPassword = New("Wrong Password")
)

type Error interface {

	// Error implements the error interface.
	Error() string

	// Msg returns error message
	Msg() string

	// Err returns wrapped error
	Err() Error
}

var _ Error = (*customError)(nil)

// customError struct represents a Mainflux error
type customError struct {
	msg string
	err Error
}

func (ce *customError) Error() string {
	if ce == nil {
		return ""
	}
	if ce.err == nil {
		return ce.msg
	}
	return ce.msg + " : " + ce.err.Error()
}

func (ce *customError) Msg() string {
	return ce.msg
}

func (ce *customError) Err() Error {
	return ce.err
}

// Contains inspects if e2 error is contained in any layer of e1 error
func Contains(e1 error, e2 error) bool {
	if e1 == nil || e2 == nil {
		return e2 == e1
	}
	ce, ok := e1.(Error)
	if ok {
		if ce.Msg() == e2.Error() {
			return true
		}
		return Contains(ce.Err(), e2)
	}
	return e1.Error() == e2.Error()
}

// Wrap returns an Error that wrap err with wrapper
func Wrap(wrapper error, err error) error {
	if wrapper == nil || err == nil {
		return wrapper
	}
	if w, ok := wrapper.(Error); ok {
		return &customError{
			msg: w.Msg(),
			err: cast(err),
		}
	}
	return &customError{
		msg: wrapper.Error(),
		err: cast(err),
	}
}

func cast(err error) Error {
	if err == nil {
		return nil
	}
	if e, ok := err.(Error); ok {
		return e
	}
	return &customError{
		msg: err.Error(),
		err: nil,
	}
}

// New returns an Error that formats as the given text.
func New(text string) Error {
	return &customError{
		msg: text,
		err: nil,
	}
}
