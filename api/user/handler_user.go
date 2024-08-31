package user

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/online-bnsp/backend/middleware/auth"
	repo "github.com/online-bnsp/backend/repo/generated"
	"github.com/online-bnsp/backend/util"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req UserRequest

	// Decode JSON request body
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Printf("error parsing request: %v", err)
		util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, "Error parsing request", struct{}{}).WriteResponse(w, r)
		return
	}

	// Validate request
	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		var errMsg strings.Builder
		for _, err := range err.(validator.ValidationErrors) {
			errMsg.WriteString(err.Field() + " is invalid; ")
		}
		log.Printf("validation error: %v", errMsg.String())
		util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, "Invalid input: "+errMsg.String(), struct{}{}).WriteResponse(w, r)
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("error hashing password: %v", err)
		util.NewResponse(http.StatusInternalServerError, http.StatusInternalServerError, "Error hashing password", struct{}{}).WriteResponse(w, r)
		return
	}

	fmt.Println(req)

	// Store user in the database
	err = h.db.CreateUser(r.Context(), repo.CreateUserParams{
		Email:    req.Email,
		Password: string(hashedPassword),
		Nama:     req.Nama,
		Role:     req.Role,
		Photo:    util.SqlString(req.Photo),
	})

	if err != nil {
		// Handle duplicate key error
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			log.Printf("duplicate key error: %v", err)
			util.NewResponse(http.StatusConflict, http.StatusConflict, "Email already exists", struct{}{}).WriteResponse(w, r)
			return
		}

		// Log and respond for other errors
		log.Printf("error creating user: %v", err)
		util.NewResponse(http.StatusInternalServerError, http.StatusInternalServerError, "Error creating user", struct{}{}).WriteResponse(w, r)
		return
	}

	// Create a response with the user data
	responseData := map[string]interface{}{
		"email": req.Email,
		"nama":  req.Nama,
	}

	util.NewResponse(http.StatusOK, http.StatusOK, "User registered successfully", responseData).WriteResponse(w, r)
}

