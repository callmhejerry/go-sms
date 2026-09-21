-- students
ALTER TABLE students
RENAME CONSTRAINT students_tenant_admission_number_unique
TO students_tenant_id_admission_number_key;


-- users
ALTER TABLE users
RENAME CONSTRAINT users_tenant_id_email_unique
TO  users_tenant_id_email_key;


-- roles
ALTER TABLE roles
RENAME CONSTRAINT roles_tenant_id_name_unique
TO roles_tenant_id_name_key;


-- academic sessions
ALTER TABLE academic_sessions
RENAME CONSTRAINT academic_sessions_tenant_id_name_unique
TO academic_sessions_tenant_id_name_key;


-- classes
ALTER TABLE classes
RENAME CONSTRAINT classes_tenant_id_name_unique
TO classes_tenant_id_name_key;


-- class_arms
ALTER TABLE class_arms
RENAME CONSTRAINT class_arms_tenant_id_class_id_name_unique
TO class_arms_tenant_id_class_id_name_key;
