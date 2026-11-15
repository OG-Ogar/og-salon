package handlers

import "net/http"

// HairstylistRouter dispatches every /hairstylist/ URL.
//
// Hair stylists share the barber account system and workspace. Only the
// branded sign-in page lives here; account recovery and registration
// reuse the existing barber pages so there is one account per provider.
func HairstylistRouter(w http.ResponseWriter, r *http.Request) {
	parts := pathParts(r.URL.Path)

	if len(parts) == 1 {
		http.Redirect(w, r, "/hairstylist/login", http.StatusSeeOther)

		return
	}

	if len(parts) == 2 {
		switch parts[1] {
		case "login":
			HairstylistLoginHandler(w, r)

			return

		case "register":
			http.Redirect(w, r, "/barber/register", http.StatusSeeOther)

			return

		case "forgot-password":
			http.Redirect(w, r, "/barber/forgot-password", http.StatusSeeOther)

			return

		case "reset-password":
			target := "/barber/reset-password"
			if r.URL.RawQuery != "" {
				target += "?" + r.URL.RawQuery
			}
			http.Redirect(w, r, target, http.StatusSeeOther)

			return
		}
	}

	journeyError(w, "Page not found.", http.StatusNotFound)
}
