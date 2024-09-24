package user

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
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

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req UserRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println("error parsing request:", err)
		resp := util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, "Error parsing request", struct{}{})
		resp.WriteResponse(w, r)
		return
	}

	// Validate request
	if err := h.validate.Struct(req); err != nil {
		log.Println("error validation request:", err)
		resp := util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, err.Error(), struct{}{})
		resp.WriteResponse(w, r)
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("error hashing password: %v", err)
		util.NewResponse(http.StatusInternalServerError, http.StatusInternalServerError, "Error hashing password", struct{}{}).WriteResponse(w, r)
		return
	}

	err = h.db.CreateUser(r.Context(), repo.CreateUserParams{
		Nama:     req.Nama,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     req.Role,
		Photo:    util.SqlString("static/default.png"),
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

	util.NewResponse(http.StatusOK, http.StatusOK, "User registered successfully", struct{}{}).WriteResponse(w, r)
}

func (h *Handler) UserInfo(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id")
	if userID == "" {
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}

	// Panggil metode yang mengeksekusi query GetCartByUserID
	userIDInt, ok := userID.(int32)
	if !ok {
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}

	user, err := h.db.GetUserByID(r.Context(), userIDInt)
	if err != nil {
		log.Printf("error getting user info: %v", err)
		util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, "Error getting user info", struct{}{}).WriteResponse(w, r)
		return
	}

	var res User
	res.UserID = user.UserID
	res.Nama = user.Nama
	res.Email = user.Email
	res.Role = user.Role
	res.Photo = user.Photo.String

	util.NewResponse(http.StatusOK, http.StatusOK, "User info successfully requested", res).WriteResponse(w, r)
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req UserRequest

	// Parse form data
	err := r.ParseMultipartForm(10 << 20) // batasan ukuran file (10MB)
	if err != nil {
		log.Printf("error parsing form data: %v", err)
		util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, "Error parsing form data", struct{}{}).WriteResponse(w, r)
		return
	}

	// Ambil data dari form
	req.Nama = r.FormValue("nama")
	req.Email = r.FormValue("email")
	req.Password = r.FormValue("password")
	req.Role = r.FormValue("role")

	// Ambil file photo dari form
	file, handler, err := r.FormFile("photo")
	if err != nil {
		log.Printf("error retrieving the file: %v", err)
		util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, "Error retrieving the file", struct{}{}).WriteResponse(w, r)
		return
	}
	defer file.Close()

	// Simpan file photo
	basePath, _ := os.Getwd()

	publicPath := path.Join(basePath, "public")

	photoPath := path.Join(publicPath, handler.Filename)
	dst, err := os.Create(photoPath)
	if err != nil {
		log.Printf("error saving the file: %v", err)
		util.NewResponse(http.StatusInternalServerError, http.StatusInternalServerError, "Error saving the file", struct{}{}).WriteResponse(w, r)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		log.Printf("error copying the file: %v", err)
		util.NewResponse(http.StatusInternalServerError, http.StatusInternalServerError, "Error copying the file", struct{}{}).WriteResponse(w, r)
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

	// Store user in the database
	err = h.db.CreateUser(r.Context(), repo.CreateUserParams{
		Email:    req.Email,
		Password: string(hashedPassword),
		Nama:     req.Nama,
		Role:     req.Role,
		Photo:    util.SqlString(photoPath), // Simpan path ke foto
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

func (h *Handler) GetAllUserByStudent(w http.ResponseWriter, r *http.Request) {
	// Mengambil semua user dengan role "student" dari database
	data, err := h.db.GetAllUserByStudent(r.Context(), "student")
	if err != nil {
		log.Println("error fetching student data:", err)
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

	userID := r.Context().Value("user_id")
	if userID == "" {
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}

	// Panggil metode yang mengeksekusi query GetCartByUserID
	userIDInt, ok := userID.(int32)
	if !ok {
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}

	var req UpdateUserRequest

	// Parse form data
	err := r.ParseMultipartForm(10 << 20) // batasan ukuran file (10MB)
	if err != nil {
		log.Printf("error parsing form data: %v", err)
		util.NewResponse(http.StatusBadRequest, http.StatusBadRequest, "Error parsing form data", struct{}{}).WriteResponse(w, r)
		return
	}

	// Ambil data dari form
	req.Nama = r.FormValue("nama")
	req.Email = r.FormValue("email")

	// Ambil file photo dari form
	file, handler, err := r.FormFile("photo")
	photoPath := path.Join("public", "default.png")
	if err == nil {
		defer file.Close()

		// Simpan file photo
		basePath, _ := os.Getwd()

		publicPath := path.Join(basePath, "public")

		photoPath = path.Join(publicPath, handler.Filename)
		dst, err := os.Create(photoPath)
		if err != nil {
			log.Printf("error saving the file: %v", err)
			util.NewResponse(http.StatusInternalServerError, http.StatusInternalServerError, "Error saving the file", struct{}{}).WriteResponse(w, r)
			return
		}
		defer dst.Close()

		if _, err := io.Copy(dst, file); err != nil {
			log.Printf("error copying the file: %v", err)
			util.NewResponse(http.StatusInternalServerError, http.StatusInternalServerError, "Error copying the file", struct{}{}).WriteResponse(w, r)
			return
		}

		photoPath = path.Join("static", handler.Filename)
	}

	// Update the user in the database
	err = h.db.UpdateUser(ctx, repo.UpdateUserParams{
		UserID: int32(userIDInt),
		Nama:   req.Nama,
		Email:  req.Email,
		Photo:  util.SqlString(photoPath),
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
	var resp *util.Response

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
		util.NewResponse(http.StatusNotFound, http.StatusNotFound, "User atau password salah", struct{}{}).WriteResponse(w, r)
		return
	}

	bcrypt.CompareHashAndPassword([]byte(data.Password), []byte(req.Password))
	if err != nil {
		log.Println("error no user:", err)
		util.NewResponse(http.StatusNotFound, http.StatusNotFound, "User atau password salah", struct{}{}).WriteResponse(w, r)
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
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "elearning",
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		log.Println("error creating token:", err)
		util.NewResponse(http.StatusNotFound, http.StatusNotFound, "Sesi tidak dapat dibuat untuk user", struct{}{}).WriteResponse(w, r)
		return
	}

	// Set token in cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   false, // Set to true if using HTTPS
		Path:     "/",
	})

	// Success response
	resp = util.NewResponse(http.StatusOK, http.StatusOK, "Login berhasil", map[string]interface{}{})
	resp.WriteResponse(w, r)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Now(),
		HttpOnly: true,
		Secure:   false,
		Path:     "/",
	})

	resp := util.NewResponse(http.StatusOK, http.StatusOK, "Logout berhasil", map[string]interface{}{})
	resp.WriteResponse(w, r)
}

func (h *Handler) ReadCookieAndVerifyToken(w http.ResponseWriter, r *http.Request) {

	// Deklarasi kunci rahasia JWT
	var jwtKey = []byte("my_secret_key") // Ganti "my_secret_key" dengan kunci rahasia Anda
	// Mendapatkan cookie bernama "token"
	cookie, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			// Jika cookie salah
			http.Error(w, "Cookie salah", http.StatusUnauthorized)
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

func (h *Handler) GetUserInfo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID := ctx.Value("user_id")
	role := ctx.Value("role")

	if (userID == nil) || (role == nil) {
		util.NewResponse(http.StatusUnauthorized, http.StatusUnauthorized, "Harap login terlebih dahulu", struct{}{}).WriteResponse(w, r)
		return
	}

	userIDString, _ := userID.(int32)
	roleString, _ := role.(string)

	res := make(map[string]string)
	res["user_id"] = fmt.Sprintf("%d", userIDString)
	res["role"] = roleString

	util.NewResponse(http.StatusOK, http.StatusOK, "OK", res).WriteResponse(w, r)
}
