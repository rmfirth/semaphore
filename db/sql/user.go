package sql

import (
	"database/sql"
	"github.com/ansible-semaphore/semaphore/db"
	"github.com/masterminds/squirrel"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func (d *SqlDb) CreateUserWithoutPassword(user db.User) (newUser db.User, err error) {

	err = db.ValidateUser(user)
	if err != nil {
		return
	}

	user.Password = ""
	user.Created = db.GetParsedTime(time.Now())

	err = d.sql.Insert(&user)

	if err != nil {
		return
	}

	newUser = user
	return
}

func (d *SqlDb) CreateUser(user db.UserWithPwd) (newUser db.User, err error) {

	err = db.ValidateUser(user.User)
	if err != nil {
		return
	}

	pwdHash, err := bcrypt.GenerateFromPassword([]byte(user.Pwd), 11)

	if err != nil {
		return
	}

	user.Password = string(pwdHash)
	user.Created = db.GetParsedTime(time.Now())

	err = d.sql.Insert(&user.User)

	if err != nil {
		return
	}

	newUser = user.User
	return
}

func (d *SqlDb) DeleteUser(userID int) error {
	res, err := d.exec("delete from `user` where id=?", userID)
	return validateMutationResult(res, err)
}

func (d *SqlDb) UpdateUser(user db.UserWithPwd) error {
	var err error

	if user.Pwd != "" {
		var pwdHash []byte
		pwdHash, err = bcrypt.GenerateFromPassword([]byte(user.Pwd), 11)
		if err != nil {
			return err
		}
		_, err = d.exec(
			"update user set name=?, username=?, email=?, alert=?, admin=?, password=? where id=?",
			user.Name,
			user.Username,
			user.Email,
			user.Alert,
			user.Admin,
			string(pwdHash),
			user.ID)
	} else {
		_, err = d.exec(
			"update `user` set name=?, username=?, email=?, alert=?, admin=? where id=?",
			user.Name,
			user.Username,
			user.Email,
			user.Alert,
			user.Admin,
			user.ID)
	}

	return err
}

func (d *SqlDb) SetUserPassword(userID int, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 11)
	if err != nil {
		return err
	}
	_, err = d.exec(
		"update `user` set password=? where id=?",
		string(hash), userID)
	return err
}

func (d *SqlDb) CreateProjectUser(projectUser db.ProjectUser) (newProjectUser db.ProjectUser, err error) {
	_, err = d.exec(
		"insert into project__user (project_id, user_id, `admin`) values (?, ?, ?)",
		projectUser.ProjectID,
		projectUser.UserID,
		projectUser.Admin)

	if err != nil {
		return
	}

	newProjectUser = projectUser
	return
}

func (d *SqlDb) GetProjectUser(projectID, userID int) (db.ProjectUser, error) {
	var user db.ProjectUser

	err := d.selectOne(&user,
		"select * from project__user where project_id=? and user_id=?",
		projectID,
		userID)

	if err == sql.ErrNoRows {
		err = db.ErrNotFound
	}

	return user, err
}

func (d *SqlDb) GetProjectUsers(projectID int, params db.RetrieveQueryParams) (users []db.User, err error) {
	q := squirrel.Select("u.*").Column("pu.admin").
		From("project__user as pu").
		LeftJoin("`user` as u on pu.user_id=u.id").
		Where("pu.project_id=?", projectID)

	sortDirection := "ASC"
	if params.SortInverted {
		sortDirection = "DESC"
	}

	switch params.SortBy {
	case "name", "username", "email":
		q = q.OrderBy("u." + params.SortBy + " " + sortDirection)
	case "admin":
		q = q.OrderBy("pu." + params.SortBy + " " + sortDirection)
	default:
		q = q.OrderBy("u.name " + sortDirection)
	}

	query, args, err := q.ToSql()

	if err != nil {
		return
	}

	_, err = d.selectAll(&users, query, args...)

	return
}

func (d *SqlDb) UpdateProjectUser(projectUser db.ProjectUser) error {
	_, err := d.exec(
		"update `project__user` set admin=? where user_id=? and project_id = ?",
		projectUser.Admin,
		projectUser.UserID,
		projectUser.ProjectID)

	return err
}

func (d *SqlDb) DeleteProjectUser(projectID, userID int) error {
	_, err := d.exec("delete from project__user where user_id=? and project_id=?", userID, projectID)
	return err
}

//GetUser retrieves a user from the database by ID
func (d *SqlDb) GetUser(userID int) (db.User, error) {
	var user db.User

	err := d.selectOne(&user, "select * from `user` where id=?", userID)

	if err == sql.ErrNoRows {
		err = db.ErrNotFound
	}

	return user, err
}

func (d *SqlDb) GetUsers(params db.RetrieveQueryParams) (users []db.User, err error) {
	query, args, err := getSqlForTable("user", params)

	if err != nil {
		return
	}

	_, err = d.selectAll(&users, query, args...)

	return
}

func (d *SqlDb) GetUserByLoginOrEmail(login string, email string) (existingUser db.User, err error) {
	err = d.selectOne(
		&existingUser,
		d.prepareQuery("select * from `user` where email=? or username=?"),
		email, login)

	if err == sql.ErrNoRows {
		err = db.ErrNotFound
	}

	return
}

func (d *SqlDb) CreatePlaceholderUser() error {
	user := db.User{
		Created: db.GetParsedTime(time.Now()),
	}
	return d.sql.Insert(&user)
}

func (d *SqlDb) GetPlaceholderUser() (user db.User, err error) {
	err = d.selectOne(&user, d.prepareQuery("select id from `user` where username=''"))

	if err == sql.ErrNoRows {
		err = db.ErrNotFound
	}

	return
}

//
//func (d *SqlDb) HasPlaceholderUser() (bool, error) {
//	_, err := d.getPlaceholderUser()
//
//	if err == nil {
//		return true, nil
//	}
//
//	if err == db.ErrNotFound {
//		return false, nil
//	}
//
//	return false, err
//}
//
//func (d *SqlDb) ReplacePlaceholderUser(user db.UserWithPwd) (newUser db.User, err error) {
//	placeholder, err := d.getPlaceholderUser()
//	if err != nil {
//		return
//	}
//	user.ID = placeholder.ID
//	err = d.UpdateUser(user)
//	if err != nil {
//		return
//	}
//	newUser = user.User
//	return
//}
