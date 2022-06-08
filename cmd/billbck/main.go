package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/u2lentaru/billbck/internal/api"
	"github.com/u2lentaru/billbck/internal/models"
	"github.com/u2lentaru/billbck/internal/utils"

	// "github.com/mfuentesg/go-jwtmiddleware"

	httpSwagger "github.com/swaggo/http-swagger"
	_ "github.com/u2lentaru/billbck/cmd/billbck/docs"
	"github.com/urfave/negroni"
)

//PG - server struct
type PG struct {
	dbpool *pgxpool.Pool
}

// @title Billing Backend Server
// @version 1.0
// @description This is a backend server.
// @termsOfService http://swagger.io/terms/

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

//posterc.kz:44475 localhost:8080
// @host posterc.kz:44475
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InVzZXIxIiwicGFzc3dvcmQiOiJ1c2VyMSJ9.-qgJjYhayo7CT1YD1xLB36Xytf1HprRBeLbi5tZcOPE
func main() {
	ctx := context.Background()
	url := "postgres://postgres:postgres@" + os.Getenv("DB_HOST") + ":5432/postgres"
	// url := "postgres://postgres:postgres@" + os.Getenv("DB_HOST") + ":5432/billing"

	//ApiKeyAuth Bearer

	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		log.Fatal(err)
	}

	cfg.MaxConns = 8
	cfg.MinConns = 1

	dbpool, err := pgxpool.ConnectConfig(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer dbpool.Close()

	rows, err := dbpool.Query(ctx, "SELECT version();")
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	for rows.Next() {
		v := ""
		err = rows.Scan(&v)

		if err != nil {
			log.Println("failed to scan row:", err)
		}

		log.Println("version:", v)
	}

	pgs := PG{dbpool}
	apg := api.APG{Dbpool: dbpool}
	route := mux.NewRouter()

	// route.HandleFunc("/", handleRoot).Methods("GET", "OPTIONS")
	route.HandleFunc("/", pgs.handleLogin).Methods("GET", "OPTIONS")
	route.HandleFunc("/admin/", pgs.handleAdmin).Methods("GET", "OPTIONS")
	// route.HandleFunc("/login/", pgs.handleLogin).Methods("GET", "OPTIONS")
	route.HandleFunc("/form_types", apg.HandleFormTypes).Methods("GET", "OPTIONS")
	route.HandleFunc("/form_types_add", apg.HandleAddFormType).Methods("POST", "OPTIONS")
	route.HandleFunc("/form_types_upd", apg.HandleUpdFormType).Methods("POST", "OPTIONS")
	route.HandleFunc("/form_types_del", apg.HandleDelFormType).Methods("POST", "OPTIONS")
	route.HandleFunc("/form_types/{id:[0-9]+}", apg.HandleGetFormType).Methods("GET", "OPTIONS")
	route.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler).Methods("GET", "OPTIONS")
	route.HandleFunc("/sub_types", apg.HandleSubTypes).Methods("GET", "OPTIONS")
	route.HandleFunc("/sub_types_add", apg.HandleAddSubType).Methods("POST", "OPTIONS")
	route.HandleFunc("/sub_types_upd", apg.HandleUpdSubType).Methods("POST", "OPTIONS")
	route.HandleFunc("/sub_types_del", apg.HandleDelSubType).Methods("POST", "OPTIONS")
	route.HandleFunc("/sub_types/{id:[0-9]+}", apg.HandleGetSubType).Methods("GET", "OPTIONS")
	route.HandleFunc("/subjects", apg.HandleSubjects).Methods("GET", "OPTIONS")
	route.HandleFunc("/subjects_add", apg.HandleAddSubject).Methods("POST", "OPTIONS")
	route.HandleFunc("/subjects_upd", apg.HandleUpdSubject).Methods("POST", "OPTIONS")
	route.HandleFunc("/subjects_del", apg.HandleDelSubject).Methods("POST", "OPTIONS")
	route.HandleFunc("/subjects/{id:[0-9]+}", apg.HandleGetSubject).Methods("GET", "OPTIONS")
	route.HandleFunc("/subjects_hist/{id:[0-9]+}", apg.HandleGetSubjectHist).Methods("GET", "OPTIONS")
	route.HandleFunc("/positions", apg.HandlePositions).Methods("GET", "OPTIONS")
	route.HandleFunc("/positions_add", apg.HandleAddPosition).Methods("POST", "OPTIONS")
	route.HandleFunc("/positions_upd", apg.HandleUpdPosition).Methods("POST", "OPTIONS")
	route.HandleFunc("/positions_del", apg.HandleDelPosition).Methods("POST", "OPTIONS")
	route.HandleFunc("/positions/{id:[0-9]+}", apg.HandleGetPosition).Methods("GET", "OPTIONS")
	route.HandleFunc("/banks", apg.HandleBanks).Methods("GET", "OPTIONS")
	route.HandleFunc("/banks_add", apg.HandleAddBank).Methods("POST", "OPTIONS")
	route.HandleFunc("/banks_upd", apg.HandleUpdBank).Methods("POST", "OPTIONS")
	route.HandleFunc("/banks_del", apg.HandleDelBank).Methods("POST", "OPTIONS")
	route.HandleFunc("/banks/{id:[0-9]+}", apg.HandleGetBank).Methods("GET", "OPTIONS")
	route.HandleFunc("/org_info", apg.HandleOrgInfos).Methods("GET", "OPTIONS")
	route.HandleFunc("/org_info_add", apg.HandleAddOrgInfo).Methods("POST", "OPTIONS")
	route.HandleFunc("/org_info_upd", apg.HandleUpdOrgInfo).Methods("POST", "OPTIONS")
	route.HandleFunc("/org_info_del", apg.HandleDelOrgInfo).Methods("POST", "OPTIONS")
	route.HandleFunc("/org_info/{id:[0-9]+}", apg.HandleGetOrgInfo).Methods("GET", "OPTIONS")
	route.HandleFunc("/sub_banks", apg.HandleSubBanks).Methods("GET", "OPTIONS")
	route.HandleFunc("/sub_banks_add", apg.HandleAddSubBank).Methods("POST", "OPTIONS")
	route.HandleFunc("/sub_banks_upd", apg.HandleUpdSubBank).Methods("POST", "OPTIONS")
	route.HandleFunc("/sub_banks_del", apg.HandleDelSubBank).Methods("POST", "OPTIONS")
	route.HandleFunc("/sub_banks/{id:[0-9]+}", apg.HandleGetSubBank).Methods("GET", "OPTIONS")
	route.HandleFunc("/sub_banks_setactive/{id:[0-9]+}", apg.HandleGetSubBankSetActive).Methods("POST", "OPTIONS")
	route.HandleFunc("/building_types", apg.HandleBuildingTypes).Methods("GET", "OPTIONS")
	route.HandleFunc("/building_types_add", apg.HandleAddBuildingType).Methods("POST", "OPTIONS")
	route.HandleFunc("/building_types_upd", apg.HandleUpdBuildingType).Methods("POST", "OPTIONS")
	route.HandleFunc("/building_types_del", apg.HandleDelBuildingType).Methods("POST", "OPTIONS")
	route.HandleFunc("/building_types/{id:[0-9]+}", apg.HandleGetBuildingType).Methods("GET", "OPTIONS")
	route.HandleFunc("/areas", apg.HandleAreas).Methods("GET", "OPTIONS")
	route.HandleFunc("/areas_add", apg.HandleAddArea).Methods("POST", "OPTIONS")
	route.HandleFunc("/areas_upd", apg.HandleUpdArea).Methods("POST", "OPTIONS")
	route.HandleFunc("/areas_del", apg.HandleDelArea).Methods("POST", "OPTIONS")
	route.HandleFunc("/areas/{id:[0-9]+}", apg.HandleGetArea).Methods("GET", "OPTIONS")
	route.HandleFunc("/ksk", apg.HandleKsk).Methods("GET", "OPTIONS")
	route.HandleFunc("/ksk_add", apg.HandleAddKsk).Methods("POST", "OPTIONS")
	route.HandleFunc("/ksk_upd", apg.HandleUpdKsk).Methods("POST", "OPTIONS")
	route.HandleFunc("/ksk_del", apg.HandleDelKsk).Methods("POST", "OPTIONS")
	route.HandleFunc("/ksk/{id:[0-9]+}", apg.HandleGetKsk).Methods("GET", "OPTIONS")
	route.HandleFunc("/sectors", apg.HandleSectors).Methods("GET", "OPTIONS")
	route.HandleFunc("/sectors_add", apg.HandleAddSector).Methods("POST", "OPTIONS")
	route.HandleFunc("/sectors_upd", apg.HandleUpdSector).Methods("POST", "OPTIONS")
	route.HandleFunc("/sectors_del", apg.HandleDelSector).Methods("POST", "OPTIONS")
	route.HandleFunc("/sectors/{id:[0-9]+}", apg.HandleGetSector).Methods("GET", "OPTIONS")
	route.HandleFunc("/connectors", apg.HandleConnectors).Methods("GET", "OPTIONS")
	route.HandleFunc("/connectors_add", apg.HandleAddConnector).Methods("POST", "OPTIONS")
	route.HandleFunc("/connectors_upd", apg.HandleUpdConnector).Methods("POST", "OPTIONS")
	route.HandleFunc("/connectors_del", apg.HandleDelConnector).Methods("POST", "OPTIONS")
	route.HandleFunc("/connectors/{id:[0-9]+}", apg.HandleGetConnector).Methods("GET", "OPTIONS")
	route.HandleFunc("/input_types", apg.HandleInputTypes).Methods("GET", "OPTIONS")
	route.HandleFunc("/input_types_add", apg.HandleAddInputType).Methods("POST", "OPTIONS")
	route.HandleFunc("/input_types_upd", apg.HandleUpdInputType).Methods("POST", "OPTIONS")
	route.HandleFunc("/input_types_del", apg.HandleDelInputType).Methods("POST", "OPTIONS")
	route.HandleFunc("/input_types/{id:[0-9]+}", apg.HandleGetInputType).Methods("GET", "OPTIONS")
	route.HandleFunc("/reliabilities", apg.HandleReliabilities).Methods("GET", "OPTIONS")
	route.HandleFunc("/reliabilities_add", apg.HandleAddReliability).Methods("POST", "OPTIONS")
	route.HandleFunc("/reliabilities_upd", apg.HandleUpdReliability).Methods("POST", "OPTIONS")
	route.HandleFunc("/reliabilities_del", apg.HandleDelReliability).Methods("POST", "OPTIONS")
	route.HandleFunc("/reliabilities/{id:[0-9]+}", apg.HandleGetReliability).Methods("GET", "OPTIONS")
	route.HandleFunc("/voltages", apg.HandleVoltages).Methods("GET", "OPTIONS")
	route.HandleFunc("/voltages_add", apg.HandleAddVoltage).Methods("POST", "OPTIONS")
	route.HandleFunc("/voltages_upd", apg.HandleUpdVoltage).Methods("POST", "OPTIONS")
	route.HandleFunc("/voltages_del", apg.HandleDelVoltage).Methods("POST", "OPTIONS")
	route.HandleFunc("/voltages/{id:[0-9]+}", apg.HandleGetVoltage).Methods("GET", "OPTIONS")
	route.HandleFunc("/eso", apg.HandleEso).Methods("GET", "OPTIONS")
	route.HandleFunc("/eso_add", apg.HandleAddEso).Methods("POST", "OPTIONS")
	route.HandleFunc("/eso_upd", apg.HandleUpdEso).Methods("POST", "OPTIONS")
	route.HandleFunc("/eso_del", apg.HandleDelEso).Methods("POST", "OPTIONS")
	route.HandleFunc("/eso/{id:[0-9]+}", apg.HandleGetEso).Methods("GET", "OPTIONS")
	route.HandleFunc("/customergroups", apg.HandleCustomerGroups).Methods("GET", "OPTIONS")
	route.HandleFunc("/customergroups_add", apg.HandleAddCustomerGroup).Methods("POST", "OPTIONS")
	route.HandleFunc("/customergroups_upd", apg.HandleUpdCustomerGroup).Methods("POST", "OPTIONS")
	route.HandleFunc("/customergroups_del", apg.HandleDelCustomerGroup).Methods("POST", "OPTIONS")
	route.HandleFunc("/customergroups/{id:[0-9]+}", apg.HandleGetCustomerGroup).Methods("GET", "OPTIONS")
	route.HandleFunc("/objtypes", apg.HandleObjTypes).Methods("GET", "OPTIONS")
	route.HandleFunc("/objtypes_add", apg.HandleAddObjType).Methods("POST", "OPTIONS")
	route.HandleFunc("/objtypes_upd", apg.HandleUpdObjType).Methods("POST", "OPTIONS")
	route.HandleFunc("/objtypes_del", apg.HandleDelObjType).Methods("POST", "OPTIONS")
	route.HandleFunc("/objtypes/{id:[0-9]+}", apg.HandleGetObjType).Methods("GET", "OPTIONS")
	route.HandleFunc("/uzo", apg.HandleUzo).Methods("GET", "OPTIONS")
	route.HandleFunc("/uzo_add", apg.HandleAddUzo).Methods("POST", "OPTIONS")
	route.HandleFunc("/uzo_upd", apg.HandleUpdUzo).Methods("POST", "OPTIONS")
	route.HandleFunc("/uzo_del", apg.HandleDelUzo).Methods("POST", "OPTIONS")
	route.HandleFunc("/uzo/{id:[0-9]+}", apg.HandleGetUzo).Methods("GET", "OPTIONS")
	route.HandleFunc("/putypes", apg.HandlePuTypes).Methods("GET", "OPTIONS")
	route.HandleFunc("/putypes_add", apg.HandleAddPuType).Methods("POST", "OPTIONS")
	route.HandleFunc("/putypes_upd", apg.HandleUpdPuType).Methods("POST", "OPTIONS")
	route.HandleFunc("/putypes_del", apg.HandleDelPuType).Methods("POST", "OPTIONS")
	route.HandleFunc("/putypes/{id:[0-9]+}", apg.HandleGetPuType).Methods("GET", "OPTIONS")
	route.HandleFunc("/acttypes", apg.HandleActTypes).Methods("GET", "OPTIONS")
	route.HandleFunc("/acttypes_add", apg.HandleAddActType).Methods("POST", "OPTIONS")
	route.HandleFunc("/acttypes_upd", apg.HandleUpdActType).Methods("POST", "OPTIONS")
	route.HandleFunc("/acttypes_del", apg.HandleDelActType).Methods("POST", "OPTIONS")
	route.HandleFunc("/acttypes/{id:[0-9]+}", apg.HandleGetActType).Methods("GET", "OPTIONS")
	route.HandleFunc("/cashdesks", apg.HandleCashdesks).Methods("GET", "OPTIONS")
	route.HandleFunc("/cashdesks_add", apg.HandleAddCashdesk).Methods("POST", "OPTIONS")
	route.HandleFunc("/cashdesks_upd", apg.HandleUpdCashdesk).Methods("POST", "OPTIONS")
	route.HandleFunc("/cashdesks_del", apg.HandleDelCashdesk).Methods("POST", "OPTIONS")
	route.HandleFunc("/cashdesks/{id:[0-9]+}", apg.HandleGetCashdesk).Methods("GET", "OPTIONS")
	route.HandleFunc("/paymenttypes", apg.HandlePaymentTypes).Methods("GET", "OPTIONS")
	route.HandleFunc("/paymenttypes_add", apg.HandleAddPaymentType).Methods("POST", "OPTIONS")
	route.HandleFunc("/paymenttypes_upd", apg.HandleUpdPaymentType).Methods("POST", "OPTIONS")
	route.HandleFunc("/paymenttypes_del", apg.HandleDelPaymentType).Methods("POST", "OPTIONS")
	route.HandleFunc("/paymenttypes/{id:[0-9]+}", apg.HandleGetPaymentType).Methods("GET", "OPTIONS")
	route.HandleFunc("/chargetypes", apg.HandleChargeTypes).Methods("GET", "OPTIONS")
	route.HandleFunc("/chargetypes_add", apg.HandleAddChargeType).Methods("POST", "OPTIONS")
	route.HandleFunc("/chargetypes_upd", apg.HandleUpdChargeType).Methods("POST", "OPTIONS")
	route.HandleFunc("/chargetypes_del", apg.HandleDelChargeType).Methods("POST", "OPTIONS")
	route.HandleFunc("/chargetypes/{id:[0-9]+}", apg.HandleGetChargeType).Methods("GET", "OPTIONS")
	route.HandleFunc("/tariffgroups", apg.HandleTariffGroups).Methods("GET", "OPTIONS")
	route.HandleFunc("/tariffgroups_add", apg.HandleAddTariffGroup).Methods("POST", "OPTIONS")
	route.HandleFunc("/tariffgroups_upd", apg.HandleUpdTariffGroup).Methods("POST", "OPTIONS")
	route.HandleFunc("/tariffgroups_del", apg.HandleDelTariffGroup).Methods("POST", "OPTIONS")
	route.HandleFunc("/tariffgroups/{id:[0-9]+}", apg.HandleGetTariffGroup).Methods("GET", "OPTIONS")
	route.HandleFunc("/tariffs", apg.HandleTariffs).Methods("GET", "OPTIONS")
	route.HandleFunc("/tariffs_add", apg.HandleAddTariff).Methods("POST", "OPTIONS")
	route.HandleFunc("/tariffs_upd", apg.HandleUpdTariff).Methods("POST", "OPTIONS")
	route.HandleFunc("/tariffs_del", apg.HandleDelTariff).Methods("POST", "OPTIONS")
	route.HandleFunc("/tariffs/{id:[0-9]+}", apg.HandleGetTariff).Methods("GET", "OPTIONS")
	route.HandleFunc("/rp", apg.HandleRp).Methods("GET", "OPTIONS")
	route.HandleFunc("/rp_add", apg.HandleAddRp).Methods("POST", "OPTIONS")
	route.HandleFunc("/rp_upd", apg.HandleUpdRp).Methods("POST", "OPTIONS")
	route.HandleFunc("/rp_del", apg.HandleDelRp).Methods("POST", "OPTIONS")
	route.HandleFunc("/rp/{id:[0-9]+}", apg.HandleGetRp).Methods("GET", "OPTIONS")
	route.HandleFunc("/tp", apg.HandleTp).Methods("GET", "OPTIONS")
	route.HandleFunc("/tp_add", apg.HandleAddTp).Methods("POST", "OPTIONS")
	route.HandleFunc("/tp_upd", apg.HandleUpdTp).Methods("POST", "OPTIONS")
	route.HandleFunc("/tp_del", apg.HandleDelTp).Methods("POST", "OPTIONS")
	route.HandleFunc("/tp/{id:[0-9]+}", apg.HandleGetTp).Methods("GET", "OPTIONS")
	route.HandleFunc("/grp", apg.HandleGRp).Methods("GET", "OPTIONS")
	route.HandleFunc("/grp_add", apg.HandleAddGRp).Methods("POST", "OPTIONS")
	route.HandleFunc("/grp_upd", apg.HandleUpdGRp).Methods("POST", "OPTIONS")
	route.HandleFunc("/grp_del", apg.HandleDelGRp).Methods("POST", "OPTIONS")
	route.HandleFunc("/grp/{id:[0-9]+}", apg.HandleGetGRp).Methods("GET", "OPTIONS")
	route.HandleFunc("/streets", apg.HandleStreets).Methods("GET", "OPTIONS")
	route.HandleFunc("/streets_add", apg.HandleAddStreet).Methods("POST", "OPTIONS")
	route.HandleFunc("/streets_upd", apg.HandleUpdStreet).Methods("POST", "OPTIONS")
	route.HandleFunc("/streets_del", apg.HandleDelStreet).Methods("POST", "OPTIONS")
	route.HandleFunc("/streets/{id:[0-9]+}", apg.HandleGetStreet).Methods("GET", "OPTIONS")
	route.HandleFunc("/houses", apg.HandleHouses).Methods("GET", "OPTIONS")
	route.HandleFunc("/houses_add", apg.HandleAddHouse).Methods("POST", "OPTIONS")
	route.HandleFunc("/houses_upd", apg.HandleUpdHouse).Methods("POST", "OPTIONS")
	route.HandleFunc("/houses_del", apg.HandleDelHouse).Methods("POST", "OPTIONS")
	route.HandleFunc("/houses/{id:[0-9]+}", apg.HandleGetHouse).Methods("GET", "OPTIONS")
	route.HandleFunc("/cities", apg.HandleCities).Methods("GET", "OPTIONS")
	route.HandleFunc("/cities_add", apg.HandleAddCity).Methods("POST", "OPTIONS")
	route.HandleFunc("/cities_upd", apg.HandleUpdCity).Methods("POST", "OPTIONS")
	route.HandleFunc("/cities_del", apg.HandleDelCity).Methods("POST", "OPTIONS")
	route.HandleFunc("/cities/{id:[0-9]+}", apg.HandleGetCity).Methods("GET", "OPTIONS")
	route.HandleFunc("/contracts", apg.HandleContracts).Methods("GET", "OPTIONS")
	route.HandleFunc("/contracts_add", apg.HandleAddContract).Methods("POST", "OPTIONS")
	route.HandleFunc("/contracts_upd", apg.HandleUpdContract).Methods("POST", "OPTIONS")
	route.HandleFunc("/contracts_del", apg.HandleDelContract).Methods("POST", "OPTIONS")
	route.HandleFunc("/contracts/{id:[0-9]+}", apg.HandleGetContract).Methods("GET", "OPTIONS")
	route.HandleFunc("/contracts_getobject/{id:[0-9]+}", apg.HandleGetContractObject).Methods("GET", "OPTIONS")
	route.HandleFunc("/contracts_hist/{id:[0-9]+}", apg.HandleGetContractHist).Methods("GET", "OPTIONS")
	route.HandleFunc("/objects", apg.HandleObjects).Methods("GET", "OPTIONS")
	route.HandleFunc("/objects_add", apg.HandleAddObject).Methods("POST", "OPTIONS")
	route.HandleFunc("/objects_upd", apg.HandleUpdObject).Methods("POST", "OPTIONS")
	route.HandleFunc("/objects_del", apg.HandleDelObject).Methods("POST", "OPTIONS")
	route.HandleFunc("/objects/{id:[0-9]+}", apg.HandleGetObject).Methods("GET", "OPTIONS")
	route.HandleFunc("/objects_getcontract/{id:[0-9]+}", apg.HandleGetObjectContract).Methods("GET", "OPTIONS")
	route.HandleFunc("/objects_mff/{hid:[0-9]+}", apg.HandleGetObjectMff).Methods("GET", "OPTIONS")
	route.HandleFunc("/objcontracts", apg.HandleObjContracts).Methods("GET", "OPTIONS")
	route.HandleFunc("/objcontracts_add", apg.HandleAddObjContract).Methods("POST", "OPTIONS")
	route.HandleFunc("/objcontracts_upd", apg.HandleUpdObjContract).Methods("POST", "OPTIONS")
	route.HandleFunc("/objcontracts_del", apg.HandleDelObjContract).Methods("POST", "OPTIONS")
	route.HandleFunc("/objcontracts/{id:[0-9]+}", apg.HandleGetObjContract).Methods("GET", "OPTIONS")
	route.HandleFunc("/acts", apg.HandleActs).Methods("GET", "OPTIONS")
	route.HandleFunc("/acts_add", apg.HandleAddAct).Methods("POST", "OPTIONS")
	route.HandleFunc("/acts_upd", apg.HandleUpdAct).Methods("POST", "OPTIONS")
	route.HandleFunc("/acts_del", apg.HandleDelAct).Methods("POST", "OPTIONS")
	route.HandleFunc("/acts_activate", apg.HandleActActivate).Methods("GET", "OPTIONS")
	route.HandleFunc("/acts/{id:[0-9]+}", apg.HandleGetAct).Methods("GET", "OPTIONS")
	route.HandleFunc("/actdetails", apg.HandleActDetails).Methods("GET", "OPTIONS")
	route.HandleFunc("/actdetails_add", apg.HandleAddActDetail).Methods("POST", "OPTIONS")
	route.HandleFunc("/actdetails_upd", apg.HandleUpdActDetail).Methods("POST", "OPTIONS")
	route.HandleFunc("/actdetails_del", apg.HandleDelActDetail).Methods("POST", "OPTIONS")
	route.HandleFunc("/actdetails/{id:[0-9]+}", apg.HandleGetActDetail).Methods("GET", "OPTIONS")
	route.HandleFunc("/pu", apg.HandlePu).Methods("GET", "OPTIONS")
	route.HandleFunc("/pu_add", apg.HandleAddPu).Methods("POST", "OPTIONS")
	route.HandleFunc("/pu_obj_add", apg.HandleAddPu).Methods("POST", "OPTIONS")
	route.HandleFunc("/pu_upd", apg.HandleUpdPu).Methods("POST", "OPTIONS")
	route.HandleFunc("/pu_del", apg.HandleDelPu).Methods("POST", "OPTIONS")
	route.HandleFunc("/pu/{id:[0-9]+}", apg.HandleGetPu).Methods("GET", "OPTIONS")
	route.HandleFunc("/pu_obj", apg.HandlePuObj).Methods("GET", "OPTIONS")
	route.HandleFunc("/puvalues", apg.HandlePuValues).Methods("GET", "OPTIONS")
	route.HandleFunc("/puvalues_add", apg.HandleAddPuValue).Methods("POST", "OPTIONS")
	route.HandleFunc("/puvalues_upd", apg.HandleUpdPuValue).Methods("POST", "OPTIONS")
	route.HandleFunc("/puvalues_del", apg.HandleDelPuValue).Methods("POST", "OPTIONS")
	route.HandleFunc("/puvalues/{id:[0-9]+}", apg.HandleGetPuValue).Methods("GET", "OPTIONS")
	route.HandleFunc("/puvalues_askue_prev", apg.HandlePuValuesAskuePreview).Methods("POST", "OPTIONS")
	route.HandleFunc("/puvalues_askue", apg.HandlePuValuesAskue).Methods("POST", "OPTIONS")
	route.HandleFunc("/balance", apg.HandleBalance).Methods("GET", "OPTIONS")
	route.HandleFunc("/balance/{id:[0-9]+}/{tid:[0-9]+}", apg.HandleGetBalance).Methods("GET", "OPTIONS")
	route.HandleFunc("/balance_sum", apg.HandleBalanceSum).Methods("GET", "OPTIONS")
	route.HandleFunc("/balance_sum1", apg.HandleBalanceSum1).Methods("GET", "OPTIONS")
	route.HandleFunc("/balance_sum0", apg.HandleBalanceSum0).Methods("GET", "OPTIONS")
	route.HandleFunc("/balance_tab1", apg.HandleBalanceTab1).Methods("GET", "OPTIONS")
	route.HandleFunc("/balance_tab0", apg.HandleBalanceTab0).Methods("GET", "OPTIONS")
	route.HandleFunc("/balance_branch", apg.HandleBalanceBranch).Methods("GET", "OPTIONS")
	route.HandleFunc("/subpu", apg.HandleSubPu).Methods("GET", "OPTIONS")
	route.HandleFunc("/subpu_add", apg.HandleAddSubPu).Methods("POST", "OPTIONS")
	route.HandleFunc("/subpu_upd", apg.HandleUpdSubPu).Methods("POST", "OPTIONS")
	route.HandleFunc("/subpu_del", apg.HandleDelSubPu).Methods("POST", "OPTIONS")
	route.HandleFunc("/subpu/{id:[0-9]+}", apg.HandleGetSubPu).Methods("GET", "OPTIONS")
	route.HandleFunc("/subpu_prl", apg.HandlePrlSubPu).Methods("GET", "OPTIONS")
	route.HandleFunc("/staff", apg.HandleStaff).Methods("GET", "OPTIONS")
	route.HandleFunc("/staff_add", apg.HandleAddStaff).Methods("POST", "OPTIONS")
	route.HandleFunc("/staff_upd", apg.HandleUpdStaff).Methods("POST", "OPTIONS")
	route.HandleFunc("/staff_del", apg.HandleDelStaff).Methods("POST", "OPTIONS")
	route.HandleFunc("/staff/{id:[0-9]+}", apg.HandleGetStaff).Methods("GET", "OPTIONS")
	route.HandleFunc("/tgu", apg.HandleTgu).Methods("GET", "OPTIONS")
	route.HandleFunc("/tgu_add", apg.HandleAddTgu).Methods("POST", "OPTIONS")
	route.HandleFunc("/tgu_upd", apg.HandleUpdTgu).Methods("POST", "OPTIONS")
	route.HandleFunc("/tgu_del", apg.HandleDelTgu).Methods("POST", "OPTIONS")
	route.HandleFunc("/tgu/{id:[0-9]+}", apg.HandleGetTgu).Methods("GET", "OPTIONS")
	route.HandleFunc("/conclusions", apg.HandleConclusions).Methods("GET", "OPTIONS")
	route.HandleFunc("/conclusions_add", apg.HandleAddConclusion).Methods("POST", "OPTIONS")
	route.HandleFunc("/conclusions_upd", apg.HandleUpdConclusion).Methods("POST", "OPTIONS")
	route.HandleFunc("/conclusions_del", apg.HandleDelConclusion).Methods("POST", "OPTIONS")
	route.HandleFunc("/conclusions/{id:[0-9]+}", apg.HandleGetConclusion).Methods("GET", "OPTIONS")
	route.HandleFunc("/shutdowntypes", apg.HandleShutdownTypes).Methods("GET", "OPTIONS")
	route.HandleFunc("/shutdowntypes_add", apg.HandleAddShutdownType).Methods("POST", "OPTIONS")
	route.HandleFunc("/shutdowntypes_upd", apg.HandleUpdShutdownType).Methods("POST", "OPTIONS")
	route.HandleFunc("/shutdowntypes_del", apg.HandleDelShutdownType).Methods("POST", "OPTIONS")
	route.HandleFunc("/shutdowntypes/{id:[0-9]+}", apg.HandleGetShutdownType).Methods("GET", "OPTIONS")
	route.HandleFunc("/violations", apg.HandleViolations).Methods("GET", "OPTIONS")
	route.HandleFunc("/violations_add", apg.HandleAddViolation).Methods("POST", "OPTIONS")
	route.HandleFunc("/violations_upd", apg.HandleUpdViolation).Methods("POST", "OPTIONS")
	route.HandleFunc("/violations_del", apg.HandleDelViolation).Methods("POST", "OPTIONS")
	route.HandleFunc("/violations/{id:[0-9]+}", apg.HandleGetViolation).Methods("GET", "OPTIONS")
	route.HandleFunc("/reasons", apg.HandleReasons).Methods("GET", "OPTIONS")
	route.HandleFunc("/reasons_add", apg.HandleAddReason).Methods("POST", "OPTIONS")
	route.HandleFunc("/reasons_upd", apg.HandleUpdReason).Methods("POST", "OPTIONS")
	route.HandleFunc("/reasons_del", apg.HandleDelReason).Methods("POST", "OPTIONS")
	route.HandleFunc("/reasons/{id:[0-9]+}", apg.HandleGetReason).Methods("GET", "OPTIONS")
	route.HandleFunc("/seals", apg.HandleSeals).Methods("GET", "OPTIONS")
	route.HandleFunc("/seals_add", apg.HandleAddSeal).Methods("POST", "OPTIONS")
	route.HandleFunc("/seals_upd", apg.HandleUpdSeal).Methods("POST", "OPTIONS")
	route.HandleFunc("/seals_del", apg.HandleDelSeal).Methods("POST", "OPTIONS")
	route.HandleFunc("/seals/{id:[0-9]+}", apg.HandleGetSeal).Methods("GET", "OPTIONS")
	route.HandleFunc("/sealtypes", apg.HandleSealTypes).Methods("GET", "OPTIONS")
	route.HandleFunc("/sealtypes_add", apg.HandleAddSealType).Methods("POST", "OPTIONS")
	route.HandleFunc("/sealtypes_upd", apg.HandleUpdSealType).Methods("POST", "OPTIONS")
	route.HandleFunc("/sealtypes_del", apg.HandleDelSealType).Methods("POST", "OPTIONS")
	route.HandleFunc("/sealtypes/{id:[0-9]+}", apg.HandleGetSealType).Methods("GET", "OPTIONS")
	route.HandleFunc("/sealcolours", apg.HandleSealColours).Methods("GET", "OPTIONS")
	route.HandleFunc("/sealcolours_add", apg.HandleAddSealColour).Methods("POST", "OPTIONS")
	route.HandleFunc("/sealcolours_upd", apg.HandleUpdSealColour).Methods("POST", "OPTIONS")
	route.HandleFunc("/sealcolours_del", apg.HandleDelSealColour).Methods("POST", "OPTIONS")
	route.HandleFunc("/sealcolours/{id:[0-9]+}", apg.HandleGetSealColour).Methods("GET", "OPTIONS")
	route.HandleFunc("/sealstatuses", apg.HandleSealStatuses).Methods("GET", "OPTIONS")
	route.HandleFunc("/sealstatuses_add", apg.HandleAddSealStatus).Methods("POST", "OPTIONS")
	route.HandleFunc("/sealstatuses_upd", apg.HandleUpdSealStatus).Methods("POST", "OPTIONS")
	route.HandleFunc("/sealstatuses_del", apg.HandleDelSealStatus).Methods("POST", "OPTIONS")
	route.HandleFunc("/sealstatuses/{id:[0-9]+}", apg.HandleGetSealStatus).Methods("GET", "OPTIONS")
	route.HandleFunc("/charges", apg.HandleCharges).Methods("GET", "OPTIONS")
	route.HandleFunc("/charges_add", apg.HandleAddCharge).Methods("POST", "OPTIONS")
	route.HandleFunc("/charges_upd", apg.HandleUpdCharge).Methods("POST", "OPTIONS")
	route.HandleFunc("/charges_del", apg.HandleDelCharge).Methods("POST", "OPTIONS")
	route.HandleFunc("/charges/{id:[0-9]+}", apg.HandleGetCharge).Methods("GET", "OPTIONS")
	route.HandleFunc("/charges_run/{id:[0-9]+}", apg.HandleChargeRun).Methods("GET", "OPTIONS")
	route.HandleFunc("/payments", apg.HandlePayments).Methods("GET", "OPTIONS")
	route.HandleFunc("/payments_add", apg.HandleAddPayment).Methods("POST", "OPTIONS")
	route.HandleFunc("/payments_upd", apg.HandleUpdPayment).Methods("POST", "OPTIONS")
	route.HandleFunc("/payments_del", apg.HandleDelPayment).Methods("POST", "OPTIONS")
	route.HandleFunc("/payments/{id:[0-9]+}", apg.HandleGetPayment).Methods("GET", "OPTIONS")
	route.HandleFunc("/equipmenttypes", apg.HandleEquipmentTypes).Methods("GET", "OPTIONS")
	route.HandleFunc("/equipmenttypes_add", apg.HandleAddEquipmentType).Methods("POST", "OPTIONS")
	route.HandleFunc("/equipmenttypes_upd", apg.HandleUpdEquipmentType).Methods("POST", "OPTIONS")
	route.HandleFunc("/equipmenttypes_del", apg.HandleDelEquipmentType).Methods("POST", "OPTIONS")
	route.HandleFunc("/equipmenttypes/{id:[0-9]+}", apg.HandleGetEquipmentType).Methods("GET", "OPTIONS")
	route.HandleFunc("/equipment", apg.HandleEquipment).Methods("GET", "OPTIONS")
	route.HandleFunc("/equipment_add", apg.HandleAddEquipment).Methods("POST", "OPTIONS")
	route.HandleFunc("/equipment_addlist", apg.HandleAddEquipmentList).Methods("POST", "OPTIONS")
	route.HandleFunc("/equipment_upd", apg.HandleUpdEquipment).Methods("POST", "OPTIONS")
	route.HandleFunc("/equipment_del", apg.HandleDelEquipment).Methods("POST", "OPTIONS")
	route.HandleFunc("/equipment_delobj", apg.HandleDelObjEquipment).Methods("POST", "OPTIONS")
	route.HandleFunc("/equipment/{id:[0-9]+}", apg.HandleGetEquipment).Methods("GET", "OPTIONS")
	route.HandleFunc("/calculationtypes", apg.HandleCalculationTypes).Methods("GET", "OPTIONS")
	route.HandleFunc("/calculationtypes_add", apg.HandleAddCalculationType).Methods("POST", "OPTIONS")
	route.HandleFunc("/calculationtypes_upd", apg.HandleUpdCalculationType).Methods("POST", "OPTIONS")
	route.HandleFunc("/calculationtypes_del", apg.HandleDelCalculationType).Methods("POST", "OPTIONS")
	route.HandleFunc("/calculationtypes/{id:[0-9]+}", apg.HandleGetCalculationType).Methods("GET", "OPTIONS")
	route.HandleFunc("/results", apg.HandleResults).Methods("GET", "OPTIONS")
	route.HandleFunc("/results_add", apg.HandleAddResult).Methods("POST", "OPTIONS")
	route.HandleFunc("/results_upd", apg.HandleUpdResult).Methods("POST", "OPTIONS")
	route.HandleFunc("/results_del", apg.HandleDelResult).Methods("POST", "OPTIONS")
	route.HandleFunc("/results/{id:[0-9]+}", apg.HandleGetResult).Methods("GET", "OPTIONS")
	route.HandleFunc("/servicetypes", apg.HandleServiceTypes).Methods("GET", "OPTIONS")
	route.HandleFunc("/servicetypes_add", apg.HandleAddServiceType).Methods("POST", "OPTIONS")
	route.HandleFunc("/servicetypes_upd", apg.HandleUpdServiceType).Methods("POST", "OPTIONS")
	route.HandleFunc("/servicetypes_del", apg.HandleDelServiceType).Methods("POST", "OPTIONS")
	route.HandleFunc("/servicetypes/{id:[0-9]+}", apg.HandleGetServiceType).Methods("GET", "OPTIONS")
	route.HandleFunc("/requesttypes", apg.HandleRequestTypes).Methods("GET", "OPTIONS")
	route.HandleFunc("/requesttypes_add", apg.HandleAddRequestType).Methods("POST", "OPTIONS")
	route.HandleFunc("/requesttypes_upd", apg.HandleUpdRequestType).Methods("POST", "OPTIONS")
	route.HandleFunc("/requesttypes_del", apg.HandleDelRequestType).Methods("POST", "OPTIONS")
	route.HandleFunc("/requesttypes/{id:[0-9]+}", apg.HandleGetRequestType).Methods("GET", "OPTIONS")
	route.HandleFunc("/claimtypes", apg.HandleClaimTypes).Methods("GET", "OPTIONS")
	route.HandleFunc("/claimtypes_add", apg.HandleAddClaimType).Methods("POST", "OPTIONS")
	route.HandleFunc("/claimtypes_upd", apg.HandleUpdClaimType).Methods("POST", "OPTIONS")
	route.HandleFunc("/claimtypes_del", apg.HandleDelClaimType).Methods("POST", "OPTIONS")
	route.HandleFunc("/claimtypes/{id:[0-9]+}", apg.HandleGetClaimType).Methods("GET", "OPTIONS")
	route.HandleFunc("/requests", apg.HandleRequests).Methods("GET", "OPTIONS")
	route.HandleFunc("/requests_add", apg.HandleAddRequest).Methods("POST", "OPTIONS")
	route.HandleFunc("/requests_upd", apg.HandleUpdRequest).Methods("POST", "OPTIONS")
	route.HandleFunc("/requests_del", apg.HandleDelRequest).Methods("POST", "OPTIONS")
	route.HandleFunc("/requests/{id:[0-9]+}", apg.HandleGetRequest).Methods("GET", "OPTIONS")
	route.HandleFunc("/objstatuses", apg.HandleObjStatuses).Methods("GET", "OPTIONS")
	route.HandleFunc("/objstatuses_add", apg.HandleAddObjStatus).Methods("POST", "OPTIONS")
	route.HandleFunc("/objstatuses_upd", apg.HandleUpdObjStatus).Methods("POST", "OPTIONS")
	route.HandleFunc("/objstatuses_del", apg.HandleDelObjStatus).Methods("POST", "OPTIONS")
	route.HandleFunc("/objstatuses/{id:[0-9]+}", apg.HandleGetObjStatus).Methods("GET", "OPTIONS")
	route.HandleFunc("/users", apg.HandleUsers).Methods("GET", "OPTIONS")
	route.HandleFunc("/users_add", apg.HandleAddUser).Methods("POST", "OPTIONS")
	route.HandleFunc("/users_upd", apg.HandleUpdUser).Methods("POST", "OPTIONS")
	route.HandleFunc("/users_del", apg.HandleDelUser).Methods("POST", "OPTIONS")
	route.HandleFunc("/users/{id:[0-9]+}", apg.HandleGetUser).Methods("GET", "OPTIONS")
	route.HandleFunc("/contractmots", apg.HandleContractMots).Methods("GET", "OPTIONS")
	route.HandleFunc("/contractmots_add", apg.HandleAddContractMot).Methods("POST", "OPTIONS")
	route.HandleFunc("/contractmots_upd", apg.HandleUpdContractMot).Methods("POST", "OPTIONS")
	route.HandleFunc("/contractmots_del", apg.HandleDelContractMot).Methods("POST", "OPTIONS")
	route.HandleFunc("/contractmots/{id:[0-9]+}", apg.HandleGetContractMot).Methods("GET", "OPTIONS")
	route.HandleFunc("/requestkinds", apg.HandleRequestKinds).Methods("GET", "OPTIONS")
	route.HandleFunc("/requestkinds_add", apg.HandleAddRequestKind).Methods("POST", "OPTIONS")
	route.HandleFunc("/requestkinds_upd", apg.HandleUpdRequestKind).Methods("POST", "OPTIONS")
	route.HandleFunc("/requestkinds_del", apg.HandleDelRequestKind).Methods("POST", "OPTIONS")
	route.HandleFunc("/requestkinds/{id:[0-9]+}", apg.HandleGetRequestKind).Methods("GET", "OPTIONS")
	route.HandleFunc("/periods", apg.HandlePeriods).Methods("GET", "OPTIONS")
	route.HandleFunc("/periods_add", apg.HandleAddPeriod).Methods("POST", "OPTIONS")
	route.HandleFunc("/periods_upd", apg.HandleUpdPeriod).Methods("POST", "OPTIONS")
	route.HandleFunc("/periods_del", apg.HandleDelPeriod).Methods("POST", "OPTIONS")
	route.HandleFunc("/periods/{id:[0-9]+}", apg.HandleGetPeriod).Methods("GET", "OPTIONS")
	route.HandleFunc("/transtypes", apg.HandleTransTypes).Methods("GET", "OPTIONS")
	route.HandleFunc("/transtypes_add", apg.HandleAddTransType).Methods("POST", "OPTIONS")
	route.HandleFunc("/transtypes_upd", apg.HandleUpdTransType).Methods("POST", "OPTIONS")
	route.HandleFunc("/transtypes_del", apg.HandleDelTransType).Methods("POST", "OPTIONS")
	route.HandleFunc("/transtypes/{id:[0-9]+}", apg.HandleGetTransType).Methods("GET", "OPTIONS")
	route.HandleFunc("/transcurr", apg.HandleTransCurr).Methods("GET", "OPTIONS")
	route.HandleFunc("/transcurr_add", apg.HandleAddTransCurr).Methods("POST", "OPTIONS")
	route.HandleFunc("/transcurr_upd", apg.HandleUpdTransCurr).Methods("POST", "OPTIONS")
	route.HandleFunc("/transcurr_del", apg.HandleDelTransCurr).Methods("POST", "OPTIONS")
	route.HandleFunc("/transcurr/{id:[0-9]+}", apg.HandleGetTransCurr).Methods("GET", "OPTIONS")
	route.HandleFunc("/transvolt", apg.HandleTransVolt).Methods("GET", "OPTIONS")
	route.HandleFunc("/transvolt_add", apg.HandleAddTransVolt).Methods("POST", "OPTIONS")
	route.HandleFunc("/transvolt_upd", apg.HandleUpdTransVolt).Methods("POST", "OPTIONS")
	route.HandleFunc("/transvolt_del", apg.HandleDelTransVolt).Methods("POST", "OPTIONS")
	route.HandleFunc("/transvolt/{id:[0-9]+}", apg.HandleGetTransVolt).Methods("GET", "OPTIONS")
	route.HandleFunc("/transpwrtypes", apg.HandleTransPwrTypes).Methods("GET", "OPTIONS")
	route.HandleFunc("/transpwrtypes_add", apg.HandleAddTransPwrType).Methods("POST", "OPTIONS")
	route.HandleFunc("/transpwrtypes_upd", apg.HandleUpdTransPwrType).Methods("POST", "OPTIONS")
	route.HandleFunc("/transpwrtypes_del", apg.HandleDelTransPwrType).Methods("POST", "OPTIONS")
	route.HandleFunc("/transpwrtypes/{id:[0-9]+}", apg.HandleGetTransPwrType).Methods("GET", "OPTIONS")
	route.HandleFunc("/transpwr", apg.HandleTransPwr).Methods("GET", "OPTIONS")
	route.HandleFunc("/transpwr_add", apg.HandleAddTransPwr).Methods("POST", "OPTIONS")
	route.HandleFunc("/transpwr_upd", apg.HandleUpdTransPwr).Methods("POST", "OPTIONS")
	route.HandleFunc("/transpwr_del", apg.HandleDelTransPwr).Methods("POST", "OPTIONS")
	route.HandleFunc("/transpwr/{id:[0-9]+}", apg.HandleGetTransPwr).Methods("GET", "OPTIONS")
	route.HandleFunc("/objtranscurr", apg.HandleObjTransCurr).Methods("GET", "OPTIONS")
	route.HandleFunc("/objtranscurr_add", apg.HandleAddObjTransCurr).Methods("POST", "OPTIONS")
	route.HandleFunc("/objtranscurr_upd", apg.HandleUpdObjTransCurr).Methods("POST", "OPTIONS")
	route.HandleFunc("/objtranscurr_del", apg.HandleDelObjTransCurr).Methods("POST", "OPTIONS")
	route.HandleFunc("/objtranscurr/{id:[0-9]+}", apg.HandleGetObjTransCurr).Methods("GET", "OPTIONS")
	route.HandleFunc("/objtranscurr_obj", apg.HandleObjTransCurrByObj).Methods("GET", "OPTIONS")
	route.HandleFunc("/objtransvolt", apg.HandleObjTransVolt).Methods("GET", "OPTIONS")
	route.HandleFunc("/objtransvolt_add", apg.HandleAddObjTransVolt).Methods("POST", "OPTIONS")
	route.HandleFunc("/objtransvolt_upd", apg.HandleUpdObjTransVolt).Methods("POST", "OPTIONS")
	route.HandleFunc("/objtransvolt_del", apg.HandleDelObjTransVolt).Methods("POST", "OPTIONS")
	route.HandleFunc("/objtransvolt/{id:[0-9]+}", apg.HandleGetObjTransVolt).Methods("GET", "OPTIONS")
	route.HandleFunc("/objtransvolt_obj", apg.HandleObjTransVoltByObj).Methods("GET", "OPTIONS")
	route.HandleFunc("/objtranspwr", apg.HandleObjTransPwr).Methods("GET", "OPTIONS")
	route.HandleFunc("/objtranspwr_add", apg.HandleAddObjTransPwr).Methods("POST", "OPTIONS")
	route.HandleFunc("/objtranspwr_upd", apg.HandleUpdObjTransPwr).Methods("POST", "OPTIONS")
	route.HandleFunc("/objtranspwr_del", apg.HandleDelObjTransPwr).Methods("POST", "OPTIONS")
	route.HandleFunc("/objtranspwr/{id:[0-9]+}", apg.HandleGetObjTransPwr).Methods("GET", "OPTIONS")
	route.HandleFunc("/objtranspwr_obj", apg.HandleObjTransPwrByObj).Methods("GET", "OPTIONS")
	route.HandleFunc("/cableresistances", apg.HandleCableResistances).Methods("GET", "OPTIONS")
	route.HandleFunc("/cableresistances_add", apg.HandleAddCableResistance).Methods("POST", "OPTIONS")
	route.HandleFunc("/cableresistances_upd", apg.HandleUpdCableResistance).Methods("POST", "OPTIONS")
	route.HandleFunc("/cableresistances_del", apg.HandleDelCableResistance).Methods("POST", "OPTIONS")
	route.HandleFunc("/cableresistances/{id:[0-9]+}", apg.HandleGetCableResistance).Methods("GET", "OPTIONS")
	route.HandleFunc("/objlines", apg.HandleObjLines).Methods("GET", "OPTIONS")
	route.HandleFunc("/objlines_add", apg.HandleAddObjLine).Methods("POST", "OPTIONS")
	route.HandleFunc("/objlines_upd", apg.HandleUpdObjLine).Methods("POST", "OPTIONS")
	route.HandleFunc("/objlines_del", apg.HandleDelObjLine).Methods("POST", "OPTIONS")
	route.HandleFunc("/objlines/{id:[0-9]+}", apg.HandleGetObjLine).Methods("GET", "OPTIONS")
	route.HandleFunc("/objlines_obj", apg.HandleObjLinesByObj).Methods("GET", "OPTIONS")
	route.HandleFunc("/askuetypes", apg.HandleAskueTypes).Methods("GET", "OPTIONS")
	route.HandleFunc("/askuetypes_add", apg.HandleAddAskueType).Methods("POST", "OPTIONS")
	route.HandleFunc("/askuetypes_upd", apg.HandleUpdAskueType).Methods("POST", "OPTIONS")
	route.HandleFunc("/askuetypes_del", apg.HandleDelAskueType).Methods("POST", "OPTIONS")
	route.HandleFunc("/askuetypes/{id:[0-9]+}", apg.HandleGetAskueType).Methods("GET", "OPTIONS")
	route.HandleFunc("/puvaluetypes", apg.HandlePuValueTypes).Methods("GET", "OPTIONS")
	route.HandleFunc("/puvaluetypes_add", apg.HandleAddPuValueType).Methods("POST", "OPTIONS")
	route.HandleFunc("/puvaluetypes_upd", apg.HandleUpdPuValueType).Methods("POST", "OPTIONS")
	route.HandleFunc("/puvaluetypes_del", apg.HandleDelPuValueType).Methods("POST", "OPTIONS")
	route.HandleFunc("/puvaluetypes/{id:[0-9]+}", apg.HandleGetPuValueType).Methods("GET", "OPTIONS")
	route.HandleFunc("/ordertypes", apg.HandleOrderTypes).Methods("GET", "OPTIONS")
	route.HandleFunc("/ordertypes_add", apg.HandleAddOrderType).Methods("POST", "OPTIONS")
	route.HandleFunc("/ordertypes_upd", apg.HandleUpdOrderType).Methods("POST", "OPTIONS")
	route.HandleFunc("/ordertypes_del", apg.HandleDelOrderType).Methods("POST", "OPTIONS")
	route.HandleFunc("/ordertypes/{id:[0-9]+}", apg.HandleGetOrderType).Methods("GET", "OPTIONS")
	route.HandleFunc("/distributionzones", apg.HandleDistributionZones).Methods("GET", "OPTIONS")
	route.HandleFunc("/distributionzones_add", apg.HandleAddDistributionZone).Methods("POST", "OPTIONS")
	route.HandleFunc("/distributionzones_upd", apg.HandleUpdDistributionZone).Methods("POST", "OPTIONS")
	route.HandleFunc("/distributionzones_del", apg.HandleDelDistributionZone).Methods("POST", "OPTIONS")
	route.HandleFunc("/distributionzones/{id:[0-9]+}", apg.HandleGetDistributionZone).Methods("GET", "OPTIONS")
	route.HandleFunc("/contracttypes", apg.HandleContractTypes).Methods("GET", "OPTIONS")
	route.HandleFunc("/contracttypes_add", apg.HandleAddContractType).Methods("POST", "OPTIONS")
	route.HandleFunc("/contracttypes_upd", apg.HandleUpdContractType).Methods("POST", "OPTIONS")
	route.HandleFunc("/contracttypes_del", apg.HandleDelContractType).Methods("POST", "OPTIONS")
	route.HandleFunc("/contracttypes/{id:[0-9]+}", apg.HandleGetContractType).Methods("GET", "OPTIONS")

	// n := negroni.New(negroni.HandlerFunc(utils.AuthValidate))
	n := negroni.New(negroni.HandlerFunc(utils.MWSetupResponse))
	// n.Use(negroni.HandlerFunc(utils.AuthValidate))
	// n.UseHandler(route)

	log.Println("Server is listening at http://localhost:8080/")

	//go run . noauth - run without utils.AuthValidate middleware
	//docker run ... -e NOAUTH="TRUE" - run without utils.AuthValidate middleware
	noauth := false
	sna, ok := os.LookupEnv("NOAUTH")

	if !ok {
		noauth = false
	} else {
		noauth = (sna == "TRUE")
	}

	if !((len(os.Args) > 1 && os.Args[1] == "noauth") || noauth) {
		// log.Fatal(http.ListenAndServe(":8080", route))
		n.Use(negroni.HandlerFunc(utils.AuthValidate))
	}
	// } else {
	// 	// log.Fatal(http.ListenAndServe(":8080", n))
	// }

	n.UseHandler(route)
	log.Fatal(http.ListenAndServe(":8080", n))
}

