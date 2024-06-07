package services

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// Register new use.
func (k *Happy) CreateUser(w http.ResponseWriter, r *http.Request) {
	zap.S().Debugln("Create new user:", r.URL.Host)
	/*
		var user oapi.NewUser
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			zap.S().Errorln("Can't decode json")
			// If can't decode 400
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Set hash as user password.
		hash, err := HashPassword(user.Password)
		if err != nil {
			zap.S().Errorln("Error creating hash from password")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Add user to database.
		_, err = k.stor.AddUser(r.Context(), user.Login, hash, user.Email)
		if err != nil {
			var pgErr *pq.Error
			// If URL exist in DataBase
			if errors.As(err, &pgErr) && pgerrcode.UniqueViolation == pgErr.Code {
				errt := "User's login is used"
				zap.S().Infoln(errt, err)
				http.Error(w, errt, http.StatusConflict)
				return
			}
			return
		}

		// set status code 201
		w.WriteHeader(http.StatusCreated)
	*/
}

// Validate user in Keeper, if sucsess it return user's id.
func (k *Happy) Login(w http.ResponseWriter, r *http.Request) {
	/*
		var nuser oapi.NewUser
		if err := json.NewDecoder(r.Body).Decode(&nuser); err != nil {
			// If can't decode 400
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Get User from storage
		dbUser, err := k.stor.GetByLogin(r.Context(), nuser.Login)
		if err != nil {
			zap.S().Infoln("User not found by login. ", err)
			http.Error(w, "Wrong login or password or otp", http.StatusUnauthorized)
			return
		}
		// Check OTP is correct.
		valid := totp.Validate(nuser.Otp, dbUser.Secret)
		if !valid {
			zap.S().Infoln("User enter wrong OTP. ")
			http.Error(w, "Wrong login or password or otp", http.StatusUnauthorized)
			return
		}

		// Check pass is correct
		err = k.CheckPassword(nuser.Password, dbUser.PassHash)
		if err != nil {
			http.Error(w, "Wrong login or password or otp", http.StatusUnauthorized)
			return
		}

		// Create jwt with access.
		// For user with login = admin grand rigts to admin handlers
		access := ""
		if nuser.Login == "admin" {
			access = "admin"
		}
		allowAll, err := k.ua.CreateJWSWithClaims(dbUser.UUID.String(), []string{access})
		if err != nil {
			zap.S().Errorln("Error creating jwt string: ", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Add("Content-Type", "text/plain")
		w.Header().Add("Authorization", config.AuthPrefix+string(allowAll))

		// set status code 200
		w.WriteHeader(http.StatusOK)

		_, err = w.Write([]byte("User loged in."))
		if err != nil {
			zap.S().Errorln("Can't write to response in LoginUser handler", err)
		}
	*/
}

// HashPassword returns the bcrypt hash of the password.
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedPassword), nil
}
