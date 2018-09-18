package middleware

// ValueMap sends context as a key-value between the chaining of Middleware.
type ValueMap map[string]interface{}

// Set is a helper method for setting ValueMap pointer's value.
func (v *ValueMap) Set(key string, value interface{}) {
	(*v)[key] = value
}

// Get is a helper method for getting ValueMap pointer's value.
func (v *ValueMap) Get(key string) interface{} {
	return (*v)[key]
}