func (s *PG) handleAdmin(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprintf("It's handleAdmin!\n")))

	ctx := context.Background()
	v := ""

	err := s.dbpool.QueryRow(ctx, "SELECT version();").Scan(&v)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write([]byte(fmt.Sprintf(v, "\n")))
	return
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprintf("It's handleRoot!\n")))

	return
}

/*
// handleLogin godoc
// @Summary List user forms
// @Description Get user forms
// @Tags login
// @Produce  json
// @Success 200 {array} []string
// @Failure 500
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Router / [get]*/
func (s *PG) handleLogin(w http.ResponseWriter, r *http.Request) {
	type Claims struct {
		Username string `json:"username"`
		Password string `json:"password"`
		jwt.StandardClaims
	}

	utils.SetupResponse(&w)

	if (*r).Method == "OPTIONS" {
		return
	}

	ctx := context.Background()

	var jwtSecretKey = []byte("jwt_secret_key")

	claims := &Claims{}

	bearerToken := r.Header.Get("Authorization")

	if bearerToken == "" {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(fmt.Sprintf("Token is empty!\n")))

		return
	}

	tokenString := strings.Split(bearerToken, " ")[1]

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecretKey, nil
	})
	if err != nil {
		http.Error(w, "Error parse token (login): "+err.Error(), 500)
		return
	}

	if !token.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(fmt.Sprintf("Error: %s\n", err)))

		return
	}

	// out_arr := []string{}
	out_arr := []models.LoginForm{}

	rows, err := s.dbpool.Query(ctx, "SELECT * from func_get_user_forms($1);", claims.Username)
	// rows, err := s.dbpool.Query(ctx, "SELECT to_jsonb(func_get_user_forms($1));", claims.Username)

	if err != nil {
		http.Error(w, "Error SELECT * from func_get_user_forms: "+err.Error(), 500)
		return
	}

	defer rows.Close()

	// w.Write([]byte("["))

	for rows.Next() {
		f := models.LoginForm{}
		err = rows.Scan(&f.Form, &f.Rights, &f.UserId)
		// err = rows.Scan(&f)

		// w.Write([]byte(f + "\n"))

		if err != nil {
			log.Println("failed to scan row:", err)
		}

		out_arr = append(out_arr, f)
	}
	// w.Write([]byte("]"))

	output, err := json.Marshal(out_arr)
	if err != nil {
		http.Error(w, "Error marshal output: "+err.Error(), 500)
		return
	}

	w.Write(output)

	return

}
