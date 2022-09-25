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
	"github.com/u2lentaru/billbck/internal/adapters/db/pgsql"
	"github.com/u2lentaru/billbck/internal/models"
	"github.com/u2lentaru/billbck/internal/services"
	"github.com/u2lentaru/billbck/internal/utils"
)

type ifContractService interface {
	GetList(ctx context.Context, pg, pgs int, gs1, gs2, gs3, gs4 string, gs5, gs6 int, gs7, gs8, gs9, gs10, gs11, gs12, gs13, gs14 string, ord int, dsc bool) (models.Contract_count, error)
	Add(ctx context.Context, ea models.Contract) (int, error)
	Upd(ctx context.Context, eu models.Contract) (int, error)
	Del(ctx context.Context, ed models.IdClose) (int, error)
	GetOne(ctx context.Context, i int) (models.Contract_count, error)
	GetObj(ctx context.Context, i int, a string) (models.ObjContract_count, error)
	GetHist(ctx context.Context, i int) (string, error)
}

// HandleContracts godoc
// @Summary List contracts
// @Description get contract list
// @Tags contracts
// @Produce  json
// @Param page query int false "page number"
// @Param page_size query int false "page size"
// @Param contractnumber query string false "contractnumber search pattern"
// @Param personalaccount query string false "personalaccount search pattern"
// @Param active query boolean false "active contracts"
// @Param custname query string false "customer name search pattern"
// @Param custid query int false "customer id"
// @Param custtype query int false "customer type"
// @Param custaddr query string false "customer address search pattern"
// @Param consname query string false "сonsignee name search pattern"
// @Param esoname query string false "eso name search pattern"
// @Param esocn query string false "eso contract number search pattern"
// @Param areaname query string false "area name search pattern"
// @Param cgname query string false "customer group name search pattern"
// @Param motru query string false "motive of termination ru search pattern"
// @Param motkz query string false "motive of termination kz search pattern"
// @Param ordering query string false "order by {contractnumber|personalaccount|startdate|custname|esoname|enddate}"
// @Param desc query boolean false "descending order {true|false}"
// @Success 200 {object} models.Contract_count
// @Failure 500
// @Router /contracts [get]
func HandleContracts(w http.ResponseWriter, r *http.Request) {
	var gs ifContractService
	gs = services.NewContractService(pgsql.ContractStorage{})
	ctx := context.Background()
	auth := utils.GetAuth(r)

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
	gs1s, ok := query["contractnumber"]
	if ok && len(gs1s) > 0 {
		//case insensitive
		gs1 = strings.ToUpper(gs1s[0])
		//quotes
		re := regexp.MustCompile(`'`)
		gs1 = string(re.ReplaceAll([]byte(gs1), []byte("''")))
	}

	gs2 := ""
	gs2s, ok := query["personalaccount"]
	if ok && len(gs2s) > 0 {
		//case insensitive
		gs2 = strings.ToUpper(gs2s[0])
		//quotes
		re := regexp.MustCompile(`'`)
		gs2 = string(re.ReplaceAll([]byte(gs2), []byte("''")))
	}

	gs3 := ""
	gs3s, ok := query["active"]
	if ok && len(gs3s) > 0 {
		if gs3s[0] == "true" || gs3s[0] == "false" {
			gs3 = gs3s[0]
		} else {
			gs3 = ""
		}
	}

	gs4 := ""
	gs4s, ok := query["custname"]
	if ok && len(gs4s) > 0 {
		//case insensitive
		gs4 = strings.ToUpper(gs4s[0])
		//quotes
		re := regexp.MustCompile(`'`)
		gs4 = string(re.ReplaceAll([]byte(gs4), []byte("''")))
	}

	gs5 := 0
	gs5s, ok := query["custid"]
	if ok && len(gs5s) > 0 {
		t, err := strconv.Atoi(gs5s[0])
		if err == nil {
			gs5 = t
		}
	}

	gs6 := 0
	gs6s, ok := query["custtype"]
	if ok && len(gs6s) > 0 {
		t, err := strconv.Atoi(gs6s[0])
		if err == nil {
			gs6 = t
		}
	}

	gs7 := ""
	gs7s, ok := query["custaddr"]
	if ok && len(gs7s) > 0 {
		//case insensitive
		gs7 = strings.ToUpper(gs7s[0])
		//quotes
		re := regexp.MustCompile(`'`)
		gs7 = string(re.ReplaceAll([]byte(gs7), []byte("''")))
	}

	gs8 := ""
	gs8s, ok := query["consname"]
	if ok && len(gs8s) > 0 {
		//case insensitive
		gs8 = strings.ToUpper(gs8s[0])
		//quotes
		re := regexp.MustCompile(`'`)
		gs8 = string(re.ReplaceAll([]byte(gs8), []byte("''")))
	}

	gs9 := ""
	gs9s, ok := query["esoname"]
	if ok && len(gs9s) > 0 {
		//case insensitive
		gs9 = strings.ToUpper(gs9s[0])
		//quotes
		re := regexp.MustCompile(`'`)
		gs9 = string(re.ReplaceAll([]byte(gs9), []byte("''")))
	}

	gs10 := ""
	gs10s, ok := query["esocn"]
	if ok && len(gs10s) > 0 {
		//case insensitive
		gs10 = strings.ToUpper(gs10s[0])
		//quotes
		re := regexp.MustCompile(`'`)
		gs10 = string(re.ReplaceAll([]byte(gs10), []byte("''")))
	}

	gs11 := ""
	gs11s, ok := query["areaname"]
	if ok && len(gs11s) > 0 {
		//case insensitive
		gs11 = strings.ToUpper(gs11s[0])
		//quotes
		re := regexp.MustCompile(`'`)
		gs11 = string(re.ReplaceAll([]byte(gs11), []byte("''")))
	}

	gs12 := ""
	gs12s, ok := query["cgname"]
	if ok && len(gs12s) > 0 {
		//case insensitive
		gs12 = strings.ToUpper(gs12s[0])
		//quotes
		re := regexp.MustCompile(`'`)
		gs12 = string(re.ReplaceAll([]byte(gs12), []byte("''")))
	}

	gs13 := ""
	gs13s, ok := query["motru"]
	if ok && len(gs13s) > 0 {
		//case insensitive
		gs13 = strings.ToUpper(gs13s[0])
		//quotes
		re := regexp.MustCompile(`'`)
		gs13 = string(re.ReplaceAll([]byte(gs13), []byte("''")))
	}

	gs14 := ""
	gs14s, ok := query["motkz"]
	if ok && len(gs14s) > 0 {
		//case insensitive
		gs14 = strings.ToUpper(gs14s[0])
		//quotes
		re := regexp.MustCompile(`'`)
		gs14 = string(re.ReplaceAll([]byte(gs14), []byte("''")))
	}

	ord := 1
	ords, ok := query["ordering"]
	if !ok || len(ords) == 0 {
		ord = 1
	} else if ords[0] == "contractnumber" {
		ord = 5
	} else if ords[0] == "personalaccount" {
		ord = 3
	} else if ords[0] == "startdate" {
		ord = 6
	} else if ords[0] == "custname" {
		ord = 15
	} else if ords[0] == "esoname" {
		ord = 18
	} else if ords[0] == "enddate" {
		ord = 7
	}

	dsc := false
	dscs, ok := query["desc"]
	if ok && len(dscs) > 0 {
		if !(dscs[0] == "0") {
			dsc = true
		}
	}

	out_arr, err := gs.GetList(ctx, pg, pgs, gs1, gs2, gs3, gs4, gs5, gs6, gs7, gs8, gs9, gs10, gs11, gs12, gs13, gs14, ord, dsc)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	out_arr.Auth = auth
	out_count, err := json.Marshal(out_arr)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(out_count)

	return
}

