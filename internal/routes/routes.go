package routes

import (
	"github.com/gorilla/mux"
	"github.com/u2lentaru/billbck/internal/api"
)

func AddRoutes(r *mux.Router, apg *api.APG) {
	r.HandleFunc("/form_types", api.HandleFormTypes).Methods("GET", "OPTIONS")
	r.HandleFunc("/form_types_add", api.HandleAddFormType).Methods("POST", "OPTIONS")
	r.HandleFunc("/form_types_upd", api.HandleUpdFormType).Methods("POST", "OPTIONS")
	r.HandleFunc("/form_types_del", api.HandleDelFormType).Methods("POST", "OPTIONS")
	r.HandleFunc("/form_types/{id:[0-9]+}", api.HandleGetFormType).Methods("GET", "OPTIONS")
	// r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler).Methods("GET", "OPTIONS")
	r.HandleFunc("/sub_types", api.HandleSubTypes).Methods("GET", "OPTIONS")
	r.HandleFunc("/sub_types_add", api.HandleAddSubType).Methods("POST", "OPTIONS")
	r.HandleFunc("/sub_types_upd", api.HandleUpdSubType).Methods("POST", "OPTIONS")
	r.HandleFunc("/sub_types_del", api.HandleDelSubType).Methods("POST", "OPTIONS")
	r.HandleFunc("/sub_types/{id:[0-9]+}", api.HandleGetSubType).Methods("GET", "OPTIONS")
	r.HandleFunc("/subjects", api.HandleSubjects).Methods("GET", "OPTIONS")
	r.HandleFunc("/subjects_add", api.HandleAddSubject).Methods("POST", "OPTIONS")
	r.HandleFunc("/subjects_upd", api.HandleUpdSubject).Methods("POST", "OPTIONS")
	r.HandleFunc("/subjects_del", api.HandleDelSubject).Methods("POST", "OPTIONS")
	r.HandleFunc("/subjects/{id:[0-9]+}", api.HandleGetSubject).Methods("GET", "OPTIONS")
	r.HandleFunc("/subjects_hist/{id:[0-9]+}", api.HandleGetSubjectHist).Methods("GET", "OPTIONS")
	r.HandleFunc("/positions", api.HandlePositions).Methods("GET", "OPTIONS")
	r.HandleFunc("/positions_add", api.HandleAddPosition).Methods("POST", "OPTIONS")
	r.HandleFunc("/positions_upd", api.HandleUpdPosition).Methods("POST", "OPTIONS")
	r.HandleFunc("/positions_del", api.HandleDelPosition).Methods("POST", "OPTIONS")
	r.HandleFunc("/positions/{id:[0-9]+}", api.HandleGetPosition).Methods("GET", "OPTIONS")
	r.HandleFunc("/banks", api.HandleBanks).Methods("GET", "OPTIONS")
	r.HandleFunc("/banks_add", api.HandleAddBank).Methods("POST", "OPTIONS")
	r.HandleFunc("/banks_upd", api.HandleUpdBank).Methods("POST", "OPTIONS")
	r.HandleFunc("/banks_del", api.HandleDelBank).Methods("POST", "OPTIONS")
	r.HandleFunc("/banks/{id:[0-9]+}", api.HandleGetBank).Methods("GET", "OPTIONS")
	r.HandleFunc("/org_info", api.HandleOrgInfos).Methods("GET", "OPTIONS")
	r.HandleFunc("/org_info_add", api.HandleAddOrgInfo).Methods("POST", "OPTIONS")
	r.HandleFunc("/org_info_upd", api.HandleUpdOrgInfo).Methods("POST", "OPTIONS")
	r.HandleFunc("/org_info_del", api.HandleDelOrgInfo).Methods("POST", "OPTIONS")
	r.HandleFunc("/org_info/{id:[0-9]+}", api.HandleGetOrgInfo).Methods("GET", "OPTIONS")
	r.HandleFunc("/sub_banks", api.HandleSubBanks).Methods("GET", "OPTIONS")
	r.HandleFunc("/sub_banks_add", api.HandleAddSubBank).Methods("POST", "OPTIONS")
	r.HandleFunc("/sub_banks_upd", api.HandleUpdSubBank).Methods("POST", "OPTIONS")
	r.HandleFunc("/sub_banks_del", api.HandleDelSubBank).Methods("POST", "OPTIONS")
	r.HandleFunc("/sub_banks/{id:[0-9]+}", api.HandleGetSubBank).Methods("GET", "OPTIONS")
	r.HandleFunc("/sub_banks_setactive/{id:[0-9]+}", api.HandleGetSubBankSetActive).Methods("POST", "OPTIONS")
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
	r.HandleFunc("/ksk", api.HandleKsk).Methods("GET", "OPTIONS")
	r.HandleFunc("/ksk_add", api.HandleAddKsk).Methods("POST", "OPTIONS")
	r.HandleFunc("/ksk_upd", api.HandleUpdKsk).Methods("POST", "OPTIONS")
	r.HandleFunc("/ksk_del", api.HandleDelKsk).Methods("POST", "OPTIONS")
	r.HandleFunc("/ksk/{id:[0-9]+}", api.HandleGetKsk).Methods("GET", "OPTIONS")
	r.HandleFunc("/sectors", api.HandleSectors).Methods("GET", "OPTIONS")
	r.HandleFunc("/sectors_add", api.HandleAddSector).Methods("POST", "OPTIONS")
	r.HandleFunc("/sectors_upd", api.HandleUpdSector).Methods("POST", "OPTIONS")
	r.HandleFunc("/sectors_del", api.HandleDelSector).Methods("POST", "OPTIONS")
	r.HandleFunc("/sectors/{id:[0-9]+}", api.HandleGetSector).Methods("GET", "OPTIONS")
	r.HandleFunc("/connectors", api.HandleConnectors).Methods("GET", "OPTIONS")
	r.HandleFunc("/connectors_add", api.HandleAddConnector).Methods("POST", "OPTIONS")
	r.HandleFunc("/connectors_upd", api.HandleUpdConnector).Methods("POST", "OPTIONS")
	r.HandleFunc("/connectors_del", api.HandleDelConnector).Methods("POST", "OPTIONS")
	r.HandleFunc("/connectors/{id:[0-9]+}", api.HandleGetConnector).Methods("GET", "OPTIONS")
	r.HandleFunc("/input_types", api.HandleInputTypes).Methods("GET", "OPTIONS")
	r.HandleFunc("/input_types_add", api.HandleAddInputType).Methods("POST", "OPTIONS")
	r.HandleFunc("/input_types_upd", api.HandleUpdInputType).Methods("POST", "OPTIONS")
	r.HandleFunc("/input_types_del", api.HandleDelInputType).Methods("POST", "OPTIONS")
	r.HandleFunc("/input_types/{id:[0-9]+}", api.HandleGetInputType).Methods("GET", "OPTIONS")
	r.HandleFunc("/reliabilities", api.HandleReliabilities).Methods("GET", "OPTIONS")
	r.HandleFunc("/reliabilities_add", api.HandleAddReliability).Methods("POST", "OPTIONS")
	r.HandleFunc("/reliabilities_upd", api.HandleUpdReliability).Methods("POST", "OPTIONS")
	r.HandleFunc("/reliabilities_del", api.HandleDelReliability).Methods("POST", "OPTIONS")
	r.HandleFunc("/reliabilities/{id:[0-9]+}", api.HandleGetReliability).Methods("GET", "OPTIONS")
	r.HandleFunc("/voltages", apg.HandleVoltages).Methods("GET", "OPTIONS")
	r.HandleFunc("/voltages_add", apg.HandleAddVoltage).Methods("POST", "OPTIONS")
	r.HandleFunc("/voltages_upd", apg.HandleUpdVoltage).Methods("POST", "OPTIONS")
	r.HandleFunc("/voltages_del", apg.HandleDelVoltage).Methods("POST", "OPTIONS")
	r.HandleFunc("/voltages/{id:[0-9]+}", apg.HandleGetVoltage).Methods("GET", "OPTIONS")
	r.HandleFunc("/eso", api.HandleEso).Methods("GET", "OPTIONS")
	r.HandleFunc("/eso_add", api.HandleAddEso).Methods("POST", "OPTIONS")
	r.HandleFunc("/eso_upd", api.HandleUpdEso).Methods("POST", "OPTIONS")
	r.HandleFunc("/eso_del", api.HandleDelEso).Methods("POST", "OPTIONS")
	r.HandleFunc("/eso/{id:[0-9]+}", api.HandleGetEso).Methods("GET", "OPTIONS")
	r.HandleFunc("/customergroups", api.HandleCustomerGroups).Methods("GET", "OPTIONS")
	r.HandleFunc("/customergroups_add", api.HandleAddCustomerGroup).Methods("POST", "OPTIONS")
	r.HandleFunc("/customergroups_upd", api.HandleUpdCustomerGroup).Methods("POST", "OPTIONS")
	r.HandleFunc("/customergroups_del", api.HandleDelCustomerGroup).Methods("POST", "OPTIONS")
	r.HandleFunc("/customergroups/{id:[0-9]+}", api.HandleGetCustomerGroup).Methods("GET", "OPTIONS")
	r.HandleFunc("/objtypes", api.HandleObjTypes).Methods("GET", "OPTIONS")
	r.HandleFunc("/objtypes_add", api.HandleAddObjType).Methods("POST", "OPTIONS")
	r.HandleFunc("/objtypes_upd", api.HandleUpdObjType).Methods("POST", "OPTIONS")
	r.HandleFunc("/objtypes_del", api.HandleDelObjType).Methods("POST", "OPTIONS")
	r.HandleFunc("/objtypes/{id:[0-9]+}", api.HandleGetObjType).Methods("GET", "OPTIONS")
	r.HandleFunc("/uzo", apg.HandleUzo).Methods("GET", "OPTIONS")
	r.HandleFunc("/uzo_add", apg.HandleAddUzo).Methods("POST", "OPTIONS")
	r.HandleFunc("/uzo_upd", apg.HandleUpdUzo).Methods("POST", "OPTIONS")
	r.HandleFunc("/uzo_del", apg.HandleDelUzo).Methods("POST", "OPTIONS")
	r.HandleFunc("/uzo/{id:[0-9]+}", apg.HandleGetUzo).Methods("GET", "OPTIONS")
	r.HandleFunc("/putypes", api.HandlePuTypes).Methods("GET", "OPTIONS")
	r.HandleFunc("/putypes_add", api.HandleAddPuType).Methods("POST", "OPTIONS")
	r.HandleFunc("/putypes_upd", api.HandleUpdPuType).Methods("POST", "OPTIONS")
	r.HandleFunc("/putypes_del", api.HandleDelPuType).Methods("POST", "OPTIONS")
	r.HandleFunc("/putypes/{id:[0-9]+}", api.HandleGetPuType).Methods("GET", "OPTIONS")
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
	r.HandleFunc("/paymenttypes", api.HandlePaymentTypes).Methods("GET", "OPTIONS")
	r.HandleFunc("/paymenttypes_add", api.HandleAddPaymentType).Methods("POST", "OPTIONS")
	r.HandleFunc("/paymenttypes_upd", api.HandleUpdPaymentType).Methods("POST", "OPTIONS")
	r.HandleFunc("/paymenttypes_del", api.HandleDelPaymentType).Methods("POST", "OPTIONS")
	r.HandleFunc("/paymenttypes/{id:[0-9]+}", api.HandleGetPaymentType).Methods("GET", "OPTIONS")
	r.HandleFunc("/chargetypes", api.HandleChargeTypes).Methods("GET", "OPTIONS")
	r.HandleFunc("/chargetypes_add", api.HandleAddChargeType).Methods("POST", "OPTIONS")
	r.HandleFunc("/chargetypes_upd", api.HandleUpdChargeType).Methods("POST", "OPTIONS")
	r.HandleFunc("/chargetypes_del", api.HandleDelChargeType).Methods("POST", "OPTIONS")
	r.HandleFunc("/chargetypes/{id:[0-9]+}", api.HandleGetChargeType).Methods("GET", "OPTIONS")
	r.HandleFunc("/tariffgroups", api.HandleTariffGroups).Methods("GET", "OPTIONS")
	r.HandleFunc("/tariffgroups_add", api.HandleAddTariffGroup).Methods("POST", "OPTIONS")
	r.HandleFunc("/tariffgroups_upd", api.HandleUpdTariffGroup).Methods("POST", "OPTIONS")
	r.HandleFunc("/tariffgroups_del", api.HandleDelTariffGroup).Methods("POST", "OPTIONS")
	r.HandleFunc("/tariffgroups/{id:[0-9]+}", api.HandleGetTariffGroup).Methods("GET", "OPTIONS")
	r.HandleFunc("/tariffs", api.HandleTariffs).Methods("GET", "OPTIONS")
	r.HandleFunc("/tariffs_add", api.HandleAddTariff).Methods("POST", "OPTIONS")
	r.HandleFunc("/tariffs_upd", api.HandleUpdTariff).Methods("POST", "OPTIONS")
	r.HandleFunc("/tariffs_del", api.HandleDelTariff).Methods("POST", "OPTIONS")
	r.HandleFunc("/tariffs/{id:[0-9]+}", api.HandleGetTariff).Methods("GET", "OPTIONS")
	r.HandleFunc("/rp", api.HandleRp).Methods("GET", "OPTIONS")
	r.HandleFunc("/rp_add", api.HandleAddRp).Methods("POST", "OPTIONS")
	r.HandleFunc("/rp_upd", api.HandleUpdRp).Methods("POST", "OPTIONS")
	r.HandleFunc("/rp_del", api.HandleDelRp).Methods("POST", "OPTIONS")
	r.HandleFunc("/rp/{id:[0-9]+}", api.HandleGetRp).Methods("GET", "OPTIONS")
	r.HandleFunc("/tp", apg.HandleTp).Methods("GET", "OPTIONS")
	r.HandleFunc("/tp_add", apg.HandleAddTp).Methods("POST", "OPTIONS")
	r.HandleFunc("/tp_upd", apg.HandleUpdTp).Methods("POST", "OPTIONS")
	r.HandleFunc("/tp_del", apg.HandleDelTp).Methods("POST", "OPTIONS")
	r.HandleFunc("/tp/{id:[0-9]+}", apg.HandleGetTp).Methods("GET", "OPTIONS")
	r.HandleFunc("/grp", api.HandleGRp).Methods("GET", "OPTIONS")
	r.HandleFunc("/grp_add", api.HandleAddGRp).Methods("POST", "OPTIONS")
	r.HandleFunc("/grp_upd", api.HandleUpdGRp).Methods("POST", "OPTIONS")
	r.HandleFunc("/grp_del", api.HandleDelGRp).Methods("POST", "OPTIONS")
	r.HandleFunc("/grp/{id:[0-9]+}", api.HandleGetGRp).Methods("GET", "OPTIONS")
	r.HandleFunc("/streets", api.HandleStreets).Methods("GET", "OPTIONS")
	r.HandleFunc("/streets_add", api.HandleAddStreet).Methods("POST", "OPTIONS")
	r.HandleFunc("/streets_upd", api.HandleUpdStreet).Methods("POST", "OPTIONS")
	r.HandleFunc("/streets_del", api.HandleDelStreet).Methods("POST", "OPTIONS")
	r.HandleFunc("/streets/{id:[0-9]+}", api.HandleGetStreet).Methods("GET", "OPTIONS")
	r.HandleFunc("/houses", api.HandleHouses).Methods("GET", "OPTIONS")
	r.HandleFunc("/houses_add", api.HandleAddHouse).Methods("POST", "OPTIONS")
	r.HandleFunc("/houses_upd", api.HandleUpdHouse).Methods("POST", "OPTIONS")
	r.HandleFunc("/houses_del", api.HandleDelHouse).Methods("POST", "OPTIONS")
	r.HandleFunc("/houses/{id:[0-9]+}", api.HandleGetHouse).Methods("GET", "OPTIONS")
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
	r.HandleFunc("/objects", api.HandleObjects).Methods("GET", "OPTIONS")
	r.HandleFunc("/objects_add", api.HandleAddObject).Methods("POST", "OPTIONS")
	r.HandleFunc("/objects_upd", api.HandleUpdObject).Methods("POST", "OPTIONS")
	r.HandleFunc("/objects_del", api.HandleDelObject).Methods("POST", "OPTIONS")
	r.HandleFunc("/objects/{id:[0-9]+}", api.HandleGetObject).Methods("GET", "OPTIONS")
	r.HandleFunc("/objects_getcontract/{id:[0-9]+}", api.HandleGetObjectContract).Methods("GET", "OPTIONS")
	r.HandleFunc("/objects_mff/{hid:[0-9]+}", api.HandleGetObjectMff).Methods("GET", "OPTIONS")
	r.HandleFunc("/objcontracts", api.HandleObjContracts).Methods("GET", "OPTIONS")
	r.HandleFunc("/objcontracts_add", api.HandleAddObjContract).Methods("POST", "OPTIONS")
	r.HandleFunc("/objcontracts_upd", api.HandleUpdObjContract).Methods("POST", "OPTIONS")
	r.HandleFunc("/objcontracts_del", api.HandleDelObjContract).Methods("POST", "OPTIONS")
	r.HandleFunc("/objcontracts/{id:[0-9]+}", api.HandleGetObjContract).Methods("GET", "OPTIONS")
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
	r.HandleFunc("/pu", api.HandlePu).Methods("GET", "OPTIONS")
	r.HandleFunc("/pu_add", api.HandleAddPu).Methods("POST", "OPTIONS")
	// r.HandleFunc("/pu_obj_add", api.HandleAddPu).Methods("POST", "OPTIONS")
	r.HandleFunc("/pu_upd", api.HandleUpdPu).Methods("POST", "OPTIONS")
	r.HandleFunc("/pu_del", api.HandleDelPu).Methods("POST", "OPTIONS")
	r.HandleFunc("/pu/{id:[0-9]+}", api.HandleGetPu).Methods("GET", "OPTIONS")
	r.HandleFunc("/pu_obj", api.HandlePuObj).Methods("GET", "OPTIONS")
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
	r.HandleFunc("/subpu", api.HandleSubPu).Methods("GET", "OPTIONS")
	r.HandleFunc("/subpu_add", api.HandleAddSubPu).Methods("POST", "OPTIONS")
	r.HandleFunc("/subpu_upd", api.HandleUpdSubPu).Methods("POST", "OPTIONS")
	r.HandleFunc("/subpu_del", api.HandleDelSubPu).Methods("POST", "OPTIONS")
	r.HandleFunc("/subpu/{id:[0-9]+}", api.HandleGetSubPu).Methods("GET", "OPTIONS")
	r.HandleFunc("/subpu_prl", api.HandlePrlSubPu).Methods("GET", "OPTIONS")
	r.HandleFunc("/staff", api.HandleStaff).Methods("GET", "OPTIONS")
	r.HandleFunc("/staff_add", api.HandleAddStaff).Methods("POST", "OPTIONS")
	r.HandleFunc("/staff_upd", api.HandleUpdStaff).Methods("POST", "OPTIONS")
	r.HandleFunc("/staff_del", api.HandleDelStaff).Methods("POST", "OPTIONS")
	r.HandleFunc("/staff/{id:[0-9]+}", api.HandleGetStaff).Methods("GET", "OPTIONS")
	r.HandleFunc("/tgu", api.HandleTgu).Methods("GET", "OPTIONS")
	r.HandleFunc("/tgu_add", api.HandleAddTgu).Methods("POST", "OPTIONS")
	r.HandleFunc("/tgu_upd", api.HandleUpdTgu).Methods("POST", "OPTIONS")
	r.HandleFunc("/tgu_del", api.HandleDelTgu).Methods("POST", "OPTIONS")
	r.HandleFunc("/tgu/{id:[0-9]+}", api.HandleGetTgu).Methods("GET", "OPTIONS")
	r.HandleFunc("/conclusions", api.HandleConclusions).Methods("GET", "OPTIONS")
	r.HandleFunc("/conclusions_add", api.HandleAddConclusion).Methods("POST", "OPTIONS")
	r.HandleFunc("/conclusions_upd", api.HandleUpdConclusion).Methods("POST", "OPTIONS")
	r.HandleFunc("/conclusions_del", api.HandleDelConclusion).Methods("POST", "OPTIONS")
	r.HandleFunc("/conclusions/{id:[0-9]+}", api.HandleGetConclusion).Methods("GET", "OPTIONS")
	r.HandleFunc("/shutdowntypes", api.HandleShutdownTypes).Methods("GET", "OPTIONS")
	r.HandleFunc("/shutdowntypes_add", api.HandleAddShutdownType).Methods("POST", "OPTIONS")
	r.HandleFunc("/shutdowntypes_upd", api.HandleUpdShutdownType).Methods("POST", "OPTIONS")
	r.HandleFunc("/shutdowntypes_del", api.HandleDelShutdownType).Methods("POST", "OPTIONS")
	r.HandleFunc("/shutdowntypes/{id:[0-9]+}", api.HandleGetShutdownType).Methods("GET", "OPTIONS")
	r.HandleFunc("/violations", apg.HandleViolations).Methods("GET", "OPTIONS")
	r.HandleFunc("/violations_add", apg.HandleAddViolation).Methods("POST", "OPTIONS")
	r.HandleFunc("/violations_upd", apg.HandleUpdViolation).Methods("POST", "OPTIONS")
	r.HandleFunc("/violations_del", apg.HandleDelViolation).Methods("POST", "OPTIONS")
	r.HandleFunc("/violations/{id:[0-9]+}", apg.HandleGetViolation).Methods("GET", "OPTIONS")
	r.HandleFunc("/reasons", api.HandleReasons).Methods("GET", "OPTIONS")
	r.HandleFunc("/reasons_add", api.HandleAddReason).Methods("POST", "OPTIONS")
	r.HandleFunc("/reasons_upd", api.HandleUpdReason).Methods("POST", "OPTIONS")
	r.HandleFunc("/reasons_del", api.HandleDelReason).Methods("POST", "OPTIONS")
	r.HandleFunc("/reasons/{id:[0-9]+}", api.HandleGetReason).Methods("GET", "OPTIONS")
	r.HandleFunc("/seals", api.HandleSeals).Methods("GET", "OPTIONS")
	r.HandleFunc("/seals_add", api.HandleAddSeal).Methods("POST", "OPTIONS")
	r.HandleFunc("/seals_upd", api.HandleUpdSeal).Methods("POST", "OPTIONS")
	r.HandleFunc("/seals_del", api.HandleDelSeal).Methods("POST", "OPTIONS")
	r.HandleFunc("/seals/{id:[0-9]+}", api.HandleGetSeal).Methods("GET", "OPTIONS")
	r.HandleFunc("/sealtypes", api.HandleSealTypes).Methods("GET", "OPTIONS")
	r.HandleFunc("/sealtypes_add", api.HandleAddSealType).Methods("POST", "OPTIONS")
	r.HandleFunc("/sealtypes_upd", api.HandleUpdSealType).Methods("POST", "OPTIONS")
	r.HandleFunc("/sealtypes_del", api.HandleDelSealType).Methods("POST", "OPTIONS")
	r.HandleFunc("/sealtypes/{id:[0-9]+}", api.HandleGetSealType).Methods("GET", "OPTIONS")
	r.HandleFunc("/sealcolours", api.HandleSealColours).Methods("GET", "OPTIONS")
	r.HandleFunc("/sealcolours_add", api.HandleAddSealColour).Methods("POST", "OPTIONS")
	r.HandleFunc("/sealcolours_upd", api.HandleUpdSealColour).Methods("POST", "OPTIONS")
	r.HandleFunc("/sealcolours_del", api.HandleDelSealColour).Methods("POST", "OPTIONS")
	r.HandleFunc("/sealcolours/{id:[0-9]+}", api.HandleGetSealColour).Methods("GET", "OPTIONS")
	r.HandleFunc("/sealstatuses", api.HandleSealStatuses).Methods("GET", "OPTIONS")
	r.HandleFunc("/sealstatuses_add", api.HandleAddSealStatus).Methods("POST", "OPTIONS")
	r.HandleFunc("/sealstatuses_upd", api.HandleUpdSealStatus).Methods("POST", "OPTIONS")
	r.HandleFunc("/sealstatuses_del", api.HandleDelSealStatus).Methods("POST", "OPTIONS")
	r.HandleFunc("/sealstatuses/{id:[0-9]+}", api.HandleGetSealStatus).Methods("GET", "OPTIONS")
	r.HandleFunc("/charges", api.HandleCharges).Methods("GET", "OPTIONS")
	r.HandleFunc("/charges_add", api.HandleAddCharge).Methods("POST", "OPTIONS")
	r.HandleFunc("/charges_upd", api.HandleUpdCharge).Methods("POST", "OPTIONS")
	r.HandleFunc("/charges_del", api.HandleDelCharge).Methods("POST", "OPTIONS")
	r.HandleFunc("/charges/{id:[0-9]+}", api.HandleGetCharge).Methods("GET", "OPTIONS")
	r.HandleFunc("/charges_run/{id:[0-9]+}", apg.HandleChargeRun).Methods("GET", "OPTIONS")
	r.HandleFunc("/payments", api.HandlePayments).Methods("GET", "OPTIONS")
	r.HandleFunc("/payments_add", api.HandleAddPayment).Methods("POST", "OPTIONS")
	r.HandleFunc("/payments_upd", api.HandleUpdPayment).Methods("POST", "OPTIONS")
	r.HandleFunc("/payments_del", api.HandleDelPayment).Methods("POST", "OPTIONS")
	r.HandleFunc("/payments/{id:[0-9]+}", api.HandleGetPayment).Methods("GET", "OPTIONS")
	r.HandleFunc("/equipmenttypes", api.HandleEquipmentTypes).Methods("GET", "OPTIONS")
	r.HandleFunc("/equipmenttypes_add", api.HandleAddEquipmentType).Methods("POST", "OPTIONS")
	r.HandleFunc("/equipmenttypes_upd", api.HandleUpdEquipmentType).Methods("POST", "OPTIONS")
	r.HandleFunc("/equipmenttypes_del", api.HandleDelEquipmentType).Methods("POST", "OPTIONS")
	r.HandleFunc("/equipmenttypes/{id:[0-9]+}", api.HandleGetEquipmentType).Methods("GET", "OPTIONS")
	r.HandleFunc("/equipment", api.HandleEquipment).Methods("GET", "OPTIONS")
	r.HandleFunc("/equipment_add", api.HandleAddEquipment).Methods("POST", "OPTIONS")
	r.HandleFunc("/equipment_addlist", api.HandleAddEquipmentList).Methods("POST", "OPTIONS")
	r.HandleFunc("/equipment_upd", api.HandleUpdEquipment).Methods("POST", "OPTIONS")
	r.HandleFunc("/equipment_del", api.HandleDelEquipment).Methods("POST", "OPTIONS")
	r.HandleFunc("/equipment_delobj", api.HandleDelObjEquipment).Methods("POST", "OPTIONS")
	r.HandleFunc("/equipment/{id:[0-9]+}", api.HandleGetEquipment).Methods("GET", "OPTIONS")
	r.HandleFunc("/calculationtypes", api.HandleCalculationTypes).Methods("GET", "OPTIONS")
	r.HandleFunc("/calculationtypes_add", api.HandleAddCalculationType).Methods("POST", "OPTIONS")
	r.HandleFunc("/calculationtypes_upd", api.HandleUpdCalculationType).Methods("POST", "OPTIONS")
	r.HandleFunc("/calculationtypes_del", api.HandleDelCalculationType).Methods("POST", "OPTIONS")
	r.HandleFunc("/calculationtypes/{id:[0-9]+}", api.HandleGetCalculationType).Methods("GET", "OPTIONS")
	r.HandleFunc("/results", api.HandleResults).Methods("GET", "OPTIONS")
	r.HandleFunc("/results_add", api.HandleAddResult).Methods("POST", "OPTIONS")
	r.HandleFunc("/results_upd", api.HandleUpdResult).Methods("POST", "OPTIONS")
	r.HandleFunc("/results_del", api.HandleDelResult).Methods("POST", "OPTIONS")
	r.HandleFunc("/results/{id:[0-9]+}", api.HandleGetResult).Methods("GET", "OPTIONS")
	r.HandleFunc("/servicetypes", api.HandleServiceTypes).Methods("GET", "OPTIONS")
	r.HandleFunc("/servicetypes_add", api.HandleAddServiceType).Methods("POST", "OPTIONS")
	r.HandleFunc("/servicetypes_upd", api.HandleUpdServiceType).Methods("POST", "OPTIONS")
	r.HandleFunc("/servicetypes_del", api.HandleDelServiceType).Methods("POST", "OPTIONS")
	r.HandleFunc("/servicetypes/{id:[0-9]+}", api.HandleGetServiceType).Methods("GET", "OPTIONS")
	r.HandleFunc("/requesttypes", api.HandleRequestTypes).Methods("GET", "OPTIONS")
	r.HandleFunc("/requesttypes_add", api.HandleAddRequestType).Methods("POST", "OPTIONS")
	r.HandleFunc("/requesttypes_upd", api.HandleUpdRequestType).Methods("POST", "OPTIONS")
	r.HandleFunc("/requesttypes_del", api.HandleDelRequestType).Methods("POST", "OPTIONS")
	r.HandleFunc("/requesttypes/{id:[0-9]+}", api.HandleGetRequestType).Methods("GET", "OPTIONS")
	r.HandleFunc("/claimtypes", api.HandleClaimTypes).Methods("GET", "OPTIONS")
	r.HandleFunc("/claimtypes_add", api.HandleAddClaimType).Methods("POST", "OPTIONS")
	r.HandleFunc("/claimtypes_upd", api.HandleUpdClaimType).Methods("POST", "OPTIONS")
	r.HandleFunc("/claimtypes_del", api.HandleDelClaimType).Methods("POST", "OPTIONS")
	r.HandleFunc("/claimtypes/{id:[0-9]+}", api.HandleGetClaimType).Methods("GET", "OPTIONS")
	r.HandleFunc("/requests", api.HandleRequests).Methods("GET", "OPTIONS")
	r.HandleFunc("/requests_add", api.HandleAddRequest).Methods("POST", "OPTIONS")
	r.HandleFunc("/requests_upd", api.HandleUpdRequest).Methods("POST", "OPTIONS")
	r.HandleFunc("/requests_del", api.HandleDelRequest).Methods("POST", "OPTIONS")
	r.HandleFunc("/requests/{id:[0-9]+}", api.HandleGetRequest).Methods("GET", "OPTIONS")
	r.HandleFunc("/objstatuses", api.HandleObjStatuses).Methods("GET", "OPTIONS")
	r.HandleFunc("/objstatuses_add", api.HandleAddObjStatus).Methods("POST", "OPTIONS")
	r.HandleFunc("/objstatuses_upd", api.HandleUpdObjStatus).Methods("POST", "OPTIONS")
	r.HandleFunc("/objstatuses_del", api.HandleDelObjStatus).Methods("POST", "OPTIONS")
	r.HandleFunc("/objstatuses/{id:[0-9]+}", api.HandleGetObjStatus).Methods("GET", "OPTIONS")
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
	r.HandleFunc("/requestkinds", api.HandleRequestKinds).Methods("GET", "OPTIONS")
	r.HandleFunc("/requestkinds_add", api.HandleAddRequestKind).Methods("POST", "OPTIONS")
	r.HandleFunc("/requestkinds_upd", api.HandleUpdRequestKind).Methods("POST", "OPTIONS")
	r.HandleFunc("/requestkinds_del", api.HandleDelRequestKind).Methods("POST", "OPTIONS")
	r.HandleFunc("/requestkinds/{id:[0-9]+}", api.HandleGetRequestKind).Methods("GET", "OPTIONS")
	r.HandleFunc("/periods", api.HandlePeriods).Methods("GET", "OPTIONS")
	r.HandleFunc("/periods_add", api.HandleAddPeriod).Methods("POST", "OPTIONS")
	r.HandleFunc("/periods_upd", api.HandleUpdPeriod).Methods("POST", "OPTIONS")
	r.HandleFunc("/periods_del", api.HandleDelPeriod).Methods("POST", "OPTIONS")
	r.HandleFunc("/periods/{id:[0-9]+}", api.HandleGetPeriod).Methods("GET", "OPTIONS")
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
	r.HandleFunc("/objtranscurr", api.HandleObjTransCurr).Methods("GET", "OPTIONS")
	r.HandleFunc("/objtranscurr_add", api.HandleAddObjTransCurr).Methods("POST", "OPTIONS")
	r.HandleFunc("/objtranscurr_upd", api.HandleUpdObjTransCurr).Methods("POST", "OPTIONS")
	r.HandleFunc("/objtranscurr_del", api.HandleDelObjTransCurr).Methods("POST", "OPTIONS")
	r.HandleFunc("/objtranscurr/{id:[0-9]+}", api.HandleGetObjTransCurr).Methods("GET", "OPTIONS")
	r.HandleFunc("/objtranscurr_obj", api.HandleObjTransCurrByObj).Methods("GET", "OPTIONS")
	r.HandleFunc("/objtransvolt", api.HandleObjTransVolt).Methods("GET", "OPTIONS")
	r.HandleFunc("/objtransvolt_add", api.HandleAddObjTransVolt).Methods("POST", "OPTIONS")
	r.HandleFunc("/objtransvolt_upd", api.HandleUpdObjTransVolt).Methods("POST", "OPTIONS")
	r.HandleFunc("/objtransvolt_del", api.HandleDelObjTransVolt).Methods("POST", "OPTIONS")
	r.HandleFunc("/objtransvolt/{id:[0-9]+}", api.HandleGetObjTransVolt).Methods("GET", "OPTIONS")
	r.HandleFunc("/objtransvolt_obj", api.HandleObjTransVoltByObj).Methods("GET", "OPTIONS")
	r.HandleFunc("/objtranspwr", api.HandleObjTransPwr).Methods("GET", "OPTIONS")
	r.HandleFunc("/objtranspwr_add", api.HandleAddObjTransPwr).Methods("POST", "OPTIONS")
	r.HandleFunc("/objtranspwr_upd", api.HandleUpdObjTransPwr).Methods("POST", "OPTIONS")
	r.HandleFunc("/objtranspwr_del", api.HandleDelObjTransPwr).Methods("POST", "OPTIONS")
	r.HandleFunc("/objtranspwr/{id:[0-9]+}", api.HandleGetObjTransPwr).Methods("GET", "OPTIONS")
	r.HandleFunc("/objtranspwr_obj", api.HandleObjTransPwrByObj).Methods("GET", "OPTIONS")
	r.HandleFunc("/cableresistances", api.HandleCableResistances).Methods("GET", "OPTIONS")
	r.HandleFunc("/cableresistances_add", api.HandleAddCableResistance).Methods("POST", "OPTIONS")
	r.HandleFunc("/cableresistances_upd", api.HandleUpdCableResistance).Methods("POST", "OPTIONS")
	r.HandleFunc("/cableresistances_del", api.HandleDelCableResistance).Methods("POST", "OPTIONS")
	r.HandleFunc("/cableresistances/{id:[0-9]+}", api.HandleGetCableResistance).Methods("GET", "OPTIONS")
	r.HandleFunc("/objlines", api.HandleObjLines).Methods("GET", "OPTIONS")
	r.HandleFunc("/objlines_add", api.HandleAddObjLine).Methods("POST", "OPTIONS")
	r.HandleFunc("/objlines_upd", api.HandleUpdObjLine).Methods("POST", "OPTIONS")
	r.HandleFunc("/objlines_del", api.HandleDelObjLine).Methods("POST", "OPTIONS")
	r.HandleFunc("/objlines/{id:[0-9]+}", api.HandleGetObjLine).Methods("GET", "OPTIONS")
	r.HandleFunc("/objlines_obj", api.HandleObjLinesByObj).Methods("GET", "OPTIONS")
	r.HandleFunc("/askuetypes", api.HandleAskueTypes).Methods("GET", "OPTIONS")
	r.HandleFunc("/askuetypes_add", api.HandleAddAskueType).Methods("POST", "OPTIONS")
	r.HandleFunc("/askuetypes_upd", api.HandleUpdAskueType).Methods("POST", "OPTIONS")
	r.HandleFunc("/askuetypes_del", api.HandleDelAskueType).Methods("POST", "OPTIONS")
	r.HandleFunc("/askuetypes/{id:[0-9]+}", api.HandleGetAskueType).Methods("GET", "OPTIONS")
	r.HandleFunc("/puvaluetypes", api.HandlePuValueTypes).Methods("GET", "OPTIONS")
	r.HandleFunc("/puvaluetypes_add", api.HandleAddPuValueType).Methods("POST", "OPTIONS")
	r.HandleFunc("/puvaluetypes_upd", api.HandleUpdPuValueType).Methods("POST", "OPTIONS")
	r.HandleFunc("/puvaluetypes_del", api.HandleDelPuValueType).Methods("POST", "OPTIONS")
	r.HandleFunc("/puvaluetypes/{id:[0-9]+}", api.HandleGetPuValueType).Methods("GET", "OPTIONS")
	r.HandleFunc("/ordertypes", api.HandleOrderTypes).Methods("GET", "OPTIONS")
	r.HandleFunc("/ordertypes_add", api.HandleAddOrderType).Methods("POST", "OPTIONS")
	r.HandleFunc("/ordertypes_upd", api.HandleUpdOrderType).Methods("POST", "OPTIONS")
	r.HandleFunc("/ordertypes_del", api.HandleDelOrderType).Methods("POST", "OPTIONS")
	r.HandleFunc("/ordertypes/{id:[0-9]+}", api.HandleGetOrderType).Methods("GET", "OPTIONS")
	r.HandleFunc("/distributionzones", api.HandleDistributionZones).Methods("GET", "OPTIONS")
	r.HandleFunc("/distributionzones_add", api.HandleAddDistributionZone).Methods("POST", "OPTIONS")
	r.HandleFunc("/distributionzones_upd", api.HandleUpdDistributionZone).Methods("POST", "OPTIONS")
	r.HandleFunc("/distributionzones_del", api.HandleDelDistributionZone).Methods("POST", "OPTIONS")
	r.HandleFunc("/distributionzones/{id:[0-9]+}", api.HandleGetDistributionZone).Methods("GET", "OPTIONS")
	r.HandleFunc("/contracttypes", api.HandleContractTypes).Methods("GET", "OPTIONS")
	r.HandleFunc("/contracttypes_add", api.HandleAddContractType).Methods("POST", "OPTIONS")
	r.HandleFunc("/contracttypes_upd", api.HandleUpdContractType).Methods("POST", "OPTIONS")
	r.HandleFunc("/contracttypes_del", api.HandleDelContractType).Methods("POST", "OPTIONS")
	r.HandleFunc("/contracttypes/{id:[0-9]+}", api.HandleGetContractType).Methods("GET", "OPTIONS")
}
