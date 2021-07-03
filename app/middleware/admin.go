package middleware

/*
func Admin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		v, ok := context.GetOk(r,"permission")
		if !ok {
			failedAuth(w, http.StatusUnauthorized, `{error: true, message: "Unauthorized"}`)
		}
		if p, ok := v.(int); ok && p > 0 {
			next.ServeHTTP(w, r)
		} else {
			failedAuth(w, http.StatusUnauthorized, `{error: true, message: "Unauthorized"}`)
		}
	}
}
*/
