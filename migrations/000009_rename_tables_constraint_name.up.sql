-- students
ALTER TABLE students
RENAME CONSTRAINT students_tenant_id_admission_number_key
TO students_tenant_id_admission_number_unique;


-- users
ALTER TABLE users
RENAME CONSTRAINT users_tenant_id_email_key
TO users_tenant_id_email_unique;


-- roles
ALTER TABLE roles
RENAME CONSTRAINT roles_tenant_id_name_key
TO roles_tenant_id_name_unique;


-- academic sessions
ALTER TABLE academic_sessions
RENAME CONSTRAINT academic_sessions_tenant_id_name_key
TO academic_sessions_tenant_id_name_unique;


-- classes
ALTER TABLE classes
RENAME CONSTRAINT classes_tenant_id_name_key
TO classes_tenant_id_name_unique;


-- class_arms
ALTER TABLE class_arms
RENAME CONSTRAINT class_arms_tenant_id_class_id_name_key
TO class_arms_tenant_id_class_id_name_unique;