// HandleAddContract godoc
// @Summary Add contract
// @Description add contract
// @Tags contracts
// @Accept json
// @Produce  json
// @Param a body models.AddContract true "New contract. Significant params: PersonalAccount, BarCode, ContractNumber, Startdate, Customer.SubId, Consignee.SubId, EsoContractNumber, Eso.Id, Area.Id, CustomerGroup.Id, Notes(n)"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /contracts_add [post]
func HandleAddContract(w http.ResponseWriter, r *http.Request) {
	var gs ifContractService
	gs = services.NewContractService(pgsql.ContractStorage{})
	ctx := context.Background()

	a := models.Contract{}
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

	ai, err := gs.Add(ctx, a)

	if err != nil {
		log.Println("Failed execute ifContractService.Add: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ai})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return

}

// HandleUpdContract godoc
// @Summary Update contract
// @Description update contract
// @Tags contracts
// @Accept json
// @Produce  json
// @Param u body models.Contract true "Update contract. Significant params: Id, ContractNumber, Startdate, Enddate(n), Customer.SubId, Consignee.SubId, EsoContractNumber, Eso.Id, Area.Id, CustomerGroup.Id, ContractMot.Id(n), Notes(n), MotNotes(n)"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /contracts_upd [post]
func HandleUpdContract(w http.ResponseWriter, r *http.Request) {
	var gs ifContractService
	gs = services.NewContractService(pgsql.ContractStorage{})
	ctx := context.Background()

	u := models.Contract{}
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

	ui, err := gs.Upd(ctx, u)

	if err != nil {
		log.Println("Failed execute ifContractService.Upd: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: ui})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleDelContract godoc
// @Summary Close contracts
// @Description close contracts
// @Tags contracts
// @Accept json
// @Produce  json
// @Param d body models.IdClose true "Close contract. Significant params: Id, Enddate, ContractMot.Id, MotNotes(n)"
// @Success 200 {object} models.Json_id
// @Failure 500
// @Router /contracts_del [post]
func HandleDelContract(w http.ResponseWriter, r *http.Request) {
	var gs ifContractService
	gs = services.NewContractService(pgsql.ContractStorage{})
	ctx := context.Background()

	d := models.IdClose{}
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

	res, err := gs.Del(ctx, d)
	if err != nil {
		log.Println("Failed execute ifContractService.Del: ", err)
	}

	output, err := json.Marshal(models.Json_id{Id: res})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)

	return
}

// HandleGetContract godoc
// @Summary Get contract
// @Description get contract
// @Tags contracts
// @Produce  json
// @Param id path int true "Contract by id"
// @Success 200 {object} models.Contract_count
// @Failure 500
// @Router /contracts/{id} [get]
func HandleGetContract(w http.ResponseWriter, r *http.Request) {
	var gs ifContractService
	gs = services.NewContractService(pgsql.ContractStorage{})
	ctx := context.Background()
	auth := utils.GetAuth(r)

	vars := mux.Vars(r)
	i, err := strconv.Atoi(vars["id"])
	if err != nil {
		i = 0
	}

	out_arr, err := gs.GetOne(ctx, i)
	if err != nil {
		log.Println("Failed execute ifContractService.GetOne: ", err)
	}

	out_arr.Auth = auth
	out_count, err := json.Marshal(out_arr)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(out_count)

	return
}

// HandleGetContractObject godoc
// @Summary Get contracts objects
// @Description get contracts objects
// @Tags contracts
// @Produce  json
// @Param id path int true "Contract by id"
// @Param active query boolean false "active objects"
// @Success 200 {object} models.ObjContract
// @Failure 500
// @Router /contracts_getobject/{id} [get]
func HandleGetContractObject(w http.ResponseWriter, r *http.Request) {
	var gs ifContractService
	gs = services.NewContractService(pgsql.ContractStorage{})
	ctx := context.Background()
	vars := mux.Vars(r)
	i, err := strconv.Atoi(vars["id"])
	if err != nil {
		i = 0
	}

	out_arr := models.ObjContract_count{}

	query := r.URL.Query()

	dsc := ""
	dscs, ok := query["active"]
	if ok && len(dscs) > 0 {
		if dscs[0] == "true" || dscs[0] == "false" {
			dsc = dscs[0]
		} else {
			dsc = ""
		}
	}

	out_arr, err = gs.GetObj(ctx, i, dsc)
	if err != nil {
		log.Println("Failed execute ifContractService.GetOne: ", err)
	}

	// out_count, err := json.Marshal(out_arr)

	out_count, err := json.Marshal(out_arr)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(out_count)

	return
}

// HandleGetContractHist godoc
// @Summary Contract history
// @Description get contract history
// @Tags contracts
// @Produce  json
// @Param id path int true "Contract history by id"
// @Success 200 {object} string
// @Failure 500
// @Router /contracts_hist/{id} [get]
func HandleGetContractHist(w http.ResponseWriter, r *http.Request) {
	var gs ifContractService
	gs = services.NewContractService(pgsql.ContractStorage{})
	ctx := context.Background()

	vars := mux.Vars(r)
	i, err := strconv.Atoi(vars["id"])
	if err != nil {
		i = 0
	}

	hs, err := gs.GetHist(ctx, i)
	if err != nil {
		log.Println("Failed execute ifContractService.GetOne: ", err)
	}

	w.Write([]byte(hs))

	return
}

/////////////////////////////////////////////////////////////////////////////////////////////
// package api

// import (
// 	"context"
// 	"database/sql"
// 	"encoding/json"
// 	"io/ioutil"
// 	"log"
// 	"net/http"
// 	"regexp"
// 	"strconv"
// 	"strings"
// 	"time"

// 	"github.com/gorilla/mux"
// 	"github.com/jackc/pgx/v4"
// 	"github.com/u2lentaru/billbck/internal/models"
// 	"github.com/u2lentaru/billbck/internal/utils"
// )

// // HandleContracts godoc
// // @Summary List contracts
// // @Description get contract list
// // @Tags contracts
// // @Produce  json
// // @Param page query int false "page number"
// // @Param page_size query int false "page size"
// // @Param contractnumber query string false "contractnumber search pattern"
// // @Param personalaccount query string false "personalaccount search pattern"
// // @Param active query boolean false "active contracts"
// // @Param custname query string false "customer name search pattern"
// // @Param custid query int false "customer id"
// // @Param custtype query int false "customer type"
// // @Param custaddr query string false "customer address search pattern"
// // @Param consname query string false "сonsignee name search pattern"
// // @Param esoname query string false "eso name search pattern"
// // @Param esocn query string false "eso contract number search pattern"
// // @Param areaname query string false "area name search pattern"
// // @Param cgname query string false "customer group name search pattern"
// // @Param motru query string false "motive of termination ru search pattern"
// // @Param motkz query string false "motive of termination kz search pattern"
// // @Param ordering query string false "order by {contractnumber|personalaccount|startdate|custname|esoname|enddate}"
// // @Param desc query boolean false "descending order {true|false}"
// // @Success 200 {object} models.Contract_count
// // @Failure 500
// // @Router /contracts [get]
// func (s *APG) HandleContracts(w http.ResponseWriter, r *http.Request) {
// 	gs := models.Contract{}
// 	ctx := context.Background()

// 	query := r.URL.Query()

// 	pg := 1
// 	spg, ok := query["page"]

// 	if ok && len(spg) > 0 {
// 		if pgt, err := strconv.Atoi(spg[0]); err != nil {
// 			pg = 1
// 		} else {
// 			pg = pgt
// 		}
// 	}

// 	pgs := 20
// 	spgs, ok := query["page_size"]
// 	if ok && len(spgs) > 0 {
// 		if pgst, err := strconv.Atoi(spgs[0]); err != nil {
// 			pgs = 20
// 		} else {
// 			pgs = pgst
// 		}
// 	}

// 	gs1 := ""
// 	gs1s, ok := query["contractnumber"]
// 	if ok && len(gs1s) > 0 {
// 		//case insensitive
// 		gs1 = strings.ToUpper(gs1s[0])
// 		//quotes
// 		re := regexp.MustCompile(`'`)
// 		gs1 = string(re.ReplaceAll([]byte(gs1), []byte("''")))
// 	}

// 	gs2 := ""
// 	gs2s, ok := query["personalaccount"]
// 	if ok && len(gs2s) > 0 {
// 		//case insensitive
// 		gs2 = strings.ToUpper(gs2s[0])
// 		//quotes
// 		re := regexp.MustCompile(`'`)
// 		gs2 = string(re.ReplaceAll([]byte(gs2), []byte("''")))
// 	}

// 	// gs3 := time.Now().Format("2006-01-02")
// 	// // log.Println(gs3)
// 	// gs3s, ok := query["actualdate"]
// 	// if ok && len(gs3s) > 0 {
// 	// 	//case insensitive
// 	// 	gs3 = strings.ToUpper(gs3s[0])
// 	// 	//quotes
// 	// 	re := regexp.MustCompile(`'`)
// 	// 	gs3 = string(re.ReplaceAll([]byte(gs3), []byte("''")))
// 	// }

// 	gs3 := ""
// 	gs3s, ok := query["active"]
// 	if ok && len(gs3s) > 0 {
// 		if gs3s[0] == "true" || gs3s[0] == "false" {
// 			gs3 = gs3s[0]
// 		} else {
// 			gs3 = ""
// 		}
// 	}

// 	gs4 := ""
// 	gs4s, ok := query["custname"]
// 	if ok && len(gs4s) > 0 {
// 		//case insensitive
// 		gs4 = strings.ToUpper(gs4s[0])
// 		//quotes
// 		re := regexp.MustCompile(`'`)
// 		gs4 = string(re.ReplaceAll([]byte(gs4), []byte("''")))
// 	}

// 	gs5 := 0
// 	gs5s, ok := query["custid"]
// 	if ok && len(gs5s) > 0 {
// 		t, err := strconv.Atoi(gs5s[0])
// 		if err == nil {
// 			gs5 = t
// 		}
// 	}

// 	gs6 := 0
// 	gs6s, ok := query["custtype"]
// 	if ok && len(gs6s) > 0 {
// 		t, err := strconv.Atoi(gs6s[0])
// 		if err == nil {
// 			gs6 = t
// 		}
// 	}

// 	gs7 := ""
// 	gs7s, ok := query["custaddr"]
// 	if ok && len(gs7s) > 0 {
// 		//case insensitive
// 		gs7 = strings.ToUpper(gs7s[0])
// 		//quotes
// 		re := regexp.MustCompile(`'`)
// 		gs7 = string(re.ReplaceAll([]byte(gs7), []byte("''")))
// 	}

// 	gs8 := ""
// 	gs8s, ok := query["consname"]
// 	if ok && len(gs8s) > 0 {
// 		//case insensitive
// 		gs8 = strings.ToUpper(gs8s[0])
// 		//quotes
// 		re := regexp.MustCompile(`'`)
// 		gs8 = string(re.ReplaceAll([]byte(gs8), []byte("''")))
// 	}

// 	gs9 := ""
// 	gs9s, ok := query["esoname"]
// 	if ok && len(gs9s) > 0 {
// 		//case insensitive
// 		gs9 = strings.ToUpper(gs9s[0])
// 		//quotes
// 		re := regexp.MustCompile(`'`)
// 		gs9 = string(re.ReplaceAll([]byte(gs9), []byte("''")))
// 	}

// 	gs10 := ""
// 	gs10s, ok := query["esocn"]
// 	if ok && len(gs10s) > 0 {
// 		//case insensitive
// 		gs10 = strings.ToUpper(gs10s[0])
// 		//quotes
// 		re := regexp.MustCompile(`'`)
// 		gs10 = string(re.ReplaceAll([]byte(gs10), []byte("''")))
// 	}

// 	gs11 := ""
// 	gs11s, ok := query["areaname"]
// 	if ok && len(gs11s) > 0 {
// 		//case insensitive
// 		gs11 = strings.ToUpper(gs11s[0])
// 		//quotes
// 		re := regexp.MustCompile(`'`)
// 		gs11 = string(re.ReplaceAll([]byte(gs11), []byte("''")))
// 	}

// 	gs12 := ""
// 	gs12s, ok := query["cgname"]
// 	if ok && len(gs12s) > 0 {
// 		//case insensitive
// 		gs12 = strings.ToUpper(gs12s[0])
// 		//quotes
// 		re := regexp.MustCompile(`'`)
// 		gs12 = string(re.ReplaceAll([]byte(gs12), []byte("''")))
// 	}

// 	gs13 := ""
// 	gs13s, ok := query["motru"]
// 	if ok && len(gs13s) > 0 {
// 		//case insensitive
// 		gs13 = strings.ToUpper(gs13s[0])
// 		//quotes
// 		re := regexp.MustCompile(`'`)
// 		gs13 = string(re.ReplaceAll([]byte(gs13), []byte("''")))
// 	}

// 	gs14 := ""
// 	gs14s, ok := query["motkz"]
// 	if ok && len(gs14s) > 0 {
// 		//case insensitive
// 		gs14 = strings.ToUpper(gs14s[0])
// 		//quotes
// 		re := regexp.MustCompile(`'`)
// 		gs14 = string(re.ReplaceAll([]byte(gs14), []byte("''")))
// 	}

// 	gsc := 0
// 	err := s.Dbpool.QueryRow(ctx, "SELECT * from func_contracts_cnt($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14);", gs1, gs2, utils.NullableString(gs3),
// 		gs4, utils.NullableInt(int32(gs5)), utils.NullableInt(int32(gs6)), gs7, gs8, gs9, gs10, gs11, gs12, utils.NullableString(gs13),
// 		utils.NullableString(gs14)).Scan(&gsc)

// 	if err != nil {
// 		http.Error(w, err.Error(), 500)
// 		return
// 	}

// 	out_arr := make([]models.Contract, 0,
// 		func() int {
// 			if gsc < pgs {
// 				return gsc
// 			} else {
// 				return pgs
// 			}
// 		}())

// 	ord := 1
// 	ords, ok := query["ordering"]
// 	if !ok || len(ords) == 0 {
// 		ord = 1
// 	} else if ords[0] == "contractnumber" {
// 		ord = 5
// 	} else if ords[0] == "personalaccount" {
// 		ord = 3
// 	} else if ords[0] == "startdate" {
// 		ord = 6
// 	} else if ords[0] == "custname" {
// 		ord = 15
// 	} else if ords[0] == "esoname" {
// 		ord = 18
// 	} else if ords[0] == "enddate" {
// 		ord = 7
// 	}

// 	dsc := false
// 	dscs, ok := query["desc"]
// 	if ok && len(dscs) > 0 {
// 		if !(dscs[0] == "0") {
// 			dsc = true
// 		}
// 	}

// 	rows, err := s.Dbpool.Query(ctx, "SELECT * from func_contracts_get($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18);",
// 		pg, pgs, gs1, gs2, utils.NullableString(gs3), gs4, utils.NullableInt(int32(gs5)), utils.NullableInt(int32(gs6)), gs7, gs8, gs9,
// 		gs10, gs11, gs12, utils.NullableString(gs13), utils.NullableString(gs14), ord, dsc)
// 	if err != nil {
// 		http.Error(w, err.Error(), 500)
// 		return
// 	}

// 	defer rows.Close()

// 	var rmi sql.NullInt32
// 	var rmr, rmk sql.NullString

// 	for rows.Next() {
// 		err = rows.Scan(&gs.Id, &gs.BarCode, &gs.PersonalAccount, &gs.Id, &gs.ContractNumber, &gs.Startdate, &gs.Enddate, &gs.Customer.SubId,
// 			&gs.Consignee.SubId, &gs.EsoContractNumber, &gs.Eso.Id, &gs.Area.Id, &gs.CustomerGroup.Id, &rmi, &gs.Notes, &gs.MotNotes,
// 			&gs.Customer.SubName, &gs.Customer.SubAddr, &gs.Consignee.SubName, &gs.Eso.EsoName, &gs.Area.AreaName,
// 			&gs.CustomerGroup.CustomerGroupName, &rmr, &rmk)

// 		gs.ContractMot.Id = int(rmi.Int32)
// 		gs.ContractMot.ContractMotNameRu = rmr.String
// 		gs.ContractMot.ContractMotNameKz = rmk.String

// 		if err != nil {
// 			log.Println("failed to scan row:", err)
// 		}

// 		out_arr = append(out_arr, gs)
// 	}

// 	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}

// 	out_count, err := json.Marshal(models.Contract_count{Values: out_arr, Count: gsc, Auth: auth})
// 	if err != nil {
// 		http.Error(w, err.Error(), 500)
// 		return
// 	}

// 	w.Write(out_count)

// 	return
// }

// // HandleAddContract godoc
// // @Summary Add contract
// // @Description add contract
// // @Tags contracts
// // @Accept json
// // @Produce  json
// // @Param a body models.AddContract true "New contract. Significant params: PersonalAccount, BarCode, ContractNumber, Startdate, Customer.SubId, Consignee.SubId, EsoContractNumber, Eso.Id, Area.Id, CustomerGroup.Id, Notes(n)"
// // @Success 200 {object} models.Json_id
// // @Failure 500
// // @Router /contracts_add [post]
// func (s *APG) HandleAddContract(w http.ResponseWriter, r *http.Request) {
// 	a := models.AddContract{}
// 	body, err := ioutil.ReadAll(r.Body)

// 	defer r.Body.Close()

// 	if err != nil {
// 		http.Error(w, err.Error(), 500)
// 		return
// 	}

// 	err = json.Unmarshal(body, &a)
// 	if err != nil {
// 		http.Error(w, err.Error(), 500)
// 		return
// 	}

// 	ai := 0
// 	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_contracts_add($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11);",
// 		a.PersonalAccount, a.BarCode, a.ContractNumber, a.Startdate, a.Customer.SubId, a.Consignee.SubId, a.EsoContractNumber, a.Eso.Id,
// 		a.Area.Id, a.CustomerGroup.Id, a.Notes).Scan(&ai)

// 	if err != nil {
// 		log.Println("Failed execute func_contracts_add: ", err)
// 	}

// 	output, err := json.Marshal(models.Json_id{Id: ai})
// 	if err != nil {
// 		http.Error(w, err.Error(), 500)
// 		return
// 	}
// 	w.Write(output)

// 	return

// }

// // HandleUpdContract godoc
// // @Summary Update contract
// // @Description update contract
// // @Tags contracts
// // @Accept json
// // @Produce  json
// // @Param u body models.Contract true "Update contract. Significant params: Id, ContractNumber, Startdate, Enddate(n), Customer.SubId, Consignee.SubId, EsoContractNumber, Eso.Id, Area.Id, CustomerGroup.Id, ContractMot.Id(n), Notes(n), MotNotes(n)"
// // @Success 200 {object} models.Json_id
// // @Failure 500
// // @Router /contracts_upd [post]
// func (s *APG) HandleUpdContract(w http.ResponseWriter, r *http.Request) {
// 	u := models.Contract{}
// 	body, err := ioutil.ReadAll(r.Body)

// 	defer r.Body.Close()

// 	if err != nil {
// 		http.Error(w, err.Error(), 500)
// 		return
// 	}

// 	err = json.Unmarshal(body, &u)
// 	if err != nil {
// 		http.Error(w, err.Error(), 500)
// 		return
// 	}

// 	ui := 0
// 	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_contracts_upd($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13);", u.Id,
// 		u.ContractNumber, u.Startdate, u.Enddate, u.Customer.SubId, u.Consignee.SubId, u.EsoContractNumber, u.Eso.Id, u.Area.Id,
// 		u.CustomerGroup.Id, utils.NullableInt(int32(u.ContractMot.Id)), u.Notes, u.MotNotes).Scan(&ui)

// 	if err != nil {
// 		log.Println("Failed execute func_contracts_upd: ", err)
// 	}

// 	output, err := json.Marshal(models.Json_id{Id: ui})

// 	if err != nil {
// 		http.Error(w, err.Error(), 500)
// 		return
// 	}
// 	w.Write(output)

// 	return
// }

// // HandleDelContract godoc
// // @Summary Close contracts
// // @Description close contracts
// // @Tags contracts
// // @Accept json
// // @Produce  json
// // @Param d body models.IdClose true "Close contract. Significant params: Id, Enddate, ContractMot.Id, MotNotes(n)"
// // @Success 200 {object} models.Json_id
// // @Failure 500
// // @Router /contracts_del [post]
// func (s *APG) HandleDelContract(w http.ResponseWriter, r *http.Request) {
// 	d := models.IdClose{}
// 	body, err := ioutil.ReadAll(r.Body)

// 	// query := r.URL.Query()

// 	// gs3 := time.Now().Format("2006-01-02")
// 	// // log.Println(gs3)
// 	// gs3s, ok := query["actualdate"]
// 	// if ok && len(gs3s) > 0 {
// 	// 	//case insensitive
// 	// 	gs3 = strings.ToUpper(gs3s[0])
// 	// 	//quotes
// 	// 	re := regexp.MustCompile(`'`)
// 	// 	gs3 = string(re.ReplaceAll([]byte(gs3), []byte("''")))
// 	// }

// 	defer r.Body.Close()

// 	if err != nil {
// 		http.Error(w, err.Error(), 500)
// 		return
// 	}

// 	err = json.Unmarshal(body, &d)
// 	if err != nil {
// 		http.Error(w, err.Error(), 500)
// 		return
// 	}

// 	// res := []int{}
// 	i := 0
// 	_, err = time.Parse("2006-01-02", d.CloseDate)
// 	if err != nil {
// 		d.CloseDate = time.Now().Format("2006-01-02")
// 	}
// 	// for _, id := range d.Ids {
// 	err = s.Dbpool.QueryRow(context.Background(), "SELECT func_contracts_del($1,$2,$3,$4);", d.Id, d.CloseDate, d.ContractMot.Id,
// 		utils.NullableString(d.MotNotes)).Scan(&i)
// 	// res = append(res, i)

// 	if err != nil {
// 		log.Println("Failed execute func_contracts_del: ", err)
// 	}
// 	// }

// 	output, err := json.Marshal(models.Json_id{Id: i})
// 	if err != nil {
// 		http.Error(w, err.Error(), 500)
// 		return
// 	}
// 	w.Write(output)

// 	return
// }

// // HandleGetContract godoc
// // @Summary Get contract
// // @Description get contract
// // @Tags contracts
// // @Produce  json
// // @Param id path int true "Contract by id"
// // @Success 200 {object} models.Contract_count
// // @Failure 500
// // @Router /contracts/{id} [get]
// func (s *APG) HandleGetContract(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	i := vars["id"]
// 	g := models.Contract{}

// 	out_arr := []models.Contract{}

// 	// query := r.URL.Query()

// 	// gs3 := time.Now().Format("2006-01-02")
// 	// // log.Println(gs3)
// 	// gs3s, ok := query["actualdate"]
// 	// if ok && len(gs3s) > 0 {
// 	// 	//case insensitive
// 	// 	gs3 = strings.ToUpper(gs3s[0])
// 	// 	//quotes
// 	// 	re := regexp.MustCompile(`'`)
// 	// 	gs3 = string(re.ReplaceAll([]byte(gs3), []byte("''")))
// 	// }

// 	var rmi sql.NullInt32
// 	var rmr, rmk sql.NullString

// 	err := s.Dbpool.QueryRow(context.Background(), "SELECT * from func_contract_get($1);", i).Scan(&g.Id, &g.BarCode, &g.PersonalAccount,
// 		&g.Id, &g.ContractNumber, &g.Startdate, &g.Enddate, &g.Customer.SubId, &g.Consignee.SubId, &g.EsoContractNumber, &g.Eso.Id,
// 		&g.Area.Id, &g.CustomerGroup.Id, &rmi, &g.Notes, &g.MotNotes, &g.Customer.SubName, &g.Customer.SubAddr, &g.Consignee.SubName,
// 		&g.Eso.EsoName, &g.Area.AreaName, &g.CustomerGroup.CustomerGroupName, &rmr, &rmk)

// 	g.ContractMot.Id = int(rmi.Int32)
// 	g.ContractMot.ContractMotNameRu = rmr.String
// 	g.ContractMot.ContractMotNameKz = rmk.String

// 	if err != nil && err != pgx.ErrNoRows {
// 		log.Println("Failed execute from func_contract_get: ", err)
// 	}

// 	out_arr = append(out_arr, g)

// 	// output, err := json.Marshal(g)
// 	// if err != nil {
// 	// 	http.Error(w, err.Error(), 500)
// 	// 	return
// 	// }
// 	// w.Write(output)

// 	auth := models.Auth{Create: true, Read: true, Update: true, Delete: true}

// 	out_count, err := json.Marshal(models.Contract_count{Values: out_arr, Count: 1, Auth: auth})

// 	if err != nil {
// 		http.Error(w, err.Error(), 500)
// 		return
// 	}

// 	w.Write(out_count)

// 	return
// }

// // HandleGetContractObject godoc
// // @Summary Get contracts objects
// // @Description get contracts objects
// // @Tags contracts
// // @Produce  json
// // @Param id path int true "Contract by id"
// // @Param active query boolean false "active objects"
// // @Success 200 {object} models.ObjContract
// // @Failure 500
// // @Router /contracts_getobject/{id} [get]
// func (s *APG) HandleGetContractObject(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	i := vars["id"]

// 	gs := models.ObjContract{}
// 	ctx := context.Background()
// 	out_arr := []models.ObjContract{}

// 	query := r.URL.Query()

// 	dsc := ""
// 	dscs, ok := query["active"]
// 	if ok && len(dscs) > 0 {
// 		if dscs[0] == "true" || dscs[0] == "false" {
// 			dsc = dscs[0]
// 		} else {
// 			dsc = ""
// 		}
// 	}

// 	rows, err := s.Dbpool.Query(ctx, "SELECT * from func_contract_getobject($1,$2);", i, utils.NullableString(dsc))
// 	if err != nil {
// 		http.Error(w, err.Error(), 500)
// 		return
// 	}

// 	defer rows.Close()

// 	for rows.Next() {
// 		err = rows.Scan(&gs.Id, &gs.Contract.Id, &gs.Object.Id, &gs.Startdate, &gs.Enddate, &gs.Object.ObjectName, &gs.Object.RegQty,
// 			&gs.Object.FlatNumber, &gs.Object.House.Id, &gs.Object.House.HouseNumber, &gs.Object.House.BuildingNumber,
// 			&gs.Object.House.Street.City.CityName, &gs.Object.House.Street.Id, &gs.Object.House.Street.StreetName, &gs.Object.TariffGroup.Id,
// 			&gs.Object.TariffGroup.TariffGroupName, &gs.Contract.ContractNumber, &gs.Contract.Startdate, &gs.Contract.Enddate,
// 			&gs.Contract.Customer.SubId, &gs.Contract.Customer.SubName, &gs.Contract.Consignee.SubId, &gs.Contract.Consignee.SubName)
// 		if err != nil {
// 			log.Println("failed to scan row:", err)
// 		}

// 		out_arr = append(out_arr, gs)
// 	}

// 	gsc := 0
// 	err = s.Dbpool.QueryRow(ctx, "SELECT * from func_contract_getobject_cnt($1,$2);", i, utils.NullableString(dsc)).Scan(&gsc)

// 	if err != nil {
// 		http.Error(w, err.Error(), 500)
// 		return
// 	}

// 	out_count, err := json.Marshal(models.ObjContract_count{Values: out_arr, Count: gsc})
// 	if err != nil {
// 		http.Error(w, err.Error(), 500)
// 		return
// 	}

// 	w.Write(out_count)

// 	return
// }

// // HandleGetContractHist godoc
// // @Summary Contract history
// // @Description get contract history
// // @Tags contracts
// // @Produce  json
// // @Param id path int true "Contract history by id"
// // @Success 200 {object} string
// // @Failure 500
// // @Router /contracts_hist/{id} [get]
// func (s *APG) HandleGetContractHist(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	i := vars["id"]
// 	hist_arr := []string{}

// 	h := ""
// 	rows, err := s.Dbpool.Query(context.Background(), "SELECT * from func_contracts_hist($1);", i)

// 	if err != nil {
// 		http.Error(w, err.Error(), 500)
// 		return
// 	}

// 	defer rows.Close()

// 	qa := false
// 	w.Write([]byte("["))

// 	for rows.Next() {

// 		if qa {
// 			w.Write([]byte(","))
// 		}
// 		qa = true

// 		err = rows.Scan(&h)
// 		w.Write([]byte(h))

// 		if err != nil {
// 			log.Println("failed to scan row:", err)
// 		}
// 		hist_arr = append(hist_arr, h)
// 	}

// 	w.Write([]byte("]"))

// 	return
// }
