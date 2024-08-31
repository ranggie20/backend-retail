CREATE TABLE  "users" (
    user_id SERIAL PRIMARY KEY,
    nama VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(225) NOT NULL,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(225) NOT NULL,
    photo TEXT
);

CREATE TABLE teachers (
  teacher_id serial PRIMARY KEY,
  user_id serial,
  teacher_name VARCHAR (255) NOT NULL,
  created_at timestamp,
  deleted_at timestamp,
  updated_at timestamp
);

CREATE TABLE courses (
  course_id SERIAL PRIMARY KEY,
  course_name VARCHAR(255) NOT NULL,
  course_description TEXT NOT NULL,
  category_id INTEGER NULL,
  price INTEGER NOT NULL,
  thumbnail TEXT,
  created_at TIMESTAMP,
  deleted_at TIMESTAMP,
  updated_at TIMESTAMP
);

CREATE TABLE categories (
  category_id SERIAL PRIMARY KEY,
  category_name VARCHAR(255) NOT NULL,
  icon VARCHAR(255) NOT NULL,
  created_at TIMESTAMP,
  updated_at TIMESTAMP
);

CREATE TABLE courses_video (
  course_video_id SERIAL PRIMARY KEY,
  course_id INTEGER,
  course_video_name VARCHAR(255) NOT NULL,
  path_video TEXT NOT NULL,
  created_at TIMESTAMP,
  deleted_at TIMESTAMP,
  updated_at TIMESTAMP
);

CREATE TABLE transaction_history (
  transaction_history_id SERIAL PRIMARY KEY,
  subscription_id INTEGER,
  quantity INTEGER NOT NULL,
  total_amount INTEGER NOT NULL,
  is_paid VARCHAR(255) NOT NULL,
  subcriptions_start_date TIMESTAMP,
  proof TEXT,
  created_at TIMESTAMP,
  deleted_at TIMESTAMP,
  updated_at TIMESTAMP
);

CREATE TABLE wishlist (
  wishlist_id SERIAL PRIMARY KEY,
  user_id INTEGER,
  course_id INTEGER,
  created_at TIMESTAMP,
  deleted_at TIMESTAMP,
  updated_at TIMESTAMP
);

CREATE TABLE notification (
  notification_id SERIAL PRIMARY KEY,
  user_id INTEGER,
  course_id INTEGER,
  message TEXT,
  is_read VARCHAR(255) NOT NULL,
  created_at TIMESTAMP,
  updated_at TIMESTAMP
);

CREATE TABLE cart (
  cart_id SERIAL PRIMARY KEY,
  user_id INTEGER,
  course_id INTEGER,
  course_name VARCHAR(255) NOT NULL,
  price INTEGER,
  quantity INTEGER,
  total_amount INTEGER
);


CREATE TABLE subscriptions (
  subscription_id SERIAL PRIMARY KEY,
  user_id INTEGER,
  course_id INTEGER,
  payment_id INTEGER,
  is_correct VARCHAR(255) NOT NULL,
  created_at TIMESTAMP,
  updated_at TIMESTAMP
);

CREATE TABLE payment (
  payment_id SERIAL PRIMARY KEY,
  user_id INTEGER,
  payment_method_id INTEGER,
  payment_status_id INTEGER,
  total_amount INTEGER,
  payment_date TIMESTAMP
);

CREATE TABLE payment_method (
  payment_method_id SERIAL PRIMARY KEY,
  payment_method_name VARCHAR(255) NOT NULL,
  created_at TIMESTAMP
);

CREATE TABLE payment_status (
  payment_status_id SERIAL PRIMARY KEY,
  payment_status_name VARCHAR(255) NOT NULL,
  created_at TIMESTAMP
);
