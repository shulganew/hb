package services

import (
	"crypto/hmac"
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"fmt"
	"io/fs"
	"net/http"
	"regexp"
	"strings"

	"github.com/shulganew/hb.git/internal/api/oapi"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

//go:embed templates/register/*
var staticFiles embed.FS

// Register form.
func (k *Happy) NewUserForm(w http.ResponseWriter, r *http.Request, file string) {
	fsys, err := fs.Sub(staticFiles, "templates/register")
	if err != nil {
		zap.S().Errorln("Path file error:", err)
	}
	http.ServeFileFS(w, r, fsys, file)
}

// Register new use.
func (k *Happy) CreateUser(w http.ResponseWriter, r *http.Request) {
	tguser := r.FormValue("tg")
	name := r.FormValue("uname")
	pw := r.FormValue("psw")
	pwr := r.FormValue("pswr")
	hbd := r.FormValue("hb")

	// Remove injections.
	tguser = regexp.MustCompile(`[^a-zA-Z0-9@\s]+`).ReplaceAllString(tguser, "")
	name = regexp.MustCompile(`[^a-zA-Z0-9\s]+`).ReplaceAllString(name, "")
	pw = regexp.MustCompile(`[^a-zA-Z0-9#@\s]+`).ReplaceAllString(pw, "")
	hbd = regexp.MustCompile(`[^0-9#@\s-]+`).ReplaceAllString(hbd, "")

	zap.S().Debugf("New Registration! %s  %s  %s \n", tguser, name, hbd)
	// Check password.
	if pw != pwr {
		http.Redirect(w, r, Answer("Password missmatch!"), http.StatusSeeOther)
	}

	// Check tg user at prefix.
	tguser = strings.TrimPrefix(tguser, "@")

	// Set hash as user password.
	pwh, err := HashPassword(pw)
	if err != nil {
		zap.S().Errorln("Error creating hash from password")
		http.Redirect(w, r, Answer("Error creating hash from password: "+err.Error()), http.StatusSeeOther)
		return
	}

	err = k.stor.AddUser(r.Context(), tguser, name, pwh, hbd)
	if err != nil {
		zap.S().Errorln(err)
		http.Redirect(w, r, Answer("Error adding user: "+err.Error()), http.StatusSeeOther)
	}

	http.Redirect(w, r, Answer("User added!"), http.StatusSeeOther)
}

// Validate user in Keeper, if sucsess it return user's id.
func (k *Happy) Login(w http.ResponseWriter, r *http.Request, params oapi.LoginParams) {
	zap.S().Debugln("Get loing request: ", r.URL.RawQuery)
	isValid := k.validateTG(params)
	if isValid {
		k.memory.Add(params.Username)
		http.Redirect(w, r, Answer("User Loged in!"), http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/hb/register/reg.html", http.StatusSeeOther)
	}
}

// Validate telegram user auth request.
func (k *Happy) validateTG(params oapi.LoginParams) (isValid bool) {
	// Construct data_check_string.
	data := strings.Join([]string{"auth_date=" + params.AuthDate, "first_name=" + params.FirstName, "id=" + params.Id, "last_name=" + params.LastName, "photo_url=" + params.PhotoUrl, "username=" + params.Username}, "\n")

	// Get secret key.
	hasher := sha256.New()
	hasher.Write([]byte(k.conf.Bot))
	key := hasher.Sum(nil)

	// Get HMAC 256.
	sig := hmac.New(sha256.New, key)
	sig.Write([]byte(data))

	// Calculate message hash.
	cHash := hex.EncodeToString(sig.Sum(nil))

	// Check is valid telegram hash.
	if cHash == params.Hash {
		zap.S().Debugln("Valid")
		return true
	}
	zap.S().Debugln("Not Valid")
	return
}

// HashPassword returns the bcrypt hash of the password.
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedPassword), nil
}

// Answer page constructor.
func Answer(ans string) string {
	var sb strings.Builder
	sb.WriteString("/hb/register/status.html?status=")
	sb.WriteString(ans)
	return sb.String()
}
