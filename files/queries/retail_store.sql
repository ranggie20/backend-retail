
-- name: CreateUser :exec
INSERT INTO "users" (
    nama,
    email,
    password,
    role,
    photo
) VALUES (
    $1, $2, $3, $4, $5
);

-- name: GetAllUser :many
SELECT * FROM "users";

-- name: GetUserByID :one
SELECT * FROM "users" WHERE user_id = $1;

-- name: GetAllUserByTeacher :many
SELECT * FROM "users" WHERE role = $1;

-- name: UpdateUser :exec
UPDATE "users" SET nama = $1, email = $2, password = $3, role = $4, photo = $5 WHERE user_id = $6;

-- name: DeleteUser :exec
DELETE FROM "users" WHERE user_id = $1;

-- name: Login :one
SELECT * FROM "users" WHERE nama = $1;

-- name: CreateTeacher :exec
INSERT INTO teachers (
    teacher_id,
    user_id,    
    teacher_name,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5
);

-- name: GetAllTeacher :many
SELECT * FROM teachers;

-- name: GetTeacherByID :one
SELECT * FROM teachers WHERE teacher_id = $1;

-- name: UpdateTeacher :exec
UPDATE teachers SET teacher_id = $1, user_id = $2, teacher_name = $3 WHERE teacher_id = $4;

-- name: DeleteTeacher :exec
DELETE FROM teachers WHERE teacher_id = $1;


-- name: CreateCourse :exec
INSERT INTO courses (
  course_name,
  course_description,
  category_id,
  price,
  thumbnail,
  created_at,
  deleted_at,
  updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
);

-- name: GetCoursePrice :many
SELECT course_id, course_name, course_description , price, thumbnail
FROM courses
ORDER BY price DESC;


-- name: GetMyCoursePage :many 
SELECT c.course_id, c.course_name, c.course_description
FROM subscriptions s
LEFT JOIN  courses c  ON c.course_id = s.course_id WHERE user_id = $1;


-- name: GetAllCourse :many
SELECT * FROM courses;

-- name: GetCourseByID :one
SELECT * FROM courses WHERE course_id = $1;

-- name: GetCourseByNew :one
SELECT * FROM courses ORDER BY created_at DESC;

-- name: UpdateCourse :exec
UPDATE courses SET course_id = $1, course_name = $2,course_description = $3, category_id = $4, price = $5, thumbnail  = $6, created_at = $7, deleted_at = $8, updated_at= $9 WHERE course_id = $10;

-- name: DeleteCourse :exec
DELETE FROM courses WHERE course_id = $1;


-- name: CreateCategory :exec
INSERT INTO categories (
    category_name,
    icon,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4
);

-- name: GetAllCategories :many
SELECT * FROM categories;

-- name: GetCategory :many
SELECT c.course_id, c.course_name, cr.category_id, cr.category_name 
FROM courses c
JOIN categories cr ON c.category_id = cr.category_id 
WHERE cr.category_name;


-- name: GetCategoryByID :one
SELECT * FROM categories WHERE category_id = $1;

-- name: UpdateCategory :exec
 UPDATE categories SET category_id = $1, category_name = $2, icon = $3, updated_at = $4 WHERE category_id = $5;

-- name: DeleteCategory :exec
DELETE FROM categories WHERE category_id = $1;

-- name: CreateCourseVideo :exec
INSERT INTO courses_video (
    course_id,
    course_video_name,
    path_video,
    created_at,
    deleted_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6
);

-- name: GetCourse :one
 SELECT 
 	course_id,
 	thumbnail, 
 	course_name, 
 	price
 FROM courses c ;


-- name: GetAllCourseVideos :many
SELECT * FROM courses_video;

-- name: GetCourseVideo :many
SELECT c.course_id, c.course_name, c.course_description, cat.category_name,
       cv.course_video_id , cv.course_video_name , cv.path_video 
FROM courses c
LEFT JOIN categories cat ON c.category_id = cat.category_id 
LEFT JOIN courses_video cv  ON c.course_id = cv.course_id;


-- name: GetCourseVideoByID :one
SELECT * FROM courses_video WHERE course_video_id = $1;

-- name: UpdateCourseVideo :exec
UPDATE courses_video SET course_id = $1, course_video_name = $2, path_video = $3, updated_at = $4 WHERE course_video_id = $5;

-- name: DeleteCourseVideo :exec
DELETE FROM courses_video WHERE course_video_id = $1;


-- name: CreateWishlist :exec
INSERT INTO wishlist (
    user_id,
    course_id,
    created_at,
    deleted_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5
);

-- name: GetAllWishlists :many
SELECT * FROM wishlist;

-- name: GetWishlistByID :one
SELECT * FROM wishlist WHERE wishlist_id = $1;

-- name: UpdateWishlist :exec
UPDATE wishlist SET user_id = $2, course_id = $3, updated_at = $4 WHERE wishlist_id = $1;

