package routes

import (
	"github.com/gorilla/mux"
	"github.com/u2lentaru/billbck/internal/api"
)

func AddRoutes(r *mux.Router, apg *api.APG) {
	r.HandleFunc("/form_types", apg.HandleFormTypes).Methods("GET", "OPTIONS")
	r.HandleFunc("/form_types_add", apg.HandleAddFormType).Methods("POST", "OPTIONS")
	r.HandleFunc("/form_types_upd", apg.HandleUpdFormType).Methods("POST", "OPTIONS")
	r.HandleFunc("/form_types_del", apg.HandleDelFormType).Methods("POST", "OPTIONS")
	r.HandleFunc("/form_types/{id:[0-9]+}", apg.HandleGetFormType).Methods("GET", "OPTIONS")
	// r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler).Methods("GET", "OPTIONS")
	r.HandleFunc("/sub_types", apg.HandleSubTypes).Methods("GET", "OPTIONS")
	r.HandleFunc("/sub_types_add", apg.HandleAddSubType).Methods("POST", "OPTIONS")
	r.HandleFunc("/sub_types_upd", apg.HandleUpdSubType).Methods("POST", "OPTIONS")
	r.HandleFunc("/sub_types_del", apg.HandleDelSubType).Methods("POST", "OPTIONS")
	r.HandleFunc("/sub_types/{id:[0-9]+}", apg.HandleGetSubType).Methods("GET", "OPTIONS")
	r.HandleFunc("/subjects", apg.HandleSubjects).Methods("GET", "OPTIONS")
	r.HandleFunc("/subjects_add", apg.HandleAddSubject).Methods("POST", "OPTIONS")
	r.HandleFunc("/subjects_upd", apg.HandleUpdSubject).Methods("POST", "OPTIONS")
	r.HandleFunc("/subjects_del", apg.HandleDelSubject).Methods("POST", "OPTIONS")
	r.HandleFunc("/subjects/{id:[0-9]+}", apg.HandleGetSubject).Methods("GET", "OPTIONS")
	r.HandleFunc("/subjects_hist/{id:[0-9]+}", apg.HandleGetSubjectHist).Methods("GET", "OPTIONS")
	r.HandleFunc("/positions", apg.HandlePositions).Methods("GET", "OPTIONS")
	r.HandleFunc("/positions_add", apg.HandleAddPosition).Methods("POST", "OPTIONS")
	r.HandleFunc("/positions_upd", apg.HandleUpdPosition).Methods("POST", "OPTIONS")
	r.HandleFunc("/positions_del", apg.HandleDelPosition).Methods("POST", "OPTIONS")
	r.HandleFunc("/positions/{id:[0-9]+}", apg.HandleGetPosition).Methods("GET", "OPTIONS")
	r.HandleFunc("/banks", api.HandleBanks).Methods("GET", "OPTIONS")
	r.HandleFunc("/banks_add", api.HandleAddBank).Methods("POST", "OPTIONS")
	r.HandleFunc("/banks_upd", api.HandleUpdBank).Methods("POST", "OPTIONS")
	r.HandleFunc("/banks_del", api.HandleDelBank).Methods("POST", "OPTIONS")
	r.HandleFunc("/banks/{id:[0-9]+}", api.HandleGetBank).Methods("GET", "OPTIONS")
	r.HandleFunc("/org_info", apg.HandleOrgInfos).Methods("GET", "OPTIONS")
	r.HandleFunc("/org_info_add", apg.HandleAddOrgInfo).Methods("POST", "OPTIONS")
	r.HandleFunc("/org_info_upd", apg.HandleUpdOrgInfo).Methods("POST", "OPTIONS")
	r.HandleFunc("/org_info_del", apg.HandleDelOrgInfo).Methods("POST", "OPTIONS")
	r.HandleFunc("/org_info/{id:[0-9]+}", apg.HandleGetOrgInfo).Methods("GET", "OPTIONS")
	r.HandleFunc("/sub_banks", apg.HandleSubBanks).Methods("GET", "OPTIONS")
	r.HandleFunc("/sub_banks_add", apg.HandleAddSubBank).Methods("POST", "OPTIONS")
	r.HandleFunc("/sub_banks_upd", apg.HandleUpdSubBank).Methods("POST", "OPTIONS")
	r.HandleFunc("/sub_banks_del", apg.HandleDelSubBank).Methods("POST", "OPTIONS")
	r.HandleFunc("/sub_banks/{id:[0-9]+}", apg.HandleGetSubBank).Methods("GET", "OPTIONS")
	r.HandleFunc("/sub_banks_setactive/{id:[0-9]+}", apg.HandleGetSubBankSetActive).Methods("POST", "OPTIONS")
	r.HandleFunc("/building_types", api.HandleBuildingTypes).Methods("GET", "OPTIONS")
	r.HandleFunc("/building_types_add", api.HandleAddBuildingType).Methods("POST", "OPTIONS")
	r.HandleFunc("/building_types_upd", api.HandleUpdBuildingType).Methods("POST", "OPTIONS")
	r.HandleFunc("/building_types_del", api.HandleDelBuildingType).Methods("POST", "OPTIONS")
	r.HandleFunc("/building_types/{id:[0-9]+}", api.HandleGetBuildingType).Methods("GET", "OPTIONS")
	r.HandleFunc("/areas", api.HandleAreas).Methods("GET", "OPTIONS")
	r.HandleFunc("/areas_add", api.HandleAddArea).Methods("POST", "OPTIONS")
	r.HandleFunc("/areas_upd", api.HandleUpdArea).Methods("POST", "OPTIONS")
	r.HandleFunc("/areas_del", api.HandleDelArea).Methods("POST", "OPTIONS")
	r.HandleFunc("/areas/{id:[0-9]+}", api.HandleGetArea).Methods("GET", "OPTIONS")
	r.HandleFunc("/ksk", apg.HandleKsk).Methods("GET", "OPTIONS")
	r.HandleFunc("/ksk_add", apg.HandleAddKsk).Methods("POST", "OPTIONS")
	r.HandleFunc("/ksk_upd", apg.HandleUpdKsk).Methods("POST", "OPTIONS")
	r.HandleFunc("/ksk_del", apg.HandleDelKsk).Methods("POST", "OPTIONS")
	r.HandleFunc("/ksk/{id:[0-9]+}", apg.HandleGetKsk).Methods("GET", "OPTIONS")
	r.HandleFunc("/sectors", apg.HandleSectors).Methods("GET", "OPTIONS")
	r.HandleFunc("/sectors_add", apg.HandleAddSector).Methods("POST", "OPTIONS")
	r.HandleFunc("/sectors_upd", apg.HandleUpdSector).Methods("POST", "OPTIONS")
	r.HandleFunc("/sectors_del", apg.HandleDelSector).Methods("POST", "OPTIONS")
	r.HandleFunc("/sectors/{id:[0-9]+}", apg.HandleGetSector).Methods("GET", "OPTIONS")
	r.HandleFunc("/connectors", api.HandleConnectors).Methods("GET", "OPTIONS")
	r.HandleFunc("/connectors_add", api.HandleAddConnector).Methods("POST", "OPTIONS")
	r.HandleFunc("/connectors_upd", api.HandleUpdConnector).Methods("POST", "OPTIONS")
	r.HandleFunc("/connectors_del", api.HandleDelConnector).Methods("POST", "OPTIONS")
	r.HandleFunc("/connectors/{id:[0-9]+}", api.HandleGetConnector).Methods("GET", "OPTIONS")
	r.HandleFunc("/input_types", apg.HandleInputTypes).Methods("GET", "OPTIONS")
	r.HandleFunc("/input_types_add", apg.HandleAddInputType).Methods("POST", "OPTIONS")
	r.HandleFunc("/input_types_upd", apg.HandleUpdInputType).Methods("POST", "OPTIONS")
	r.HandleFunc("/input_types_del", apg.HandleDelInputType).Methods("POST", "OPTIONS")
	r.HandleFunc("/input_types/{id:[0-9]+}", apg.HandleGetInputType).Methods("GET", "OPTIONS")
	r.HandleFunc("/reliabilities", apg.HandleReliabilities).Methods("GET", "OPTIONS")
	r.HandleFunc("/reliabilities_add", apg.HandleAddReliability).Methods("POST", "OPTIONS")
	r.HandleFunc("/reliabilities_upd", apg.HandleUpdReliability).Methods("POST", "OPTIONS")
	r.HandleFunc("/reliabilities_del", apg.HandleDelReliability).Methods("POST", "OPTIONS")
	r.HandleFunc("/reliabilities/{id:[0-9]+}", apg.HandleGetReliability).Methods("GET", "OPTIONS")
	r.HandleFunc("/voltages", apg.HandleVoltages).Methods("GET", "OPTIONS")
	r.HandleFunc("/voltages_add", apg.HandleAddVoltage).Methods("POST", "OPTIONS")
	r.HandleFunc("/voltages_upd", apg.HandleUpdVoltage).Methods("POST", "OPTIONS")
	r.HandleFunc("/voltages_del", apg.HandleDelVoltage).Methods("POST", "OPTIONS")
	r.HandleFunc("/voltages/{id:[0-9]+}", apg.HandleGetVoltage).Methods("GET", "OPTIONS")
	r.HandleFunc("/eso", apg.HandleEso).Methods("GET", "OPTIONS")
	r.HandleFunc("/eso_add", apg.HandleAddEso).Methods("POST", "OPTIONS")
	r.HandleFunc("/eso_upd", apg.HandleUpdEso).Methods("POST", "OPTIONS")
	r.HandleFunc("/eso_del", apg.HandleDelEso).Methods("POST", "OPTIONS")
	r.HandleFunc("/eso/{id:[0-9]+}", apg.HandleGetEso).Methods("GET", "OPTIONS")
	r.HandleFunc("/customergroups", api.HandleCustomerGroups).Methods("GET", "OPTIONS")
	r.HandleFunc("/customergroups_add", api.HandleAddCustomerGroup).Methods("POST", "OPTIONS")
	r.HandleFunc("/customergroups_upd", api.HandleUpdCustomerGroup).Methods("POST", "OPTIONS")
	r.HandleFunc("/customergroups_del", api.HandleDelCustomerGroup).Methods("POST", "OPTIONS")
	r.HandleFunc("/customergroups/{id:[0-9]+}", api.HandleGetCustomerGroup).Methods("GET", "OPTIONS")
	r.HandleFunc("/objtypes", apg.HandleObjTypes).Methods("GET", "OPTIONS")
	r.HandleFunc("/objtypes_add", apg.HandleAddObjType).Methods("POST", "OPTIONS")
	r.HandleFunc("/objtypes_upd", apg.HandleUpdObjType).Methods("POST", "OPTIONS")
	r.HandleFunc("/objtypes_del", apg.HandleDelObjType).Methods("POST", "OPTIONS")
	r.HandleFunc("/objtypes/{id:[0-9]+}", apg.HandleGetObjType).Methods("GET", "OPTIONS")
	r.HandleFunc("/uzo", apg.HandleUzo).Methods("GET", "OPTIONS")
	r.HandleFunc("/uzo_add", apg.HandleAddUzo).Methods("POST", "OPTIONS")
	r.HandleFunc("/uzo_upd", apg.HandleUpdUzo).Methods("POST", "OPTIONS")
	r.HandleFunc("/uzo_del", apg.HandleDelUzo).Methods("POST", "OPTIONS")
	r.HandleFunc("/uzo/{id:[0-9]+}", apg.HandleGetUzo).Methods("GET", "OPTIONS")
	r.HandleFunc("/putypes", apg.HandlePuTypes).Methods("GET", "OPTIONS")
	r.HandleFunc("/putypes_add", apg.HandleAddPuType).Methods("POST", "OPTIONS")
	r.HandleFunc("/putypes_upd", apg.HandleUpdPuType).Methods("POST", "OPTIONS")
	r.HandleFunc("/putypes_del", apg.HandleDelPuType).Methods("POST", "OPTIONS")
	r.HandleFunc("/putypes/{id:[0-9]+}", apg.HandleGetPuType).Methods("GET", "OPTIONS")
	r.HandleFunc("/acttypes", api.HandleActTypes).Methods("GET", "OPTIONS")
	r.HandleFunc("/acttypes_add", api.HandleAddActType).Methods("POST", "OPTIONS")
	r.HandleFunc("/acttypes_upd", api.HandleUpdActType).Methods("POST", "OPTIONS")
	r.HandleFunc("/acttypes_del", api.HandleDelActType).Methods("POST", "OPTIONS")
	r.HandleFunc("/acttypes/{id:[0-9]+}", api.HandleGetActType).Methods("GET", "OPTIONS")
	r.HandleFunc("/cashdesks", api.HandleCashdesks).Methods("GET", "OPTIONS")
	r.HandleFunc("/cashdesks_add", api.HandleAddCashdesk).Methods("POST", "OPTIONS")
	r.HandleFunc("/cashdesks_upd", api.HandleUpdCashdesk).Methods("POST", "OPTIONS")
	r.HandleFunc("/cashdesks_del", api.HandleDelCashdesk).Methods("POST", "OPTIONS")
	r.HandleFunc("/cashdesks/{id:[0-9]+}", api.HandleGetCashdesk).Methods("GET", "OPTIONS")
	r.HandleFunc("/paymenttypes", apg.HandlePaymentTypes).Methods("GET", "OPTIONS")
	r.HandleFunc("/paymenttypes_add", apg.HandleAddPaymentType).Methods("POST", "OPTIONS")
	r.HandleFunc("/paymenttypes_upd", apg.HandleUpdPaymentType).Methods("POST", "OPTIONS")
	r.HandleFunc("/paymenttypes_del", apg.HandleDelPaymentType).Methods("POST", "OPTIONS")
	r.HandleFunc("/paymenttypes/{id:[0-9]+}", apg.HandleGetPaymentType).Methods("GET", "OPTIONS")
	r.HandleFunc("/chargetypes", api.HandleChargeTypes).Methods("GET", "OPTIONS")
	r.HandleFunc("/chargetypes_add", api.HandleAddChargeType).Methods("POST", "OPTIONS")
	r.HandleFunc("/chargetypes_upd", api.HandleUpdChargeType).Methods("POST", "OPTIONS")
	r.HandleFunc("/chargetypes_del", api.HandleDelChargeType).Methods("POST", "OPTIONS")
	r.HandleFunc("/chargetypes/{id:[0-9]+}", api.HandleGetChargeType).Methods("GET", "OPTIONS")
	r.HandleFunc("/tariffgroups", apg.HandleTariffGroups).Methods("GET", "OPTIONS")
	r.HandleFunc("/tariffgroups_add", apg.HandleAddTariffGroup).Methods("POST", "OPTIONS")
	r.HandleFunc("/tariffgroups_upd", apg.HandleUpdTariffGroup).Methods("POST", "OPTIONS")
	r.HandleFunc("/tariffgroups_del", apg.HandleDelTariffGroup).Methods("POST", "OPTIONS")
	r.HandleFunc("/tariffgroups/{id:[0-9]+}", apg.HandleGetTariffGroup).Methods("GET", "OPTIONS")
	r.HandleFunc("/tariffs", apg.HandleTariffs).Methods("GET", "OPTIONS")
	r.HandleFunc("/tariffs_add", apg.HandleAddTariff).Methods("POST", "OPTIONS")
	r.HandleFunc("/tariffs_upd", apg.HandleUpdTariff).Methods("POST", "OPTIONS")
	r.HandleFunc("/tariffs_del", apg.HandleDelTariff).Methods("POST", "OPTIONS")
	r.HandleFunc("/tariffs/{id:[0-9]+}", apg.HandleGetTariff).Methods("GET", "OPTIONS")
	r.HandleFunc("/rp", apg.HandleRp).Methods("GET", "OPTIONS")
	r.HandleFunc("/rp_add", apg.HandleAddRp).Methods("POST", "OPTIONS")
	r.HandleFunc("/rp_upd", apg.HandleUpdRp).Methods("POST", "OPTIONS")
	r.HandleFunc("/rp_del", apg.HandleDelRp).Methods("POST", "OPTIONS")
	r.HandleFunc("/rp/{id:[0-9]+}", apg.HandleGetRp).Methods("GET", "OPTIONS")
	r.HandleFunc("/tp", apg.HandleTp).Methods("GET", "OPTIONS")
	r.HandleFunc("/tp_add", apg.HandleAddTp).Methods("POST", "OPTIONS")
	r.HandleFunc("/tp_upd", apg.HandleUpdTp).Methods("POST", "OPTIONS")
	r.HandleFunc("/tp_del", apg.HandleDelTp).Methods("POST", "OPTIONS")
	r.HandleFunc("/tp/{id:[0-9]+}", apg.HandleGetTp).Methods("GET", "OPTIONS")
	r.HandleFunc("/grp", apg.HandleGRp).Methods("GET", "OPTIONS")
	r.HandleFunc("/grp_add", apg.HandleAddGRp).Methods("POST", "OPTIONS")
	r.HandleFunc("/grp_upd", apg.HandleUpdGRp).Methods("POST", "OPTIONS")
	r.HandleFunc("/grp_del", apg.HandleDelGRp).Methods("POST", "OPTIONS")
	r.HandleFunc("/grp/{id:[0-9]+}", apg.HandleGetGRp).Methods("GET", "OPTIONS")
	r.HandleFunc("/streets", apg.HandleStreets).Methods("GET", "OPTIONS")
	r.HandleFunc("/streets_add", apg.HandleAddStreet).Methods("POST", "OPTIONS")
	r.HandleFunc("/streets_upd", apg.HandleUpdStreet).Methods("POST", "OPTIONS")
	r.HandleFunc("/streets_del", apg.HandleDelStreet).Methods("POST", "OPTIONS")
	r.HandleFunc("/streets/{id:[0-9]+}", apg.HandleGetStreet).Methods("GET", "OPTIONS")
	r.HandleFunc("/houses", apg.HandleHouses).Methods("GET", "OPTIONS")
	r.HandleFunc("/houses_add", apg.HandleAddHouse).Methods("POST", "OPTIONS")
	r.HandleFunc("/houses_upd", apg.HandleUpdHouse).Methods("POST", "OPTIONS")
	r.HandleFunc("/houses_del", apg.HandleDelHouse).Methods("POST", "OPTIONS")
	r.HandleFunc("/houses/{id:[0-9]+}", apg.HandleGetHouse).Methods("GET", "OPTIONS")
	r.HandleFunc("/cities", api.HandleCities).Methods("GET", "OPTIONS")
	r.HandleFunc("/cities_add", api.HandleAddCity).Methods("POST", "OPTIONS")
	r.HandleFunc("/cities_upd", api.HandleUpdCity).Methods("POST", "OPTIONS")
	r.HandleFunc("/cities_del", api.HandleDelCity).Methods("POST", "OPTIONS")
	r.HandleFunc("/cities/{id:[0-9]+}", api.HandleGetCity).Methods("GET", "OPTIONS")
	r.HandleFunc("/contracts", api.HandleContracts).Methods("GET", "OPTIONS")
	r.HandleFunc("/contracts_add", api.HandleAddContract).Methods("POST", "OPTIONS")
	r.HandleFunc("/contracts_upd", api.HandleUpdContract).Methods("POST", "OPTIONS")
	r.HandleFunc("/contracts_del", api.HandleDelContract).Methods("POST", "OPTIONS")
	r.HandleFunc("/contracts/{id:[0-9]+}", api.HandleGetContract).Methods("GET", "OPTIONS")
	r.HandleFunc("/contracts_getobject/{id:[0-9]+}", api.HandleGetContractObject).Methods("GET", "OPTIONS")
	r.HandleFunc("/contracts_hist/{id:[0-9]+}", api.HandleGetContractHist).Methods("GET", "OPTIONS")
	r.HandleFunc("/objects", apg.HandleObjects).Methods("GET", "OPTIONS")
	r.HandleFunc("/objects_add", apg.HandleAddObject).Methods("POST", "OPTIONS")
	r.HandleFunc("/objects_upd", apg.HandleUpdObject).Methods("POST", "OPTIONS")
	r.HandleFunc("/objects_del", apg.HandleDelObject).Methods("POST", "OPTIONS")
	r.HandleFunc("/objects/{id:[0-9]+}", apg.HandleGetObject).Methods("GET", "OPTIONS")
	r.HandleFunc("/objects_getcontract/{id:[0-9]+}", apg.HandleGetObjectContract).Methods("GET", "OPTIONS")
	r.HandleFunc("/objects_mff/{hid:[0-9]+}", apg.HandleGetObjectMff).Methods("GET", "OPTIONS")
	r.HandleFunc("/objcontracts", apg.HandleObjContracts).Methods("GET", "OPTIONS")
	r.HandleFunc("/objcontracts_add", apg.HandleAddObjContract).Methods("POST", "OPTIONS")
	r.HandleFunc("/objcontracts_upd", apg.HandleUpdObjContract).Methods("POST", "OPTIONS")
	r.HandleFunc("/objcontracts_del", apg.HandleDelObjContract).Methods("POST", "OPTIONS")
	r.HandleFunc("/objcontracts/{id:[0-9]+}", apg.HandleGetObjContract).Methods("GET", "OPTIONS")
	r.HandleFunc("/acts", api.HandleActs).Methods("GET", "OPTIONS")
	r.HandleFunc("/acts_add", api.HandleAddAct).Methods("POST", "OPTIONS")
	r.HandleFunc("/acts_upd", api.HandleUpdAct).Methods("POST", "OPTIONS")
	r.HandleFunc("/acts_del", api.HandleDelAct).Methods("POST", "OPTIONS")
	r.HandleFunc("/acts_activate", api.HandleActActivate).Methods("GET", "OPTIONS")
	r.HandleFunc("/acts/{id:[0-9]+}", apg.HandleGetAct).Methods("GET", "OPTIONS")
	r.HandleFunc("/actdetails", api.HandleActDetails).Methods("GET", "OPTIONS")
	r.HandleFunc("/actdetails_add", api.HandleAddActDetail).Methods("POST", "OPTIONS")
	r.HandleFunc("/actdetails_upd", api.HandleUpdActDetail).Methods("POST", "OPTIONS")
	r.HandleFunc("/actdetails_del", api.HandleDelActDetail).Methods("POST", "OPTIONS")
	r.HandleFunc("/actdetails/{id:[0-9]+}", api.HandleGetActDetail).Methods("GET", "OPTIONS")
	r.HandleFunc("/pu", apg.HandlePu).Methods("GET", "OPTIONS")
	r.HandleFunc("/pu_add", apg.HandleAddPu).Methods("POST", "OPTIONS")
	r.HandleFunc("/pu_obj_add", apg.HandleAddPu).Methods("POST", "OPTIONS")
	r.HandleFunc("/pu_upd", apg.HandleUpdPu).Methods("POST", "OPTIONS")
	r.HandleFunc("/pu_del", apg.HandleDelPu).Methods("POST", "OPTIONS")
	r.HandleFunc("/pu/{id:[0-9]+}", apg.HandleGetPu).Methods("GET", "OPTIONS")
	r.HandleFunc("/pu_obj", apg.HandlePuObj).Methods("GET", "OPTIONS")
	r.HandleFunc("/puvalues", apg.HandlePuValues).Methods("GET", "OPTIONS")
	r.HandleFunc("/puvalues_add", apg.HandleAddPuValue).Methods("POST", "OPTIONS")
	r.HandleFunc("/puvalues_upd", apg.HandleUpdPuValue).Methods("POST", "OPTIONS")
	r.HandleFunc("/puvalues_del", apg.HandleDelPuValue).Methods("POST", "OPTIONS")
	r.HandleFunc("/puvalues/{id:[0-9]+}", apg.HandleGetPuValue).Methods("GET", "OPTIONS")
	r.HandleFunc("/puvalues_askue_prev", apg.HandlePuValuesAskuePreview).Methods("POST", "OPTIONS")
	r.HandleFunc("/puvalues_askue", apg.HandlePuValuesAskue).Methods("POST", "OPTIONS")
	r.HandleFunc("/balance", apg.HandleBalance).Methods("GET", "OPTIONS")
	r.HandleFunc("/balance/{id:[0-9]+}/{tid:[0-9]+}", apg.HandleGetBalance).Methods("GET", "OPTIONS")
	r.HandleFunc("/balance_sum", apg.HandleBalanceSum).Methods("GET", "OPTIONS")
	r.HandleFunc("/balance_sum1", apg.HandleBalanceSum1).Methods("GET", "OPTIONS")
	r.HandleFunc("/balance_sum0", apg.HandleBalanceSum0).Methods("GET", "OPTIONS")
	r.HandleFunc("/balance_tab1", apg.HandleBalanceTab1).Methods("GET", "OPTIONS")
	r.HandleFunc("/balance_tab0", apg.HandleBalanceTab0).Methods("GET", "OPTIONS")
	r.HandleFunc("/balance_branch", apg.HandleBalanceBranch).Methods("GET", "OPTIONS")
	r.HandleFunc("/subpu", apg.HandleSubPu).Methods("GET", "OPTIONS")
	r.HandleFunc("/subpu_add", apg.HandleAddSubPu).Methods("POST", "OPTIONS")
	r.HandleFunc("/subpu_upd", apg.HandleUpdSubPu).Methods("POST", "OPTIONS")
	r.HandleFunc("/subpu_del", apg.HandleDelSubPu).Methods("POST", "OPTIONS")
	r.HandleFunc("/subpu/{id:[0-9]+}", apg.HandleGetSubPu).Methods("GET", "OPTIONS")
	r.HandleFunc("/subpu_prl", apg.HandlePrlSubPu).Methods("GET", "OPTIONS")
	r.HandleFunc("/staff", apg.HandleStaff).Methods("GET", "OPTIONS")
	r.HandleFunc("/staff_add", apg.HandleAddStaff).Methods("POST", "OPTIONS")
	r.HandleFunc("/staff_upd", apg.HandleUpdStaff).Methods("POST", "OPTIONS")
	r.HandleFunc("/staff_del", apg.HandleDelStaff).Methods("POST", "OPTIONS")
	r.HandleFunc("/staff/{id:[0-9]+}", apg.HandleGetStaff).Methods("GET", "OPTIONS")
	r.HandleFunc("/tgu", apg.HandleTgu).Methods("GET", "OPTIONS")
	r.HandleFunc("/tgu_add", apg.HandleAddTgu).Methods("POST", "OPTIONS")
	r.HandleFunc("/tgu_upd", apg.HandleUpdTgu).Methods("POST", "OPTIONS")
	r.HandleFunc("/tgu_del", apg.HandleDelTgu).Methods("POST", "OPTIONS")
	r.HandleFunc("/tgu/{id:[0-9]+}", apg.HandleGetTgu).Methods("GET", "OPTIONS")
	r.HandleFunc("/conclusions", api.HandleConclusions).Methods("GET", "OPTIONS")
	r.HandleFunc("/conclusions_add", api.HandleAddConclusion).Methods("POST", "OPTIONS")
	r.HandleFunc("/conclusions_upd", api.HandleUpdConclusion).Methods("POST", "OPTIONS")
	r.HandleFunc("/conclusions_del", api.HandleDelConclusion).Methods("POST", "OPTIONS")
	r.HandleFunc("/conclusions/{id:[0-9]+}", api.HandleGetConclusion).Methods("GET", "OPTIONS")
	r.HandleFunc("/shutdowntypes", apg.HandleShutdownTypes).Methods("GET", "OPTIONS")
	r.HandleFunc("/shutdowntypes_add", apg.HandleAddShutdownType).Methods("POST", "OPTIONS")
	r.HandleFunc("/shutdowntypes_upd", apg.HandleUpdShutdownType).Methods("POST", "OPTIONS")
	r.HandleFunc("/shutdowntypes_del", apg.HandleDelShutdownType).Methods("POST", "OPTIONS")
	r.HandleFunc("/shutdowntypes/{id:[0-9]+}", apg.HandleGetShutdownType).Methods("GET", "OPTIONS")
	r.HandleFunc("/violations", apg.HandleViolations).Methods("GET", "OPTIONS")
	r.HandleFunc("/violations_add", apg.HandleAddViolation).Methods("POST", "OPTIONS")
	r.HandleFunc("/violations_upd", apg.HandleUpdViolation).Methods("POST", "OPTIONS")
	r.HandleFunc("/violations_del", apg.HandleDelViolation).Methods("POST", "OPTIONS")
	r.HandleFunc("/violations/{id:[0-9]+}", apg.HandleGetViolation).Methods("GET", "OPTIONS")
	r.HandleFunc("/reasons", apg.HandleReasons).Methods("GET", "OPTIONS")
	r.HandleFunc("/reasons_add", apg.HandleAddReason).Methods("POST", "OPTIONS")
	r.HandleFunc("/reasons_upd", apg.HandleUpdReason).Methods("POST", "OPTIONS")
	r.HandleFunc("/reasons_del", apg.HandleDelReason).Methods("POST", "OPTIONS")
	r.HandleFunc("/reasons/{id:[0-9]+}", apg.HandleGetReason).Methods("GET", "OPTIONS")
	r.HandleFunc("/seals", apg.HandleSeals).Methods("GET", "OPTIONS")
	r.HandleFunc("/seals_add", apg.HandleAddSeal).Methods("POST", "OPTIONS")
	r.HandleFunc("/seals_upd", apg.HandleUpdSeal).Methods("POST", "OPTIONS")
	r.HandleFunc("/seals_del", apg.HandleDelSeal).Methods("POST", "OPTIONS")
	r.HandleFunc("/seals/{id:[0-9]+}", apg.HandleGetSeal).Methods("GET", "OPTIONS")
	r.HandleFunc("/sealtypes", apg.HandleSealTypes).Methods("GET", "OPTIONS")
	r.HandleFunc("/sealtypes_add", apg.HandleAddSealType).Methods("POST", "OPTIONS")
	r.HandleFunc("/sealtypes_upd", apg.HandleUpdSealType).Methods("POST", "OPTIONS")
	r.HandleFunc("/sealtypes_del", apg.HandleDelSealType).Methods("POST", "OPTIONS")
	r.HandleFunc("/sealtypes/{id:[0-9]+}", apg.HandleGetSealType).Methods("GET", "OPTIONS")
	r.HandleFunc("/sealcolours", apg.HandleSealColours).Methods("GET", "OPTIONS")
	r.HandleFunc("/sealcolours_add", apg.HandleAddSealColour).Methods("POST", "OPTIONS")
	r.HandleFunc("/sealcolours_upd", apg.HandleUpdSealColour).Methods("POST", "OPTIONS")
	r.HandleFunc("/sealcolours_del", apg.HandleDelSealColour).Methods("POST", "OPTIONS")
	r.HandleFunc("/sealcolours/{id:[0-9]+}", apg.HandleGetSealColour).Methods("GET", "OPTIONS")
	r.HandleFunc("/sealstatuses", apg.HandleSealStatuses).Methods("GET", "OPTIONS")
	r.HandleFunc("/sealstatuses_add", apg.HandleAddSealStatus).Methods("POST", "OPTIONS")
	r.HandleFunc("/sealstatuses_upd", apg.HandleUpdSealStatus).Methods("POST", "OPTIONS")
	r.HandleFunc("/sealstatuses_del", apg.HandleDelSealStatus).Methods("POST", "OPTIONS")
	r.HandleFunc("/sealstatuses/{id:[0-9]+}", apg.HandleGetSealStatus).Methods("GET", "OPTIONS")
	r.HandleFunc("/charges", api.HandleCharges).Methods("GET", "OPTIONS")
	r.HandleFunc("/charges_add", api.HandleAddCharge).Methods("POST", "OPTIONS")
	r.HandleFunc("/charges_upd", api.HandleUpdCharge).Methods("POST", "OPTIONS")
	r.HandleFunc("/charges_del", api.HandleDelCharge).Methods("POST", "OPTIONS")
	r.HandleFunc("/charges/{id:[0-9]+}", api.HandleGetCharge).Methods("GET", "OPTIONS")
	r.HandleFunc("/charges_run/{id:[0-9]+}", apg.HandleChargeRun).Methods("GET", "OPTIONS")
	r.HandleFunc("/payments", apg.HandlePayments).Methods("GET", "OPTIONS")
	r.HandleFunc("/payments_add", apg.HandleAddPayment).Methods("POST", "OPTIONS")
	r.HandleFunc("/payments_upd", apg.HandleUpdPayment).Methods("POST", "OPTIONS")
	r.HandleFunc("/payments_del", apg.HandleDelPayment).Methods("POST", "OPTIONS")
	r.HandleFunc("/payments/{id:[0-9]+}", apg.HandleGetPayment).Methods("GET", "OPTIONS")
	r.HandleFunc("/equipmenttypes", apg.HandleEquipmentTypes).Methods("GET", "OPTIONS")
	r.HandleFunc("/equipmenttypes_add", apg.HandleAddEquipmentType).Methods("POST", "OPTIONS")
	r.HandleFunc("/equipmenttypes_upd", apg.HandleUpdEquipmentType).Methods("POST", "OPTIONS")
	r.HandleFunc("/equipmenttypes_del", apg.HandleDelEquipmentType).Methods("POST", "OPTIONS")
	r.HandleFunc("/equipmenttypes/{id:[0-9]+}", apg.HandleGetEquipmentType).Methods("GET", "OPTIONS")
	r.HandleFunc("/equipment", apg.HandleEquipment).Methods("GET", "OPTIONS")
	r.HandleFunc("/equipment_add", apg.HandleAddEquipment).Methods("POST", "OPTIONS")
	r.HandleFunc("/equipment_addlist", apg.HandleAddEquipmentList).Methods("POST", "OPTIONS")
	r.HandleFunc("/equipment_upd", apg.HandleUpdEquipment).Methods("POST", "OPTIONS")
	r.HandleFunc("/equipment_del", apg.HandleDelEquipment).Methods("POST", "OPTIONS")
	r.HandleFunc("/equipment_delobj", apg.HandleDelObjEquipment).Methods("POST", "OPTIONS")
	r.HandleFunc("/equipment/{id:[0-9]+}", apg.HandleGetEquipment).Methods("GET", "OPTIONS")
	r.HandleFunc("/calculationtypes", api.HandleCalculationTypes).Methods("GET", "OPTIONS")
	r.HandleFunc("/calculationtypes_add", api.HandleAddCalculationType).Methods("POST", "OPTIONS")
	r.HandleFunc("/calculationtypes_upd", api.HandleUpdCalculationType).Methods("POST", "OPTIONS")
	r.HandleFunc("/calculationtypes_del", api.HandleDelCalculationType).Methods("POST", "OPTIONS")
	r.HandleFunc("/calculationtypes/{id:[0-9]+}", api.HandleGetCalculationType).Methods("GET", "OPTIONS")
	r.HandleFunc("/results", apg.HandleResults).Methods("GET", "OPTIONS")
	r.HandleFunc("/results_add", apg.HandleAddResult).Methods("POST", "OPTIONS")
	r.HandleFunc("/results_upd", apg.HandleUpdResult).Methods("POST", "OPTIONS")
	r.HandleFunc("/results_del", apg.HandleDelResult).Methods("POST", "OPTIONS")
	r.HandleFunc("/results/{id:[0-9]+}", apg.HandleGetResult).Methods("GET", "OPTIONS")
	r.HandleFunc("/servicetypes", apg.HandleServiceTypes).Methods("GET", "OPTIONS")
	r.HandleFunc("/servicetypes_add", apg.HandleAddServiceType).Methods("POST", "OPTIONS")
	r.HandleFunc("/servicetypes_upd", apg.HandleUpdServiceType).Methods("POST", "OPTIONS")
	r.HandleFunc("/servicetypes_del", apg.HandleDelServiceType).Methods("POST", "OPTIONS")
	r.HandleFunc("/servicetypes/{id:[0-9]+}", apg.HandleGetServiceType).Methods("GET", "OPTIONS")
	r.HandleFunc("/requesttypes", apg.HandleRequestTypes).Methods("GET", "OPTIONS")
	r.HandleFunc("/requesttypes_add", apg.HandleAddRequestType).Methods("POST", "OPTIONS")
	r.HandleFunc("/requesttypes_upd", apg.HandleUpdRequestType).Methods("POST", "OPTIONS")
	r.HandleFunc("/requesttypes_del", apg.HandleDelRequestType).Methods("POST", "OPTIONS")
	r.HandleFunc("/requesttypes/{id:[0-9]+}", apg.HandleGetRequestType).Methods("GET", "OPTIONS")
	r.HandleFunc("/claimtypes", api.HandleClaimTypes).Methods("GET", "OPTIONS")
	r.HandleFunc("/claimtypes_add", api.HandleAddClaimType).Methods("POST", "OPTIONS")
	r.HandleFunc("/claimtypes_upd", api.HandleUpdClaimType).Methods("POST", "OPTIONS")
	r.HandleFunc("/claimtypes_del", api.HandleDelClaimType).Methods("POST", "OPTIONS")
	r.HandleFunc("/claimtypes/{id:[0-9]+}", api.HandleGetClaimType).Methods("GET", "OPTIONS")
	r.HandleFunc("/requests", apg.HandleRequests).Methods("GET", "OPTIONS")
	r.HandleFunc("/requests_add", apg.HandleAddRequest).Methods("POST", "OPTIONS")
	r.HandleFunc("/requests_upd", apg.HandleUpdRequest).Methods("POST", "OPTIONS")
	r.HandleFunc("/requests_del", apg.HandleDelRequest).Methods("POST", "OPTIONS")
	r.HandleFunc("/requests/{id:[0-9]+}", apg.HandleGetRequest).Methods("GET", "OPTIONS")
	r.HandleFunc("/objstatuses", apg.HandleObjStatuses).Methods("GET", "OPTIONS")
	r.HandleFunc("/objstatuses_add", apg.HandleAddObjStatus).Methods("POST", "OPTIONS")
	r.HandleFunc("/objstatuses_upd", apg.HandleUpdObjStatus).Methods("POST", "OPTIONS")
	r.HandleFunc("/objstatuses_del", apg.HandleDelObjStatus).Methods("POST", "OPTIONS")
	r.HandleFunc("/objstatuses/{id:[0-9]+}", apg.HandleGetObjStatus).Methods("GET", "OPTIONS")
	r.HandleFunc("/users", apg.HandleUsers).Methods("GET", "OPTIONS")
	r.HandleFunc("/users_add", apg.HandleAddUser).Methods("POST", "OPTIONS")
	r.HandleFunc("/users_upd", apg.HandleUpdUser).Methods("POST", "OPTIONS")
	r.HandleFunc("/users_del", apg.HandleDelUser).Methods("POST", "OPTIONS")
	r.HandleFunc("/users/{id:[0-9]+}", apg.HandleGetUser).Methods("GET", "OPTIONS")
	r.HandleFunc("/contractmots", api.HandleContractMots).Methods("GET", "OPTIONS")
	r.HandleFunc("/contractmots_add", api.HandleAddContractMot).Methods("POST", "OPTIONS")
	r.HandleFunc("/contractmots_upd", api.HandleUpdContractMot).Methods("POST", "OPTIONS")
	r.HandleFunc("/contractmots_del", api.HandleDelContractMot).Methods("POST", "OPTIONS")
	r.HandleFunc("/contractmots/{id:[0-9]+}", api.HandleGetContractMot).Methods("GET", "OPTIONS")
	r.HandleFunc("/requestkinds", apg.HandleRequestKinds).Methods("GET", "OPTIONS")
	r.HandleFunc("/requestkinds_add", apg.HandleAddRequestKind).Methods("POST", "OPTIONS")
	r.HandleFunc("/requestkinds_upd", apg.HandleUpdRequestKind).Methods("POST", "OPTIONS")
	r.HandleFunc("/requestkinds_del", apg.HandleDelRequestKind).Methods("POST", "OPTIONS")
	r.HandleFunc("/requestkinds/{id:[0-9]+}", apg.HandleGetRequestKind).Methods("GET", "OPTIONS")
	r.HandleFunc("/periods", apg.HandlePeriods).Methods("GET", "OPTIONS")
	r.HandleFunc("/periods_add", apg.HandleAddPeriod).Methods("POST", "OPTIONS")
	r.HandleFunc("/periods_upd", apg.HandleUpdPeriod).Methods("POST", "OPTIONS")
	r.HandleFunc("/periods_del", apg.HandleDelPeriod).Methods("POST", "OPTIONS")
	r.HandleFunc("/periods/{id:[0-9]+}", apg.HandleGetPeriod).Methods("GET", "OPTIONS")
	r.HandleFunc("/transtypes", apg.HandleTransTypes).Methods("GET", "OPTIONS")
	r.HandleFunc("/transtypes_add", apg.HandleAddTransType).Methods("POST", "OPTIONS")
	r.HandleFunc("/transtypes_upd", apg.HandleUpdTransType).Methods("POST", "OPTIONS")
	r.HandleFunc("/transtypes_del", apg.HandleDelTransType).Methods("POST", "OPTIONS")
	r.HandleFunc("/transtypes/{id:[0-9]+}", apg.HandleGetTransType).Methods("GET", "OPTIONS")
	r.HandleFunc("/transcurr", apg.HandleTransCurr).Methods("GET", "OPTIONS")
	r.HandleFunc("/transcurr_add", apg.HandleAddTransCurr).Methods("POST", "OPTIONS")
	r.HandleFunc("/transcurr_upd", apg.HandleUpdTransCurr).Methods("POST", "OPTIONS")
	r.HandleFunc("/transcurr_del", apg.HandleDelTransCurr).Methods("POST", "OPTIONS")
	r.HandleFunc("/transcurr/{id:[0-9]+}", apg.HandleGetTransCurr).Methods("GET", "OPTIONS")
	r.HandleFunc("/transvolt", apg.HandleTransVolt).Methods("GET", "OPTIONS")
	r.HandleFunc("/transvolt_add", apg.HandleAddTransVolt).Methods("POST", "OPTIONS")
	r.HandleFunc("/transvolt_upd", apg.HandleUpdTransVolt).Methods("POST", "OPTIONS")
	r.HandleFunc("/transvolt_del", apg.HandleDelTransVolt).Methods("POST", "OPTIONS")
	r.HandleFunc("/transvolt/{id:[0-9]+}", apg.HandleGetTransVolt).Methods("GET", "OPTIONS")
	r.HandleFunc("/transpwrtypes", apg.HandleTransPwrTypes).Methods("GET", "OPTIONS")
	r.HandleFunc("/transpwrtypes_add", apg.HandleAddTransPwrType).Methods("POST", "OPTIONS")
	r.HandleFunc("/transpwrtypes_upd", apg.HandleUpdTransPwrType).Methods("POST", "OPTIONS")
	r.HandleFunc("/transpwrtypes_del", apg.HandleDelTransPwrType).Methods("POST", "OPTIONS")
	r.HandleFunc("/transpwrtypes/{id:[0-9]+}", apg.HandleGetTransPwrType).Methods("GET", "OPTIONS")
	r.HandleFunc("/transpwr", apg.HandleTransPwr).Methods("GET", "OPTIONS")
	r.HandleFunc("/transpwr_add", apg.HandleAddTransPwr).Methods("POST", "OPTIONS")
	r.HandleFunc("/transpwr_upd", apg.HandleUpdTransPwr).Methods("POST", "OPTIONS")
	r.HandleFunc("/transpwr_del", apg.HandleDelTransPwr).Methods("POST", "OPTIONS")
	r.HandleFunc("/transpwr/{id:[0-9]+}", apg.HandleGetTransPwr).Methods("GET", "OPTIONS")
	r.HandleFunc("/objtranscurr", apg.HandleObjTransCurr).Methods("GET", "OPTIONS")
	r.HandleFunc("/objtranscurr_add", apg.HandleAddObjTransCurr).Methods("POST", "OPTIONS")
	r.HandleFunc("/objtranscurr_upd", apg.HandleUpdObjTransCurr).Methods("POST", "OPTIONS")
	r.HandleFunc("/objtranscurr_del", apg.HandleDelObjTransCurr).Methods("POST", "OPTIONS")
	r.HandleFunc("/objtranscurr/{id:[0-9]+}", apg.HandleGetObjTransCurr).Methods("GET", "OPTIONS")
	r.HandleFunc("/objtranscurr_obj", apg.HandleObjTransCurrByObj).Methods("GET", "OPTIONS")
	r.HandleFunc("/objtransvolt", apg.HandleObjTransVolt).Methods("GET", "OPTIONS")
	r.HandleFunc("/objtransvolt_add", apg.HandleAddObjTransVolt).Methods("POST", "OPTIONS")
	r.HandleFunc("/objtransvolt_upd", apg.HandleUpdObjTransVolt).Methods("POST", "OPTIONS")
	r.HandleFunc("/objtransvolt_del", apg.HandleDelObjTransVolt).Methods("POST", "OPTIONS")
	r.HandleFunc("/objtransvolt/{id:[0-9]+}", apg.HandleGetObjTransVolt).Methods("GET", "OPTIONS")
	r.HandleFunc("/objtransvolt_obj", apg.HandleObjTransVoltByObj).Methods("GET", "OPTIONS")
	r.HandleFunc("/objtranspwr", apg.HandleObjTransPwr).Methods("GET", "OPTIONS")
	r.HandleFunc("/objtranspwr_add", apg.HandleAddObjTransPwr).Methods("POST", "OPTIONS")
	r.HandleFunc("/objtranspwr_upd", apg.HandleUpdObjTransPwr).Methods("POST", "OPTIONS")
	r.HandleFunc("/objtranspwr_del", apg.HandleDelObjTransPwr).Methods("POST", "OPTIONS")
	r.HandleFunc("/objtranspwr/{id:[0-9]+}", apg.HandleGetObjTransPwr).Methods("GET", "OPTIONS")
	r.HandleFunc("/objtranspwr_obj", apg.HandleObjTransPwrByObj).Methods("GET", "OPTIONS")
	r.HandleFunc("/cableresistances", api.HandleCableResistances).Methods("GET", "OPTIONS")
	r.HandleFunc("/cableresistances_add", api.HandleAddCableResistance).Methods("POST", "OPTIONS")
	r.HandleFunc("/cableresistances_upd", api.HandleUpdCableResistance).Methods("POST", "OPTIONS")
	r.HandleFunc("/cableresistances_del", api.HandleDelCableResistance).Methods("POST", "OPTIONS")
	r.HandleFunc("/cableresistances/{id:[0-9]+}", api.HandleGetCableResistance).Methods("GET", "OPTIONS")
	r.HandleFunc("/objlines", apg.HandleObjLines).Methods("GET", "OPTIONS")
	r.HandleFunc("/objlines_add", apg.HandleAddObjLine).Methods("POST", "OPTIONS")
	r.HandleFunc("/objlines_upd", apg.HandleUpdObjLine).Methods("POST", "OPTIONS")
	r.HandleFunc("/objlines_del", apg.HandleDelObjLine).Methods("POST", "OPTIONS")
	r.HandleFunc("/objlines/{id:[0-9]+}", apg.HandleGetObjLine).Methods("GET", "OPTIONS")
	r.HandleFunc("/objlines_obj", apg.HandleObjLinesByObj).Methods("GET", "OPTIONS")
	r.HandleFunc("/askuetypes", api.HandleAskueTypes).Methods("GET", "OPTIONS")
	r.HandleFunc("/askuetypes_add", api.HandleAddAskueType).Methods("POST", "OPTIONS")
	r.HandleFunc("/askuetypes_upd", api.HandleUpdAskueType).Methods("POST", "OPTIONS")
	r.HandleFunc("/askuetypes_del", api.HandleDelAskueType).Methods("POST", "OPTIONS")
	r.HandleFunc("/askuetypes/{id:[0-9]+}", api.HandleGetAskueType).Methods("GET", "OPTIONS")
	r.HandleFunc("/puvaluetypes", apg.HandlePuValueTypes).Methods("GET", "OPTIONS")
	r.HandleFunc("/puvaluetypes_add", apg.HandleAddPuValueType).Methods("POST", "OPTIONS")
	r.HandleFunc("/puvaluetypes_upd", apg.HandleUpdPuValueType).Methods("POST", "OPTIONS")
	r.HandleFunc("/puvaluetypes_del", apg.HandleDelPuValueType).Methods("POST", "OPTIONS")
	r.HandleFunc("/puvaluetypes/{id:[0-9]+}", apg.HandleGetPuValueType).Methods("GET", "OPTIONS")
	r.HandleFunc("/ordertypes", apg.HandleOrderTypes).Methods("GET", "OPTIONS")
	r.HandleFunc("/ordertypes_add", apg.HandleAddOrderType).Methods("POST", "OPTIONS")
	r.HandleFunc("/ordertypes_upd", apg.HandleUpdOrderType).Methods("POST", "OPTIONS")
	r.HandleFunc("/ordertypes_del", apg.HandleDelOrderType).Methods("POST", "OPTIONS")
	r.HandleFunc("/ordertypes/{id:[0-9]+}", apg.HandleGetOrderType).Methods("GET", "OPTIONS")
	r.HandleFunc("/distributionzones", apg.HandleDistributionZones).Methods("GET", "OPTIONS")
	r.HandleFunc("/distributionzones_add", apg.HandleAddDistributionZone).Methods("POST", "OPTIONS")
	r.HandleFunc("/distributionzones_upd", apg.HandleUpdDistributionZone).Methods("POST", "OPTIONS")
	r.HandleFunc("/distributionzones_del", apg.HandleDelDistributionZone).Methods("POST", "OPTIONS")
	r.HandleFunc("/distributionzones/{id:[0-9]+}", apg.HandleGetDistributionZone).Methods("GET", "OPTIONS")
	r.HandleFunc("/contracttypes", api.HandleContractTypes).Methods("GET", "OPTIONS")
	r.HandleFunc("/contracttypes_add", api.HandleAddContractType).Methods("POST", "OPTIONS")
	r.HandleFunc("/contracttypes_upd", api.HandleUpdContractType).Methods("POST", "OPTIONS")
	r.HandleFunc("/contracttypes_del", api.HandleDelContractType).Methods("POST", "OPTIONS")
	r.HandleFunc("/contracttypes/{id:[0-9]+}", api.HandleGetContractType).Methods("GET", "OPTIONS")
}
