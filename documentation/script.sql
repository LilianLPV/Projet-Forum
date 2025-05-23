CREATE TABLE account (
    user_id INT PRIMARY KEY,
    user_name VARCHAR(255) NOT NULL,
    user_password VARCHAR(512) NOT NULL,
    user_email VARCHAR(255) NOT NULL,
    user_profile_picture VARCHAR(255),
    user_bio VARCHAR(255),
    user_last_connexion DATETIME NOT NULL,
);


CREATE TABLE feed (
    feed_id VARCHAR(50) PRIMARY KEY,
    feed_title VARCHAR(100) NOT NULL,
    feed_description VARCHAR(255) NOT NULL,
    feed_status VARCHAR(50) NOT NULL,
    feed_creation_date DATETIME NOT NULL,
    user_id INT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES account(user_id)
);

CREATE TABLE tag (
    tag_id INT PRIMARY KEY,
    tag_name VARCHAR(50) NOT NULL,
    tag_type VARCHAR(50) NOT NULL
);


CREATE TABLE post (
    post_id INT PRIMARY KEY,
    post_date DATETIME NOT NULL,
    post_content VARCHAR(255) NOT NULL,
    feed_id VARCHAR(50) NOT NULL,
    user_id INT NOT NULL,
    FOREIGN KEY (feed_id) REFERENCES feed(feed_id),
    FOREIGN KEY (user_id) REFERENCES account(user_id)
);


CREATE TABLE picture (
    picture_id INT PRIMARY KEY,
    picture_blob VARCHAR(255),
    post_id INT NOT NULL,
    FOREIGN KEY (post_id) REFERENCES post(post_id)
);


CREATE TABLE comment (
    comment_id VARCHAR(50) PRIMARY KEY,
    comment_date DATETIME NOT NULL,
    comment_content VARCHAR(255) NOT NULL,
    parent_comment_id VARCHAR(50) NOT NULL,
    post_id INT NOT NULL,
    FOREIGN KEY (parent_comment_id) REFERENCES comment(comment_id),
    FOREIGN KEY (post_id) REFERENCES post(post_id)
);


CREATE TABLE add_tag (
    feed_id VARCHAR(50) NOT NULL,
    tag_id INT NOT NULL,
    PRIMARY KEY (feed_id, tag_id),
    FOREIGN KEY (feed_id) REFERENCES feed(feed_id),
    FOREIGN KEY (tag_id) REFERENCES tag(tag_id)
);


CREATE TABLE reaction (
    user_id INT NOT NULL,
    post_id INT NOT NULL,
    PRIMARY KEY (user_id, post_id),
    FOREIGN KEY (user_id) REFERENCES account(user_id),
    FOREIGN KEY (post_id) REFERENCES post(post_id)
);