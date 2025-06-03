CREATE TABLE account(
                        user_id INT,
                        user_name VARCHAR(255) NOT NULL,
                        user_password VARCHAR(512) NOT NULL,
                        user_email VARCHAR(255) NOT NULL,
                        user_profile_picture VARCHAR(255) NOT NULL,
                        user_bio VARCHAR(255) NOT NULL,
                        user_last_connection DATETIME NOT NULL,
                        PRIMARY KEY(user_id)
);

CREATE TABLE feed(
                     feed_id INT,
                     feed_title VARCHAR(255) NOT NULL,
                     feed_description VARCHAR(255) NOT NULL,
                     feed_state VARCHAR(255) NOT NULL,
                     feed_creation_date DATETIME NOT NULL ,
                     user_id INT NOT NULL,
                     PRIMARY KEY(feed_id),
                     FOREIGN KEY(user_id) REFERENCES account(user_id)
);

CREATE TABLE tag(
                    tag_id INT,
                    tag_name VARCHAR(255) NOT NULL,
                    tag_type VARCHAR(255) NOT NULL ,
                    PRIMARY KEY(tag_id)
);

CREATE TABLE post(
                     post_id INT ,
                     post_date DATETIME NOT NULL ,
                     post_content VARCHAR(255) NOT NULL,
                     feed_id VARCHAR INT NOT NULL,
                     user_id INT NOT NULL,
                     PRIMARY KEY(post_id),
                     FOREIGN KEY(feed_id) REFERENCES feed(feed_id),
                     FOREIGN KEY(user_id) REFERENCES account(user_id)
);

CREATE TABLE picture(
                        picture_id INT ,
                        picture_blob VARCHAR(255) NOT NULL ,
                        post_id INT NOT NULL,
                        PRIMARY KEY(picture_id),
                        FOREIGN KEY(post_id) REFERENCES post(post_id)
);

CREATE TABLE comment(
                        comment_id INT  ,
                        comment_date DATETIME NOT NULL ,
                        comment_content VARCHAR(255),
                        comment_id_1 VARCHAR INT NOT NULL,
                        post_id INT NOT NULL,
                        user_id INT NOT NULL,
                        PRIMARY KEY(comment_id),
                        FOREIGN KEY(comment_id_1) REFERENCES comment(comment_id),
                        FOREIGN KEY(post_id) REFERENCES post(post_id)
);

CREATE TABLE add_tag(
                        feed_id INT ,
                        tag_id INT NOT NULL,
                        PRIMARY KEY(feed_id, tag_id),
                        FOREIGN KEY(feed_id) REFERENCES feed(feed_id),
                        FOREIGN KEY(tag_id) REFERENCES tag(tag_id)
);

CREATE TABLE reaction(
                         user_id INT NOT NULL ,
                         post_id INT NOT NULL,
                         PRIMARY KEY(user_id, post_id),
                         FOREIGN KEY(user_id) REFERENCES account(user_id),
                         FOREIGN KEY(post_id) REFERENCES post(post_id)
);


