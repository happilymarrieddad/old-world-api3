CREATE CONSTRAINT constraint_user_unq_email FOR (u:User) REQUIRE u.email IS UNIQUE;
