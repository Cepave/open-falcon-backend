package gin

// ViableParamValue defines the value of param and whether or not the value is viable
//
// In order to indicate the value of paramter of bool, integer... ,
// this object has a boolean filed to indicate whether or not the value is sensible.
type ViableParamValue struct {
	// Whether or not the parameter is viable(not empty)
	Viable bool
	// The value of parameter
	Value interface{}
	// The error of processing value of parameter
	Error error
}
