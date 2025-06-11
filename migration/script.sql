CREATE TABLE account(
                        user_id INT AUTO_INCREMENT NOT NULL,
                        user_name VARCHAR(255) NOT NULL,
                        user_password VARCHAR(512) NOT NULL,
                        user_email VARCHAR(255) NOT NULL,
                        user_profile_picture VARCHAR(255) NOT NULL DEFAULT "none",
                        user_bio VARCHAR(255) NOT NULL DEFAULT "this is my bio.",
                        user_last_connection DATETIME NOT NULL,
                        PRIMARY KEY(user_id)
);

CREATE TABLE feed(
                     feed_id INT AUTO_INCREMENT NOT NULL,
                     feed_title VARCHAR(255) NOT NULL,
                     feed_description VARCHAR(255) NOT NULL,
                     feed_state VARCHAR(255) NOT NULL,
                     feed_creation_date DATETIME NOT NULL ,
                     user_id INT NOT NULL,
                     PRIMARY KEY(feed_id),
                     FOREIGN KEY(user_id) REFERENCES account(user_id)
);

CREATE TABLE tag(
                    tag_id INT AUTO_INCREMENT NOT NULL,
                    tag_name VARCHAR(255) NOT NULL,
                    tag_type VARCHAR(255) NOT NULL ,
                    PRIMARY KEY(tag_id)
);

CREATE TABLE post(
                     post_id INT AUTO_INCREMENT NOT NULL,
                     post_date DATETIME NOT NULL ,
                     post_content VARCHAR(255) NOT NULL,
                     feed_id INT NOT NULL,
                     user_id INT NOT NULL,
                     PRIMARY KEY(post_id),
                     FOREIGN KEY(feed_id) REFERENCES feed(feed_id),
                     FOREIGN KEY(user_id) REFERENCES account(user_id)
);

CREATE TABLE picture(
                        picture_id INT AUTO_INCREMENT NOT NULL,
                        picture_blob VARCHAR(255) NOT NULL ,
                        post_id INT NOT NULL,
                        PRIMARY KEY(picture_id),
                        FOREIGN KEY(post_id) REFERENCES post(post_id)
);

CREATE TABLE comment(
                        comment_id INT AUTO_INCREMENT NOT NULL,
                        comment_date DATETIME NOT NULL ,
                        comment_content VARCHAR(255),
                        comment_id_1 INT NOT NULL,
                        post_id INT NOT NULL,
                        user_id INT NOT NULL,
                        PRIMARY KEY(comment_id),
                        FOREIGN KEY(comment_id_1) REFERENCES comment(comment_id),
                        FOREIGN KEY(post_id) REFERENCES post(post_id)
);

CREATE TABLE add_tag(
                        feed_id INT AUTO_INCREMENT NOT NULL,
                        tag_id INT NOT NULL,
                        PRIMARY KEY(feed_id, tag_id),
                        FOREIGN KEY(feed_id) REFERENCES feed(feed_id),
                        FOREIGN KEY(tag_id) REFERENCES tag(tag_id)
);

CREATE TABLE reaction(
                         user_id INT NOT NULL ,
                         post_id INT NOT NULL,
                         reaction BOOLEAN NOT NULL,
                         PRIMARY KEY(user_id, post_id),
                         FOREIGN KEY(user_id) REFERENCES account(user_id),
                         FOREIGN KEY(post_id) REFERENCES post(post_id)
);


CREATE TABLE reaction_comment(
                                 user_id INT NOT NULL ,
                                 comment_id INT NOT NULL,
                                 reaction BOOLEAN NOT NULL,
                                 PRIMARY KEY(user_id, comment_id),
                                 FOREIGN KEY(user_id) REFERENCES account(user_id),
                                 FOREIGN KEY(comment_id) REFERENCES comment(comment_id)
);
;