func (h *Handler) GetAllUser(w http.ResponseWriter, r *http.Request) {
	data, err := h.db.GetAllUser(r.Context())
	if err != nil {
		log.Println("error fetching all data item:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var res []User
	for _, d := range data {
		res = append(res, User{
			UserID:   d.UserID,
			Nama:     d.Nama,
			Email:    d.Email,
			Password: d.Password, // Password biasanya tidak di-return, pertimbangkan untuk menghilangkan ini
			Role:     d.Role,
			Photo:    d.Photo.String,
		})
	}

	util.NewResponse(http.StatusOK, http.StatusOK, "", res).WriteResponse(w, r)
}

func (h *Handler) GetAllUserByTeacher(w http.ResponseWriter, r *http.Request) {
	// Mengambil semua user dengan role "teacher" dari database
	data, err := h.db.GetAllUserByTeacher(r.Context(), "teacher")
	if err != nil {
		log.Println("error fetching teacher data:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Menyiapkan slice untuk response
	var res []User
	for _, d := range data {
		res = append(res, User{
			UserID: d.UserID,
			Nama:   d.Nama,
			Email:  d.Email,
			Role:   d.Role,
			Photo:  d.Photo.String,
		})
	}

	// Mengirim response dengan status OK
	util.NewResponse(http.StatusOK, http.StatusOK, "", res).WriteResponse(w, r)
}

func (h *Handler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	// Mendapatkan parameter ID dari URL
	vars := chi.URLParam(r, "id")

	// Parsing ID dari string ke int32
	id, err := strconv.ParseInt(vars, 10, 32)
	if err != nil {
		log.Println("error parsing ID:", err)
		util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, "Invalid ID format", struct{}{}).WriteResponse(w, r)
		return
	}

	// Mendapatkan data UserInfo dari database berdasarkan ID
	data, err := h.db.GetUserByID(r.Context(), int32(id))
	if err == sql.ErrNoRows {
		util.NewResponse(http.StatusUnauthorized, http.StatusUnauthorized, "Invalid ID", struct{}{}).WriteResponse(w, r)
		return
	} else if err != nil {
		log.Println("error fetching user info by ID:", err)
		util.NewResponse(http.StatusInternalServerError, http.StatusInternalServerError, "Internal server error", struct{}{}).WriteResponse(w, r)
		return
	}

	// Mapping data UserInfo ke struktur respons
	res := User{
		UserID: data.UserID,
		Nama:   data.Nama,
		Email:  data.Email,
		Role:   data.Role,
	}

	// Mengirimkan respons
	util.NewResponse(http.StatusOK, http.StatusOK, "", res).WriteResponse(w, r)
}
func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	resp := util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, "Bad request", map[string]interface{}{})

	// Get the user ID from the URL parameters
	userID := chi.URLParam(r, "id")
	if userID == "" {
		resp = util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, "User ID is required", struct{}{})
		resp.WriteResponse(w, r)
		return
	}

	// Convert userID to int32
	id, err := strconv.Atoi(userID)
	if err != nil {
		log.Println("invalid user ID:", err)
		resp = util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, "Invalid user ID", struct{}{})
		resp.WriteResponse(w, r)
		return
	}

	var req UserRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println("error parsing request:", err)
		resp = util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, "Error parsing request", struct{}{})
		resp.WriteResponse(w, r)
		return
	}

	// Validate request
	if err := h.validate.Struct(req); err != nil {
		log.Println("error validation request:", err)
		resp = util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, err.Error(), struct{}{})
		resp.WriteResponse(w, r)
		return
	}

	// Update the user in the database
	err = h.db.UpdateUser(ctx, repo.UpdateUserParams{
		UserID:   int32(id),
		Nama:     req.Nama,
		Email:    req.Email,
		Password: req.Password,
		Role:     req.Role,
		Photo:    util.SqlString(req.Photo),
	})

	if err != nil {
		log.Println("error updating user in db:", err)
		resp = util.NewResponse(http.StatusInternalServerError, http.StatusInternalServerError, "Try again later", struct{}{})
		resp.WriteResponse(w, r)
		return
	}

	resp.Status = http.StatusOK
	resp.Code = http.StatusOK
	resp.Message = "User updated successfully"
	resp.WriteResponse(w, r)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	resp := util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, "Bad request", map[string]interface{}{})

	var jwtKey = []byte("your_secret_key") // Replace with your secret key
	var req LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println("error parsing request:", err)
		resp = util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, "Error parsing request", struct{}{})
		resp.WriteResponse(w, r)
		return
	}

	data, err := h.db.Login(ctx, req.Nama)
	if err != nil {
		log.Println("error no user:", err)
		http.Error(w, "User salah", http.StatusNotFound)
		return
	}

	bcrypt.CompareHashAndPassword([]byte(data.Password), []byte(req.Password))
	if err != nil {
		log.Println("error no user:", err)
		http.Error(w, "password salah", http.StatusNotFound)
		return
	}

	// Determine role from user data
	role := data.Role

	// Create JWT claims
	claims := &auth.Claims{
		Username: req.Nama,
		UserID:   data.UserID,
		Role:     role, // Ensure role is included
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(15 * time.Minute).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "your-app-name",
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		log.Println("error creating token:", err)
		http.Error(w, "Gagal membuat token", http.StatusInternalServerError)
		return
	}

	// Set token in cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  time.Now().Add(15 * time.Minute),
		HttpOnly: true,
		Secure:   false, // Set to true if using HTTPS
		Path:     "/",
	})

	// Success response
	resp = util.NewResponse(http.StatusOK, http.StatusOK, "Login berhasil", map[string]interface{}{})
	resp.WriteResponse(w, r)
}

func (h *Handler) ReadCookieAndVerifyToken(w http.ResponseWriter, r *http.Request) {

	// Deklarasi kunci rahasia JWT
	var jwtKey = []byte("my_secret_key") // Ganti "my_secret_key" dengan kunci rahasia Anda
	// Mendapatkan cookie bernama "token"
	cookie, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			// Jika cookie tidak ditemukan
			http.Error(w, "Cookie tidak ditemukan", http.StatusUnauthorized)
			return
		}
		// Jika terjadi error lain saat membaca cookie
		http.Error(w, "Gagal membaca cookie", http.StatusBadRequest)
		return
	}

	// Mendapatkan nilai token dari cookie
	tokenStr := cookie.Value

	// Klaim untuk menyimpan data yang akan diekstrak dari token JWT
	claims := &Claims{}

	// Memverifikasi token JWT
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil // Menggunakan jwtKey yang sama dengan yang digunakan untuk menandatangani token
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			http.Error(w, "Token tidak valid", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Gagal memverifikasi token", http.StatusBadRequest)
		return
	}

	if !token.Valid {
		http.Error(w, "Token tidak valid", http.StatusUnauthorized)
		return
	}

	// Token valid, Anda dapat menggunakan klaim yang ada di dalam token
	fmt.Fprintf(w, "Welcome %s!", claims.Username)
}
