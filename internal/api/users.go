package api

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4"
	"github.com/u2lentaru/billbck/internal/models"
)

// HandleUsers godoc
// @Summary Users list
// @Description get user list
// @Tags users
// @Produce  json
// @Param page query int false "page number"
// @Param page_size query int false "page size"
// @Param username query string false "username search pattern"
// @Param ordering query string false "order by {id|username}"
// @Param desc query boolean false "descending order {true|false}"
// @Success 200 {array} models.User_count
// @Failure 500
// @Router /users [get]
func (s *APG) HandleUsers(w http.ResponseWriter, r *http.Request) {
	gs := models.User{}
	ctx := context.Background()
	out_arr := []models.User{}

	query := r.URL.Query()

	pg := 1
	spg, ok := query["page"]

	if ok && len(spg) > 0 {
		if pgt, err := strconv.Atoi(spg[0]); err != nil {
			pg = 1
		} else {
			pg = pgt
		}
	}

	pgs := 20
	spgs, ok := query["page_size"]
	if ok && len(spgs) > 0 {
		if pgst, err := strconv.Atoi(spgs[0]); err != nil {
			pgs = 20
		} else {
			pgs = pgst
		}
	}

	gs1 := ""
	gs1s, ok := query["username"]
	if ok && len(gs1s) > 0 {
		//case insensitive
		gs1 = strings.ToUpper(gs1s[0])
		//quotes
		re := regexp.MustCompile(`'`)
		gs1 = string(re.ReplaceAll([]byte(gs1), []byte("''")))
	}

	ord := 1
	ords, ok := query["ordering"]
	if !ok || len(ords) == 0 {
		ord = 1
	} else if ords[0] == "username" {
		ord = 2
	}

	dsc := false
	dscs, ok := query["desc"]
	if ok && len(dscs) > 0 {
		if !(dscs[0] == "0") {
			dsc = true
		}
	}

	rows, err := s.Dbpool.Query(ctx, "SELECT * from func_users_get($1,$2,$3,$4,$5);", pg, pgs, gs1, ord, dsc)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&gs.Id, &gs.UserName, &gs.OrgInfo.Id, &gs.Lang.Id, &gs.ChangePass, &gs.Position.Id, &gs.UserFullName,
			&gs.Created, &gs.Closed, &gs.OrgInfo.OIName, &gs.Lang.LangName, &gs.Position.PositionName)
		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, gs)
	}

	gsc := 0
	err = s.Dbpool.QueryRow(ctx, "SELECT * from func_users_cnt($1);", gs1).Scan(&gsc)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}
	out_count, err := json.Marshal(models.User_count{Values: out_arr, Count: gsc, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(out_count)

	return
}

// HandleAddUser godoc
// @Summary Add user
// @Description add user
// @Tags users
// @Accept json
// @Produce  json
// @Param a body models.AddUser true "New user. Significant params: UserName, OrgInfo.Id, Lang.Id, ChangePass, Position.Id, UserFullName, Created"
// @Success 200 {array} models.Json_id
// @Failure 500
// @Router /users_add [post]
func (s *APG) HandleAddUser(w http.ResponseWriter, r *http.Request) {
	a := models.AddUser{}
	body, err := ioutil.ReadAll(r.Body)

	defer r.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = json.Unmarshal(body, &a)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	ai := 0
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_users_add($1,$2,$3,$4,$5,$6,$7);", a.UserName, a.OrgInfo.Id, a.Lang.Id,
		a.ChangePass, a.Position.Id, a.UserFullName, a.Created).Scan(&ai)

	if err != nil {
		log.Println("Failed execute func_users_add: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ai})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleUpdUser godoc
// @Summary Update user
// @Description update user
// @Tags users
// @Accept json
// @Produce  json
// @Param u body models.User true "Update user. Significant params: Id, UserName, OrgInfo.Id, Lang.Id, ChangePass, Position.Id, UserFullName, Created, Closed(n)"
// @Success 200 {array} models.Json_id
// @Failure 500
// @Router /users_upd [post]
func (s *APG) HandleUpdUser(w http.ResponseWriter, r *http.Request) {
	u := models.User{}
	body, err := ioutil.ReadAll(r.Body)

	defer r.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = json.Unmarshal(body, &u)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	ui := 0
	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_users_upd($1,$2,$3,$4,$5,$6,$7,$8,$9);", u.Id, u.UserName, u.OrgInfo.Id,
		u.Lang.Id, u.ChangePass, u.Position.Id, u.UserFullName, u.Created, u.Closed).Scan(&ui)

	if err != nil {
		log.Println("Failed execute func_users_upd: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ui})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleDelUser godoc
// @Summary Delete users
// @Description delete users
// @Tags users
// @Accept json
// @Produce  json
// @Param d body models.Json_ids true "Delete users"
// @Success 200 {array} models.Json_ids
// @Failure 500
// @Router /users_del [post]
func (s *APG) HandleDelUser(w http.ResponseWriter, r *http.Request) {
	d := models.Json_ids{}
	body, err := ioutil.ReadAll(r.Body)

	defer r.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = json.Unmarshal(body, &d)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	res := []int{}
	i := 0
	for _, id := range d.Ids {
		err = s.Dbpool.QueryRow(context.Background(), "SELECT func_users_del($1);", id).Scan(&i)
		res = append(res, i)

		if err != nil {
			log.Println("Failed execute func_users_del: ", err)
		}
	}

	output, err := json.Marshal(models.Json_ids{Ids: res})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleGetUser godoc
// @Summary Get user
// @Description get user
// @Tags users
// @Produce  json
// @Param id path int true "User by id"
// @Success 200 {array} models.User_count
// @Failure 500
// @Router /users/{id} [get]
func (s *APG) HandleGetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	i := vars["id"]
	g := models.User{}
	out_arr := []models.User{}

	err := s.Dbpool.QueryRow(context.Background(), "SELECT * from func_user_get($1);", i).Scan(&g.Id, &g.UserName, &g.OrgInfo.Id, &g.Lang.Id,
		&g.ChangePass, &g.Position.Id, &g.UserFullName, &g.Created, &g.Closed, &g.OrgInfo.OIName, &g.Lang.LangName, &g.Position.PositionName)

	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed execute from func_user_get: ", err)
	}

	out_arr = append(out_arr, g)
	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}

	// output, err := json.Marshal(g)
	out_count, err := json.Marshal(models.User_count{Values: out_arr, Count: 1, Auth: auth})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(out_count)

	return
}
