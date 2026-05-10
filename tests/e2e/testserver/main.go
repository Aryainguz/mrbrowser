package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintln(w, `
<!DOCTYPE html>
<html>
<head><title>Login</title></head>
<body>
	<h1>Welcome Back</h1>
	<form id="login-form">
		<div>
			<label>Email Address</label>
			<input type="email" id="email" />
		</div>
		<div>
			<label>Password</label>
			<input type="password" id="password" />
		</div>
		<button type="button" id="submit-btn" onclick="document.body.innerHTML = '<h2>Dashboard</h2>'">Sign In</button>
	</form>
</body>
</html>
		`)
	})

	http.HandleFunc("/dynamic", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintln(w, `
<!DOCTYPE html>
<html>
<head><title>Dynamic Selectors</title></head>
<body>
	<h1>Dynamic Page</h1>
	<!-- The ID and class change randomly/frequently, making selectors fragile -->
	<button id="btn_8347923" class="action-x23" onclick="document.body.innerHTML = '<h2>Clicked!</h2>'">Download Report</button>
</body>
</html>
		`)
	})

	log.Println("Starting test server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