-- name: DeleteWishlist :exec
DELETE FROM wishlist WHERE user_id = $1;

-- name: CreateCart :exec 
INSERT INTO cart(
  user_id ,
  course_id ,
  course_name ,
  price ,
  quantity ,
  total_amount 
) VALUES (
    $1, $2, $3, $4, $5, $6
);

-- name: GetAllCart :many
SELECT * FROM cart;

-- name: GetCart :many 
SELECT
	cs.course_id,
	cs.thumbnail,
    cs.course_name,
    cs.price,
    cr.quantity,
    cr.total_amount
FROM
    cart cr
LEFT JOIN courses cs 
ON cr.course_id = cs.course_id;


-- name: UpdateCart :exec
UPDATE cart SET cart_id = $1, course_id = $2, course_name = $3, price = $4, quantity = $5, total_amount = $6 WHERE cart_id = $7;

-- name: DeleteCart :exec
DELETE FROM cart WHERE user_id = $1;


-- name: CreateNotification :exec
INSERT INTO notification (
  user_id,
  course_id,
  message,
  is_read,
  created_at,
  updated_at
) VALUES (
  $1, $2, $3, $4, $5, $6
);

-- name: GetAllNotifications :many
SELECT * FROM notification;

-- name: UpdateNotification :exec
UPDATE notification
SET user_id = $1, course_id = $2, message = $3, is_read = $4, created_at = $5, updated_at = $6 WHERE notification_id = $7;

-- name: DeleteNotification :exec
DELETE FROM notification WHERE notification_id = $1;


-- name: CreateSubscription :exec
INSERT INTO subscriptions (
    user_id,
    course_id,
    cart_id,
    is_correct,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6
);

-- name: GetAllSubscriptions :many
SELECT * FROM subscriptions;

--  query untuk melihat 6 course terpopuler/ popular course --

-- name: GetPopularCourse :many
SELECT 
    c.course_id, 
    c.course_name, 
    c.thumbnail,
    COUNT(s.course_id) AS total_enrollments
FROM 
    courses c
LEFT JOIN 
    subscriptions s ON s.course_id = c.course_id 
LEFT JOIN 
    transaction_history th ON th.subscription_id = s.subscription_id 
GROUP BY 
    c.course_id, 
    c.course_name,
    c.thumbnail
ORDER BY 
    total_enrollments DESC
LIMIT 6;



-- name: GetSubscriptionByID :one
SELECT * FROM subscriptions WHERE subscription_id = $1;


-- name: CreatePayment :exec
INSERT INTO payment (
    user_id,
    course_id,
    subscription_id,
    payment_method_id,
    payment_status_id,
    total_amount,
    payment_date
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
);

-- name: GetPayment :many
 SELECT 
    c.cart_id, 
    p.payment_id, 
    c.course_name, 
    c.quantity, 
    p.total_amount, 
    pm.payment_method_name, 
    p.payment_date
  FROM payment p 
  LEFT JOIN payment_method pm 
    ON p.payment_method_id = pm.payment_method_id 
  LEFT JOIN subscriptions s 
    ON s.subscription_id = p.subscription_id 
  LEFT JOIN cart c  
    ON c.cart_id = s.cart_id;

-- name: GetPaymentHistory :many
SELECT p.payment_id, c.course_id, c.course_name, p.total_amount, ps.payment_status_name, th.subcriptions_start_date 
FROM transaction_history th
LEFT JOIN subscriptions s 
ON th.subscription_id = s.subscription_id
left join payment p 
on p.subscription_id = s.subscription_id 
left join payment_status ps 
on ps.id = p.payment_status_id 
LEFT JOIN courses c 
ON s.course_id = c.course_id
ORDER BY th.subcriptions_start_date DESC;



-- name: GetAllPayment :many
SELECT * FROM payment;

-- name: GetPaymentByID :one
SELECT * FROM payment WHERE payment_id = $1;

-- name: CreatePaymentMethod :exec
INSERT INTO payment_method (
    payment_method_name,
    created_at
) VALUES (
    $1, $2
);

-- name: GetAllPaymentMethod :many
SELECT * FROM payment_method;

-- name: GetPaymentMethodByID :one
SELECT * FROM payment_method WHERE payment_method_id = $1;


-- name: CreateTransactionHistory :exec
INSERT INTO transaction_history (
    subscriptions_id,
    quantity,
    total_amount,
    is_paid,
    subcriptions_start_date,
    proof,
    created_at,
    updated_at,
    deleted_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
);

-- name: GetAllTransactionHistory :many
SELECT * FROM transaction_history;


-- name: CreatePaymentStatus :exec
INSERT INTO payment_status (
    payment_status_name,
    created_at
) VALUES (
    $1, $2
);

-- name: GetAllPaymentStatus :many
SELECT * FROM payment_status;

-- name: GetPaymentStatusByID :one
SELECT * FROM payment_status WHERE payment_status_id = $1;










